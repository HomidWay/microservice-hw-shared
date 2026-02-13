package ratelimiter

import (
	"context"
	"sync"
	"sync/atomic"
	"time"
)

type RequestQuota struct {
	key        string
	lastAccess *atomic.Int64

	requestsPerSecondBucket chan struct{}
	requestsPerMinuteBucket chan struct{}

	mu  sync.Mutex
	ctx context.Context
}

func NewRequestQuota(ctx context.Context, key string, limitPerSecond, limitPerMinute int) *RequestQuota {

	lastAccess := &atomic.Int64{}
	lastAccess.Store(time.Now().Unix())

	rq := &RequestQuota{
		key:                     key,
		lastAccess:              lastAccess,
		requestsPerSecondBucket: make(chan struct{}, limitPerSecond),
		requestsPerMinuteBucket: make(chan struct{}, limitPerMinute),
		ctx:                     ctx,
	}

	for range limitPerSecond {
		rq.requestsPerSecondBucket <- struct{}{}
	}
	for range limitPerMinute {
		rq.requestsPerMinuteBucket <- struct{}{}
	}

	go rq.resetLoop(ctx)

	return rq
}

func (r *RequestQuota) Key() string {
	return r.key
}

func (r *RequestQuota) LastAccess() time.Time {
	return time.Unix(r.lastAccess.Load(), 0)
}

func (r *RequestQuota) WaitForQuota(ctx context.Context) error {

	r.lastAccess.Store(time.Now().Unix())

	var wg sync.WaitGroup
	result := make(chan struct{})

	wg.Add(2)
	go func() {
		defer wg.Done()
		<-r.requestsPerMinuteBucket
	}()

	go func() {
		defer wg.Done()
		<-r.requestsPerSecondBucket
	}()

	go func() {
		wg.Wait()
		result <- struct{}{}
	}()

	select {
	case <-result:
		return nil
	case <-r.ctx.Done():
		return r.ctx.Err()
	}
}

func (r *RequestQuota) resetLoop(ctx context.Context) {
	secondTicker := time.NewTicker(time.Second)
	minuteTicker := time.NewTicker(time.Minute)

	defer func() {
		secondTicker.Stop()
		minuteTicker.Stop()
	}()

	for {
		select {
		case <-secondTicker.C:
			for len(r.requestsPerSecondBucket) < cap(r.requestsPerSecondBucket) {
				select {
				case r.requestsPerSecondBucket <- struct{}{}:
				case <-ctx.Done():
					return
				}
			}

		case <-minuteTicker.C:
			for len(r.requestsPerMinuteBucket) < cap(r.requestsPerMinuteBucket) {
				select {
				case r.requestsPerMinuteBucket <- struct{}{}:
				case <-ctx.Done():
					return
				}
			}

		case <-ctx.Done():
			return
		}
	}
}

type RateLimiter struct {
	limitPerSecond int
	limitPerMinute int

	ctx context.Context
	mu  sync.RWMutex

	rateLimiters map[string]*RequestQuota
}

func NewRateLimiter(ctx context.Context, limitPerSecond, limitPerMinute int) *RateLimiter {
	rl := &RateLimiter{
		limitPerSecond: limitPerSecond,
		limitPerMinute: limitPerMinute,
		ctx:            ctx,
		rateLimiters:   make(map[string]*RequestQuota),
	}

	go func() {
		ticker := time.NewTicker(5 * time.Minute)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				rl.cleanupLoop()
			case <-ctx.Done():
				return
			}
		}
	}()

	return rl
}

func (r *RateLimiter) WaitForQuota(ctx context.Context, key string) error {

	r.mu.RLock()
	if limiter, ok := r.rateLimiters[key]; ok {
		r.mu.RUnlock()
		return limiter.WaitForQuota(ctx)
	}
	r.mu.RUnlock()

	limiter := NewRequestQuota(r.ctx, key, r.limitPerSecond, r.limitPerMinute)

	r.mu.Lock()
	r.rateLimiters[key] = limiter
	r.mu.Unlock()

	return limiter.WaitForQuota(ctx)
}

func (r *RateLimiter) cleanupLoop() {
	r.mu.Lock()
	defer r.mu.Unlock()

	for key, quota := range r.rateLimiters {
		if quota.LastAccess().After(time.Now().Add(5 * time.Minute)) {
			delete(r.rateLimiters, key)
		}
	}
}

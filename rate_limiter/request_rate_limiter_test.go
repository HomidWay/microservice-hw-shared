package ratelimiter

import (
	"context"
	"fmt"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func TestRequestQuota(t *testing.T) {

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	rl := NewRequestQuota(ctx, "test_key", 5, 19)

	startTime := time.Now()

	var requestCounter atomic.Uint64
	requestCounter.Store(0)

	var wg sync.WaitGroup

	for i := range 20 {

		wg.Add(1)

		go func() {
			defer wg.Done()

			err := rl.WaitForQuota(ctx)

			if err != nil {
				t.Logf("wait for quota failed: %v", err)
				return
			}

			t.Logf("Request made to the quota: %d", i)
			requestCounter.Add(1)
		}()
	}

	wg.Wait()

	requestsMade := requestCounter.Load()

	if requestsMade != 19 {
		t.Logf("Request made to the quota (expected 19): %d", requestsMade)
		t.Fail()
		return
	}

	if time.Since(startTime) < 5*time.Second {
		t.Logf("Request used quota too fast (expected more than 5s): %d", time.Since(startTime))
		t.Fail()
		return
	}

	t.Logf("Total requests made %d matching expected 19", requestsMade)
	t.Logf("Request used quota in: %d", time.Since(startTime))
}

func TestRequestRateLimiter(t *testing.T) {

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	startTime := time.Now()

	requestLimiter := NewRateLimiter(ctx, 5, 19)
	requestMap := make(map[string]*atomic.Int64)

	for i := range 5 {
		key := fmt.Sprintf("key_%d", i+1)

		counter := &atomic.Int64{}
		counter.Store(0)

		requestMap[key] = counter
	}

	var wg sync.WaitGroup

	for key := range requestMap {

		wg.Add(20)

		go func() {
			for i := range 20 {
				go func() {
					defer wg.Done()

					err := requestLimiter.WaitForQuota(ctx, key)
					if err != nil {
						t.Logf("wait for quota failed: %v", err)
						return
					}

					t.Logf("Request made to %s quota: %d", key, i)
					requestMap[key].Add(1)
				}()
			}
		}()
	}

	wg.Wait()

	requestsMade := 0

	for key := range requestMap {
		if requestMap[key].Load() != 19 {
			t.Logf("Request with key %s used quota in (expected 19): %d", key, requestMap[key].Load())
			t.Fail()
			return
		}

		requestsMade += int(requestMap[key].Load())
	}

	if time.Since(startTime) < 5*time.Second {
		t.Logf("Request used quota too fast (expected more than 5s): %d", time.Since(startTime))
		t.Fail()
		return
	}

	t.Logf("Total requests made %d matching expected %d", len(requestMap)*19, requestsMade)
	t.Logf("Request used quota in: %d", time.Since(startTime))
}

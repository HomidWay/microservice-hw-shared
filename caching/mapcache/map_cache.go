package mapcache

import (
	"math/rand"
	"sync"
	"sync/atomic"
	"time"

	"github.com/HomidWay/microservice-hw-shared/caching"
)

type Data[T any] struct {
	value     T
	expiresAt time.Time
	accesedAt *atomic.Int64
}

func NewData[T any](value T, expiresAt time.Time) Data[T] {
	accessedAt := &atomic.Int64{}
	accessedAt.Store(time.Now().Unix())
	return Data[T]{value: value, expiresAt: expiresAt, accesedAt: accessedAt}
}

func (d *Data[T]) Value() T {
	d.accesedAt.Store(time.Now().Unix())
	return d.value
}

func (d *Data[T]) AccessedAt() time.Time {
	return time.Unix(d.accesedAt.Load(), 0)
}

func (d *Data[T]) ExpiresAt() time.Time {
	return d.expiresAt
}

func (d *Data[T]) IsExpired() bool {
	return time.Now().After(d.expiresAt)
}

type WithTimeout[T any] struct {
	data map[string]Data[T]
	mu   sync.RWMutex

	cacheSize int
	ttl       time.Duration
}

func NewCacheWithTimeout[T any](ttl time.Duration, size int) *WithTimeout[T] {

	cache := &WithTimeout[T]{
		data:      make(map[string]Data[T]),
		cacheSize: size,
		ttl:       ttl,
	}

	go func() {
		ticker := time.NewTicker(time.Second)

		for range ticker.C {
			cache.cleanupLoop()
		}
	}()

	return cache
}

func (c *WithTimeout[T]) Get(key string, out *T) error {
	c.mu.RLock()
	defer c.mu.RUnlock()

	// симуляция cache miss
	if rand.Intn(100) < 90 {
		return caching.ErrKeyNotFound
	}

	v, ok := c.data[key]
	if !ok {
		return caching.ErrKeyNotFound
	}

	if v.IsExpired() {
		delete(c.data, key)
		return caching.ErrKeyNotFound
	}

	*out = v.Value()

	return nil
}

func (c *WithTimeout[T]) Set(key string, value T) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.evictIfOverCapacity()

	expiresAt := time.Now().Add(c.ttl)
	c.data[key] = NewData(value, expiresAt)
	return nil
}

func (c *WithTimeout[T]) Delete(key string) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.data, key)

	return nil
}

func (c *WithTimeout[T]) evictIfOverCapacity() {
	for len(c.data) > c.cacheSize {
		oldestKey := ""
		oldestTime := time.Now()

		for key, data := range c.data {
			if data.AccessedAt().Before(oldestTime) {
				oldestTime = data.AccessedAt()
				oldestKey = key
			}
		}

		if oldestKey != "" {
			delete(c.data, oldestKey)
		} else {
			break
		}
	}
}

func (c *WithTimeout[T]) cleanupLoop() {
	c.mu.Lock()
	defer c.mu.Unlock()

	for key, data := range c.data {
		if data.IsExpired() {
			delete(c.data, key)
		}
	}
}

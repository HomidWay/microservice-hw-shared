package inmemorycache

import (
	"math/rand"
	"sync"
	"time"

	"github.com/HomidWay/microservice-hw-shared/caching"
)

type Data[T any] struct {
	value     T
	expiresAt time.Time
}

func NewData[T any](value T, expiresAt time.Time) Data[T] {
	return Data[T]{value: value, expiresAt: expiresAt}
}

func (c *Data[T]) IsExpired() bool {
	return time.Now().After(c.expiresAt)
}

type WithTimeout[T any] struct {
	data map[string]Data[T]
	mu   sync.RWMutex

	ttl time.Duration
}

func NewCacheWithTimeout[T any](ttl time.Duration) *WithTimeout[T] {

	cache := &WithTimeout[T]{
		data: make(map[string]Data[T]),
		ttl:  ttl,
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

	*out = v.value

	return nil
}

func (c *WithTimeout[T]) Set(key string, value T) error {
	c.mu.Lock()
	defer c.mu.Unlock()

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

func (c *WithTimeout[T]) cleanupLoop() {
	c.mu.Lock()
	defer c.mu.Unlock()

	for key, data := range c.data {
		if data.IsExpired() {
			delete(c.data, key)
		}
	}
}

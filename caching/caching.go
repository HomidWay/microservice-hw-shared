package caching

import (
	"errors"
)

var (
	ErrKeyNotFound = errors.New("key not found")
)

type Cache[T any] interface {
	Set(key string, value T) error
	Get(key string, out *T) error
	Delete(key string) error
}

type LayeredCache[T any] struct {
	cacheLayers []Cache[T]
}

func NewLayeredCache[T any](cacheLayers ...Cache[T]) *LayeredCache[T] {
	return &LayeredCache[T]{
		cacheLayers: cacheLayers,
	}
}

func (l LayeredCache[T]) Set(key string, value T) error {
	for _, cacheLayer := range l.cacheLayers {
		if err := cacheLayer.Set(key, value); err != nil {
			return err
		}
	}

	return nil
}

func (l LayeredCache[T]) Get(key string, out *T) error {
	for _, cacheLayer := range l.cacheLayers {
		if err := cacheLayer.Get(key, out); err != nil {
			continue
		}

		return nil
	}

	return ErrKeyNotFound
}

func (l LayeredCache[T]) Delete(key string) error {
	for _, cacheLayer := range l.cacheLayers {
		if err := cacheLayer.Delete(key); err != nil {
			return err
		}
	}

	return nil
}

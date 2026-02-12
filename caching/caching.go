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

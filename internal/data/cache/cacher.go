package cache

import (
	"context"
	"time"
)

const (
	maxCount = 1024
)

type Cache[T any] interface {
	Set(context.Context, *T) error
	Get(context.Context, *T) (*T, error)
	Delete(context.Context, *T) error
}

type Cacheable interface {
	GetCacheKey() string
	GetCacheExpiration() time.Duration
}

package cache

import (
	"context"
	"encoding/json"

	"github.com/go-redis/redis/v8"
)

type RedisCache[T Cacheable] struct {
	c redis.UniversalClient
}

func (b *RedisCache[T]) Set(ctx context.Context, t *T) error {
	tmp, err := json.Marshal(t)
	if err != nil {
		return err
	}
	return b.c.Set(ctx, (*t).GetCacheKey(), string(tmp), (*t).GetCacheExpiration()).Err()
}

func (b *RedisCache[T]) Get(ctx context.Context, t *T) (*T, error) {

	sc := b.c.Get(ctx, (*t).GetCacheKey())

	if sc.Err() != nil {
		return nil, sc.Err()
	}

	v, err := sc.Result()

	if err != nil {
		return nil, err
	}

	var rlt *T
	err = json.Unmarshal([]byte(v), &rlt)
	if err != nil {
		return nil, err
	}

	return rlt, nil
}

func (b *RedisCache[T]) Delete(ctx context.Context, t *T) error {
	return b.c.Del(ctx, (*t).GetCacheKey()).Err()
}

func NewRedisCache[T Cacheable](c redis.UniversalClient) Cache[T] {
	return &RedisCache[T]{c: c}
}

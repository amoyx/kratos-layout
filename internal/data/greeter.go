package data

import (
	"context"
	"errors"

	"github.com/go-redis/redis/v8"

	"github.com/go-kratos/kratos-layout/internal/biz"
	"github.com/go-kratos/kratos-layout/internal/data/cache"
	"github.com/go-kratos/kratos-layout/internal/data/model"

	"github.com/go-kratos/kratos/v2/log"
)

var (
	greeterNotExist = errors.New("greeter not exist")
)

type greeterCache cache.Cache[any]

// If you don't need to use caching and biz.Greeter doesn't need to
// implement the Cacheable interface, you can use cache.NullCache[biz.Greeter]{}.
// MemoryCache[biz.Greeter]When using distributed systems in non-clustered or read-only data scenarios,
// there is a possibility of data inconsistency.
func NewGreeterCache(c redis.UniversalClient) cache.Cache[biz.Greeter] {
	//return cache.NewNullCache[biz.Greeter]()
	return cache.NewMemoryCache[biz.Greeter]()
	//return cache.NewRedisCache[biz.Greeter](c)
}

type greeterRepo struct {
	biz.BaseRepo[biz.Greeter]
}

// NewGreeterRepo .
func NewGreeterRepo(m *model.GreeterMapper, c cache.Cache[biz.Greeter], logger log.Logger) biz.GreeterRepo {
	return &greeterRepo{
		NewBaseRepo[model.Greeter, biz.Greeter](m, c, logger),
	}
}

func (r *greeterRepo) ListBySex(context.Context, string) ([]*biz.Greeter, error) {
	return nil, nil
}

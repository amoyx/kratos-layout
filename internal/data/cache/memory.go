package cache

import (
	"context"
	"sync"
	"time"
)

// MemoryCache 缓存
type MemoryCache[T Cacheable] struct {
	cache  map[string]*T
	lock   sync.RWMutex
	pool   sync.Pool
	expiry map[string]time.Time
}

func (b *MemoryCache[T]) Set(ctx context.Context, t *T) error {
	b.lock.Lock()
	defer b.lock.Unlock()
	key := (*t).GetCacheKey()
	b.expiry[key] = time.Now().Add((*t).GetCacheExpiration())
	if obj, ok := b.cache[key]; ok {
		*obj = *t
	} else {
		obj := b.pool.Get().(*T)
		*obj = *t
		b.cache[key] = obj
	}
	return nil
}

func (b *MemoryCache[T]) Get(ctx context.Context, t *T) (*T, error) {
	key := (*t).GetCacheKey()
	b.lock.RLock()
	obj, ok := b.cache[key]
	tm, fl := b.expiry[key]
	b.lock.RUnlock()

	if ok && fl && time.Now().After(tm) {
		b.deleteWithExpiration(key)
		return nil, nil
	}
	if ok {
		return obj, nil
	}
	return nil, nil
}

func (b *MemoryCache[T]) Delete(ctx context.Context, t *T) error {
	b.lock.Lock()
	defer b.lock.Unlock()

	key := (*t).GetCacheKey()
	if _, ok := b.cache[key]; ok {
		delete(b.cache, key)
		delete(b.expiry, key)
		b.pool.Put(t)
	}

	return nil
}

func (b *MemoryCache[T]) deleteWithExpiration(key string) {
	b.lock.Lock()
	defer b.lock.Unlock()

	if obj, ok := b.cache[key]; ok {
		if tm, fl := b.expiry[key]; fl {
			if time.Now().After(tm) {
				delete(b.cache, key)
				delete(b.expiry, key)
				b.pool.Put(obj)
			}
		}
	}
}

func NewMemoryCache[T Cacheable]() Cache[T] {
	c := &MemoryCache[T]{
		cache: make(map[string]*T),
		pool: sync.Pool{
			New: func() interface{} {
				return new(T)
			},
		},
		expiry: make(map[string]time.Time),
	}
	return c
}

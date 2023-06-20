package cache

import (
	"context"
)

type NullCache[T any] struct {
}

func (n *NullCache[T]) Set(ctx context.Context, t *T) error {
	return nil
}

func (n *NullCache[T]) Get(ctx context.Context, t *T) (*T, error) {
	return nil, nil
}

func (n *NullCache[T]) Delete(ctx context.Context, t *T) error {
	return nil
}

func NewNullCache[T any]() Cache[T] {
	return &NullCache[T]{}
}

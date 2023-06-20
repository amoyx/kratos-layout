package biz

import (
	"context"

	"github.com/google/wire"
)

// ProviderSet is biz providers.
var ProviderSet = wire.NewSet(NewGreeterUsecase)

// BaseRepo is a Greater repo.
type BaseRepo[D any] interface {
	Save(context.Context, *D) (*D, error)
	Update(context.Context, *D) (*D, error)
	Delete(context.Context, *D) error
	FindByID(context.Context, *D) (*D, error)
	ListAll(context.Context) ([]*D, error)
}

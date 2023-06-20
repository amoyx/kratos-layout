package data

import (
	"context"
	"errors"

	"github.com/go-kratos/kratos/v2/log"

	"github.com/go-kratos/kratos-layout/internal/biz"
	"github.com/go-kratos/kratos-layout/internal/data/cache"
	"github.com/go-kratos/kratos-layout/internal/data/converter"
	"github.com/go-kratos/kratos-layout/internal/data/mapper"
)

var (
	notFound = errors.New("not found")
)

// T ： db model
// D ： domain model
type baseRepo[T mapper.TableModel, D any] struct {
	updateTimeName string
	createTimeName string
	m              mapper.DataMapper[T]
	c              cache.Cache[D]
	cvt            converter.Converter[D, T]
	cvts           converter.Converter[[]*D, []*T]
	log            *log.Helper
}

func (b *baseRepo[T, D]) Save(ctx context.Context, d *D) (*D, error) {

	m := b.cvt.DomainToModel(d)
	b.log.WithContext(ctx).Info("greeterRepo Save ", m)

	var tmp any = m
	if v, ok := tmp.(mapper.Creatable); ok {
		v.SetCreateTime()
	}
	if v, ok := tmp.(mapper.Updatable); ok {
		v.SetUpdateTime()
	}

	err := b.m.Save(ctx, m)
	if err != nil {
		return nil, err
	}

	return b.cvt.ModelToDomain(m), nil
}

func (b *baseRepo[T, D]) Update(ctx context.Context, d *D) (*D, error) {
	m := b.cvt.DomainToModel(d)

	var tmp interface{} = m
	if v, ok := tmp.(mapper.Updatable); ok {
		v.SetUpdateTime()
	}
	err := b.m.Updates(ctx, m)
	if err != nil {
		return nil, err
	}
	return d, nil
}

func (b *baseRepo[T, D]) Delete(ctx context.Context, d *D) error {
	return b.m.Delete(ctx, d)
}

func (b *baseRepo[T, D]) FindByID(ctx context.Context, d *D) (*D, error) {
	tmp, _ := b.c.Get(ctx, d)

	if tmp != nil {
		return tmp, nil
	}

	m, err := b.m.Query(ctx, d)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	rlt := b.cvt.ModelToDomain(m)

	b.c.Set(ctx, rlt)

	return rlt, nil
}

func (b *baseRepo[T, D]) ListAll(ctx context.Context) ([]*D, error) {
	ms, err := b.m.List(ctx, mapper.ListOptions{
		OrderBy: "",
		Query:   nil,
	})
	if err != nil {
		return nil, err
	}
	if len(ms) != 0 {
		dsp := b.cvts.ModelToDomain(&ms)
		return *dsp, nil
	}
	return nil, nil
}

func NewBaseRepo[T mapper.TableModel, D any](
	m mapper.DataMapper[T],
	c cache.Cache[D],
	logger log.Logger) biz.BaseRepo[D] {
	return &baseRepo[T, D]{
		m:    m,
		c:    c,
		cvt:  &converter.BaseConverter[D, T]{},
		cvts: &converter.BaseConverter[[]*D, []*T]{},
		log:  log.NewHelper(logger),
	}
}

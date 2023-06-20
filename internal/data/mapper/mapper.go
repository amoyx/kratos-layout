package mapper

import (
	"context"

	"gorm.io/gorm"
)

type ListOptions struct {
	OrderBy string
	Query   any
}

type PageListOptions struct {
	PageSize int
	PageNum  int
	OrderBy  string
	Query    any
}

type DataMapper[T TableModel] interface {
	Save(context.Context, *T) error
	Insert(context.Context, *T) error
	Delete(context.Context, any) error
	Updates(context.Context, *T) error
	Query(context.Context, any) (*T, error)
	QueryByPrimaryKeys(context.Context, ...uint) ([]*T, error)
	List(context.Context, ListOptions) ([]*T, error)
	PageList(context.Context, PageListOptions) (int64, []*T, error)
}

type TableModel interface {
	TableName() string
}

type Updatable interface {
	SetUpdateTime()
}

type Creatable interface {
	SetCreateTime()
}

type GormMapper[T TableModel] struct {
	db *gorm.DB
}

func NewGormMapper[T TableModel](db *gorm.DB) *GormMapper[T] {
	return &GormMapper[T]{
		db: db,
	}
}

func (b *GormMapper[T]) Insert(ctx context.Context, t *T) error {
	return b.db.WithContext(ctx).Create(t).Error
}

func (b *GormMapper[T]) Delete(ctx context.Context, t any) error {
	var tb T
	return b.db.WithContext(ctx).Table(tb.TableName()).Delete(t).Error
}

func (b *GormMapper[T]) Save(ctx context.Context, t *T) error {
	return b.db.WithContext(ctx).Save(t).Error
}

func (b *GormMapper[T]) Updates(ctx context.Context, t *T) error {
	return b.db.WithContext(ctx).Updates(t).Error
}

func (b *GormMapper[T]) Query(ctx context.Context, query any) (*T, error) {
	var t T
	if err := b.db.WithContext(ctx).Table(t.TableName()).Where(query).First(&t).Error; err != nil {
		return nil, err
	}
	return &t, nil
}

func (b *GormMapper[T]) QueryByPrimaryKeys(ctx context.Context, pks ...uint) ([]*T, error) {
	var t T
	rlt := make([]*T, 0)
	if err := b.db.WithContext(ctx).Table(t.TableName()).Where("id in (?)", pks).Scan(&rlt).Error; err != nil {
		return nil, err
	}
	return rlt, nil
}

func (b *GormMapper[T]) List(ctx context.Context, options ListOptions) ([]*T, error) {
	var t T
	rlt := make([]*T, 0)

	if err := b.db.WithContext(ctx).Table(t.TableName()).Where(options.Query).Order(options.OrderBy).Scan(&rlt).Error; err != nil {
		return nil, err
	}
	return rlt, nil
}

func (b *GormMapper[T]) PageList(ctx context.Context, options PageListOptions) (int64, []*T, error) {
	var t T
	rlt := make([]*T, 0)
	count := int64(0)

	err := b.db.WithContext(ctx).Table(t.TableName()).Where(options.Query).Count(&count).Error

	if err != nil {
		return 0, nil, err
	}

	err = b.db.WithContext(ctx).Table(t.TableName()).Where(options.Query).
		Limit(options.PageSize).Offset((options.PageNum - 1) * options.PageSize).
		Order(options.OrderBy).Scan(&rlt).Error

	if err != nil {
		return 0, nil, err
	}
	return count, rlt, nil
}

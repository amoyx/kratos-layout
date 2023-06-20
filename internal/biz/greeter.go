package biz

import (
	"context"
	"strconv"
	"time"

	"github.com/go-kratos/kratos/v2/log"
)

// Greeter is a Greeter model.
type Greeter struct {
	Id   int64
	Name string
	Sex  string
}

func (g Greeter) GetCacheExpiration() time.Duration {
	return 100 * time.Second
}

func (g Greeter) GetCacheKey() string {
	return "helloworld:greeter_" + strconv.FormatInt(g.Id, 10)
}

// GreeterRepo is a Greater repo.
type GreeterRepo interface {
	BaseRepo[Greeter]
	ListBySex(context.Context, string) ([]*Greeter, error)
}

// GreeterUsecase is a Greeter usecase.
type GreeterUsecase struct {
	repo GreeterRepo
	log  *log.Helper
}

// NewGreeterUsecase new a Greeter usecase.
func NewGreeterUsecase(repo GreeterRepo, logger log.Logger) *GreeterUsecase {
	return &GreeterUsecase{repo: repo, log: log.NewHelper(logger)}
}

// Create creates a new greeter.
func (uc *GreeterUsecase) Create(ctx context.Context, g *Greeter) (*Greeter, error) {
	uc.log.WithContext(ctx).Info(g)
	rlt, err := uc.repo.Save(ctx, g)
	if err != nil {
		uc.log.WithContext(ctx).Error(err)
		return nil, err
	}
	return rlt, nil
}

// Get retrieves a greeter by ID.
func (uc *GreeterUsecase) Get(ctx context.Context, g *Greeter) (*Greeter, error) {
	uc.log.WithContext(ctx).Info(g)
	return uc.repo.FindByID(ctx, g)
}

// Update retrieves a greeter by greeter.
func (uc *GreeterUsecase) Update(ctx context.Context, g *Greeter) (*Greeter, error) {
	uc.log.WithContext(ctx).Info(g)
	rlt, err := uc.repo.Update(ctx, g)
	if err != nil {
		uc.log.WithContext(ctx).Error(err)
		return nil, err
	}
	return rlt, nil
}

// Delete deletes a greeter by ID.
func (uc *GreeterUsecase) Delete(ctx context.Context, g *Greeter) error {
	uc.log.WithContext(ctx).Info(g)
	return uc.repo.Delete(ctx, g)
}

// List lists greeters with optional filtering and pagination.
func (uc *GreeterUsecase) List(ctx context.Context) ([]*Greeter, error) {
	uc.log.WithContext(ctx).Info("list")
	return uc.repo.ListAll(ctx)
}

// QueryBySex query greeters with optional filtering and pagination.
func (uc *GreeterUsecase) QueryBySex(ctx context.Context, sex string) ([]*Greeter, error) {
	uc.log.WithContext(ctx).Info(sex)
	return uc.repo.ListBySex(ctx, sex)
}

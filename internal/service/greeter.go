package service

import (
	"context"

	"github.com/go-kratos/kratos/v2/log"

	v1 "github.com/go-kratos/kratos-layout/api/helloworld/v1"
	"github.com/go-kratos/kratos-layout/internal/assembler"
	"github.com/go-kratos/kratos-layout/internal/biz"
)

// GreeterService is a greeter service.
type GreeterService struct {
	v1.UnimplementedUserServiceServer
	uc  *biz.GreeterUsecase
	ba  assembler.BaseAssembler[v1.CreateUserRequest, biz.Greeter]
	bu  assembler.BaseAssembler[v1.UpdateUserRequest, biz.Greeter]
	ab  assembler.BaseAssembler[v1.UserResponse, biz.Greeter]
	abs assembler.BaseAssembler[[]*v1.UserResponse, []*biz.Greeter]
	log *log.Helper
}

func (g *GreeterService) GetUser(ctx context.Context, request *v1.GetUserRequest) (*v1.UserResponse, error) {
	g.log.WithContext(ctx).Info(request)
	dm := &biz.Greeter{Id: request.Id}
	rlt, err := g.uc.Get(ctx, dm)
	if err != nil {
		g.log.WithContext(ctx).Error(err)
		return nil, v1.ErrorUserNotFound(err.Error())
	}
	return g.ab.DomainToDto(rlt), nil
}

func (g *GreeterService) CreateUser(ctx context.Context, request *v1.CreateUserRequest) (*v1.UserResponse, error) {
	g.log.WithContext(ctx).Info(request)
	dm := g.ba.DTOToDomain(request)
	rlt, err := g.uc.Create(ctx, dm)
	if err != nil {
		g.log.WithContext(ctx).Error(err)
		return nil, v1.ErrorCreateUserFailed(err.Error())
	}
	return g.ab.DomainToDto(rlt), nil
}

func (g *GreeterService) UpdateUser(ctx context.Context, request *v1.UpdateUserRequest) (*v1.UserResponse, error) {
	g.log.WithContext(ctx).Info(request)
	dm := g.bu.DTOToDomain(request)
	rlt, err := g.uc.Update(ctx, dm)
	if err != nil {
		g.log.WithContext(ctx).Error(err)
		return nil, v1.ErrorUpdateUserFailed(err.Error())
	}
	return g.ab.DomainToDto(rlt), nil
}

func (g *GreeterService) DeleteUser(ctx context.Context, request *v1.DeleteUserRequest) (*v1.EmptyResponse, error) {
	g.log.WithContext(ctx).Info(request)
	dm := &biz.Greeter{Id: request.Id}
	err := g.uc.Delete(ctx, dm)
	if err != nil {
		g.log.WithContext(ctx).Error(err)
		return nil, v1.ErrorDeleteUserFailed(err.Error())
	}
	return &v1.EmptyResponse{}, nil
}

func (g *GreeterService) ListUsers(ctx context.Context, request *v1.ListUsersRequest) (*v1.ListUsersResponse, error) {
	g.log.WithContext(ctx).Info(request)
	tmp, err := g.uc.List(ctx)
	if err != nil {
		g.log.WithContext(ctx).Error(err)
		return nil, nil
	}

	rlt := g.abs.DomainToDto(&tmp)

	return &v1.ListUsersResponse{
		Users: *rlt,
	}, nil
}

// NewGreeterService new a greeter service.
func NewGreeterService(uc *biz.GreeterUsecase, logger log.Logger) *GreeterService {
	return &GreeterService{
		uc:  uc,
		log: log.NewHelper(logger),
	}
}

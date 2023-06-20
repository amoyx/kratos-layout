package service

import (
	"context"

	"github.com/go-kratos/kratos/v2/log"

	v2 "github.com/go-kratos/kratos-layout/api/helloworld/v2"
	"github.com/go-kratos/kratos-layout/internal/assembler"
	"github.com/go-kratos/kratos-layout/internal/biz"
)

// GreeterServiceV2 is a greeter service.
type GreeterServiceV2 struct {
	v2.UnimplementedUserServiceServer
	uc  *biz.GreeterUsecase
	ba  assembler.BaseAssembler[v2.CreateUserRequest, biz.Greeter]
	ab  assembler.BaseAssembler[v2.UserResponse, biz.Greeter]
	bu  assembler.BaseAssembler[v2.UpdateUserRequest, biz.Greeter]
	abs assembler.BaseAssembler[[]*v2.UserResponse, []*biz.Greeter]
	log *log.Helper
}

// NewGreeterServiceV2 new a greeter service.
func NewGreeterServiceV2(uc *biz.GreeterUsecase, logger log.Logger) *GreeterServiceV2 {
	return &GreeterServiceV2{
		uc:  uc,
		log: log.NewHelper(logger),
	}
}

func (g *GreeterServiceV2) GetUser(ctx context.Context, request *v2.GetUserRequest) (*v2.UserResponse, error) {
	g.log.WithContext(ctx).Info(request)
	dm := &biz.Greeter{Id: request.Id}
	rlt, err := g.uc.Get(ctx, dm)
	if err != nil {
		g.log.WithContext(ctx).Error(err)
		return nil, v2.ErrorUserNotFound(err.Error())
	}
	return g.ab.DomainToDto(rlt), nil
}

func (g *GreeterServiceV2) CreateUser(ctx context.Context, request *v2.CreateUserRequest) (*v2.UserResponse, error) {
	g.log.WithContext(ctx).Info(request)
	dm := g.ba.DTOToDomain(request)
	rlt, err := g.uc.Create(ctx, dm)
	if err != nil {
		g.log.WithContext(ctx).Error(err)
		return nil, v2.ErrorCreateUserFailed(err.Error())
	}
	return g.ab.DomainToDto(rlt), nil
}

func (g *GreeterServiceV2) UpdateUser(ctx context.Context, request *v2.UpdateUserRequest) (*v2.UserResponse, error) {
	g.log.WithContext(ctx).Info(request)
	dm := g.bu.DTOToDomain(request)
	rlt, err := g.uc.Update(ctx, dm)
	if err != nil {
		g.log.WithContext(ctx).Error(err)
		return nil, v2.ErrorUpdateUserFailed(err.Error())
	}
	return g.ab.DomainToDto(rlt), nil
}

func (g *GreeterServiceV2) DeleteUser(ctx context.Context, request *v2.DeleteUserRequest) (*v2.EmptyResponse, error) {
	g.log.WithContext(ctx).Info(request)
	dm := &biz.Greeter{Id: request.Id}
	err := g.uc.Delete(ctx, dm)
	if err != nil {
		g.log.WithContext(ctx).Error(err)
		return nil, v2.ErrorDeleteUserFailed(err.Error())
	}
	return &v2.EmptyResponse{}, nil
}

func (g *GreeterServiceV2) ListUsers(ctx context.Context, request *v2.ListUsersRequest) (*v2.ListUsersResponse, error) {
	g.log.WithContext(ctx).Info(request)
	tmp, err := g.uc.List(ctx)
	if err != nil {
		g.log.WithContext(ctx).Error(err)
		return nil, nil
	}

	rlt := g.abs.DomainToDto(&tmp)

	return &v2.ListUsersResponse{
		Users: *rlt,
	}, nil
}

func (g *GreeterServiceV2) QueryUsersBySex(ctx context.Context, request *v2.QueryUsersBySexRequest) (*v2.QueryUsersBySexResponse, error) {
	g.log.WithContext(ctx).Info(request)
	tmp, err := g.uc.QueryBySex(ctx, request.Sex)
	if err != nil {
		g.log.WithContext(ctx).Error(err)
		return nil, err
	}

	rlt := g.abs.DomainToDto(&tmp)

	return &v2.QueryUsersBySexResponse{
		Users: *rlt,
	}, nil
}

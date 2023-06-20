package main

import (
	"context"

	"github.com/go-kratos/kratos/contrib/registry/eureka/v2"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/registry"
	"github.com/pkg/errors"

	"github.com/go-kratos/kratos-layout/internal/conf"
)

type registrySet []registry.Registrar

func NewRegistrySet(application *conf.Application) registrySet {
	rs := make(registrySet, 0, 4)

	//1.registry eureka
	if application.Eureka != nil {
		rs = append(rs, NewEurekaRegistrar(application.Eureka))
	}

	//2.other ...

	return rs
}

// When an error occurs in registry.Registrar within app.Run() in Kratos,
// it is thrown as an error, but many third-party components have primitive
// errors that are difficult to pinpoint. Here, rewrapping registry.Registrar
// can help quickly locate the error and facilitate troubleshooting.
// the usage of error handling should be determined based on practical needs.
type EurekaRegistrar struct {
	*eureka.Registry
}

func NewEurekaRegistrar(conf *conf.Eureka) registry.Registrar {
	rg, err := eureka.New([]string{conf.Address},
		eureka.WithEurekaPath(conf.Path),
	)
	if err != nil {
		panic(err)
	}
	return &EurekaRegistrar{
		rg,
	}
}

func (e *EurekaRegistrar) Register(ctx context.Context, service *registry.ServiceInstance) error {
	err := e.Registry.Register(ctx, service)
	if err != nil {
		err = errors.Wrap(err, "eureka registrar error")
		log.Error(err)
		return err
	}
	return nil
}

func (e *EurekaRegistrar) Deregister(ctx context.Context, service *registry.ServiceInstance) error {
	err := e.Registry.Register(ctx, service)
	if err != nil {
		err = errors.Wrap(err, "eureka deregister error")
		log.Error(err)
		return err
	}
	return nil
}

package data

import (
	"context"
	"strings"

	"github.com/google/wire"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"github.com/go-redis/redis/v8"

	"github.com/go-kratos/kratos-layout/internal/conf"
	"github.com/go-kratos/kratos-layout/internal/data/model"
)

// ProviderSet is service providers.
var ProviderSet = wire.NewSet(
	NewGreeterRepo,
)

// MapperSet is service providers.
var MapperSet = wire.NewSet(
	model.NewGreeterMapper,
	NewDB,
)

var CacheSet = wire.NewSet(
	NewGreeterCache,
	NewCache,
)

func NewDB(c *conf.Data) *gorm.DB {
	db, err := gorm.Open(mysql.Open(c.Database.Source), &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true,
	})
	if err != nil {
		panic("failed to connect database")
	}
	return db
}

func NewCache(c *conf.Data) redis.UniversalClient {
	if strings.ToLower(c.Redis.Type) == "cluster" {
		Client := redis.NewClusterClient(&redis.ClusterOptions{
			Addrs:    strings.Split(c.Redis.Addr, ","), //set redis cluster url
			Password: c.Redis.Password,                 //set password
		})

		_, err := Client.Ping(context.Background()).Result()

		if err != nil {
			panic(err)
		}
		return Client
	}

	Client := redis.NewClient(&redis.Options{
		Addr:     c.Redis.Addr,     //set redis cluster url
		Password: c.Redis.Password, //set password
		DB:       int(c.Redis.Index),
	})

	_, err := Client.Ping(context.Background()).Result()

	if err != nil {
		panic(err)
	}
	return Client
}

func NewFakeCache(c *conf.Data) redis.UniversalClient {
	return &redis.Client{}
}

func NewFakeDB(c *conf.Data) *gorm.DB {
	return &gorm.DB{}
}

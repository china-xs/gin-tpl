// Package redis
// @author: xs
// @date: 2022/3/10
// @Description: redis 简单处理事redis
package redis

import (
	"context"
	"github.com/go-redis/redis/v8"
	"github.com/google/wire"
	"github.com/spf13/viper"
	"time"
)

type Options struct {
	Addr     string `yaml:"addr"`
	Password string `yaml:"password"`
	DB       int    `yaml:"db"`
	PoolSize int32  `yaml:"poolSize"`
}

var ProviderSet = wire.NewSet(New, NewOps)

func NewOps(v *viper.Viper) (*Options, error) {
	var (
		err error
		o   = new(Options)
	)
	o.PoolSize = 100
	o.DB = 0

	if err = v.UnmarshalKey("redis", o); err != nil {
		return nil, err
	}
	return o, err
}

func New(o *Options) (*redis.Client, error) {
	rdb := redis.NewClient(&redis.Options{
		//Addr:     "localhost:16379",
		Addr:     o.Addr,
		Password: o.Password, // no password set
		DB:       o.DB,       // use default DB
		PoolSize: 100,        // 连接池大小
	})
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_, err := rdb.Ping(ctx).Result()
	if err != nil {
		return nil, err
	}
	return rdb, nil
}

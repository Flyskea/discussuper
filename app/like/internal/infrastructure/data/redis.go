package data

import (
	"context"
	"fmt"
	"harmoni/internal/conf"

	"github.com/google/wire"
	"github.com/redis/go-redis/v9"
)

var ProviderSet = wire.NewSet(
	NewRedis,
	wire.Bind(new(redis.UniversalClient), new(*redis.Client)),
)

func NewRedis(conf *conf.Redis) (*redis.Client, func(), error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:         fmt.Sprintf("%s:%d", conf.GetIp(), conf.GetPort()),
		Password:     conf.Password,
		DB:           int(conf.Database),
		PoolSize:     int(conf.GetPoolSize()),
		ReadTimeout:  conf.GetReadTimeout().AsDuration(),
		WriteTimeout: conf.GetWriteTimeout().AsDuration(),
	})

	cleanFunc := func() {
		rdb.Close()
	}

	ctx, cancel := context.WithTimeout(context.Background(), rdb.Options().ReadTimeout)
	defer cancel()
	_, err := rdb.Ping(ctx).Result()
	if err != nil {
		return nil, cleanFunc, err
	}

	return rdb, cleanFunc, nil
}

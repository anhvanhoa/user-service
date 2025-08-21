package bootstrap

import (
	"cms-server/domain/service/cache"
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisConfig struct {
	Addr        string
	Password    string
	DB          int
	Network     string
	MaxIdle     int
	MaxActive   int
	IdleTimeout int
}

type RedisIntance struct {
	config RedisConfig
	client *redis.Client
	ctx    context.Context
}

func NewRedis(c RedisConfig) cache.RedisConfigImpl {
	var ctx = context.Background()
	client := redis.NewClient(&redis.Options{
		Addr:         c.Addr,
		Password:     c.Password,
		DB:           c.DB,
		Network:      c.Network,
		PoolSize:     c.MaxActive,
		MinIdleConns: c.MaxIdle,
		PoolTimeout:  time.Duration(c.IdleTimeout) * time.Second,
	})

	ri := &RedisIntance{
		config: c,
		client: client,
		ctx:    ctx,
	}
	_, err := ri.client.Ping(ctx).Result()
	if err != nil {
		panic("Failed to connect to Redis: " + err.Error())
	}

	return ri
}

func (ri *RedisIntance) Set(key string, value []byte, time time.Duration) error {
	err := ri.client.Set(ri.ctx, key, value, time).Err()
	if err != nil {
		return err
	}
	return nil
}

func (ri *RedisIntance) Get(key string) ([]byte, error) {
	value, err := ri.client.Get(ri.ctx, key).Result()
	if err != nil {
		return nil, err
	}
	return []byte(value), nil
}

func (ri *RedisIntance) Delete(key string) error {
	err := ri.client.Del(ri.ctx, key).Err()
	if err != nil {
		return err
	}
	return nil
}

func (ri *RedisIntance) Reset() error {
	err := ri.client.FlushAll(ri.ctx).Err()
	if err != nil {
		return err
	}
	return nil
}

func (ri *RedisIntance) Close() error {
	err := ri.client.Close()
	if err != nil {
		return err
	}
	return nil
}

// func (ri *RedisIntance) Ping() error {
// 	_, err := ri.client.Ping(ri.ctx).Result()
// 	if err != nil {
// 		return err
// 	}
// 	return nil
// }

func NewRedisConfig(
	addr, password string, db int, network string,
	maxIde, maxActive, idleTimeout int,
) RedisConfig {
	return RedisConfig{
		Addr:        addr,
		Password:    password,
		DB:          db,
		Network:     network,
		MaxIdle:     maxIde,
		MaxActive:   maxActive,
		IdleTimeout: idleTimeout,
	}
}

package caches

import (
	"fmt"
	"time"

	"github.com/go-redis/redis"
	"github.com/pkg/errors"
	"github.com/xgxw/foundation-go"
)

const (
	defaultRedisHost = "127.0.0.1"
	defaultRedisPort = 6379
)

type RedisOptions struct {
	Host     string `yaml:"host" mapstructure:"host"`
	Port     int    `yaml:"port" mapstructure:"port"`
	Password string `yaml:"password" mapstructure:"password"`
	DB       int    `yaml:"db" mapstructure:"db"`
	PoolSize int    `yaml:"pool_size" mapstructure:"pool_size"`
}

func (opts *RedisOptions) loadDefaults() {
	if opts.Host == "" {
		opts.Host = defaultRedisHost
	}
	if opts.Port == 0 {
		opts.Port = defaultRedisPort
	}
}

type RedisCache struct {
	client *redis.Client
}

var _ foundation.CacheAccessor = new(RedisCache)

func NewRedisCache(opts RedisOptions) (*RedisCache, error) {
	opts.loadDefaults()

	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", opts.Host, opts.Port),
		Password: opts.Password,
		PoolSize: opts.PoolSize,
		DB:       opts.DB,
	})
	if _, err := client.Ping().Result(); err != nil {
		return nil, err
	}

	cache := &RedisCache{client: client}
	return cache, nil
}

func (c *RedisCache) Get(key string) (s string, err error) {
	s, err = c.client.Get(key).Result()
	if err == redis.Nil {
		err = nil
	}
	if err != nil {
		err = errors.Wrap(err, "redis get failed")
	}
	return
}

func (c *RedisCache) Set(key string, value string, expiresIn int) (err error) {
	err = c.client.Set(key, value, time.Second*time.Duration(expiresIn)).Err()
	if err != nil {
		err = errors.Wrap(err, "redis set failed")
	}
	return
}

func (c *RedisCache) Delete(key string) (err error) {
	err = c.client.Del(key).Err()
	if err != nil {
		err = errors.Wrap(err, "redis del failed")
	}
	return
}

func (c *RedisCache) IncrBy(key string, value int64) (result int64, err error) {
	result, err = c.client.IncrBy(key, value).Result()
	if err != nil {
		err = errors.Wrap(err, "redis incr failed")
	}
	return
}

func (c *RedisCache) Exists(key string) (bool, error) {
	result, err := c.client.Exists(key).Result()
	if err != nil {
		return false, errors.Wrap(err, "redis exists failed")
	}
	return result > 0, nil
}

func (c *RedisCache) Redis() *redis.Client {
	return c.client
}

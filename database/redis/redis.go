// Package redis
//
// @author: xwc1125
package redis

import (
	"context"
	"os"
	"sync"

	"github.com/chain5j/logger"
	"github.com/go-redis/redis/v8"
	"github.com/spf13/viper"
	"github.com/xwc1125/xwc1125-pkg/utils/stringutil"
)

type Mode int

const (
	Normal Mode = iota
	Guard
	Cluster
)

type DB struct {
	log    logger.Logger
	config RedisConfig
	*redis.Client
	cluster *redis.ClusterClient
}

var (
	cache     *DB
	cacheOnce sync.Once
)

func Cache() *DB {
	cacheOnce.Do(func() {
		var config RedisConfig
		err := viper.UnmarshalKey("redis", &config)
		if err != nil {
			logger.Crit("get redis config err", "err", err)
		}
		cache, err = New(config)
		if err != nil {
			logger.Crit("new redis err", "err", err)
			os.Exit(1)
		}
	})
	return cache
}

func New(config RedisConfig) (*DB, error) {
	var (
		client  *redis.Client
		cluster *redis.ClusterClient
	)
	switch config.Mode {
	case Normal:
		opt := &redis.Options{
			Addr: config.Addr[0],
			DB:   config.DB, // use default DB
		}
		if !stringutil.IsEmpty(config.Password) {
			opt.Password = config.Password
		}
		client = redis.NewClient(opt)
	case Guard:
		// 连接哨兵模式
		opt := &redis.FailoverOptions{
			MasterName:    config.MasterName,
			SentinelAddrs: config.Addr,
		}
		if !stringutil.IsEmpty(config.Password) {
			opt.Password = config.Password
		}
		client = redis.NewFailoverClient(opt)
	case Cluster:
		opt := &redis.ClusterOptions{
			Addrs: config.Addr,
		}
		cluster = redis.NewClusterClient(opt)
		if !stringutil.IsEmpty(config.Password) {
			opt.Password = config.Password
		}
	}

	var logger = logger.Log("redis")
	_, err := client.Ping(context.Background()).Result()
	if err != nil {
		logger.Error("redis连接失败", "err", err)
		return nil, err
	}
	logger.Debug("redis连接成功")

	return &DB{
		log:     logger,
		config:  config,
		Client:  client,
		cluster: cluster,
	}, nil
}

func (r *DB) Redis() *redis.Client {
	return r.Client
}

func (r *DB) ClusterRedis() *redis.ClusterClient {
	return r.cluster
}

// Cmd 命令
func (r *DB) Cmd() redis.Cmdable {
	if r.Client != nil {
		return r.Client
	}
	if r.cluster != nil {
		return r.cluster
	}
	return nil
}

// PubSub 订阅
func (r *DB) PubSub() redis.UniversalClient {
	if r.Client != nil {
		return r.Client
	}
	if r.cluster != nil {
		return r.cluster
	}
	return nil
}

func (r *DB) Ping() bool {
	_, err := r.Client.Ping(context.Background()).Result()
	if err != nil {
		logger.Info("redis ping", "err", err)
		return false
	}
	return true
}

func (r *DB) IsExist(key string) bool {
	v, err := r.Client.Exists(context.Background(), key).Result()
	if err != nil {
		return false
	}

	return v != 0
}

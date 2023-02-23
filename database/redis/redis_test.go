// Package redis
//
// @author: xwc1125
package redis

import (
	"fmt"
	"log"
	"testing"
	"time"
)

func TestNew(t *testing.T) {
	redis, err := New(RedisConfig{
		Mode: 0,
		Addr: []string{"127.0.0.1:6379"},
		DB:   1,
	})
	if err != nil {
		panic(err)
	}
	redis.Cmd().Set("1", "2", 0)
	cmd := redis.Cmd().Get("1")
	result, err := cmd.Result()
	fmt.Println(result)

	redis.Redis().Subscribe()
}

func TestNew2(t *testing.T) {
	redis, err := New(RedisConfig{
		Mode: 0,
		Addr: []string{"127.0.0.1:6379"},
		DB:   1,
	})
	if err != nil {
		panic(err)
	}

	_, err = redis.Cmd().Ping().Result()
	if err != nil {
		log.Printf("redis连接失败,错误信息:%v\n", err)
		return
	}
	log.Println("redis连接成功")

	// 发布订阅
	pubsub := redis.PubSub().Subscribe("channel1")

	// 发布前先订阅
	_, err = pubsub.Receive()
	if err != nil {
		panic(err)
	}

	// Go channel接收信息
	ch := pubsub.Channel()

	// 发布消息
	err = redis.Cmd().Publish("channel1", "hello").Err()
	if err != nil {
		panic(err)
	}

	time.AfterFunc(time.Second, func() {
		// 订阅和通道都关闭
		_ = pubsub.Close()
	})

	// 消费
	for msg := range ch {
		log.Println(msg.Channel, msg.Payload)
	}
}

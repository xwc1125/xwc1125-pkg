// Package captcha
//
// @author: xwc1125
package captcha

import (
	"github.com/chain5j/chain5j-pkg/crypto/hashalg"
	"github.com/chain5j/chain5j-pkg/util/hexutil"
	"github.com/mojocn/base64Captcha"
	"github.com/xwc1125/xwc1125-pkg/database/redis"
)

var (
	_                base64Captcha.Store = new(RedisStore)
	CaptchaPrefixKey                     = "captcha_"
)

type Store base64Captcha.Store

type RedisStore struct {
	db *redis.DB
}

func NewRedisStoreByConfig(config redis.RedisConfig) (*RedisStore, error) {
	db, err := redis.New(config)
	if err != nil {
		return nil, err
	}
	return &RedisStore{
		db: db,
	}, nil
}
func NewRedisStore(db *redis.DB) (*RedisStore, error) {
	return &RedisStore{
		db: db,
	}, nil
}

func (s *RedisStore) Set(id string, value string) error {
	return s.db.Set(getKey(id), valToHex(value), base64Captcha.Expiration).Err()
}

func getKey(id string) string {
	return CaptchaPrefixKey + id
}

func valToHex(value string) string {
	return hexutil.Encode(hashalg.Sha256([]byte(value)))
}
func (s *RedisStore) Get(id string, clear bool) string {
	realId := getKey(id)
	if clear {
		defer s.db.Del(realId)
	}
	result, err := s.db.Get(realId).Result()
	if err != nil {
		return ""
	}
	return result
}

func (s *RedisStore) Verify(id, answer string, clear bool) bool {
	v := s.Get(id, clear)
	return v == valToHex(answer)
}

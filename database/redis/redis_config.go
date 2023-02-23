// Package redis
//
// @author: xwc1125
package redis

// RedisConfig reids配置
type RedisConfig struct {
	Mode        Mode     `json:"is_guard" mapstructure:"is_guard"`
	Addr        []string `json:"addr" mapstructure:"addr"`
	Username    string   `json:"username" mapstructure:"username" yaml:"username"`
	Password    string   `json:"password" mapstructure:"password"`
	DB          int      `json:"db" mapstructure:"db"`
	MasterName  string   `json:"master_name" mapstructure:"master_name"`
	CacheExpire int64    `json:"cache_expire" mapstructure:"cache_expire"` // 过期时间[秒]
}

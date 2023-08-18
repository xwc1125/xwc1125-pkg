// Package email
package email

type Config struct {
	Addr      string `json:"addr" mapstructure:"addr" yaml:"addr"`
	Port      string `json:"port" mapstructure:"port" yaml:"port"`
	Username  string `json:"username" mapstructure:"username" yaml:"username"`
	Password  string `json:"password" mapstructure:"password" yaml:"password"`
	PoolCount int    `json:"pool_count" mapstructure:"pool_count" yaml:"pool_count"`
}

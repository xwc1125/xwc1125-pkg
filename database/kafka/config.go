// Package kafka
//
// @author: xwc1125
// @date: 2021/3/23
package kafka

// KafkaConfig kafka配置
type KafkaConfig struct {
	IsAsync      bool     `json:"is_async" mapstructure:"is_async"`           // 是否为异步
	Addrs        []string `json:"addrs" mapstructure:"addrs"`                 // kafka服务的地址
	GroupId      string   `json:"group_id" mapstructure:"group_id"`           // groupId
	Topic        []string `json:"topic" mapstructure:"topic"`                 // topic
	IsLog        bool     `json:"is_log" mapstructure:"is_log"`               // 是否打印日志
	KafkaVersion string   `json:"kafka_version" mapstructure:"kafka_version"` // kafka版本
	Strategy     string   `json:"strategy" mapstructure:"strategy"`           // group Rebalance策略
	SASLEnable   bool     `json:"sasl_enable" mapstructure:"sasl_enable"`     // 是否开启SASL
	SASLUser     string   `json:"sasl_user" mapstructure:"sasl_user"`         // SASL用户名
	SASLPassword string   `json:"sasl_password" mapstructure:"sasl_password"` // SASL密码
}

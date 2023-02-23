// Package ipfilter
//
// @author: xwc1125
package ipfilter

// ConnFilter 连接过滤器接口定义
type ConnFilter interface {
	// OnConnected 客户端连接建立
	// 返回false则关闭连接，同时返回需要关闭连接的原因
	OnConnected(ip string) (bool, string)
	GetAbnConn(ip string) int // 获取异常连接数
}

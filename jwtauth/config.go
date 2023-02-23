// Package jwtauth
//
// @author: xwc1125
package jwtauth

import (
	"fmt"

	"github.com/casbin/casbin/v2/util"
)

const (
	DefaultJwtContextKey = "key-jwt"
	AuthorizationKEY     = "Authorization"
)

// JWTConfig jwt配置信息
type JWTConfig struct {
	AuthorizationKey string `json:"authorization_key" mapstructure:"authorization_key" yaml:"authorization_key"` // 请求中header中token key
	ParamTokenKey    string `json:"param_token_key" mapstructure:"param_token_key" yaml:"param_token_key"`       // 请求中params中token key
	CookieTokenKey   string `json:"cookie_token_key" mapstructure:"cookie_token_key" yaml:"cookie_token_key"`    // 请求中cookie中token key
	JwtContextKey    string `json:"jwt_context_key" mapstructure:"jwt_context_key" yaml:"jwt_context_key"`       // 缓存token的key

	PriKeyPath          string   `json:"pri_key_path" mapstructure:"pri_key_path"`                                                   // rsa私钥
	PubKeyPath          string   `json:"pub_key_path" mapstructure:"pub_key_path"`                                                   // rsa公钥
	Secret              string   `json:"secret" mapstructure:"secret"`                                                               // token 密钥
	Timeout             int64    `json:"timeout" mapstructure:"timeout"`                                                             // token 过期时间 单位：秒
	RefreshTimeout      int64    `json:"refresh_timeout" mapstructure:"refresh_timeout" yaml:"refresh_timeout"`                      // refreshToken 过期时间 单位：秒
	IgnoreURLs          []string `json:"ignore_urls" mapstructure:"ignore_urls"`                                                     // 忽略的url
	OnlineKey           string   `json:"online_key" mapstructure:"online_key"`                                                       // 在线标识
	EnableAuthOnOptions bool     `json:"enable_auth_on_options" mapstructure:"enable_auth_on_options" yaml:"enable_auth_on_options"` // 是否启动OPTIONS方法的所有请求都将使用身份验证
	Debug               bool     `json:"debug" mapstructure:"debug" yaml:"debug"`                                                    // 是否启动debug，进行日志输出
}

func (c *JWTConfig) IsIgnore(reqPath, reqMethod string) bool {
	reqMPath := fmt.Sprintf("%s:%s", reqMethod, reqPath)
	for _, path := range c.IgnoreURLs {
		if util.KeyMatch2(reqMPath, path) {
			return true
		}
	}
	return false
}

// Package database
//
// @author: xwc1125
package database

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"strings"

	"github.com/xwc1125/xwc1125-pkg/secure/password"
)

type MysqlConfig struct {
	Driver          string `json:"driver" mapstructure:"driver"`
	Url             string `json:"url" mapstructure:"url"`
	Username        string `json:"username" mapstructure:"username"`
	Password        string `json:"password" mapstructure:"password"`
	Secret          string `json:"secret" mapstructure:"secret"`
	MaxIdleConns    int    `json:"max_idle_conns" mapstructure:"max_idle_conns"`
	MaxOpenConns    int    `json:"max_open_conns" mapstructure:"max_open_conns"`
	ConnMaxLifetime int    `json:"conn_max_lifetime" mapstructure:"conn_max_lifetime"`
	ConnMaxIdleTime int    `json:"conn_max_idle_time" mapstructure:"conn_max_idle_time"`
	PrefixTable     string `json:"prefix_table" mapstructure:"prefix_table"`
	PrefixColumn    string `json:"prefix_column" mapstructure:"prefix_column"`
	ShowSQL         bool   `json:"show_sql" mapstructure:"show_sql"`
	LogLevel        int    `json:"log_level" mapstructure:"log_level"`
}

// DSNPrint dns pring
func (d MysqlConfig) DSNPrint() string {
	return fmt.Sprintf("%s:@%s", d.Username, d.Url)
}

// DSN datasource
func (d MysqlConfig) DSN() string {
	var buf bytes.Buffer
	username := strings.TrimSpace(d.Username)
	if username != "" {
		buf.WriteString(d.Username + ":")
	}
	// 需要对密码进行加解密
	pwd := strings.TrimSpace(d.Password)
	pwdBase64, err := base64.StdEncoding.DecodeString(pwd)
	if err == nil {
		key, err := password.DecryptKey(pwdBase64, d.Secret)
		if err == nil && key != nil {
			pwd = string(key.PrivateKey)
		}
	}
	if pwd != "" {
		buf.WriteString(pwd)
		if d.Driver == "mysql" {
			buf.WriteString("@")
		}
	}
	buf.WriteString(d.Url)
	return buf.String()
}

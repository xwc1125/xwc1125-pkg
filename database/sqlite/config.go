// Package sqlite
//
// @author: xwc1125
package sqlite

type SqliteConfig struct {
	Datasource      string `json:"datasource" mapstructure:"datasource"`
	MaxIdleConns    int    `json:"max_idle_conns" mapstructure:"max_idle_conns"`
	MaxOpenConns    int    `json:"max_open_conns" mapstructure:"max_open_conns"`
	ConnMaxLifetime int    `json:"conn_max_lifetime" mapstructure:"conn_max_lifetime"`
	ConnMaxIdleTime int    `json:"conn_max_idle_time" mapstructure:"conn_max_idle_time"`
}

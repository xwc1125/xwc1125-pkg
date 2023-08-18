// Package config
package config

type RemoteConfig interface {
	ReadConfig(result interface{}) error
	WatchConfig() <-chan bool
	AllSettings() map[string]interface{}
	Stop() error
}

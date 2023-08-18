// Package apollo
package apollo

type Config struct {
	Debug      bool   `json:"debug" mapstructure:"debug" yaml:"debug"`
	Provider   string `json:"provider" mapstructure:"provider" yaml:"provider"`
	AppID      string `json:"app_id" mapstructure:"app_id" yaml:"app_id"`
	Endpoint   string `json:"endpoint" mapstructure:"endpoint" yaml:"endpoint"`
	Cluster    string `json:"cluster" mapstructure:"cluster" yaml:"cluster"`
	Namespace  string `json:"namespace" mapstructure:"namespace" yaml:"namespace"`
	ConfigType string `json:"config_type" mapstructure:"config_type" yaml:"config_type"`
	Secret     string `json:"secret" mapstructure:"secret" yaml:"secret"`
	ReleaseKey string `json:"release_key" mapstructure:"release_key" yaml:"release_key"`
	Label      string `json:"label" mapstructure:"label" yaml:"label"`
}

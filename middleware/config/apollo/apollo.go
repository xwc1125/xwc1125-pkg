// Package apollo
package apollo

import (
	"github.com/chain5j/logger"
	remote "github.com/shima-park/agollo/viper-remote"
	"github.com/spf13/viper"
)

var (
	configPath string
)

type Apollo struct {
	log          logger.Logger
	apolloConfig *Config
	app          *viper.Viper
	provider     *remote.Extender
}

func New(apolloConfig *Config) (*Apollo, error) {
	log := logger.Log("apollo")
	provider, err := remote.NewApolloProvider(remote.ApolloConfig{
		Endpoint:   apolloConfig.Endpoint,
		AppID:      apolloConfig.AppID,
		ConfigType: apolloConfig.ConfigType,
		Namespace:  apolloConfig.Namespace,
	})
	if err != nil {
		log.Error("new apollo provider err", "err", err)
		return nil, err
	}

	app := provider.GetViper()
	// 根据namespace实际格式设置对应type
	if len(apolloConfig.Secret) > 0 {
		err = app.AddSecureRemoteProvider(apolloConfig.Provider, apolloConfig.Endpoint, apolloConfig.Namespace, apolloConfig.Secret)
	} else {
		err = app.AddRemoteProvider(apolloConfig.Provider, apolloConfig.Endpoint, apolloConfig.Namespace)
	}
	if err != nil {
		log.Error("viper remote provider err", "err", err)
		return nil, err
	}
	client := &Apollo{
		log:          log,
		apolloConfig: apolloConfig,
		app:          app,
		provider:     provider,
	}
	return client, nil
}

func (a *Apollo) AllSettings() map[string]interface{} {
	return a.app.AllSettings()
}

func (a *Apollo) ReadConfig(result interface{}) error {
	err := a.app.ReadRemoteConfig()
	if err != nil {
		a.log.Error("read apollo config err", "err", err)
		return err
	}
	if a.apolloConfig.Debug {
		a.log.Info("read remote config", "settings", a.app.AllSettings())
	}

	// 直接反序列化到结构体中
	err = a.app.Unmarshal(&result)
	if err != nil {
		a.log.Error("unmarshal remote config to result err", "err", err)
		return err
	}
	return nil
}

func (a *Apollo) WatchConfig() <-chan bool {
	return a.provider.WatchRemoteConfigOnChannel()
}
func (a *Apollo) Stop() error {
	return nil
}

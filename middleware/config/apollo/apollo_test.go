// Package apollo
package apollo

import (
	"fmt"
	"sync"
	"testing"

	"github.com/mitchellh/mapstructure"
	"github.com/shima-park/agollo"
	remote "github.com/shima-park/agollo/viper-remote"
)

func TestApollo(t *testing.T) {
	var apolloConfig = Config{
		Debug:      true,
		Provider:   "apollo",
		AppID:      "test2",
		Endpoint:   "http://127.0.0.1:8080",
		Cluster:    "default",
		Namespace:  "app-test.yaml",
		ConfigType: "yaml",
	}
	apollo, err := New(&apolloConfig)
	if err != nil {
		panic(err)
	}
	var demoConf Config
	err = apollo.ReadConfig(&demoConf)
	if err != nil {
		panic(err)
	}
	fmt.Println(demoConf)
	respChan := apollo.WatchConfig()
	for {
		<-respChan
		// on changed and notify
		fmt.Println("app.AllSettings==>:", apollo.AllSettings())
	}
}
func TestApollo2(t *testing.T) {
	remoteProvider, _ := remote.NewApolloProvider(remote.ApolloConfig{
		Endpoint:   "localhost:8080",
		AppID:      "test1",
		ConfigType: "yaml",
		Namespace:  "app-test.yaml",
	})

	v := remoteProvider.GetViper()
	// sync read remote config
	_ = v.ReadRemoteConfig()
	fmt.Println("app.AllSettings:", v.AllSettings())

	respChan := remoteProvider.WatchRemoteConfigOnChannel()
	go func(rc <-chan bool) {
		for {
			<-rc
			// on changed and notify
			fmt.Println("app.AllSettings:", v.AllSettings())
		}
	}(respChan)
	var wg sync.WaitGroup
	wg.Add(1)
	wg.Wait()
}

func TestApolloWatch(t *testing.T) {
	a, err := agollo.New("localhost:8080", "test1", agollo.FailTolerantOnBackupExists(), agollo.PreloadNamespaces("", "app-test.yaml"))
	if err != nil {
		panic(err)
	}

	errorCh := a.Start() // Start后会启动goroutine监听变化，并更新agollo对象内的配置cache
	get := a.Get("")
	fmt.Println(get)
	watchCh := a.Watch()
	for {
		select {
		case err := <-errorCh:
			panic(err)
		case resp := <-watchCh:
			fmt.Println(
				"Namespace:", resp.Namespace,
				"OldValue:", resp.OldValue,
				"NewValue:", resp.NewValue,
				"Error:", resp.Error,
			)
			var demoConf Config
			err := mapstructure.Decode(resp.NewValue, &demoConf)
			if err != nil {
				panic(err)
			}
			fmt.Println(demoConf)
		}
	}
}

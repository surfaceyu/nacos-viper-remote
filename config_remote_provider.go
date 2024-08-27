package nacos_viper_remote

import (
	"bytes"
	"fmt"

	"github.com/spf13/viper"
)

type ViperRemoteProvider struct {
	configType string
	configSet  string
}

func NewRemoteProvider(configType string) *ViperRemoteProvider {
	return &ViperRemoteProvider{
		configType: configType,
		configSet:  "yoyogo.cloud.discovery.metadata"}
}

func (provider *ViperRemoteProvider) GetProvider(runtimeViper *viper.Viper) *viper.Viper {
	var option *Option
	err := runtimeViper.Sub(provider.configSet).Unmarshal(&option)
	if err != nil {
		panic(err)
	}
	SetOptions(option)
	remote_viper := viper.New()
	_ = remote_viper.AddRemoteProvider("nacos", "localhost", "")
	if provider.configType == "" {
		provider.configType = "yaml"
	}
	remote_viper.SetConfigType(provider.configType)
	err = remote_viper.ReadRemoteConfig()
	if err == nil {
		fmt.Println("used remote viper")
		return remote_viper
	} else {
		panic(err)
	}
}

func (provider *ViperRemoteProvider) WatchRemoteConfigOnChannel(remoteViper *viper.Viper) <-chan bool {
	updater := make(chan bool)

	respChan, _ := viper.RemoteConfig.WatchChannel(DefaultRemoteProvider())
	go func(rc <-chan *viper.RemoteResponse) {
		for {
			b := <-rc
			reader := bytes.NewReader(b.Value)
			_ = remoteViper.ReadConfig(reader)
			// configuration on changed
			updater <- true
		}
	}(respChan)

	return updater
}

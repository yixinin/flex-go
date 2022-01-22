package config

import "github.com/yixinin/flex/registry"

type Config struct {
	App  string
	Etcd registry.EtcdConfig `mapstructure:"etcd"`
}

var conf *Config

func Init() {
	conf = new(Config)
}

func GetConfig() *Config {
	return conf
}

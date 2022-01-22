package config

import (
	"context"
	"fmt"
	"os"

	"github.com/pelletier/go-toml"
	"github.com/yixinin/flex/registry"
	"github.com/yixinin/flex/topic"
)

type Config struct {
	App      string              `mapstructure:"app"`
	LogLevel string              `mapstructure:"level"`
	Etcd     registry.EtcdConfig `mapstructure:"etcd"`
	Router   string              `mapstructure:"router"`
	Buffer   string              `mapstructure:"buffer"`
	Topics   []topic.Config      `mapstructure:"topics"`
}

var conf *Config

func Init(ctx context.Context) {
	conf = new(Config)

	buf, err := os.ReadFile("./config/example.toml")
	if err != nil {
		panic(err)
	}
	err = toml.Unmarshal(buf, conf)
	if err != nil {
		panic(err)
	}
	fmt.Printf("read config:%+v\n", conf)
}

func GetConfig() *Config {
	return conf
}

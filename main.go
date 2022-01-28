package main

import (
	"context"
	"strings"

	"github.com/sirupsen/logrus"
	"github.com/yixinin/flex/config"
	"github.com/yixinin/flex/logger"
	"github.com/yixinin/flex/registry"
)

func main() {
	var rawCtx = context.Background()
	config.Init(context.Background())
	log := logrus.New()
	log.Level = func() logrus.Level {
		switch strings.ToLower(config.GetConfig().LogLevel) {
		case "info":
			return logrus.InfoLevel
		case "warn":
			return logrus.WarnLevel
		case "error":
			return logrus.ErrorLevel
		}
		return logrus.DebugLevel
	}()
	logger.Init(logger.FromLogrus(log))

	registry.Init(config.GetConfig().Etcd)
	go registry.RegisterAddr(rawCtx, config.GetConfig().Port)

	m := NewManager(rawCtx, config.GetConfig().Topics)
	m.Run(rawCtx)
}

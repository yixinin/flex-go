package main

import (
	"context"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/yixinin/flex/config"
	"github.com/yixinin/flex/logger"
	"github.com/yixinin/flex/registry"
)

func main() {
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
	var ctx, cancel = context.WithCancel(context.Background())
	defer cancel()
	var delayCtx, delayCancel = context.WithCancel(context.Background())
	defer delayCancel()

	registry.Init(config.GetConfig().Etcd)
	go registry.RegisterAddr(ctx)

	m := NewManager(delayCtx)
	m.Run(ctx)
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, os.Interrupt, syscall.SIGTERM)

	<-sigs
	cancel()
	logger.Info(ctx, "flex will start exit after 30s ...")
	select {
	case <-time.After(30 * time.Second):
		delayCancel()
	}
	logger.Info(ctx, "flex wait all job finish ...")
	m.Wait()
	logger.Info(ctx, "flex exited")

}

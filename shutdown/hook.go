package shutdown

import (
	"context"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/yixinin/flex/logger"
)

type Runner interface {
	Run(ctx context.Context) error
}

type ShutdownHooker interface {
	Runner
	BeforeShutdown()
	ShutDown()
	AfterShutdown()
}

type ShutDownHook struct {
	hookers []ShutdownHooker
	quit    chan struct{}
}

func NewShutdownHook(quit chan struct{}, hooks ...ShutdownHooker) *ShutDownHook {
	return &ShutDownHook{
		hookers: hooks,
		quit:    quit,
	}
}

func (h *ShutDownHook) Add(hooker ShutdownHooker) {
	h.hookers = append(h.hookers, hooker)
}

func (h *ShutDownHook) Run(ctx context.Context) {
	var wg sync.WaitGroup
	for _, task := range h.hookers {
		wg.Add(1)
		go func(task Runner) {
			defer func() {
				recover()
				wg.Done()
			}()
			err := task.Run(ctx)
			if err != nil {
				logger.Error(ctx, err)
			}
		}(task)
	}

	if h.quit == nil {
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		<-quit
	} else {
		<-h.quit
	}

	for _, hooker := range h.hookers {
		hooker.BeforeShutdown()
	}
	// waitting..
	wg.Wait()

	for _, hooker := range h.hookers {
		hooker.ShutDown()
	}

	logger.Info(ctx, "handle after shutdown ...")
	for _, hooker := range h.hookers {
		hooker.AfterShutdown()
	}
}

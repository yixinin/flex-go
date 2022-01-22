package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	var ctx, cancel = context.WithCancel(context.Background())
	defer cancel()

	m := NewManager()
	m.Run(ctx)
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, os.Interrupt, syscall.SIGTERM)

	<-sigs
	cancel()

	m.Wait()
}

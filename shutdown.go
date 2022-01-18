package main

import (
	"os"
	"os/signal"
	"syscall"
)

// ShutdownHandler waits for SIGINT/SIGTERM signals and calls cancelFunc
func ShutdownHandler(cancelFunc func()) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)

	<-c
	cancelFunc()
}

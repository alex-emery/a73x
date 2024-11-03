package main

import (
	"os"
	"os/signal"
	"syscall"

	"git.sr.ht/~a73x/home/web"
	"go.uber.org/zap"
)

func main() {
	logger, _ := zap.NewProduction()
	server, err := web.New(logger)
	if err != nil {
		logger.Fatal("failed to start webserver", zap.Error(err))
	}

	done := make(chan os.Signal, 0)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	logger.Info("Starting web server")

	go func() {
		if err := server.ListenAndServe(); err != nil {
			logger.Error("fail running web server", zap.Error(err))
		}
		done <- nil
	}()

	<-done
	logger.Info("Stopping web server")
}

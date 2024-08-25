package main

import (
	"context"
	"errors"
	"os"
	"os/signal"
	"syscall"
	"time"

	"git.sr.ht/~a73x/home/tui"
	"git.sr.ht/~a73x/home/web"

	"github.com/charmbracelet/log"
	"github.com/charmbracelet/ssh"
)

const (
	port = "2222"
)

func main() {
	host := os.Getenv("DOMAIN")
	if host == "" {
		host = "localhost"
	}
	server, err := web.New()

	s, err := tui.New(host, port)

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	log.Info("Starting SSH server", "host", host, "port", port)
	go func() {
		if err := server.ListenAndServe(); err != nil {
			log.Error("fail running web server", "error", err)
		}
		done <- nil
	}()
	go func() {
		if err = s.ListenAndServe(); err != nil && !errors.Is(err, ssh.ErrServerClosed) {
			log.Error("Could not start server", "error", err)
			done <- nil
		}
	}()

	<-done
	log.Info("Stopping SSH server")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer func() { cancel() }()
	if err := s.Shutdown(ctx); err != nil && !errors.Is(err, ssh.ErrServerClosed) {
		log.Error("Could not stop server", "error", err)
	}
}

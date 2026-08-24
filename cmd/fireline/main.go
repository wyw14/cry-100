package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/wyw14/cry-100/internal/api"
	"github.com/wyw14/cry-100/internal/app"
)

func main() {
	config := app.DefaultConfig()
	if !config.Valid() {
		log.Fatal("invalid FireLine configuration")
	}
	runtime, err := app.New(config.DataDir)
	if err != nil {
		log.Fatal(err)
	}
	runtime.Seed()
	server := &http.Server{Addr: config.Address, Handler: api.NewServer(runtime).Router(), ReadHeaderTimeout: 5 * time.Second}
	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt, syscall.SIGTERM)
	<-signals
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(config.ShutdownSeconds)*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		log.Printf("shutdown HTTP server: %v", err)
	}
	if err := runtime.SaveSnapshot(); err != nil {
		log.Printf("save recovery snapshot: %v", err)
	}
}

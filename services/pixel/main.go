package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"github.com/mikestefanello/gcpworkflow/config"
	"github.com/mikestefanello/gcpworkflow/pkg/handlers"
	"github.com/mikestefanello/gcpworkflow/pkg/repo"
)

func main() {
	// Load configuration
	cfg, err := config.GetConfig()
	if err != nil {
		panic(fmt.Sprintf("unable to configuration: %v", err))
	}

	// Initialize web
	e := echo.New()
	e.Logger.SetLevel(log.INFO)

	// Pixel repository
	pr, err := repo.NewPubSubRepo(cfg)
	if err != nil {
		panic(fmt.Sprintf("unable to initialize pixel repo: %v", err))
	}

	// Build the router
	handlers.BuildRouter(e, cfg, pr)

	// Start the server
	go func() {
		srv := http.Server{
			Addr:         fmt.Sprintf("%s:%d", cfg.HTTP.Hostname, cfg.HTTP.Port),
			Handler:      e,
			ReadTimeout:  cfg.HTTP.ReadTimeout,
			WriteTimeout: cfg.HTTP.WriteTimeout,
			IdleTimeout:  cfg.HTTP.IdleTimeout,
		}

		if err := e.StartServer(&srv); err != http.ErrServerClosed {
			e.Logger.Fatalf("shutting down the server: v", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with a timeout of 10 seconds.
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	signal.Notify(quit, os.Kill)
	<-quit
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := e.Shutdown(ctx); err != nil {
		e.Logger.Fatal(err)
	}
}

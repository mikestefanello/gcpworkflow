package web

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/labstack/echo/v4"
	echomw "github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
	"github.com/mikestefanello/gcpworkflow/pkg/config"
)

func Init(cfg config.Config) *echo.Echo {
	// Initialize web
	e := echo.New()
	e.Logger.SetLevel(log.INFO)

	// Global middleware
	e.Use(
		echomw.RemoveTrailingSlashWithConfig(echomw.TrailingSlashConfig{
			RedirectCode: http.StatusMovedPermanently,
		}),
		echomw.Recover(),
		echomw.Secure(),
		echomw.RequestID(),
		echomw.Logger(),
		echomw.TimeoutWithConfig(echomw.TimeoutConfig{
			Timeout: cfg.App.Timeout,
		}),
		echomw.CORSWithConfig(echomw.CORSConfig{
			AllowMethods: []string{http.MethodGet},
		}),
	)

	// Error handler
	eh := new(errorHandler)
	e.HTTPErrorHandler = eh.handle

	return e
}

func Start(cfg config.Config, e *echo.Echo) {
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

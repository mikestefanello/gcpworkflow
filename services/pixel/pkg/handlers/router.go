package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
	echomw "github.com/labstack/echo/v4/middleware"
	"github.com/mikestefanello/gcpworkflow/config"
	"github.com/mikestefanello/gcpworkflow/pkg/middleware"
	"github.com/mikestefanello/gcpworkflow/pkg/pixel"
)

func BuildRouter(e *echo.Echo, cfg config.Config, pr pixel.Repository) {
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

	// Pixel route
	ph := newPixelHandler(pr)
	e.GET("/", ph.get, middleware.CacheControl(0)).Name = "pixel.get"
}

package main

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/mikestefanello/gcpworkflow/pkg/config"
	"github.com/mikestefanello/gcpworkflow/pkg/pixel"
	"github.com/mikestefanello/gcpworkflow/pkg/web"
)

func main() {
	// Load configuration
	cfg, err := config.GetConfig()
	if err != nil {
		panic(fmt.Sprintf("unable to configuration: %v", err))
	}

	// Initialize web
	e := web.Init(cfg)

	// Cleaner route
	e.POST("/", cleaner)

	// Start the server
	web.Start(cfg, e)
}

func cleaner(ctx echo.Context) error {
	var p pixel.Pixel

	if err := ctx.Bind(&p); err != nil {
		return err
	}

	p.ExtraField = "extra field added by cleaner"
	
	return ctx.JSON(http.StatusOK, p)
}

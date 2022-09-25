package main

import (
	"fmt"

	"github.com/mikestefanello/gcpworkflow/pkg/config"
	"github.com/mikestefanello/gcpworkflow/pkg/middleware"
	"github.com/mikestefanello/gcpworkflow/pkg/web"
	"github.com/mikestefanello/gcpworkflow/services/pixel/pkg/handlers"
	"github.com/mikestefanello/gcpworkflow/services/pixel/pkg/repo"
)

func main() {
	// Load configuration
	cfg, err := config.GetConfig()
	if err != nil {
		panic(fmt.Sprintf("unable to configuration: %v", err))
	}

	// Initialize web
	e := web.Init(cfg)

	// Pixel repository
	pr, err := repo.NewPubSubRepo(cfg)
	if err != nil {
		panic(fmt.Sprintf("unable to initialize pixel repo: %v", err))
	}

	// Pixel route
	ph := handlers.NewPixelHandler(pr)
	e.GET("/", ph.Get, middleware.CacheControl(0)).Name = "pixel.get"

	// Start the server
	web.Start(cfg, e)
}

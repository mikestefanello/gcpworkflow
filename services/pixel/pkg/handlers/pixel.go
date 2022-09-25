package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/mikestefanello/gcpworkflow/pkg/pixel"
)

// TrackingPixel is the bytes that make up the gif tracking pixel
var TrackingPixel = []byte{
	71, 73, 70, 56, 57, 97, 1, 0, 1, 0, 128, 0, 0, 0, 0, 0,
	255, 255, 255, 33, 249, 4, 1, 0, 0, 0, 0, 44, 0, 0, 0, 0,
	1, 0, 1, 0, 0, 2, 1, 68, 0, 59,
}

type pixelHandler struct {
	repo pixel.Repository
}

func NewPixelHandler(repo pixel.Repository) *pixelHandler {
	return &pixelHandler{
		repo: repo,
	}
}

func (h *pixelHandler) Get(ctx echo.Context) error {
	var p pixel.Pixel

	if err := ctx.Bind(&p); err != nil {
		return err
	}

	defer func() {
		if err := h.repo.Store(ctx.Request().Context(), p); err != nil {
			ctx.Logger().Error(err)
		}
	}()

	return ctx.Blob(http.StatusOK, "image/gif", TrackingPixel)
}

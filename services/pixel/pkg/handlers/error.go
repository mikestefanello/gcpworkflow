package handlers

import (
	"context"
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"
)

type (
	errorHandler struct{}

	errorResponse struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
	}
)

// handle handles all API errors
func (e *errorHandler) handle(err error, ctx echo.Context) {
	// Nothing to handle if the response has been committed or context canceled
	if ctx.Response().Committed || errors.Is(err, context.Canceled) {
		return
	}

	var resp errorResponse

	// If the error is an echo.HTTPError, we can derive the status code and message
	if he, ok := err.(*echo.HTTPError); ok {
		if he.Message != nil {
			resp.Message = he.Message.(string)
		}
		resp.Code = he.Code
	} else {
		// Provide fallback defaults
		resp.Code = http.StatusInternalServerError
		resp.Message = http.StatusText(http.StatusInternalServerError)
	}

	// Log the error
	if resp.Code >= 500 {
		ctx.Logger().Error(err)
	} else {
		ctx.Logger().Info(err)
	}

	if err = ctx.JSON(resp.Code, resp); err != nil {
		ctx.Logger().Error(err)
	}
}

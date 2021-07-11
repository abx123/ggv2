package middleware

import (
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.uber.org/zap"
)

func RequestID() echo.MiddlewareFunc {
	uuid := uuid.New().String()
	logger := zap.L().With(zap.String("rqId", uuid))
	zap.ReplaceGlobals(logger)

	// return middleware.RequestID()
	return middleware.RequestIDWithConfig(middleware.RequestIDConfig{
		Generator: func() string {
			return uuid
		},
	})
}

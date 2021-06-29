package middleware

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func Logger() echo.MiddlewareFunc {
	return middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "id=${id}, method=${method}, uri=${uri}, status=${status}, error=${error}, latency=${latency}ns\n",
	})
}

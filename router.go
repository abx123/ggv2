package main

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo-contrib/prometheus"
	"github.com/labstack/echo/v4"

	// "github.com/labstack/echo/v4/middleware"

	"ggv2/handler"
	"ggv2/handler/middleware"
	// "ggv2/middleware"
)

type router struct {
	Port int
	Conn *sqlx.DB
}

func NewRouter(port int, conn *sqlx.DB) *router {
	return &router{
		Port: port,
		Conn: conn,
	}
}

func (router *router) InitRouter() *echo.Echo {
	// c := InjectController()
	handler := handler.NewHandler(router.Conn)
	r := echo.New()

	// Middleware
	// r.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
	// 	Format: "method=${method}, uri=${uri}, status=${status}, error=${error}, latency=${latency}ns\n",
	// }))
	r.Pre(middleware.RequestID())
	r.Use(middleware.Logger())
	r.Use(middleware.Recover())
	// r.Use(middleware.RequestID())
	// cors := middleware.Cors()
	r.Use(middleware.Cors())
	// r.Use(apmechov4.Middleware())
	//CORS
	// r.Use(middleware.CORSWithConfig(middleware.CORSConfig{
	// 	AllowOrigins: []string{"*"},
	// 	AllowMethods: []string{echo.GET, echo.HEAD, echo.PUT, echo.PATCH, echo.POST, echo.DELETE, echo.OPTIONS},
	// }))

	p := prometheus.NewPrometheus("GGv2", nil)
	p.Use(r)
	// middleware.Cors
	// r.Use(middleware.Cors)

	// Healthcheck
	r.GET("/ping", handler.Ping)

	// // Empty tables
	// r.GET("/empty_tables", (c).EmptyTables)

	// // Tables
	// r.GET("/tables", (c).GetTables)
	r.GET("/table/:id", handler.GetTable)
	r.PUT("/table", handler.CreateTable)

	// // Guest List
	// r.POST("/guest_list/:name", (c).AddToGuestList)
	r.GET("/guest_list", handler.GetGuestList)

	// // Guest Arrives
	// r.PUT("/guests/:name", (c).GuestArrived)

	// // Guest Leaves
	// r.DELETE("/guests/:name", (c).GuestDepart)

	// // List Arrived Guest
	r.GET("/guests", handler.ListArrivedGuest)

	// // Empty Seats
	// r.GET("/seats_empty", (c).GetEmptySeatsCount)

	r.Start(fmt.Sprintf(":%d", router.Port))
	return r
}

package main

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo-contrib/prometheus"
	"github.com/labstack/echo/v4"

	"ggv2/handler"
	"ggv2/handler/middleware"
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
	handler := handler.NewHandler(router.Conn)
	r := echo.New()

	// Middleware
	r.Pre(middleware.RequestID())
	r.Use(middleware.Logger())
	r.Use(middleware.Recover())
	r.Use(middleware.Cors())

	p := prometheus.NewPrometheus("GGv2", nil)
	p.Use(r)

	// Healthcheck
	r.GET("/ping", handler.Ping)

	// // Empty tables
	r.GET("/empty_tables", handler.EmptyTables)

	// // Tables
	r.GET("/tables", handler.GetTables)
	r.GET("/table/:id", handler.GetTable)
	r.PUT("/table", handler.CreateTable)

	// // Guest List
	r.POST("/guest_list/:name", handler.AddToGuestList)
	r.GET("/guest_list", handler.GetGuestList)

	// // Guest Arrives
	r.PUT("/guests/:name", handler.GuestArrived)

	// // Guest Leaves
	r.DELETE("/guests/:name", handler.GuestDepart)

	// // List Arrived Guest
	r.GET("/guests", handler.ListArrivedGuest)

	// // Empty Seats
	r.GET("/seats_empty", handler.GetEmptySeatsCount)

	r.Start(fmt.Sprintf(":%d", router.Port))
	return r
}

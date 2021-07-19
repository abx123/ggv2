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
	th := handler.NewTableHandler(router.Conn)
	gh := handler.NewGuestHandler(router.Conn)
	r := echo.New()

	// Middleware
	r.Use(middleware.Middleware())

	p := prometheus.NewPrometheus("GGv2", nil)
	p.Use(r)

	// Healthcheck
	r.GET("/ping", gh.Ping)

	// // Empty tables
	r.GET("/empty_tables", th.EmptyTables)

	// // Tables
	r.GET("/tables", th.GetTables)
	r.GET("/table/:id", th.GetTable)
	r.PUT("/table", th.CreateTable)

	// // Guest List
	r.POST("/guest_list/:name", gh.AddToGuestList)
	r.GET("/guest_list", gh.GetGuestList)

	// // Guest Arrives
	r.PUT("/guests/:name", gh.GuestArrived)

	// // Guest Leaves
	r.DELETE("/guests/:name", gh.GuestDepart)

	// // List Arrived Guest
	r.GET("/guests", gh.ListArrivedGuest)

	// // Empty Seats
	r.GET("/seats_empty", th.GetEmptySeatsCount)

	r.Start(fmt.Sprintf(":%d", router.Port))
	return r
}

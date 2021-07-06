package services

import (
	"context"
	"ggv2/entities"
)

type DbService interface {
	GetTable(context.Context, int64) (*entities.Table, error)
	ListTables(context.Context, int64, int64) ([]*entities.Table, error)
	CreateTable(context.Context, int64) (*entities.Table, error)
	GetEmptySeatsCount(context.Context) (int, error)
	AddToGuestList(context.Context, int64, int64, string) error
	ListRSVPGuests(context.Context, int64, int64) ([]*entities.Guest, error)
	GuestDepart(context.Context, string) error
	GuestArrival(context.Context, int64, string) error
	ListArrivedGuests(context.Context, int64, int64) ([]*entities.Guest, error)
	EmptyTables(context.Context) error
}

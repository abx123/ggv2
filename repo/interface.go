package repo

import (
	"context"
	"ggv2/entities"
)

type DbRepo interface {
	GetTable(context.Context, int64) (*entities.Table, error)

	CreateTable(context.Context, *entities.Table) (*entities.Table, error)
	ListTables(context.Context, int64, int64) ([]*entities.Table, error)
	EmptyTables(context.Context) error
	GetEmptySeatsCount(context.Context) (int, error)
	GetGuestByName(context.Context, *entities.Guest) (*entities.Guest, error)
	AddToGuestList(context.Context, *entities.Guest) error
	ListGuests(context.Context, int64, int64) ([]*entities.Guest, error)
	GuestArrived(context.Context, *entities.Guest) error
	ListArrivedGuests(context.Context, int64, int64) ([]*entities.Guest, error)
	GuestDepart(context.Context, *entities.Guest) error
}

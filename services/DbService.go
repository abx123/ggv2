package services

import (
	"context"
	"fmt"
	"ggv2/constant"
	"ggv2/entities"
	"ggv2/repo"
)

type DBService struct {
	repo repo.DbRepo
}

func NewDbService(r *repo.DBRepo) *DBService {
	return &DBService{
		repo: r,
	}
}

// GetTable returns detail of a single table.
func (svc *DBService) GetTable(ctx context.Context, id int64) (*entities.Table, error) {
	table, err := svc.repo.GetTable(ctx, id)
	if err != nil {
		// DB operation returns error
		return nil, err
	}
	// All ok, return result
	return table, nil
}

func (svc *DBService) ListTables(ctx context.Context, limit, offset int64) ([]*entities.Table, error) {
	tables, err := svc.repo.ListTables(ctx, limit, offset)
	if err != nil {
		return nil, err
	}
	return tables, nil
}

func (svc *DBService) CreateTable(ctx context.Context, capacity int64) (*entities.Table, error) {
	table := &entities.Table{
		Capacity:          capacity,
		AvailableCapacity: capacity,
		PlannedCapacity:   capacity,
	}
	table, err := svc.repo.CreateTable(ctx, table)
	if err != nil {
		reqId := ctx.Value(constant.ContextKeyRequestID)
		fmt.Println("hre2", err.Error(), ctx)
		fmt.Println("reqId:", reqId)
		return nil, err
	}
	return table, nil
}

func (svc *DBService) AddToGuestList(ctx context.Context, accompanyingGuests, tableID int64, name string) error {
	guest := &entities.Guest{
		Name:        name,
		TotalGuests: accompanyingGuests + 1,
		TableID:     tableID,
	}
	err := svc.repo.AddToGuestList(ctx, guest)

	return err
}

func (svc *DBService) GuestDepart(ctx context.Context, name string) error {
	guest := &entities.Guest{
		Name: name,
	}
	err := svc.repo.GuestDepart(ctx, guest)
	return err
}

func (svc *DBService) GuestArrival(ctx context.Context, accompanyingGuests int64, name string) error {
	guest := &entities.Guest{
		Name:               name,
		TotalArrivedGuests: accompanyingGuests + 1,
	}
	err := svc.repo.GuestArrived(ctx, guest)
	return err
}

func (svc *DBService) ListArrivedGuests(ctx context.Context, limit, offset int64) ([]*entities.Guest, error) {
	guests, err := svc.repo.ListArrivedGuests(ctx, limit, offset)
	if err != nil {
		return nil, err
	}
	return guests, nil
}

func (svc *DBService) ListRSVPGuests(ctx context.Context, limit, offset int64) ([]*entities.Guest, error) {
	guests, err := svc.repo.ListGuests(ctx, limit, offset)
	if err != nil {
		return nil, err
	}
	return guests, nil
}

func (svc *DBService) EmptyTables(ctx context.Context) error {
	err := svc.repo.EmptyTables(ctx)
	return err
}

func (svc *DBService) GetEmptySeatsCount(ctx context.Context) (int, error) {
	count, err := svc.repo.GetEmptySeatsCount(ctx)
	if err != nil {
		return 0, err
	}
	return count, nil
}

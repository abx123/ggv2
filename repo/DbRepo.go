package repo

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/jmoiron/sqlx"

	"go.uber.org/zap"

	"ggv2/constant"
	"ggv2/entities"
)

type ContextKey string

type DBRepo struct {
	db *sqlx.DB
}

const ContextKeyRequestID ContextKey = "requestID"

func NewDbRepo(db *sqlx.DB) *DBRepo {
	return &DBRepo{
		db: db,
	}
}

// GetTable returns detail of a single table.
func (r *DBRepo) GetTable(ctx context.Context, id int64) (*entities.Table, error) {
	reqId := ctx.Value(ContextKeyRequestID)
	logger := zap.L().With(zap.String("rqId", fmt.Sprintf("%v", reqId)))
	table := entities.Table{}
	// Execute Statement
	err := r.db.Get(&table, "SELECT * FROM `table` WHERE tableid=?", id)
	if err != nil {
		logger.Error("db returns error", zap.Error(err))
		if errors.Is(err, sql.ErrNoRows) {
			return nil, constant.ErrTableNotFound
		}
		// Error paring statement result into struct
		return nil, constant.ErrDBErr
	}
	return &table, nil
}

func (r *DBRepo) CreateTable(ctx context.Context, table *entities.Table) (*entities.Table, error) {
	reqId := ctx.Value(constant.ContextKeyRequestID)
	logger := zap.L().With(zap.String("rqId", fmt.Sprintf("%v", reqId)))
	// Execute Statement
	res, err := r.db.Exec("INSERT INTO `table` (capacity, pcapacity, acapacity) VALUES(?, ?, ?)", table.Capacity, table.Capacity, table.Capacity)
	if err != nil {
		logger.Error("db returns error", zap.Error(err))
		// Error paring statement result into struct
		return nil, constant.ErrDBErr
	}
	id, err := res.LastInsertId()
	if err != nil {
		logger.Error("db returns error", zap.Error(err))
		// Error getting ID of newly created record
		return nil, constant.ErrDBErr
	}
	table.TableID = id

	return table, nil
}

func (r *DBRepo) ListTables(ctx context.Context, limit, offset int64) ([]*entities.Table, error) {
	tables := []*entities.Table{}
	reqId := ctx.Value(ContextKeyRequestID)
	logger := zap.L().With(zap.String("rqId", fmt.Sprintf("%v", reqId)))
	err := r.db.Select(&tables, "SELECT * FROM `table` LIMIT ? OFFSET ?", limit, offset)
	if err != nil {
		logger.Error("db returns error", zap.Error(err))
		if errors.Is(err, sql.ErrNoRows) {
			return nil, constant.ErrTableNotFound
		}
		return nil, constant.ErrDBErr
	}
	return tables, nil
}

func (r *DBRepo) EmptyTables(ctx context.Context) error {
	reqId := ctx.Value(ContextKeyRequestID)
	logger := zap.L().With(zap.String("rqId", fmt.Sprintf("%v", reqId)))
	_, err := r.db.Exec("TRUNCATE TABLE `table`;")
	if err != nil {
		logger.Error("db returns error", zap.Error(err))
		return constant.ErrDBErr
	}
	_, err = r.db.Exec("TRUNCATE TABLE `guests`;")
	if err != nil {
		logger.Error("db returns error", zap.Error(err))
		return constant.ErrDBErr
	}
	return nil
}

// GetEmptySeatsCount calculate current total unoccupied seats.
func (r *DBRepo) GetEmptySeatsCount(ctx context.Context) (int, error) {
	c := 0
	reqId := ctx.Value(ContextKeyRequestID)
	logger := zap.L().With(zap.String("rqId", fmt.Sprintf("%v", reqId)))
	// Execute query
	err := r.db.Get(&c, "SELECT SUM(acapacity) FROM `table`")
	if err != nil {
		logger.Error("db returns error", zap.Error(err))
		// Error executing query
		return c, constant.ErrDBErr
	}
	// All ok, return ok
	return c, nil
}

func (r *DBRepo) GetGuestByName(ctx context.Context, g *entities.Guest) (*entities.Guest, error) {
	guest := entities.Guest{}
	reqId := ctx.Value(ContextKeyRequestID)
	logger := zap.L().With(zap.String("rqId", fmt.Sprintf("%v", reqId)))
	// Execute Statement
	err := r.db.Get(&guest, "SELECT * FROM `guests` WHERE name = ?", g.Name)
	if err != nil {
		logger.Error("db returns error", zap.Error(err))
		if errors.Is(err, sql.ErrNoRows) {
			return nil, constant.ErrGuestNotFound
		}
		// Error paring statement result into struct
		return nil, constant.ErrDBErr
	}
	return &guest, nil
}

func (r *DBRepo) ListGuests(ctx context.Context, limit, offset int64) ([]*entities.Guest, error) {
	guests := []*entities.Guest{}
	reqId := ctx.Value(ContextKeyRequestID)
	logger := zap.L().With(zap.String("rqId", fmt.Sprintf("%v", reqId)))

	err := r.db.Select(&guests, "SELECT * FROM `guests` LIMIT ? OFFSET ?", limit, offset)
	if err != nil {
		logger.Error("db returns error", zap.Error(err))
		if errors.Is(err, sql.ErrNoRows) {
			return nil, constant.ErrGuestNotFound
		}
		return nil, constant.ErrDBErr
	}
	return guests, nil
}

func (r *DBRepo) ListArrivedGuests(ctx context.Context, limit, offset int64) ([]*entities.Guest, error) {
	guests := []*entities.Guest{}
	reqId := ctx.Value(ContextKeyRequestID)
	logger := zap.L().With(zap.String("rqId", fmt.Sprintf("%v", reqId)))

	err := r.db.Select(&guests, "SELECT * FROM `guests` WHERE total_arrived_guests > 0 LIMIT ? OFFSET ?", limit, offset)
	if err != nil {
		logger.Error("db returns error", zap.Error(err))
		if errors.Is(err, sql.ErrNoRows) {
			return nil, constant.ErrGuestNotFound
		}
		return nil, constant.ErrDBErr
	}
	return guests, nil
}

func (r *DBRepo) AddToGuestList(ctx context.Context, guest *entities.Guest) error {
	reqId := ctx.Value(ContextKeyRequestID)
	logger := zap.L().With(zap.String("rqId", fmt.Sprintf("%v", reqId)))
	rsvpGuest, err := r.GetGuestByName(ctx, guest)
	if err != nil && err != constant.ErrGuestNotFound {
		logger.Error("db returns error", zap.Error(err))
		// Error getting guest RSVP record, returning error
		return constant.ErrDBErr
	}
	if rsvpGuest != nil {
		return constant.ErrGuestAlreadyRSVP
	}
	table, err := r.GetTable(ctx, guest.TableID)
	if err != nil {
		logger.Error("db returns error", zap.Error(err))
		// Error getting guest RSVP record, returning error
		return constant.ErrDBErr
	}
	// Table capacity less than number of guests
	if table.PlannedCapacity < guest.TotalGuests {
		return constant.ErrTableIsFull
		// return constant.ErrTableIsFull
	}
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		logger.Error("db returns error", zap.Error(err))
		// Error starting transaction
		return constant.ErrDBErr
	}
	_, err = tx.ExecContext(ctx, "INSERT INTO `guests` (total_guests, tableid, name) VALUES(?, ?, ?)", guest.TotalGuests, guest.TableID, guest.Name)
	if err != nil {
		logger.Error("db returns error", zap.Error(err))
		// Error creating RSVP record for guest
		tx.Rollback()
		return constant.ErrDBErr
	}
	// Calculate new capacity
	table.PlannedCapacity -= guest.TotalGuests
	// Update table capacity information
	res, err := tx.ExecContext(ctx, "UPDATE `table` SET pcapacity=?, version = version + 1 WHERE tableid = ? AND version = ?", table.PlannedCapacity, table.TableID, table.Version)
	if err != nil {
		logger.Error("db returns error", zap.Error(err))
		// Error updating table capacity information
		tx.Rollback()
		return constant.ErrDBErr
	}
	c, err := res.RowsAffected()
	if err != nil {
		logger.Error("db returns error", zap.Error(err))
		// Error getting optimistic lock data for table
		tx.Rollback()
		return constant.ErrDBErr
	}
	if c != 1 {
		// Unable to secure optimistic lock for table
		tx.Rollback()
		return constant.ErrFailedOptimisticLock
		// return constant.ErrFailedOptimisticLock
	}

	// All ok, commiting transaction
	err = tx.Commit()
	if err != nil {
		logger.Error("db returns error", zap.Error(err))
		// Error commiting transaction
		return constant.ErrDBErr
	}

	return nil
}

func (r *DBRepo) GuestArrived(ctx context.Context, guest *entities.Guest) error {
	reqId := ctx.Value(ContextKeyRequestID)
	logger := zap.L().With(zap.String("rqId", fmt.Sprintf("%v", reqId)))
	guestArrival, err := r.GetGuestByName(ctx, guest)
	if err != nil {
		if err == constant.ErrGuestNotFound {
			return constant.ErrGuestNeverRSVP
		}
		logger.Error("db returns error", zap.Error(err))
		// Error getting guest arrival record, returning error
		return constant.ErrDBErr
	}
	if guestArrival.TotalArrivedGuests != 0 {
		return constant.ErrGuestAlreadyArrived
		// return constant.ErrGuestAlreadyArrived
	}
	table, err := r.GetTable(ctx, guestArrival.TableID)
	if err != nil {
		logger.Error("db returns error", zap.Error(err))
		// Error getting guest RSVP record, returning error
		return constant.ErrDBErr
	}
	if table.AvailableCapacity < guest.TotalArrivedGuests {
		// Table capacity less than number of guests
		// tx.Rollback()
		return constant.ErrTableIsFull
		// return constant.ErrTableIsFull
	}
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		logger.Error("db returns error", zap.Error(err))
		// Error starting transaction
		return constant.ErrDBErr
	}
	// Able to accomodate guests, checking-in guest
	res, err := tx.ExecContext(ctx, "UPDATE `guests` SET total_arrived_guests=?, version = version + 1, arrivaltime=NOW() WHERE guestid = ? AND version = ?", guest.TotalArrivedGuests, guestArrival.ID, guestArrival.Version)
	if err != nil {
		logger.Error("db returns error", zap.Error(err))
		// Error creating check-in record for guest
		tx.Rollback()
		return constant.ErrDBErr
	}
	c, err := res.RowsAffected()
	if err != nil {
		logger.Error("db returns error", zap.Error(err))
		// Error getting optimistic lock data for table
		tx.Rollback()
		return constant.ErrDBErr
	}
	if c != 1 {
		// Unable to secure optimistic lock for table
		tx.Rollback()
		return constant.ErrFailedOptimisticLock
		// return constant.ErrFailedOptimisticLock
	}
	// Calculate new capacity
	table.AvailableCapacity -= guest.TotalArrivedGuests
	// Update table capacity information
	// res, err := tx.ExecContext(ctx, "UPDATE available_tables SET available_capacity=?, version = version + 1 WHERE id = ? AND version = ?", row.AvailableCapacity, row.TableID, row.Version)
	res, err = tx.ExecContext(ctx, "UPDATE `table` SET acapacity=?, version = version + 1 WHERE tableid = ? AND version = ?", table.AvailableCapacity, table.TableID, table.Version)
	if err != nil {
		logger.Error("db returns error", zap.Error(err))
		// Error updating table capacity information
		tx.Rollback()
		return constant.ErrDBErr
	}
	c, err = res.RowsAffected()
	if err != nil {
		logger.Error("db returns error", zap.Error(err))
		// Error getting optimistic lock data for table
		tx.Rollback()
		return constant.ErrDBErr
	}
	if c != 1 {
		// Unable to secure optimistic lock for table
		tx.Rollback()
		return constant.ErrFailedOptimisticLock
		// return constant.ErrFailedOptimisticLock
	}
	// All ok, commiting transaction
	err = tx.Commit()
	if err != nil {
		logger.Error("db returns error", zap.Error(err))
		// Error commiting transaction
		return constant.ErrDBErr
	}
	return nil
}

func (r *DBRepo) GuestDepart(ctx context.Context, guest *entities.Guest) error {
	reqId := ctx.Value(ContextKeyRequestID)
	logger := zap.L().With(zap.String("rqId", fmt.Sprintf("%v", reqId)))
	guestArrival, err := r.GetGuestByName(ctx, guest)
	if err != nil {
		logger.Error("db returns error", zap.Error(err))
		// Error getting guest arrival record, returning error
		return constant.ErrDBErr
	}
	if guestArrival.TotalArrivedGuests == 0 {
		return constant.ErrGuestNotArrived
	}
	table, err := r.GetTable(ctx, guestArrival.TableID)
	if err != nil {
		logger.Error("db returns error", zap.Error(err))
		// Error getting guest RSVP record, returning error
		return constant.ErrDBErr
	}
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		logger.Error("db returns error", zap.Error(err))
		// Error starting transaction
		return constant.ErrDBErr
	}
	// Able to accomodate guests, checking-in guest
	res, err := tx.ExecContext(ctx, "UPDATE `guests` SET total_arrived_guests=0, version = version + 1, arrivaltime=\"\" WHERE guestid = ? AND version = ?", guestArrival.ID, guestArrival.Version)
	if err != nil {
		logger.Error("db returns error", zap.Error(err))
		// Error creating check-in record for guest
		tx.Rollback()
		return constant.ErrDBErr
	}
	c, err := res.RowsAffected()
	if err != nil {
		logger.Error("db returns error", zap.Error(err))
		// Error getting optimistic lock data for table
		tx.Rollback()
		return constant.ErrDBErr
	}
	if c != 1 {
		// Unable to secure optimistic lock for table
		tx.Rollback()
		return constant.ErrFailedOptimisticLock
		// constant.ErrFailedOptimisticLock
	}
	table.AvailableCapacity += guestArrival.TotalArrivedGuests
	res, err = tx.ExecContext(ctx, "UPDATE `table` SET acapacity=?, version = version + 1 WHERE tableid = ? AND version = ?", table.AvailableCapacity, table.TableID, table.Version)
	if err != nil {
		logger.Error("db returns error", zap.Error(err))
		// Error updating table capacity information
		tx.Rollback()
		return constant.ErrDBErr
	}
	c, err = res.RowsAffected()
	if err != nil {
		logger.Error("db returns error", zap.Error(err))
		// Error getting optimistic lock data for table
		tx.Rollback()
		return constant.ErrDBErr
	}
	if c != 1 {
		// Unable to secure optimistic lock for table
		tx.Rollback()
		return constant.ErrFailedOptimisticLock
	}
	// All ok, commiting transaction
	err = tx.Commit()
	if err != nil {
		logger.Error("db returns error", zap.Error(err))
		// Error commiting transaction
		return constant.ErrDBErr
	}
	return nil
}

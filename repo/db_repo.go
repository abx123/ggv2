package repo

import (
	"context"
	"database/sql"
	"errors"

	"github.com/jmoiron/sqlx"

	"go.uber.org/zap"

	"ggv2/entities"
)

type DBRepo struct {
	db *sqlx.DB
}

var (
	errTableNotFound = errors.New("table not found")

	errGuestNotFound = errors.New("guest not found")

	errGuestAlreadyRSVP = errors.New("guest already RSVP")

	errGuestAlreadyArrived = errors.New("guest already arrived")

	errTableIsFull = errors.New("table is full")

	errFailedOptimisticLock = errors.New("unable to secure optimistic lock, please retry")

	errDBErr = errors.New("database returns error")

	errGuestNeverRSVP = errors.New("guest never rsvp")

	errGuestNotArrived = errors.New("guest not arrived")
)

func NewDbRepo(db *sqlx.DB) *DBRepo {
	return &DBRepo{
		db: db,
	}
}

// GetTable returns detail of a single table.
func (r *DBRepo) GetTable(ctx context.Context, id int64) (*entities.Table, error) {
	table := entities.Table{}
	// Execute Statement
	err := r.db.Get(&table, "SELECT * FROM `table` WHERE id=?", id)
	if err != nil {
		zap.L().Error("db returns error", zap.Error(err))
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errTableNotFound
		}
		// Error paring statement result into struct
		return nil, errDBErr
	}
	return &table, nil
}

func (r *DBRepo) CreateTable(ctx context.Context, table *entities.Table) (*entities.Table, error) {
	// Execute Statement
	res, err := r.db.Exec("INSERT INTO `table` (capacity, pcapacity, acapacity) VALUES(?, ?, ?)", table.Capacity, table.Capacity, table.Capacity)
	if err != nil {
		zap.L().Error("db returns error", zap.Error(err))
		// Error paring statement result into struct
		return nil, errDBErr
	}
	id, err := res.LastInsertId()
	if err != nil {
		zap.L().Error("db returns error", zap.Error(err))
		// Error getting ID of newly created record
		return nil, errDBErr
	}
	table.TableID = id

	return table, nil
}

func (r *DBRepo) ListTables(ctx context.Context, limit, offset int64) ([]*entities.Table, error) {
	tables := []*entities.Table{}
	err := r.db.Select(&tables, "SELECT * FROM `table` LIMIT ? OFFSET ?", limit, offset)
	if err != nil {
		zap.L().Error("db returns error", zap.Error(err))
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errTableNotFound
		}
		return nil, errDBErr
	}
	return tables, nil
}

func (r *DBRepo) EmptyTables(ctx context.Context) error {
	_, err := r.db.Exec("TRUNCATE TABLE `table`;")
	if err != nil {
		zap.L().Error("db returns error", zap.Error(err))
		return errDBErr
	}
	_, err = r.db.Exec("TRUNCATE TABLE `guests`;")
	if err != nil {
		zap.L().Error("db returns error", zap.Error(err))
		return errDBErr
	}
	return nil
}

// GetEmptySeatsCount calculate current total unoccupied seats.
func (r *DBRepo) GetEmptySeatsCount(ctx context.Context) (int, error) {
	c := 0
	// Execute query
	err := r.db.Get(&c, "SELECT SUM(acapacity) FROM `table`")
	if err != nil {
		zap.L().Error("db returns error", zap.Error(err))
		// Error executing query
		return c, errDBErr
	}
	// All ok, return ok
	return c, nil
}

func (r *DBRepo) GetGuestByName(ctx context.Context, g *entities.Guest) (*entities.Guest, error) {
	guest := entities.Guest{}
	// Execute Statement
	err := r.db.Get(&guest, "SELECT * FROM `guests` WHERE name = ?", g.Name)
	if err != nil {
		zap.L().Error("db returns error", zap.Error(err))
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errGuestNotFound
		}
		// Error paring statement result into struct
		return nil, errDBErr
	}
	return &guest, nil
}

func (r *DBRepo) ListGuests(ctx context.Context, limit, offset int64) ([]*entities.Guest, error) {
	guests := []*entities.Guest{}

	err := r.db.Select(&guests, "SELECT * FROM `guests` LIMIT ? OFFSET ?", limit, offset)
	if err != nil {
		zap.L().Error("db returns error", zap.Error(err))
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errGuestNotFound
		}
		return nil, errDBErr
	}
	return guests, nil
}

func (r *DBRepo) ListArrivedGuests(ctx context.Context, limit, offset int64) ([]*entities.Guest, error) {
	guests := []*entities.Guest{}

	err := r.db.Select(&guests, "SELECT * FROM `guests` WHERE total_arrived_guests > 0 LIMIT ? OFFSET ?", limit, offset)
	if err != nil {
		zap.L().Error("db returns error", zap.Error(err))
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errGuestNotFound
		}
		return nil, errDBErr
	}
	return guests, nil
}

func (r *DBRepo) AddToGuestList(ctx context.Context, guest *entities.Guest) error {
	rsvpGuest, err := r.GetGuestByName(ctx, guest)
	if err != nil && err != errGuestNotFound {
		zap.L().Error("db returns error", zap.Error(err))
		// Error getting guest RSVP record, returning error
		return errDBErr
	}
	if rsvpGuest != nil {
		return errGuestAlreadyRSVP
	}
	table, err := r.GetTable(ctx, guest.TableID)
	if err != nil {
		zap.L().Error("db returns error", zap.Error(err))
		// Error getting guest RSVP record, returning error
		return errDBErr
	}
	// Table capacity less than number of guests
	if table.PlannedCapacity < guest.TotalGuests {
		return errTableIsFull
	}
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		zap.L().Error("db returns error", zap.Error(err))
		// Error starting transaction
		return errDBErr
	}
	insert, err := tx.ExecContext(ctx, "INSERT INTO `guests` (total_rsvp_guests, tableid, name) VALUES(?, ?, ?)", guest.TotalGuests, guest.TableID, guest.Name)
	if err != nil {
		zap.L().Error("db returns error", zap.Error(err))
		// Error creating RSVP record for guest
		tx.Rollback()
		return errDBErr
	}
	_, err = insert.LastInsertId()
	if err != nil {
		zap.L().Error("db returns error", zap.Error(err))
		// Error getting ID of newly created record
		return errDBErr
	}
	// Calculate new capacity
	table.PlannedCapacity -= guest.TotalGuests
	// Update table capacity information
	res, err := tx.ExecContext(ctx, "UPDATE `table` SET pcapacity=?, version = version + 1 WHERE id = ? AND version = ?", table.PlannedCapacity, table.TableID, table.Version)
	if err != nil {
		zap.L().Error("db returns error", zap.Error(err))
		// Error updating table capacity information
		tx.Rollback()
		return errDBErr
	}
	c, err := res.RowsAffected()
	if err != nil {
		zap.L().Error("db returns error", zap.Error(err))
		// Error getting optimistic lock data for table
		tx.Rollback()
		return errDBErr
	}
	if c != 1 {
		// Unable to secure optimistic lock for table
		tx.Rollback()
		return errFailedOptimisticLock
	}

	// All ok, commiting transaction
	err = tx.Commit()
	if err != nil {
		zap.L().Error("db returns error", zap.Error(err))
		// Error commiting transaction
		return errDBErr
	}

	return nil
}

func (r *DBRepo) GuestArrived(ctx context.Context, guest *entities.Guest) error {
	guestArrival, err := r.GetGuestByName(ctx, guest)
	if err != nil {
		if err == errGuestNotFound {
			return errGuestNeverRSVP
		}
		zap.L().Error("db returns error", zap.Error(err))
		// Error getting guest arrival record, returning error
		return errDBErr
	}
	if guestArrival.TotalArrivedGuests != 0 {
		return errGuestAlreadyArrived
	}
	table, err := r.GetTable(ctx, guestArrival.TableID)
	if err != nil {
		zap.L().Error("db returns error", zap.Error(err))
		// Error getting guest RSVP record, returning error
		return errDBErr
	}
	if table.AvailableCapacity < guest.TotalArrivedGuests {
		// Table capacity less than number of guests
		return errTableIsFull
	}
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		zap.L().Error("db returns error", zap.Error(err))
		// Error starting transaction
		return errDBErr
	}
	// Able to accomodate guests, checking-in guest
	res, err := tx.ExecContext(ctx, "UPDATE `guests` SET total_arrived_guests=?, version = version + 1, arrivaltime=NOW() WHERE id = ? AND version = ?", guest.TotalArrivedGuests, guestArrival.ID, guestArrival.Version)
	if err != nil {
		zap.L().Error("db returns error", zap.Error(err))
		tx.Rollback()
		return errDBErr
	}
	c, err := res.RowsAffected()
	if err != nil {
		zap.L().Error("db returns error", zap.Error(err))
		// Error getting optimistic lock data for table
		tx.Rollback()
		return errDBErr
	}
	if c != 1 {
		// Unable to secure optimistic lock for table
		tx.Rollback()
		return errFailedOptimisticLock
	}
	// Calculate new capacity
	table.AvailableCapacity -= guest.TotalArrivedGuests
	// Update table capacity information
	res, err = tx.ExecContext(ctx, "UPDATE `table` SET acapacity=?, version = version + 1 WHERE id = ? AND version = ?", table.AvailableCapacity, table.TableID, table.Version)
	if err != nil {
		zap.L().Error("db returns error", zap.Error(err))
		// Error updating table capacity information
		tx.Rollback()
		return errDBErr
	}
	c, err = res.RowsAffected()
	if err != nil {
		zap.L().Error("db returns error", zap.Error(err))
		// Error getting optimistic lock data for table
		tx.Rollback()
		return errDBErr
	}
	if c != 1 {
		// Unable to secure optimistic lock for table
		tx.Rollback()
		return errFailedOptimisticLock
	}
	// All ok, commiting transaction
	err = tx.Commit()
	if err != nil {
		zap.L().Error("db returns error", zap.Error(err))
		// Error commiting transaction
		return errDBErr
	}
	return nil
}

func (r *DBRepo) GuestDepart(ctx context.Context, guest *entities.Guest) error {
	guestArrival, err := r.GetGuestByName(ctx, guest)
	if err != nil {
		zap.L().Error("db returns error", zap.Error(err))
		// Error getting guest arrival record, returning error
		return errDBErr
	}
	if guestArrival.TotalArrivedGuests == 0 {
		return errGuestNotArrived
	}
	table, err := r.GetTable(ctx, guestArrival.TableID)
	if err != nil {
		zap.L().Error("db returns error", zap.Error(err))
		// Error getting guest RSVP record, returning error
		return errDBErr
	}
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		zap.L().Error("db returns error", zap.Error(err))
		// Error starting transaction
		return errDBErr
	}
	// Able to accomodate guests, checking-in guest
	res, err := tx.ExecContext(ctx, "UPDATE `guests` SET total_arrived_guests=0, version = version + 1, arrivaltime=\"\" WHERE id = ? AND version = ?", guestArrival.ID, guestArrival.Version)
	if err != nil {
		zap.L().Error("db returns error", zap.Error(err))
		// Error creating check-in record for guest
		tx.Rollback()
		return errDBErr
	}
	c, err := res.RowsAffected()
	if err != nil {
		zap.L().Error("db returns error", zap.Error(err))
		// Error getting optimistic lock data for table
		tx.Rollback()
		return errDBErr
	}
	if c != 1 {
		// Unable to secure optimistic lock for table
		tx.Rollback()
		return errFailedOptimisticLock
	}
	table.AvailableCapacity += guestArrival.TotalArrivedGuests
	res, err = tx.ExecContext(ctx, "UPDATE `table` SET acapacity=?, version = version + 1 WHERE id = ? AND version = ?", table.AvailableCapacity, table.TableID, table.Version)
	if err != nil {
		zap.L().Error("db returns error", zap.Error(err))
		// Error updating table capacity information
		tx.Rollback()
		return errDBErr
	}
	c, err = res.RowsAffected()
	if err != nil {
		zap.L().Error("db returns error", zap.Error(err))
		// Error getting optimistic lock data for table
		tx.Rollback()
		return errDBErr
	}
	if c != 1 {
		// Unable to secure optimistic lock for table
		tx.Rollback()
		return errFailedOptimisticLock
	}
	// All ok, commiting transaction
	err = tx.Commit()
	if err != nil {
		zap.L().Error("db returns error", zap.Error(err))
		// Error commiting transaction
		return errDBErr
	}
	return nil
}

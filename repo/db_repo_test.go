package repo

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"regexp"
	"testing"

	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	sqlxmock "github.com/zhashkevych/go-sqlxmock"

	"ggv2/entities"
)

func NewMockDb() (*sqlx.DB, sqlxmock.Sqlmock) {
	db, mock, err := sqlxmock.Newx()
	if err != nil {
		log.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	return db, mock
}

func TestCreateTable(t *testing.T) {
	query := regexp.QuoteMeta("INSERT INTO `table` (capacity, pcapacity, acapacity) VALUES(?, ?, ?)")
	type TestCase struct {
		name          string
		desc          string
		input         *entities.Table
		err           error
		dbErr         bool
		expRes        *entities.Table
		expErr        error
		lastInsertErr bool
	}
	testcases := []TestCase{
		{
			name: "Happy case",
			desc: "Db return record",
			input: &entities.Table{
				Capacity: int64(7),
			},
			expRes: &entities.Table{
				TableID:  99,
				Capacity: int64(7),
				Version:  0,
			},
		},
		{
			name: "Sad case",
			desc: "Db return error",
			input: &entities.Table{
				Capacity: int64(7),
			},
			err:    fmt.Errorf("mock error"),
			expErr: errDBErr,
			dbErr:  true,
		},
		{
			name: "Sad case",
			desc: "LastInsertId return error",
			input: &entities.Table{
				Capacity: int64(7),
			},
			err:           fmt.Errorf("LastInsertId error"),
			expErr:        errDBErr,
			lastInsertErr: true,
		},
	}

	for _, v := range testcases {
		db, mock := NewMockDb()
		repo := NewDbRepo(db)
		if v.dbErr {
			mock.ExpectExec(query).WillReturnError(v.err)
		}
		if v.lastInsertErr {
			mock.ExpectExec(query).WillReturnResult(sqlxmock.NewErrorResult(v.err))
		}
		mock.ExpectExec(query).WillReturnResult(sqlxmock.NewResult(99, 1))
		actRes, actErr := repo.CreateTable(context.Background(), v.input)
		assert.Equal(t, v.expErr, actErr)
		assert.Equal(t, v.expRes, actRes)
	}
}

func TestGetEmptySeatsCount(t *testing.T) {
	query := regexp.QuoteMeta("SELECT SUM(acapacity) FROM `table`")
	rows := sqlxmock.NewRows([]string{"SUM(acapacity)"}).AddRow(99)
	type TestCase struct {
		name   string
		desc   string
		err    error
		expRes int
		expErr error
	}
	testcases := []TestCase{
		{
			name:   "Happy case",
			desc:   "Db return record",
			expRes: 99,
		},
		{
			name:   "Sad case",
			desc:   "Db return error",
			err:    fmt.Errorf("mock error"),
			expRes: 0,
			expErr: errDBErr,
		},
	}
	for _, v := range testcases {
		db, mock := NewMockDb()
		repo := NewDbRepo(db)
		if v.err != nil {
			mock.ExpectQuery(query).WillReturnError(v.err)
		} else {
			mock.ExpectQuery(query).WillReturnRows(rows)
		}
		actRes, actErr := repo.GetEmptySeatsCount(context.Background())
		assert.Equal(t, v.expErr, actErr)
		assert.Equal(t, v.expRes, actRes)
	}
}

func TestListTables(t *testing.T) {
	query := regexp.QuoteMeta("SELECT * FROM `table` LIMIT ? OFFSET ?")
	rows := sqlxmock.NewRows([]string{"id", "capacity", "acapacity", "pcapacity", "version"}).AddRow(1, 6, 6, 6, 0).AddRow(2, 10, 8, 6, 3)
	type TestCase struct {
		name   string
		desc   string
		err    error
		dbErr  bool
		expRes []*entities.Table
		expErr error
	}
	testcases := []TestCase{
		{
			name: "Happy case",
			desc: "Db return record",
			expRes: []*entities.Table{
				{
					TableID:           1,
					Capacity:          6,
					AvailableCapacity: 6,
					PlannedCapacity:   6,
					Version:           0,
				},
				{
					TableID:           2,
					Capacity:          10,
					AvailableCapacity: 8,
					PlannedCapacity:   6,
					Version:           3,
				},
			},
		},
		{
			name:   "Sad case",
			desc:   "Db return error",
			err:    fmt.Errorf("mock error"),
			dbErr:  true,
			expErr: errDBErr,
		},
		{
			name:   "Sad case",
			desc:   "no table found",
			dbErr:  true,
			err:    sql.ErrNoRows,
			expErr: errTableNotFound,
		},
	}
	for _, v := range testcases {
		db, mock := NewMockDb()
		repo := NewDbRepo(db)
		if v.dbErr {
			mock.ExpectQuery(query).WillReturnError(v.err)
		}
		mock.ExpectQuery(query).WillReturnRows(rows)
		actRes, actErr := repo.ListTables(context.Background(), 10, 0)
		assert.Equal(t, v.expErr, actErr)
		assert.Equal(t, v.expRes, actRes)
	}
}

func TestEmptyTables(t *testing.T) {
	type TestCase struct {
		name             string
		desc             string
		err              error
		expErr           error
		truncateTableErr bool
		truncateGuestErr bool
	}
	testcases := []TestCase{
		{
			name: "Happy case",
			desc: "tables truncated",
		},
		{
			name:             "Sad case",
			desc:             "Truncate `table` return error",
			err:              fmt.Errorf("mock error"),
			expErr:           errDBErr,
			truncateTableErr: true,
		},
		{
			name:             "Sad case",
			desc:             "Truncate `guests` return error",
			err:              fmt.Errorf("mock error"),
			expErr:           errDBErr,
			truncateGuestErr: true,
		},
	}
	for _, v := range testcases {
		db, mock := NewMockDb()
		repo := NewDbRepo(db)
		if v.truncateTableErr {
			mock.ExpectExec(regexp.QuoteMeta("TRUNCATE TABLE `table`;")).WillReturnError(v.err)
		}
		mock.ExpectExec(regexp.QuoteMeta("TRUNCATE TABLE `table`;")).WillReturnResult(sqlxmock.NewResult(1, 1))
		if v.truncateGuestErr {
			mock.ExpectExec(regexp.QuoteMeta("TRUNCATE TABLE `guests`;")).WillReturnError(v.err)
		}
		mock.ExpectExec(regexp.QuoteMeta("TRUNCATE TABLE `guests`;")).WillReturnResult(sqlxmock.NewResult(1, 1))

		actErr := repo.EmptyTables(context.Background())
		assert.Equal(t, v.expErr, actErr)
	}
}

func TestGetGuestByName(t *testing.T) {
	query := regexp.QuoteMeta("SELECT * FROM `guests` WHERE name = ?")
	rows := sqlxmock.NewRows([]string{"id", "name", "total_rsvp_guests", "total_arrived_guests", "tableid"}).AddRow(1, "dummy", 2, 4, 3)
	type TestCase struct {
		name   string
		desc   string
		err    error
		expErr error
		expRes *entities.Guest
	}
	testcases := []TestCase{
		{
			name: "Happy case",
			desc: "Db return record",
			expRes: &entities.Guest{
				ID:                 1,
				Name:               "dummy",
				TotalGuests:        2,
				TableID:            3,
				TotalArrivedGuests: 4,
			},
		},
		{
			name:   "Sad case",
			desc:   "Db return error",
			err:    fmt.Errorf("mock error"),
			expErr: errDBErr,
		},
		{
			name:   "Sad case",
			desc:   "guest not found",
			err:    sql.ErrNoRows,
			expErr: errGuestNotFound,
		},
	}
	for _, v := range testcases {
		db, mock := NewMockDb()
		repo := NewDbRepo(db)
		if v.err != nil {
			mock.ExpectQuery(query).WillReturnError(v.err)
		} else {
			mock.ExpectQuery(query).WillReturnRows(rows)
		}
		actRes, actErr := repo.GetGuestByName(context.Background(), &entities.Guest{Name: "dummy"})
		assert.Equal(t, v.expErr, actErr)
		assert.Equal(t, v.expRes, actRes)
	}
}

func TestGetTable(t *testing.T) {
	query := regexp.QuoteMeta("SELECT * FROM `table` WHERE id=?")
	rows := sqlxmock.NewRows([]string{"id", "capacity", "acapacity", "pcapacity"}).AddRow(1, 2, 3, 4)
	type TestCase struct {
		name   string
		desc   string
		err    error
		expRes *entities.Table
		expErr error
	}
	testcases := []TestCase{
		{
			name: "Happy case",
			desc: "Db return record",
			expRes: &entities.Table{
				TableID:           1,
				Capacity:          2,
				AvailableCapacity: 3,
				PlannedCapacity:   4,
			},
		},
		{
			name:   "Sad case",
			desc:   "Db return error",
			err:    fmt.Errorf("mock error"),
			expErr: errDBErr,
		},
		{
			name:   "Sad case",
			desc:   "table not found",
			err:    sql.ErrNoRows,
			expErr: errTableNotFound,
		},
	}
	for _, v := range testcases {
		db, mock := NewMockDb()
		repo := NewDbRepo(db)
		if v.err != nil {
			mock.ExpectQuery(query).WillReturnError(v.err)
		} else {
			mock.ExpectQuery(query).WillReturnRows(rows)
		}
		actRes, actErr := repo.GetTable(context.Background(), 1)
		assert.Equal(t, v.expErr, actErr)
		assert.Equal(t, v.expRes, actRes)
	}
}

func TestListArrivedGuest(t *testing.T) {
	query := regexp.QuoteMeta("SELECT * FROM `guests` WHERE total_arrived_guests > 0 LIMIT ? OFFSET ?")
	row := sqlxmock.NewRows([]string{"id", "name", "total_rsvp_guests", "total_arrived_guests", "version", "arrivaltime", "tableid"}).AddRow(1, "dummy", 2, 3, 4, "2021-06-04 04:06:44", 5)
	type TestCase struct {
		name   string
		desc   string
		err    error
		expRes []*entities.Guest
		expErr error
	}
	testcases := []TestCase{
		{
			name: "Happy case",
			desc: "Db return record",
			expRes: []*entities.Guest{
				{
					ID:                 1,
					Name:               "dummy",
					TotalGuests:        2,
					TotalArrivedGuests: 3,
					Version:            4,
					ArrivalTime:        "2021-06-04 04:06:44",
					TableID:            5,
				},
			},
		},
		{
			name:   "Sad case",
			desc:   "Db return error",
			err:    fmt.Errorf("mock error"),
			expErr: errDBErr,
		},
		{
			name:   "Sad case",
			desc:   "no guest found",
			err:    sql.ErrNoRows,
			expErr: errGuestNotFound,
		},
	}
	for _, v := range testcases {
		db, mock := NewMockDb()
		repo := NewDbRepo(db)
		if v.err != nil {
			mock.ExpectQuery(query).WillReturnError(v.err)
		} else {
			mock.ExpectQuery(query).WillReturnRows(row)
		}
		actRes, actErr := repo.ListArrivedGuests(context.Background(), 10, 0)
		assert.Equal(t, v.expErr, actErr)
		assert.Equal(t, v.expRes, actRes)
	}
}

func TestListGuests(t *testing.T) {
	query := regexp.QuoteMeta("SELECT * FROM `guests` LIMIT ? OFFSET ?")
	row := sqlxmock.NewRows([]string{"id", "name", "total_rsvp_guests", "total_arrived_guests", "version", "arrivaltime", "tableid"}).AddRow(1, "dummy", 2, 3, 4, "2021-06-04 04:06:44", 5)
	type TestCase struct {
		name   string
		desc   string
		err    error
		expRes []*entities.Guest
		expErr error
	}
	testcases := []TestCase{
		{
			name: "Happy case",
			desc: "Db return record",
			expRes: []*entities.Guest{
				{
					ID:                 1,
					Name:               "dummy",
					TotalGuests:        2,
					TotalArrivedGuests: 3,
					Version:            4,
					ArrivalTime:        "2021-06-04 04:06:44",
					TableID:            5,
				},
			},
		},
		{
			name:   "Sad case",
			desc:   "Db return error",
			err:    fmt.Errorf("mock error"),
			expErr: errDBErr,
		},
		{
			name:   "Sad case",
			desc:   "guest not found",
			err:    sql.ErrNoRows,
			expErr: errGuestNotFound,
		},
	}
	for _, v := range testcases {
		db, mock := NewMockDb()
		repo := NewDbRepo(db)
		if v.err != nil {
			mock.ExpectQuery(query).WillReturnError(v.err)
		} else {
			mock.ExpectQuery(query).WillReturnRows(row)
		}
		actRes, actErr := repo.ListGuests(context.TODO(), 10, 0)
		assert.Equal(t, v.expErr, actErr)
		assert.Equal(t, v.expRes, actRes)
	}
}

func TestGuestArrived(t *testing.T) {
	getTableQuery := regexp.QuoteMeta("SELECT * FROM `table` WHERE id=?")
	updateGuestQuery := regexp.QuoteMeta("UPDATE `guests` SET total_arrived_guests=?, version = version + 1, arrivaltime=NOW() WHERE id = ? AND version = ?")
	updateTableQuery := regexp.QuoteMeta("UPDATE `table` SET acapacity=?, version = version + 1 WHERE id = ? AND version = ?")
	getGuestByNameQuery := regexp.QuoteMeta("SELECT * FROM `guests` WHERE name = ?")
	type TestCase struct {
		name                         string
		desc                         string
		err                          error
		getTableRows                 *sqlxmock.Rows
		getGuestByNameRows           *sqlxmock.Rows
		getGuestByNameErr            bool
		getTableErr                  bool
		beginTxErr                   bool
		updateGuestErr               bool
		updateGuestRowsAffectedErr   bool
		updateGuestOptimisticLockErr bool
		updateTableErr               bool
		updateTableRowsAffectedErr   bool
		updateTableOptimisticLockErr bool
		commitErr                    bool
		expErr                       error
	}
	testcases := []TestCase{
		{
			name:               "Happy case",
			desc:               "all ok, no error",
			getGuestByNameRows: sqlxmock.NewRows([]string{"id", "name", "total_rsvp_guests", "total_arrived_guests", "version", "arrivaltime", "tableid"}).AddRow(1, "dummy", 2, 0, 4, "", 5),
			getTableRows:       sqlxmock.NewRows([]string{"id", "capacity", "acapacity", "pcapacity", "version"}).AddRow(1, 6, 6, 3, 2),
		},
		{
			name:               "Sad case",
			desc:               "commit error",
			err:                fmt.Errorf("dummy error"),
			commitErr:          true,
			getGuestByNameRows: sqlxmock.NewRows([]string{"id", "name", "total_rsvp_guests", "total_arrived_guests", "version", "arrivaltime", "tableid"}).AddRow(1, "dummy", 2, 0, 4, "", 5),
			getTableRows:       sqlxmock.NewRows([]string{"id", "capacity", "acapacity", "pcapacity", "version"}).AddRow(1, 6, 6, 3, 2),
			expErr:             errDBErr,
		},
		{
			name:               "Sad case",
			desc:               "update table error",
			err:                fmt.Errorf("dummy error"),
			updateTableErr:     true,
			getGuestByNameRows: sqlxmock.NewRows([]string{"id", "name", "total_rsvp_guests", "total_arrived_guests", "version", "arrivaltime", "tableid"}).AddRow(1, "dummy", 2, 0, 4, "", 5),
			getTableRows:       sqlxmock.NewRows([]string{"id", "capacity", "acapacity", "pcapacity", "version"}).AddRow(1, 6, 6, 3, 2),
			expErr:             errDBErr,
		},
		{
			name:                       "Sad case",
			desc:                       "update table rows affected error",
			err:                        fmt.Errorf("rows affected err"),
			updateTableRowsAffectedErr: true,
			getGuestByNameRows:         sqlxmock.NewRows([]string{"id", "name", "total_rsvp_guests", "total_arrived_guests", "version", "arrivaltime", "tableid"}).AddRow(1, "dummy", 2, 0, 4, "", 5),
			getTableRows:               sqlxmock.NewRows([]string{"id", "capacity", "acapacity", "pcapacity", "version"}).AddRow(1, 6, 6, 3, 2),
			expErr:                     errDBErr,
		},
		{
			name:                         "Sad case",
			desc:                         "update table optimistic lock error",
			updateTableOptimisticLockErr: true,
			getGuestByNameRows:           sqlxmock.NewRows([]string{"id", "name", "total_rsvp_guests", "total_arrived_guests", "version", "arrivaltime", "tableid"}).AddRow(1, "dummy", 2, 0, 4, "", 5),
			getTableRows:                 sqlxmock.NewRows([]string{"id", "capacity", "acapacity", "pcapacity", "version"}).AddRow(1, 6, 6, 3, 2),
			expErr:                       errFailedOptimisticLock,
		},
		{
			name:               "Sad case",
			desc:               "update guest error",
			err:                fmt.Errorf("dummy error"),
			updateGuestErr:     true,
			getGuestByNameRows: sqlxmock.NewRows([]string{"id", "name", "total_rsvp_guests", "total_arrived_guests", "version", "arrivaltime", "tableid"}).AddRow(1, "dummy", 2, 0, 4, "", 5),
			getTableRows:       sqlxmock.NewRows([]string{"id", "capacity", "acapacity", "pcapacity", "version"}).AddRow(1, 6, 6, 3, 2),
			expErr:             errDBErr,
		},
		{
			name:                       "Sad case",
			desc:                       "update guest rows affected error",
			err:                        fmt.Errorf("rows affected err"),
			updateGuestRowsAffectedErr: true,
			getGuestByNameRows:         sqlxmock.NewRows([]string{"id", "name", "total_rsvp_guests", "total_arrived_guests", "version", "arrivaltime", "tableid"}).AddRow(1, "dummy", 2, 0, 4, "", 5),
			getTableRows:               sqlxmock.NewRows([]string{"id", "capacity", "acapacity", "pcapacity", "version"}).AddRow(1, 6, 6, 3, 2),
			expErr:                     errDBErr,
		},
		{
			name:                         "Sad case",
			desc:                         "update guest optimistic lock error",
			updateGuestOptimisticLockErr: true,
			getGuestByNameRows:           sqlxmock.NewRows([]string{"id", "name", "total_rsvp_guests", "total_arrived_guests", "version", "arrivaltime", "tableid"}).AddRow(1, "dummy", 2, 0, 4, "", 5),
			getTableRows:                 sqlxmock.NewRows([]string{"id", "capacity", "acapacity", "pcapacity", "version"}).AddRow(1, 6, 6, 3, 2),
			expErr:                       errFailedOptimisticLock,
		},
		{
			name:               "Sad case",
			desc:               "begin tx error",
			err:                fmt.Errorf("tx error"),
			beginTxErr:         true,
			getGuestByNameRows: sqlxmock.NewRows([]string{"id", "name", "total_rsvp_guests", "total_arrived_guests", "version", "arrivaltime", "tableid"}).AddRow(1, "dummy", 2, 0, 4, "", 5),
			getTableRows:       sqlxmock.NewRows([]string{"id", "capacity", "acapacity", "pcapacity", "version"}).AddRow(1, 6, 6, 3, 2),
			expErr:             errDBErr,
		},
		{
			name:               "Sad case",
			desc:               "get table returns db error",
			err:                sql.ErrNoRows,
			getGuestByNameRows: sqlxmock.NewRows([]string{"id", "name", "total_rsvp_guests", "total_arrived_guests", "version", "arrivaltime", "tableid"}).AddRow(1, "dummy", 2, 0, 4, "", 5),
			getTableRows:       sqlxmock.NewRows([]string{"id", "capacity", "acapacity", "pcapacity", "version"}).AddRow(1, 6, 6, 3, 2),
			getTableErr:        true,
			expErr:             errDBErr,
		},
		{
			name:               "Sad case",
			desc:               "table cannot accomodate guest",
			getGuestByNameRows: sqlxmock.NewRows([]string{"id", "name", "total_rsvp_guests", "total_arrived_guests", "version", "arrivaltime", "tableid"}).AddRow(1, "dummy", 2, 0, 4, "", 5),
			getTableRows:       sqlxmock.NewRows([]string{"id", "capacity", "acapacity", "pcapacity", "version"}).AddRow(1, 6, 0, 3, 2),
			expErr:             errTableIsFull,
		},
		{
			name:               "Sad case",
			desc:               "guest already arrived",
			err:                nil,
			getGuestByNameRows: sqlxmock.NewRows([]string{"id", "name", "total_rsvp_guests", "total_arrived_guests", "version", "arrivaltime", "tableid"}).AddRow(1, "dummy", 2, 3, 4, "2021-06-04 04:06:44", 5),
			getTableRows:       sqlxmock.NewRows([]string{}),
			expErr:             errGuestAlreadyArrived,
		},
		{
			name:               "Sad case",
			desc:               "guest not found",
			err:                sql.ErrNoRows,
			getGuestByNameRows: sqlxmock.NewRows([]string{}),
			getGuestByNameErr:  true,
			getTableRows:       sqlxmock.NewRows([]string{}),
			expErr:             errGuestNeverRSVP,
		},
		{
			name:               "Sad case",
			desc:               "get guest return db error",
			err:                fmt.Errorf("dummy error"),
			getGuestByNameRows: sqlxmock.NewRows([]string{}),
			getGuestByNameErr:  true,
			getTableRows:       sqlxmock.NewRows([]string{}),
			expErr:             errDBErr,
		},
	}
	for _, v := range testcases {
		db, mock := NewMockDb()
		repo := NewDbRepo(db)
		if v.getGuestByNameErr {
			mock.ExpectQuery(getGuestByNameQuery).WillReturnError(v.err)
		}
		mock.ExpectQuery(getGuestByNameQuery).WillReturnRows(v.getGuestByNameRows)
		if v.getTableErr {
			mock.ExpectQuery(getTableQuery).WillReturnError(v.err)
		}
		mock.ExpectQuery(getTableQuery).WillReturnRows(v.getTableRows)
		if v.beginTxErr {
			mock.ExpectBegin().WillReturnError(v.err)
		}
		mock.ExpectBegin()

		if v.updateGuestErr {
			mock.ExpectExec(updateGuestQuery).WillReturnError(v.err)
			mock.ExpectRollback()
		}
		if v.updateGuestRowsAffectedErr {
			mock.ExpectExec(updateGuestQuery).WillReturnResult(sqlxmock.NewErrorResult(v.err))
			mock.ExpectRollback()
		}
		if v.updateGuestOptimisticLockErr {
			mock.ExpectExec(updateGuestQuery).WillReturnResult(sqlxmock.NewResult(0, 0))
			mock.ExpectRollback()
		}
		mock.ExpectExec(updateGuestQuery).WillReturnResult(sqlxmock.NewResult(1, 1))
		if v.updateTableErr {
			mock.ExpectExec(updateTableQuery).WillReturnError(v.err)
			mock.ExpectRollback()
		}
		if v.updateTableRowsAffectedErr {
			mock.ExpectExec(updateTableQuery).WillReturnResult(sqlxmock.NewErrorResult(v.err))
			mock.ExpectRollback()
		}
		if v.updateTableOptimisticLockErr {
			mock.ExpectExec(updateTableQuery).WillReturnResult(sqlxmock.NewResult(0, 0))
			mock.ExpectRollback()
		}
		mock.ExpectExec(updateTableQuery).WillReturnResult(sqlxmock.NewResult(1, 1))
		if v.commitErr {
			mock.ExpectCommit().WillReturnError(v.err)
		}
		mock.ExpectCommit()
		actErr := repo.GuestArrived(context.Background(), &entities.Guest{Name: "dummy", TotalArrivedGuests: 1})
		assert.Equal(t, v.expErr, actErr)
	}
}

func TestAddToGuestList(t *testing.T) {
	getGuestByNameQuery := regexp.QuoteMeta("SELECT * FROM `guests` WHERE name = ?")
	getTableQuery := regexp.QuoteMeta("SELECT * FROM `table` WHERE id=?")
	insertGuestQuery := regexp.QuoteMeta("INSERT INTO `guests` (total_rsvp_guests, tableid, name) VALUES(?, ?, ?)")
	updateTableQuery := regexp.QuoteMeta("UPDATE `table` SET pcapacity=?, version = version + 1 WHERE id = ? AND version = ?")
	type TestCase struct {
		name                 string
		desc                 string
		err                  error
		getTableRows         *sqlxmock.Rows
		beginTxErr           bool
		optimisticLockErr    bool
		rowsAffectedErr      bool
		commitErr            bool
		expErr               error
		getGuestByNameErr    bool
		getTableErr          bool
		insertGuestErr       bool
		insertGuestLastIdErr bool
		updateTableErr       bool
		getGuestByNameRow    *sqlxmock.Rows
	}
	testcases := []TestCase{
		{
			name:              "Happy case",
			desc:              "all ok, no error",
			getGuestByNameRow: sqlxmock.NewRows([]string{}),
			getTableRows:      sqlxmock.NewRows([]string{"id", "capacity", "acapacity", "pcapacity", "version"}).AddRow(1, 6, 6, 6, 2),
		},
		{
			name:              "Sad case",
			desc:              "get guest returns error",
			err:               fmt.Errorf("mock error"),
			getGuestByNameRow: sqlxmock.NewRows([]string{}),
			getGuestByNameErr: true,
			expErr:            errDBErr,
			getTableRows:      sqlxmock.NewRows([]string{}),
		},
		{
			name:              "Sad case",
			desc:              "get guest returns error",
			err:               fmt.Errorf("mock error"),
			getGuestByNameRow: sqlxmock.NewRows([]string{"id", "name", "total_rsvp_guests", "total_arrived_guests", "version", "arrivaltime", "tableid"}).AddRow(1, "dummy", 2, 3, 4, "2021-06-04 04:06:44", 5),
			expErr:            errGuestAlreadyRSVP,
			getTableRows:      sqlxmock.NewRows([]string{}),
		},
		{
			name:              "Sad case",
			desc:              "get table returns error",
			getGuestByNameRow: sqlxmock.NewRows([]string{}),
			getTableRows:      sqlxmock.NewRows([]string{}),
			getTableErr:       true,
			err:               fmt.Errorf("mock error"),
			expErr:            errDBErr,
		},

		{
			name:              "Sad case",
			desc:              "table cannot accomodate",
			err:               fmt.Errorf("mock error"),
			getGuestByNameRow: sqlxmock.NewRows([]string{}),
			expErr:            errTableIsFull,
			getTableRows:      sqlxmock.NewRows([]string{"id", "capacity", "acapacity", "pcapacity", "version"}).AddRow(1, 6, 6, 0, 2),
		},
		{
			name:              "Sad case",
			desc:              "begin tx return error",
			err:               fmt.Errorf("mock error"),
			getGuestByNameRow: sqlxmock.NewRows([]string{}),
			expErr:            errDBErr,
			getTableRows:      sqlxmock.NewRows([]string{"id", "capacity", "acapacity", "pcapacity", "version"}).AddRow(1, 6, 6, 6, 2),
			beginTxErr:        true,
		},
		{
			name:              "Sad case",
			desc:              "insert guest return error",
			err:               fmt.Errorf("mock error"),
			getGuestByNameRow: sqlxmock.NewRows([]string{}),
			expErr:            errDBErr,
			getTableRows:      sqlxmock.NewRows([]string{"id", "capacity", "acapacity", "pcapacity", "version"}).AddRow(1, 6, 6, 6, 2),
			insertGuestErr:    true,
		},
		{
			name:                 "Sad case",
			desc:                 "last insert id return error",
			err:                  fmt.Errorf("last insert id error"),
			getGuestByNameRow:    sqlxmock.NewRows([]string{}),
			expErr:               errDBErr,
			getTableRows:         sqlxmock.NewRows([]string{"id", "capacity", "acapacity", "pcapacity", "version"}).AddRow(1, 6, 6, 6, 2),
			insertGuestLastIdErr: true,
		},
		{
			name:              "Sad case",
			desc:              "update table return error",
			err:               fmt.Errorf("mock error"),
			getGuestByNameRow: sqlxmock.NewRows([]string{}),
			expErr:            errDBErr,
			getTableRows:      sqlxmock.NewRows([]string{"id", "capacity", "acapacity", "pcapacity", "version"}).AddRow(1, 6, 6, 6, 2),
			updateTableErr:    true,
		},
		{
			name:              "Sad case",
			desc:              "rows affected return error",
			err:               fmt.Errorf("mock error"),
			getGuestByNameRow: sqlxmock.NewRows([]string{}),
			expErr:            errDBErr,
			getTableRows:      sqlxmock.NewRows([]string{"id", "capacity", "acapacity", "pcapacity", "version"}).AddRow(1, 6, 6, 6, 2),
			rowsAffectedErr:   true,
		},
		{
			name:              "Sad case",
			desc:              "optimistic lock error",
			getGuestByNameRow: sqlxmock.NewRows([]string{}),
			expErr:            errFailedOptimisticLock,
			getTableRows:      sqlxmock.NewRows([]string{"id", "capacity", "acapacity", "pcapacity", "version"}).AddRow(1, 6, 6, 6, 2),
			rowsAffectedErr:   true,
		},
		{
			name:              "Sad case",
			desc:              "commit error",
			err:               fmt.Errorf("mock error"),
			getGuestByNameRow: sqlxmock.NewRows([]string{}),
			expErr:            errDBErr,
			getTableRows:      sqlxmock.NewRows([]string{"id", "capacity", "acapacity", "pcapacity", "version"}).AddRow(1, 6, 6, 6, 2),
			commitErr:         true,
		},
	}
	for _, v := range testcases {
		db, mock := NewMockDb()
		repo := NewDbRepo(db)
		if v.getGuestByNameErr {
			mock.ExpectQuery(getGuestByNameQuery).WillReturnError(v.err)
		}
		mock.ExpectQuery(getGuestByNameQuery).WillReturnRows(v.getGuestByNameRow)
		if v.getTableErr {
			mock.ExpectQuery(getTableQuery).WillReturnError(v.err)
		}
		mock.ExpectQuery(getTableQuery).WillReturnRows(v.getTableRows)
		if v.beginTxErr {
			mock.ExpectBegin().WillReturnError(v.err)
		}
		mock.ExpectBegin()
		if v.insertGuestErr {
			mock.ExpectExec(insertGuestQuery).WillReturnError(v.err)
			mock.ExpectRollback()
		}
		if v.insertGuestLastIdErr {
			mock.ExpectExec(insertGuestQuery).WillReturnResult(sqlxmock.NewErrorResult(v.err))
			mock.ExpectRollback()
		}
		mock.ExpectExec(insertGuestQuery).WillReturnResult(sqlxmock.NewResult(1, 1))
		if v.updateTableErr {
			mock.ExpectExec(updateTableQuery).WillReturnError(v.err)
			mock.ExpectRollback()
		}
		if v.rowsAffectedErr {
			mock.ExpectExec(updateTableQuery).WillReturnResult(sqlxmock.NewErrorResult(v.err))
			mock.ExpectRollback()
		}
		if v.optimisticLockErr {
			mock.ExpectExec(updateTableQuery).WillReturnResult(sqlxmock.NewResult(0, 0))
			mock.ExpectRollback()
		}
		mock.ExpectExec(updateTableQuery).WillReturnResult(sqlxmock.NewResult(1, 1))
		if v.commitErr {
			mock.ExpectCommit().WillReturnError(v.err)
		}
		mock.ExpectCommit()
		actErr := repo.AddToGuestList(context.Background(), &entities.Guest{Name: "dummy", TableID: 1, TotalGuests: 2})
		assert.Equal(t, v.expErr, actErr)
	}
}

func TestGuestDepart(t *testing.T) {
	getGuestByNameQuery := regexp.QuoteMeta("SELECT * FROM `guests` WHERE name = ?")
	getTableQuery := regexp.QuoteMeta("SELECT * FROM `table` WHERE id=?")
	updateGuestQuery := regexp.QuoteMeta("UPDATE `guests` SET total_arrived_guests=0, version = version + 1, arrivaltime=\"\" WHERE id = ? AND version = ?")
	updateTableQuery := regexp.QuoteMeta("UPDATE `table` SET acapacity=?, version = version + 1 WHERE id = ? AND version = ?")
	type TestCase struct {
		name                         string
		desc                         string
		err                          error
		getGuestByNameErr            bool
		getGuestByNameRow            *sqlxmock.Rows
		getTableErr                  bool
		getTableRow                  *sqlxmock.Rows
		updateGuestErr               bool
		updateGuestRowAffectedErr    bool
		updateGuestOptimisticLockErr bool
		updateTableErr               bool
		updateTableRowAffectedErr    bool
		updateTableOptimisticLockErr bool
		beginTxErr                   bool
		commitErr                    bool
		expErr                       error
	}
	testcases := []TestCase{

		{
			name:              "Happy case",
			desc:              "all ok, no error",
			getTableRow:       sqlxmock.NewRows([]string{"id", "capacity", "acapacity", "pcapacity", "version"}).AddRow(1, 6, 6, 6, 2),
			getGuestByNameRow: sqlxmock.NewRows([]string{"id", "name", "total_rsvp_guests", "total_arrived_guests", "version", "arrivaltime", "tableid"}).AddRow(1, "dummy", 2, 3, 4, "2021-06-04 04:06:44", 5),
		},

		{
			name:              "Sad case",
			desc:              "get guest returns error",
			err:               fmt.Errorf("mock error"),
			getGuestByNameErr: true,
			getGuestByNameRow: sqlxmock.NewRows([]string{}),
			getTableRow:       sqlxmock.NewRows([]string{}),
			expErr:            errDBErr,
		},
		{
			name:              "Sad case",
			desc:              "guest not arrived",
			getGuestByNameRow: sqlxmock.NewRows([]string{"id", "name", "total_rsvp_guests", "total_arrived_guests", "version", "arrivaltime", "tableid"}).AddRow(1, "dummy", 2, 0, 4, "", 5),
			getTableRow:       sqlxmock.NewRows([]string{}),
			expErr:            errGuestNotArrived,
		},
		{
			name:              "Sad case",
			desc:              "get table returns error",
			err:               fmt.Errorf("mock error"),
			getTableErr:       true,
			getTableRow:       sqlxmock.NewRows([]string{}),
			getGuestByNameRow: sqlxmock.NewRows([]string{"id", "name", "total_rsvp_guests", "total_arrived_guests", "version", "arrivaltime", "tableid"}).AddRow(1, "dummy", 2, 3, 4, "2021-06-04 04:06:44", 5),
			expErr:            errDBErr,
		},
		{
			name:              "Sad case",
			desc:              "begin tx returns error",
			err:               fmt.Errorf("mock error"),
			beginTxErr:        true,
			getTableRow:       sqlxmock.NewRows([]string{"id", "capacity", "acapacity", "pcapacity", "version"}).AddRow(1, 6, 6, 6, 2),
			getGuestByNameRow: sqlxmock.NewRows([]string{"id", "name", "total_rsvp_guests", "total_arrived_guests", "version", "arrivaltime", "tableid"}).AddRow(1, "dummy", 2, 3, 4, "2021-06-04 04:06:44", 5),
			expErr:            errDBErr,
		},
		{
			name:              "Sad case",
			desc:              "update guest returns error",
			err:               fmt.Errorf("mock error"),
			updateGuestErr:    true,
			getTableRow:       sqlxmock.NewRows([]string{"id", "capacity", "acapacity", "pcapacity", "version"}).AddRow(1, 6, 6, 6, 2),
			getGuestByNameRow: sqlxmock.NewRows([]string{"id", "name", "total_rsvp_guests", "total_arrived_guests", "version", "arrivaltime", "tableid"}).AddRow(1, "dummy", 2, 3, 4, "2021-06-04 04:06:44", 5),
			expErr:            errDBErr,
		},
		{
			name:                      "Sad case",
			desc:                      "update guest rows affected returns error",
			err:                       fmt.Errorf("mock error"),
			updateGuestRowAffectedErr: true,
			getTableRow:               sqlxmock.NewRows([]string{"id", "capacity", "acapacity", "pcapacity", "version"}).AddRow(1, 6, 6, 6, 2),
			getGuestByNameRow:         sqlxmock.NewRows([]string{"id", "name", "total_rsvp_guests", "total_arrived_guests", "version", "arrivaltime", "tableid"}).AddRow(1, "dummy", 2, 3, 4, "2021-06-04 04:06:44", 5),
			expErr:                    errDBErr,
		},
		{
			name:                         "Sad case",
			desc:                         "update guest optimistic lock error",
			updateGuestOptimisticLockErr: true,
			getTableRow:                  sqlxmock.NewRows([]string{"id", "capacity", "acapacity", "pcapacity", "version"}).AddRow(1, 6, 6, 6, 2),
			getGuestByNameRow:            sqlxmock.NewRows([]string{"id", "name", "total_rsvp_guests", "total_arrived_guests", "version", "arrivaltime", "tableid"}).AddRow(1, "dummy", 2, 3, 4, "2021-06-04 04:06:44", 5),
			expErr:                       errFailedOptimisticLock,
		},

		{
			name:              "Sad case",
			desc:              "update guest returns error",
			err:               fmt.Errorf("mock error"),
			updateTableErr:    true,
			getTableRow:       sqlxmock.NewRows([]string{"id", "capacity", "acapacity", "pcapacity", "version"}).AddRow(1, 6, 6, 6, 2),
			getGuestByNameRow: sqlxmock.NewRows([]string{"id", "name", "total_rsvp_guests", "total_arrived_guests", "version", "arrivaltime", "tableid"}).AddRow(1, "dummy", 2, 3, 4, "2021-06-04 04:06:44", 5),
			expErr:            errDBErr,
		},
		{
			name:                      "Sad case",
			desc:                      "update guest rows affected returns error",
			err:                       fmt.Errorf("mock error"),
			updateTableRowAffectedErr: true,
			getTableRow:               sqlxmock.NewRows([]string{"id", "capacity", "acapacity", "pcapacity", "version"}).AddRow(1, 6, 6, 6, 2),
			getGuestByNameRow:         sqlxmock.NewRows([]string{"id", "name", "total_rsvp_guests", "total_arrived_guests", "version", "arrivaltime", "tableid"}).AddRow(1, "dummy", 2, 3, 4, "2021-06-04 04:06:44", 5),
			expErr:                    errDBErr,
		},
		{
			name:                         "Sad case",
			desc:                         "update guest optimistic lock error",
			updateTableOptimisticLockErr: true,
			getTableRow:                  sqlxmock.NewRows([]string{"id", "capacity", "acapacity", "pcapacity", "version"}).AddRow(1, 6, 6, 6, 2),
			getGuestByNameRow:            sqlxmock.NewRows([]string{"id", "name", "total_rsvp_guests", "total_arrived_guests", "version", "arrivaltime", "tableid"}).AddRow(1, "dummy", 2, 3, 4, "2021-06-04 04:06:44", 5),
			expErr:                       errFailedOptimisticLock,
		},
		{
			name:              "Sad case",
			desc:              "commit returns error",
			err:               fmt.Errorf("mock error"),
			commitErr:         true,
			getTableRow:       sqlxmock.NewRows([]string{"id", "capacity", "acapacity", "pcapacity", "version"}).AddRow(1, 6, 6, 6, 2),
			getGuestByNameRow: sqlxmock.NewRows([]string{"id", "name", "total_rsvp_guests", "total_arrived_guests", "version", "arrivaltime", "tableid"}).AddRow(1, "dummy", 2, 3, 4, "2021-06-04 04:06:44", 5),
			expErr:            errDBErr,
		},
	}
	for _, v := range testcases {
		db, mock := NewMockDb()
		repo := NewDbRepo(db)
		if v.getGuestByNameErr {
			mock.ExpectQuery(getGuestByNameQuery).WillReturnError(v.err)
		}
		mock.ExpectQuery(getGuestByNameQuery).WillReturnRows(v.getGuestByNameRow)
		if v.getTableErr {
			mock.ExpectQuery(getTableQuery).WillReturnError(v.err)
		}
		mock.ExpectQuery(getTableQuery).WillReturnRows(v.getTableRow)
		if v.beginTxErr {
			mock.ExpectBegin().WillReturnError(v.err)
		}
		mock.ExpectBegin()
		if v.updateGuestErr {
			mock.ExpectExec(updateGuestQuery).WillReturnError(v.err)
			mock.ExpectRollback()
		}
		if v.updateGuestRowAffectedErr {
			mock.ExpectExec(updateGuestQuery).WillReturnResult(sqlxmock.NewErrorResult(v.err))
			mock.ExpectRollback()
		}
		if v.updateGuestOptimisticLockErr {
			mock.ExpectExec(updateGuestQuery).WillReturnResult(sqlxmock.NewResult(0, 0))
			mock.ExpectRollback()
		}
		mock.ExpectExec(updateGuestQuery).WillReturnResult(sqlxmock.NewResult(1, 1))

		if v.updateTableErr {
			mock.ExpectExec(updateTableQuery).WillReturnError(v.err)
			mock.ExpectRollback()
		}
		if v.updateTableRowAffectedErr {
			mock.ExpectExec(updateTableQuery).WillReturnResult(sqlxmock.NewErrorResult(v.err))
			mock.ExpectRollback()
		}
		if v.updateTableOptimisticLockErr {
			mock.ExpectExec(updateTableQuery).WillReturnResult(sqlxmock.NewResult(0, 0))
			mock.ExpectRollback()
		}
		mock.ExpectExec(updateTableQuery).WillReturnResult(sqlxmock.NewResult(1, 1))

		if v.commitErr {
			mock.ExpectCommit().WillReturnError(v.err)
		}
		mock.ExpectCommit()
		actErr := repo.GuestDepart(context.Background(), &entities.Guest{Name: "dummy"})
		assert.Equal(t, v.expErr, actErr)
	}
}

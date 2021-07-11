package services

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"ggv2/entities"
	"ggv2/repo/mocks"
)

func TestCreateTable(t *testing.T) {
	type TestCase struct {
		name string
		desc string
		err  error
		res  *entities.Table
	}
	testcases := []TestCase{
		{
			name: "Happy case",
			desc: "all ok",
			res: &entities.Table{
				TableID:           1,
				Capacity:          7,
				AvailableCapacity: 7,
				PlannedCapacity:   7,
				Version:           0,
			},
		},
		{
			name: "Sad case",
			desc: "repo return error",
			err:  fmt.Errorf("mock error"),
		},
	}

	for _, v := range testcases {
		repo := new(mocks.DbRepo)
		dbService := &DBService{repo: repo}
		repo.On("CreateTable", context.Background(), &entities.Table{Capacity: 7, AvailableCapacity: 7, PlannedCapacity: 7}).Return(v.res, v.err)
		actT, actErr := dbService.CreateTable(context.Background(), 7)
		assert.Equal(t, v.res, actT)
		assert.Equal(t, v.err, actErr)
	}
}

func TestListTables(t *testing.T) {
	type TestCase struct {
		name string
		desc string
		err  error
		res  []*entities.Table
	}
	testcases := []TestCase{
		{
			name: "Happy case",
			desc: "all ok",
			res: []*entities.Table{
				{
					TableID:           1,
					Capacity:          7,
					AvailableCapacity: 7,
					PlannedCapacity:   7,
					Version:           0,
				},
			},
		},
		{
			name: "Sad case",
			desc: "repo return error",
			err:  fmt.Errorf("mock error"),
		},
	}
	for _, v := range testcases {
		repo := new(mocks.DbRepo)
		dbService := &DBService{repo: repo}
		repo.On("ListTables", context.Background(), int64(10), int64(0)).Return(v.res, v.err)
		actRes, actErr := dbService.ListTables(context.Background(), 10, 0)
		assert.Equal(t, v.res, actRes)
		assert.Equal(t, v.err, actErr)
	}
}

func TestGetEmptySeatsCount(t *testing.T) {
	type TestCase struct {
		name string
		desc string
		err  error
		res  int
	}
	testcases := []TestCase{
		{
			name: "Happy case",
			desc: "all ok",
			res:  99,
		},
		{
			name: "Sad case",
			desc: "repo return error",
			err:  fmt.Errorf("mock error"),
			res:  0,
		},
	}
	for _, v := range testcases {
		repo := new(mocks.DbRepo)
		dbService := &DBService{repo: repo}
		repo.On("GetEmptySeatsCount", context.Background()).Return(v.res, v.err)
		actRes, actErr := dbService.GetEmptySeatsCount(context.Background())
		assert.Equal(t, v.res, actRes)
		assert.Equal(t, v.err, actErr)
	}
}

func TestAddToGuestList(t *testing.T) {
	type TestCase struct {
		name string
		desc string
		err  error
	}
	testcases := []TestCase{
		{
			name: "Happy case",
			desc: "all ok",
		},
		{
			name: "Sad case",
			desc: "repo return error",
			err:  fmt.Errorf("mock error"),
		},
	}
	for _, v := range testcases {
		repo := new(mocks.DbRepo)
		dbService := &DBService{repo: repo}
		repo.On("AddToGuestList", context.Background(), &entities.Guest{Name: "dummy", TableID: 1, TotalGuests: 3}).Return(v.err)
		actErr := dbService.AddToGuestList(context.Background(), 2, 1, "dummy")
		assert.Equal(t, v.err, actErr)
	}
}

func TestListRSVPGuests(t *testing.T) {
	type TestCase struct {
		name string
		desc string
		err  error
		res  []*entities.Guest
	}
	testcases := []TestCase{
		{
			name: "Happy case",
			desc: "all ok",
			res: []*entities.Guest{
				{
					ID:          1,
					Name:        "dummy",
					TableID:     2,
					TotalGuests: 3,
					ArrivalTime: "2021-06-04 04:06:44",
					Version:     4,
				},
			},
		},
		{
			name: "Sad case",
			desc: "repo return error",
			err:  fmt.Errorf("mock error"),
		},
	}
	for _, v := range testcases {
		repo := new(mocks.DbRepo)
		dbService := &DBService{repo: repo}
		repo.On("ListGuests", context.Background(), int64(10), int64(0)).Return(v.res, v.err)
		actRes, actErr := dbService.ListRSVPGuests(context.Background(), 10, 0)
		assert.Equal(t, v.err, actErr)
		assert.Equal(t, v.res, actRes)
	}
}

func TestGuestDepart(t *testing.T) {
	type TestCase struct {
		name string
		desc string
		err  error
	}
	testcases := []TestCase{
		{
			name: "Happy case",
			desc: "all ok",
		},
		{
			name: "Sad case",
			desc: "repo return error",
			err:  fmt.Errorf("mock error"),
		},
	}
	for _, v := range testcases {
		repo := new(mocks.DbRepo)
		dbService := &DBService{repo: repo}
		repo.On("GuestDepart", context.Background(), &entities.Guest{Name: "dummy"}).Return(v.err)
		actErr := dbService.GuestDepart(context.Background(), "dummy")
		assert.Equal(t, v.err, actErr)
	}
}

func TestGuestArrived(t *testing.T) {
	type TestCase struct {
		name string
		desc string
		err  error
	}
	testcases := []TestCase{
		{
			name: "Happy case",
			desc: "all ok",
		},
		{
			name: "Sad case",
			desc: "repo return error",
			err:  fmt.Errorf("mock error"),
		},
	}
	for _, v := range testcases {
		repo := new(mocks.DbRepo)
		dbService := &DBService{repo: repo}
		repo.On("GuestArrived", context.Background(), &entities.Guest{Name: "dummy", TotalArrivedGuests: 2}).Return(v.err)
		actErr := dbService.GuestArrival(context.Background(), 1, "dummy")
		assert.Equal(t, v.err, actErr)
	}
}

func TestEmptyTables(t *testing.T) {
	type TestCase struct {
		name string
		desc string
		err  error
	}
	testcases := []TestCase{
		{
			name: "Happy case",
			desc: "all ok",
		},
		{
			name: "Sad case",
			desc: "repo return error",
			err:  fmt.Errorf("mock error"),
		},
	}
	for _, v := range testcases {
		repo := new(mocks.DbRepo)
		dbService := &DBService{repo: repo}
		repo.On("EmptyTables", context.Background()).Return(v.err)
		actErr := dbService.EmptyTables(context.Background())
		assert.Equal(t, v.err, actErr)
	}
}

func TestListArrivedGuest(t *testing.T) {
	type TestCase struct {
		name string
		desc string
		err  error
		res  []*entities.Guest
	}
	testcases := []TestCase{
		{
			name: "Happy case",
			desc: "all ok",
			res: []*entities.Guest{
				{
					ID:                 1,
					Name:               "dummy",
					TableID:            2,
					TotalArrivedGuests: 3,
					ArrivalTime:        "2021-06-04 04:06:44",
					Version:            4,
				},
			},
		},
		{
			name: "Sad case",
			desc: "repo return error",
			err:  fmt.Errorf("mock error"),
		},
	}
	for _, v := range testcases {
		repo := new(mocks.DbRepo)
		dbService := &DBService{repo: repo}
		repo.On("ListArrivedGuests", context.Background(), int64(10), int64(0)).Return(v.res, v.err)
		actRes, actErr := dbService.ListArrivedGuests(context.Background(), 10, 0)
		assert.Equal(t, v.err, actErr)
		assert.Equal(t, v.res, actRes)
	}
}

func TestGetTable(t *testing.T) {
	type TestCase struct {
		name string
		desc string
		err  error
		res  *entities.Table
	}
	testcases := []TestCase{
		{
			name: "Happy case",
			desc: "all ok",
			res: &entities.Table{
				TableID:           1,
				Capacity:          2,
				AvailableCapacity: 3,
				PlannedCapacity:   4,
				Version:           5,
			},
		},
		{
			name: "Sad case",
			desc: "repo return error",
			err:  fmt.Errorf("mock error"),
		},
	}
	for _, v := range testcases {
		repo := new(mocks.DbRepo)
		dbService := &DBService{repo: repo}
		repo.On("GetTable", context.Background(), int64(1)).Return(v.res, v.err)
		actRes, actErr := dbService.GetTable(context.Background(), 1)
		assert.Equal(t, v.err, actErr)
		assert.Equal(t, v.res, actRes)
	}
}

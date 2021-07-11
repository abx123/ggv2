package handler

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"

	"ggv2/entities"
	"ggv2/services/mocks"
)

func TestGetTables(t *testing.T) {
	type TestCase struct {
		name     string
		desc     string
		url      string
		err      error
		expRes   []*entities.Table
		httpCode int
	}
	testcases := []TestCase{
		{
			name: "Happy case",
			desc: "All ok",
			expRes: []*entities.Table{
				{
					TableID:  1,
					Capacity: 2,
				},
			},
			url:      "http://localhost:1323/tables",
			httpCode: http.StatusOK,
		},
		{
			name:     "Sad case",
			desc:     "service returns error",
			expRes:   []*entities.Table{},
			err:      fmt.Errorf("mock error"),
			httpCode: http.StatusInternalServerError,
			url:      "http://localhost:1323/tables",
		},
		{
			name:     "Sad case",
			desc:     "invalid limit path param",
			expRes:   []*entities.Table{},
			httpCode: http.StatusBadRequest,
			url:      "http://localhost:1323/tables?limit=error",
		},
		{
			name:     "Sad case",
			desc:     "invalid offset path param",
			expRes:   []*entities.Table{},
			httpCode: http.StatusBadRequest,
			url:      "http://localhost:1323/tables?offset=error",
		},
	}
	for _, v := range testcases {
		dbSvc := new(mocks.DbService)
		ctx := context.Background()
		ctx = context.WithValue(ctx, contextKeyRequestID, "")

		dbSvc.On("ListTables", ctx, int64(10), int64(0)).Return(v.expRes, v.err)
		handler := Handler{dbSvc}
		req := httptest.NewRequest("GET", v.url, nil)
		w := httptest.NewRecorder()
		r := echo.New()
		r.GET("/tables", handler.GetTables)
		r.ServeHTTP(w, req)
		assert.Equal(t, v.httpCode, w.Code)
	}
}

func TestGetTable(t *testing.T) {
	type TestCase struct {
		name     string
		desc     string
		err      error
		expRes   *entities.Table
		httpCode int
		url      string
	}
	testcases := []TestCase{
		{
			name: "Happy case",
			desc: "All ok",
			expRes: &entities.Table{
				TableID:           1,
				Capacity:          2,
				AvailableCapacity: 3,
				PlannedCapacity:   4,
				Version:           5,
			},
			httpCode: http.StatusOK,
			url:      "http://localhost:1323/table/1",
		},
		{
			name:     "Sad case",
			desc:     "service returns error",
			expRes:   &entities.Table{},
			err:      fmt.Errorf("mock error"),
			httpCode: http.StatusInternalServerError,
			url:      "http://localhost:1323/table/1",
		},
		{
			name:     "Sad case",
			desc:     "invalid request param",
			expRes:   &entities.Table{},
			err:      fmt.Errorf("mock error"),
			httpCode: http.StatusBadRequest,
			url:      "http://localhost:1323/table/invalid",
		},
	}
	for _, v := range testcases {
		dbSvc := new(mocks.DbService)
		handler := Handler{dbSvc}
		ctx := context.Background()
		ctx = context.WithValue(ctx, contextKeyRequestID, "")
		dbSvc.On("GetTable", ctx, int64(1)).Return(v.expRes, v.err)
		req := httptest.NewRequest("GET", v.url, nil)
		w := httptest.NewRecorder()
		r := echo.New()
		r.GET("/table/:id", handler.GetTable)
		r.ServeHTTP(w, req)
		assert.Equal(t, v.httpCode, w.Code)
	}
}

func TestCreateTables(t *testing.T) {
	type TestCase struct {
		name     string
		desc     string
		err      error
		expRes   *entities.Table
		httpCode int
		capacity string
	}
	testcases := []TestCase{
		{
			name: "Happy case",
			desc: "All ok",
			expRes: &entities.Table{
				TableID:           1,
				Capacity:          2,
				AvailableCapacity: 3,
				PlannedCapacity:   4,
				Version:           5,
			},
			httpCode: http.StatusCreated,
			capacity: "5",
		},
		{
			name:     "Sad case",
			desc:     "service returns error",
			expRes:   &entities.Table{},
			err:      fmt.Errorf("mock error"),
			httpCode: http.StatusInternalServerError,
			capacity: "5",
		},
		{
			name:     "Sad case",
			desc:     "capacity < 1",
			expRes:   &entities.Table{},
			err:      fmt.Errorf("mock error"),
			httpCode: http.StatusBadRequest,
			capacity: "-1",
		},
		{
			name:     "Sad case",
			desc:     "invalid form data",
			expRes:   &entities.Table{},
			err:      fmt.Errorf("mock error"),
			httpCode: http.StatusBadRequest,
			capacity: "error",
		},
	}
	for _, v := range testcases {
		dbSvc := new(mocks.DbService)
		ctx := context.Background()
		ctx = context.WithValue(ctx, contextKeyRequestID, "")
		dbSvc.On("CreateTable", ctx, int64(5)).Return(v.expRes, v.err)
		handler := Handler{dbSvc}
		form := url.Values{}
		if v.capacity != "" {
			form.Add("capacity", v.capacity)
		}
		req := httptest.NewRequest(http.MethodPut, "http://localhost:1323/table", strings.NewReader(form.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		w := httptest.NewRecorder()
		r := echo.New()
		r.PUT("/table", handler.CreateTable)
		r.ServeHTTP(w, req)
		assert.Equal(t, v.httpCode, w.Code)
	}
}

func TestEmptyTables(t *testing.T) {
	type TestCase struct {
		name     string
		desc     string
		err      error
		httpCode int
	}
	testcases := []TestCase{
		{
			name:     "Happy case",
			desc:     "All ok",
			httpCode: http.StatusOK,
		},
		{
			name:     "Sad case",
			desc:     "service returns error",
			err:      fmt.Errorf("mock error"),
			httpCode: http.StatusInternalServerError,
		},
	}
	for _, v := range testcases {
		dbSvc := new(mocks.DbService)
		ctx := context.Background()
		ctx = context.WithValue(ctx, contextKeyRequestID, "")
		dbSvc.On("EmptyTables", ctx).Return(v.err)
		handler := Handler{dbSvc}
		req := httptest.NewRequest(http.MethodGet, "http://localhost:1323/empty_tables", nil)
		w := httptest.NewRecorder()
		r := echo.New()
		r.GET("/empty_tables", handler.EmptyTables)
		r.ServeHTTP(w, req)
		assert.Equal(t, v.httpCode, w.Code)
	}
}

func TestGetEmptySeatsCount(t *testing.T) {
	type TestCase struct {
		name     string
		desc     string
		err      error
		expRes   int
		httpCode int
	}
	testcases := []TestCase{
		{
			name:     "Happy case",
			desc:     "All ok",
			expRes:   99,
			httpCode: http.StatusOK,
		},
		{
			name:     "Sad case",
			desc:     "service returns error",
			expRes:   0,
			err:      fmt.Errorf("mock error"),
			httpCode: http.StatusInternalServerError,
		},
	}
	for _, v := range testcases {
		dbSvc := new(mocks.DbService)
		ctx := context.Background()
		ctx = context.WithValue(ctx, contextKeyRequestID, "")
		dbSvc.On("GetEmptySeatsCount", ctx).Return(v.expRes, v.err)
		handler := Handler{dbSvc}
		req := httptest.NewRequest(http.MethodGet, "http://localhost:1323/seats_empty", nil)
		w := httptest.NewRecorder()
		r := echo.New()
		r.GET("/seats_empty", handler.GetEmptySeatsCount)
		r.ServeHTTP(w, req)
		assert.Equal(t, v.httpCode, w.Code)
	}
}

func TestAddToGuestList(t *testing.T) {
	type TestCase struct {
		name               string
		desc               string
		err                error
		httpCode           int
		table              string
		accompanyingGuests string
	}
	testcases := []TestCase{
		{
			name:               "Happy case",
			desc:               "All ok",
			httpCode:           http.StatusCreated,
			table:              "1",
			accompanyingGuests: "2",
		},
		{
			name:               "Sad case",
			desc:               "invalid form data",
			httpCode:           http.StatusBadRequest,
			table:              "invalid",
			accompanyingGuests: "2",
		},
		{
			name:               "Sad case",
			desc:               "tableid < 1",
			httpCode:           http.StatusBadRequest,
			table:              "-5",
			accompanyingGuests: "2",
		},
		{
			name:               "Sad case",
			desc:               "accompanyingGuests < 0",
			httpCode:           http.StatusBadRequest,
			table:              "1",
			accompanyingGuests: "-2",
		},
		{
			name:               "Sad case",
			desc:               "Table not found",
			httpCode:           http.StatusNotFound,
			err:                fmt.Errorf("table not found"),
			table:              "1",
			accompanyingGuests: "2",
		},
		{
			name:               "Sad case",
			desc:               "db error",
			httpCode:           http.StatusInternalServerError,
			err:                fmt.Errorf("mock error"),
			table:              "1",
			accompanyingGuests: "2",
		},
	}
	for _, v := range testcases {
		dbSvc := new(mocks.DbService)
		ctx := context.Background()
		ctx = context.WithValue(ctx, contextKeyRequestID, "")
		dbSvc.On("AddToGuestList", ctx, int64(2), int64(1), "dummy").Return(v.err)
		handler := Handler{dbSvc}
		form := url.Values{}
		if v.table != "" {
			form.Add("table", v.table)
		}
		if v.accompanyingGuests != "" {
			form.Add("accompanying_guests", v.accompanyingGuests)
		}

		req := httptest.NewRequest(http.MethodPost, "http://localhost:1323/guest_list/dummy", strings.NewReader(form.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		w := httptest.NewRecorder()
		r := echo.New()
		r.POST("/guest_list/:name", handler.AddToGuestList)
		r.ServeHTTP(w, req)
		assert.Equal(t, v.httpCode, w.Code)
	}
}

func TestPing(t *testing.T) {
	dbSvc := new(mocks.DbService)
	handler := Handler{dbSvc}
	req := httptest.NewRequest(http.MethodGet, "http://localhost:1323/ping", nil)
	w := httptest.NewRecorder()
	r := echo.New()
	r.GET("/ping", handler.Ping)
	r.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)
}

func TestGetGuestList(t *testing.T) {
	type TestCase struct {
		name     string
		desc     string
		err      error
		expRes   []*entities.Guest
		url      string
		httpCode int
	}
	testcases := []TestCase{
		{
			name: "Happy case",
			desc: "All ok",
			expRes: []*entities.Guest{
				{
					ID:          1,
					Name:        "dummy",
					TableID:     2,
					TotalGuests: 3,
					ArrivalTime: "2021-06-04 04:06:44",
					Version:     4,
				},
			},
			url:      "http://localhost:1323/guest_list",
			httpCode: http.StatusOK,
		},
		{
			name:     "Sad case",
			desc:     "service returns error",
			expRes:   []*entities.Guest{},
			url:      "http://localhost:1323/guest_list",
			err:      fmt.Errorf("mock error"),
			httpCode: http.StatusInternalServerError,
		},
		{
			name:     "Sad case",
			desc:     "invalid limit path parameter",
			expRes:   []*entities.Guest{},
			url:      "http://localhost:1323/guest_list?limit=error",
			httpCode: http.StatusBadRequest,
		},
		{
			name:     "Sad case",
			desc:     "invalid offset path parameter",
			expRes:   []*entities.Guest{},
			url:      "http://localhost:1323/guest_list?offset=error",
			httpCode: http.StatusBadRequest,
		},
	}
	for _, v := range testcases {
		dbSvc := new(mocks.DbService)
		ctx := context.Background()
		ctx = context.WithValue(ctx, contextKeyRequestID, "")
		dbSvc.On("ListRSVPGuests", ctx, int64(10), int64(0)).Return(v.expRes, v.err)
		handler := Handler{dbSvc}
		req := httptest.NewRequest(http.MethodGet, v.url, nil)
		w := httptest.NewRecorder()
		r := echo.New()
		r.GET("/guest_list", handler.GetGuestList)
		r.ServeHTTP(w, req)
		assert.Equal(t, v.httpCode, w.Code)
	}
}

func TestGuestArrived(t *testing.T) {
	type TestCase struct {
		name               string
		desc               string
		err                error
		httpCode           int
		accompanyingGuests string
	}
	testcases := []TestCase{
		{
			name:               "Happy case",
			desc:               "All ok",
			httpCode:           http.StatusCreated,
			accompanyingGuests: "2",
		},
		{
			name:               "Sad case",
			desc:               "invalid form data",
			httpCode:           http.StatusBadRequest,
			accompanyingGuests: "invalid",
		},
		{
			name:               "Sad case",
			desc:               "accompanyingGuests < 0",
			httpCode:           http.StatusBadRequest,
			accompanyingGuests: "-2",
		},
		{
			name:               "Sad case",
			desc:               "no rsvp/alrady checked-in error",
			httpCode:           http.StatusNotFound,
			err:                fmt.Errorf("guest did not register"),
			accompanyingGuests: "2",
		},
		{
			name:               "Sad case",
			desc:               "db error",
			httpCode:           http.StatusInternalServerError,
			err:                fmt.Errorf("mock error"),
			accompanyingGuests: "2",
		},
	}
	for _, v := range testcases {
		dbSvc := new(mocks.DbService)
		ctx := context.Background()
		ctx = context.WithValue(ctx, contextKeyRequestID, "")
		dbSvc.On("GuestArrival", ctx, int64(2), "dummy").Return(v.err)
		handler := Handler{dbSvc}
		form := url.Values{}
		if v.accompanyingGuests != "" {
			form.Add("accompanying_guests", v.accompanyingGuests)
		}
		req := httptest.NewRequest(http.MethodPut, "http://localhost:1323/guests/dummy", strings.NewReader(form.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		w := httptest.NewRecorder()
		r := echo.New()
		r.PUT("/guests/:name", handler.GuestArrived)
		r.ServeHTTP(w, req)
		assert.Equal(t, v.httpCode, w.Code)
	}
}

func TestListArrivedGuest(t *testing.T) {
	type TestCase struct {
		name     string
		desc     string
		err      error
		expRes   []*entities.Guest
		url      string
		httpCode int
	}
	testcases := []TestCase{
		{
			name: "Happy case",
			desc: "All ok",
			expRes: []*entities.Guest{
				{
					ID:                 1,
					Name:               "dummy",
					TableID:            2,
					TotalArrivedGuests: 3,
					ArrivalTime:        "2021-06-04 04:06:44",
					Version:            4,
				},
			},
			url:      "http://localhost:1323/guests",
			httpCode: http.StatusOK,
		},
		{
			name:     "Sad case",
			desc:     "service returns error",
			expRes:   []*entities.Guest{},
			err:      fmt.Errorf("mock error"),
			httpCode: http.StatusInternalServerError,
			url:      "http://localhost:1323/guests",
		},
		{
			name:     "Sad case",
			desc:     "invalid limit path parameter",
			expRes:   []*entities.Guest{},
			httpCode: http.StatusBadRequest,
			url:      "http://localhost:1323/guests?limit=error",
		},
		{
			name:     "Sad case",
			desc:     "invalid offset path parameter",
			expRes:   []*entities.Guest{},
			httpCode: http.StatusBadRequest,
			url:      "http://localhost:1323/guests?offset=error",
		},
	}
	for _, v := range testcases {
		dbSvc := new(mocks.DbService)
		ctx := context.Background()
		ctx = context.WithValue(ctx, contextKeyRequestID, "")
		dbSvc.On("ListArrivedGuests", ctx, int64(10), int64(0)).Return(v.expRes, v.err)
		handler := Handler{dbSvc}
		req := httptest.NewRequest(http.MethodGet, v.url, nil)
		w := httptest.NewRecorder()
		r := echo.New()
		r.GET("/guests", handler.ListArrivedGuest)
		r.ServeHTTP(w, req)
		assert.Equal(t, v.httpCode, w.Code)
	}
}

func TestGuestDepart(t *testing.T) {
	type TestCase struct {
		name     string
		desc     string
		err      error
		httpCode int
	}
	testcases := []TestCase{
		{
			name:     "Happy case",
			desc:     "All ok",
			httpCode: http.StatusAccepted,
		},
		{
			name:     "Sad case",
			desc:     "service returns error",
			err:      fmt.Errorf("mock error"),
			httpCode: http.StatusInternalServerError,
		},
		{
			name:     "Sad case",
			desc:     "guest not found error",
			err:      errGuestNotFound,
			httpCode: http.StatusNotFound,
		},
	}
	for _, v := range testcases {
		dbSvc := new(mocks.DbService)
		ctx := context.Background()
		ctx = context.WithValue(ctx, contextKeyRequestID, "")
		dbSvc.On("GuestDepart", ctx, "dummy").Return(v.err)
		handler := Handler{dbSvc}
		req := httptest.NewRequest(http.MethodDelete, "http://localhost:1323/guests/dummy", nil)
		w := httptest.NewRecorder()
		r := echo.New()
		r.DELETE("/guests/:name", handler.GuestDepart)
		r.ServeHTTP(w, req)
		assert.Equal(t, v.httpCode, w.Code)
	}
}

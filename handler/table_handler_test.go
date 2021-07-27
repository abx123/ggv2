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
		
		

		dbSvc.On("ListTables", context.Background(), int64(10), int64(0)).Return(v.expRes, v.err)
		th := TableHandler{dbSvc}
		req := httptest.NewRequest("GET", v.url, nil)
		w := httptest.NewRecorder()
		r := echo.New()
		r.GET("/tables", th.GetTables)
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
		th := TableHandler{dbSvc}
		
		
		dbSvc.On("GetTable", context.Background(), int64(1)).Return(v.expRes, v.err)
		req := httptest.NewRequest("GET", v.url, nil)
		w := httptest.NewRecorder()
		r := echo.New()
		r.GET("/table/:id", th.GetTable)
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
		
		
		dbSvc.On("CreateTable", context.Background(), int64(5)).Return(v.expRes, v.err)
		th := TableHandler{dbSvc}
		form := url.Values{}
		if v.capacity != "" {
			form.Add("capacity", v.capacity)
		}
		req := httptest.NewRequest(http.MethodPut, "http://localhost:1323/table", strings.NewReader(form.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		w := httptest.NewRecorder()
		r := echo.New()
		r.PUT("/table", th.CreateTable)
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
		
		
		dbSvc.On("EmptyTables", context.Background()).Return(v.err)
		th := TableHandler{dbSvc}
		req := httptest.NewRequest(http.MethodGet, "http://localhost:1323/empty_tables", nil)
		w := httptest.NewRecorder()
		r := echo.New()
		r.GET("/empty_tables", th.EmptyTables)
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
		
		
		dbSvc.On("GetEmptySeatsCount", context.Background()).Return(v.expRes, v.err)
		th := TableHandler{dbSvc}
		req := httptest.NewRequest(http.MethodGet, "http://localhost:1323/seats_empty", nil)
		w := httptest.NewRecorder()
		r := echo.New()
		r.GET("/seats_empty", th.GetEmptySeatsCount)
		r.ServeHTTP(w, req)
		assert.Equal(t, v.httpCode, w.Code)
	}
}

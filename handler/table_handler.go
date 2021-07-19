package handler

import (
	"net/http"
	"strconv"

	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"

	"go.uber.org/zap"

	"ggv2/handler/presenter"
	"ggv2/repo"
	"ggv2/services"
)

type createTableRequest struct {
	Capacity int64 `json:"capacity" form:"capacity"`
}
type putCreateTableResponse struct {
	Table *presenter.Table `json:"table"`
}

type getSeatsEmptyResponse struct {
	SeatsEmpty int64 `json:"seats_empty"`
}

type TableHandler struct {
	dbSvc services.DbService
}

func NewTableHandler(conn *sqlx.DB) *TableHandler {
	dbRepo := repo.NewDbRepo(conn)
	dbSvc := services.NewDbService(dbRepo)

	return &TableHandler{
		dbSvc: dbSvc,
	}
}

// GetTable handles GET /table/:id"
func (th *TableHandler) GetTable(c echo.Context) (err error) {
	reqID := c.Response().Header().Get(echo.HeaderXRequestID)
	id := c.Param("id")
	tableId, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		// Invalid request parameter
		zap.L().Error(errInvalidRequest.Error(), zap.Error(err))
		return c.JSON(http.StatusBadRequest, presenter.ErrResp(reqID, errInvalidRequest))
	}
	// Query database
	res, err := th.dbSvc.GetTable(c.Request().Context(), tableId)
	if err != nil {
		// Error while querying database
		return c.JSON(http.StatusInternalServerError, presenter.ErrResp(reqID, err))
	}
	// Return ok
	return c.JSON(http.StatusOK, &presenter.Table{
		// TableID:  res.TableID,
		Capacity: res.Capacity,
	})
}

// GetTables handles GET /tables
func (th *TableHandler) GetTables(c echo.Context) (err error) {
	reqID := c.Response().Header().Get(echo.HeaderXRequestID)
	limit, offset, err := getLimitAndOffest(c)
	if err != nil {
		zap.L().Error(errInvalidRequest.Error(), zap.Error(err))
		return c.JSON(http.StatusBadRequest, presenter.ErrResp(reqID, err))
	}
	// Query database
	data, err := th.dbSvc.ListTables(c.Request().Context(), limit, offset)

	if err != nil {
		// Error while querying database
		return c.JSON(http.StatusInternalServerError, presenter.ErrResp(reqID, err))
	}
	// Map response fields
	var tables []*presenter.Table
	for _, d := range data {
		tables = append(tables, &presenter.Table{
			TableID:  d.TableID,
			Capacity: d.Capacity,
		})
	}

	// Return ok
	return c.JSON(http.StatusOK, tables)
}

// CreateTables handles PUT /table
func (th *TableHandler) CreateTable(c echo.Context) (err error) {
	reqID := c.Response().Header().Get(echo.HeaderXRequestID)
	r := new(createTableRequest)
	if err = c.Bind(r); err != nil {
		// Invalid request parameter
		zap.L().Error(errInvalidRequest.Error(), zap.Error(err))
		return c.JSON(http.StatusBadRequest, presenter.ErrResp(reqID, errInvalidRequest))
	}
	if r.Capacity < 1 {
		// Invalid request parameter
		zap.L().Error(errInvalidRequest.Error(), zap.Error(errCapacityLessThanOne))
		return c.JSON(http.StatusBadRequest, presenter.ErrResp(reqID, errCapacityLessThanOne))
	}
	res := &putCreateTableResponse{}
	// Query database
	data, err := th.dbSvc.CreateTable(c.Request().Context(), r.Capacity)
	if err != nil {
		// Error while querying database
		return c.JSON(http.StatusInternalServerError, presenter.ErrResp(reqID, err))
	}
	// Map response fields
	res.Table = &presenter.Table{
		TableID:  data.TableID,
		Capacity: data.Capacity,
	}
	// Return ok
	return c.JSON(http.StatusCreated, res)
}

// Init handles GET /init
func (th *TableHandler) EmptyTables(c echo.Context) (err error) {
	reqID := c.Response().Header().Get(echo.HeaderXRequestID)
	// Query database
	err = th.dbSvc.EmptyTables(c.Request().Context())
	if err != nil {
		// Error while querying database
		return c.JSON(http.StatusInternalServerError, presenter.ErrResp(reqID, err))
	}
	// Return ok
	return c.JSON(http.StatusOK, "Tables emptied!")
}

// GetEmptySeatsCount handles GET /seats_empty
func (th *TableHandler) GetEmptySeatsCount(c echo.Context) (err error) {
	res := &getSeatsEmptyResponse{}
	reqID := c.Response().Header().Get(echo.HeaderXRequestID)
	// Query database
	count, err := th.dbSvc.GetEmptySeatsCount(c.Request().Context())
	if err != nil {
		// Error while querying database
		return c.JSON(http.StatusInternalServerError, presenter.ErrResp(reqID, err))
	}
	// Map response fields
	res.SeatsEmpty = int64(count)
	// Return ok
	return c.JSON(http.StatusOK, res)
}

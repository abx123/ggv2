package handler

import (
	"context"
	"errors"
	"net/http"
	"strconv"

	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"

	"go.uber.org/zap"

	"ggv2/handler/presenter"
	"ggv2/repo"
	"ggv2/services"
)

type contextKey string

const contextKeyRequestID contextKey = "requestID"

var (
	errInvalidRequest                = errors.New("invalid request parameter")
	errCapacityLessThanOne           = errors.New("capacity cannot be less than 1")
	errAccompanyingGuestLessThanZero = errors.New("accompanying guest cannot be less than 0")
	errTableNotFound                 = errors.New("table not found")
	errGuestNotFound                 = errors.New("guest not found")
)

type errResp struct {
	ReqId  string `json:"requestID"`
	ErrMsg string `json:"message"`
}

type Handler struct {
	dbSvc services.DbService
}

type putGuestArrivesRequest struct {
	AccompanyingGuests int64 `json:"accompanying_guests" form:"accompanying_guests"`
}
type putGuestArrivesResponse struct {
	Name string `json:"name"`
}

type postGuestListRequest struct {
	Table              int64 `json:"table" form:"table"`
	AccompanyingGuests int64 `json:"accompanying_guests" form:"accompanying_guests"`
}
type postGuestListResponse struct {
	Name string `json:"name"`
}

type createTableRequest struct {
	Capacity int64 `json:"capacity" form:"capacity"`
}
type putCreateTableResponse struct {
	Table *presenter.Table `json:"table"`
}

type getSeatsEmptyResponse struct {
	SeatsEmpty int64 `json:"seats_empty"`
}

type getGuestListResponse struct {
	Guests []*presenter.Guest `json:"guests"`
}

func NewHandler(conn *sqlx.DB) *Handler {
	dbRepo := repo.NewDbRepo(conn)
	dbSvc := services.NewDbService(dbRepo)

	return &Handler{
		dbSvc: dbSvc,
	}
}

// GetTable handles GET /table/:id"
func (h *Handler) GetTable(c echo.Context) (err error) {
	reqID := c.Response().Header().Get(echo.HeaderXRequestID)
	ctx := context.WithValue(c.Request().Context(), contextKeyRequestID, c.Response().Header().Get(echo.HeaderXRequestID))
	id := c.Param("id")
	tableId, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		// Invalid request parameter
		zap.L().Error(errInvalidRequest.Error(), zap.Error(err))
		return c.JSON(http.StatusBadRequest, returnErr(reqID, errInvalidRequest))
	}
	// Query database
	res, err := h.dbSvc.GetTable(ctx, tableId)
	if err != nil {
		// Error while querying database
		return c.JSON(http.StatusInternalServerError, returnErr(reqID, err))
	}
	// Return ok
	return c.JSON(http.StatusOK, &presenter.Table{
		TableID:  res.TableID,
		Capacity: res.Capacity,
	})
}

// GetTables handles GET /tables
func (con *Handler) GetTables(c echo.Context) (err error) {
	reqID := c.Response().Header().Get(echo.HeaderXRequestID)
	ctx := context.WithValue(c.Request().Context(), contextKeyRequestID, c.Response().Header().Get(echo.HeaderXRequestID))
	limit, offset, err := getLimitAndOffest(c)
	if err != nil {
		zap.L().Error(errInvalidRequest.Error(), zap.Error(err))
		return c.JSON(http.StatusBadRequest, returnErr(reqID, err))
	}
	// Query database
	data, err := con.dbSvc.ListTables(ctx, limit, offset)

	if err != nil {
		// Error while querying database
		return c.JSON(http.StatusInternalServerError, returnErr(reqID, err))
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
func (con *Handler) CreateTable(c echo.Context) (err error) {
	reqID := c.Response().Header().Get(echo.HeaderXRequestID)
	ctx := context.WithValue(c.Request().Context(), contextKeyRequestID, c.Response().Header().Get(echo.HeaderXRequestID))
	r := new(createTableRequest)
	if err = c.Bind(r); err != nil {
		// Invalid request parameter
		zap.L().Error(errInvalidRequest.Error(), zap.Error(err))
		return c.JSON(http.StatusBadRequest, returnErr(reqID, errInvalidRequest))
	}
	if r.Capacity < 1 {
		// Invalid request parameter
		zap.L().Error(errInvalidRequest.Error(), zap.Error(errCapacityLessThanOne))
		return c.JSON(http.StatusBadRequest, returnErr(reqID, errCapacityLessThanOne))
	}
	res := &putCreateTableResponse{}
	// Query database
	data, err := con.dbSvc.CreateTable(ctx, r.Capacity)
	if err != nil {
		// Error while querying database
		return c.JSON(http.StatusInternalServerError, returnErr(reqID, err))
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
func (con *Handler) EmptyTables(c echo.Context) (err error) {
	reqID := c.Response().Header().Get(echo.HeaderXRequestID)
	ctx := context.WithValue(c.Request().Context(), contextKeyRequestID, c.Response().Header().Get(echo.HeaderXRequestID))
	// Query database
	err = con.dbSvc.EmptyTables(ctx)
	if err != nil {
		// Error while querying database
		return c.JSON(http.StatusInternalServerError, returnErr(reqID, err))
	}
	// Return ok
	return c.JSON(http.StatusOK, "Tables emptied!")
}

// GetEmptySeatsCount handles GET /seats_empty
func (con *Handler) GetEmptySeatsCount(c echo.Context) (err error) {
	res := &getSeatsEmptyResponse{}
	reqID := c.Response().Header().Get(echo.HeaderXRequestID)
	ctx := context.WithValue(c.Request().Context(), contextKeyRequestID, c.Response().Header().Get(echo.HeaderXRequestID))
	// Query database
	count, err := con.dbSvc.GetEmptySeatsCount(ctx)
	if err != nil {
		// Error while querying database
		return c.JSON(http.StatusInternalServerError, returnErr(reqID, err))
	}
	// Map response fields
	res.SeatsEmpty = int64(count)
	// Return ok
	return c.JSON(http.StatusOK, res)
}

// AddToGuestList handles POST /guest_list/:name
func (con *Handler) AddToGuestList(c echo.Context) (err error) {
	// Get and validate request parameter
	r := &postGuestListRequest{}
	name := c.Param("name")
	reqID := c.Response().Header().Get(echo.HeaderXRequestID)
	ctx := context.WithValue(c.Request().Context(), contextKeyRequestID, c.Response().Header().Get(echo.HeaderXRequestID))
	if err = c.Bind(r); err != nil {
		// Invalid request parameter
		zap.L().Error(errInvalidRequest.Error(), zap.Error(err))
		return c.JSON(http.StatusBadRequest, returnErr(reqID, errInvalidRequest))
	}
	if r.Table < 1 {
		// Invalid request parameter
		zap.L().Error(errInvalidRequest.Error(), zap.Error(errCapacityLessThanOne))
		return c.JSON(http.StatusBadRequest, returnErr(reqID, errCapacityLessThanOne))
	}
	if r.AccompanyingGuests < 0 {
		// Invalid request parameter
		zap.L().Error(errInvalidRequest.Error(), zap.Error(errAccompanyingGuestLessThanZero))
		return c.JSON(http.StatusBadRequest, returnErr(reqID, errAccompanyingGuestLessThanZero))
	}

	// Query database
	err = con.dbSvc.AddToGuestList(ctx, r.AccompanyingGuests, r.Table, name)
	if err != nil {
		// Error while querying database
		if err == errTableNotFound {

			return c.JSON(http.StatusNotFound, returnErr(reqID, errTableNotFound))
		}
		return c.JSON(http.StatusInternalServerError, returnErr(reqID, err))
	}
	// Return ok
	return c.JSON(http.StatusCreated, postGuestListResponse{Name: name})
}

// Healthcheck handles GET /
func (con *Handler) Ping(c echo.Context) (err error) {
	// Server is up and running, return OK!
	return c.String(http.StatusOK, "Pong")
}

// GetGuestList handles GET /guest_list
func (con *Handler) GetGuestList(c echo.Context) (err error) {
	reqID := c.Response().Header().Get(echo.HeaderXRequestID)
	ctx := context.WithValue(c.Request().Context(), contextKeyRequestID, c.Response().Header().Get(echo.HeaderXRequestID))
	res := getGuestListResponse{}
	limit, offset, err := getLimitAndOffest(c)
	if err != nil {
		zap.L().Error(errInvalidRequest.Error(), zap.Error(err))
		return c.JSON(http.StatusBadRequest, returnErr(reqID, err))
	}
	// Query database
	data, err := con.dbSvc.ListRSVPGuests(ctx, limit, offset)
	if err != nil {
		// Error while querying database
		return c.JSON(http.StatusInternalServerError, returnErr(reqID, err))
	}
	// Map response fields
	var guests []*presenter.Guest
	for _, d := range data {
		guests = append(guests, &presenter.Guest{
			ID:                 d.ID,
			Name:               d.Name,
			TableID:            d.TableID,
			AccompanyingGuests: d.TotalGuests,
		})
	}
	res.Guests = guests
	// Return ok
	return c.JSON(http.StatusOK, res)
}

// GuestArrived handles PUT /guests/:name
func (con *Handler) GuestArrived(c echo.Context) (err error) {

	reqID := c.Response().Header().Get(echo.HeaderXRequestID)
	ctx := context.WithValue(c.Request().Context(), contextKeyRequestID, c.Response().Header().Get(echo.HeaderXRequestID))
	// Get and validate request parameter
	r := putGuestArrivesRequest{}
	name := c.Param("name")
	if err = c.Bind(r); err != nil {
		// Invalid request parameter
		zap.L().Error(errInvalidRequest.Error(), zap.Error(err))
		return c.JSON(http.StatusBadRequest, returnErr(reqID, errInvalidRequest))
	}
	if r.AccompanyingGuests < 0 {
		// Invalid request parameter
		zap.L().Error(errAccompanyingGuestLessThanZero.Error(), zap.Error(errAccompanyingGuestLessThanZero))
		return c.JSON(http.StatusBadRequest, returnErr(reqID, errAccompanyingGuestLessThanZero))
	}
	res := putGuestArrivesResponse{}
	// Query database
	err = con.dbSvc.GuestArrival(ctx, r.AccompanyingGuests, name)
	if err != nil {
		// Error while querying database
		if err.Error() == "guest did not register" || err.Error() == "guest already arrived" {
			return c.JSON(http.StatusNotFound, returnErr(reqID, err))
		}
		return c.JSON(http.StatusInternalServerError, returnErr(reqID, err))
	}
	// Map response fields
	res.Name = name
	// Return ok
	return c.JSON(http.StatusCreated, res)
}

// ListArrivedGuest handles GET /guests
func (con *Handler) ListArrivedGuest(c echo.Context) (err error) {
	reqID := c.Response().Header().Get(echo.HeaderXRequestID)
	ctx := context.WithValue(c.Request().Context(), contextKeyRequestID, c.Response().Header().Get(echo.HeaderXRequestID))
	res := getGuestListResponse{}
	limit, offset, err := getLimitAndOffest(c)
	if err != nil {
		zap.L().Error(errInvalidRequest.Error(), zap.Error(err))
		return c.JSON(http.StatusBadRequest, returnErr(reqID, err))
	}
	// Query database
	data, err := con.dbSvc.ListArrivedGuests(ctx, limit, offset)
	if err != nil {
		// Error while querying database
		return c.JSON(http.StatusInternalServerError, returnErr(reqID, err))
	}
	// Map response fields
	var guests []*presenter.Guest
	for _, d := range data {
		guests = append(guests, &presenter.Guest{
			ID:                 d.ID,
			Name:               d.Name,
			TableID:            d.TableID,
			AccompanyingGuests: d.TotalArrivedGuests,
			ArrivalTime:        d.ArrivalTime,
		})
	}
	res.Guests = guests
	// Return ok
	return c.JSON(http.StatusOK, res)
}

// GuestDepart handles DELETE /guests/:name
func (con *Handler) GuestDepart(c echo.Context) (err error) {
	reqID := c.Response().Header().Get(echo.HeaderXRequestID)
	ctx := context.WithValue(c.Request().Context(), contextKeyRequestID, c.Response().Header().Get(echo.HeaderXRequestID))
	// Get and validate request parameter
	name := c.Param("name")
	// Query database
	err = con.dbSvc.GuestDepart(ctx, name)
	if err != nil {
		// Error while querying database
		if err == errGuestNotFound || err == errTableNotFound {
			return c.JSON(http.StatusNotFound, returnErr(reqID, err))
		}
		return c.JSON(http.StatusInternalServerError, returnErr(reqID, err))
	}
	// Return ok
	return c.JSON(http.StatusAccepted, "OK!")
}

func returnErr(reqID string, err error) *errResp {
	return &errResp{
		ReqId:  reqID,
		ErrMsg: err.Error(),
	}
}

func getLimitAndOffest(c echo.Context) (int64, int64, error) {
	strlimit := c.QueryParam("limit")
	stroffset := c.QueryParam("offset")
	var limit int64 = 10
	var offset int64
	var err error
	if strlimit != "" {
		limit, err = strconv.ParseInt(strlimit, 10, 64)
		if err != nil {
			// zap.L().Error(errInvalidRequest.Error(), zap.Error(err))
			// return c.JSON(http.StatusBadRequest, returnErr(reqID, err))
			return 0, 0, err
		}
	}
	if stroffset != "" {
		offset, err = strconv.ParseInt(stroffset, 10, 64)
		if err != nil {
			// zap.L().Error(errInvalidRequest.Error(), zap.Error(err))
			// return c.JSON(http.StatusBadRequest, returnErr(reqID, err))
			return 0, 0, err
		}
	}

	return limit, offset, nil
}

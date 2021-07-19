package handler

import (
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

var (
	errInvalidRequest                = errors.New("invalid request parameter")
	errCapacityLessThanOne           = errors.New("capacity cannot be less than 1")
	errAccompanyingGuestLessThanZero = errors.New("accompanying guest cannot be less than 0")
	errTableNotFound                 = errors.New("table not found")
	errGuestNotFound                 = errors.New("guest not found")
)

type GuestHandler struct {
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

type getGuestListResponse struct {
	Guests []*presenter.Guest `json:"guests"`
}

func NewGuestHandler(conn *sqlx.DB) *GuestHandler {
	dbRepo := repo.NewDbRepo(conn)
	dbSvc := services.NewDbService(dbRepo)

	return &GuestHandler{
		dbSvc: dbSvc,
	}
}

// AddToGuestList handles POST /guest_list/:name
func (con *GuestHandler) AddToGuestList(c echo.Context) (err error) {
	// Get and validate request parameter
	r := &postGuestListRequest{}
	name := c.Param("name")
	reqID := c.Response().Header().Get(echo.HeaderXRequestID)
	if err = c.Bind(r); err != nil {
		// Invalid request parameter
		zap.L().Error(errInvalidRequest.Error(), zap.Error(err))
		return c.JSON(http.StatusBadRequest, presenter.ErrResp(reqID, errInvalidRequest))
	}
	if r.Table < 1 {
		// Invalid request parameter
		zap.L().Error(errInvalidRequest.Error(), zap.Error(errCapacityLessThanOne))
		return c.JSON(http.StatusBadRequest, presenter.ErrResp(reqID, errCapacityLessThanOne))
	}
	if r.AccompanyingGuests < 0 {
		// Invalid request parameter
		// c.Response().Header().Get(echo.HeaderXRequestID)
		zap.L().Error(errInvalidRequest.Error(), zap.Error(errAccompanyingGuestLessThanZero), zap.String("rqId", reqID))
		return c.JSON(http.StatusBadRequest, presenter.ErrResp(reqID, errAccompanyingGuestLessThanZero))
	}

	// Query database
	err = con.dbSvc.AddToGuestList(c.Request().Context(), r.AccompanyingGuests, r.Table, name)
	if err != nil {
		// Error while querying database
		if err == errTableNotFound {

			return c.JSON(http.StatusNotFound, presenter.ErrResp(reqID, errTableNotFound))
		}
		return c.JSON(http.StatusInternalServerError, presenter.ErrResp(reqID, err))
	}
	// Return ok
	return c.JSON(http.StatusCreated, postGuestListResponse{Name: name})
}

// Healthcheck handles GET /
func (con *GuestHandler) Ping(c echo.Context) (err error) {
	// Server is up and running, return OK!
	return c.String(http.StatusOK, "Pong")
}

// GetGuestList handles GET /guest_list
func (con *GuestHandler) GetGuestList(c echo.Context) (err error) {
	reqID := c.Response().Header().Get(echo.HeaderXRequestID)
	res := getGuestListResponse{}
	limit, offset, err := getLimitAndOffest(c)
	if err != nil {
		zap.L().Error(errInvalidRequest.Error(), zap.Error(err))
		return c.JSON(http.StatusBadRequest, presenter.ErrResp(reqID, err))
	}
	// Query database
	data, err := con.dbSvc.ListRSVPGuests(c.Request().Context(), limit, offset)
	if err != nil {
		// Error while querying database
		return c.JSON(http.StatusInternalServerError, presenter.ErrResp(reqID, err))
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
func (con *GuestHandler) GuestArrived(c echo.Context) (err error) {

	reqID := c.Response().Header().Get(echo.HeaderXRequestID)
	// Get and validate request parameter
	r := putGuestArrivesRequest{}
	name := c.Param("name")
	if err = c.Bind(r); err != nil {
		// Invalid request parameter
		zap.L().Error(errInvalidRequest.Error(), zap.Error(err))
		return c.JSON(http.StatusBadRequest, presenter.ErrResp(reqID, errInvalidRequest))
	}
	if r.AccompanyingGuests < 0 {
		// Invalid request parameter
		zap.L().Error(errAccompanyingGuestLessThanZero.Error(), zap.Error(errAccompanyingGuestLessThanZero))
		return c.JSON(http.StatusBadRequest, presenter.ErrResp(reqID, errAccompanyingGuestLessThanZero))
	}
	res := putGuestArrivesResponse{}
	// Query database
	err = con.dbSvc.GuestArrival(c.Request().Context(), r.AccompanyingGuests, name)
	if err != nil {
		// Error while querying database
		if err.Error() == "guest did not register" || err.Error() == "guest already arrived" {
			return c.JSON(http.StatusNotFound, presenter.ErrResp(reqID, err))
		}
		return c.JSON(http.StatusInternalServerError, presenter.ErrResp(reqID, err))
	}
	// Map response fields
	res.Name = name
	// Return ok
	return c.JSON(http.StatusCreated, res)
}

// ListArrivedGuest handles GET /guests
func (con *GuestHandler) ListArrivedGuest(c echo.Context) (err error) {
	reqID := c.Response().Header().Get(echo.HeaderXRequestID)
	res := getGuestListResponse{}
	limit, offset, err := getLimitAndOffest(c)
	if err != nil {
		zap.L().Error(errInvalidRequest.Error(), zap.Error(err))
		return c.JSON(http.StatusBadRequest, presenter.ErrResp(reqID, err))
	}
	// Query database
	data, err := con.dbSvc.ListArrivedGuests(c.Request().Context(), limit, offset)
	if err != nil {
		// Error while querying database
		return c.JSON(http.StatusInternalServerError, presenter.ErrResp(reqID, err))
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
func (con *GuestHandler) GuestDepart(c echo.Context) (err error) {
	reqID := c.Response().Header().Get(echo.HeaderXRequestID)
	// Get and validate request parameter
	name := c.Param("name")
	// Query database
	err = con.dbSvc.GuestDepart(c.Request().Context(), name)
	if err != nil {
		// Error while querying database
		if err == errGuestNotFound || err == errTableNotFound {
			return c.JSON(http.StatusNotFound, presenter.ErrResp(reqID, err))
		}
		return c.JSON(http.StatusInternalServerError, presenter.ErrResp(reqID, err))
	}
	// Return ok
	return c.JSON(http.StatusAccepted, "OK!")
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
			return 0, 0, err
		}
	}
	if stroffset != "" {
		offset, err = strconv.ParseInt(stroffset, 10, 64)
		if err != nil {
			return 0, 0, err
		}
	}

	return limit, offset, nil
}

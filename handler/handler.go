package handler

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"

	"go.uber.org/zap"

	"ggv2/constant"
	"ggv2/entities"
	"ggv2/handler/presenter"
	"ggv2/repo"
	"ggv2/services"
)

type errResp struct {
	ReqId  string `json:"requestID"`
	ErrMsg string `json:"message"`
}

type Handler struct {
	DbSvc services.Service
}

func NewHandler(conn *sqlx.DB) *Handler {
	dbRepo := repo.NewDbRepo(conn)
	dbSvc := services.NewDbService(dbRepo)

	return &Handler{
		DbSvc: dbSvc,
	}
}

// GetTable handles GET /table/:id"
func (h *Handler) GetTable(c echo.Context) (err error) {
	ctx := c.Request().Context()
	reqID := c.Response().Header().Get(echo.HeaderXRequestID)
	ctx = context.WithValue(ctx, constant.ContextKeyRequestID, reqID)
	id := c.Param("id")
	logger := zap.L().With(zap.String("rqId", fmt.Sprintf("%v", reqID)))
	logRequest(c, "GetTable")
	tableId, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		errResp := &errResp{
			ReqId:  reqID,
			ErrMsg: err.Error(),
		}
		// Invalid request parameter
		logger.Error(constant.ErrInvalidRequest.Error(), zap.Error(err))
		return c.JSON(http.StatusBadRequest, errResp)
	}
	// Query database
	res, err := h.DbSvc.GetTable(ctx, tableId)
	if err != nil {
		// Error while querying database
		errResp := &errResp{
			ReqId:  reqID,
			ErrMsg: err.Error(),
		}
		return c.JSON(http.StatusInternalServerError, errResp)
	}
	// Return ok
	return c.JSON(http.StatusOK, res)
}

// GetTables handles GET /tables
func (con *Handler) GetTables(c echo.Context) (err error) {

	limit := c.QueryParam("limit")
	offset := c.QueryParam("offset")
	ctx := c.Request().Context()
	reqID := c.Response().Header().Get(echo.HeaderXRequestID)
	ctx = context.WithValue(ctx, constant.ContextKeyRequestID, reqID)
	var iLimit int64 = 10
	var iOffset int64
	logger := zap.L().With(zap.String("rqId", fmt.Sprintf("%v", reqID)))
	logRequest(c, "GetTables")
	if limit != "" {
		iLimit, err = strconv.ParseInt(limit, 10, 64)
		if err != nil {
			errResp := &errResp{
				ReqId:  reqID,
				ErrMsg: err.Error(),
			}
			logger.Error(constant.ErrInvalidRequest.Error(), zap.Error(err))
			return c.JSON(http.StatusBadRequest, errResp)
		}
	}
	if offset != "" {
		iOffset, err = strconv.ParseInt(offset, 10, 64)
		if err != nil {
			errResp := &errResp{
				ReqId:  reqID,
				ErrMsg: err.Error(),
			}
			logger.Error(constant.ErrInvalidRequest.Error(), zap.Error(err))
			return c.JSON(http.StatusBadRequest, errResp)
		}
	}
	// Query database
	res, err := con.DbSvc.ListTables(ctx, iLimit, iOffset)

	if err != nil {
		// Error while querying database
		errResp := &errResp{
			ReqId:  reqID,
			ErrMsg: err.Error(),
		}
		return c.JSON(http.StatusInternalServerError, errResp)
	}
	// Return ok
	return c.JSON(http.StatusOK, res)
}

// CreateTables handles PUT /table
func (con *Handler) CreateTable(c echo.Context) (err error) {
	type createTableRequest struct {
		Capacity int64 `json:"capacity" form:"capacity"`
	}
	type PutCreateTableResponse struct {
		Table *entities.Table `json:"table"`
	}
	ctx := c.Request().Context()
	reqID := c.Response().Header().Get(echo.HeaderXRequestID)
	ctx = context.WithValue(ctx, constant.ContextKeyRequestID, reqID)
	logger := zap.L().With(zap.String("rqId", fmt.Sprintf("%v", reqID)))
	r := new(createTableRequest)
	logRequest(c, "CreateTable")
	if err = c.Bind(r); err != nil {
		// Invalid request parameter
		errResp := &errResp{
			ReqId:  reqID,
			ErrMsg: constant.ErrInvalidRequest.Error(),
		}
		logger.Error(constant.ErrInvalidRequest.Error(), zap.Error(err))
		return c.JSON(http.StatusBadRequest, errResp)
	}
	if r.Capacity < 1 {
		// Invalid request parameter
		errResp := &errResp{
			ReqId:  reqID,
			ErrMsg: constant.ErrCapacityLessThanOne.Error(),
		}
		logger.Error(constant.ErrInvalidRequest.Error(), zap.Error(constant.ErrCapacityLessThanOne))
		return c.JSON(http.StatusBadRequest, errResp)
	}
	res := &PutCreateTableResponse{}
	// Query database
	table, err := con.DbSvc.CreateTable(ctx, r.Capacity)
	if err != nil {
		// Error while querying database
		errResp := &errResp{
			ReqId:  reqID,
			ErrMsg: err.Error(),
		}
		return c.JSON(http.StatusInternalServerError, errResp)
	}
	// Map response fields
	res.Table = table
	// Return ok
	return c.JSON(http.StatusCreated, res)
}

// Init handles GET /init
func (con *Handler) EmptyTables(c echo.Context) (err error) {
	ctx := c.Request().Context()
	reqID := c.Response().Header().Get(echo.HeaderXRequestID)
	ctx = context.WithValue(ctx, constant.ContextKeyRequestID, reqID)
	logRequest(c, "EmptyTables")
	// Query database
	err = con.DbSvc.EmptyTables(ctx)
	if err != nil {
		// Error while querying database
		errResp := &errResp{
			ReqId:  reqID,
			ErrMsg: err.Error(),
		}
		return c.JSON(http.StatusInternalServerError, errResp)
	}
	// Return ok
	return c.JSON(http.StatusOK, "Tables emptied!")
}

// GetEmptySeatsCount handles GET /seats_empty
func (con *Handler) GetEmptySeatsCount(c echo.Context) (err error) {
	type GetSeatsEmptyResponse struct {
		SeatsEmpty int64 `json:"seats_empty"`
	}
	res := GetSeatsEmptyResponse{}
	ctx := c.Request().Context()
	reqID := c.Response().Header().Get(echo.HeaderXRequestID)
	ctx = context.WithValue(ctx, constant.ContextKeyRequestID, reqID)
	logRequest(c, "GetEmptySeatsCount")
	// Query database
	count, err := con.DbSvc.GetEmptySeatsCount(ctx)
	if err != nil {
		// Error while querying database
		errResp := &errResp{
			ReqId:  reqID,
			ErrMsg: err.Error(),
		}
		return c.JSON(http.StatusInternalServerError, errResp)
	}
	// Map response fields
	res.SeatsEmpty = int64(count)
	// Return ok
	return c.JSON(http.StatusOK, res)
}

// AddToGuestList handles POST /guest_list/:name
func (con *Handler) AddToGuestList(c echo.Context) (err error) {
	type PostGuestListRequest struct {
		Table              int64 `json:"table" form:"table"`
		AccompanyingGuests int64 `json:"accompanying_guests" form:"accompanying_guests"`
	}
	type PostGuestListResponse struct {
		Name string `json:"name"`
	}
	// Get and validate request parameter
	r := new(PostGuestListRequest)
	n := c.Param("name")
	ctx := c.Request().Context()
	reqID := c.Response().Header().Get(echo.HeaderXRequestID)
	ctx = context.WithValue(ctx, constant.ContextKeyRequestID, reqID)
	logger := zap.L().With(zap.String("rqId", fmt.Sprintf("%v", reqID)))
	logRequest(c, "AddToGuestList")
	if err = c.Bind(r); err != nil {
		// Invalid request parameter
		errResp := &errResp{
			ReqId:  reqID,
			ErrMsg: err.Error(),
		}
		logger.Error(constant.ErrInvalidRequest.Error(), zap.Error(err))
		return c.JSON(http.StatusBadRequest, errResp)
	}
	if r.Table < 1 {
		// Invalid request parameter
		errResp := &errResp{
			ReqId:  reqID,
			ErrMsg: constant.ErrTableNotFound.Error(),
		}
		logger.Error(constant.ErrInvalidRequest.Error(), zap.Error(err))
		return c.JSON(http.StatusBadRequest, errResp)
	}
	if r.AccompanyingGuests < 0 {
		// Invalid request parameter
		errResp := &errResp{
			ReqId:  reqID,
			ErrMsg: constant.ErrAccompanyingGuestLessThanZero.Error(),
		}
		logger.Error(constant.ErrInvalidRequest.Error(), zap.Error(constant.ErrAccompanyingGuestLessThanZero))
		return c.JSON(http.StatusBadRequest, errResp)
	}
	res := PostGuestListResponse{}
	// Query database
	err = con.DbSvc.AddToGuestList(ctx, r.AccompanyingGuests, r.Table, n)
	if err != nil {
		errResp := &errResp{
			ReqId:  reqID,
			ErrMsg: err.Error(),
		}
		// Error while querying database
		if err.Error() == "table not found" {
			
			return c.JSON(http.StatusNotFound, errResp)
		}
		return c.JSON(http.StatusInternalServerError, errResp)
	}
	// Map response fields
	res.Name = n
	// Return ok
	return c.JSON(http.StatusCreated, res)
}

// Healthcheck handles GET /
func (con *Handler) Ping(c echo.Context) (err error) {
	// Server is up and running, return OK!
	return c.String(http.StatusOK, "Pong")
}

// GetGuestList handles GET /guest_list
func (con *Handler) GetGuestList(c echo.Context) (err error) {
	type GetGuestListResponse struct {
		Guests []*presenter.Guest `json:"guests"`
	}
	ctx := c.Request().Context()
	reqID := c.Response().Header().Get(echo.HeaderXRequestID)
	ctx = context.WithValue(ctx, constant.ContextKeyRequestID, reqID)
	res := GetGuestListResponse{}
	limit := c.QueryParam("limit")
	offset := c.QueryParam("offset")
	logger := zap.L().With(zap.String("rqId", fmt.Sprintf("%v", reqID)))
	logRequest(c, "GetGuestList")
	
	var iLimit int64 = 10
	var iOffset int64
	if limit != "" {
		iLimit, err = strconv.ParseInt(limit, 10, 64)
		if err != nil {
			errResp := &errResp{
				ReqId:  reqID,
				ErrMsg: err.Error(),
			}
			logger.Error(constant.ErrInvalidRequest.Error(), zap.Error(err))
			return c.JSON(http.StatusBadRequest, errResp)
		}
	}
	if offset != "" {
		iOffset, err = strconv.ParseInt(offset, 10, 64)
		if err != nil {
			errResp := &errResp{
				ReqId:  reqID,
				ErrMsg: err.Error(),
			}
			logger.Error(constant.ErrInvalidRequest.Error(), zap.Error(err))
			return c.JSON(http.StatusBadRequest, errResp)
		}
	}
	// Query database
	data, err := con.DbSvc.ListRSVPGuests(ctx, iLimit, iOffset)
	if err != nil {
		// Error while querying database
		errResp := &errResp{
			ReqId:  reqID,
			ErrMsg: err.Error(),
		}
		return c.JSON(http.StatusInternalServerError, errResp)
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
	type PutGuestArrivesRequest struct {
		AccompanyingGuests int64 `json:"accompanying_guests" form:"accompanying_guests"`
	}
	type PutGuestArrivesResponse struct {
		Name string `json:"name"`
	}
	ctx := c.Request().Context()
	reqID := c.Response().Header().Get(echo.HeaderXRequestID)
	ctx = context.WithValue(ctx, constant.ContextKeyRequestID, reqID)
	logger := zap.L().With(zap.String("rqId", fmt.Sprintf("%v", reqID)))
	logRequest(c, "GuestArrived")
	// Get and validate request parameter
	r := new(PutGuestArrivesRequest)
	name := c.Param("name")
	if err = c.Bind(r); err != nil {
		// Invalid request parameter
		errResp := &errResp{
			ReqId:  reqID,
			ErrMsg: err.Error(),
		}
		logger.Error(constant.ErrInvalidRequest.Error(), zap.Error(err))
		return c.JSON(http.StatusBadRequest, errResp)
	}
	if r.AccompanyingGuests < 0 {
		// Invalid request parameter
		errResp := &errResp{
			ReqId:  reqID,
			ErrMsg: constant.ErrAccompanyingGuestLessThanZero.Error(),
		}
		logger.Error(constant.ErrAccompanyingGuestLessThanZero.Error(), zap.Error(constant.ErrAccompanyingGuestLessThanZero))
		return c.JSON(http.StatusBadRequest, errResp)
	}
	res := PutGuestArrivesResponse{}
	// Query database
	err = con.DbSvc.GuestArrival(ctx, r.AccompanyingGuests, name)
	if err != nil {
		// Error while querying database
		errResp := &errResp{
			ReqId:  reqID,
			ErrMsg: err.Error(),
		}
		if err.Error() == "guest did not register" || err.Error() == "guest already arrived" {
			return c.JSON(http.StatusNotFound, errResp)
		}
		return c.JSON(http.StatusInternalServerError, errResp)
	}
	// Map response fields
	res.Name = name
	// Return ok
	return c.JSON(http.StatusCreated, res)
}

// ListArrivedGuest handles GET /guests
func (con *Handler) ListArrivedGuest(c echo.Context) (err error) {
	type GetGuestListResponse struct {
		Guests []*presenter.Guest `json:"guests"`
	}
	ctx := c.Request().Context()
	reqID := c.Response().Header().Get(echo.HeaderXRequestID)
	ctx = context.WithValue(ctx, constant.ContextKeyRequestID, reqID)
	res := GetGuestListResponse{}
	limit := c.QueryParam("limit")
	offset := c.QueryParam("offset")
	logger := zap.L().With(zap.String("rqId", fmt.Sprintf("%v", reqID)))
	logRequest(c, "GuestArrived")
	var iLimit int64 = 10
	var iOffset int64
	if limit != "" {
		iLimit, err = strconv.ParseInt(limit, 10, 64)
		if err != nil {
			errResp := &errResp{
				ReqId:  reqID,
				ErrMsg: err.Error(),
			}
			logger.Error(constant.ErrInvalidRequest.Error(), zap.Error(err))
			return c.JSON(http.StatusBadRequest, errResp)
		}
	}
	if offset != "" {
		iOffset, err = strconv.ParseInt(offset, 10, 64)
		if err != nil {
			errResp := &errResp{
				ReqId:  reqID,
				ErrMsg: err.Error(),
			}
			logger.Error(constant.ErrInvalidRequest.Error(), zap.Error(err))
			return c.JSON(http.StatusBadRequest, errResp)
		}
	}
	// Query database
	data, err := con.DbSvc.ListArrivedGuests(ctx, iLimit, iOffset)
	if err != nil {
		// Error while querying database
		errResp := &errResp{
			ReqId:  reqID,
			ErrMsg: err.Error(),
		}
		return c.JSON(http.StatusInternalServerError, errResp)
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
	ctx := c.Request().Context()
	reqID := c.Response().Header().Get(echo.HeaderXRequestID)
	ctx = context.WithValue(ctx, constant.ContextKeyRequestID, reqID)
	// logger := zap.L().With(zap.String("rqId", fmt.Sprintf("%v", reqID)))
	logRequest(c, "GuestDepart")
	// Get and validate request parameter
	name := c.Param("name")
	// Query database
	err = con.DbSvc.GuestDepart(ctx, name)
	if err != nil {
		errResp := &errResp{
			ReqId:  reqID,
			ErrMsg: err.Error(),
		}
		// Error while querying database
		if err.Error() == "guest not found" || err.Error() == "table not found" {
			return c.JSON(http.StatusNotFound, errResp)
		}
		return c.JSON(http.StatusInternalServerError, errResp)
	}
	// Return ok
	return c.JSON(http.StatusAccepted, "OK!")
}

func logRequest(c echo.Context, name string){

	zf := []zap.Field{}
	reqID := c.Response().Header().Get(echo.HeaderXRequestID)
	logger := zap.L().With(zap.String("rqId", fmt.Sprintf("%v", reqID)))
	qp := c.QueryParams()
	fp, _ := c.FormParams()
	pn := c.ParamNames()
	pathParams := []string{}

	if c.Request().URL.RawQuery != ""{
		zf = append(zf, zap.String("QueryString", c.Request().URL.RawQuery))
	}
	if fmt.Sprintf("%v", fp) != fmt.Sprintf("%v", qp){
		zf = append(zf, zap.String("FormData", fmt.Sprintf("%s", fp)))
	}
	for _, v := range pn{
		pathParams = append(pathParams, fmt.Sprintf("%s=%s", v, c.Param(v)))
	}
	if len(pathParams) > 0{
		zf = append(zf, zap.String("PathParam", strings.Join(pathParams, "&")))
	}
	logger.Info(name, zf...)
}
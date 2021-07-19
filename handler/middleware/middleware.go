package middleware

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"runtime"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"

	"go.uber.org/zap"
)

type bodyDumpResponseWriter struct {
	io.Writer
	http.ResponseWriter
}

func Middleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) (err error) {

			uuid := uuid.New().String()
			zap.ReplaceGlobals(initLogger())
			logger := zap.L().With(zap.String("rqId", uuid))
			zap.ReplaceGlobals(logger)

			resBody := new(bytes.Buffer)
			mw := io.MultiWriter(c.Response().Writer, resBody)
			writer := &bodyDumpResponseWriter{Writer: mw, ResponseWriter: c.Response().Writer}
			c.Response().Writer = writer

			req := c.Request()
			res := c.Response()
			start := time.Now()

			// Request ID
			rid := req.Header.Get(echo.HeaderXRequestID)
			if rid == "" {
				res.Header().Set(echo.HeaderXRequestID, uuid)
			}

			// CORS
			res.Header().Set(echo.HeaderAccessControlAllowOrigin, "*")
			res.Header().Set(echo.HeaderAccessControlAllowMethods, strings.Join([]string{echo.GET, echo.HEAD, echo.PUT, echo.PATCH, echo.POST, echo.DELETE, echo.OPTIONS}, ","))

			if err = next(c); err != nil {
				c.Error(err)
			}
			stop := time.Now()
			if req.RequestURI != "/metrics" {
				// Log Request
				zf := []zap.Field{}
				qp := c.QueryParams()
				fp, _ := c.FormParams()
				pn := c.ParamNames()
				pathParams := []string{}
				zf = append(zf, zap.String("method", req.Method))
				zf = append(zf, zap.String("uri", req.RequestURI))
				if c.Request().URL.RawQuery != "" {
					zf = append(zf, zap.String("QueryString", c.Request().URL.RawQuery))
				}
				if fmt.Sprintf("%v", fp) != fmt.Sprintf("%v", qp) {
					zf = append(zf, zap.String("FormData", fmt.Sprintf("%s", fp)))
				}
				for _, v := range pn {
					pathParams = append(pathParams, fmt.Sprintf("%s=%s", v, c.Param(v)))
				}
				if len(pathParams) > 0 {
					zf = append(zf, zap.String("PathParam", strings.Join(pathParams, "&")))
				}
				logger.Info("RqLog:", zf...)

				// Log Response
				zf = []zap.Field{}
				zf = append(zf, zap.String("status", fmt.Sprintf("%d", res.Status)))
				zf = append(zf, zap.String("latency", stop.Sub(start).String()))
				zf = append(zf, zap.String("rsBody", resBody.String()))

				logger.Info("RsLog:", zf...)

			}

			// Recover
			defer func() {
				if r := recover(); r != nil {
					err, ok := r.(error)
					if !ok {
						err = fmt.Errorf("%v", r)
					}
					stack := make([]byte, 4<<10)
					length := runtime.Stack(stack, !false)
					msg := fmt.Sprintf("[PANIC RECOVER] %v %s\n", err, stack[:length])
					c.Logger().Print(msg)
					c.Error(err)
				}
			}()

			return
		}
	}
}

func (w *bodyDumpResponseWriter) WriteHeader(code int) {
	w.ResponseWriter.WriteHeader(code)
}

func (w *bodyDumpResponseWriter) Write(b []byte) (int, error) {
	return w.Writer.Write(b)
}

func (w *bodyDumpResponseWriter) Flush() {
	w.ResponseWriter.(http.Flusher).Flush()
}

func (w *bodyDumpResponseWriter) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	return w.ResponseWriter.(http.Hijacker).Hijack()
}

func initLogger() *zap.Logger {
	const (
		logPath = "./logs/ggv2.log"
	)
	os.OpenFile(logPath, os.O_RDONLY|os.O_CREATE, 0666)
	c := zap.NewProductionConfig()
	c.OutputPaths = []string{"stdout", logPath}
	l, err := c.Build()
	if err != nil {
		panic(err)
	}
	return l
}

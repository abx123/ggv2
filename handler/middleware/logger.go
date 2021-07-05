package middleware

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

type bodyDumpResponseWriter struct {
	io.Writer
	http.ResponseWriter
}

func Logger() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) (err error) {

			resBody := new(bytes.Buffer)
			mw := io.MultiWriter(c.Response().Writer, resBody)
			writer := &bodyDumpResponseWriter{Writer: mw, ResponseWriter: c.Response().Writer}
			c.Response().Writer = writer

			req := c.Request()
			res := c.Response()
			start := time.Now()
			logger := zap.L().With(zap.String("rqId", fmt.Sprintf("%v", req.Header.Get(echo.HeaderXRequestID))))

			if err = next(c); err != nil {
				c.Error(err)
			}
			stop := time.Now()
			// fmt.Printf("id=%s, method=%s, uri=%s, status=%d, latency=%s, resBody=%s\n", res.Header().Get(echo.HeaderXRequestID), req.Method, req.RequestURI, res.Status, stop.Sub(start).String(), resBody.String())

			// request
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

			// response
			zf = []zap.Field{}
			zf = append(zf, zap.String("status", fmt.Sprintf("%d", res.Status)))
			zf = append(zf, zap.String("latency", stop.Sub(start).String()))
			zf = append(zf, zap.String("rsBody", resBody.String()))

			logger.Info("RsLog:", zf...)


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

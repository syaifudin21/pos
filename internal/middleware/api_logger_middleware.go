package middleware

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/msyaifudin/pos/pkg/elasticsearch"
)

type bodyDumpResponseWriter struct {
	http.ResponseWriter
	body *bytes.Buffer
}

func (w *bodyDumpResponseWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

func APILoggerMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		start := time.Now()

		// Baca request body
		var reqBody string
		if c.Request().Body != nil {
			bodyBytes, _ := ioutil.ReadAll(c.Request().Body)
			reqBody = string(bodyBytes)
			// Reset body supaya handler tetap bisa baca
			c.Request().Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))
		}

		// Bungkus response writer
		origWriter := c.Response().Writer
		buf := new(bytes.Buffer)
		c.Response().Writer = &bodyDumpResponseWriter{ResponseWriter: origWriter, body: buf}

		err := next(c)

		// Ambil response body
		respBody := buf.String()

		// Ambil user dari context jika ada
		user := ""
		if u := c.Get("user"); u != nil {
			user = toString(u)
		}

		dur := time.Since(start).Milliseconds()

		logData := elasticsearch.APILog{
			Method:     c.Request().Method,
			Path:       c.Path(),
			Status:     c.Response().Status,
			DurationMs: dur,
			User:       user,
			Error:      errString(err),
			Extra: map[string]interface{}{
				"request":  reqBody,
				"response": respBody,
			},
		}
		elasticsearch.LogAPI(logData)

		return err
	}
}

func toString(v interface{}) string {
	switch t := v.(type) {
	case string:
		return t
	case []byte:
		return string(t)
	default:
		return ""
	}
}

func errString(err error) string {
	if err == nil {
		return ""
	}
	return err.Error()
}

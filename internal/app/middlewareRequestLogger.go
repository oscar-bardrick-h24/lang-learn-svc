package app

import (
	"net/http"
	"time"

	"github.com/sirupsen/logrus"
)

type responseWriterRecorder struct {
	http.ResponseWriter
	status int
	body   string
}

func (w *responseWriterRecorder) WriteHeader(statusCode int) {
	w.status = statusCode
	w.ResponseWriter.WriteHeader(statusCode)
}

func (w *responseWriterRecorder) Write(bs []byte) (int, error) {
	w.body += string(bs)
	return w.ResponseWriter.Write(bs)
}

func (app *App) MiddlewareRequestResponseLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		reqID := r.Header.Get(reqIDHeader)

		app.logger.WithFields(
			logrus.Fields{
				"type":             "request",
				"path":             r.URL.Path,
				"query_parameters": r.URL.RawQuery,
				"method":           r.Method,
				"request_id":       reqID,
				"src_host":         r.RemoteAddr,
			},
		).Info("incoming request")

		start := time.Now().UTC()

		customResponseWriter := &responseWriterRecorder{ResponseWriter: w}
		next.ServeHTTP(customResponseWriter, r)

		end := time.Now().UTC()

		respLogFields := logrus.Fields{
			"type":          "response",
			"status":        customResponseWriter.status,
			"response_time": end.Sub(start).Milliseconds(),
			"path":          r.URL.Path,
			"method":        r.Method,
			"request_id":    reqID,
			"tgt_host":      r.Host,
		}

		if app.env == "development" {
			respLogFields["body"] = customResponseWriter.body
		}

		app.logger.WithFields(respLogFields).Info("outgoing response")
	})
}

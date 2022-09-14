package app

import (
	"net/http"

	"github.com/OJOMB/graffiti-berlin-svc/internal/pkg/contextual"
)

const reqIDHeader = "X-Request-ID"

func (app *App) MiddlewareRequestID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		reqID := r.Header.Get(reqIDHeader)
		if reqID == "" {
			generatedID, err := app.idTool.New()
			if err != nil {
				app.logger.Errorf("failed to generate request ID: %s", err.Error())
				reqID = "unknown"
			} else {
				reqID = generatedID
			}
		}

		r = r.WithContext(contextual.SetRequestID(r.Context(), reqID))
		next.ServeHTTP(w, r)
	})
}

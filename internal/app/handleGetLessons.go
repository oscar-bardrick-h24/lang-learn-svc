package app

import (
	"encoding/json"
	"net/http"
)

const handleGetLessons = "handleGetLessons"

func (app *App) handleGetLessons() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		lessons, dErr := app.service.GetLessons(r.Context())
		if dErr != nil {
			apperr := app.newAppErrFromDomainErr(dErr)
			http.Error(w, apperr.Error(), apperr.Code())
			return
		}

		respBodyBytes, err := json.Marshal(lessons)
		if err != nil {
			app.logger.WithField(appHandler, handleGetLessons).WithError(err).Error("failed to marshal json response")
			apperr := newAppErr("failed to marshal json response", http.StatusInternalServerError)
			http.Error(w, apperr.Error(), apperr.Code())
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write(respBodyBytes)
	}
}

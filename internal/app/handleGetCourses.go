package app

import (
	"encoding/json"
	"net/http"
)

const handleGetCourses = "handleGetCourses"

func (app *App) handleGetCourses() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		createdByID := r.URL.Query().Get("createdby")

		languages, dErr := app.service.GetCourses(r.Context(), createdByID)
		if dErr != nil {
			apperr := app.newAppErrFromDomainErr(dErr)
			http.Error(w, apperr.Error(), apperr.Code())
			return
		}

		langsBytes, err := json.Marshal(languages)
		if err != nil {
			apperr := newAppErr("failed to marshal languages", http.StatusInternalServerError)
			http.Error(w, apperr.Error(), apperr.Code())
			return
		}

		app.logger.WithField(appHandler, handleGetCourses).Info("all courses successfully retrieved")
		w.Write(langsBytes)
	}
}

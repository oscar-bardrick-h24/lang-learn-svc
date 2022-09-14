package app

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

const handleGetCourse = "handleGetCourse"

func (app *App) handleGetCourse() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		courseID := mux.Vars(r)[urlVarCourseID]

		course, dErr := app.service.GetCourse(r.Context(), courseID)
		if dErr != nil {
			apperr := app.newAppErrFromDomainErr(dErr)
			http.Error(w, apperr.Error(), apperr.Code())
			return
		}

		courseJSON, err := json.Marshal(course)
		if err != nil {
			apperr := newAppErr("could not marshal course", http.StatusInternalServerError)
			http.Error(w, apperr.Error(), apperr.Code())
			return
		}

		app.logger.WithField(appHandler, handleGetCourse).Debugf("course %s retrieved", courseID)
		w.Header().Set("Content-Type", "application/json")
		w.Write(courseJSON)
	}
}

package app

import (
	"net/http"

	"github.com/gorilla/mux"
)

const handleDeleteCourse = "handleDeleteCourse"

func (app *App) handleDeleteCourse() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		courseID := mux.Vars(r)[urlVarCourseID]

		dErr := app.service.DeleteCourse(r.Context(), courseID)
		if dErr != nil {
			apperr := app.newAppErrFromDomainErr(dErr)
			http.Error(w, apperr.Error(), apperr.Code())
			return
		}

		app.logger.WithField(appHandler, handleDeleteCourse).Debugf("course %s successfully deleted", courseID)
		w.WriteHeader(http.StatusOK)
	}
}

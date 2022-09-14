package app

import (
	"net/http"

	"github.com/gorilla/mux"
)

const handleUserEnroll = "handleUserEnroll"

func (app *App) handleUserEnroll() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID := mux.Vars(r)[urlVarUserID]
		courseID := mux.Vars(r)[urlVarCourseID]

		dErr := app.service.EnrollUser(r.Context(), userID, courseID)
		if dErr != nil {
			apperr := app.newAppErrFromDomainErr(dErr)
			http.Error(w, apperr.Error(), apperr.Code())
			return
		}

		app.logger.WithField(appHandler, handleUserEnroll).Debugf("user %s successfully enrolled in course %s", userID, courseID)
		w.WriteHeader(http.StatusOK)
	}
}

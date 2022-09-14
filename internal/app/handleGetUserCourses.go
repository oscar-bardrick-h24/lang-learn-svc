package app

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

const handleGetUserCourses = "handleGetUserCourses"

func (app *App) handleGetUserCourses() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID := mux.Vars(r)[urlVarUserID]

		courses, dErr := app.service.GetUserCourses(r.Context(), userID)
		if dErr != nil {
			apperr := app.newAppErrFromDomainErr(dErr)
			http.Error(w, apperr.Error(), apperr.Code())
			return
		}

		coursesBytes, err := json.Marshal(courses)
		if err != nil {
			apperr := newAppErr("failed to marshal courses", http.StatusInternalServerError)
			http.Error(w, apperr.Error(), apperr.Code())
			return
		}

		app.logger.WithField(appHandler, handleGetUserCourses).Debugf("all enrolled courses successfully retrieved for user %s", userID)
		w.Write(coursesBytes)
	}
}

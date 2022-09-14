package app

import (
	"net/http"

	"github.com/gorilla/mux"
)

const handleDeleteUser = "handleDeleteUser"

func (app *App) handleDeleteUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		userID := vars[urlVarUserID]

		if dErr := app.service.DeleteUser(r.Context(), userID); dErr != nil {
			apperr := app.newAppErrFromDomainErr(dErr)
			http.Error(w, apperr.Error(), apperr.Code())
			return
		}

		app.logger.WithField(appHandler, handleDeleteUser).Infof("user %s deleted", userID)
		w.WriteHeader(http.StatusNoContent)
	}
}

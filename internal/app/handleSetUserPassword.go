package app

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/gorilla/mux"
)

const handleSetUserPassword = "handleSetUserPassword"

type setPasswordReq struct {
	Password string `json:"password"`
}

// handleSetUserPassword handles the request to change a user's password.
func (app *App) handleSetUserPassword() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		userID := vars[urlVarUserID]

		reqBodyBytes, err := ioutil.ReadAll(r.Body)
		if err != nil {
			apperr := newAppErr("request body unreadable", http.StatusBadRequest)
			http.Error(w, apperr.Error(), apperr.Code())
			return
		}

		defer r.Body.Close()

		var passwordData setPasswordReq
		if err := json.Unmarshal(reqBodyBytes, &passwordData); err != nil {
			apperr := newAppErr("invalid json in request body", http.StatusBadRequest)
			http.Error(w, apperr.Error(), apperr.Code())
			return
		}

		if dErr := app.service.SetUserPassword(r.Context(), userID, passwordData.Password); dErr != nil {
			apperr := app.newAppErrFromDomainErr(dErr)
			http.Error(w, apperr.Error(), apperr.Code())
			return
		}

		app.logger.WithField(appHandler, handleSetUserPassword).Debugf("user %s profile pic set", userID)
		w.WriteHeader(http.StatusNoContent)
	}
}

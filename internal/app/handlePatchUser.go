package app

import (
	"io/ioutil"
	"net/http"

	"github.com/gorilla/mux"
)

const handlePatchUser = "handlePatchUser"

// handlePatchUser handles PATCH requests to /users/{id} in accordance with JSON PATCH RFC6902
// https://datatracker.ietf.org/doc/html/rfc6902/
// handlePatchUser will only patch User Attributes. Attempts to patch other User fields like password will be ignored
func (app *App) handlePatchUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		userID := vars[urlVarUserID]

		patch, err := ioutil.ReadAll(r.Body)
		if err != nil {
			apperr := newAppErr("request body unreadable", http.StatusBadRequest)
			http.Error(w, apperr.Error(), apperr.Code())
			return
		}

		defer r.Body.Close()

		if dErr := app.service.PatchUser(r.Context(), userID, patch); dErr != nil {
			app.logger.WithField(appHandler, handlePatchUser).WithError(dErr).Error("failed to update user attributes")
			apperr := app.newAppErrFromDomainErr(dErr)
			http.Error(w, apperr.Error(), apperr.Code())
			return
		}

		// PATCH does not return a body
		// TODO: could implement content negotiation to return resource if requested via Accept header
		w.WriteHeader(http.StatusNoContent)
	}
}

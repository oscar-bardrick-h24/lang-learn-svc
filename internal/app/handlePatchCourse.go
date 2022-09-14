package app

import (
	"io/ioutil"
	"net/http"

	"github.com/gorilla/mux"
)

const handlePatchCourse = "handlePatchCourse"

// handlePatchCourse handles PATCH requests to /users/{id} in accordance with JSON PATCH RFC6902
// https://datatracker.ietf.org/doc/html/rfc6902/
// only patches certain patchable Course fields. Attempts to patch unpatchable fields like 'ID' will be ignored
func (app *App) handlePatchCourse() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		courseID := vars[urlVarCourseID]

		patch, err := ioutil.ReadAll(r.Body)
		if err != nil {
			apperr := newAppErr("request body unreadable", http.StatusBadRequest)
			http.Error(w, apperr.Error(), apperr.Code())
			return
		}

		defer r.Body.Close()

		if dErr := app.service.PatchCourse(r.Context(), courseID, patch); dErr != nil {
			app.logger.WithField(appHandler, handlePatchCourse).WithError(dErr).Error("failed to patch course due to service error")
			apperr := app.newAppErrFromDomainErr(dErr)
			http.Error(w, apperr.Error(), apperr.Code())
			return
		}

		// PATCH does not return a body
		// TODO: could implement content negotiation to return resource if requested via Accept header
		w.WriteHeader(http.StatusNoContent)
	}
}

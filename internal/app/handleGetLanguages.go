package app

import (
	"encoding/json"
	"net/http"
)

const handleGetLanguages = "handleGetLanguages"

func (app *App) handleGetLanguages() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		languages, dErr := app.service.GetLanguages(r.Context())
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

		app.logger.WithField(appHandler, handleGetLanguages).Info("all languages successfully retrieved")
		w.Write(langsBytes)
	}
}

package app

import (
	"net/http"
)

const handleDeleteLanguages = "handleDeleteLanguages"

func (app *App) handleDeleteLanguages() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if dErr := app.service.DeleteLanguages(r.Context()); dErr != nil {
			apperr := app.newAppErrFromDomainErr(dErr)
			http.Error(w, apperr.Error(), apperr.Code())
			return
		}

		app.logger.WithField(appHandler, handleDeleteLanguages).Info("All Languages deleted")
		w.WriteHeader(http.StatusNoContent)
	}
}

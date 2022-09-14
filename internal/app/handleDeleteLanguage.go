package app

import (
	"net/http"

	"github.com/gorilla/mux"
)

const handleDeleteLanguage = "handleDeleteLanguage"

func (app *App) handleDeleteLanguage() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		LangCode := vars[urlVarLanguageCode]

		if dErr := app.service.DeleteLanguage(r.Context(), LangCode); dErr != nil {
			apperr := app.newAppErrFromDomainErr(dErr)
			http.Error(w, apperr.Error(), apperr.Code())
			return
		}

		app.logger.WithField(appHandler, handleDeleteLanguage).Infof("Language %s deleted", LangCode)
		w.WriteHeader(http.StatusNoContent)
	}
}

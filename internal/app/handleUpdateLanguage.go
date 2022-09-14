package app

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/gorilla/mux"
)

const handleUpdateLanguage = "handleUpdateLanguage"

type updateLanguageReq struct {
	Name string `json:"name"`
}

func (app *App) handleUpdateLanguage() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		langCode := vars[urlVarLanguageCode]

		reqBodyBytes, err := ioutil.ReadAll(r.Body)
		if err != nil {
			apperr := newAppErr("request body unreadable", http.StatusBadRequest)
			http.Error(w, apperr.Error(), apperr.Code())
			return
		}

		defer r.Body.Close()

		var langData updateLanguageReq
		if err := json.Unmarshal(reqBodyBytes, &langData); err != nil {
			apperr := newAppErr("invalid json in request body", http.StatusBadRequest)
			http.Error(w, apperr.Error(), apperr.Code())
			return
		}

		// set the profile pic url
		if dErr := app.service.UpdateLanguage(r.Context(), langCode, langData.Name); dErr != nil {
			apperr := app.newAppErrFromDomainErr(dErr)
			http.Error(w, apperr.Error(), apperr.Code())
			return
		}

		app.logger.WithField(appHandler, handleUpdateLanguage).Infof("language %s updated", langCode)
		w.WriteHeader(http.StatusNoContent)
	}
}

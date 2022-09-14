package app

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

const handlerCreateLanguage = "handleCreateLanguage"

type createLanguageReq struct {
	Name string `json:"name"`
	Code string `json:"code"`
}

func (app *App) handleCreateLanguage() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		reqBodyBytes, err := ioutil.ReadAll(r.Body)
		if err != nil {
			apperr := newAppErr("request body unreadable", http.StatusBadRequest)
			http.Error(w, apperr.Error(), apperr.Code())
			return
		}

		defer r.Body.Close()

		var langData createLanguageReq
		if err := json.Unmarshal(reqBodyBytes, &langData); err != nil {
			apperr := newAppErr("invalid json in request body", http.StatusBadRequest)
			http.Error(w, apperr.Error(), apperr.Code())
			return
		}

		if dErr := app.service.CreateLanguage(r.Context(), langData.Code, langData.Name); dErr != nil {
			apperr := app.newAppErrFromDomainErr(dErr)
			http.Error(w, apperr.Error(), apperr.Code())
			return
		}

		app.logger.WithField(appHandler, handlerCreateLanguage).Infof("language %s created", langData.Code)
		w.WriteHeader(http.StatusNoContent)
	}
}

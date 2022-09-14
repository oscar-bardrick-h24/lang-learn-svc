package app

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/gorilla/mux"
)

const handleUpdateLesson = "handleUpdateLesson"

type updateLessonReq struct {
	Title    string `json:"title"`
	Text     string `json:"text"`
	Language string `json:"language"`
}

func (app *App) handleUpdateLesson() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		lessonID := mux.Vars(r)[urlVarLessonID]

		reqBodyBytes, err := ioutil.ReadAll(r.Body)
		if err != nil {
			apperr := newAppErr("request body unreadable", http.StatusBadRequest)
			http.Error(w, apperr.Error(), apperr.Code())
			return
		}

		defer r.Body.Close()

		var lessonData updateLessonReq
		if err := json.Unmarshal(reqBodyBytes, &lessonData); err != nil {
			apperr := newAppErr("invalid json in request body", http.StatusBadRequest)
			http.Error(w, apperr.Error(), apperr.Code())
			return
		}

		if dErr := app.service.UpdateLesson(r.Context(), lessonID, lessonData.Title, lessonData.Text, lessonData.Language); dErr != nil {
			apperr := app.newAppErrFromDomainErr(dErr)
			http.Error(w, apperr.Error(), apperr.Code())
			return
		}

		app.logger.WithField(appHandler, handleUpdateLesson).Infof("lesson %s updated", lessonID)
		w.WriteHeader(http.StatusNoContent)
	}
}

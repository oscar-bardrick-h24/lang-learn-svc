package app

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

const handleGetLesson = "handleGetLesson"

func (app *App) handleGetLesson() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		lessonID := mux.Vars(r)[urlVarLessonID]

		lesson, dErr := app.service.GetLesson(r.Context(), lessonID)
		if dErr != nil {
			apperr := app.newAppErrFromDomainErr(dErr)
			http.Error(w, apperr.Error(), apperr.Code())
			return
		}

		respBodyBytes, err := json.Marshal(lesson)
		if err != nil {
			apperr := newAppErr("could not marshal lesson", http.StatusInternalServerError)
			http.Error(w, apperr.Error(), apperr.Code())
			return
		}

		app.logger.WithField(appHandler, handleGetLesson).Debugf("lesson %s retrieved", lessonID)
		w.Header().Set("Content-Type", "application/json")
		w.Write(respBodyBytes)
	}
}

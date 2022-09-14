package app

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

const handleCreateLesson = "handleCreateLesson"

func (app *App) handleCreateLesson() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		reqBodyBytes, err := ioutil.ReadAll(r.Body)
		if err != nil {
			apperr := newAppErr("request body unreadable", http.StatusBadRequest)
			http.Error(w, apperr.Error(), apperr.Code())
			return
		}

		defer r.Body.Close()

		var lessonData createLessonReq
		if err := json.Unmarshal(reqBodyBytes, &lessonData); err != nil {
			apperr := newAppErr("invalid json in request body", http.StatusBadRequest)
			http.Error(w, apperr.Error(), apperr.Code())
			return
		}

		course, dErr := app.service.CreateLesson(r.Context(), lessonData.Title, lessonData.Text, lessonData.Language)
		if dErr != nil {
			apperr := app.newAppErrFromDomainErr(dErr)
			http.Error(w, apperr.Error(), apperr.Code())
			return
		}

		respBodyBytes, err := json.Marshal(course)
		if err != nil {
			app.logger.WithField(appHandler, handleCreateCourse).WithError(err).Error("failed to marshal json response")
			apperr := newAppErr("failed to marshal json response", http.StatusInternalServerError)
			http.Error(w, apperr.Error(), apperr.Code())
			return
		}

		w.WriteHeader(http.StatusCreated)
		w.Write(respBodyBytes)
	}
}

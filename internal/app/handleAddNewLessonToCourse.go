package app

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/gorilla/mux"
)

const handleAddNewLessonToCourse = "handleAddNewLessonToCourse"

type createLessonReq struct {
	Title    string `json:"title"`
	Text     string `json:"text"`
	Language string `json:"language"`
}

func (app *App) handleAppendNewLessonToCourse() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		courseID := mux.Vars(r)[urlVarCourseID]

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

		newLesson, dErr := app.service.AppendNewLessonToCourse(r.Context(), courseID, lessonData.Title, lessonData.Text, lessonData.Language)
		if dErr != nil {
			apperr := app.newAppErrFromDomainErr(dErr)
			http.Error(w, apperr.Error(), apperr.Code())
			return
		}

		respBodyBytes, err := json.Marshal(newLesson)
		if err != nil {
			app.logger.WithField(appHandler, handlerCreateUser).WithError(err).Error("failed to marshal json response")
			apperr := newAppErr("failed to marshal json response", http.StatusInternalServerError)
			http.Error(w, apperr.Error(), apperr.Code())
			return
		}

		app.logger.WithField(appHandler, handleAddNewLessonToCourse).Debugf("lesson %s successfully added to course %s", newLesson.ID, courseID)
		w.WriteHeader(http.StatusCreated)
		w.Write(respBodyBytes)
	}
}

package app

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/OJOMB/graffiti-berlin-svc/internal/pkg/domain"
)

const handleCreateCourse = "handleCreateCourse"

type createCourseReq struct {
	Title   string          `json:"title"`
	Lessons []domain.Lesson `json:"lessons"`
}

func (app *App) handleCreateCourse() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		reqBodyBytes, err := ioutil.ReadAll(r.Body)
		if err != nil {
			apperr := newAppErr("request body unreadable", http.StatusBadRequest)
			http.Error(w, apperr.Error(), apperr.Code())
			return
		}

		defer r.Body.Close()

		var courseData createCourseReq
		if err := json.Unmarshal(reqBodyBytes, &courseData); err != nil {
			apperr := newAppErr("invalid json in request body", http.StatusBadRequest)
			http.Error(w, apperr.Error(), apperr.Code())
			return
		}

		course, dErr := app.service.CreateCourse(r.Context(), courseData.Title, courseData.Lessons)
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

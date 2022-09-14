package app

import "net/http"

const handleDeleteLesson = "handleDeleteLesson"

func (app *App) handleDeleteLesson() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		lessonID := r.URL.Query().Get(urlVarLessonID)
		if lessonID == "" {
			apperr := newAppErr("missing id query parameter", http.StatusBadRequest)
			http.Error(w, apperr.Error(), apperr.Code())
			return
		}

		dErr := app.service.DeleteLesson(r.Context(), lessonID)
		if dErr != nil {
			apperr := app.newAppErrFromDomainErr(dErr)
			http.Error(w, apperr.Error(), apperr.Code())
			return
		}

		app.logger.WithField(appHandler, handleDeleteLesson).Infof("Lesson %s deleted", lessonID)
		w.WriteHeader(http.StatusNoContent)
	}
}

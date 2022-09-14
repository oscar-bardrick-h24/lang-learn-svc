package app

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

const (
	urlVarLanguageCode = "languageCode"
	urlVarLessonID     = "lessonID"
	urlVarUserID       = "userID"
	urlVarCourseID     = "courseID"
)

func (app *App) routes() {
	// swagger-ui served on /docs endpoint
	// no real need for middlewares here
	if app.env != envProd {
		fs := http.FileServer(http.Dir("./api/OpenAPI/"))
		app.router.PathPrefix("/docs").Handler(http.StripPrefix("/docs", fs))
	}

	///////////////////////////////////////
	// APP ROUTER                        //
	// handles routing app functionality //
	// no authentication required	     //
	///////////////////////////////////////
	appMux := app.router.PathPrefix("").Subrouter()

	appMux.HandleFunc("/auth", app.handleAuthenticate()).Methods(http.MethodPost)
	appMux.HandleFunc("/ping", app.handlePing()).Methods(http.MethodGet)

	appMux.Use(app.MiddlewareRequestID, app.MiddlewareRequestResponseLogger)

	////////////////////////////////////////////////
	// API ROUTER                                 //
	// handles routing api (domain) functionality //
	// some routes are protected by JWT auth      //
	////////////////////////////////////////////////

	// unauthenticated routes
	apiV1MuxUnauthed := app.router.PathPrefix("/v1").Subrouter()
	apiV1MuxUnauthed.HandleFunc("/users", app.handleCreateUser()).Methods(http.MethodPost)

	// authenticated routes
	apiV1MuxAuthed := app.router.PathPrefix("/v1").Subrouter()

	///////////
	// Users //
	///////////
	apiV1MuxAuthed.HandleFunc(fmt.Sprintf("/users/{%s}", urlVarUserID), app.handleGetUser()).Methods(http.MethodGet)
	apiV1MuxAuthed.HandleFunc(fmt.Sprintf("/users/{%s}", urlVarUserID), app.handlePatchUser()).Methods(http.MethodPatch)
	apiV1MuxAuthed.HandleFunc(fmt.Sprintf("/users/{%s}", urlVarUserID), app.handleDeleteUser()).Methods(http.MethodDelete)
	apiV1MuxAuthed.HandleFunc(fmt.Sprintf("/users/{%s}/courses", urlVarUserID), app.handleGetUserCourses()).Methods(http.MethodGet)
	apiV1MuxAuthed.HandleFunc(fmt.Sprintf("/users/{%s}/courses/{%s}", urlVarUserID, urlVarCourseID), app.handleUserEnroll()).Methods(http.MethodPut) // PUT since endpoint is idempotent
	apiV1MuxAuthed.HandleFunc(fmt.Sprintf("/users/{%s}/profilePic", urlVarUserID), app.handleSetProfilePic()).Methods(http.MethodPut)
	apiV1MuxAuthed.HandleFunc(fmt.Sprintf("/users/{%s}/password", urlVarUserID), app.handleSetUserPassword()).Methods(http.MethodPut)

	// Languages
	apiV1MuxAuthed.HandleFunc("/languages", app.handleCreateLanguage()).Methods(http.MethodPost)
	apiV1MuxAuthed.HandleFunc("/languages", app.handleGetLanguages()).Methods(http.MethodGet)
	apiV1MuxAuthed.HandleFunc("/languages", app.handleDeleteLanguages()).Methods(http.MethodDelete)
	apiV1MuxAuthed.HandleFunc(fmt.Sprintf("/languages/{%s}", urlVarLanguageCode), app.handleUpdateLanguage()).Methods(http.MethodPut)
	apiV1MuxAuthed.HandleFunc(fmt.Sprintf("/languages/{%s}", urlVarLanguageCode), app.handleDeleteLanguage()).Methods(http.MethodDelete)

	// Courses
	apiV1MuxAuthed.HandleFunc("/courses", app.handleCreateCourse()).Methods(http.MethodPost)
	apiV1MuxAuthed.HandleFunc("/courses", app.handleGetCourses()).Methods(http.MethodGet)
	apiV1MuxAuthed.HandleFunc(fmt.Sprintf("/courses/{%s}", urlVarCourseID), app.handleGetCourse()).Methods(http.MethodGet)
	apiV1MuxAuthed.HandleFunc(fmt.Sprintf("/courses/{%s}", urlVarCourseID), app.handlePatchCourse()).Methods(http.MethodPatch)
	apiV1MuxAuthed.HandleFunc(fmt.Sprintf("/courses/{%s}", urlVarCourseID), app.handleDeleteCourse()).Methods(http.MethodDelete)
	apiV1MuxAuthed.HandleFunc(fmt.Sprintf("/courses/{%s}/lessons", urlVarCourseID), app.handleAppendNewLessonToCourse()).Methods(http.MethodPost)

	// Lessons
	apiV1MuxAuthed.HandleFunc("/lessons", app.handleGetLessons()).Methods(http.MethodGet)
	apiV1MuxAuthed.HandleFunc(fmt.Sprintf("/lessons/{%s}", urlVarLessonID), app.handleGetLesson()).Methods(http.MethodGet)
	apiV1MuxAuthed.HandleFunc(fmt.Sprintf("/lessons/{%s}", urlVarLessonID), app.handleDeleteLesson()).Methods(http.MethodDelete)
	apiV1MuxAuthed.HandleFunc(fmt.Sprintf("/lessons/{%s}", urlVarLessonID), app.handleUpdateLesson()).Methods(http.MethodPut)

	mws := []mux.MiddlewareFunc{app.MiddlewareRequestID, app.MiddlewareRequestResponseLogger}
	if app.env != envDev {
		mws = append(mws, app.MiddlewareAuth)
	}
	apiV1MuxAuthed.Use(mws...)
}

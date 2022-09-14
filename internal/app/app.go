package app

import (
	"net"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

const (
	componentApp = "App"
	appHandler   = "handler"

	envDev  = "development"
	envStag = "staging"
	envProd = "production"
)

type App struct {
	router  *mux.Router
	logger  *logrus.Entry
	addr    net.Addr
	version string
	env     string

	tokenAuth TokenAuth
	idTool    IDTool
	service   Service
}

func New(r *mux.Router, logger *logrus.Logger, addr net.Addr, version string, env string, auth TokenAuth, idTool IDTool, svc Service) *App {
	return &App{
		router:    r,
		logger:    logger.WithField("component", componentApp),
		addr:      addr,
		version:   version,
		env:       env,
		tokenAuth: auth,
		idTool:    idTool,
		service:   svc,
	}
}

func (app *App) Run() {
	app.routes()
	app.logger.Fatal(http.ListenAndServe(app.addr.String(), app.router))
}

package main

import (
	"database/sql"
	"fmt"
	"net"
	"os"
	"strconv"
	"time"

	"github.com/OJOMB/lang-learn-svc/internal/app"
	"github.com/OJOMB/lang-learn-svc/internal/pkg/auth"
	"github.com/OJOMB/lang-learn-svc/internal/pkg/domain"
	"github.com/OJOMB/lang-learn-svc/internal/pkg/passwords"
	"github.com/OJOMB/lang-learn-svc/internal/pkg/repo"
	"github.com/OJOMB/lang-learn-svc/internal/pkg/uuidv4"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"

	_ "github.com/lib/pq"
)

const (
	versionEnv     = "SVC_VERSION"
	portEnv        = "SVC_PORT"
	environmentEnv = "SVC_ENVIRONMENT"
	dbHostEnv      = "DB_HOST"
	dbPortEnv      = "DB_PORT"
	dbUserEnv      = "DB_USER"
	dbPasswordEnv  = "DB_PASSWORD"

	defaultVersion     = "v0.0.0"
	defaultPort        = 8080
	defaultHost        = "0.0.0.0"
	defaultEnvironment = "dev"
	defaultDBPort      = 3306
	defaultDBUser      = "langlearnsvc"
	defaultDBPassword  = "simple"

	dbName  = "langlearndb"
	appName = "lang-learn-svc"

	passwordGeneratorCost = 15
)

func main() {
	logger := logrus.New()
	logger.SetFormatter(&logrus.JSONFormatter{})

	version := os.Getenv(versionEnv)
	if version == "" {
		logger.Info("failed to retrieve app version number from env...using default")
		version = defaultVersion
	}

	logger.Infof("app version number %s", version)

	var appPort int
	var err error
	portStr := os.Getenv(portEnv)
	if portStr == "" {
		logger.Infof("failed to retrieve app port number from env...using default %d", defaultPort)
		appPort = defaultPort
	} else {
		appPort, err = strconv.Atoi(portStr)
		if err != nil {
			logger.WithError(err).Fatalf("retrieved invalid service port number from env: %s", portStr)
		}
	}

	environment := os.Getenv(environmentEnv)
	if environment == "" {
		logger.Info("failed to retrieve app environment from env...using default")
		environment = defaultEnvironment
	}

	dbHost := os.Getenv(dbHostEnv)
	if version == "" {
		logger.Info("failed to retrieve DB host from env...using default")
		dbHost = defaultHost
	}

	var dbPort int
	dbPortStr := os.Getenv(dbPortEnv)
	if dbPortStr == "" {
		logger.Info("failed to retrieve DB host from env...using default")
		dbPort = defaultDBPort
	} else {
		dbPort, err = strconv.Atoi(dbPortStr)
		if err != nil {
			logger.WithError(err).Fatalf("retrieved invalid DB port number from env: %s", portStr)
		}
	}

	dbUser := os.Getenv(dbUserEnv)
	if dbUser == "" {
		logger.Info("failed to retrieve DB host from env...using default")
		dbUser = defaultDBUser
	}

	dbPassword := os.Getenv(dbPasswordEnv)
	if dbUser == "" {
		logger.Info("failed to retrieve DB host from env...using default")
		dbPassword = defaultDBPassword
	}

	logger.Infof("connecting to DB @ %s:%d as %s", dbHost, dbPort, dbUser)

	// default loc - so UTC
	dbCnxnStr := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable", dbUser, dbPassword, dbHost, dbPort, dbName)
	db, err := sql.Open("postgres", dbCnxnStr)
	if err != nil {
		logger.WithError(err).Fatalf("failed to establish connection to DB @ %s:%d", dbHost, dbPort)
	}

	defer db.Close()

	logger.Info("successfully connected to DB")

	uuidTool := uuidv4.NewGenerator()

	server := app.New(
		mux.NewRouter(),
		logger,
		&net.TCPAddr{IP: net.ParseIP(defaultHost), Port: appPort},
		version,
		environment,
		// obviously hardcoding JWT secret is not secure, just for demo purposes
		auth.NewJWTTool("secretKey", 12*time.Hour, appName, uuidTool),
		uuidTool,
		domain.NewService(
			logger,
			repo.NewSQLRepo(db, logger),
			uuidTool,
			passwords.NewGenerator(passwordGeneratorCost),
		),
	)

	server.Run()
}

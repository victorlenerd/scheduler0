package http_server

import (
	"github.com/go-http-utils/logger"
	"github.com/go-pg/pg"
	"github.com/gorilla/mux"
	"github.com/robfig/cron"
	"github.com/unrolled/secure"
	"log"
	"net/http"
	"os"
	"scheduler0/server/db"
	"scheduler0/server/http_server/controllers/credential"
	"scheduler0/server/http_server/controllers/execution"
	"scheduler0/server/http_server/controllers/job"
	"scheduler0/server/http_server/controllers/project"
	"scheduler0/server/http_server/middlewares"
	"scheduler0/server/process"
	"scheduler0/utils"
)

// Start this will start the http server
func Start() {
	conn, err := db.OpenConnection()
	if err != nil {
		panic(err)
	}

	dbConnection := conn.(*pg.DB)

	// SetupDB logging
	log.SetFlags(0)
	log.SetOutput(new(utils.LogWriter))
	jobProcessor := process.JobProcessor{
		DBConnection: dbConnection,
		Cron: cron.New(),
		RecoveredJobs: []process.RecoveredJob{},
	}

	// Set time zone, create database and run db
	db.CreateModelTables(dbConnection)

	// StartJobs process to execute cron-server jobs
	go jobProcessor.StartJobs()

	// HTTP router setup
	router := mux.NewRouter()

	// Security middleware
	secureMiddleware := secure.New(secure.Options{FrameDeny: true})

	// Initialize controllers
	executionController := execution.Controller{DBConnection: dbConnection}
	jobController := job.Controller{
		DBConnection: dbConnection,
		JobProcessor: &jobProcessor,
	}
	projectController := project.Controller{DBConnection: dbConnection}
	credentialController := credential.Controller{DBConnection: dbConnection}

	// Mount middleware
	middleware := middlewares.MiddlewareType{}

	router.Use(secureMiddleware.Handler)
	router.Use(mux.CORSMethodMiddleware(router))
	router.Use(middleware.ContextMiddleware)
	router.Use(middleware.AuthMiddleware(dbConnection))

	// Executions Endpoint
	router.HandleFunc("/executions", executionController.List).Methods(http.MethodGet)

	// Credentials Endpoint
	router.HandleFunc("/credentials", credentialController.CreateOne).Methods(http.MethodPost)
	router.HandleFunc("/credentials", credentialController.List).Methods(http.MethodGet)
	router.HandleFunc("/credentials/{uuid}", credentialController.GetOne).Methods(http.MethodGet)
	router.HandleFunc("/credentials/{uuid}", credentialController.UpdateOne).Methods(http.MethodPut)
	router.HandleFunc("/credentials/{uuid}", credentialController.DeleteOne).Methods(http.MethodDelete)

	// Job Endpoint
	router.HandleFunc("/jobs", jobController.CreateOne).Methods(http.MethodPost)
	router.HandleFunc("/jobs", jobController.List).Methods(http.MethodGet)
	router.HandleFunc("/jobs/{uuid}", jobController.GetOne).Methods(http.MethodGet)
	router.HandleFunc("/jobs/{uuid}", jobController.UpdateOne).Methods(http.MethodPut)
	router.HandleFunc("/jobs/{uuid}", jobController.DeleteOne).Methods(http.MethodDelete)

	// Projects Endpoint
	router.HandleFunc("/projects", projectController.CreateOne).Methods(http.MethodPost)
	router.HandleFunc("/projects", projectController.List).Methods(http.MethodGet)
	router.HandleFunc("/projects/{uuid}", projectController.GetOne).Methods(http.MethodGet)
	router.HandleFunc("/projects/{uuid}", projectController.UpdateOne).Methods(http.MethodPut)
	router.HandleFunc("/projects/{uuid}", projectController.DeleteOne).Methods(http.MethodDelete)

	router.PathPrefix("/api-docs/").Handler(http.StripPrefix("/api-docs/", http.FileServer(http.Dir("./server/http_server/api-docs/"))))

	log.Println("Server is running on port", os.Getenv(utils.PortEnv))
	err = http.ListenAndServe(utils.GetPort(), logger.Handler(router, os.Stdout, logger.CombineLoggerType))
	utils.CheckErr(err)
}

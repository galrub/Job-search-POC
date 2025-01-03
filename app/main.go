package main

import (
	"fmt"
	"net/http"

	"github.com/galrub/go/jobSearch/config"
	"github.com/galrub/go/jobSearch/handlers"
	"github.com/galrub/go/jobSearch/internal/logger"
	"github.com/galrub/go/jobSearch/internal/middleware"
	"github.com/galrub/go/jobSearch/internal/pool"
	_ "github.com/joho/godotenv/autoload"
)

func Favicon(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./web/images/favicon.ico")
}

func main() {
	logger.InitStaticLogger()
	fmt.Println("initializing App")
	pool.InitStatic()

	devMod := config.DevMode()

	mux := http.NewServeMux()
	// static content serving
	staticFs := http.FileServer(http.Dir("./web/static"))
	mux.Handle("GET /static/", http.StripPrefix("/static/", staticFs))
	mux.HandleFunc("GET /favicon.ico", Favicon)

	// jobs fragment
	if devMod {
		println("starting in Dev mode")
		mux.Handle("GET /x/jobs", middleware.CreateCleanChainForFunc(handlers.GetJobListFragment))

		// delete Job
		mux.Handle("DELETE /x/jobs/{id}", middleware.CreateCleanChainForFunc(handlers.DeleteJobWithJobsListFragment))

		// main page
		mux.Handle("GET /", middleware.CreateCleanChainForFunc(handlers.GetMainPage))

		// Login fragment
		mux.Handle("POST /x/login", middleware.CreateCleanChainForFunc(handlers.GetLoginFragment))

		// insert job
		mux.Handle("POST /x/jobs/{id}", middleware.CreateCleanChainForFunc(handlers.CreateOrUpdateJobWithJobsListFragment))

		// job edit, new job
		mux.Handle("GET /x/newJob", middleware.CreateCleanChainForFunc(handlers.PrepareJobEditForNew))

		// job edit exiting job
		mux.Handle("GET /x/jobs/{id}", middleware.CreateCleanChainForFunc(handlers.PrepareJobEditForExiting))

	} else {
		mux.Handle("GET /x/jobs", middleware.CreateSecureChainForFunc(handlers.GetJobListFragment))

		// delete Job
		mux.Handle("DELETE /x/jobs/{id}", middleware.CreateSecureChainForFunc(handlers.DeleteJobWithJobsListFragment))

		// main page
		mux.Handle("GET /", middleware.CreateCleanChainForFunc(handlers.GetMainPage))

		// Login fragment
		mux.Handle("POST /x/login", middleware.CreateCleanChainForFunc(handlers.GetLoginFragment))

		// insert job
		mux.Handle("POST /x/jobs", middleware.CreateSecureChainForFunc(handlers.CreateOrUpdateJobWithJobsListFragment))

		// API
		mux.Handle("POST /api/v1/Login", middleware.CreateSecureChainForFunc(handlers.Login))

		// job edit, new job
		mux.Handle("GET /x/newJob", middleware.CreateSecureChainForFunc(handlers.PrepareJobEditForNew))

		// job edit exiting job
		mux.Handle("GET /x/jobs/{id}", middleware.CreateSecureChainForFunc(handlers.PrepareJobEditForExiting))

	}
	// mux is ready, starting the server
	fmt.Println("Starting Server")
	err := http.ListenAndServe(":8080", mux)
	if err != nil {
		logger.LOG.Err(err).Msg("failed to start server")
	}
}

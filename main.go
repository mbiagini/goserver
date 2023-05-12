package main

import (
	"fmt"
	"goserver/apierrors"
	"goserver/config"
	"goserver/utils/gslog"
	"goserver/utils/gsmiddleware"
	"goserver/utils/gsrender"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	_ "github.com/swaggo/http-swagger/example/go-chi/docs"
)

// @title Go Chi Server
// @version 1.0.0
// @description Servidor que utiliza el framework chi y expone una API REST.

// @contact.name RedFoxSoft
// @contact.email support@redfoxsoft.com
func main() {

	r := chi.NewRouter()

	// Load configuration from external file.
	err := config.LoadConfiguration("./resources/config.json")
	if err != nil {
		fmt.Println("Error found in app configuration")
		fmt.Println(err.Error())
		os.Exit(1)
	}

	// Add middlewares.
	gslog.Server("Setting middlewares")
	r.Use(gsmiddleware.MetricsHandler)
	r.Use(gsmiddleware.TraceID)
	r.Use(gsmiddleware.Recoverer(apierrors.New(apierrors.ERR_NOT_DEFINED), gsrender.JSON))
	r.Use(gsmiddleware.HttpLogHandler)
	
	// Configure routes.
	gslog.Server("Setting routes") 
	Routes(r)

	// Start server.
	gslog.Server("Starting server")
	http.ListenAndServe("127.0.0.1:8080", r)

}
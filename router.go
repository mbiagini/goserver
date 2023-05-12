package main

import (
	"goserver/apierrors"
	"goserver/presentation/controller"
	"goserver/utils/gslog"
	"goserver/utils/gsmiddleware"
	"goserver/utils/gsrender"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func Routes(r *chi.Mux) {

	// Metrics.
	r.Handle("/metrics", promhttp.Handler())

	// Basepath.
	r.Route("/go-server/v1", func(r chi.Router) {

		// Health.
		r.Get("/health", controller.CheckHealth)

		// Users.
		r.Route("/users", func(r chi.Router) {
			r.Post("/", controller.PostUser)
			r.Get("/", controller.GetUsers)
			r.Get("/{id}", controller.GetUserById)
		})

		// Questions.
		r.Route("/identity-validation", func(r chi.Router) {
			r.Post("/questions", controller.GetQuestions)
		})
	})

	// No matching path.
	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		// Trace ID.
		traceID := gsmiddleware.GetTraceID(r.Context())

		err := apierrors.New(apierrors.OPERATION_NOT_DEFINED)
		gslog.ErrorFrom(err, traceID)
		gsrender.WriteJSON(w, http.StatusNotFound, err)
	})
}
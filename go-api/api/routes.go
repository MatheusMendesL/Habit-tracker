package api

import (
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

// define the routes to the project and calls the handle to check
func ControlRoutes() *chi.Mux {
	r := chi.NewRouter()

	r.Use(middleware.Recoverer)
	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)

	StatusRoute(r)
	UserRoutes(r)

	return r
}

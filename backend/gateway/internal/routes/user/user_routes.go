package user

import (
	userHandler "gateway/internal/handlers/user"
	"gateway/internal/middlewares"

	"github.com/go-chi/chi/v5"
)

func UserRoutes(r chi.Router, h *userHandler.UserHandler) {
	r.Route("/user", func(r chi.Router) {
		r.Use(middlewares.RequireAuth)
		r.Get("/GetUserByID/{id}", h.GetUserByID)
	})
}

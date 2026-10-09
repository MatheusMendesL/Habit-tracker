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
		r.Post("/GetUsersByIDs", h.GetUsersByIDs)
		r.Post("/SearchUsers", h.SearchUsers)
		r.Put("/EditUser/{id}", h.EditUser)
		r.Put("/EditPassword/{id}", h.EditPassword)
		r.Delete("/DeleteUser/{id}", h.DeleteUser)
	})
}

package social

import (
	socialHandler "gateway/internal/handlers/social"
	"gateway/internal/middlewares"

	"github.com/go-chi/chi/v5"
)

func SocialRoutes(r chi.Router, h *socialHandler.SocialHandler) {
	r.Route("/social", func(r chi.Router) {
		r.Use(middlewares.RequireAuth)
		r.Post("/StartFollowing", h.StartFollowing)
		r.Post("/Unfollow", h.Unfollow)
		r.Get("/ListFollowers/{id}", h.ListFollowers)
		r.Get("/ListFollowing/{id}", h.ListFollowing)
	})
}

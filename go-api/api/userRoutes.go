package api

import (
	"go-api/response"

	"github.com/go-chi/chi"
)

func UserRoutes(r chi.Router) {
	r.Route("/user", func(r chi.Router) {
		r.Get("/get_all_users", response.GetAllUsers)
		r.Get("/get_by_id/{id:[0-9]+}", response.GetUserById)
	})
}

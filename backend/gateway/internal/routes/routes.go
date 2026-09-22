package routes

import (
	"gateway/internal/middlewares"
	"gateway/internal/routes/habit"
	"gateway/internal/routes/social"
	"gateway/internal/routes/stats"
	"gateway/internal/routes/user"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func ControlRoutes() *chi.Mux {
	r := chi.NewRouter()

	r.Use(middleware.Recoverer)
	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middlewares.Cors)

	r.Route("/api/v1", func(r chi.Router) {
		user.UserRoutes(r)
		stats.StatsRoutes(r)
		social.SocialRoutes(r)
		habit.HabitRoutes(r)
		habit.RoutineRoutes(r)
	})

	return r
}

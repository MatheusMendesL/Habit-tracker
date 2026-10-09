package routes

import (
	habitHandler "gateway/internal/handlers/habit"
	socialHandler "gateway/internal/handlers/social"
	statsHandler "gateway/internal/handlers/stats"
	userHandler "gateway/internal/handlers/user"
	"gateway/internal/middlewares"
	"gateway/internal/routes/habit"
	"gateway/internal/routes/social"
	"gateway/internal/routes/stats"
	"gateway/internal/routes/user"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func ControlRoutes(
	userHandler *userHandler.UserHandler,
	socialHandler *socialHandler.SocialHandler,
	habitHandler *habitHandler.HabitHandler,
	routineHandler *habitHandler.RoutineHandler,
	statsHandler *statsHandler.StatsHandler,
) *chi.Mux {
	r := chi.NewRouter()

	r.Use(middleware.Recoverer)
	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middlewares.Cors)

	r.Route("/api/v1", func(r chi.Router) {
		user.UserRoutes(r, userHandler)
		stats.StatsRoutes(r, statsHandler)
		social.SocialRoutes(r, socialHandler)
		habit.HabitRoutes(r, habitHandler)
		habit.RoutineRoutes(r, routineHandler)
	})

	return r
}

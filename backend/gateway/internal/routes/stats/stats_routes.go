package stats

import (
	statsHandler "gateway/internal/handlers/stats"
	"gateway/internal/middlewares"

	"github.com/go-chi/chi/v5"
)

func StatsRoutes(r chi.Router, h *statsHandler.StatsHandler) {
	r.Group(func(r chi.Router) {
		r.Use(middlewares.RequireAuth)
		r.Post("/stats", h.CreateUserStats)
		r.Get("/stats", h.GetUserStats)
		r.Delete("/stats", h.DeleteUserStats)
		r.Post("/stats/habits/{id}/completions", h.RegisterHabitCompletion)
		r.Delete("/stats/habits/{id}/completions", h.UndoHabitCompletion)
		r.Post("/stats/routines/{id}/completions", h.RegisterRoutineCompletion)
		r.Delete("/stats/routines/{id}/completions", h.UndoRoutineCompletion)
	})
}

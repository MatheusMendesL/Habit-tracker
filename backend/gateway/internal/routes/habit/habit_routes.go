package habit

import (
	habitHandler "gateway/internal/handlers/habit"
	"gateway/internal/middlewares"

	"github.com/go-chi/chi/v5"
)

func HabitRoutes(r chi.Router, h *habitHandler.HabitHandler) {
	r.Route("/habit", func(r chi.Router) {
		r.Use(middlewares.RequireAuth)
		r.Post("/", h.CreateHabit)
		r.Get("/", h.ListHabitsByUser)
		r.Get("/shared/{id}", h.GetSharedHabitByID)
		r.Get("/shared/routine/{id}", h.ListSharedHabitsByRoutine)
		r.Get("/{id}", h.GetHabitByID)
		r.Get("/routine/{id}", h.ListHabitsByRoutine)
		r.Put("/{id}", h.EditHabit)
		r.Delete("/{id}", h.DeleteHabit)
		r.Post("/{id}/complete", h.MarkHabitCompleted)
		r.Delete("/{id}/complete", h.UnmarkHabitCompleted)
		r.Get("/{id}/logs", h.GetHabitLogs)
	})
}

package habit

import (
	habitHandler "gateway/internal/handlers/habit"
	"gateway/internal/middlewares"

	"github.com/go-chi/chi/v5"
)

func RoutineRoutes(r chi.Router, h *habitHandler.RoutineHandler) {
	r.Route("/routine", func(r chi.Router) {
		r.Use(middlewares.RequireAuth)
		r.Post("/", h.CreateRoutine)
		r.Get("/", h.ListRoutinesByUser)
		r.Get("/{id}", h.GetRoutineByID)
		r.Put("/{id}", h.EditRoutine)
		r.Delete("/{id}", h.DeleteRoutine)
		r.Post("/{id}/habits", h.AddHabitToRoutine)
		r.Delete("/{id}/habits/{habitID}", h.RemoveHabitFromRoutine)
		r.Post("/{id}/complete", h.MarkRoutineCompleted)
		r.Delete("/{id}/complete", h.UnmarkRoutineCompleted)
		r.Get("/{id}/logs", h.GetRoutineLogs)
	})
}

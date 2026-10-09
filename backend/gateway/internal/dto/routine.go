package dto

import "time"

type Routine struct {
	ID        string    `json:"id"`
	UserID    string    `json:"user_id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at,omitempty"`
}

type CreateRoutineRequest struct {
	Name string `json:"name"`
}

type EditRoutineRequest struct {
	Name *string `json:"name,omitempty"`
}

type RoutineResponse struct {
	Routine *Routine `json:"routine"`
}

type RoutinesResponse struct {
	Routines []*Routine `json:"routines"`
}

type RoutineActionResponse struct {
	Success bool `json:"success"`
}

type RoutineCompletionRequest struct {
	CompletedAt time.Time `json:"completed_at"`
}

type RoutineLog struct {
	RoutineID   string    `json:"routine_id"`
	CompletedAt time.Time `json:"completed_at"`
}

type RoutineLogsResponse struct {
	Logs []*RoutineLog `json:"logs"`
}

type RoutineHabitRequest struct {
	HabitID string `json:"habit_id"`
}

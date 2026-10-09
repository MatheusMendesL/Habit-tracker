package dto

import "time"

type Habit struct {
	ID          string    `json:"id"`
	UserID      string    `json:"user_id"`
	Name        string    `json:"name"`
	Description string    `json:"description,omitempty"`
	ImageURL    string    `json:"image_url,omitempty"`
	CreatedAt   time.Time `json:"created_at,omitempty"`
}

type CreateHabitRequest struct {
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	ImageURL    string `json:"image_url,omitempty"`
}

type EditHabitRequest struct {
	Name        *string `json:"name,omitempty"`
	Description *string `json:"description,omitempty"`
	ImageURL    *string `json:"image_url,omitempty"`
}

type HabitResponse struct {
	Habit *Habit `json:"habit"`
}

type HabitsResponse struct {
	Habits []*Habit `json:"habits"`
}

type HabitActionResponse struct {
	Success bool `json:"success"`
}

type HabitCompletionRequest struct {
	CompletedAt time.Time `json:"completed_at"`
}

type HabitLog struct {
	HabitID     string    `json:"habit_id"`
	CompletedAt time.Time `json:"completed_at"`
}

type HabitLogsResponse struct {
	Logs []*HabitLog `json:"logs"`
}

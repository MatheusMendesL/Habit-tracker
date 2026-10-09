package dto

import "time"

type UserStats struct {
	UserID               string    `json:"user_id"`
	CompletedHabits      int32     `json:"completed_habits"`
	CompletedRoutines    int32     `json:"completed_routines"`
	CurrentHabitStreak   int32     `json:"current_habit_streak"`
	LongestHabitStreak   int32     `json:"longest_habit_streak"`
	CurrentRoutineStreak int32     `json:"current_routine_streak"`
	LongestRoutineStreak int32     `json:"longest_routine_streak"`
	CreatedAt            time.Time `json:"created_at,omitempty"`
	UpdatedAt            time.Time `json:"updated_at,omitempty"`
}

type UserStatsResponse struct {
	Stats *UserStats `json:"stats"`
}

type StatsActionResponse struct {
	Success bool `json:"success"`
}

type StatsCompletionRequest struct {
	CompletedAt time.Time `json:"completed_at"`
}

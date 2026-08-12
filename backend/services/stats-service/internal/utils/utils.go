package utils

import (
	pbStats "shared/pb/stats"
	"stats-service/db"

	"google.golang.org/protobuf/types/known/timestamppb"
)

func ToProtoStats(stats db.UserStats) *pbStats.UserStats {
	return &pbStats.UserStats{
		UserId:               stats.UserID.String(),
		CompletedHabits:      stats.CompletedHabits,
		CompletedRoutines:    stats.CompletedRoutines,
		CurrentHabitStreak:   stats.CurrentHabitStreak,
		LongestHabitStreak:   stats.LongestHabitStreak,
		CurrentRoutineStreak: stats.CurrentRoutineStreak,
		LongestRoutineStreak: stats.LongestRoutineStreak,
		CreatedAt:            timestamppb.New(stats.CreatedAt),
		UpdatedAt:            timestamppb.New(stats.UpdatedAt),
	}
}

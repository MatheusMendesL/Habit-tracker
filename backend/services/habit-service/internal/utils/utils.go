package utils

import (
	"database/sql"
	"habit-service/db"
	pbHabit "shared/pb/habit"

	"google.golang.org/protobuf/types/known/timestamppb"
)

func ToNullString(param string) sql.NullString {
	return sql.NullString{
		String: param,
		Valid:  true,
	}
}

func NullStringToString(ns sql.NullString) string {
	if ns.Valid {
		return ns.String
	}
	return ""
}

func ToProtoHabit(habit db.Habit) *pbHabit.Habit {
	return &pbHabit.Habit{
		Id:          habit.ID.String(),
		UserId:      habit.UserID.String(),
		Name:        habit.Name,
		Description: NullStringToString(habit.Description),
		ImageUrl:    NullStringToString(habit.ImageUrl),
		CreatedAt:   timestamppb.New(habit.CreatedAt),
	}
}

func ToProtoRoutine(routine db.Routine) *pbHabit.Routine {
	return &pbHabit.Routine{
		Id:        routine.ID.String(),
		UserId:    routine.UserID.String(),
		Name:      routine.Name,
		CreatedAt: timestamppb.New(routine.CreatedAt),
	}
}

func ToProtoLog(log db.GetHabitLogsRow) *pbHabit.HabitLog {
	return &pbHabit.HabitLog{
		HabitId:     log.HabitID.String(),
		CompletedAt: timestamppb.New(log.CompletedAt),
	}
}

func ToProtoRoutineLog(log db.GetRoutineLogsRow) *pbHabit.RoutineLog {
	return &pbHabit.RoutineLog{
		RoutineId:   log.RoutineID.String(),
		CompletedAt: timestamppb.New(log.CompletedAt),
	}
}

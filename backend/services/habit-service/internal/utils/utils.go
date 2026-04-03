package utils

import (
	"database/sql"
	"habit-service/db"
	pbHabit "shared/pb/habit"
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
		Id:          habit.ID,
		UserId:      habit.UserID,
		Name:        habit.Name,
		Description: NullStringToString(habit.Description),
		ImageUrl:    NullStringToString(habit.ImageUrl),
	}
}

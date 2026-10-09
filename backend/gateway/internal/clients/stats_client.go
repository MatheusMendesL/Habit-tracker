package clients

import (
	"context"
	"fmt"
	"time"

	"gateway/internal/dto"
	"gateway/internal/helper"

	"google.golang.org/protobuf/types/known/timestamppb"
	pb "shared/pb/stats"
)

type StatsClient struct {
	client pb.StatsServiceClient
}

func NewStatsClient(client pb.StatsServiceClient) *StatsClient {
	return &StatsClient{client: client}
}

func statsFromProto(stats *pb.UserStats) *dto.UserStats {
	if stats == nil {
		return nil
	}
	result := &dto.UserStats{
		UserID:               stats.UserId,
		CompletedHabits:      stats.CompletedHabits,
		CompletedRoutines:    stats.CompletedRoutines,
		CurrentHabitStreak:   stats.CurrentHabitStreak,
		LongestHabitStreak:   stats.LongestHabitStreak,
		CurrentRoutineStreak: stats.CurrentRoutineStreak,
		LongestRoutineStreak: stats.LongestRoutineStreak,
	}
	if stats.CreatedAt != nil && stats.CreatedAt.IsValid() {
		result.CreatedAt = stats.CreatedAt.AsTime()
	}
	if stats.UpdatedAt != nil && stats.UpdatedAt.IsValid() {
		result.UpdatedAt = stats.UpdatedAt.AsTime()
	}
	return result
}

func (c *StatsClient) CreateUserStats(ctx context.Context, userID string) (*dto.UserStatsResponse, error) {
	res, err := c.client.CreateUserStats(ctx, &pb.CreateUserStatsRequest{UserId: userID})
	if err != nil {
		return nil, err
	}
	if res == nil || res.Stats == nil {
		return nil, fmt.Errorf("empty response from stats service")
	}
	return &dto.UserStatsResponse{Stats: statsFromProto(res.Stats)}, nil
}

func (c *StatsClient) GetUserStats(ctx context.Context, userID string) (*dto.UserStatsResponse, error) {
	res, err := c.client.GetUserStats(ctx, &pb.GetUserStatsRequest{UserId: userID})
	if err != nil {
		return nil, err
	}
	if res == nil || res.Stats == nil {
		return nil, fmt.Errorf("empty response from stats service")
	}
	return &dto.UserStatsResponse{Stats: statsFromProto(res.Stats)}, nil
}

func (c *StatsClient) DeleteUserStats(ctx context.Context, userID string) (*dto.StatsActionResponse, error) {
	res, err := c.client.DeleteUserStats(ctx, &pb.DeleteUserStatsRequest{UserId: userID})
	if err != nil {
		return nil, err
	}
	if res == nil {
		return nil, fmt.Errorf("empty response from stats service")
	}
	return &dto.StatsActionResponse{Success: res.Success}, nil
}

func (c *StatsClient) RegisterHabitCompletion(ctx context.Context, userID, habitID string, completedAt time.Time) (*dto.StatsActionResponse, error) {
	timestamp, err := statsCompletionTimestamp(completedAt)
	if err != nil {
		return nil, err
	}
	res, err := c.client.RegisterHabitCompletion(ctx, &pb.RegisterHabitCompletionRequest{
		UserId: userID, HabitId: habitID, CompletedAt: timestamp,
	})
	if err != nil {
		return nil, err
	}
	if res == nil {
		return nil, fmt.Errorf("empty response from stats service")
	}
	return &dto.StatsActionResponse{Success: res.Success}, nil
}

func (c *StatsClient) UndoHabitCompletion(ctx context.Context, userID, habitID string, completedAt time.Time) (*dto.StatsActionResponse, error) {
	timestamp, err := statsCompletionTimestamp(completedAt)
	if err != nil {
		return nil, err
	}
	res, err := c.client.UndoHabitCompletion(ctx, &pb.UndoHabitCompletionRequest{
		UserId: userID, HabitId: habitID, CompletedAt: timestamp,
	})
	if err != nil {
		return nil, err
	}
	if res == nil {
		return nil, fmt.Errorf("empty response from stats service")
	}
	return &dto.StatsActionResponse{Success: res.Success}, nil
}

func (c *StatsClient) RegisterRoutineCompletion(ctx context.Context, userID, routineID string, completedAt time.Time) (*dto.StatsActionResponse, error) {
	timestamp, err := statsCompletionTimestamp(completedAt)
	if err != nil {
		return nil, err
	}
	res, err := c.client.RegisterRoutineCompletion(ctx, &pb.RegisterRoutineCompletionRequest{
		UserId: userID, RoutineId: routineID, CompletedAt: timestamp,
	})
	if err != nil {
		return nil, err
	}
	if res == nil {
		return nil, fmt.Errorf("empty response from stats service")
	}
	return &dto.StatsActionResponse{Success: res.Success}, nil
}

func (c *StatsClient) UndoRoutineCompletion(ctx context.Context, userID, routineID string, completedAt time.Time) (*dto.StatsActionResponse, error) {
	timestamp, err := statsCompletionTimestamp(completedAt)
	if err != nil {
		return nil, err
	}
	res, err := c.client.UndoRoutineCompletion(ctx, &pb.UndoRoutineCompletionRequest{
		UserId: userID, RoutineId: routineID, CompletedAt: timestamp,
	})
	if err != nil {
		return nil, err
	}
	if res == nil {
		return nil, fmt.Errorf("empty response from stats service")
	}
	return &dto.StatsActionResponse{Success: res.Success}, nil
}

func statsCompletionTimestamp(completedAt time.Time) (*timestamppb.Timestamp, error) {
	if completedAt.IsZero() {
		return nil, fmt.Errorf("completed_at is required")
	}
	return helper.TimestampToProto(completedAt)
}

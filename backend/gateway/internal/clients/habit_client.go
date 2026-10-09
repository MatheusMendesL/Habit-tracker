package clients

import (
	"context"
	"fmt"
	"gateway/internal/helper"
	"time"

	"gateway/internal/dto"

	pb "shared/pb/habit"
)

type HabitClient struct {
	client pb.HabitServiceClient
}

func NewHabitClient(client pb.HabitServiceClient) *HabitClient {
	return &HabitClient{client: client}
}

func habitFromProto(habit *pb.Habit) *dto.Habit {
	if habit == nil {
		return nil
	}
	result := &dto.Habit{
		ID:          habit.Id,
		UserID:      habit.UserId,
		Name:        habit.Name,
		Description: habit.Description,
		ImageURL:    habit.ImageUrl,
	}
	if habit.CreatedAt != nil && habit.CreatedAt.IsValid() {
		result.CreatedAt = habit.CreatedAt.AsTime()
	}
	return result
}

func (c *HabitClient) CreateHabit(ctx context.Context, userID string, request *dto.CreateHabitRequest) (*dto.HabitResponse, error) {
	if request == nil {
		return nil, fmt.Errorf("invalid request")
	}
	res, err := c.client.CreateHabit(ctx, &pb.CreateHabitRequest{Habit: &pb.Habit{
		UserId:      userID,
		Name:        request.Name,
		Description: request.Description,
		ImageUrl:    request.ImageURL,
	}})
	if err != nil {
		return nil, err
	}
	if res == nil || res.Habit == nil {
		return nil, fmt.Errorf("empty response from habit service")
	}
	return &dto.HabitResponse{Habit: habitFromProto(res.Habit)}, nil
}

func (c *HabitClient) GetHabitByID(ctx context.Context, habitID string) (*dto.HabitResponse, error) {
	res, err := c.client.GetHabitByID(ctx, &pb.GetHabitByIDRequest{HabitId: habitID})
	if err != nil {
		return nil, err
	}
	if res == nil || res.Habit == nil {
		return nil, fmt.Errorf("empty response from habit service")
	}
	return &dto.HabitResponse{Habit: habitFromProto(res.Habit)}, nil
}

func (c *HabitClient) ListHabitsByUser(ctx context.Context, userID string) (*dto.HabitsResponse, error) {
	res, err := c.client.ListHabitsByUser(ctx, &pb.ListHabitsByUserRequest{UserId: userID})
	if err != nil {
		return nil, err
	}
	if res == nil {
		return nil, fmt.Errorf("empty response from habit service")
	}
	habits := make([]*dto.Habit, 0, len(res.Habits))
	for _, habit := range res.Habits {
		if mapped := habitFromProto(habit); mapped != nil {
			habits = append(habits, mapped)
		}
	}
	return &dto.HabitsResponse{Habits: habits}, nil
}

func (c *HabitClient) ListHabitsByRoutine(ctx context.Context, routineID string) (*dto.HabitsResponse, error) {
	res, err := c.client.ListHabitsByRoutine(ctx, &pb.ListHabitsByRoutineRequest{RoutineId: routineID})
	if err != nil {
		return nil, err
	}
	if res == nil {
		return nil, fmt.Errorf("empty response from habit service")
	}
	habits := make([]*dto.Habit, 0, len(res.Habits))
	for _, habit := range res.Habits {
		if mapped := habitFromProto(habit); mapped != nil {
			habits = append(habits, mapped)
		}
	}
	return &dto.HabitsResponse{Habits: habits}, nil
}

func (c *HabitClient) EditHabit(ctx context.Context, habitID string, request *dto.EditHabitRequest) (*dto.HabitResponse, error) {
	if request == nil {
		return nil, fmt.Errorf("invalid request")
	}
	res, err := c.client.EditHabit(ctx, &pb.EditHabitRequest{
		HabitId:     habitID,
		Name:        request.Name,
		Description: request.Description,
		ImageUrl:    request.ImageURL,
	})
	if err != nil {
		return nil, err
	}
	if res == nil || res.Habit == nil {
		return nil, fmt.Errorf("empty response from habit service")
	}
	return &dto.HabitResponse{Habit: habitFromProto(res.Habit)}, nil
}

func (c *HabitClient) DeleteHabit(ctx context.Context, habitID string) (*dto.HabitActionResponse, error) {
	res, err := c.client.DeleteHabit(ctx, &pb.DeleteHabitRequest{HabitId: habitID})
	if err != nil {
		return nil, err
	}
	if res == nil {
		return nil, fmt.Errorf("empty response from habit service")
	}
	return &dto.HabitActionResponse{Success: res.Success}, nil
}

func (c *HabitClient) MarkHabitCompleted(ctx context.Context, habitID string, completedAt time.Time) (*dto.HabitActionResponse, error) {
	if completedAt.IsZero() {
		return nil, fmt.Errorf("completed_at is required")
	}
	timestamp, err := helper.TimestampToProto(completedAt)
	if err != nil {
		return nil, err
	}
	res, err := c.client.MarkHabitCompleted(ctx, &pb.MarkHabitCompletedRequest{
		HabitId:     habitID,
		CompletedAt: timestamp,
	})
	if err != nil {
		return nil, err
	}
	if res == nil {
		return nil, fmt.Errorf("empty response from habit service")
	}
	return &dto.HabitActionResponse{Success: res.Success}, nil
}

func (c *HabitClient) UnmarkHabitCompleted(ctx context.Context, habitID string, completedAt time.Time) (*dto.HabitActionResponse, error) {
	if completedAt.IsZero() {
		return nil, fmt.Errorf("completed_at is required")
	}
	timestamp, err := helper.TimestampToProto(completedAt)
	if err != nil {
		return nil, err
	}
	res, err := c.client.UnmarkHabitCompleted(ctx, &pb.UnmarkHabitCompletedRequest{
		HabitId:     habitID,
		CompletedAt: timestamp,
	})
	if err != nil {
		return nil, err
	}
	if res == nil {
		return nil, fmt.Errorf("empty response from habit service")
	}
	return &dto.HabitActionResponse{Success: res.Success}, nil
}

func (c *HabitClient) GetHabitLogs(ctx context.Context, habitID string, startDate, endDate time.Time) (*dto.HabitLogsResponse, error) {
	if startDate.IsZero() || endDate.IsZero() {
		return nil, fmt.Errorf("start_date and end_date are required")
	}
	startTimestamp, err := helper.TimestampToProto(startDate)
	if err != nil {
		return nil, err
	}
	endTimestamp, err := helper.TimestampToProto(endDate)
	if err != nil {
		return nil, err
	}
	res, err := c.client.GetHabitLogs(ctx, &pb.GetHabitLogsRequest{
		HabitId:   habitID,
		StartDate: startTimestamp,
		EndDate:   endTimestamp,
	})
	if err != nil {
		return nil, err
	}
	if res == nil {
		return nil, fmt.Errorf("empty response from habit service")
	}
	logs := make([]*dto.HabitLog, 0, len(res.Logs))
	for _, log := range res.Logs {
		if log == nil || log.CompletedAt == nil || !log.CompletedAt.IsValid() {
			continue
		}
		logs = append(logs, &dto.HabitLog{HabitID: log.HabitId, CompletedAt: log.CompletedAt.AsTime()})
	}
	return &dto.HabitLogsResponse{Logs: logs}, nil
}

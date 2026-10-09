package clients

import (
	"context"
	"fmt"
	"gateway/internal/helper"
	"time"

	"gateway/internal/dto"

	pb "shared/pb/habit"
)

type RoutineClient struct {
	client pb.RoutineServiceClient
}

func NewRoutineClient(client pb.RoutineServiceClient) *RoutineClient {
	return &RoutineClient{client: client}
}

func routineFromProto(routine *pb.Routine) *dto.Routine {
	if routine == nil {
		return nil
	}
	result := &dto.Routine{ID: routine.Id, UserID: routine.UserId, Name: routine.Name}
	if routine.CreatedAt != nil && routine.CreatedAt.IsValid() {
		result.CreatedAt = routine.CreatedAt.AsTime()
	}
	return result
}

func (c *RoutineClient) CreateRoutine(ctx context.Context, userID string, request *dto.CreateRoutineRequest) (*dto.RoutineResponse, error) {
	if request == nil {
		return nil, fmt.Errorf("invalid request")
	}
	res, err := c.client.CreateRoutine(ctx, &pb.CreateRoutineRequest{Routine: &pb.Routine{
		UserId: userID,
		Name:   request.Name,
	}})
	if err != nil {
		return nil, err
	}
	if res == nil || res.Routine == nil {
		return nil, fmt.Errorf("empty response from routine service")
	}
	return &dto.RoutineResponse{Routine: routineFromProto(res.Routine)}, nil
}

func (c *RoutineClient) GetRoutineByID(ctx context.Context, routineID string) (*dto.RoutineResponse, error) {
	res, err := c.client.GetRoutineByID(ctx, &pb.GetRoutineByIDRequest{RoutineId: routineID})
	if err != nil {
		return nil, err
	}
	if res == nil || res.Routine == nil {
		return nil, fmt.Errorf("empty response from routine service")
	}
	return &dto.RoutineResponse{Routine: routineFromProto(res.Routine)}, nil
}

func (c *RoutineClient) EditRoutine(ctx context.Context, routineID string, request *dto.EditRoutineRequest) (*dto.RoutineResponse, error) {
	if request == nil {
		return nil, fmt.Errorf("invalid request")
	}
	res, err := c.client.EditRoutine(ctx, &pb.EditRoutineRequest{RoutineId: routineID, Name: request.Name})
	if err != nil {
		return nil, err
	}
	if res == nil || res.Routine == nil {
		return nil, fmt.Errorf("empty response from routine service")
	}
	return &dto.RoutineResponse{Routine: routineFromProto(res.Routine)}, nil
}

func (c *RoutineClient) DeleteRoutine(ctx context.Context, routineID string) (*dto.RoutineActionResponse, error) {
	res, err := c.client.DeleteRoutine(ctx, &pb.DeleteRoutineRequest{RoutineId: routineID})
	if err != nil {
		return nil, err
	}
	if res == nil {
		return nil, fmt.Errorf("empty response from routine service")
	}
	return &dto.RoutineActionResponse{Success: res.Success}, nil
}

func (c *RoutineClient) AddHabitToRoutine(ctx context.Context, routineID string, request *dto.RoutineHabitRequest) (*dto.RoutineActionResponse, error) {
	if request == nil {
		return nil, fmt.Errorf("invalid request")
	}
	res, err := c.client.AddHabitToRoutine(ctx, &pb.AddHabitToRoutineRequest{RoutineId: routineID, HabitId: request.HabitID})
	if err != nil {
		return nil, err
	}
	if res == nil {
		return nil, fmt.Errorf("empty response from routine service")
	}
	return &dto.RoutineActionResponse{Success: res.Success}, nil
}

func (c *RoutineClient) RemoveHabitFromRoutine(ctx context.Context, routineID string, request *dto.RoutineHabitRequest) (*dto.RoutineActionResponse, error) {
	if request == nil {
		return nil, fmt.Errorf("invalid request")
	}
	res, err := c.client.RemoveHabitFromRoutine(ctx, &pb.RemoveHabitFromRoutineRequest{RoutineId: routineID, HabitId: request.HabitID})
	if err != nil {
		return nil, err
	}
	if res == nil {
		return nil, fmt.Errorf("empty response from routine service")
	}
	return &dto.RoutineActionResponse{Success: res.Success}, nil
}

func (c *RoutineClient) ListRoutinesByUser(ctx context.Context, userID string) (*dto.RoutinesResponse, error) {
	res, err := c.client.ListRoutinesByUser(ctx, &pb.ListRoutinesByUserRequest{UserId: userID})
	if err != nil {
		return nil, err
	}
	if res == nil {
		return nil, fmt.Errorf("empty response from routine service")
	}
	routines := make([]*dto.Routine, 0, len(res.Routines))
	for _, routine := range res.Routines {
		if mapped := routineFromProto(routine); mapped != nil {
			routines = append(routines, mapped)
		}
	}
	return &dto.RoutinesResponse{Routines: routines}, nil
}

func (c *RoutineClient) MarkRoutineCompleted(ctx context.Context, routineID string, completedAt time.Time) (*dto.RoutineActionResponse, error) {
	if completedAt.IsZero() {
		return nil, fmt.Errorf("completed_at is required")
	}
	timestamp, err := helper.TimestampToProto(completedAt)
	if err != nil {
		return nil, err
	}
	res, err := c.client.MarkRoutineCompleted(ctx, &pb.MarkRoutineCompletedRequest{RoutineId: routineID, CompletedAt: timestamp})
	if err != nil {
		return nil, err
	}
	if res == nil {
		return nil, fmt.Errorf("empty response from routine service")
	}
	return &dto.RoutineActionResponse{Success: res.Success}, nil
}

func (c *RoutineClient) UnmarkRoutineCompleted(ctx context.Context, routineID string, completedAt time.Time) (*dto.RoutineActionResponse, error) {
	if completedAt.IsZero() {
		return nil, fmt.Errorf("completed_at is required")
	}
	timestamp, err := helper.TimestampToProto(completedAt)
	if err != nil {
		return nil, err
	}
	res, err := c.client.UnmarkRoutineCompleted(ctx, &pb.UnmarkRoutineCompletedRequest{RoutineId: routineID, CompletedAt: timestamp})
	if err != nil {
		return nil, err
	}
	if res == nil {
		return nil, fmt.Errorf("empty response from routine service")
	}
	return &dto.RoutineActionResponse{Success: res.Success}, nil
}

func (c *RoutineClient) GetRoutineLogs(ctx context.Context, routineID string, startDate, endDate time.Time) (*dto.RoutineLogsResponse, error) {
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
	res, err := c.client.GetRoutineLogs(ctx, &pb.GetRoutineLogsRequest{RoutineId: routineID, StartDate: startTimestamp, EndDate: endTimestamp})
	if err != nil {
		return nil, err
	}
	if res == nil {
		return nil, fmt.Errorf("empty response from routine service")
	}
	logs := make([]*dto.RoutineLog, 0, len(res.Logs))
	for _, log := range res.Logs {
		if log == nil || log.CompletedAt == nil || !log.CompletedAt.IsValid() {
			continue
		}
		logs = append(logs, &dto.RoutineLog{RoutineID: log.RoutineId, CompletedAt: log.CompletedAt.AsTime()})
	}
	return &dto.RoutineLogsResponse{Logs: logs}, nil
}

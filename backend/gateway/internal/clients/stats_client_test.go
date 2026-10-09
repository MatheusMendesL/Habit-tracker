package clients

import (
	"context"
	"errors"
	"testing"
	"time"

	pb "shared/pb/stats"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type statsServiceClientMock struct {
	createUserStats           func(context.Context, *pb.CreateUserStatsRequest, ...grpc.CallOption) (*pb.CreateUserStatsResponse, error)
	getUserStats              func(context.Context, *pb.GetUserStatsRequest, ...grpc.CallOption) (*pb.GetUserStatsResponse, error)
	deleteUserStats           func(context.Context, *pb.DeleteUserStatsRequest, ...grpc.CallOption) (*pb.DeleteUserStatsResponse, error)
	registerHabitCompletion   func(context.Context, *pb.RegisterHabitCompletionRequest, ...grpc.CallOption) (*pb.RegisterHabitCompletionResponse, error)
	undoHabitCompletion       func(context.Context, *pb.UndoHabitCompletionRequest, ...grpc.CallOption) (*pb.UndoHabitCompletionResponse, error)
	registerRoutineCompletion func(context.Context, *pb.RegisterRoutineCompletionRequest, ...grpc.CallOption) (*pb.RegisterRoutineCompletionResponse, error)
	undoRoutineCompletion     func(context.Context, *pb.UndoRoutineCompletionRequest, ...grpc.CallOption) (*pb.UndoRoutineCompletionResponse, error)
}

func (m *statsServiceClientMock) CreateUserStats(ctx context.Context, req *pb.CreateUserStatsRequest, opts ...grpc.CallOption) (*pb.CreateUserStatsResponse, error) {
	return m.createUserStats(ctx, req, opts...)
}
func (m *statsServiceClientMock) GetUserStats(ctx context.Context, req *pb.GetUserStatsRequest, opts ...grpc.CallOption) (*pb.GetUserStatsResponse, error) {
	return m.getUserStats(ctx, req, opts...)
}
func (m *statsServiceClientMock) DeleteUserStats(ctx context.Context, req *pb.DeleteUserStatsRequest, opts ...grpc.CallOption) (*pb.DeleteUserStatsResponse, error) {
	return m.deleteUserStats(ctx, req, opts...)
}
func (m *statsServiceClientMock) RegisterHabitCompletion(ctx context.Context, req *pb.RegisterHabitCompletionRequest, opts ...grpc.CallOption) (*pb.RegisterHabitCompletionResponse, error) {
	return m.registerHabitCompletion(ctx, req, opts...)
}
func (m *statsServiceClientMock) UndoHabitCompletion(ctx context.Context, req *pb.UndoHabitCompletionRequest, opts ...grpc.CallOption) (*pb.UndoHabitCompletionResponse, error) {
	return m.undoHabitCompletion(ctx, req, opts...)
}
func (m *statsServiceClientMock) RegisterRoutineCompletion(ctx context.Context, req *pb.RegisterRoutineCompletionRequest, opts ...grpc.CallOption) (*pb.RegisterRoutineCompletionResponse, error) {
	return m.registerRoutineCompletion(ctx, req, opts...)
}
func (m *statsServiceClientMock) UndoRoutineCompletion(ctx context.Context, req *pb.UndoRoutineCompletionRequest, opts ...grpc.CallOption) (*pb.UndoRoutineCompletionResponse, error) {
	return m.undoRoutineCompletion(ctx, req, opts...)
}

var statsTestTime = time.Date(2026, time.October, 8, 12, 0, 0, 0, time.UTC)

func expectedStatsProto() *pb.UserStats {
	return &pb.UserStats{
		UserId:               "user-id",
		CompletedHabits:      12,
		CompletedRoutines:    7,
		CurrentHabitStreak:   4,
		LongestHabitStreak:   8,
		CurrentRoutineStreak: 3,
		LongestRoutineStreak: 9,
		CreatedAt:            timestamppb.New(statsTestTime),
		UpdatedAt:            timestamppb.New(statsTestTime.Add(time.Hour)),
	}
}

func TestStatsClient_CreateUserStats(t *testing.T) {
	mock := &statsServiceClientMock{
		createUserStats: func(_ context.Context, req *pb.CreateUserStatsRequest, _ ...grpc.CallOption) (*pb.CreateUserStatsResponse, error) {
			if req.UserId != "user-id" {
				t.Errorf("CreateUserStats() user id = %q, want %q", req.UserId, "user-id")
			}
			return &pb.CreateUserStatsResponse{Stats: expectedStatsProto()}, nil
		},
	}
	response, err := NewStatsClient(mock).CreateUserStats(context.Background(), "user-id")
	if err != nil {
		t.Fatalf("CreateUserStats() returned unexpected error: %v", err)
	}
	if response == nil || response.Stats.UserID != "user-id" || response.Stats.CompletedHabits != 12 ||
		response.Stats.CompletedRoutines != 7 || response.Stats.CurrentHabitStreak != 4 ||
		response.Stats.LongestHabitStreak != 8 || response.Stats.CurrentRoutineStreak != 3 ||
		response.Stats.LongestRoutineStreak != 9 || response.Stats.CreatedAt != statsTestTime ||
		response.Stats.UpdatedAt != statsTestTime.Add(time.Hour) {
		t.Fatalf("CreateUserStats() returned unexpected result: %#v", response)
	}
}

func TestStatsClient_CreateUserStats_InvalidArgument(t *testing.T) {
	expectedErr := status.Error(codes.InvalidArgument, "invalid argument")
	mock := &statsServiceClientMock{createUserStats: func(context.Context, *pb.CreateUserStatsRequest, ...grpc.CallOption) (*pb.CreateUserStatsResponse, error) {
		return nil, expectedErr
	}}
	response, err := NewStatsClient(mock).CreateUserStats(context.Background(), "invalid-id")
	if response != nil || !errors.Is(err, expectedErr) {
		t.Fatalf("CreateUserStats() = (%#v, %v), want nil and %v", response, err, expectedErr)
	}
}

func TestStatsClient_CreateUserStats_DatabaseError(t *testing.T) {
	expectedErr := status.Error(codes.Internal, "database failure")
	mock := &statsServiceClientMock{createUserStats: func(context.Context, *pb.CreateUserStatsRequest, ...grpc.CallOption) (*pb.CreateUserStatsResponse, error) {
		return nil, expectedErr
	}}
	response, err := NewStatsClient(mock).CreateUserStats(context.Background(), "user-id")
	if response != nil || !errors.Is(err, expectedErr) {
		t.Fatalf("CreateUserStats() = (%#v, %v), want nil and %v", response, err, expectedErr)
	}
}

func TestStatsClient_GetUserStats(t *testing.T) {
	mock := &statsServiceClientMock{
		getUserStats: func(_ context.Context, req *pb.GetUserStatsRequest, _ ...grpc.CallOption) (*pb.GetUserStatsResponse, error) {
			if req.UserId != "user-id" {
				t.Errorf("GetUserStats() user id = %q, want %q", req.UserId, "user-id")
			}
			return &pb.GetUserStatsResponse{Stats: expectedStatsProto()}, nil
		},
	}
	response, err := NewStatsClient(mock).GetUserStats(context.Background(), "user-id")
	if err != nil {
		t.Fatalf("GetUserStats() returned unexpected error: %v", err)
	}
	if response == nil || response.Stats.UserID != "user-id" || response.Stats.CompletedHabits != 12 || response.Stats.CompletedRoutines != 7 {
		t.Fatalf("GetUserStats() returned unexpected result: %#v", response)
	}
}

func TestStatsClient_GetUserStats_InvalidArgument(t *testing.T) {
	expectedErr := status.Error(codes.InvalidArgument, "invalid argument")
	mock := &statsServiceClientMock{getUserStats: func(context.Context, *pb.GetUserStatsRequest, ...grpc.CallOption) (*pb.GetUserStatsResponse, error) {
		return nil, expectedErr
	}}
	response, err := NewStatsClient(mock).GetUserStats(context.Background(), "invalid-id")
	if response != nil || !errors.Is(err, expectedErr) {
		t.Fatalf("GetUserStats() = (%#v, %v), want nil and %v", response, err, expectedErr)
	}
}

func TestStatsClient_GetUserStats_UserNotFound(t *testing.T) {
	expectedErr := status.Error(codes.NotFound, "user not found")
	mock := &statsServiceClientMock{getUserStats: func(context.Context, *pb.GetUserStatsRequest, ...grpc.CallOption) (*pb.GetUserStatsResponse, error) {
		return nil, expectedErr
	}}
	response, err := NewStatsClient(mock).GetUserStats(context.Background(), "missing-user")
	if response != nil || !errors.Is(err, expectedErr) {
		t.Fatalf("GetUserStats() = (%#v, %v), want nil and %v", response, err, expectedErr)
	}
}

func TestStatsClient_GetUserStats_StatsNotFound(t *testing.T) {
	expectedErr := status.Error(codes.NotFound, "stats not found")
	mock := &statsServiceClientMock{getUserStats: func(context.Context, *pb.GetUserStatsRequest, ...grpc.CallOption) (*pb.GetUserStatsResponse, error) {
		return nil, expectedErr
	}}
	response, err := NewStatsClient(mock).GetUserStats(context.Background(), "user-without-stats")
	if response != nil || !errors.Is(err, expectedErr) {
		t.Fatalf("GetUserStats() = (%#v, %v), want nil and %v", response, err, expectedErr)
	}
}

func TestStatsClient_GetUserStats_DatabaseError(t *testing.T) {
	expectedErr := status.Error(codes.Internal, "database failure")
	mock := &statsServiceClientMock{getUserStats: func(context.Context, *pb.GetUserStatsRequest, ...grpc.CallOption) (*pb.GetUserStatsResponse, error) {
		return nil, expectedErr
	}}
	response, err := NewStatsClient(mock).GetUserStats(context.Background(), "user-id")
	if response != nil || !errors.Is(err, expectedErr) {
		t.Fatalf("GetUserStats() = (%#v, %v), want nil and %v", response, err, expectedErr)
	}
}

func TestStatsClient_DeleteUserStats(t *testing.T) {
	mock := &statsServiceClientMock{deleteUserStats: func(_ context.Context, req *pb.DeleteUserStatsRequest, _ ...grpc.CallOption) (*pb.DeleteUserStatsResponse, error) {
		if req.UserId != "user-id" {
			t.Errorf("DeleteUserStats() user id = %q, want %q", req.UserId, "user-id")
		}
		return &pb.DeleteUserStatsResponse{Success: true}, nil
	}}
	response, err := NewStatsClient(mock).DeleteUserStats(context.Background(), "user-id")
	if err != nil {
		t.Fatalf("DeleteUserStats() returned unexpected error: %v", err)
	}
	if response == nil || !response.Success {
		t.Fatalf("DeleteUserStats() returned %#v, want success", response)
	}
}

func TestStatsClient_DeleteUserStats_InvalidArgument(t *testing.T) {
	expectedErr := status.Error(codes.InvalidArgument, "invalid argument")
	mock := &statsServiceClientMock{deleteUserStats: func(context.Context, *pb.DeleteUserStatsRequest, ...grpc.CallOption) (*pb.DeleteUserStatsResponse, error) {
		return nil, expectedErr
	}}
	response, err := NewStatsClient(mock).DeleteUserStats(context.Background(), "invalid-id")
	if response != nil || !errors.Is(err, expectedErr) {
		t.Fatalf("DeleteUserStats() = (%#v, %v), want nil and %v", response, err, expectedErr)
	}
}

func TestStatsClient_DeleteUserStats_UserNotFound(t *testing.T) {
	expectedErr := status.Error(codes.NotFound, "user not found")
	mock := &statsServiceClientMock{deleteUserStats: func(context.Context, *pb.DeleteUserStatsRequest, ...grpc.CallOption) (*pb.DeleteUserStatsResponse, error) {
		return nil, expectedErr
	}}
	response, err := NewStatsClient(mock).DeleteUserStats(context.Background(), "missing-user")
	if response != nil || !errors.Is(err, expectedErr) {
		t.Fatalf("DeleteUserStats() = (%#v, %v), want nil and %v", response, err, expectedErr)
	}
}

func TestStatsClient_DeleteUserStats_StatsNotFound(t *testing.T) {
	expectedErr := status.Error(codes.NotFound, "stats not found")
	mock := &statsServiceClientMock{deleteUserStats: func(context.Context, *pb.DeleteUserStatsRequest, ...grpc.CallOption) (*pb.DeleteUserStatsResponse, error) {
		return nil, expectedErr
	}}
	response, err := NewStatsClient(mock).DeleteUserStats(context.Background(), "user-without-stats")
	if response != nil || !errors.Is(err, expectedErr) {
		t.Fatalf("DeleteUserStats() = (%#v, %v), want nil and %v", response, err, expectedErr)
	}
}

func TestStatsClient_DeleteUserStats_DatabaseError(t *testing.T) {
	expectedErr := status.Error(codes.Internal, "database failure")
	mock := &statsServiceClientMock{deleteUserStats: func(context.Context, *pb.DeleteUserStatsRequest, ...grpc.CallOption) (*pb.DeleteUserStatsResponse, error) {
		return nil, expectedErr
	}}
	response, err := NewStatsClient(mock).DeleteUserStats(context.Background(), "user-id")
	if response != nil || !errors.Is(err, expectedErr) {
		t.Fatalf("DeleteUserStats() = (%#v, %v), want nil and %v", response, err, expectedErr)
	}
}

func TestStatsClient_RegisterHabitCompletion(t *testing.T) {
	expectedRequest := &pb.RegisterHabitCompletionRequest{UserId: "user-id", HabitId: "habit-id", CompletedAt: timestamppb.New(statsTestTime)}
	mock := &statsServiceClientMock{registerHabitCompletion: func(_ context.Context, req *pb.RegisterHabitCompletionRequest, _ ...grpc.CallOption) (*pb.RegisterHabitCompletionResponse, error) {
		if !proto.Equal(req, expectedRequest) {
			t.Errorf("RegisterHabitCompletion() request = %v, want %v", req, expectedRequest)
		}
		return &pb.RegisterHabitCompletionResponse{Success: true}, nil
	}}
	response, err := NewStatsClient(mock).RegisterHabitCompletion(context.Background(), "user-id", "habit-id", statsTestTime)
	if err != nil {
		t.Fatalf("RegisterHabitCompletion() returned unexpected error: %v", err)
	}
	if response == nil || !response.Success {
		t.Fatalf("RegisterHabitCompletion() returned %#v, want success", response)
	}
}

func TestStatsClient_RegisterHabitCompletion_InvalidArgument(t *testing.T) {
	expectedErr := status.Error(codes.InvalidArgument, "invalid argument")
	mock := &statsServiceClientMock{registerHabitCompletion: func(context.Context, *pb.RegisterHabitCompletionRequest, ...grpc.CallOption) (*pb.RegisterHabitCompletionResponse, error) {
		return nil, expectedErr
	}}
	response, err := NewStatsClient(mock).RegisterHabitCompletion(context.Background(), "invalid-user", "habit-id", statsTestTime)
	if response != nil || !errors.Is(err, expectedErr) {
		t.Fatalf("RegisterHabitCompletion() = (%#v, %v), want nil and %v", response, err, expectedErr)
	}
}

func TestStatsClient_RegisterHabitCompletion_UserNotFound(t *testing.T) {
	expectedErr := status.Error(codes.NotFound, "user not found")
	mock := &statsServiceClientMock{registerHabitCompletion: func(context.Context, *pb.RegisterHabitCompletionRequest, ...grpc.CallOption) (*pb.RegisterHabitCompletionResponse, error) {
		return nil, expectedErr
	}}
	response, err := NewStatsClient(mock).RegisterHabitCompletion(context.Background(), "missing-user", "habit-id", statsTestTime)
	if response != nil || !errors.Is(err, expectedErr) {
		t.Fatalf("RegisterHabitCompletion() = (%#v, %v), want nil and %v", response, err, expectedErr)
	}
}

func TestStatsClient_RegisterHabitCompletion_DatabaseError(t *testing.T) {
	expectedErr := status.Error(codes.Internal, "database failure")
	mock := &statsServiceClientMock{registerHabitCompletion: func(context.Context, *pb.RegisterHabitCompletionRequest, ...grpc.CallOption) (*pb.RegisterHabitCompletionResponse, error) {
		return nil, expectedErr
	}}
	response, err := NewStatsClient(mock).RegisterHabitCompletion(context.Background(), "user-id", "habit-id", statsTestTime)
	if response != nil || !errors.Is(err, expectedErr) {
		t.Fatalf("RegisterHabitCompletion() = (%#v, %v), want nil and %v", response, err, expectedErr)
	}
}

func TestStatsClient_UndoHabitCompletion(t *testing.T) {
	expectedRequest := &pb.UndoHabitCompletionRequest{UserId: "user-id", HabitId: "habit-id", CompletedAt: timestamppb.New(statsTestTime)}
	mock := &statsServiceClientMock{undoHabitCompletion: func(_ context.Context, req *pb.UndoHabitCompletionRequest, _ ...grpc.CallOption) (*pb.UndoHabitCompletionResponse, error) {
		if !proto.Equal(req, expectedRequest) {
			t.Errorf("UndoHabitCompletion() request = %v, want %v", req, expectedRequest)
		}
		return &pb.UndoHabitCompletionResponse{Success: true}, nil
	}}
	response, err := NewStatsClient(mock).UndoHabitCompletion(context.Background(), "user-id", "habit-id", statsTestTime)
	if err != nil {
		t.Fatalf("UndoHabitCompletion() returned unexpected error: %v", err)
	}
	if response == nil || !response.Success {
		t.Fatalf("UndoHabitCompletion() returned %#v, want success", response)
	}
}

func TestStatsClient_UndoHabitCompletion_InvalidArgument(t *testing.T) {
	expectedErr := status.Error(codes.InvalidArgument, "invalid argument")
	mock := &statsServiceClientMock{undoHabitCompletion: func(context.Context, *pb.UndoHabitCompletionRequest, ...grpc.CallOption) (*pb.UndoHabitCompletionResponse, error) {
		return nil, expectedErr
	}}
	response, err := NewStatsClient(mock).UndoHabitCompletion(context.Background(), "invalid-user", "habit-id", statsTestTime)
	if response != nil || !errors.Is(err, expectedErr) {
		t.Fatalf("UndoHabitCompletion() = (%#v, %v), want nil and %v", response, err, expectedErr)
	}
}

func TestStatsClient_UndoHabitCompletion_UserNotFound(t *testing.T) {
	expectedErr := status.Error(codes.NotFound, "user not found")
	mock := &statsServiceClientMock{undoHabitCompletion: func(context.Context, *pb.UndoHabitCompletionRequest, ...grpc.CallOption) (*pb.UndoHabitCompletionResponse, error) {
		return nil, expectedErr
	}}
	response, err := NewStatsClient(mock).UndoHabitCompletion(context.Background(), "missing-user", "habit-id", statsTestTime)
	if response != nil || !errors.Is(err, expectedErr) {
		t.Fatalf("UndoHabitCompletion() = (%#v, %v), want nil and %v", response, err, expectedErr)
	}
}

func TestStatsClient_UndoHabitCompletion_DatabaseError(t *testing.T) {
	expectedErr := status.Error(codes.Internal, "database failure")
	mock := &statsServiceClientMock{undoHabitCompletion: func(context.Context, *pb.UndoHabitCompletionRequest, ...grpc.CallOption) (*pb.UndoHabitCompletionResponse, error) {
		return nil, expectedErr
	}}
	response, err := NewStatsClient(mock).UndoHabitCompletion(context.Background(), "user-id", "habit-id", statsTestTime)
	if response != nil || !errors.Is(err, expectedErr) {
		t.Fatalf("UndoHabitCompletion() = (%#v, %v), want nil and %v", response, err, expectedErr)
	}
}

func TestStatsClient_RegisterRoutineCompletion(t *testing.T) {
	expectedRequest := &pb.RegisterRoutineCompletionRequest{UserId: "user-id", RoutineId: "routine-id", CompletedAt: timestamppb.New(statsTestTime)}
	mock := &statsServiceClientMock{registerRoutineCompletion: func(_ context.Context, req *pb.RegisterRoutineCompletionRequest, _ ...grpc.CallOption) (*pb.RegisterRoutineCompletionResponse, error) {
		if !proto.Equal(req, expectedRequest) {
			t.Errorf("RegisterRoutineCompletion() request = %v, want %v", req, expectedRequest)
		}
		return &pb.RegisterRoutineCompletionResponse{Success: true}, nil
	}}
	response, err := NewStatsClient(mock).RegisterRoutineCompletion(context.Background(), "user-id", "routine-id", statsTestTime)
	if err != nil {
		t.Fatalf("RegisterRoutineCompletion() returned unexpected error: %v", err)
	}
	if response == nil || !response.Success {
		t.Fatalf("RegisterRoutineCompletion() returned %#v, want success", response)
	}
}

func TestStatsClient_RegisterRoutineCompletion_InvalidArgument(t *testing.T) {
	expectedErr := status.Error(codes.InvalidArgument, "invalid argument")
	mock := &statsServiceClientMock{registerRoutineCompletion: func(context.Context, *pb.RegisterRoutineCompletionRequest, ...grpc.CallOption) (*pb.RegisterRoutineCompletionResponse, error) {
		return nil, expectedErr
	}}
	response, err := NewStatsClient(mock).RegisterRoutineCompletion(context.Background(), "invalid-user", "routine-id", statsTestTime)
	if response != nil || !errors.Is(err, expectedErr) {
		t.Fatalf("RegisterRoutineCompletion() = (%#v, %v), want nil and %v", response, err, expectedErr)
	}
}

func TestStatsClient_RegisterRoutineCompletion_UserNotFound(t *testing.T) {
	expectedErr := status.Error(codes.NotFound, "user not found")
	mock := &statsServiceClientMock{registerRoutineCompletion: func(context.Context, *pb.RegisterRoutineCompletionRequest, ...grpc.CallOption) (*pb.RegisterRoutineCompletionResponse, error) {
		return nil, expectedErr
	}}
	response, err := NewStatsClient(mock).RegisterRoutineCompletion(context.Background(), "missing-user", "routine-id", statsTestTime)
	if response != nil || !errors.Is(err, expectedErr) {
		t.Fatalf("RegisterRoutineCompletion() = (%#v, %v), want nil and %v", response, err, expectedErr)
	}
}

func TestStatsClient_RegisterRoutineCompletion_DatabaseError(t *testing.T) {
	expectedErr := status.Error(codes.Internal, "database failure")
	mock := &statsServiceClientMock{registerRoutineCompletion: func(context.Context, *pb.RegisterRoutineCompletionRequest, ...grpc.CallOption) (*pb.RegisterRoutineCompletionResponse, error) {
		return nil, expectedErr
	}}
	response, err := NewStatsClient(mock).RegisterRoutineCompletion(context.Background(), "user-id", "routine-id", statsTestTime)
	if response != nil || !errors.Is(err, expectedErr) {
		t.Fatalf("RegisterRoutineCompletion() = (%#v, %v), want nil and %v", response, err, expectedErr)
	}
}

func TestStatsClient_UndoRoutineCompletion(t *testing.T) {
	expectedRequest := &pb.UndoRoutineCompletionRequest{UserId: "user-id", RoutineId: "routine-id", CompletedAt: timestamppb.New(statsTestTime)}
	mock := &statsServiceClientMock{undoRoutineCompletion: func(_ context.Context, req *pb.UndoRoutineCompletionRequest, _ ...grpc.CallOption) (*pb.UndoRoutineCompletionResponse, error) {
		if !proto.Equal(req, expectedRequest) {
			t.Errorf("UndoRoutineCompletion() request = %v, want %v", req, expectedRequest)
		}
		return &pb.UndoRoutineCompletionResponse{Success: true}, nil
	}}
	response, err := NewStatsClient(mock).UndoRoutineCompletion(context.Background(), "user-id", "routine-id", statsTestTime)
	if err != nil {
		t.Fatalf("UndoRoutineCompletion() returned unexpected error: %v", err)
	}
	if response == nil || !response.Success {
		t.Fatalf("UndoRoutineCompletion() returned %#v, want success", response)
	}
}

func TestStatsClient_UndoRoutineCompletion_InvalidArgument(t *testing.T) {
	expectedErr := status.Error(codes.InvalidArgument, "invalid argument")
	mock := &statsServiceClientMock{undoRoutineCompletion: func(context.Context, *pb.UndoRoutineCompletionRequest, ...grpc.CallOption) (*pb.UndoRoutineCompletionResponse, error) {
		return nil, expectedErr
	}}
	response, err := NewStatsClient(mock).UndoRoutineCompletion(context.Background(), "invalid-user", "routine-id", statsTestTime)
	if response != nil || !errors.Is(err, expectedErr) {
		t.Fatalf("UndoRoutineCompletion() = (%#v, %v), want nil and %v", response, err, expectedErr)
	}
}

func TestStatsClient_UndoRoutineCompletion_UserNotFound(t *testing.T) {
	expectedErr := status.Error(codes.NotFound, "user not found")
	mock := &statsServiceClientMock{undoRoutineCompletion: func(context.Context, *pb.UndoRoutineCompletionRequest, ...grpc.CallOption) (*pb.UndoRoutineCompletionResponse, error) {
		return nil, expectedErr
	}}
	response, err := NewStatsClient(mock).UndoRoutineCompletion(context.Background(), "missing-user", "routine-id", statsTestTime)
	if response != nil || !errors.Is(err, expectedErr) {
		t.Fatalf("UndoRoutineCompletion() = (%#v, %v), want nil and %v", response, err, expectedErr)
	}
}

func TestStatsClient_UndoRoutineCompletion_DatabaseError(t *testing.T) {
	expectedErr := status.Error(codes.Internal, "database failure")
	mock := &statsServiceClientMock{undoRoutineCompletion: func(context.Context, *pb.UndoRoutineCompletionRequest, ...grpc.CallOption) (*pb.UndoRoutineCompletionResponse, error) {
		return nil, expectedErr
	}}
	response, err := NewStatsClient(mock).UndoRoutineCompletion(context.Background(), "user-id", "routine-id", statsTestTime)
	if response != nil || !errors.Is(err, expectedErr) {
		t.Fatalf("UndoRoutineCompletion() = (%#v, %v), want nil and %v", response, err, expectedErr)
	}
}

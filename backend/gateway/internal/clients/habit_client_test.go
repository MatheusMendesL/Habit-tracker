package clients

import (
	"context"
	"errors"
	"testing"
	"time"

	"gateway/internal/dto"

	pb "shared/pb/habit"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type habitServiceClientMock struct {
	createHabit          func(context.Context, *pb.CreateHabitRequest, ...grpc.CallOption) (*pb.CreateHabitResponse, error)
	getHabitByID         func(context.Context, *pb.GetHabitByIDRequest, ...grpc.CallOption) (*pb.GetHabitByIDResponse, error)
	listHabitsByUser     func(context.Context, *pb.ListHabitsByUserRequest, ...grpc.CallOption) (*pb.ListHabitsByUserResponse, error)
	listHabitsByRoutine  func(context.Context, *pb.ListHabitsByRoutineRequest, ...grpc.CallOption) (*pb.ListHabitsByRoutineResponse, error)
	editHabit            func(context.Context, *pb.EditHabitRequest, ...grpc.CallOption) (*pb.EditHabitResponse, error)
	deleteHabit          func(context.Context, *pb.DeleteHabitRequest, ...grpc.CallOption) (*pb.DeleteHabitResponse, error)
	markHabitCompleted   func(context.Context, *pb.MarkHabitCompletedRequest, ...grpc.CallOption) (*pb.MarkHabitCompletedResponse, error)
	unmarkHabitCompleted func(context.Context, *pb.UnmarkHabitCompletedRequest, ...grpc.CallOption) (*pb.UnmarkHabitCompletedResponse, error)
	getHabitLogs         func(context.Context, *pb.GetHabitLogsRequest, ...grpc.CallOption) (*pb.GetHabitLogsResponse, error)
}

func (m *habitServiceClientMock) CreateHabit(ctx context.Context, req *pb.CreateHabitRequest, opts ...grpc.CallOption) (*pb.CreateHabitResponse, error) {
	return m.createHabit(ctx, req, opts...)
}
func (m *habitServiceClientMock) GetHabitByID(ctx context.Context, req *pb.GetHabitByIDRequest, opts ...grpc.CallOption) (*pb.GetHabitByIDResponse, error) {
	return m.getHabitByID(ctx, req, opts...)
}
func (m *habitServiceClientMock) ListHabitsByUser(ctx context.Context, req *pb.ListHabitsByUserRequest, opts ...grpc.CallOption) (*pb.ListHabitsByUserResponse, error) {
	return m.listHabitsByUser(ctx, req, opts...)
}
func (m *habitServiceClientMock) ListHabitsByRoutine(ctx context.Context, req *pb.ListHabitsByRoutineRequest, opts ...grpc.CallOption) (*pb.ListHabitsByRoutineResponse, error) {
	return m.listHabitsByRoutine(ctx, req, opts...)
}
func (m *habitServiceClientMock) EditHabit(ctx context.Context, req *pb.EditHabitRequest, opts ...grpc.CallOption) (*pb.EditHabitResponse, error) {
	return m.editHabit(ctx, req, opts...)
}
func (m *habitServiceClientMock) DeleteHabit(ctx context.Context, req *pb.DeleteHabitRequest, opts ...grpc.CallOption) (*pb.DeleteHabitResponse, error) {
	return m.deleteHabit(ctx, req, opts...)
}
func (m *habitServiceClientMock) MarkHabitCompleted(ctx context.Context, req *pb.MarkHabitCompletedRequest, opts ...grpc.CallOption) (*pb.MarkHabitCompletedResponse, error) {
	return m.markHabitCompleted(ctx, req, opts...)
}
func (m *habitServiceClientMock) UnmarkHabitCompleted(ctx context.Context, req *pb.UnmarkHabitCompletedRequest, opts ...grpc.CallOption) (*pb.UnmarkHabitCompletedResponse, error) {
	return m.unmarkHabitCompleted(ctx, req, opts...)
}
func (m *habitServiceClientMock) GetHabitLogs(ctx context.Context, req *pb.GetHabitLogsRequest, opts ...grpc.CallOption) (*pb.GetHabitLogsResponse, error) {
	return m.getHabitLogs(ctx, req, opts...)
}

var testTime = time.Date(2026, time.October, 8, 12, 0, 0, 0, time.UTC)

func TestHabitClient_CreateHabit(t *testing.T) {
	client := NewHabitClient(&habitServiceClientMock{
		createHabit: func(_ context.Context, req *pb.CreateHabitRequest, _ ...grpc.CallOption) (*pb.CreateHabitResponse, error) {
			expected := &pb.CreateHabitRequest{Habit: &pb.Habit{UserId: "user-id", Name: "read", Description: "daily", ImageUrl: "image"}}
			if !proto.Equal(req, expected) {
				t.Errorf("CreateHabit() request = %v, want %v", req, expected)
			}
			return &pb.CreateHabitResponse{Habit: &pb.Habit{Id: "habit-id", UserId: "user-id", Name: "read", Description: "daily", ImageUrl: "image", CreatedAt: timestamppb.New(testTime)}}, nil
		},
	})
	response, err := client.CreateHabit(context.Background(), "user-id", &dto.CreateHabitRequest{Name: "read", Description: "daily", ImageURL: "image"})
	if err != nil {
		t.Fatalf("CreateHabit() returned unexpected error: %v", err)
	}
	if response.Habit.ID != "habit-id" || response.Habit.CreatedAt != testTime {
		t.Fatalf("CreateHabit() returned unexpected result: %#v", response)
	}
}

func TestHabitClient_GetHabitByID(t *testing.T) {
	client := NewHabitClient(&habitServiceClientMock{
		getHabitByID: func(_ context.Context, req *pb.GetHabitByIDRequest, _ ...grpc.CallOption) (*pb.GetHabitByIDResponse, error) {
			expected := &pb.GetHabitByIDRequest{HabitId: "habit-id"}
			if !proto.Equal(req, expected) {
				t.Errorf("GetHabitByID() request = %v, want %v", req, expected)
			}
			return &pb.GetHabitByIDResponse{Habit: &pb.Habit{Id: "habit-id", UserId: "user-id", Name: "read", Description: "daily", ImageUrl: "image", CreatedAt: timestamppb.New(testTime)}}, nil
		},
	})
	response, err := client.GetHabitByID(context.Background(), "habit-id")
	if err != nil {
		t.Fatalf("GetHabitByID() returned unexpected error: %v", err)
	}
	if response.Habit.ID != "habit-id" || response.Habit.Name != "read" {
		t.Fatalf("GetHabitByID() returned unexpected result: %#v", response)
	}
}

func TestHabitClient_ListHabitsByUser(t *testing.T) {
	client := NewHabitClient(&habitServiceClientMock{
		listHabitsByUser: func(_ context.Context, req *pb.ListHabitsByUserRequest, _ ...grpc.CallOption) (*pb.ListHabitsByUserResponse, error) {
			expected := &pb.ListHabitsByUserRequest{UserId: "user-id"}
			if !proto.Equal(req, expected) {
				t.Errorf("ListHabitsByUser() request = %v, want %v", req, expected)
			}
			return &pb.ListHabitsByUserResponse{Habits: []*pb.Habit{{Id: "habit-id", UserId: "user-id", Name: "read", CreatedAt: timestamppb.New(testTime)}}}, nil
		},
	})
	response, err := client.ListHabitsByUser(context.Background(), "user-id")
	if err != nil {
		t.Fatalf("ListHabitsByUser() returned unexpected error: %v", err)
	}
	if len(response.Habits) != 1 || response.Habits[0].ID != "habit-id" {
		t.Fatalf("ListHabitsByUser() returned unexpected result: %#v", response)
	}
}

func TestHabitClient_ListHabitsByRoutine(t *testing.T) {
	client := NewHabitClient(&habitServiceClientMock{
		listHabitsByRoutine: func(_ context.Context, req *pb.ListHabitsByRoutineRequest, _ ...grpc.CallOption) (*pb.ListHabitsByRoutineResponse, error) {
			expected := &pb.ListHabitsByRoutineRequest{RoutineId: "routine-id"}
			if !proto.Equal(req, expected) {
				t.Errorf("ListHabitsByRoutine() request = %v, want %v", req, expected)
			}
			return &pb.ListHabitsByRoutineResponse{Habits: []*pb.Habit{{Id: "habit-id", UserId: "user-id", Name: "read", CreatedAt: timestamppb.New(testTime)}}}, nil
		},
	})
	response, err := client.ListHabitsByRoutine(context.Background(), "routine-id")
	if err != nil {
		t.Fatalf("ListHabitsByRoutine() returned unexpected error: %v", err)
	}
	if len(response.Habits) != 1 || response.Habits[0].ID != "habit-id" {
		t.Fatalf("ListHabitsByRoutine() returned unexpected result: %#v", response)
	}
}

func TestHabitClient_EditHabit(t *testing.T) {
	name, description, imageURL := "read", "daily", "image"
	client := NewHabitClient(&habitServiceClientMock{
		editHabit: func(_ context.Context, req *pb.EditHabitRequest, _ ...grpc.CallOption) (*pb.EditHabitResponse, error) {
			expected := &pb.EditHabitRequest{HabitId: "habit-id", Name: &name, Description: &description, ImageUrl: &imageURL}
			if !proto.Equal(req, expected) {
				t.Errorf("EditHabit() request = %v, want %v", req, expected)
			}
			return &pb.EditHabitResponse{Habit: &pb.Habit{Id: "habit-id", UserId: "user-id", Name: "read", Description: "daily", ImageUrl: "image", CreatedAt: timestamppb.New(testTime)}}, nil
		},
	})
	response, err := client.EditHabit(context.Background(), "habit-id", &dto.EditHabitRequest{Name: &name, Description: &description, ImageURL: &imageURL})
	if err != nil {
		t.Fatalf("EditHabit() returned unexpected error: %v", err)
	}
	if response.Habit.ID != "habit-id" || response.Habit.Description != "daily" {
		t.Fatalf("EditHabit() returned unexpected result: %#v", response)
	}
}

func TestHabitClient_DeleteHabit(t *testing.T) {
	client := NewHabitClient(&habitServiceClientMock{
		deleteHabit: func(_ context.Context, req *pb.DeleteHabitRequest, _ ...grpc.CallOption) (*pb.DeleteHabitResponse, error) {
			expected := &pb.DeleteHabitRequest{HabitId: "habit-id"}
			if !proto.Equal(req, expected) {
				t.Errorf("DeleteHabit() request = %v, want %v", req, expected)
			}
			return &pb.DeleteHabitResponse{Success: true}, nil
		},
	})
	response, err := client.DeleteHabit(context.Background(), "habit-id")
	if err != nil {
		t.Fatalf("DeleteHabit() returned unexpected error: %v", err)
	}
	if !response.Success {
		t.Fatalf("DeleteHabit() returned unexpected result: %#v", response)
	}
}

func TestHabitClient_MarkHabitCompleted(t *testing.T) {
	client := NewHabitClient(&habitServiceClientMock{
		markHabitCompleted: func(_ context.Context, req *pb.MarkHabitCompletedRequest, _ ...grpc.CallOption) (*pb.MarkHabitCompletedResponse, error) {
			expected := &pb.MarkHabitCompletedRequest{HabitId: "habit-id", CompletedAt: timestamppb.New(testTime)}
			if !proto.Equal(req, expected) {
				t.Errorf("MarkHabitCompleted() request = %v, want %v", req, expected)
			}
			return &pb.MarkHabitCompletedResponse{Success: true}, nil
		},
	})
	response, err := client.MarkHabitCompleted(context.Background(), "habit-id", testTime)
	if err != nil {
		t.Fatalf("MarkHabitCompleted() returned unexpected error: %v", err)
	}
	if !response.Success {
		t.Fatalf("MarkHabitCompleted() returned unexpected result: %#v", response)
	}
}

func TestHabitClient_UnmarkHabitCompleted(t *testing.T) {
	client := NewHabitClient(&habitServiceClientMock{
		unmarkHabitCompleted: func(_ context.Context, req *pb.UnmarkHabitCompletedRequest, _ ...grpc.CallOption) (*pb.UnmarkHabitCompletedResponse, error) {
			expected := &pb.UnmarkHabitCompletedRequest{HabitId: "habit-id", CompletedAt: timestamppb.New(testTime)}
			if !proto.Equal(req, expected) {
				t.Errorf("UnmarkHabitCompleted() request = %v, want %v", req, expected)
			}
			return &pb.UnmarkHabitCompletedResponse{Success: true}, nil
		},
	})
	response, err := client.UnmarkHabitCompleted(context.Background(), "habit-id", testTime)
	if err != nil {
		t.Fatalf("UnmarkHabitCompleted() returned unexpected error: %v", err)
	}
	if !response.Success {
		t.Fatalf("UnmarkHabitCompleted() returned unexpected result: %#v", response)
	}
}

func TestHabitClient_GetHabitLogs(t *testing.T) {
	client := NewHabitClient(&habitServiceClientMock{
		getHabitLogs: func(_ context.Context, req *pb.GetHabitLogsRequest, _ ...grpc.CallOption) (*pb.GetHabitLogsResponse, error) {
			expected := &pb.GetHabitLogsRequest{HabitId: "habit-id", StartDate: timestamppb.New(testTime.Add(-time.Hour)), EndDate: timestamppb.New(testTime)}
			if !proto.Equal(req, expected) {
				t.Errorf("GetHabitLogs() request = %v, want %v", req, expected)
			}
			return &pb.GetHabitLogsResponse{Logs: []*pb.HabitLog{{HabitId: "habit-id", CompletedAt: timestamppb.New(testTime)}}}, nil
		},
	})
	response, err := client.GetHabitLogs(context.Background(), "habit-id", testTime.Add(-time.Hour), testTime)
	if err != nil {
		t.Fatalf("GetHabitLogs() returned unexpected error: %v", err)
	}
	if len(response.Logs) != 1 || response.Logs[0].HabitID != "habit-id" || response.Logs[0].CompletedAt != testTime {
		t.Fatalf("GetHabitLogs() returned unexpected result: %#v", response)
	}
}

func TestHabitClient_GetHabitByID_InvalidArgument(t *testing.T) {
	expectedErr := status.Error(codes.InvalidArgument, "invalid argument")
	client := NewHabitClient(&habitServiceClientMock{getHabitByID: func(context.Context, *pb.GetHabitByIDRequest, ...grpc.CallOption) (*pb.GetHabitByIDResponse, error) {
		return nil, expectedErr
	}})
	response, err := client.GetHabitByID(context.Background(), "invalid-id")
	if response != nil || !errors.Is(err, expectedErr) {
		t.Fatalf("GetHabitByID() = (%#v, %v), want nil and InvalidArgument", response, err)
	}
}

func TestHabitClient_GetHabitByID_NotFound(t *testing.T) {
	expectedErr := status.Error(codes.NotFound, "habit not found")
	client := NewHabitClient(&habitServiceClientMock{getHabitByID: func(context.Context, *pb.GetHabitByIDRequest, ...grpc.CallOption) (*pb.GetHabitByIDResponse, error) {
		return nil, expectedErr
	}})
	response, err := client.GetHabitByID(context.Background(), "missing-id")
	if response != nil || !errors.Is(err, expectedErr) {
		t.Fatalf("GetHabitByID() = (%#v, %v), want nil and NotFound", response, err)
	}
}

func TestHabitClient_GetHabitByID_DatabaseError(t *testing.T) {
	expectedErr := status.Error(codes.Internal, "database failure")
	client := NewHabitClient(&habitServiceClientMock{getHabitByID: func(context.Context, *pb.GetHabitByIDRequest, ...grpc.CallOption) (*pb.GetHabitByIDResponse, error) {
		return nil, expectedErr
	}})
	response, err := client.GetHabitByID(context.Background(), "habit-id")
	if response != nil || !errors.Is(err, expectedErr) {
		t.Fatalf("GetHabitByID() = (%#v, %v), want nil and Internal", response, err)
	}
}

func TestHabitClient_CreateHabit_InvalidArgument(t *testing.T) {
	expectedErr := status.Error(codes.InvalidArgument, "invalid argument")
	client := NewHabitClient(&habitServiceClientMock{createHabit: func(context.Context, *pb.CreateHabitRequest, ...grpc.CallOption) (*pb.CreateHabitResponse, error) {
		return nil, expectedErr
	}})
	response, err := client.CreateHabit(context.Background(), "invalid-user-id", &dto.CreateHabitRequest{Name: "read"})
	if response != nil || !errors.Is(err, expectedErr) {
		t.Fatalf("CreateHabit() = (%#v, %v), want nil and InvalidArgument", response, err)
	}
}

func TestHabitClient_CreateHabit_UserNotFound(t *testing.T) {
	expectedErr := status.Error(codes.NotFound, "user not found")
	client := NewHabitClient(&habitServiceClientMock{createHabit: func(context.Context, *pb.CreateHabitRequest, ...grpc.CallOption) (*pb.CreateHabitResponse, error) {
		return nil, expectedErr
	}})
	response, err := client.CreateHabit(context.Background(), "missing-user-id", &dto.CreateHabitRequest{Name: "read"})
	if response != nil || !errors.Is(err, expectedErr) {
		t.Fatalf("CreateHabit() = (%#v, %v), want nil and NotFound", response, err)
	}
}

func TestHabitClient_CreateHabit_DatabaseError(t *testing.T) {
	expectedErr := status.Error(codes.Internal, "database failure")
	client := NewHabitClient(&habitServiceClientMock{createHabit: func(context.Context, *pb.CreateHabitRequest, ...grpc.CallOption) (*pb.CreateHabitResponse, error) {
		return nil, expectedErr
	}})
	response, err := client.CreateHabit(context.Background(), "user-id", &dto.CreateHabitRequest{Name: "read"})
	if response != nil || !errors.Is(err, expectedErr) {
		t.Fatalf("CreateHabit() = (%#v, %v), want nil and Internal", response, err)
	}
}

func TestHabitClient_EditHabit_InvalidArgument(t *testing.T) {
	expectedErr := status.Error(codes.InvalidArgument, "invalid argument")
	name := "read"
	client := NewHabitClient(&habitServiceClientMock{editHabit: func(context.Context, *pb.EditHabitRequest, ...grpc.CallOption) (*pb.EditHabitResponse, error) {
		return nil, expectedErr
	}})
	response, err := client.EditHabit(context.Background(), "invalid-id", &dto.EditHabitRequest{Name: &name})
	if response != nil || !errors.Is(err, expectedErr) {
		t.Fatalf("EditHabit() = (%#v, %v), want nil and InvalidArgument", response, err)
	}
}

func TestHabitClient_EditHabit_HabitNotFound(t *testing.T) {
	expectedErr := status.Error(codes.NotFound, "habit not found")
	name := "read"
	client := NewHabitClient(&habitServiceClientMock{editHabit: func(context.Context, *pb.EditHabitRequest, ...grpc.CallOption) (*pb.EditHabitResponse, error) {
		return nil, expectedErr
	}})
	response, err := client.EditHabit(context.Background(), "missing-id", &dto.EditHabitRequest{Name: &name})
	if response != nil || !errors.Is(err, expectedErr) {
		t.Fatalf("EditHabit() = (%#v, %v), want nil and NotFound", response, err)
	}
}

func TestHabitClient_EditHabit_DatabaseError(t *testing.T) {
	expectedErr := status.Error(codes.Internal, "database failure")
	name := "read"
	client := NewHabitClient(&habitServiceClientMock{editHabit: func(context.Context, *pb.EditHabitRequest, ...grpc.CallOption) (*pb.EditHabitResponse, error) {
		return nil, expectedErr
	}})
	response, err := client.EditHabit(context.Background(), "habit-id", &dto.EditHabitRequest{Name: &name})
	if response != nil || !errors.Is(err, expectedErr) {
		t.Fatalf("EditHabit() = (%#v, %v), want nil and Internal", response, err)
	}
}

func TestHabitClient_DeleteHabit_InvalidArgument(t *testing.T) {
	expectedErr := status.Error(codes.InvalidArgument, "invalid argument")
	client := NewHabitClient(&habitServiceClientMock{deleteHabit: func(context.Context, *pb.DeleteHabitRequest, ...grpc.CallOption) (*pb.DeleteHabitResponse, error) {
		return nil, expectedErr
	}})
	response, err := client.DeleteHabit(context.Background(), "invalid-id")
	if response != nil || !errors.Is(err, expectedErr) {
		t.Fatalf("DeleteHabit() = (%#v, %v), want nil and InvalidArgument", response, err)
	}
}

func TestHabitClient_DeleteHabit_HabitNotFound(t *testing.T) {
	expectedErr := status.Error(codes.NotFound, "habit not found")
	client := NewHabitClient(&habitServiceClientMock{deleteHabit: func(context.Context, *pb.DeleteHabitRequest, ...grpc.CallOption) (*pb.DeleteHabitResponse, error) {
		return nil, expectedErr
	}})
	response, err := client.DeleteHabit(context.Background(), "missing-id")
	if response != nil || !errors.Is(err, expectedErr) {
		t.Fatalf("DeleteHabit() = (%#v, %v), want nil and NotFound", response, err)
	}
}

func TestHabitClient_DeleteHabit_DatabaseError(t *testing.T) {
	expectedErr := status.Error(codes.Internal, "database failure")
	client := NewHabitClient(&habitServiceClientMock{deleteHabit: func(context.Context, *pb.DeleteHabitRequest, ...grpc.CallOption) (*pb.DeleteHabitResponse, error) {
		return nil, expectedErr
	}})
	response, err := client.DeleteHabit(context.Background(), "habit-id")
	if response != nil || !errors.Is(err, expectedErr) {
		t.Fatalf("DeleteHabit() = (%#v, %v), want nil and Internal", response, err)
	}
}

func TestHabitClient_ListHabitsByUser_InvalidArgument(t *testing.T) {
	expectedErr := status.Error(codes.InvalidArgument, "invalid argument")
	client := NewHabitClient(&habitServiceClientMock{listHabitsByUser: func(context.Context, *pb.ListHabitsByUserRequest, ...grpc.CallOption) (*pb.ListHabitsByUserResponse, error) {
		return nil, expectedErr
	}})
	response, err := client.ListHabitsByUser(context.Background(), "invalid-id")
	if response != nil || !errors.Is(err, expectedErr) {
		t.Fatalf("ListHabitsByUser() = (%#v, %v), want nil and InvalidArgument", response, err)
	}
}

func TestHabitClient_ListHabitsByUser_UserNotFound(t *testing.T) {
	expectedErr := status.Error(codes.NotFound, "user not found")
	client := NewHabitClient(&habitServiceClientMock{listHabitsByUser: func(context.Context, *pb.ListHabitsByUserRequest, ...grpc.CallOption) (*pb.ListHabitsByUserResponse, error) {
		return nil, expectedErr
	}})
	response, err := client.ListHabitsByUser(context.Background(), "missing-user-id")
	if response != nil || !errors.Is(err, expectedErr) {
		t.Fatalf("ListHabitsByUser() = (%#v, %v), want nil and NotFound", response, err)
	}
}

func TestHabitClient_ListHabitsByUser_DatabaseError(t *testing.T) {
	expectedErr := status.Error(codes.Internal, "database failure")
	client := NewHabitClient(&habitServiceClientMock{listHabitsByUser: func(context.Context, *pb.ListHabitsByUserRequest, ...grpc.CallOption) (*pb.ListHabitsByUserResponse, error) {
		return nil, expectedErr
	}})
	response, err := client.ListHabitsByUser(context.Background(), "user-id")
	if response != nil || !errors.Is(err, expectedErr) {
		t.Fatalf("ListHabitsByUser() = (%#v, %v), want nil and Internal", response, err)
	}
}

func TestHabitClient_ListHabitsByRoutine_InvalidArgument(t *testing.T) {
	expectedErr := status.Error(codes.InvalidArgument, "invalid argument")
	client := NewHabitClient(&habitServiceClientMock{listHabitsByRoutine: func(context.Context, *pb.ListHabitsByRoutineRequest, ...grpc.CallOption) (*pb.ListHabitsByRoutineResponse, error) {
		return nil, expectedErr
	}})
	response, err := client.ListHabitsByRoutine(context.Background(), "invalid-id")
	if response != nil || !errors.Is(err, expectedErr) {
		t.Fatalf("ListHabitsByRoutine() = (%#v, %v), want nil and InvalidArgument", response, err)
	}
}

func TestHabitClient_ListHabitsByRoutine_RoutineNotFound(t *testing.T) {
	expectedErr := status.Error(codes.NotFound, "routine not found")
	client := NewHabitClient(&habitServiceClientMock{listHabitsByRoutine: func(context.Context, *pb.ListHabitsByRoutineRequest, ...grpc.CallOption) (*pb.ListHabitsByRoutineResponse, error) {
		return nil, expectedErr
	}})
	response, err := client.ListHabitsByRoutine(context.Background(), "missing-routine-id")
	if response != nil || !errors.Is(err, expectedErr) {
		t.Fatalf("ListHabitsByRoutine() = (%#v, %v), want nil and NotFound", response, err)
	}
}

func TestHabitClient_ListHabitsByRoutine_DatabaseError(t *testing.T) {
	expectedErr := status.Error(codes.Internal, "database failure")
	client := NewHabitClient(&habitServiceClientMock{listHabitsByRoutine: func(context.Context, *pb.ListHabitsByRoutineRequest, ...grpc.CallOption) (*pb.ListHabitsByRoutineResponse, error) {
		return nil, expectedErr
	}})
	response, err := client.ListHabitsByRoutine(context.Background(), "routine-id")
	if response != nil || !errors.Is(err, expectedErr) {
		t.Fatalf("ListHabitsByRoutine() = (%#v, %v), want nil and Internal", response, err)
	}
}

func TestHabitClient_MarkHabitCompleted_InvalidArgument(t *testing.T) {
	expectedErr := status.Error(codes.InvalidArgument, "invalid argument")
	client := NewHabitClient(&habitServiceClientMock{markHabitCompleted: func(context.Context, *pb.MarkHabitCompletedRequest, ...grpc.CallOption) (*pb.MarkHabitCompletedResponse, error) {
		return nil, expectedErr
	}})
	response, err := client.MarkHabitCompleted(context.Background(), "invalid-id", testTime)
	if response != nil || !errors.Is(err, expectedErr) {
		t.Fatalf("MarkHabitCompleted() = (%#v, %v), want nil and InvalidArgument", response, err)
	}
}

func TestHabitClient_MarkHabitCompleted_HabitNotFound(t *testing.T) {
	expectedErr := status.Error(codes.NotFound, "habit not found")
	client := NewHabitClient(&habitServiceClientMock{markHabitCompleted: func(context.Context, *pb.MarkHabitCompletedRequest, ...grpc.CallOption) (*pb.MarkHabitCompletedResponse, error) {
		return nil, expectedErr
	}})
	response, err := client.MarkHabitCompleted(context.Background(), "missing-id", testTime)
	if response != nil || !errors.Is(err, expectedErr) {
		t.Fatalf("MarkHabitCompleted() = (%#v, %v), want nil and NotFound", response, err)
	}
}

func TestHabitClient_MarkHabitCompleted_DatabaseError(t *testing.T) {
	expectedErr := status.Error(codes.Internal, "database failure")
	client := NewHabitClient(&habitServiceClientMock{markHabitCompleted: func(context.Context, *pb.MarkHabitCompletedRequest, ...grpc.CallOption) (*pb.MarkHabitCompletedResponse, error) {
		return nil, expectedErr
	}})
	response, err := client.MarkHabitCompleted(context.Background(), "habit-id", testTime)
	if response != nil || !errors.Is(err, expectedErr) {
		t.Fatalf("MarkHabitCompleted() = (%#v, %v), want nil and Internal", response, err)
	}
}

func TestHabitClient_UnmarkHabitCompleted_InvalidArgument(t *testing.T) {
	expectedErr := status.Error(codes.InvalidArgument, "invalid argument")
	client := NewHabitClient(&habitServiceClientMock{unmarkHabitCompleted: func(context.Context, *pb.UnmarkHabitCompletedRequest, ...grpc.CallOption) (*pb.UnmarkHabitCompletedResponse, error) {
		return nil, expectedErr
	}})
	response, err := client.UnmarkHabitCompleted(context.Background(), "invalid-id", testTime)
	if response != nil || !errors.Is(err, expectedErr) {
		t.Fatalf("UnmarkHabitCompleted() = (%#v, %v), want nil and InvalidArgument", response, err)
	}
}

func TestHabitClient_UnmarkHabitCompleted_HabitNotFound(t *testing.T) {
	expectedErr := status.Error(codes.NotFound, "habit not found")
	client := NewHabitClient(&habitServiceClientMock{unmarkHabitCompleted: func(context.Context, *pb.UnmarkHabitCompletedRequest, ...grpc.CallOption) (*pb.UnmarkHabitCompletedResponse, error) {
		return nil, expectedErr
	}})
	response, err := client.UnmarkHabitCompleted(context.Background(), "missing-id", testTime)
	if response != nil || !errors.Is(err, expectedErr) {
		t.Fatalf("UnmarkHabitCompleted() = (%#v, %v), want nil and NotFound", response, err)
	}
}

func TestHabitClient_UnmarkHabitCompleted_DatabaseError(t *testing.T) {
	expectedErr := status.Error(codes.Internal, "database failure")
	client := NewHabitClient(&habitServiceClientMock{unmarkHabitCompleted: func(context.Context, *pb.UnmarkHabitCompletedRequest, ...grpc.CallOption) (*pb.UnmarkHabitCompletedResponse, error) {
		return nil, expectedErr
	}})
	response, err := client.UnmarkHabitCompleted(context.Background(), "habit-id", testTime)
	if response != nil || !errors.Is(err, expectedErr) {
		t.Fatalf("UnmarkHabitCompleted() = (%#v, %v), want nil and Internal", response, err)
	}
}

func TestHabitClient_GetHabitLogs_InvalidArgument(t *testing.T) {
	expectedErr := status.Error(codes.InvalidArgument, "invalid argument")
	client := NewHabitClient(&habitServiceClientMock{getHabitLogs: func(context.Context, *pb.GetHabitLogsRequest, ...grpc.CallOption) (*pb.GetHabitLogsResponse, error) {
		return nil, expectedErr
	}})
	response, err := client.GetHabitLogs(context.Background(), "invalid-id", testTime.Add(-time.Hour), testTime)
	if response != nil || !errors.Is(err, expectedErr) {
		t.Fatalf("GetHabitLogs() = (%#v, %v), want nil and InvalidArgument", response, err)
	}
}

func TestHabitClient_GetHabitLogs_HabitNotFound(t *testing.T) {
	expectedErr := status.Error(codes.NotFound, "habit not found")
	client := NewHabitClient(&habitServiceClientMock{getHabitLogs: func(context.Context, *pb.GetHabitLogsRequest, ...grpc.CallOption) (*pb.GetHabitLogsResponse, error) {
		return nil, expectedErr
	}})
	response, err := client.GetHabitLogs(context.Background(), "missing-id", testTime.Add(-time.Hour), testTime)
	if response != nil || !errors.Is(err, expectedErr) {
		t.Fatalf("GetHabitLogs() = (%#v, %v), want nil and NotFound", response, err)
	}
}

func TestHabitClient_GetHabitLogs_DatabaseError(t *testing.T) {
	expectedErr := status.Error(codes.Internal, "database failure")
	client := NewHabitClient(&habitServiceClientMock{getHabitLogs: func(context.Context, *pb.GetHabitLogsRequest, ...grpc.CallOption) (*pb.GetHabitLogsResponse, error) {
		return nil, expectedErr
	}})
	response, err := client.GetHabitLogs(context.Background(), "habit-id", testTime.Add(-time.Hour), testTime)
	if response != nil || !errors.Is(err, expectedErr) {
		t.Fatalf("GetHabitLogs() = (%#v, %v), want nil and Internal", response, err)
	}
}

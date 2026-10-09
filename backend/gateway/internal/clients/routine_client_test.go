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

type routineServiceClientMock struct {
	createRoutine          func(context.Context, *pb.CreateRoutineRequest, ...grpc.CallOption) (*pb.CreateRoutineResponse, error)
	getRoutineByID         func(context.Context, *pb.GetRoutineByIDRequest, ...grpc.CallOption) (*pb.GetRoutineByIDResponse, error)
	editRoutine            func(context.Context, *pb.EditRoutineRequest, ...grpc.CallOption) (*pb.EditRoutineResponse, error)
	deleteRoutine          func(context.Context, *pb.DeleteRoutineRequest, ...grpc.CallOption) (*pb.DeleteRoutineResponse, error)
	addHabitToRoutine      func(context.Context, *pb.AddHabitToRoutineRequest, ...grpc.CallOption) (*pb.AddHabitToRoutineResponse, error)
	removeHabitFromRoutine func(context.Context, *pb.RemoveHabitFromRoutineRequest, ...grpc.CallOption) (*pb.RemoveHabitFromRoutineResponse, error)
	listRoutinesByUser     func(context.Context, *pb.ListRoutinesByUserRequest, ...grpc.CallOption) (*pb.ListRoutinesByUserResponse, error)
	markRoutineCompleted   func(context.Context, *pb.MarkRoutineCompletedRequest, ...grpc.CallOption) (*pb.MarkRoutineCompletedResponse, error)
	unmarkRoutineCompleted func(context.Context, *pb.UnmarkRoutineCompletedRequest, ...grpc.CallOption) (*pb.UnmarkRoutineCompletedResponse, error)
	getRoutineLogs         func(context.Context, *pb.GetRoutineLogsRequest, ...grpc.CallOption) (*pb.GetRoutineLogsResponse, error)
}

func (m *routineServiceClientMock) CreateRoutine(ctx context.Context, req *pb.CreateRoutineRequest, opts ...grpc.CallOption) (*pb.CreateRoutineResponse, error) {
	return m.createRoutine(ctx, req, opts...)
}
func (m *routineServiceClientMock) GetRoutineByID(ctx context.Context, req *pb.GetRoutineByIDRequest, opts ...grpc.CallOption) (*pb.GetRoutineByIDResponse, error) {
	return m.getRoutineByID(ctx, req, opts...)
}
func (m *routineServiceClientMock) EditRoutine(ctx context.Context, req *pb.EditRoutineRequest, opts ...grpc.CallOption) (*pb.EditRoutineResponse, error) {
	return m.editRoutine(ctx, req, opts...)
}
func (m *routineServiceClientMock) DeleteRoutine(ctx context.Context, req *pb.DeleteRoutineRequest, opts ...grpc.CallOption) (*pb.DeleteRoutineResponse, error) {
	return m.deleteRoutine(ctx, req, opts...)
}
func (m *routineServiceClientMock) AddHabitToRoutine(ctx context.Context, req *pb.AddHabitToRoutineRequest, opts ...grpc.CallOption) (*pb.AddHabitToRoutineResponse, error) {
	return m.addHabitToRoutine(ctx, req, opts...)
}
func (m *routineServiceClientMock) RemoveHabitFromRoutine(ctx context.Context, req *pb.RemoveHabitFromRoutineRequest, opts ...grpc.CallOption) (*pb.RemoveHabitFromRoutineResponse, error) {
	return m.removeHabitFromRoutine(ctx, req, opts...)
}
func (m *routineServiceClientMock) ListRoutinesByUser(ctx context.Context, req *pb.ListRoutinesByUserRequest, opts ...grpc.CallOption) (*pb.ListRoutinesByUserResponse, error) {
	return m.listRoutinesByUser(ctx, req, opts...)
}
func (m *routineServiceClientMock) MarkRoutineCompleted(ctx context.Context, req *pb.MarkRoutineCompletedRequest, opts ...grpc.CallOption) (*pb.MarkRoutineCompletedResponse, error) {
	return m.markRoutineCompleted(ctx, req, opts...)
}
func (m *routineServiceClientMock) UnmarkRoutineCompleted(ctx context.Context, req *pb.UnmarkRoutineCompletedRequest, opts ...grpc.CallOption) (*pb.UnmarkRoutineCompletedResponse, error) {
	return m.unmarkRoutineCompleted(ctx, req, opts...)
}
func (m *routineServiceClientMock) GetRoutineLogs(ctx context.Context, req *pb.GetRoutineLogsRequest, opts ...grpc.CallOption) (*pb.GetRoutineLogsResponse, error) {
	return m.getRoutineLogs(ctx, req, opts...)
}

var routineTestTime = time.Date(2026, time.October, 8, 12, 0, 0, 0, time.UTC)

func TestRoutineClient_CreateRoutine(t *testing.T) {
	client := NewRoutineClient(&routineServiceClientMock{
		createRoutine: func(_ context.Context, req *pb.CreateRoutineRequest, _ ...grpc.CallOption) (*pb.CreateRoutineResponse, error) {
			expected := &pb.CreateRoutineRequest{Routine: &pb.Routine{UserId: "user-id", Name: "morning"}}
			if !proto.Equal(req, expected) {
				t.Errorf("CreateRoutine() request = %v, want %v", req, expected)
			}
			return &pb.CreateRoutineResponse{Routine: &pb.Routine{Id: "routine-id", UserId: "user-id", Name: "morning", CreatedAt: timestamppb.New(routineTestTime)}}, nil
		},
	})
	response, err := client.CreateRoutine(context.Background(), "user-id", &dto.CreateRoutineRequest{Name: "morning"})
	if err != nil {
		t.Fatalf("CreateRoutine() returned unexpected error: %v", err)
	}
	if response.Routine.ID != "routine-id" || response.Routine.CreatedAt != routineTestTime {
		t.Fatalf("CreateRoutine() returned unexpected result: %#v", response)
	}
}

func TestRoutineClient_GetRoutineByID(t *testing.T) {
	client := NewRoutineClient(&routineServiceClientMock{
		getRoutineByID: func(_ context.Context, req *pb.GetRoutineByIDRequest, _ ...grpc.CallOption) (*pb.GetRoutineByIDResponse, error) {
			expected := &pb.GetRoutineByIDRequest{RoutineId: "routine-id"}
			if !proto.Equal(req, expected) {
				t.Errorf("GetRoutineByID() request = %v, want %v", req, expected)
			}
			return &pb.GetRoutineByIDResponse{Routine: &pb.Routine{Id: "routine-id", UserId: "user-id", Name: "morning", CreatedAt: timestamppb.New(routineTestTime)}}, nil
		},
	})
	response, err := client.GetRoutineByID(context.Background(), "routine-id")
	if err != nil {
		t.Fatalf("GetRoutineByID() returned unexpected error: %v", err)
	}
	if response.Routine.ID != "routine-id" || response.Routine.Name != "morning" {
		t.Fatalf("GetRoutineByID() returned unexpected result: %#v", response)
	}
}

func TestRoutineClient_EditRoutine(t *testing.T) {
	name := "morning"
	client := NewRoutineClient(&routineServiceClientMock{
		editRoutine: func(_ context.Context, req *pb.EditRoutineRequest, _ ...grpc.CallOption) (*pb.EditRoutineResponse, error) {
			expected := &pb.EditRoutineRequest{RoutineId: "routine-id", Name: &name}
			if !proto.Equal(req, expected) {
				t.Errorf("EditRoutine() request = %v, want %v", req, expected)
			}
			return &pb.EditRoutineResponse{Routine: &pb.Routine{Id: "routine-id", UserId: "user-id", Name: "morning", CreatedAt: timestamppb.New(routineTestTime)}}, nil
		},
	})
	response, err := client.EditRoutine(context.Background(), "routine-id", &dto.EditRoutineRequest{Name: &name})
	if err != nil {
		t.Fatalf("EditRoutine() returned unexpected error: %v", err)
	}
	if response.Routine.ID != "routine-id" || response.Routine.Name != "morning" {
		t.Fatalf("EditRoutine() returned unexpected result: %#v", response)
	}
}

func TestRoutineClient_DeleteRoutine(t *testing.T) {
	client := NewRoutineClient(&routineServiceClientMock{
		deleteRoutine: func(_ context.Context, req *pb.DeleteRoutineRequest, _ ...grpc.CallOption) (*pb.DeleteRoutineResponse, error) {
			expected := &pb.DeleteRoutineRequest{RoutineId: "routine-id"}
			if !proto.Equal(req, expected) {
				t.Errorf("DeleteRoutine() request = %v, want %v", req, expected)
			}
			return &pb.DeleteRoutineResponse{Success: true}, nil
		},
	})
	response, err := client.DeleteRoutine(context.Background(), "routine-id")
	if err != nil {
		t.Fatalf("DeleteRoutine() returned unexpected error: %v", err)
	}
	if !response.Success {
		t.Fatalf("DeleteRoutine() returned unexpected result: %#v", response)
	}
}

func TestRoutineClient_AddHabitToRoutine(t *testing.T) {
	client := NewRoutineClient(&routineServiceClientMock{
		addHabitToRoutine: func(_ context.Context, req *pb.AddHabitToRoutineRequest, _ ...grpc.CallOption) (*pb.AddHabitToRoutineResponse, error) {
			expected := &pb.AddHabitToRoutineRequest{RoutineId: "routine-id", HabitId: "habit-id"}
			if !proto.Equal(req, expected) {
				t.Errorf("AddHabitToRoutine() request = %v, want %v", req, expected)
			}
			return &pb.AddHabitToRoutineResponse{Success: true}, nil
		},
	})
	response, err := client.AddHabitToRoutine(context.Background(), "routine-id", &dto.RoutineHabitRequest{HabitID: "habit-id"})
	if err != nil {
		t.Fatalf("AddHabitToRoutine() returned unexpected error: %v", err)
	}
	if !response.Success {
		t.Fatalf("AddHabitToRoutine() returned unexpected result: %#v", response)
	}
}

func TestRoutineClient_RemoveHabitFromRoutine(t *testing.T) {
	client := NewRoutineClient(&routineServiceClientMock{
		removeHabitFromRoutine: func(_ context.Context, req *pb.RemoveHabitFromRoutineRequest, _ ...grpc.CallOption) (*pb.RemoveHabitFromRoutineResponse, error) {
			expected := &pb.RemoveHabitFromRoutineRequest{RoutineId: "routine-id", HabitId: "habit-id"}
			if !proto.Equal(req, expected) {
				t.Errorf("RemoveHabitFromRoutine() request = %v, want %v", req, expected)
			}
			return &pb.RemoveHabitFromRoutineResponse{Success: true}, nil
		},
	})
	response, err := client.RemoveHabitFromRoutine(context.Background(), "routine-id", &dto.RoutineHabitRequest{HabitID: "habit-id"})
	if err != nil {
		t.Fatalf("RemoveHabitFromRoutine() returned unexpected error: %v", err)
	}
	if !response.Success {
		t.Fatalf("RemoveHabitFromRoutine() returned unexpected result: %#v", response)
	}
}

func TestRoutineClient_ListRoutinesByUser(t *testing.T) {
	client := NewRoutineClient(&routineServiceClientMock{
		listRoutinesByUser: func(_ context.Context, req *pb.ListRoutinesByUserRequest, _ ...grpc.CallOption) (*pb.ListRoutinesByUserResponse, error) {
			expected := &pb.ListRoutinesByUserRequest{UserId: "user-id"}
			if !proto.Equal(req, expected) {
				t.Errorf("ListRoutinesByUser() request = %v, want %v", req, expected)
			}
			return &pb.ListRoutinesByUserResponse{Routines: []*pb.Routine{{Id: "routine-id", UserId: "user-id", Name: "morning", CreatedAt: timestamppb.New(routineTestTime)}}}, nil
		},
	})
	response, err := client.ListRoutinesByUser(context.Background(), "user-id")
	if err != nil {
		t.Fatalf("ListRoutinesByUser() returned unexpected error: %v", err)
	}
	if len(response.Routines) != 1 || response.Routines[0].ID != "routine-id" {
		t.Fatalf("ListRoutinesByUser() returned unexpected result: %#v", response)
	}
}

func TestRoutineClient_MarkRoutineCompleted(t *testing.T) {
	client := NewRoutineClient(&routineServiceClientMock{
		markRoutineCompleted: func(_ context.Context, req *pb.MarkRoutineCompletedRequest, _ ...grpc.CallOption) (*pb.MarkRoutineCompletedResponse, error) {
			expected := &pb.MarkRoutineCompletedRequest{RoutineId: "routine-id", CompletedAt: timestamppb.New(routineTestTime)}
			if !proto.Equal(req, expected) {
				t.Errorf("MarkRoutineCompleted() request = %v, want %v", req, expected)
			}
			return &pb.MarkRoutineCompletedResponse{Success: true}, nil
		},
	})
	response, err := client.MarkRoutineCompleted(context.Background(), "routine-id", routineTestTime)
	if err != nil {
		t.Fatalf("MarkRoutineCompleted() returned unexpected error: %v", err)
	}
	if !response.Success {
		t.Fatalf("MarkRoutineCompleted() returned unexpected result: %#v", response)
	}
}

func TestRoutineClient_UnmarkRoutineCompleted(t *testing.T) {
	client := NewRoutineClient(&routineServiceClientMock{
		unmarkRoutineCompleted: func(_ context.Context, req *pb.UnmarkRoutineCompletedRequest, _ ...grpc.CallOption) (*pb.UnmarkRoutineCompletedResponse, error) {
			expected := &pb.UnmarkRoutineCompletedRequest{RoutineId: "routine-id", CompletedAt: timestamppb.New(routineTestTime)}
			if !proto.Equal(req, expected) {
				t.Errorf("UnmarkRoutineCompleted() request = %v, want %v", req, expected)
			}
			return &pb.UnmarkRoutineCompletedResponse{Success: true}, nil
		},
	})
	response, err := client.UnmarkRoutineCompleted(context.Background(), "routine-id", routineTestTime)
	if err != nil {
		t.Fatalf("UnmarkRoutineCompleted() returned unexpected error: %v", err)
	}
	if !response.Success {
		t.Fatalf("UnmarkRoutineCompleted() returned unexpected result: %#v", response)
	}
}

func TestRoutineClient_GetRoutineLogs(t *testing.T) {
	client := NewRoutineClient(&routineServiceClientMock{
		getRoutineLogs: func(_ context.Context, req *pb.GetRoutineLogsRequest, _ ...grpc.CallOption) (*pb.GetRoutineLogsResponse, error) {
			expected := &pb.GetRoutineLogsRequest{RoutineId: "routine-id", StartDate: timestamppb.New(routineTestTime.Add(-time.Hour)), EndDate: timestamppb.New(routineTestTime)}
			if !proto.Equal(req, expected) {
				t.Errorf("GetRoutineLogs() request = %v, want %v", req, expected)
			}
			return &pb.GetRoutineLogsResponse{Logs: []*pb.RoutineLog{{RoutineId: "routine-id", CompletedAt: timestamppb.New(routineTestTime)}}}, nil
		},
	})
	response, err := client.GetRoutineLogs(context.Background(), "routine-id", routineTestTime.Add(-time.Hour), routineTestTime)
	if err != nil {
		t.Fatalf("GetRoutineLogs() returned unexpected error: %v", err)
	}
	if len(response.Logs) != 1 || response.Logs[0].RoutineID != "routine-id" || response.Logs[0].CompletedAt != routineTestTime {
		t.Fatalf("GetRoutineLogs() returned unexpected result: %#v", response)
	}
}

func TestRoutineClient_CreateRoutine_InvalidArgument(t *testing.T) {
	expectedErr := status.Error(codes.InvalidArgument, "invalid argument")
	client := NewRoutineClient(&routineServiceClientMock{createRoutine: func(context.Context, *pb.CreateRoutineRequest, ...grpc.CallOption) (*pb.CreateRoutineResponse, error) {
		return nil, expectedErr
	}})
	response, err := client.CreateRoutine(context.Background(), "invalid-user-id", &dto.CreateRoutineRequest{Name: "morning"})
	if response != nil || !errors.Is(err, expectedErr) {
		t.Fatalf("CreateRoutine() = (%#v, %v), want nil and InvalidArgument", response, err)
	}
}

func TestRoutineClient_CreateRoutine_UserNotFound(t *testing.T) {
	expectedErr := status.Error(codes.NotFound, "user not found")
	client := NewRoutineClient(&routineServiceClientMock{createRoutine: func(context.Context, *pb.CreateRoutineRequest, ...grpc.CallOption) (*pb.CreateRoutineResponse, error) {
		return nil, expectedErr
	}})
	response, err := client.CreateRoutine(context.Background(), "missing-user-id", &dto.CreateRoutineRequest{Name: "morning"})
	if response != nil || !errors.Is(err, expectedErr) {
		t.Fatalf("CreateRoutine() = (%#v, %v), want nil and NotFound", response, err)
	}
}

func TestRoutineClient_CreateRoutine_DatabaseError(t *testing.T) {
	expectedErr := status.Error(codes.Internal, "database failure")
	client := NewRoutineClient(&routineServiceClientMock{createRoutine: func(context.Context, *pb.CreateRoutineRequest, ...grpc.CallOption) (*pb.CreateRoutineResponse, error) {
		return nil, expectedErr
	}})
	response, err := client.CreateRoutine(context.Background(), "user-id", &dto.CreateRoutineRequest{Name: "morning"})
	if response != nil || !errors.Is(err, expectedErr) {
		t.Fatalf("CreateRoutine() = (%#v, %v), want nil and Internal", response, err)
	}
}

func TestRoutineClient_GetRoutineByID_InvalidArgument(t *testing.T) {
	expectedErr := status.Error(codes.InvalidArgument, "invalid argument")
	client := NewRoutineClient(&routineServiceClientMock{getRoutineByID: func(context.Context, *pb.GetRoutineByIDRequest, ...grpc.CallOption) (*pb.GetRoutineByIDResponse, error) {
		return nil, expectedErr
	}})
	response, err := client.GetRoutineByID(context.Background(), "invalid-id")
	if response != nil || !errors.Is(err, expectedErr) {
		t.Fatalf("GetRoutineByID() = (%#v, %v), want nil and InvalidArgument", response, err)
	}
}

func TestRoutineClient_GetRoutineByID_NotFound(t *testing.T) {
	expectedErr := status.Error(codes.NotFound, "routine not found")
	client := NewRoutineClient(&routineServiceClientMock{getRoutineByID: func(context.Context, *pb.GetRoutineByIDRequest, ...grpc.CallOption) (*pb.GetRoutineByIDResponse, error) {
		return nil, expectedErr
	}})
	response, err := client.GetRoutineByID(context.Background(), "missing-id")
	if response != nil || !errors.Is(err, expectedErr) {
		t.Fatalf("GetRoutineByID() = (%#v, %v), want nil and NotFound", response, err)
	}
}

func TestRoutineClient_GetRoutineByID_DatabaseError(t *testing.T) {
	expectedErr := status.Error(codes.Internal, "database failure")
	client := NewRoutineClient(&routineServiceClientMock{getRoutineByID: func(context.Context, *pb.GetRoutineByIDRequest, ...grpc.CallOption) (*pb.GetRoutineByIDResponse, error) {
		return nil, expectedErr
	}})
	response, err := client.GetRoutineByID(context.Background(), "routine-id")
	if response != nil || !errors.Is(err, expectedErr) {
		t.Fatalf("GetRoutineByID() = (%#v, %v), want nil and Internal", response, err)
	}
}

func TestRoutineClient_EditRoutine_InvalidArgument(t *testing.T) {
	expectedErr := status.Error(codes.InvalidArgument, "invalid argument")
	name := "morning"
	client := NewRoutineClient(&routineServiceClientMock{editRoutine: func(context.Context, *pb.EditRoutineRequest, ...grpc.CallOption) (*pb.EditRoutineResponse, error) {
		return nil, expectedErr
	}})
	response, err := client.EditRoutine(context.Background(), "invalid-id", &dto.EditRoutineRequest{Name: &name})
	if response != nil || !errors.Is(err, expectedErr) {
		t.Fatalf("EditRoutine() = (%#v, %v), want nil and InvalidArgument", response, err)
	}
}

func TestRoutineClient_EditRoutine_RoutineNotFound(t *testing.T) {
	expectedErr := status.Error(codes.NotFound, "routine not found")
	name := "morning"
	client := NewRoutineClient(&routineServiceClientMock{editRoutine: func(context.Context, *pb.EditRoutineRequest, ...grpc.CallOption) (*pb.EditRoutineResponse, error) {
		return nil, expectedErr
	}})
	response, err := client.EditRoutine(context.Background(), "missing-id", &dto.EditRoutineRequest{Name: &name})
	if response != nil || !errors.Is(err, expectedErr) {
		t.Fatalf("EditRoutine() = (%#v, %v), want nil and NotFound", response, err)
	}
}

func TestRoutineClient_EditRoutine_DatabaseError(t *testing.T) {
	expectedErr := status.Error(codes.Internal, "database failure")
	name := "morning"
	client := NewRoutineClient(&routineServiceClientMock{editRoutine: func(context.Context, *pb.EditRoutineRequest, ...grpc.CallOption) (*pb.EditRoutineResponse, error) {
		return nil, expectedErr
	}})
	response, err := client.EditRoutine(context.Background(), "routine-id", &dto.EditRoutineRequest{Name: &name})
	if response != nil || !errors.Is(err, expectedErr) {
		t.Fatalf("EditRoutine() = (%#v, %v), want nil and Internal", response, err)
	}
}

func TestRoutineClient_DeleteRoutine_InvalidArgument(t *testing.T) {
	expectedErr := status.Error(codes.InvalidArgument, "invalid argument")
	client := NewRoutineClient(&routineServiceClientMock{deleteRoutine: func(context.Context, *pb.DeleteRoutineRequest, ...grpc.CallOption) (*pb.DeleteRoutineResponse, error) {
		return nil, expectedErr
	}})
	response, err := client.DeleteRoutine(context.Background(), "invalid-id")
	if response != nil || !errors.Is(err, expectedErr) {
		t.Fatalf("DeleteRoutine() = (%#v, %v), want nil and InvalidArgument", response, err)
	}
}

func TestRoutineClient_DeleteRoutine_RoutineNotFound(t *testing.T) {
	expectedErr := status.Error(codes.NotFound, "routine not found")
	client := NewRoutineClient(&routineServiceClientMock{deleteRoutine: func(context.Context, *pb.DeleteRoutineRequest, ...grpc.CallOption) (*pb.DeleteRoutineResponse, error) {
		return nil, expectedErr
	}})
	response, err := client.DeleteRoutine(context.Background(), "missing-id")
	if response != nil || !errors.Is(err, expectedErr) {
		t.Fatalf("DeleteRoutine() = (%#v, %v), want nil and NotFound", response, err)
	}
}

func TestRoutineClient_DeleteRoutine_DatabaseError(t *testing.T) {
	expectedErr := status.Error(codes.Internal, "database failure")
	client := NewRoutineClient(&routineServiceClientMock{deleteRoutine: func(context.Context, *pb.DeleteRoutineRequest, ...grpc.CallOption) (*pb.DeleteRoutineResponse, error) {
		return nil, expectedErr
	}})
	response, err := client.DeleteRoutine(context.Background(), "routine-id")
	if response != nil || !errors.Is(err, expectedErr) {
		t.Fatalf("DeleteRoutine() = (%#v, %v), want nil and Internal", response, err)
	}
}

func TestRoutineClient_ListRoutinesByUser_InvalidArgument(t *testing.T) {
	expectedErr := status.Error(codes.InvalidArgument, "invalid argument")
	client := NewRoutineClient(&routineServiceClientMock{listRoutinesByUser: func(context.Context, *pb.ListRoutinesByUserRequest, ...grpc.CallOption) (*pb.ListRoutinesByUserResponse, error) {
		return nil, expectedErr
	}})
	response, err := client.ListRoutinesByUser(context.Background(), "invalid-id")
	if response != nil || !errors.Is(err, expectedErr) {
		t.Fatalf("ListRoutinesByUser() = (%#v, %v), want nil and InvalidArgument", response, err)
	}
}

func TestRoutineClient_ListRoutinesByUser_UserNotFound(t *testing.T) {
	expectedErr := status.Error(codes.NotFound, "user not found")
	client := NewRoutineClient(&routineServiceClientMock{listRoutinesByUser: func(context.Context, *pb.ListRoutinesByUserRequest, ...grpc.CallOption) (*pb.ListRoutinesByUserResponse, error) {
		return nil, expectedErr
	}})
	response, err := client.ListRoutinesByUser(context.Background(), "missing-user-id")
	if response != nil || !errors.Is(err, expectedErr) {
		t.Fatalf("ListRoutinesByUser() = (%#v, %v), want nil and NotFound", response, err)
	}
}

func TestRoutineClient_ListRoutinesByUser_DatabaseError(t *testing.T) {
	expectedErr := status.Error(codes.Internal, "database failure")
	client := NewRoutineClient(&routineServiceClientMock{listRoutinesByUser: func(context.Context, *pb.ListRoutinesByUserRequest, ...grpc.CallOption) (*pb.ListRoutinesByUserResponse, error) {
		return nil, expectedErr
	}})
	response, err := client.ListRoutinesByUser(context.Background(), "user-id")
	if response != nil || !errors.Is(err, expectedErr) {
		t.Fatalf("ListRoutinesByUser() = (%#v, %v), want nil and Internal", response, err)
	}
}

func TestRoutineClient_AddHabitToRoutine_InvalidArgument(t *testing.T) {
	expectedErr := status.Error(codes.InvalidArgument, "invalid argument")
	client := NewRoutineClient(&routineServiceClientMock{addHabitToRoutine: func(context.Context, *pb.AddHabitToRoutineRequest, ...grpc.CallOption) (*pb.AddHabitToRoutineResponse, error) {
		return nil, expectedErr
	}})
	response, err := client.AddHabitToRoutine(context.Background(), "invalid-id", &dto.RoutineHabitRequest{HabitID: "habit-id"})
	if response != nil || !errors.Is(err, expectedErr) {
		t.Fatalf("AddHabitToRoutine() = (%#v, %v), want nil and InvalidArgument", response, err)
	}
}

func TestRoutineClient_AddHabitToRoutine_RoutineNotFound(t *testing.T) {
	expectedErr := status.Error(codes.NotFound, "routine not found")
	client := NewRoutineClient(&routineServiceClientMock{addHabitToRoutine: func(context.Context, *pb.AddHabitToRoutineRequest, ...grpc.CallOption) (*pb.AddHabitToRoutineResponse, error) {
		return nil, expectedErr
	}})
	response, err := client.AddHabitToRoutine(context.Background(), "missing-routine-id", &dto.RoutineHabitRequest{HabitID: "habit-id"})
	if response != nil || !errors.Is(err, expectedErr) {
		t.Fatalf("AddHabitToRoutine() = (%#v, %v), want nil and NotFound", response, err)
	}
}

func TestRoutineClient_AddHabitToRoutine_HabitNotFound(t *testing.T) {
	expectedErr := status.Error(codes.NotFound, "habit not found")
	client := NewRoutineClient(&routineServiceClientMock{addHabitToRoutine: func(context.Context, *pb.AddHabitToRoutineRequest, ...grpc.CallOption) (*pb.AddHabitToRoutineResponse, error) {
		return nil, expectedErr
	}})
	response, err := client.AddHabitToRoutine(context.Background(), "routine-id", &dto.RoutineHabitRequest{HabitID: "missing-habit-id"})
	if response != nil || !errors.Is(err, expectedErr) {
		t.Fatalf("AddHabitToRoutine() = (%#v, %v), want nil and NotFound", response, err)
	}
}

func TestRoutineClient_AddHabitToRoutine_DatabaseError(t *testing.T) {
	expectedErr := status.Error(codes.Internal, "database failure")
	client := NewRoutineClient(&routineServiceClientMock{addHabitToRoutine: func(context.Context, *pb.AddHabitToRoutineRequest, ...grpc.CallOption) (*pb.AddHabitToRoutineResponse, error) {
		return nil, expectedErr
	}})
	response, err := client.AddHabitToRoutine(context.Background(), "routine-id", &dto.RoutineHabitRequest{HabitID: "habit-id"})
	if response != nil || !errors.Is(err, expectedErr) {
		t.Fatalf("AddHabitToRoutine() = (%#v, %v), want nil and Internal", response, err)
	}
}

func TestRoutineClient_RemoveHabitFromRoutine_InvalidArgument(t *testing.T) {
	expectedErr := status.Error(codes.InvalidArgument, "invalid argument")
	client := NewRoutineClient(&routineServiceClientMock{removeHabitFromRoutine: func(context.Context, *pb.RemoveHabitFromRoutineRequest, ...grpc.CallOption) (*pb.RemoveHabitFromRoutineResponse, error) {
		return nil, expectedErr
	}})
	response, err := client.RemoveHabitFromRoutine(context.Background(), "invalid-id", &dto.RoutineHabitRequest{HabitID: "habit-id"})
	if response != nil || !errors.Is(err, expectedErr) {
		t.Fatalf("RemoveHabitFromRoutine() = (%#v, %v), want nil and InvalidArgument", response, err)
	}
}

func TestRoutineClient_RemoveHabitFromRoutine_RoutineNotFound(t *testing.T) {
	expectedErr := status.Error(codes.NotFound, "routine not found")
	client := NewRoutineClient(&routineServiceClientMock{removeHabitFromRoutine: func(context.Context, *pb.RemoveHabitFromRoutineRequest, ...grpc.CallOption) (*pb.RemoveHabitFromRoutineResponse, error) {
		return nil, expectedErr
	}})
	response, err := client.RemoveHabitFromRoutine(context.Background(), "missing-routine-id", &dto.RoutineHabitRequest{HabitID: "habit-id"})
	if response != nil || !errors.Is(err, expectedErr) {
		t.Fatalf("RemoveHabitFromRoutine() = (%#v, %v), want nil and NotFound", response, err)
	}
}

func TestRoutineClient_RemoveHabitFromRoutine_HabitNotFound(t *testing.T) {
	expectedErr := status.Error(codes.NotFound, "habit not found")
	client := NewRoutineClient(&routineServiceClientMock{removeHabitFromRoutine: func(context.Context, *pb.RemoveHabitFromRoutineRequest, ...grpc.CallOption) (*pb.RemoveHabitFromRoutineResponse, error) {
		return nil, expectedErr
	}})
	response, err := client.RemoveHabitFromRoutine(context.Background(), "routine-id", &dto.RoutineHabitRequest{HabitID: "missing-habit-id"})
	if response != nil || !errors.Is(err, expectedErr) {
		t.Fatalf("RemoveHabitFromRoutine() = (%#v, %v), want nil and NotFound", response, err)
	}
}

func TestRoutineClient_RemoveHabitFromRoutine_DatabaseError(t *testing.T) {
	expectedErr := status.Error(codes.Internal, "database failure")
	client := NewRoutineClient(&routineServiceClientMock{removeHabitFromRoutine: func(context.Context, *pb.RemoveHabitFromRoutineRequest, ...grpc.CallOption) (*pb.RemoveHabitFromRoutineResponse, error) {
		return nil, expectedErr
	}})
	response, err := client.RemoveHabitFromRoutine(context.Background(), "routine-id", &dto.RoutineHabitRequest{HabitID: "habit-id"})
	if response != nil || !errors.Is(err, expectedErr) {
		t.Fatalf("RemoveHabitFromRoutine() = (%#v, %v), want nil and Internal", response, err)
	}
}

func TestRoutineClient_MarkRoutineCompleted_InvalidArgument(t *testing.T) {
	expectedErr := status.Error(codes.InvalidArgument, "invalid argument")
	client := NewRoutineClient(&routineServiceClientMock{markRoutineCompleted: func(context.Context, *pb.MarkRoutineCompletedRequest, ...grpc.CallOption) (*pb.MarkRoutineCompletedResponse, error) {
		return nil, expectedErr
	}})
	response, err := client.MarkRoutineCompleted(context.Background(), "invalid-id", routineTestTime)
	if response != nil || !errors.Is(err, expectedErr) {
		t.Fatalf("MarkRoutineCompleted() = (%#v, %v), want nil and InvalidArgument", response, err)
	}
}

func TestRoutineClient_MarkRoutineCompleted_RoutineNotFound(t *testing.T) {
	expectedErr := status.Error(codes.NotFound, "routine not found")
	client := NewRoutineClient(&routineServiceClientMock{markRoutineCompleted: func(context.Context, *pb.MarkRoutineCompletedRequest, ...grpc.CallOption) (*pb.MarkRoutineCompletedResponse, error) {
		return nil, expectedErr
	}})
	response, err := client.MarkRoutineCompleted(context.Background(), "missing-id", routineTestTime)
	if response != nil || !errors.Is(err, expectedErr) {
		t.Fatalf("MarkRoutineCompleted() = (%#v, %v), want nil and NotFound", response, err)
	}
}

func TestRoutineClient_MarkRoutineCompleted_DatabaseError(t *testing.T) {
	expectedErr := status.Error(codes.Internal, "database failure")
	client := NewRoutineClient(&routineServiceClientMock{markRoutineCompleted: func(context.Context, *pb.MarkRoutineCompletedRequest, ...grpc.CallOption) (*pb.MarkRoutineCompletedResponse, error) {
		return nil, expectedErr
	}})
	response, err := client.MarkRoutineCompleted(context.Background(), "routine-id", routineTestTime)
	if response != nil || !errors.Is(err, expectedErr) {
		t.Fatalf("MarkRoutineCompleted() = (%#v, %v), want nil and Internal", response, err)
	}
}

func TestRoutineClient_UnmarkRoutineCompleted_InvalidArgument(t *testing.T) {
	expectedErr := status.Error(codes.InvalidArgument, "invalid argument")
	client := NewRoutineClient(&routineServiceClientMock{unmarkRoutineCompleted: func(context.Context, *pb.UnmarkRoutineCompletedRequest, ...grpc.CallOption) (*pb.UnmarkRoutineCompletedResponse, error) {
		return nil, expectedErr
	}})
	response, err := client.UnmarkRoutineCompleted(context.Background(), "invalid-id", routineTestTime)
	if response != nil || !errors.Is(err, expectedErr) {
		t.Fatalf("UnmarkRoutineCompleted() = (%#v, %v), want nil and InvalidArgument", response, err)
	}
}

func TestRoutineClient_UnmarkRoutineCompleted_RoutineNotFound(t *testing.T) {
	expectedErr := status.Error(codes.NotFound, "routine not found")
	client := NewRoutineClient(&routineServiceClientMock{unmarkRoutineCompleted: func(context.Context, *pb.UnmarkRoutineCompletedRequest, ...grpc.CallOption) (*pb.UnmarkRoutineCompletedResponse, error) {
		return nil, expectedErr
	}})
	response, err := client.UnmarkRoutineCompleted(context.Background(), "missing-id", routineTestTime)
	if response != nil || !errors.Is(err, expectedErr) {
		t.Fatalf("UnmarkRoutineCompleted() = (%#v, %v), want nil and NotFound", response, err)
	}
}

func TestRoutineClient_UnmarkRoutineCompleted_DatabaseError(t *testing.T) {
	expectedErr := status.Error(codes.Internal, "database failure")
	client := NewRoutineClient(&routineServiceClientMock{unmarkRoutineCompleted: func(context.Context, *pb.UnmarkRoutineCompletedRequest, ...grpc.CallOption) (*pb.UnmarkRoutineCompletedResponse, error) {
		return nil, expectedErr
	}})
	response, err := client.UnmarkRoutineCompleted(context.Background(), "routine-id", routineTestTime)
	if response != nil || !errors.Is(err, expectedErr) {
		t.Fatalf("UnmarkRoutineCompleted() = (%#v, %v), want nil and Internal", response, err)
	}
}

func TestRoutineClient_GetRoutineLogs_InvalidArgument(t *testing.T) {
	expectedErr := status.Error(codes.InvalidArgument, "invalid argument")
	client := NewRoutineClient(&routineServiceClientMock{getRoutineLogs: func(context.Context, *pb.GetRoutineLogsRequest, ...grpc.CallOption) (*pb.GetRoutineLogsResponse, error) {
		return nil, expectedErr
	}})
	response, err := client.GetRoutineLogs(context.Background(), "invalid-id", routineTestTime.Add(-time.Hour), routineTestTime)
	if response != nil || !errors.Is(err, expectedErr) {
		t.Fatalf("GetRoutineLogs() = (%#v, %v), want nil and InvalidArgument", response, err)
	}
}

func TestRoutineClient_GetRoutineLogs_RoutineNotFound(t *testing.T) {
	expectedErr := status.Error(codes.NotFound, "routine not found")
	client := NewRoutineClient(&routineServiceClientMock{getRoutineLogs: func(context.Context, *pb.GetRoutineLogsRequest, ...grpc.CallOption) (*pb.GetRoutineLogsResponse, error) {
		return nil, expectedErr
	}})
	response, err := client.GetRoutineLogs(context.Background(), "missing-id", routineTestTime.Add(-time.Hour), routineTestTime)
	if response != nil || !errors.Is(err, expectedErr) {
		t.Fatalf("GetRoutineLogs() = (%#v, %v), want nil and NotFound", response, err)
	}
}

func TestRoutineClient_GetRoutineLogs_DatabaseError(t *testing.T) {
	expectedErr := status.Error(codes.Internal, "database failure")
	client := NewRoutineClient(&routineServiceClientMock{getRoutineLogs: func(context.Context, *pb.GetRoutineLogsRequest, ...grpc.CallOption) (*pb.GetRoutineLogsResponse, error) {
		return nil, expectedErr
	}})
	response, err := client.GetRoutineLogs(context.Background(), "routine-id", routineTestTime.Add(-time.Hour), routineTestTime)
	if response != nil || !errors.Is(err, expectedErr) {
		t.Fatalf("GetRoutineLogs() = (%#v, %v), want nil and Internal", response, err)
	}
}

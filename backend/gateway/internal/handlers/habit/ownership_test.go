package habit

import (
	"context"
	"net/http/httptest"
	"strings"
	"testing"

	"gateway/internal/clients"
	pbHabit "shared/pb/habit"

	"github.com/go-chi/chi/v5"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type habitOwnerClientMock struct {
	pbHabit.HabitServiceClient
	ownerID string
	habits  []*pbHabit.Habit
}

func (m *habitOwnerClientMock) GetHabitByID(_ context.Context, req *pbHabit.GetHabitByIDRequest, _ ...grpc.CallOption) (*pbHabit.GetHabitByIDResponse, error) {
	return &pbHabit.GetHabitByIDResponse{Habit: &pbHabit.Habit{Id: req.HabitId, UserId: m.ownerID}}, nil
}

func (m *habitOwnerClientMock) ListHabitsByRoutine(_ context.Context, req *pbHabit.ListHabitsByRoutineRequest, _ ...grpc.CallOption) (*pbHabit.ListHabitsByRoutineResponse, error) {
	return &pbHabit.ListHabitsByRoutineResponse{Habits: m.habits}, nil
}

type routineOwnerClientMock struct {
	pbHabit.RoutineServiceClient
	ownerID string
}

func (m *routineOwnerClientMock) GetRoutineByID(_ context.Context, req *pbHabit.GetRoutineByIDRequest, _ ...grpc.CallOption) (*pbHabit.GetRoutineByIDResponse, error) {
	return &pbHabit.GetRoutineByIDResponse{Routine: &pbHabit.Routine{Id: req.RoutineId, UserId: m.ownerID}}, nil
}

func TestAuthorizeHabit(t *testing.T) {
	for _, test := range []struct {
		name      string
		ownerID   string
		wantError codes.Code
	}{
		{name: "owner can access", ownerID: "user-1"},
		{name: "other user is denied", ownerID: "user-2", wantError: codes.PermissionDenied},
	} {
		t.Run(test.name, func(t *testing.T) {
			client := clients.NewHabitClient(&habitOwnerClientMock{ownerID: test.ownerID})
			_, err := authorizeHabit(context.Background(), client, "user-1", "habit-1")
			if status.Code(err) != test.wantError {
				t.Fatalf("authorizeHabit() error code = %s, want %s", status.Code(err), test.wantError)
			}
		})
	}
}

func TestAuthorizeRoutine(t *testing.T) {
	for _, test := range []struct {
		name      string
		ownerID   string
		wantError codes.Code
	}{
		{name: "owner can access", ownerID: "user-1"},
		{name: "other user is denied", ownerID: "user-2", wantError: codes.PermissionDenied},
	} {
		t.Run(test.name, func(t *testing.T) {
			client := clients.NewRoutineClient(&routineOwnerClientMock{ownerID: test.ownerID})
			_, err := authorizeRoutine(context.Background(), client, "user-1", "routine-1")
			if status.Code(err) != test.wantError {
				t.Fatalf("authorizeRoutine() error code = %s, want %s", status.Code(err), test.wantError)
			}
		})
	}
}

func TestGetSharedHabitByIDAllowsOtherOwner(t *testing.T) {
	handler := NewHabitHandler(
		clients.NewHabitClient(&habitOwnerClientMock{ownerID: "owner-2"}),
		clients.NewRoutineClient(&routineOwnerClientMock{ownerID: "owner-1"}),
	)
	request := httptest.NewRequest("GET", "/habit/shared/habit-1", nil)
	routeContext := chi.NewRouteContext()
	routeContext.URLParams.Add("id", "habit-1")
	request = request.WithContext(context.WithValue(request.Context(), chi.RouteCtxKey, routeContext))
	response := httptest.NewRecorder()

	handler.GetSharedHabitByID(response, request)

	if response.Code != 200 || !strings.Contains(response.Body.String(), `"user_id":"owner-2"`) {
		t.Fatalf("GetSharedHabitByID() = (%d, %s), want 200 and the habit response", response.Code, response.Body.String())
	}
}

func TestListSharedHabitsByRoutineReturnsOnlyRoutineOwnersHabits(t *testing.T) {
	handler := NewHabitHandler(
		clients.NewHabitClient(&habitOwnerClientMock{habits: []*pbHabit.Habit{
			{Id: "owner-habit", UserId: "routine-owner"},
			{Id: "foreign-habit", UserId: "another-owner"},
		}}),
		clients.NewRoutineClient(&routineOwnerClientMock{ownerID: "routine-owner"}),
	)
	request := httptest.NewRequest("GET", "/habit/shared/routine/routine-1", nil)
	routeContext := chi.NewRouteContext()
	routeContext.URLParams.Add("id", "routine-1")
	request = request.WithContext(context.WithValue(request.Context(), chi.RouteCtxKey, routeContext))
	response := httptest.NewRecorder()

	handler.ListSharedHabitsByRoutine(response, request)

	if response.Code != 200 || !strings.Contains(response.Body.String(), "owner-habit") || strings.Contains(response.Body.String(), "foreign-habit") {
		t.Fatalf("ListSharedHabitsByRoutine() = (%d, %s), want only the routine owner's habit", response.Code, response.Body.String())
	}
}

func TestGetSharedRoutineByIDAllowsOtherOwner(t *testing.T) {
	handler := NewRoutineHandler(
		clients.NewRoutineClient(&routineOwnerClientMock{ownerID: "owner-2"}),
		clients.NewHabitClient(&habitOwnerClientMock{}),
	)
	request := httptest.NewRequest("GET", "/routine/shared/routine-1", nil)
	routeContext := chi.NewRouteContext()
	routeContext.URLParams.Add("id", "routine-1")
	request = request.WithContext(context.WithValue(request.Context(), chi.RouteCtxKey, routeContext))
	response := httptest.NewRecorder()

	handler.GetSharedRoutineByID(response, request)

	if response.Code != 200 || !strings.Contains(response.Body.String(), `"user_id":"owner-2"`) {
		t.Fatalf("GetSharedRoutineByID() = (%d, %s), want 200 and the routine response", response.Code, response.Body.String())
	}
}

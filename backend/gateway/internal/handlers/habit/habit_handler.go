package habit

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"time"

	"gateway/internal/clients"
	"gateway/internal/dto"
	"gateway/internal/helper"
	"gateway/internal/middlewares"

	"github.com/go-chi/chi/v5"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type HabitHandler struct {
	client        *clients.HabitClient
	routineClient *clients.RoutineClient
}

func NewHabitHandler(client *clients.HabitClient, routineClient *clients.RoutineClient) *HabitHandler {
	return &HabitHandler{client: client, routineClient: routineClient}
}

func authorizeHabit(ctx context.Context, client *clients.HabitClient, userID, habitID string) (*dto.HabitResponse, error) {
	if client == nil {
		return nil, status.Error(codes.Unavailable, "habit service client is unavailable")
	}
	response, err := client.GetHabitByID(ctx, habitID)
	if err != nil {
		return nil, err
	}
	if response.Habit.UserID != userID {
		return nil, status.Error(codes.PermissionDenied, "habit does not belong to the authenticated user")
	}
	return response, nil
}

func authorizeRoutine(ctx context.Context, client *clients.RoutineClient, userID, routineID string) (*dto.RoutineResponse, error) {
	if client == nil {
		return nil, status.Error(codes.Unavailable, "routine service client is unavailable")
	}
	response, err := client.GetRoutineByID(ctx, routineID)
	if err != nil {
		return nil, err
	}
	if response.Routine.UserID != userID {
		return nil, status.Error(codes.PermissionDenied, "routine does not belong to the authenticated user")
	}
	return response, nil
}

func writeHabitError(w http.ResponseWriter, start time.Time, info map[string]string, err error) {
	statusCode := http.StatusInternalServerError
	switch status.Code(err) {
	case codes.InvalidArgument:
		statusCode = http.StatusBadRequest
	case codes.NotFound:
		statusCode = http.StatusNotFound
	case codes.Unauthenticated:
		statusCode = http.StatusUnauthorized
	case codes.PermissionDenied:
		statusCode = http.StatusForbidden
	}
	helper.Response(helper.Response_struct{Error: err.Error()}, w, statusCode, http.StatusText(statusCode), info, time.Since(start).Milliseconds())
}

func decodeHabitJSON(w http.ResponseWriter, r *http.Request, target any) (int, string) {
	r.Body = http.MaxBytesReader(w, r.Body, 4096)
	defer r.Body.Close()
	if err := json.NewDecoder(r.Body).Decode(target); err != nil {
		var maxErr *http.MaxBytesError
		if errors.As(err, &maxErr) {
			return http.StatusRequestEntityTooLarge, "Body too large"
		}
		if errors.Is(err, io.EOF) {
			return http.StatusBadRequest, "Body is empty"
		}
		return http.StatusBadRequest, "Invalid request body"
	}
	return 0, ""
}

func habitBodyError(w http.ResponseWriter, start time.Time, info map[string]string, code int, message string) {
	helper.Response(helper.Response_struct{Error: message}, w, code, message, info, time.Since(start).Milliseconds())
}

func (h *HabitHandler) CreateHabit(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	info := map[string]string{"method": r.Method, "url": r.URL.String()}
	var request dto.CreateHabitRequest
	if code, message := decodeHabitJSON(w, r, &request); code != 0 {
		habitBodyError(w, start, info, code, message)
		return
	}
	response, err := h.client.CreateHabit(r.Context(), middlewares.UserIDFromContext(r.Context()), &request)
	if err != nil {
		writeHabitError(w, start, info, err)
		return
	}
	helper.Response(helper.Response_struct{Data: response}, w, http.StatusCreated, "habit created", info, time.Since(start).Milliseconds())
}

func (h *HabitHandler) GetHabitByID(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	info := map[string]string{"method": r.Method, "url": r.URL.String()}
	response, err := authorizeHabit(r.Context(), h.client, middlewares.UserIDFromContext(r.Context()), chi.URLParam(r, "id"))
	if err != nil {
		writeHabitError(w, start, info, err)
		return
	}
	helper.Response(helper.Response_struct{Data: response}, w, http.StatusOK, "habit found", info, time.Since(start).Milliseconds())
}

func (h *HabitHandler) GetSharedHabitByID(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	info := map[string]string{"method": r.Method, "url": r.URL.String()}
	response, err := h.client.GetHabitByID(r.Context(), chi.URLParam(r, "id"))
	if err != nil {
		writeHabitError(w, start, info, err)
		return
	}
	helper.Response(helper.Response_struct{Data: response}, w, http.StatusOK, "habit found", info, time.Since(start).Milliseconds())
}

func (h *HabitHandler) ListHabitsByUser(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	info := map[string]string{"method": r.Method, "url": r.URL.String()}
	response, err := h.client.ListHabitsByUser(r.Context(), middlewares.UserIDFromContext(r.Context()))
	if err != nil {
		writeHabitError(w, start, info, err)
		return
	}
	helper.Response(helper.Response_struct{Data: response}, w, http.StatusOK, "habits found", info, time.Since(start).Milliseconds())
}

func (h *HabitHandler) ListSharedHabitsByRoutine(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	info := map[string]string{"method": r.Method, "url": r.URL.String()}
	routine, err := h.routineClient.GetRoutineByID(r.Context(), chi.URLParam(r, "id"))
	if err != nil {
		writeHabitError(w, start, info, err)
		return
	}
	response, err := h.client.ListHabitsByRoutine(r.Context(), chi.URLParam(r, "id"))
	if err != nil {
		writeHabitError(w, start, info, err)
		return
	}
	habits := make([]*dto.Habit, 0, len(response.Habits))
	for _, habit := range response.Habits {
		if habit != nil && habit.UserID == routine.Routine.UserID {
			habits = append(habits, habit)
		}
	}
	response.Habits = habits
	helper.Response(helper.Response_struct{Data: response}, w, http.StatusOK, "habits found", info, time.Since(start).Milliseconds())
}

func (h *HabitHandler) ListHabitsByRoutine(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	info := map[string]string{"method": r.Method, "url": r.URL.String()}
	routineID := chi.URLParam(r, "id")
	if _, err := authorizeRoutine(r.Context(), h.routineClient, middlewares.UserIDFromContext(r.Context()), routineID); err != nil {
		writeHabitError(w, start, info, err)
		return
	}
	response, err := h.client.ListHabitsByRoutine(r.Context(), routineID)
	if err != nil {
		writeHabitError(w, start, info, err)
		return
	}
	userID := middlewares.UserIDFromContext(r.Context())
	for _, habit := range response.Habits {
		if habit != nil && habit.UserID != userID {
			writeHabitError(w, start, info, status.Error(codes.PermissionDenied, "routine contains a habit owned by another user"))
			return
		}
	}
	helper.Response(helper.Response_struct{Data: response}, w, http.StatusOK, "habits found", info, time.Since(start).Milliseconds())
}

func (h *HabitHandler) EditHabit(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	info := map[string]string{"method": r.Method, "url": r.URL.String()}
	habitID := chi.URLParam(r, "id")
	if _, err := authorizeHabit(r.Context(), h.client, middlewares.UserIDFromContext(r.Context()), habitID); err != nil {
		writeHabitError(w, start, info, err)
		return
	}
	var request dto.EditHabitRequest
	if code, message := decodeHabitJSON(w, r, &request); code != 0 {
		habitBodyError(w, start, info, code, message)
		return
	}
	response, err := h.client.EditHabit(r.Context(), habitID, &request)
	if err != nil {
		writeHabitError(w, start, info, err)
		return
	}
	helper.Response(helper.Response_struct{Data: response}, w, http.StatusOK, "habit updated", info, time.Since(start).Milliseconds())
}

func (h *HabitHandler) DeleteHabit(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	info := map[string]string{"method": r.Method, "url": r.URL.String()}
	habitID := chi.URLParam(r, "id")
	if _, err := authorizeHabit(r.Context(), h.client, middlewares.UserIDFromContext(r.Context()), habitID); err != nil {
		writeHabitError(w, start, info, err)
		return
	}
	response, err := h.client.DeleteHabit(r.Context(), habitID)
	if err != nil {
		writeHabitError(w, start, info, err)
		return
	}
	helper.Response(helper.Response_struct{Data: response}, w, http.StatusOK, "habit deleted", info, time.Since(start).Milliseconds())
}

func (h *HabitHandler) MarkHabitCompleted(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	info := map[string]string{"method": r.Method, "url": r.URL.String()}
	habitID := chi.URLParam(r, "id")
	if _, err := authorizeHabit(r.Context(), h.client, middlewares.UserIDFromContext(r.Context()), habitID); err != nil {
		writeHabitError(w, start, info, err)
		return
	}
	var request dto.HabitCompletionRequest
	if code, message := decodeHabitJSON(w, r, &request); code != 0 {
		habitBodyError(w, start, info, code, message)
		return
	}
	if request.CompletedAt.IsZero() {
		habitBodyError(w, start, info, http.StatusBadRequest, "completed_at is required")
		return
	}
	response, err := h.client.MarkHabitCompleted(r.Context(), habitID, request.CompletedAt)
	if err != nil {
		writeHabitError(w, start, info, err)
		return
	}
	helper.Response(helper.Response_struct{Data: response}, w, http.StatusOK, "habit marked completed", info, time.Since(start).Milliseconds())
}

func (h *HabitHandler) UnmarkHabitCompleted(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	info := map[string]string{"method": r.Method, "url": r.URL.String()}
	habitID := chi.URLParam(r, "id")
	if _, err := authorizeHabit(r.Context(), h.client, middlewares.UserIDFromContext(r.Context()), habitID); err != nil {
		writeHabitError(w, start, info, err)
		return
	}
	var request dto.HabitCompletionRequest
	if code, message := decodeHabitJSON(w, r, &request); code != 0 {
		habitBodyError(w, start, info, code, message)
		return
	}
	if request.CompletedAt.IsZero() {
		habitBodyError(w, start, info, http.StatusBadRequest, "completed_at is required")
		return
	}
	response, err := h.client.UnmarkHabitCompleted(r.Context(), habitID, request.CompletedAt)
	if err != nil {
		writeHabitError(w, start, info, err)
		return
	}
	helper.Response(helper.Response_struct{Data: response}, w, http.StatusOK, "habit completion removed", info, time.Since(start).Milliseconds())
}

func (h *HabitHandler) GetHabitLogs(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	info := map[string]string{"method": r.Method, "url": r.URL.String()}
	habitID := chi.URLParam(r, "id")
	if _, err := authorizeHabit(r.Context(), h.client, middlewares.UserIDFromContext(r.Context()), habitID); err != nil {
		writeHabitError(w, start, info, err)
		return
	}
	startDate, err := time.Parse(time.RFC3339Nano, r.URL.Query().Get("start_date"))
	if err != nil {
		habitBodyError(w, start, info, http.StatusBadRequest, "valid start_date is required")
		return
	}
	endDate, err := time.Parse(time.RFC3339Nano, r.URL.Query().Get("end_date"))
	if err != nil {
		habitBodyError(w, start, info, http.StatusBadRequest, "valid end_date is required")
		return
	}
	response, err := h.client.GetHabitLogs(r.Context(), habitID, startDate, endDate)
	if err != nil {
		writeHabitError(w, start, info, err)
		return
	}
	helper.Response(helper.Response_struct{Data: response}, w, http.StatusOK, "habit logs found", info, time.Since(start).Milliseconds())
}

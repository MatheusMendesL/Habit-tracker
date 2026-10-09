package stats

import (
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

type StatsHandler struct {
	client *clients.StatsClient
}

func NewStatsHandler(client *clients.StatsClient) *StatsHandler {
	return &StatsHandler{client: client}
}

func writeStatsError(w http.ResponseWriter, start time.Time, info map[string]string, err error) {
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

func (h *StatsHandler) CreateUserStats(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	info := statsHTTPInfo(r)
	response, err := h.client.CreateUserStats(r.Context(), middlewares.UserIDFromContext(r.Context()))
	if err != nil {
		writeStatsError(w, start, info, err)
		return
	}
	helper.Response(helper.Response_struct{Data: response}, w, http.StatusCreated, "stats created", info, time.Since(start).Milliseconds())
}

func (h *StatsHandler) GetUserStats(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	info := statsHTTPInfo(r)
	response, err := h.client.GetUserStats(r.Context(), middlewares.UserIDFromContext(r.Context()))
	if err != nil {
		writeStatsError(w, start, info, err)
		return
	}
	helper.Response(helper.Response_struct{Data: response}, w, http.StatusOK, "stats found", info, time.Since(start).Milliseconds())
}

func (h *StatsHandler) DeleteUserStats(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	info := statsHTTPInfo(r)
	response, err := h.client.DeleteUserStats(r.Context(), middlewares.UserIDFromContext(r.Context()))
	if err != nil {
		writeStatsError(w, start, info, err)
		return
	}
	helper.Response(helper.Response_struct{Data: response}, w, http.StatusOK, "stats deleted", info, time.Since(start).Milliseconds())
}

func (h *StatsHandler) RegisterHabitCompletion(w http.ResponseWriter, r *http.Request) {
	h.handleCompletion(w, r, "habit", false)
}

func (h *StatsHandler) UndoHabitCompletion(w http.ResponseWriter, r *http.Request) {
	h.handleCompletion(w, r, "habit", true)
}

func (h *StatsHandler) RegisterRoutineCompletion(w http.ResponseWriter, r *http.Request) {
	h.handleCompletion(w, r, "routine", false)
}

func (h *StatsHandler) UndoRoutineCompletion(w http.ResponseWriter, r *http.Request) {
	h.handleCompletion(w, r, "routine", true)
}

func (h *StatsHandler) handleCompletion(w http.ResponseWriter, r *http.Request, itemType string, undo bool) {
	start := time.Now()
	info := statsHTTPInfo(r)
	r.Body = http.MaxBytesReader(w, r.Body, 4096)
	defer r.Body.Close()

	var request dto.StatsCompletionRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		message := "Invalid request body"
		code := http.StatusBadRequest
		var maxErr *http.MaxBytesError
		switch {
		case errors.As(err, &maxErr):
			message, code = "Body too large", http.StatusRequestEntityTooLarge
		case errors.Is(err, io.EOF):
			message = "Body is empty"
		}
		helper.Response(helper.Response_struct{Error: message}, w, code, message, info, time.Since(start).Milliseconds())
		return
	}
	if request.CompletedAt.IsZero() {
		helper.Response(helper.Response_struct{Error: "completed_at is required"}, w, http.StatusBadRequest, "invalid request body", info, time.Since(start).Milliseconds())
		return
	}

	userID, itemID := middlewares.UserIDFromContext(r.Context()), chi.URLParam(r, "id")
	var response *dto.StatsActionResponse
	var err error
	switch {
	case itemType == "habit" && !undo:
		response, err = h.client.RegisterHabitCompletion(r.Context(), userID, itemID, request.CompletedAt)
	case itemType == "habit" && undo:
		response, err = h.client.UndoHabitCompletion(r.Context(), userID, itemID, request.CompletedAt)
	case itemType == "routine" && !undo:
		response, err = h.client.RegisterRoutineCompletion(r.Context(), userID, itemID, request.CompletedAt)
	default:
		response, err = h.client.UndoRoutineCompletion(r.Context(), userID, itemID, request.CompletedAt)
	}
	if err != nil {
		writeStatsError(w, start, info, err)
		return
	}
	message := "completion registered"
	if undo {
		message = "completion undone"
	}
	helper.Response(helper.Response_struct{Data: response}, w, http.StatusOK, message, info, time.Since(start).Milliseconds())
}

func statsHTTPInfo(r *http.Request) map[string]string {
	return map[string]string{"method": r.Method, "url": r.URL.String()}
}

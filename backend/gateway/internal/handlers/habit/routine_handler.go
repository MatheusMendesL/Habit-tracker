package habit

import (
	"net/http"
	"time"

	"gateway/internal/clients"
	"gateway/internal/dto"
	"gateway/internal/helper"
	"gateway/internal/middlewares"

	"github.com/go-chi/chi/v5"
)

type RoutineHandler struct {
	client *clients.RoutineClient
}

func NewRoutineHandler(client *clients.RoutineClient) *RoutineHandler {
	return &RoutineHandler{client: client}
}

func (h *RoutineHandler) CreateRoutine(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	info := map[string]string{"method": r.Method, "url": r.URL.String()}
	var request dto.CreateRoutineRequest
	if code, message := decodeHabitJSON(w, r, &request); code != 0 {
		habitBodyError(w, start, info, code, message)
		return
	}
	response, err := h.client.CreateRoutine(r.Context(), middlewares.UserIDFromContext(r.Context()), &request)
	if err != nil {
		writeHabitError(w, start, info, err)
		return
	}
	helper.Response(helper.Response_struct{Data: response}, w, http.StatusCreated, "routine created", info, time.Since(start).Milliseconds())
}

func (h *RoutineHandler) GetRoutineByID(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	info := map[string]string{"method": r.Method, "url": r.URL.String()}
	response, err := h.client.GetRoutineByID(r.Context(), chi.URLParam(r, "id"))
	if err != nil {
		writeHabitError(w, start, info, err)
		return
	}
	helper.Response(helper.Response_struct{Data: response}, w, http.StatusOK, "routine found", info, time.Since(start).Milliseconds())
}

func (h *RoutineHandler) EditRoutine(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	info := map[string]string{"method": r.Method, "url": r.URL.String()}
	var request dto.EditRoutineRequest
	if code, message := decodeHabitJSON(w, r, &request); code != 0 {
		habitBodyError(w, start, info, code, message)
		return
	}
	response, err := h.client.EditRoutine(r.Context(), chi.URLParam(r, "id"), &request)
	if err != nil {
		writeHabitError(w, start, info, err)
		return
	}
	helper.Response(helper.Response_struct{Data: response}, w, http.StatusOK, "routine updated", info, time.Since(start).Milliseconds())
}

func (h *RoutineHandler) DeleteRoutine(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	info := map[string]string{"method": r.Method, "url": r.URL.String()}
	response, err := h.client.DeleteRoutine(r.Context(), chi.URLParam(r, "id"))
	if err != nil {
		writeHabitError(w, start, info, err)
		return
	}
	helper.Response(helper.Response_struct{Data: response}, w, http.StatusOK, "routine deleted", info, time.Since(start).Milliseconds())
}

func (h *RoutineHandler) ListRoutinesByUser(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	info := map[string]string{"method": r.Method, "url": r.URL.String()}
	response, err := h.client.ListRoutinesByUser(r.Context(), middlewares.UserIDFromContext(r.Context()))
	if err != nil {
		writeHabitError(w, start, info, err)
		return
	}
	helper.Response(helper.Response_struct{Data: response}, w, http.StatusOK, "routines found", info, time.Since(start).Milliseconds())
}

func (h *RoutineHandler) AddHabitToRoutine(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	info := map[string]string{"method": r.Method, "url": r.URL.String()}
	var request dto.RoutineHabitRequest
	if code, message := decodeHabitJSON(w, r, &request); code != 0 {
		habitBodyError(w, start, info, code, message)
		return
	}
	response, err := h.client.AddHabitToRoutine(r.Context(), chi.URLParam(r, "id"), &request)
	if err != nil {
		writeHabitError(w, start, info, err)
		return
	}
	helper.Response(helper.Response_struct{Data: response}, w, http.StatusOK, "habit added to routine", info, time.Since(start).Milliseconds())
}

func (h *RoutineHandler) RemoveHabitFromRoutine(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	info := map[string]string{"method": r.Method, "url": r.URL.String()}
	request := &dto.RoutineHabitRequest{HabitID: chi.URLParam(r, "habitID")}
	response, err := h.client.RemoveHabitFromRoutine(r.Context(), chi.URLParam(r, "id"), request)
	if err != nil {
		writeHabitError(w, start, info, err)
		return
	}
	helper.Response(helper.Response_struct{Data: response}, w, http.StatusOK, "habit removed from routine", info, time.Since(start).Milliseconds())
}

func (h *RoutineHandler) MarkRoutineCompleted(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	info := map[string]string{"method": r.Method, "url": r.URL.String()}
	var request dto.RoutineCompletionRequest
	if code, message := decodeHabitJSON(w, r, &request); code != 0 {
		habitBodyError(w, start, info, code, message)
		return
	}
	if request.CompletedAt.IsZero() {
		habitBodyError(w, start, info, http.StatusBadRequest, "completed_at is required")
		return
	}
	response, err := h.client.MarkRoutineCompleted(r.Context(), chi.URLParam(r, "id"), request.CompletedAt)
	if err != nil {
		writeHabitError(w, start, info, err)
		return
	}
	helper.Response(helper.Response_struct{Data: response}, w, http.StatusOK, "routine marked completed", info, time.Since(start).Milliseconds())
}

func (h *RoutineHandler) UnmarkRoutineCompleted(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	info := map[string]string{"method": r.Method, "url": r.URL.String()}
	var request dto.RoutineCompletionRequest
	if code, message := decodeHabitJSON(w, r, &request); code != 0 {
		habitBodyError(w, start, info, code, message)
		return
	}
	if request.CompletedAt.IsZero() {
		habitBodyError(w, start, info, http.StatusBadRequest, "completed_at is required")
		return
	}
	response, err := h.client.UnmarkRoutineCompleted(r.Context(), chi.URLParam(r, "id"), request.CompletedAt)
	if err != nil {
		writeHabitError(w, start, info, err)
		return
	}
	helper.Response(helper.Response_struct{Data: response}, w, http.StatusOK, "routine completion removed", info, time.Since(start).Milliseconds())
}

func (h *RoutineHandler) GetRoutineLogs(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	info := map[string]string{"method": r.Method, "url": r.URL.String()}
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
	response, err := h.client.GetRoutineLogs(r.Context(), chi.URLParam(r, "id"), startDate, endDate)
	if err != nil {
		writeHabitError(w, start, info, err)
		return
	}
	helper.Response(helper.Response_struct{Data: response}, w, http.StatusOK, "routine logs found", info, time.Since(start).Milliseconds())
}

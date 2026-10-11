package user

import (
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"net/http"
	"time"

	"gateway/internal/clients"
	"gateway/internal/dto"
	"gateway/internal/helper"
	"gateway/internal/middlewares"

	"github.com/go-chi/chi/v5"
)

type UserHandler struct {
	client *clients.UserClient
}

func requireAuthenticatedUser(w http.ResponseWriter, r *http.Request, userID string, start time.Time, httpInfo map[string]string) bool {
	if userID == middlewares.UserIDFromContext(r.Context()) {
		return true
	}
	durationMs := time.Since(start).Milliseconds()
	helper.Response(helper.Response_struct{Error: "user id does not match authenticated user"}, w, http.StatusForbidden, "forbidden", httpInfo, durationMs)
	return false
}

func NewUserHandler(client *clients.UserClient) *UserHandler {
	return &UserHandler{client: client}
}

func (h *UserHandler) GetUserByID(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	userID := chi.URLParam(r, "id")

	httpInfo := map[string]string{
		"method": r.Method,
		"url":    r.URL.String(),
	}

	if r.Method != http.MethodGet {
		slog.Error("Error with the method", "Error", "The method needs to be GET")
		durationMs := time.Since(start).Milliseconds()
		helper.Response(helper.Response_struct{Error: "method not allowed"}, w, http.StatusMethodNotAllowed, "method not allowed", httpInfo, durationMs)
		return
	}

	if userID == "" {
		durationMs := time.Since(start).Milliseconds()
		helper.Response(helper.Response_struct{Error: "missing user id"}, w, http.StatusBadRequest, "user not found", httpInfo, durationMs)
		return
	}

	user, err := h.client.GetUserByID(r.Context(), userID)
	if err != nil {
		durationMs := time.Since(start).Milliseconds()
		helper.Response(helper.Response_struct{Error: err.Error()}, w, http.StatusInternalServerError, "internal server error", httpInfo, durationMs)
		return
	}

	if user == nil {
		durationMs := time.Since(start).Milliseconds()
		helper.Response(helper.Response_struct{Error: "user not found"}, w, http.StatusNotFound, "user not found", httpInfo, durationMs)
		return
	}

	durationMs := time.Since(start).Milliseconds()
	helper.Response(helper.Response_struct{Data: user}, w, http.StatusOK, "user found", httpInfo, durationMs)
}

func (h *UserHandler) SearchUsers(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	httpInfo := map[string]string{
		"method": r.Method,
		"url":    r.URL.String(),
	}

	if r.Method != http.MethodPost {
		slog.Error("Error with the method", "Error", "The method needs to be POST")
		durationMs := time.Since(start).Milliseconds()
		helper.Response(helper.Response_struct{Error: "method not allowed"}, w, http.StatusMethodNotAllowed, "method not allowed", httpInfo, durationMs)
		return
	}

	r.Body = http.MaxBytesReader(w, r.Body, 1000)
	defer r.Body.Close()

	var data dto.SearchUsersRequest
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		var maxErr *http.MaxBytesError
		durationMs := time.Since(start).Milliseconds()
		if errors.As(err, &maxErr) {
			helper.Response(helper.Response_struct{Error: "Body too large"}, w, http.StatusRequestEntityTooLarge, "Body too large", httpInfo, durationMs)
			return
		}
		if errors.Is(err, io.EOF) {
			helper.Response(helper.Response_struct{Error: "Body is empty"}, w, http.StatusBadRequest, "Body is empty", httpInfo, durationMs)
			return
		}
		helper.Response(helper.Response_struct{Error: "Invalid request body"}, w, http.StatusBadRequest, "Invalid request body", httpInfo, durationMs)
		return
	}

	response, err := h.client.SearchUsers(r.Context(), &data)
	if err != nil {
		durationMs := time.Since(start).Milliseconds()
		helper.Response(helper.Response_struct{Error: err.Error()}, w, http.StatusInternalServerError, "Internal Server Error", httpInfo, durationMs)
		return
	}
	if response == nil {
		response = &dto.SearchUsersResponse{Users: []*dto.User{}}
	}

	durationMs := time.Since(start).Milliseconds()
	helper.Response(helper.Response_struct{Data: response}, w, http.StatusOK, "users found", httpInfo, durationMs)
}

func (h *UserHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	httpInfo := map[string]string{
		"method": r.Method,
		"url":    r.URL.String(),
	}

	if r.Method != http.MethodDelete {
		slog.Error("Error with the method", "Error", "The method needs to be DELETE")
		durationMs := time.Since(start).Milliseconds()
		helper.Response(helper.Response_struct{Error: "method not allowed"}, w, http.StatusMethodNotAllowed, "method not allowed", httpInfo, durationMs)
		return
	}

	userID := chi.URLParam(r, "id")
	if userID == "" {
		durationMs := time.Since(start).Milliseconds()
		helper.Response(helper.Response_struct{Error: "missing user id"}, w, http.StatusBadRequest, "user not found", httpInfo, durationMs)
		return
	}
	if !requireAuthenticatedUser(w, r, userID, start, httpInfo) {
		return
	}

	response, err := h.client.DeleteUser(r.Context(), &dto.DeleteUserRequest{ID: userID})
	if err != nil {
		durationMs := time.Since(start).Milliseconds()
		helper.Response(helper.Response_struct{Error: err.Error()}, w, http.StatusInternalServerError, "Internal Server Error", httpInfo, durationMs)
		return
	}

	durationMs := time.Since(start).Milliseconds()
	helper.Response(helper.Response_struct{Data: response}, w, http.StatusOK, "user deleted", httpInfo, durationMs)
}

func (h *UserHandler) GetUsersByIDs(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	httpInfo := map[string]string{
		"method": r.Method,
		"url":    r.URL.String(),
	}

	if r.Method != http.MethodPost {
		durationMs := time.Since(start).Milliseconds()
		helper.Response(helper.Response_struct{Error: "method not allowed"}, w, http.StatusMethodNotAllowed, "method not allowed", httpInfo, durationMs)
		return
	}

	r.Body = http.MaxBytesReader(w, r.Body, 1000)
	defer r.Body.Close()

	var data dto.GetUsersByIDsRequest
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		var maxErr *http.MaxBytesError
		durationMs := time.Since(start).Milliseconds()
		if errors.As(err, &maxErr) {
			helper.Response(helper.Response_struct{Error: "Body too large"}, w, http.StatusRequestEntityTooLarge, "Body too large", httpInfo, durationMs)
			return
		}
		if errors.Is(err, io.EOF) {
			helper.Response(helper.Response_struct{Error: "Body is empty"}, w, http.StatusBadRequest, "Body is empty", httpInfo, durationMs)
			return
		}
		helper.Response(helper.Response_struct{Error: "Invalid request body"}, w, http.StatusBadRequest, "Invalid request body", httpInfo, durationMs)
		return
	}

	response, err := h.client.GetUsersByIDs(r.Context(), &data)
	if err != nil {
		durationMs := time.Since(start).Milliseconds()
		helper.Response(helper.Response_struct{Error: err.Error()}, w, http.StatusInternalServerError, "Internal Server Error", httpInfo, durationMs)
		return
	}
	if response == nil {
		response = &dto.GetUsersByIDsResponse{Users: []*dto.User{}}
	}

	durationMs := time.Since(start).Milliseconds()
	helper.Response(helper.Response_struct{Data: response}, w, http.StatusOK, "users found", httpInfo, durationMs)
}

func (h *UserHandler) EditUser(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	httpInfo := map[string]string{
		"method": r.Method,
		"url":    r.URL.String(),
	}

	if r.Method != http.MethodPut {
		durationMs := time.Since(start).Milliseconds()
		helper.Response(helper.Response_struct{Error: "method not allowed"}, w, http.StatusMethodNotAllowed, "method not allowed", httpInfo, durationMs)
		return
	}

	userID := chi.URLParam(r, "id")
	if userID == "" {
		durationMs := time.Since(start).Milliseconds()
		helper.Response(helper.Response_struct{Error: "missing user id"}, w, http.StatusBadRequest, "missing user id", httpInfo, durationMs)
		return
	}
	if !requireAuthenticatedUser(w, r, userID, start, httpInfo) {
		return
	}

	r.Body = http.MaxBytesReader(w, r.Body, 1000)
	defer r.Body.Close()

	var data dto.EditUserRequest
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		var maxErr *http.MaxBytesError
		durationMs := time.Since(start).Milliseconds()
		if errors.As(err, &maxErr) {
			helper.Response(helper.Response_struct{Error: "Body too large"}, w, http.StatusRequestEntityTooLarge, "Body too large", httpInfo, durationMs)
			return
		}
		if errors.Is(err, io.EOF) {
			helper.Response(helper.Response_struct{Error: "Body is empty"}, w, http.StatusBadRequest, "Body is empty", httpInfo, durationMs)
			return
		}
		helper.Response(helper.Response_struct{Error: "Invalid request body"}, w, http.StatusBadRequest, "Invalid request body", httpInfo, durationMs)
		return
	}

	response, err := h.client.EditUser(r.Context(), userID, &data)
	if err != nil {
		durationMs := time.Since(start).Milliseconds()
		helper.Response(helper.Response_struct{Error: err.Error()}, w, http.StatusInternalServerError, "Internal Server Error", httpInfo, durationMs)
		return
	}
	if response == nil {
		durationMs := time.Since(start).Milliseconds()
		helper.Response(helper.Response_struct{Error: "user not found"}, w, http.StatusNotFound, "user not found", httpInfo, durationMs)
		return
	}

	durationMs := time.Since(start).Milliseconds()
	helper.Response(helper.Response_struct{Data: response}, w, http.StatusOK, "user updated", httpInfo, durationMs)
}

func (h *UserHandler) EditPassword(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	httpInfo := map[string]string{
		"method": r.Method,
		"url":    r.URL.String(),
	}

	if r.Method != http.MethodPut {
		durationMs := time.Since(start).Milliseconds()
		helper.Response(helper.Response_struct{Error: "method not allowed"}, w, http.StatusMethodNotAllowed, "method not allowed", httpInfo, durationMs)
		return
	}

	userID := chi.URLParam(r, "id")
	if userID == "" {
		durationMs := time.Since(start).Milliseconds()
		helper.Response(helper.Response_struct{Error: "missing user id"}, w, http.StatusBadRequest, "missing user id", httpInfo, durationMs)
		return
	}
	if !requireAuthenticatedUser(w, r, userID, start, httpInfo) {
		return
	}

	r.Body = http.MaxBytesReader(w, r.Body, 1000)
	defer r.Body.Close()

	var data dto.EditPasswordRequest
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		var maxErr *http.MaxBytesError
		durationMs := time.Since(start).Milliseconds()
		if errors.As(err, &maxErr) {
			helper.Response(helper.Response_struct{Error: "Body too large"}, w, http.StatusRequestEntityTooLarge, "Body too large", httpInfo, durationMs)
			return
		}
		if errors.Is(err, io.EOF) {
			helper.Response(helper.Response_struct{Error: "Body is empty"}, w, http.StatusBadRequest, "Body is empty", httpInfo, durationMs)
			return
		}
		helper.Response(helper.Response_struct{Error: "Invalid request body"}, w, http.StatusBadRequest, "Invalid request body", httpInfo, durationMs)
		return
	}
	if data.NewPassword == "" {
		durationMs := time.Since(start).Milliseconds()
		helper.Response(helper.Response_struct{Error: "new_password is required"}, w, http.StatusBadRequest, "new_password is required", httpInfo, durationMs)
		return
	}

	response, err := h.client.EditPassword(r.Context(), userID, &data)
	if err != nil {
		durationMs := time.Since(start).Milliseconds()
		helper.Response(helper.Response_struct{Error: err.Error()}, w, http.StatusInternalServerError, "Internal Server Error", httpInfo, durationMs)
		return
	}
	if response == nil {
		durationMs := time.Since(start).Milliseconds()
		helper.Response(helper.Response_struct{Error: "password update failed"}, w, http.StatusInternalServerError, "Internal Server Error", httpInfo, durationMs)
		return
	}

	durationMs := time.Since(start).Milliseconds()
	helper.Response(helper.Response_struct{Data: response}, w, http.StatusOK, "password updated", httpInfo, durationMs)
}

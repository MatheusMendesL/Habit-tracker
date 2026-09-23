package user

import (
	"net/http"
	"time"

	"gateway/internal/clients"
	"gateway/internal/helper"

	"github.com/go-chi/chi/v5"
)

type UserHandler struct {
	client *clients.UserClient
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

package social

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
)

type SocialHandler struct {
	client *clients.SocialClient
}

func NewSocialHandler(client *clients.SocialClient) *SocialHandler {
	return &SocialHandler{client: client}
}

func (h *SocialHandler) StartFollowing(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	httpInfo := map[string]string{"method": r.Method, "url": r.URL.String()}
	if r.Method != http.MethodPost {
		helper.Response(helper.Response_struct{Error: "method not allowed"}, w, http.StatusMethodNotAllowed, "method not allowed", httpInfo, time.Since(start).Milliseconds())
		return
	}

	var data dto.FollowRequest
	if err := decodeFollowRequest(w, r, &data); err != nil {
		respondSocialBodyError(w, start, httpInfo, err)
		return
	}
	data.FollowerID = middlewares.UserIDFromContext(r.Context())
	response, err := h.client.StartFollowing(r.Context(), &data)
	if err != nil {
		helper.Response(helper.Response_struct{Error: err.Error()}, w, http.StatusInternalServerError, "Internal Server Error", httpInfo, time.Since(start).Milliseconds())
		return
	}
	if response == nil {
		helper.Response(helper.Response_struct{Error: "empty response from social service"}, w, http.StatusInternalServerError, "Internal Server Error", httpInfo, time.Since(start).Milliseconds())
		return
	}
	helper.Response(helper.Response_struct{Data: response}, w, http.StatusOK, "following started", httpInfo, time.Since(start).Milliseconds())
}

func (h *SocialHandler) Unfollow(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	httpInfo := map[string]string{"method": r.Method, "url": r.URL.String()}
	if r.Method != http.MethodPost {
		helper.Response(helper.Response_struct{Error: "method not allowed"}, w, http.StatusMethodNotAllowed, "method not allowed", httpInfo, time.Since(start).Milliseconds())
		return
	}

	var data dto.FollowRequest
	if err := decodeFollowRequest(w, r, &data); err != nil {
		respondSocialBodyError(w, start, httpInfo, err)
		return
	}
	data.FollowerID = middlewares.UserIDFromContext(r.Context())
	response, err := h.client.Unfollow(r.Context(), &data)
	if err != nil {
		helper.Response(helper.Response_struct{Error: err.Error()}, w, http.StatusInternalServerError, "Internal Server Error", httpInfo, time.Since(start).Milliseconds())
		return
	}
	if response == nil {
		helper.Response(helper.Response_struct{Error: "empty response from social service"}, w, http.StatusInternalServerError, "Internal Server Error", httpInfo, time.Since(start).Milliseconds())
		return
	}
	helper.Response(helper.Response_struct{Data: response}, w, http.StatusOK, "unfollow completed", httpInfo, time.Since(start).Milliseconds())
}

func (h *SocialHandler) ListFollowers(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	httpInfo := map[string]string{"method": r.Method, "url": r.URL.String()}
	if r.Method != http.MethodGet {
		helper.Response(helper.Response_struct{Error: "method not allowed"}, w, http.StatusMethodNotAllowed, "method not allowed", httpInfo, time.Since(start).Milliseconds())
		return
	}

	userID := chi.URLParam(r, "id")
	if userID == "" {
		helper.Response(helper.Response_struct{Error: "missing user id"}, w, http.StatusBadRequest, "missing user id", httpInfo, time.Since(start).Milliseconds())
		return
	}
	response, err := h.client.ListFollowers(r.Context(), userID)
	if err != nil {
		helper.Response(helper.Response_struct{Error: err.Error()}, w, http.StatusInternalServerError, "Internal Server Error", httpInfo, time.Since(start).Milliseconds())
		return
	}
	if response == nil {
		helper.Response(helper.Response_struct{Error: "empty response from social service"}, w, http.StatusInternalServerError, "Internal Server Error", httpInfo, time.Since(start).Milliseconds())
		return
	}
	helper.Response(helper.Response_struct{Data: response}, w, http.StatusOK, "followers found", httpInfo, time.Since(start).Milliseconds())
}

func (h *SocialHandler) ListFollowing(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	httpInfo := map[string]string{"method": r.Method, "url": r.URL.String()}
	if r.Method != http.MethodGet {
		helper.Response(helper.Response_struct{Error: "method not allowed"}, w, http.StatusMethodNotAllowed, "method not allowed", httpInfo, time.Since(start).Milliseconds())
		return
	}

	userID := chi.URLParam(r, "id")
	if userID == "" {
		helper.Response(helper.Response_struct{Error: "missing user id"}, w, http.StatusBadRequest, "missing user id", httpInfo, time.Since(start).Milliseconds())
		return
	}
	response, err := h.client.ListFollowing(r.Context(), userID)
	if err != nil {
		helper.Response(helper.Response_struct{Error: err.Error()}, w, http.StatusInternalServerError, "Internal Server Error", httpInfo, time.Since(start).Milliseconds())
		return
	}
	if response == nil {
		helper.Response(helper.Response_struct{Error: "empty response from social service"}, w, http.StatusInternalServerError, "Internal Server Error", httpInfo, time.Since(start).Milliseconds())
		return
	}
	helper.Response(helper.Response_struct{Data: response}, w, http.StatusOK, "following found", httpInfo, time.Since(start).Milliseconds())
}

func decodeFollowRequest(w http.ResponseWriter, r *http.Request, data *dto.FollowRequest) error {
	r.Body = http.MaxBytesReader(w, r.Body, 1000)
	defer r.Body.Close()
	if err := json.NewDecoder(r.Body).Decode(data); err != nil {
		var maxErr *http.MaxBytesError
		if errors.As(err, &maxErr) {
			return errSocialBodyTooLarge
		}
		if errors.Is(err, io.EOF) {
			return errSocialBodyEmpty
		}
		return errSocialBodyInvalid
	}
	return nil
}

type socialBodyError string

func (e socialBodyError) Error() string { return string(e) }

const (
	errSocialBodyTooLarge socialBodyError = "Body too large"
	errSocialBodyEmpty    socialBodyError = "Body is empty"
	errSocialBodyInvalid  socialBodyError = "Invalid request body"
)

func respondSocialBodyError(w http.ResponseWriter, start time.Time, httpInfo map[string]string, err error) {
	statusCode := http.StatusBadRequest
	if errors.Is(err, errSocialBodyTooLarge) {
		statusCode = http.StatusRequestEntityTooLarge
	}
	helper.Response(helper.Response_struct{Error: err.Error()}, w, statusCode, err.Error(), httpInfo, time.Since(start).Milliseconds())
}

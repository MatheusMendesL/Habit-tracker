package api

import (
	"go-api/helper"
	"net/http"
	"os"

	"github.com/go-chi/chi"
)

func StatusRoute(r chi.Router) {
	r.Get("/api_status", func(w http.ResponseWriter, r *http.Request) {
		helper.Response(helper.Response_struct{Data: os.Getenv("API_VERSION")}, w, http.StatusOK)
	})
}

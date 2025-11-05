package api

import (
	"fmt"
	"go-api/helper"
	"go-api/response"
	"net/http"
	"os"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

// define the routes to the project and calls the handle to check
func ControlRoutes() *chi.Mux {
	r := chi.NewRouter()

	r.Use(middleware.Recoverer)
	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)

	r.Get("/api_status", func(w http.ResponseWriter, r *http.Request) {
		helper.Response(map[string]string{"api_version": os.Getenv("API_VERSION")}, w)
	})

	r.Route("/user", func(r chi.Router) {
		r.Get("/get_all_users", response.GetAllUsers)

		r.Get("/get_by_id/{id:[0-9]+}", func(w http.ResponseWriter, r *http.Request) {
			id := chi.URLParam(r, "id")
			id_codified, err := helper.Encrypt(id)
			if err != nil {
				panic(err)
			}

			id_decodified, err := helper.Decrypt(id_codified)
			if err != nil {
				panic(err)
			}

			fmt.Fprintln(w, id_codified, " Decodificado: ", id_decodified)
		})
	})

	return r
}

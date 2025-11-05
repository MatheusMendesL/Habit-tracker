package api

import (
	"fmt"
	"go-api/helper"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

// define the routes to the project and calls the handle to check
func ControlRoutes() *chi.Mux {
	r := chi.NewRouter()

	r.Use(middleware.Recoverer)
	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)

	r.Get("/life", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Hello, the api is running")
	})

	r.Route("/user", func(r chi.Router) {
		r.Get("/get_all_users", func(w http.ResponseWriter, r *http.Request) {

		})

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

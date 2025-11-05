package main

import (
	"fmt"
	"go-api/api"
	"net/http"
	"os"

	"github.com/go-chi/chi"
	"github.com/joho/godotenv"
)

// where the server will run
func main() {

	godotenv.Load()
	r := chi.NewMux()

	r.Mount("/", api.ControlRoutes())

	port := os.Getenv("PORT")
	if port == "" {
		port = ":8080"
	}
	fmt.Println("Server running on", port)

	if err := http.ListenAndServe(port, r); err != nil {
		panic(err)
	}
}

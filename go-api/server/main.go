package main

import (
	"go-api/api"
	"log/slog"
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
	l := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(l)
	slog.Info("Server is running!", "version", "1.0.0")

	if err := http.ListenAndServe(port, r); err != nil {
		slog.Error("Deu um erro na inicialização do servidor")
	}
}

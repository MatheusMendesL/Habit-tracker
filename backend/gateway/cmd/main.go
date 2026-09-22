package main

import (
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

func main() {
	logger, err := zap.NewProduction()
	if err != nil {
		panic(err)
	}
	defer logger.Sync()
	if err := runServer(logger); err != nil {
		logger.Error("Failed to execute code",
			zap.Error(err),
		)
		os.Exit(1)
	}

	logger.Info("All systems offline")
}

func runServer(logger *zap.Logger) error {
	r := chi.NewMux()

	logger.Info("Server is running!")

	if err := http.ListenAndServe(":8080", r); err != nil {
		return err
	}
	logger.Info("The API is running!")
	return nil
}

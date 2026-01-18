package main

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/Weit145/REST_API_golang/internal/config"
)

func main() {
	// Init config
	cfg := config.MustLoad()

	fmt.Println(cfg)

	// Setup logger
	log := setupLogger(cfg.Env)

	log.Info("Info logger")

	log.Debug("Debug logger")

	// This is the entry point of the application.
	// You can initialize your application here.
}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case "local":
		log = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case "production":
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	default:
		log = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	}

	return log
}

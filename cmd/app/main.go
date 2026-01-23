package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/Weit145/REST_API_golang/internal/config"
	"github.com/Weit145/REST_API_golang/internal/http-server/handler/order/create"
	"github.com/Weit145/REST_API_golang/internal/http-server/handler/order/delete"
	"github.com/Weit145/REST_API_golang/internal/http-server/handler/order/read"
	"github.com/Weit145/REST_API_golang/internal/http-server/handler/order/update"
	my_middleware "github.com/Weit145/REST_API_golang/internal/http-server/middleware"
	"github.com/Weit145/REST_API_golang/internal/lib/logger/sloger"
	"github.com/Weit145/REST_API_golang/internal/storage/sqlite"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

func main() {
	// Init config
	cfg := config.MustLoad()

	fmt.Println(cfg)

	// Setup logger
	log := setupLogger(cfg.Env)

	log.Info("Info logger")

	log.Debug("Debug logger")

	// Initialize storage
	storage, err := sqlite.New(cfg.StoragePath)
	if err != nil {
		log.Error("Failed to initialize storage", sloger.Err(err))
		os.Exit(1)
	}

	router := chi.NewRouter()

	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)
	router.Use(middleware.Logger)

	router.Route("/order", func(r chi.Router) {
		// r.Use(middleware.BasicAuth("REST_API_golang", map[string]string{
		// 	"Weit": "123",
		// }))
		r.With(my_middleware.AuthMiddleware(log)).Post("/", create.New(log, storage))
		r.Get("/{order_name}", read.New(log, storage))
		r.With(my_middleware.AuthMiddleware(log)).Delete("/", delete.New(log, storage))
		r.With(my_middleware.AuthMiddleware(log)).Put("/", update.New(log, storage))
	})

	// router.Post("/order", create.New(log, storage))
	// router.Get("/order/{order_name}", read.New(log, storage))
	// router.Delete("/order", delete.New(log, storag
	// router.Put("/order", update.New(log, storage))

	srv := &http.Server{
		Addr:         "0.0.0.0:8080",
		Handler:      router,
		ReadTimeout:  4 * time.Second,
		WriteTimeout: 4 * time.Second,
		IdleTimeout:  30 * time.Second,
	}

	if err := srv.ListenAndServe(); err != nil {
		log.Error("Failed to start HTTP server", sloger.Err(err))
		os.Exit(1)
	}
}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case "local":
		log = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case "prod":
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	default:
		log = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	}

	return log
}

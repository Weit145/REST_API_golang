package middleware

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/Weit145/REST_API_golang/internal/lib/logger/sloger"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
)

type Response struct {
	Status string `json:"status"`
	Error  string `json:"error,omitempty"`
}

func AuthMiddleware(log *slog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			const op = "cmd.middleware.AuthMiddleware"

			log := log.With(
				slog.String("op", op),
				slog.String("request_id", middleware.GetReqID(r.Context())),
			)

			name, pass, ok := r.BasicAuth()
			if !ok || name != "Weit" || pass != "123456" {
				err := errors.New("unauthorized")
				log.Error("authorization failed", sloger.Err(err))

				render.Status(r, http.StatusUnauthorized)
				render.JSON(w, r, Response{
					Status: "error",
					Error:  "Unauthorized",
				})
				return
			}

			log.Info("Client request authorized",
				slog.String("user", name),
			)

			next.ServeHTTP(w, r)
		})
	}
}

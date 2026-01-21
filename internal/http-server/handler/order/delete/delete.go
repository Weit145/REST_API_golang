package delete

import (
	"errors"
	"io"
	"log/slog"
	"net/http"

	"github.com/Weit145/REST_API_golang/internal/lib/logger/sloger"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
)

type Request struct {
	OrderName string `json:"order_name" validate:"required"`
}

type Response struct {
	Status string `json:"status"`
	Error  string `json:"error,omitempty"`
}

//go:generate go run github.com/vektra/mockery/v2@v2.53.5 --name=DeleteOrder
type DeleteOrder interface {
	DeleteOrder(name string) error
}

func New(log *slog.Logger, deleteOrder DeleteOrder) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "http-server.handler.order.delete.New"

		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		var req Request

		err := render.DecodeJSON(r.Body, &req)
		if errors.Is(err, io.EOF) {
			// Такую ошибку встретим, если получили запрос с пустым телом.
			// Обработаем её отдельно
			log.Error("request body is empty")
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, Response{
				Status: "error",
				Error:  "empty request",
			})

			return
		}
		if err != nil {
			log.Error("Failed to decode request", sloger.Err(err))
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, Response{
				Status: "error",
				Error:  "invalid request body",
			})
			return
		}
		if err = validator.New().Struct(req); err != nil {
			log.Error("Validation error", sloger.Err(err))
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, Response{
				Status: "error",
				Error:  "validation error: " + err.Error(),
			})
			return
		}
		err = deleteOrder.DeleteOrder(req.OrderName)
		if err != nil {
			log.Error("Failed to delete order", sloger.Err(err))
			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, Response{
				Status: "error",
				Error:  "failed to delete order",
			})
			return
		}
		render.Status(r, http.StatusOK)
		render.JSON(w, r, Response{
			Status: "success",
		})
	}
}

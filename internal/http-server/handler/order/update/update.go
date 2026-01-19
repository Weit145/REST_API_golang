package update

import (
	"errors"
	"io"
	"log/slog"
	"net/http"

	"github.com/Weit145/REST_API_golang/internal/lib/logger/sloger"
	"github.com/Weit145/REST_API_golang/internal/storage/sqlite"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
)

type Request struct {
	OrderName string  `json:"order_name" validate:"required"`
	Price     float64 `json:"price"`
}
type Response struct {
	Status string `json:"status"`
	Error  string `json:"error,omitempty"`
}

type UpdateOrder interface {
	UpdateOrder(order sqlite.Order) error
}

func New(log *slog.Logger, updateOrder UpdateOrder) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "http-server.handler.order.update.New"

		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		var req Request

		err := render.DecodeJSON(r.Body, &req)
		if errors.Is(err, io.EOF) {
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
		err = updateOrder.UpdateOrder(sqlite.Order{
			Name:  req.OrderName,
			Price: req.Price,
		})
		if err != nil {
			log.Error("Failed to update order", sloger.Err(err))
			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, Response{
				Status: "error",
				Error:  "failed to update order",
			})
			return
		}
		render.JSON(w, r, Response{
			Status: "success",
		})
		log.Info("Updated order",
			slog.String("order_name", req.OrderName),
			slog.Float64("price", req.Price),
		)
	}
}

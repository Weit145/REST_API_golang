package create

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

//go:generate go run github.com/vektra/mockery/v2@v2.53.5 --name=CreateOrder
type CreateOrder interface {
	CreateOrder(order sqlite.Order) error
}

func New(log *slog.Logger, createOrder CreateOrder) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "http-server.handler.order.create.New"

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

		if req.Price <= 0 {
			log.Error("Invalid price value")
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, Response{
				Status: "error",
				Error:  "price must be greater than zero",
			})
			return
		}

		if req.OrderName == "" {
			log.Error("Order name is required")
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, Response{
				Status: "error",
				Error:  "order name is required",
			})
			return
		}

		err = createOrder.CreateOrder(sqlite.Order{
			Name:  req.OrderName,
			Price: req.Price,
		})
		if err != nil {
			log.Error("Failed to create order", sloger.Err(err))
			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, Response{
				Status: "error",
				Error:  "failed to create order",
			})
			return
		}

		render.Status(r, http.StatusCreated)
		render.JSON(w, r, Response{
			Status: "success",
		})

		log.Info("Creating order",
			slog.String("order_name", req.OrderName),
			slog.Float64("price", req.Price),
		)
	}
}

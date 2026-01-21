package read

import (
	"log/slog"
	"net/http"

	"github.com/Weit145/REST_API_golang/internal/lib/logger/sloger"
	"github.com/Weit145/REST_API_golang/internal/storage/sqlite"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
)

type Request struct {
	OrderName string `json:"order_name" validate:"required"`
}

type Response struct {
	Status string       `json:"status"`
	Error  string       `json:"error,omitempty"`
	Order  sqlite.Order `json:"order,omitempty"`
}

//go:generate go run github.com/vektra/mockery/v2@v2.53.5 --name=ReadOrder
type ReadOrder interface {
	ReadOrder(name string) (sqlite.Order, error)
}

func New(log *slog.Logger, readOrder ReadOrder) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "http-server.handler.order.read.New"

		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		var req Request

		req.OrderName = chi.URLParam(r, "order_name")
		if req.OrderName == "" {
			log.Error("Failed to decode request")
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, Response{
				Status: "error",
				Error:  "invalid request body",
			})

			return
		}

		order, err := readOrder.ReadOrder(req.OrderName)
		if err != nil {
			log.Error("failed to read order", sloger.Err(err))
			render.Status(r, http.StatusNotFound)
			render.JSON(w, r, Response{
				Status: "error",
				Error:  "failed to read order",
			})

			return
		}

		render.JSON(w, r, Response{
			Status: "success",
			Order:  order,
		})

		log.Info("Read order",
			slog.String("order_name", req.OrderName),
			slog.Float64("price", order.Price),
		)
	}
}

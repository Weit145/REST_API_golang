package read_test

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Weit145/REST_API_golang/internal/http-server/handler/order/read"
	"github.com/Weit145/REST_API_golang/internal/http-server/handler/order/read/mocks"
	"github.com/go-chi/chi"

	"github.com/Weit145/REST_API_golang/internal/lib/logger/handers/slogdiscard"
	"github.com/Weit145/REST_API_golang/internal/storage/sqlite"
	"github.com/stretchr/testify/require"
)

func TestReadHander(t *testing.T) {
	cases := []struct {
		name      string
		OrderName string
		respError string
		mockError error
	}{
		{
			name:      "success read order",
			OrderName: "Potato",
			respError: "success",
		},
		{
			name:      "success read order",
			OrderName: "Order2",
			respError: "success",
		},
		{
			name:      "error not found",
			OrderName: "Order3",
			respError: "error",
			mockError: errors.New("Not found"),
		},
	}
	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			mockReadOrder := mocks.NewReadOrder(t)

			var check bool = tc.OrderName != ""
			if check {
				mockReadOrder.On("ReadOrder", tc.OrderName).Return(sqlite.Order{
					ID:    1,
					Name:  tc.OrderName,
					Price: 100,
				}, tc.mockError).Once()
			}

			r := chi.NewRouter()
			r.Get("/order/{order_name}", read.New(slogdiscard.NewDiscardLogger(), mockReadOrder))

			ts := httptest.NewServer(r)
			defer ts.Close()

			resp, err := http.Get(ts.URL + "/order/" + tc.OrderName)
			require.NoError(t, err)
			defer resp.Body.Close()

			var body read.Response
			err = json.NewDecoder(resp.Body).Decode(&body)
			require.Equal(t, tc.respError, body.Status)
			if check && tc.mockError == nil {
				require.NoError(t, err)
				require.Equal(t, 1, body.Order.ID)
				require.Equal(t, tc.OrderName, body.Order.Name)
				require.Equal(t, float64(100), body.Order.Price)
			}
		})
	}
}

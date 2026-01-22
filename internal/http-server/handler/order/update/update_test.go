package update_test

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/Weit145/REST_API_golang/internal/http-server/handler/order/update"
	"github.com/Weit145/REST_API_golang/internal/http-server/handler/order/update/mocks"
	"github.com/Weit145/REST_API_golang/internal/lib/logger/handers/slogdiscard"
	"github.com/Weit145/REST_API_golang/internal/storage/sqlite"
	"github.com/stretchr/testify/require"
)

func TestUpdateHandler(t *testing.T) {
	cases := []struct {
		name      string
		OrderName string
		Price     float64
		respError string
		mockError error
	}{
		{
			name:      "success update order",
			OrderName: "Order1",
			Price:     123.12,
			respError: "success",
		},
		{
			name:      "missing name",
			OrderName: "",
			Price:     111.11,
			respError: "error",
		},
		{
			name:      "price = 0",
			OrderName: "Order3",
			Price:     0,
			respError: "error",
		},
		{
			name:      "failed bd",
			OrderName: "Order4",
			Price:     222.2,
			respError: "error",
			mockError: errors.New("Field db"),
		},
	}
	for _, tc := range cases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			mockUpdateOrder := mocks.NewUpdateOrder(t)

			if tc.OrderName != "" && tc.Price != 0 {
				mockUpdateOrder.On("UpdateOrder", sqlite.Order{
					Name:  tc.OrderName,
					Price: tc.Price,
				}).Return(tc.mockError).Once()
			}

			handler := update.New(slogdiscard.NewDiscardLogger(), mockUpdateOrder)

			input := fmt.Sprintf(`{"order_name":"%s","price":%f}`, tc.OrderName, tc.Price)

			req, err := http.NewRequest(http.MethodPut, "/order", strings.NewReader(input))
			require.NoError(t, err)

			rr := httptest.NewRecorder()
			handler.ServeHTTP(rr, req)

			body := rr.Body.String()

			var resp update.Response
			require.NoError(t, json.Unmarshal([]byte(body), &resp))
			require.Equal(t, tc.respError, resp.Status)
		})
	}
}

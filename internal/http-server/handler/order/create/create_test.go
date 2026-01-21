package create_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/Weit145/REST_API_golang/internal/http-server/handler/order/create"
	"github.com/Weit145/REST_API_golang/internal/http-server/handler/order/create/mocks"
	"github.com/Weit145/REST_API_golang/internal/lib/logger/handers/slogdiscard"
	"github.com/Weit145/REST_API_golang/internal/storage/sqlite"
)

func TestCreateHandker(t *testing.T) {
	cases := []struct {
		name      string
		OrderName string
		Price     float64
		respError string
		mockError error
	}{
		{
			name:      "missing order name",
			OrderName: "",
			Price:     100,
			respError: "error",
		},
		{
			name:      "missing order price",
			OrderName: "Order1",
			Price:     0,
			respError: "error",
		},
		// {
		// 	name:      "failed to create order",
		// 	OrderName: "Order2",
		// 	Price:     200,
		// 	respError: "fail",
		// 	mockError: errors.New("Field to create"),
		// },
		{
			name:      "success create order",
			OrderName: "Order3",
			Price:     300,
			respError: "success",
		},
		{
			name:      "success create order ",
			OrderName: "Order4",
			Price:     400,
			respError: "success",
		},
	}
	for _, tc := range cases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			mockCreateOrder := mocks.NewCreateOrder(t)

			if tc.OrderName != "" && tc.Price != 0 {
				mockCreateOrder.On("CreateOrder", sqlite.Order{
					Name:  tc.OrderName,
					Price: tc.Price,
				}).Return(tc.mockError).Once()
			}
			handler := create.New(slogdiscard.NewDiscardLogger(), mockCreateOrder)

			input := fmt.Sprintf(`{"order_name":"%s","price":%f}`, tc.OrderName, tc.Price)

			req, err := http.NewRequest(http.MethodPost, "/orders", strings.NewReader(input))
			require.NoError(t, err)

			rr := httptest.NewRecorder()
			handler.ServeHTTP(rr, req)

			// require.Equal(t, rr.Code, http.StatusCreated)

			body := rr.Body.String()

			var resp create.Response
			require.NoError(t, json.Unmarshal([]byte(body), &resp))
			require.Equal(t, tc.respError, resp.Status)
		})
	}
}

package delete_test

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/Weit145/REST_API_golang/internal/http-server/handler/order/create"
	"github.com/Weit145/REST_API_golang/internal/http-server/handler/order/delete"
	"github.com/Weit145/REST_API_golang/internal/http-server/handler/order/delete/mocks"
	"github.com/Weit145/REST_API_golang/internal/lib/logger/handers/slogdiscard"
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
			name:      "success delete oder",
			OrderName: "Order1",
			respError: "success",
		},
		{
			name:      "error missing name",
			OrderName: "",
			respError: "error",
		},
		{
			name:      "error delete bd",
			OrderName: "Order2",
			respError: "error",
			mockError: errors.New("Field to delete"),
		},
	}
	for _, tc := range cases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			mockDeleteOrder := mocks.NewDeleteOrder(t)

			var check bool = tc.OrderName != ""

			if check {
				mockDeleteOrder.On("DeleteOrder", tc.OrderName).Return(tc.mockError).Once()
			}

			handler := delete.New(slogdiscard.NewDiscardLogger(), mockDeleteOrder)

			input := fmt.Sprintf(`{"order_name":"%s"}`, tc.OrderName)

			req, err := http.NewRequest(http.MethodDelete, "/order", strings.NewReader(input))
			require.NoError(t, err)

			rr := httptest.NewRecorder()
			handler.ServeHTTP(rr, req)

			// require.Equal(t, rr.Code, tc.respError)

			body := rr.Body.String()

			var resp create.Response
			require.NoError(t, json.Unmarshal([]byte(body), &resp))
			require.Equal(t, tc.respError, resp.Status)
		})
	}
}

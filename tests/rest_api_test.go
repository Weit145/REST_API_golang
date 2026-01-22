package tests

import (
	"net/url"
	"testing"

	"github.com/Weit145/REST_API_golang/internal/http-server/handler/order/create"
	"github.com/Weit145/REST_API_golang/internal/http-server/handler/order/delete"
	"github.com/gavv/httpexpect/v2"
)

const (
	host = "localhost:8080"
)

func TestUrlRestAPI_CRUD(t *testing.T) {
	u := url.URL{
		Scheme: "http",
		Host:   host,
	}
	e := httpexpect.Default(t, u.String())

	obj := e.POST("/order").
		WithJSON(create.Request{
			OrderName: "Test",
			Price:     123.123,
		}).
		WithBasicAuth("Weit", "123456").
		Expect().
		Status(201).
		JSON().Object()

	obj.Value("status").String().IsEqual("success")

	obj = e.GET("/order/Test").
		Expect().
		Status(200).
		JSON().Object()

	obj.Value("status").String().IsEqual("success")

	orderObj := obj.Value("order").Object()
	orderObj.Value("Name").IsEqual("Test")
	orderObj.Value("Price").IsEqual(123.123)

	obj = e.PUT("/order").
		WithJSON(create.Request{
			OrderName: "Test",
			Price:     456.456,
		}).
		WithBasicAuth("Weit", "123456").
		Expect().
		// Status(201).
		JSON().Object()

	obj.Value("status").String().IsEqual("success")

	obj = e.GET("/order/Test").
		Expect().
		Status(200).
		JSON().Object()

	obj.Value("status").String().IsEqual("success")

	orderObj = obj.Value("order").Object()
	orderObj.Value("Name").IsEqual("Test")
	orderObj.Value("Price").IsEqual(456.456)

	obj = e.DELETE("/order").
		WithJSON(delete.Request{
			OrderName: "Test",
		}).
		WithBasicAuth("Weit", "123456").
		Expect().
		Status(200).
		JSON().Object()
	obj.Value("status").String().IsEqual("success")
}

package echo_swagger

import (
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"net/http/httptest"
	"testing"
)

func performRequest(method, target string, app *echo.Echo) *httptest.ResponseRecorder {
	req := httptest.NewRequest(method, target, nil)
	rec := httptest.NewRecorder()

	app.ServeHTTP(rec,req)
	return rec
}

func TestMiddleware_Register(t *testing.T) {
	t.Run("Endpoint check", func(t *testing.T) {
		app := echo.New()

		middleware := NewMiddleware("./docs/swagger_test.json", "/")

		middleware.Register(app)

		w1 := performRequest("GET", "/docs", app)
		assert.Equal(t, 200, w1.Code)

		w2 := performRequest("GET", "/swagger.json", app)
		assert.Equal(t, 200, w2.Code)

		w3 := performRequest("GET", "/notfound", app)
		assert.Equal(t, 404, w3.Code)
	})

	t.Run("Swagger.json file is not exist", func(t *testing.T) {
		app := echo.New()

		middleware := NewMiddleware("./docs/swagger.json", "/")

		assert.Panics(t, func() {
			middleware.Register(app)
		}, "/swagger.json file is not exist")
	})

	t.Run("Swagger.json missing file", func(t *testing.T) {
		app := echo.New()

		middleware := NewMiddleware("./docs/swagger_missing_test.json", "/")

		assert.Panics(t, func() {
			middleware.Register(app)
		}, "invalid character ':' after object key:value pair")
	})
}
package echo_swagger

import (
	"github.com/gorilla/handlers"
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

func performRequestWithOrigin(method, target, origin string, app *echo.Echo) *httptest.ResponseRecorder {
	req := httptest.NewRequest(method, target, nil)
	if origin != "" {
		req.Header.Set("Origin", origin)
	}
	rec := httptest.NewRecorder()

	app.ServeHTTP(rec, req)
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

	t.Run("Default CORS does not expose wildcard to untrusted origins", func(t *testing.T) {
		app := echo.New()

		middleware := NewMiddleware("./docs/swagger_test.json", "/")
		middleware.Register(app)

		w := performRequestWithOrigin("GET", "/swagger.json", "http://evil-attacker.com", app)
		assert.Equal(t, 200, w.Code)
		assert.Empty(t, w.Header().Get("Access-Control-Allow-Origin"))
	})

	t.Run("CORS with AllowedOrigins", func(t *testing.T) {
		app := echo.New()

		middleware := NewMiddleware("./docs/swagger_test.json", "/", "https://trusted-domain.com")
		middleware.Register(app)

		// Trusted origin
		wTrusted := performRequestWithOrigin("GET", "/swagger.json", "https://trusted-domain.com", app)
		assert.Equal(t, 200, wTrusted.Code)
		assert.Equal(t, "https://trusted-domain.com", wTrusted.Header().Get("Access-Control-Allow-Origin"))

		// Untrusted origin
		wUntrusted := performRequestWithOrigin("GET", "/swagger.json", "http://evil-attacker.com", app)
		assert.Equal(t, 200, wUntrusted.Code)
		assert.Empty(t, wUntrusted.Header().Get("Access-Control-Allow-Origin"))

		// OPTIONS preflight
		req := httptest.NewRequest("OPTIONS", "/swagger.json", nil)
		req.Header.Set("Origin", "https://trusted-domain.com")
		req.Header.Set("Access-Control-Request-Method", "GET")
		rec := httptest.NewRecorder()
		app.ServeHTTP(rec, req)
		assert.Equal(t, 200, rec.Code)
		assert.Equal(t, "https://trusted-domain.com", rec.Header().Get("Access-Control-Allow-Origin"))
	})

	t.Run("CORS with SetAllowedOrigins builder", func(t *testing.T) {
		app := echo.New()

		middleware := NewMiddleware("./docs/swagger_test.json", "/").
			SetAllowedOrigins("https://dashboard.internal")
		middleware.Register(app)

		w := performRequestWithOrigin("GET", "/swagger.json", "https://dashboard.internal", app)
		assert.Equal(t, 200, w.Code)
		assert.Equal(t, "https://dashboard.internal", w.Header().Get("Access-Control-Allow-Origin"))
	})

	t.Run("CORS with SetCORSOptions builder", func(t *testing.T) {
		app := echo.New()

		middleware := NewMiddleware("./docs/swagger_test.json", "/").
			SetCORSOptions(
				handlers.AllowedOrigins([]string{"https://cors.custom.com"}),
				handlers.AllowedMethods([]string{"GET"}),
			)
		middleware.Register(app)

		w := performRequestWithOrigin("GET", "/swagger.json", "https://cors.custom.com", app)
		assert.Equal(t, 200, w.Code)
		assert.Equal(t, "https://cors.custom.com", w.Header().Get("Access-Control-Allow-Origin"))

		wBlocked := performRequestWithOrigin("GET", "/swagger.json", "http://evil-attacker.com", app)
		assert.Equal(t, 200, wBlocked.Code)
		assert.Empty(t, wBlocked.Header().Get("Access-Control-Allow-Origin"))
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
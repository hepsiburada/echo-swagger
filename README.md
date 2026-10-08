# Echo-Swagger Middleware

Echo middleware to automatically generate RESTful API documentation with Swagger 2.0.

[![Travis branch](https://img.shields.io/travis/hepsiburada/echo-swagger/master.svg)](https://travis-ci.org/hepsiburada/echo-swagger)
[![Codecov branch](https://img.shields.io/codecov/c/github/hepsiburada/echo-swagger/master.svg)](https://codecov.io/gh/hepsiburada/echo-swagger)

## Usage

```go
package main

import (
	echo_swagger "echo-swagger"
	"github.com/labstack/echo/v4"
)

func main() {
	app := echo.New()

	// Default: Secure by default (same-origin only, no wildcard CORS header)
	swaggerMiddleware := echo_swagger.NewMiddleware("./docs/swagger.json", "/")

	// Or specify allowed CORS origins:
	// swaggerMiddleware := echo_swagger.NewMiddleware("./docs/swagger.json", "/", "https://trusted-domain.com")
	// Or via builder:
	// swaggerMiddleware.SetAllowedOrigins("https://trusted-domain.com")

	swaggerMiddleware.Register(app)

	app.Start(":8080")
}
```

### CORS Configuration

By default, cross-origin access to `/swagger.json` is restricted (no wildcard `Access-Control-Allow-Origin: *` header) to protect API definitions from unauthorized cross-origin requests.

You can configure CORS using allowed origins or full gorilla/handlers CORS options:

- **Allowed Origins**:
  ```go
  swaggerMiddleware := echo_swagger.NewMiddleware("./docs/swagger.json", "/", "https://example.com")
  // Or:
  swaggerMiddleware.SetAllowedOrigins("https://example.com", "https://admin.example.com")
  ```

- **Advanced CORS Options**:
  ```go
  swaggerMiddleware.SetCORSOptions(
      handlers.AllowedOrigins([]string{"https://example.com"}),
      handlers.AllowedMethods([]string{"GET"}),
      handlers.AllowedHeaders([]string{"Content-Type"}),
  )
  ```
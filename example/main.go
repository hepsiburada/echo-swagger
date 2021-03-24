package main

import (
	echo_swagger "echo-swagger"
	"fmt"
	"github.com/labstack/echo/v4"
)

func main() {
	app := echo.New()

	swaggerMiddleware := echo_swagger.NewMiddleware("./docs/swagger_example.json", "/")

	swaggerMiddleware.Register(app)

	_ = GetHandler()

	SetupRoutes(app)

	_ = app.Start(fmt.Sprintf(":%s", "8080"))
}

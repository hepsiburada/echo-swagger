package echo_swagger

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-openapi/loads"
	"github.com/go-openapi/runtime/middleware"
	"github.com/gorilla/handlers"
	"github.com/labstack/echo/v4"
	"net/http"
	"os"
	"path"
)

type Middleware struct {
	FilePath       string
	BasePath       string
	AllowedOrigins []string
	CORSOptions    []handlers.CORSOption
}

func (s *Middleware) swaggerUIHandler(opts middleware.SwaggerUIOpts) http.Handler {
	return middleware.SwaggerUI(opts, nil)
}

func (s *Middleware) swaggerSpecFileHandler(swaggerUiHandler http.Handler) (http.Handler, error) {
	if _, err := os.Stat(s.FilePath); os.IsNotExist(err) {
		return nil, errors.New(fmt.Sprintf("%s file is not exist", s.FilePath))
	}

	specDoc, err := loads.Spec(s.FilePath)
	if err != nil {
		return nil, err
	}

	b, err := json.MarshalIndent(specDoc.Spec(), "", "  ")
	if err != nil {
		return nil, err
	}

	specHandler := middleware.Spec(s.BasePath, b, swaggerUiHandler)

	if len(s.CORSOptions) > 0 {
		return handlers.CORS(s.CORSOptions...)(specHandler), nil
	}

	if len(s.AllowedOrigins) > 0 {
		return handlers.CORS(handlers.AllowedOrigins(s.AllowedOrigins))(specHandler), nil
	}

	return specHandler, nil
}

func (s *Middleware) SetAllowedOrigins(origins ...string) *Middleware {
	s.AllowedOrigins = origins
	return s
}

func (s *Middleware) SetCORSOptions(opts ...handlers.CORSOption) *Middleware {
	s.CORSOptions = opts
	return s
}

func (s *Middleware) Register(app *echo.Echo) {
	swaggerUIOpts := middleware.SwaggerUIOpts{
		BasePath: s.BasePath,
		SpecURL:  path.Join(s.BasePath, "swagger.json"),
		Path:     "docs",
	}

	swaggerUiHandler := s.swaggerUIHandler(swaggerUIOpts)
	specFileHandler, err := s.swaggerSpecFileHandler(swaggerUiHandler)

	if err != nil {
		panic(err)
	}

	specPath := path.Join(s.BasePath, "swagger.json")
	docsPath := path.Join(s.BasePath, swaggerUIOpts.Path)

	app.GET(docsPath, echo.WrapHandler(swaggerUiHandler))
	app.GET(specPath, echo.WrapHandler(specFileHandler))
	if len(s.CORSOptions) > 0 || len(s.AllowedOrigins) > 0 {
		app.OPTIONS(specPath, echo.WrapHandler(specFileHandler))
	}
}

func NewMiddleware(fileName string, basePath string, allowedOrigins ...string) *Middleware {
	return &Middleware{
		FilePath:       fileName,
		BasePath:       basePath,
		AllowedOrigins: allowedOrigins,
	}
}

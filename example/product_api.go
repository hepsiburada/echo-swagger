package main

import (
	"github.com/labstack/echo/v4"
	"net/http"
	"net/url"
	"sync"
)

var (
	productHandlerOnce sync.Once
	productHandler     *handler
)

type handler struct {
}

// swagger:parameters getProduct
type productRequest struct {
	// Id of an product
	// In: path
	ProductId string `json:"productId"`
}

// swagger:route GET /product/{productId} Product getProduct
// responses:
//  200: ProductModel
func (h *handler) GetByProductId(context echo.Context) error {
	productId := context.Param("productId")
	productId, err := url.QueryUnescape(productId)

	if err != nil {
		return context.JSON(http.StatusBadRequest, err)
	}

	product := NewProductModel("Test Product")

	_ = context.JSON(http.StatusOK, product)

	return nil
}

func SetupRoutes(app *echo.Echo) {
	app.GET("/product/:productId", productHandler.GetByProductId)
}

func GetHandler() *handler {
	productHandlerOnce.Do(func() {
		productHandler = NewHandler()
	})

	return productHandler
}

func NewHandler() *handler {
	return &handler{}
}

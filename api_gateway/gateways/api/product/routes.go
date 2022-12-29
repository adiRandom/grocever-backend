package product

import (
	"api_gateway/services/api/products"
	"github.com/gin-gonic/gin"
	"lib/data/dto/product_processing"
	"lib/helpers"
	"lib/network/http"
)

type Router struct {
	productApiClient *products.Client
}

func NewProductRouter(productApiClient *products.Client) *Router {
	return &Router{productApiClient}
}

func (r *Router) GetRoutes(router *gin.RouterGroup) {
	router.GET("/list", r.getAllUserProducts)
}

func (r *Router) getAllUserProducts(context *gin.Context) {
	productList, err := r.productApiClient.GetProductList()
	if err != nil {
		context.JSON(err.Code, http.Response[helpers.None]{
			StatusCode: err.Code,
			Err:        err.Error(),
			Body:       helpers.None{},
		}.GetH())
		return
	}
	context.JSON(200, http.Response[product_processing.UserProductListDto]{
		StatusCode: 200,
		Body: product_processing.UserProductListDto{
			Products: productList,
		},
	}.GetH())
}

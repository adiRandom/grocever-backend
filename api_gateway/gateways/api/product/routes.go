package product

import (
	"api_gateway/services/api/products"
	"github.com/gin-gonic/gin"
	"lib/api/middleware"
	"lib/data/dto/product"
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
	userId, exists := context.Get(middleware.UserIdKey)
	if !exists {
		context.JSON(401, http.Response[helpers.None]{
			StatusCode: 401,
			Err:        "Unauthorized",
			Body:       helpers.None{},
		})
		return
	}

	productList, apiError := r.productApiClient.GetProductList(userId.(int))
	if apiError != nil {
		context.JSON(apiError.Code, http.Response[helpers.None]{
			StatusCode: apiError.Code,
			Err:        apiError.Error(),
			Body:       helpers.None{},
		}.GetH())
		return
	}
	context.JSON(200, http.Response[product.UserProductListDto]{
		StatusCode: 200,
		Body: product.UserProductListDto{
			Products: productList,
		},
		Err: "",
	}.GetH())
}

package product

import (
	"api_gateway/services/api/products"
	"github.com/gin-gonic/gin"
	"lib/data/dto/product"
	"lib/helpers"
	"lib/network/http"
	"strconv"
)

type Router struct {
	productApiClient *products.Client
}

func NewProductRouter(productApiClient *products.Client) *Router {
	return &Router{productApiClient}
}

func (r *Router) GetRoutes(router *gin.RouterGroup) {
	router.GET("/:userId/list", r.getAllUserProducts)
}

func (r *Router) getAllUserProducts(context *gin.Context) {
	userId := context.Param("userId")
	intUserId, err := strconv.Atoi(userId)
	if err != nil {
		context.JSON(500, http.Response[helpers.None]{
			StatusCode: 500,
			Err:        "Invalid user id",
			Body:       helpers.None{},
		}.GetH())
		return
	}

	productList, apiError := r.productApiClient.GetProductList(intUserId)
	if err != nil {
		context.JSON(apiError.Code, http.Response[helpers.None]{
			StatusCode: apiError.Code,
			Err:        err.Error(),
			Body:       helpers.None{},
		}.GetH())
		return
	}
	context.JSON(200, http.Response[product.UserProductListDto]{
		StatusCode: 200,
		Body: product.UserProductListDto{
			Products: productList,
		},
	}.GetH())
}

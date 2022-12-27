package product

import (
	"github.com/gin-gonic/gin"
	"lib/data/dto/product_processing"
	"lib/data/models/user_product"
	"lib/functional"
	"lib/helpers"
	"lib/network/http"
	"productProcessing/data/database/repositories"
)

type Router struct {
	repository *repositories.UserProductRepository
}

func NewProductRouter(userProductRepo *repositories.UserProductRepository) *Router {
	return &Router{
		repository: userProductRepo,
	}
}

func (r *Router) GetRoutes(router *gin.RouterGroup) {
	router.GET("/list", r.getAllUserProducts)
}

func (r *Router) getAllUserProducts(context *gin.Context) {
	products, err := r.repository.GetAll()
	if err != nil {
		context.JSON(500, http.Response[helpers.None]{
			StatusCode: 500,
			Err:        err.Error(),
			Body:       helpers.None{},
		}.GetH())
		return
	}

	context.JSON(200, http.Response[product_processing.UserProductListDto]{
		StatusCode: 200,
		Body: product_processing.UserProductListDto{
			Products: functional.Map(
				products,
				func(product user_product.Model) product_processing.UserProductDto {
					return product.ToDto()
				},
			),
		},
	}.GetH())
}

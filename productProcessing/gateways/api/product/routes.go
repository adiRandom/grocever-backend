package product

import (
	"github.com/gin-gonic/gin"
	productDtos "lib/data/dto/product"
	productModels "lib/data/models/product"
	"lib/functional"
	"lib/helpers"
	"lib/network/http"
	"productProcessing/data/database/repositories"
)

type Router struct {
	repository *repositories.PurchaseInstalmentRepository
}

func NewProductRouter(userProductRepo *repositories.PurchaseInstalmentRepository) *Router {
	return &Router{
		repository: userProductRepo,
	}
}

func (r *Router) GetRoutes(router *gin.RouterGroup) {
	router.GET("/:userId/list", r.getAllUserProducts)
}

func (r *Router) getAllUserProducts(context *gin.Context) {
	userId := context.Param("userId")
	products, err := r.repository.GetAll()
	if err != nil {
		context.JSON(500, http.Response[helpers.None]{
			StatusCode: 500,
			Err:        err.Error(),
			Body:       helpers.None{},
		}.GetH())
		return
	}

	context.JSON(200, http.Response[productDtos.UserProductListDto]{
		StatusCode: 200,
		Body: productDtos.UserProductListDto{
			Products: functional.Map(
				products,
				func(userOcrProductModel productModels.PurchaseInstalmentModel) productDtos.PurchaseInstalmentDto {
					return userOcrProductModel.ToDto()
				},
			),
		},
	}.GetH())
}

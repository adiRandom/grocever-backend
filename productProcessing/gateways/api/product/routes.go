package product

import (
	"fmt"
	"github.com/gin-gonic/gin"
	productDtos "lib/data/dto/product"
	"lib/functional"
	"lib/helpers"
	"lib/network/http"
	"productProcessing/data/database/repositories"
	"productProcessing/data/models"
	"strconv"
)

type Router struct {
	purchaseInstalmentRepository *repositories.PurchaseInstalmentRepository
	productRepository            *repositories.ProductRepository
	missLinkRepository           *repositories.MissLinkRepository
}

func NewProductRouter(
	userProductRepo *repositories.PurchaseInstalmentRepository,
	productRepo *repositories.ProductRepository,
	missLinkRepo *repositories.MissLinkRepository,
) *Router {
	return &Router{
		purchaseInstalmentRepository: userProductRepo,
		productRepository:            productRepo,
		missLinkRepository:           missLinkRepo,
	}
}

func (r *Router) GetRoutes(router *gin.RouterGroup) {
	router.GET("/:userId/list", r.getAllUserProducts)
	router.POST("/:userId", r.createPurchaseInstalment)
	router.POST("/report", r.reportMissLink)
}

func (r *Router) createPurchaseInstalment(context *gin.Context) {
	var dto productDtos.CreatePurchaseInstalmentNoOcrWithUserDto
	err := context.BindJSON(&dto)
	if err != nil {
		context.JSON(400, http.Response[helpers.None]{
			Err:        err.Error(),
			StatusCode: 400,
			Body:       helpers.None{},
		}.GetH())
		return
	}

	model, err := r.purchaseInstalmentRepository.CreatePurchaseInstalmentNoOcr(dto)
	if err != nil {
		context.JSON(500, http.Response[helpers.None]{
			Err:        err.Error(),
			StatusCode: 500,
			Body:       helpers.None{},
		}.GetH())
		return
	}

	context.JSON(200, http.Response[productDtos.PurchaseInstalmentDto]{
		Err:        "",
		Body:       model.ToDto(),
		StatusCode: 200,
	}.GetH())
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

	products, err := r.purchaseInstalmentRepository.GetUserProducts(intUserId)
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
				func(userProduct models.UserProduct) productDtos.UserProductDto {
					return userProduct.ToDto()
				},
			),
		},
	}.GetH())
}

func (r *Router) reportMissLink(context *gin.Context) {
	var dto productDtos.ReportWithUserIdDto
	err := context.BindJSON(&dto)
	if err != nil {
		context.JSON(400, http.Response[helpers.None]{
			Err:        err.Error(),
			StatusCode: 400,
			Body:       helpers.None{},
		}.GetH())
		return
	}

	isLinkingDenied, err := r.missLinkRepository.IsLinkingDenied(dto.ProductId, dto.OcrProductName)
	if err != nil {
		context.JSON(500, http.Response[helpers.None]{
			Err:        err.Error(),
			StatusCode: 500,
			Body:       helpers.None{},
		}.GetH())
		return
	}

	if isLinkingDenied {
		context.JSON(204, http.Response[helpers.None]{
			Err:        "",
			StatusCode: 204,
			Body:       helpers.None{},
		}.GetH())
		return
	}

	entity, err := r.missLinkRepository.Create(dto.ProductId, dto.OcrProductName, dto.UserId)
	if err != nil {
		context.JSON(500, http.Response[helpers.None]{
			Err:        err.Error(),
			StatusCode: 500,
			Body:       helpers.None{},
		}.GetH())
		return
	}

	shouldBreakLink, err := r.missLinkRepository.ShouldBreakProductLink(dto.ProductId, dto.OcrProductName)
	if err != nil {
		// Revert the creation of the miss link so the user can report it again
		deleteErr := r.missLinkRepository.Delete(*entity)
		if deleteErr != nil {
			context.JSON(500, http.Response[helpers.None]{
				Err:        err.Error(),
				StatusCode: 500,
				Body:       helpers.None{},
			}.GetH())
			return
		}

		context.JSON(500, http.Response[helpers.None]{
			Err:        err.Error(),
			StatusCode: 500,
			Body:       helpers.None{},
		}.GetH())
		return
	}

	if shouldBreakLink {
		go func() {
			err := r.productRepository.BreakProductLink(dto.ProductId, dto.OcrProductName)
			if err != nil {
				// Revert the creation of the miss link so the user can report it again
				deleteErr := r.missLinkRepository.Delete(*entity)
				if deleteErr != nil {
					fmt.Println(err.Error())
				}
			}
		}()
	}

	context.Status(202)
}

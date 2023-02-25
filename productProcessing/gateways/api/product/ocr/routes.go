package ocr

import (
	"github.com/gin-gonic/gin"
	"lib/data/dto/product"
	"lib/data/dto/product/ocr"
	productModel "lib/data/models/product"
	"lib/functional"
	"lib/helpers"
	"lib/network/http"
	"productProcessing/data/database/repositories"
)

type Router struct {
	ocrProductRepository         *repositories.OcrProductRepository
	purchaseInstalmentRepository *repositories.PurchaseInstalmentRepository
}

func NewOcrRouter(
	ocrProductRepo *repositories.OcrProductRepository,
	purchaseInstalmentRepository *repositories.PurchaseInstalmentRepository,
) *Router {
	return &Router{
		ocrProductRepository:         ocrProductRepo,
		purchaseInstalmentRepository: purchaseInstalmentRepository,
	}
}

func (r *Router) GetRoutes(router *gin.RouterGroup) {
	router.POST("/exists", r.doOcrProductsExist)
	router.GET("/:name/exists", r.doesOcrProductExist)
	router.POST("/instalment", r.createPurchaseInstalment)
	router.POST("/instalment/list", r.createPurchaseInstalments)
}

func (r *Router) doesOcrProductExist(context *gin.Context) {
	name := context.Param("name")
	exists, _ := r.ocrProductRepository.Exists(name)
	if exists {
		context.Status(200)
		return
	}
	context.Status(404)
}

func (r *Router) doOcrProductsExist(context *gin.Context) {
	var ocrNamesDto ocr.ProductExists
	err := context.BindJSON(&ocrNamesDto)
	if err != nil {
		context.JSON(500, http.Response[helpers.None]{
			Err:        err.Error(),
			StatusCode: 500,
			Body:       helpers.None{},
		}.GetH())
		return
	}

	exists, _ := r.ocrProductRepository.ExistsMultiple(ocrNamesDto.OcrNames)
	context.JSON(200, http.Response[ocr.ProductExistsResponse]{
		Body: ocr.ProductExistsResponse{Exists: exists},
	}.GetH())
}

func (r *Router) createPurchaseInstalment(context *gin.Context) {
	var purchaseInstalmentDto product.CreatePurchaseInstalmentDto
	err := context.BindJSON(&purchaseInstalmentDto)
	if err != nil {
		context.JSON(500, http.Response[helpers.None]{
			Err:        err.Error(),
			StatusCode: 500,
			Body:       helpers.None{},
		}.GetH())
		return
	}

	purchaseInstalment, err := r.purchaseInstalmentRepository.CreatePurchaseInstalment(purchaseInstalmentDto)
	if err != nil {
		context.JSON(500, http.Response[helpers.None]{
			Err:        err.Error(),
			StatusCode: 500,
			Body:       helpers.None{},
		}.GetH())
		return
	}

	context.JSON(200, http.Response[product.PurchaseInstalmentDto]{
		Err:        "",
		Body:       purchaseInstalment.ToDto(),
		StatusCode: 200,
	}.GetH())
	return
}

func (r *Router) createPurchaseInstalments(context *gin.Context) {
	var purchaseInstalmentsDto product.CreatePurchaseInstalmentListDto
	err := context.BindJSON(&purchaseInstalmentsDto)
	if err != nil {
		context.JSON(500, http.Response[helpers.None]{
			Err:        err.Error(),
			StatusCode: 500,
			Body:       helpers.None{},
		}.GetH())
		return
	}

	purchaseInstalments, err := r.purchaseInstalmentRepository.CreatePurchaseInstalments(purchaseInstalmentsDto)
	if err != nil {
		context.JSON(500, http.Response[helpers.None]{
			Err:        err.Error(),
			StatusCode: 500,
			Body:       helpers.None{},
		}.GetH())
		return
	}

	context.JSON(200, http.Response[[]product.PurchaseInstalmentDto]{
		Err: "",
		Body: functional.Map(
			purchaseInstalments,
			func(purchaseInstalment productModel.PurchaseInstalmentModel) product.PurchaseInstalmentDto {
				return purchaseInstalment.ToDto()
			},
		),
		StatusCode: 200,
	}.GetH())
	return
}

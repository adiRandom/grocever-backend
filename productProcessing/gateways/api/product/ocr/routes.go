package ocr

import (
	"github.com/gin-gonic/gin"
	"lib/data/dto/product/ocr"
	"lib/helpers"
	"lib/network/http"
	"productProcessing/data/database/repositories"
)

type Router struct {
	repository *repositories.OcrProductRepository
}

func NewOcrRouter(ocrProductRepo *repositories.OcrProductRepository) *Router {
	return &Router{
		repository: ocrProductRepo,
	}
}

func (r *Router) GetRoutes(router *gin.RouterGroup) {
	router.GET("/:name/exists", r.doesOcrProductExist)
	router.POST("/exists", r.doOcrProductsExist)
}

func (r *Router) doesOcrProductExist(context *gin.Context) {
	name := context.Param("name")
	exists, _ := r.repository.Exists(name)
	if exists {
		context.Status(200)
		return
	}
	context.Status(404)
}

func (r *Router) doOcrProductsExist(context *gin.Context) {
	var ocrNamesDto ocr.OcrProductExists
	err := context.BindJSON(&ocrNamesDto)
	if err != nil {
		context.JSON(500, http.Response[helpers.None]{
			Err:        err.Error(),
			StatusCode: 500,
			Body:       helpers.None{},
		}.GetH())
		return
	}

	exists, _ := r.repository.ExistsMultiple(ocrNamesDto.OcrNames)
	context.JSON(200, http.Response[ocr.OcrProductExistsResponse]{
		Body: ocr.OcrProductExistsResponse{Exists: exists},
	}.GetH())
}

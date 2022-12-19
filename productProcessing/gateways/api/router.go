package api

import (
	"github.com/gin-gonic/gin"
	dto "lib/data/dto/product_processing"
	"lib/helpers"
	"lib/network/http"
	"productProcessing/data/database/repositories"
)

type Router struct {
	engine     *gin.Engine
	repository *repositories.OcrProductRepository
}

var router *Router = nil

func GetRouter() *Router {
	if router == nil {
		router = &Router{
			engine:     gin.Default(),
			repository: repositories.GetOcrProductRepository(),
		}

		router.initEndpoints()
	}
	return router
}

func (c *Router) initEndpoints() {
	c.engine.GET("/product/ocr/:name/exists", c.doesOcrProductExist)
	c.engine.POST("/product/ocr/exists", c.doOcrProductsExist)
}

func (c *Router) Run(port string) {
	err := c.engine.Run(port)
	if err != nil {
		panic(err)
	}
}

func (c *Router) doesOcrProductExist(context *gin.Context) {
	name := context.Param("name")
	exists, _ := c.repository.Exists(name)
	if exists {
		context.Status(200)
		return
	}
	context.Status(404)
}

func (c *Router) doOcrProductsExist(context *gin.Context) {
	var ocrNamesDto dto.OcrProductExists
	err := context.BindJSON(&ocrNamesDto)
	if err != nil {
		context.JSON(500, http.Response[helpers.None]{
			Err:        err.Error(),
			StatusCode: 500,
			Body:       helpers.None{},
		}.GetH())
		return
	}

	exists, _ := c.repository.ExistsMultiple(ocrNamesDto.OcrNames)
	context.JSON(200, http.Response[dto.OcrProductExistsResponse]{
		Body: dto.OcrProductExistsResponse{Exists: exists},
	}.GetH())
}

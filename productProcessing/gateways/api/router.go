package api

import (
	"github.com/gin-gonic/gin"
	"lib/api"
	"productProcessing/data/database/repositories"
	"productProcessing/gateways/api/product"
	"productProcessing/gateways/api/product/ocr"
)

type Router struct {
	api.Router
	engine                *gin.Engine
	ocrProductRepository  *repositories.OcrProductRepository
	userProductRepository *repositories.UserProductRepository
}

var router *Router = nil

func GetRouter() *Router {
	if router == nil {
		router = &Router{
			engine:                gin.Default(),
			ocrProductRepository:  repositories.GetOcrProductRepository(),
			userProductRepository: repositories.GetUserProductRepository(),
		}

		router.initEndpoints()
	}
	return router
}

func (c *Router) initEndpoints() {
	c.Group("/product/ocr", ocr.NewOcrRouter(c.ocrProductRepository))
	c.Group("/product", product.NewProductRouter(c.userProductRepository))
}

func (c *Router) Run(port string) {
	err := c.engine.Run(port)
	if err != nil {
		panic(err)
	}
}

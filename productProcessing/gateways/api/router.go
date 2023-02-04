package api

import (
	"lib/api"
	"productProcessing/data/database/repositories"
	"productProcessing/gateways/api/product"
	"productProcessing/gateways/api/product/ocr"
)

type Router struct {
	api.Router
	ocrProductRepository  *repositories.OcrProductRepository
	userProductRepository *repositories.PurchaseInstalmentRepository
}

var router *Router = nil

func GetRouter() *Router {
	if router == nil {
		router = &Router{
			ocrProductRepository:  repositories.GetOcrProductRepository(),
			userProductRepository: repositories.GetUserProductRepository(),
		}
		router.Init()
		router.initEndpoints()
	}
	return router
}

func (c *Router) initEndpoints() {
	c.Group("/product/ocr", ocr.NewOcrRouter(c.ocrProductRepository))
	c.Group("/product", product.NewProductRouter(c.userProductRepository))
}

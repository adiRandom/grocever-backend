package api

import (
	"lib/api"
	"productProcessing/data/database/repositories"
	"productProcessing/gateways/api/product"
	"productProcessing/gateways/api/product/ocr"
)

type Router struct {
	api.Router
	ocrProductRepository         *repositories.OcrProductRepository
	purchaseInstalmentRepository *repositories.PurchaseInstalmentRepository
	productRepository            *repositories.ProductRepository
	missLinkRepository           *repositories.MissLinkRepository
}

var router *Router = nil

func GetRouter() *Router {
	if router == nil {
		router = &Router{
			ocrProductRepository: repositories.GetOcrProductRepository(
				repositories.GetMissLinkRepository(),
			),
			purchaseInstalmentRepository: repositories.GetUserProductRepository(),
			productRepository: repositories.GetProductRepository(
				repositories.GetMissLinkRepository(),
				repositories.GetOcrProductRepository(
					repositories.GetMissLinkRepository(),
				),
			),
			missLinkRepository: repositories.GetMissLinkRepository(),
		}
		router.Init()
		router.initEndpoints()
	}
	return router
}

func (c *Router) initEndpoints() {
	c.Group("/product/ocr", ocr.NewOcrRouter(c.ocrProductRepository, c.purchaseInstalmentRepository))
	c.Group("/product", product.NewProductRouter(c.purchaseInstalmentRepository, c.productRepository, c.missLinkRepository))
}

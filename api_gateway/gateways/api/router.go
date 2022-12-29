package api

import (
	"api_gateway/gateways/api/product"
	productApi "api_gateway/services/api/products"
	"lib/api"
)

type Router struct {
	api.Router
}

var router *Router = nil

func GetRouter() *Router {
	if router == nil {
		router = &Router{}
		router.Init()
		router.initEndpoints()
	}
	return router
}

func (c *Router) initEndpoints() {
	c.Group("/product", product.NewProductRouter(productApi.GetClient()))
}

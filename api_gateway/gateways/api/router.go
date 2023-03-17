package api

import (
	"api_gateway/gateways/api/auth"
	"api_gateway/gateways/api/notification"
	"api_gateway/gateways/api/ocr"
	"api_gateway/gateways/api/product"
	"api_gateway/gateways/api/store"
	authApi "api_gateway/services/api/auth"
	notificationApi "api_gateway/services/api/notification"
	ocrApi "api_gateway/services/api/ocr"
	productApi "api_gateway/services/api/products"
	storeApi "api_gateway/services/api/store"
	"lib/api"
	dto "lib/data/dto/auth"
)

type Router struct {
	api.Router
}

var router *Router = nil

func GetRouter() *Router {
	if router == nil {
		router = &Router{}
		router.
			Init().
			WithAuth(handleAuthVerification)
		router.initEndpoints()
	}
	return router
}

func (c *Router) initEndpoints() {
	c.GroupWithAuth("/product", product.NewProductRouter(productApi.GetClient()))
	c.Group("/auth", auth.NewAuthRouter(authApi.GetClient()))
	c.GroupWithAuth("/ocr", ocr.NewOcrRouter(ocrApi.GetClient()))
	c.Group("/store", store.NewStoreRouter(storeApi.GetClient()))
	c.GroupWithAuth("/notification", notification.NewNotificationRouter(notificationApi.GetClient()))
}

func handleAuthVerification(access string) (int, error) {
	res, err := authApi.GetClient().Validate(dto.ValidateRequest{AccessToken: access})
	if err != nil {
		return -1, err
	}
	return res.UserId, nil
}

func GetBaseRouter() *api.Router {
	return &GetRouter().Router
}

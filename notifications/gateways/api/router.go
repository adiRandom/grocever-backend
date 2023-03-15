package api

import (
	"lib/api"
	"notifications/data/database/repository"
	"notifications/gateways/api/token"
)

type Router struct {
	api.Router
	repository *repository.NotificationUserRepository
}

var router *Router = nil

func GetRouter() *Router {
	if router == nil {
		router = &Router{
			repository: repository.GetNotificationUserRepository(),
		}
		router.Init()
		router.initEndpoints()
	}
	return router
}

func (c *Router) initEndpoints() {
	c.Group("/token", token.NewTokenRouter(c.repository))
}

func GetBaseRouter() *api.Router {
	return &GetRouter().Router
}

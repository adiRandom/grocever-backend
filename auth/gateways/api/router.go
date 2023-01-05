package api

import (
	"auth/data/repository"
	"auth/gateways/api/auth"
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
	c.Group("/auth", auth.NewAuthRouter(repository.GetUserRepository()))
}

func GetBaseRouter() *api.Router {
	return &GetRouter().Router
}

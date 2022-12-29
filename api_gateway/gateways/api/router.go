package api

import (
	"github.com/gin-gonic/gin"
)

type Router struct {
	engine *gin.Engine
}

var router *Router = nil

func GetRouter() *Router {
	if router == nil {
		router = &Router{
			engine: gin.Default(),
		}

		router.initEndpoints()
	}
	return router
}

func (c *Router) initEndpoints() {
	c.engine.GET("/product/list", c.getProductList)
}

func (c *Router) Run(port string) {
	err := c.engine.Run(port)
	if err != nil {
		panic(err)
	}
}

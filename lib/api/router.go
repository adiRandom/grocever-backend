package api

import "github.com/gin-gonic/gin"

type Router struct {
	engine *gin.Engine
}

func (r *Router) Group(path string, g RouterGroup) *Router {
	routerGroup := r.engine.Group(path)
	g.GetRoutes(routerGroup)
	return r
}

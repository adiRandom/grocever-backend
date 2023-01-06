package api

import (
	"github.com/gin-gonic/gin"
	"lib/api/middleware"
)

type Router struct {
	engine      *gin.Engine
	authHandler middleware.AuthHandler
}

func (r *Router) Group(path string, g RouterGroup) *Router {
	routerGroup := r.engine.Group(path)
	g.GetRoutes(routerGroup)
	return r
}

func (r *Router) GroupWithAuth(path string, g RouterGroup) *Router {
	routerGroup := r.engine.Group(path)
	routerGroup.Use(middleware.AuthMiddleware(r.authHandler))
	g.GetRoutes(routerGroup)
	return r
}

func (r *Router) Run(port string) {
	println("Starting API on port: " + port)
	err := r.engine.Run(port)
	if err != nil {
		panic(err)
	}
}

func (r *Router) Init() *Router {
	r.engine = gin.Default()
	return r
}

func (r *Router) WithAuth(handler middleware.AuthHandler) *Router {
	r.authHandler = handler
	return r
}

func (r *Router) UseMiddleware(middleware gin.HandlerFunc) {
	r.engine.Use(middleware)
}

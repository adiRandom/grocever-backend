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

func (r *Router) Run(port string) {
	println("Starting API on port: " + port)
	err := r.engine.Run(port)
	if err != nil {
		panic(err)
	}
}

func (r *Router) Init() {
	r.engine = gin.Default()
}

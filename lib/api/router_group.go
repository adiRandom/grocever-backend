package api

import "github.com/gin-gonic/gin"

type RouterGroup interface {
	GetRoutes(router *gin.RouterGroup)
}

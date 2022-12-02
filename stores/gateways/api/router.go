package api

import (
	"github.com/gin-gonic/gin"
	"lib/helpers"
	"lib/network/http"
	"stores/data/database/entity"
	"stores/data/database/repository"
)

type Router struct {
	engine *gin.Engine
	repo   *repository.StoreMetadata
}

var router *Router = nil

func GetRouter() *Router {
	if router == nil {
		router = &Router{
			engine: gin.Default(),
			repo:   repository.GetStoreMetadataRepository(),
		}

		router.initEndpoints()
	}
	return router
}

func (c *Router) initEndpoints() {
	c.engine.GET("/store/list", c.getAllStores)
	c.engine.GET("/store/:name", c.getStoreByName)
}

func (c *Router) getAllStores(ctx *gin.Context) {
	stores, err := c.repo.GetAll()
	if err != nil {
		ctx.JSON(500,
			(http.Response[helpers.None]{
				Err:        err.Error(),
				StatusCode: 500,
				Body:       helpers.None{},
			}).GetH(),
		)
		return
	}
	ctx.JSON(200,
		(http.Response[[]entity.StoreMetadata]{
			Err:        err.Error(),
			StatusCode: 200,
			Body:       stores,
		}).GetH(),
	)
}

func (c *Router) getStoreByName(ctx *gin.Context) {
	name := ctx.Param("name")
	store, err := c.repo.GetByName(name)
	if err != nil {
		ctx.JSON(500,
			(http.Response[helpers.None]{
				Err:        err.Error(),
				StatusCode: 500,
				Body:       helpers.None{},
			}).GetH(),
		)
		return
	}
	ctx.JSON(200,
		(http.Response[entity.StoreMetadata]{
			Err:        err.Error(),
			StatusCode: 200,
			Body:       *store,
		}).GetH(),
	)
}

func (c *Router) Start() {
	err := c.engine.Run()
	if err != nil {
		panic(err)
		return
	}
}

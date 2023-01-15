package api

import (
	"github.com/gin-gonic/gin"
	dtos "lib/data/dto/store"
	"lib/functional"
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

	resBody := functional.Map(stores, func(store entity.StoreMetadata) dtos.MetadataDto {
		model := store.ToModel()
		return model.ToDto()
	})

	ctx.JSON(200,
		(http.Response[[]dtos.MetadataDto]{
			StatusCode: 200,
			Body:       resBody,
		}).GetH(),
	)
}

func (c *Router) getStoreByName(ctx *gin.Context) {
	name := ctx.Param("name")
	store, err := c.repo.GetByName(name)
	storeModel := store.ToModel()
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
		(http.Response[dtos.MetadataDto]{
			Err:        err.Error(),
			StatusCode: 200,
			Body:       storeModel.ToDto(),
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

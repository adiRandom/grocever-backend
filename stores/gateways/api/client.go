package api

import (
	"github.com/gin-gonic/gin"
	"lib/helpers"
	"lib/network/http"
	"stores/data/database/entity"
	"stores/data/database/repository"
)

type Client struct {
	engine *gin.Engine
	repo   *repository.StoreMetadata
}

var client *Client = nil

func GetClient() *Client {
	if client == nil {
		client = &Client{
			engine: gin.Default(),
			repo:   repository.GetStoreMetadataRepository(),
		}

		client.initEndpoints()
	}
	return client
}

func (c *Client) initEndpoints() {
	c.engine.GET("/store/list", c.getAllStores)
	c.engine.GET("/store/:name", c.getStoreByName)
}

func (c *Client) getAllStores(ctx *gin.Context) {
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

func (c *Client) getStoreByName(ctx *gin.Context) {
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

func (c *Client) Start() {
	err := c.engine.Run()
	if err != nil {
		panic(err)
		return
	}
}

package store

import (
	"api_gateway/services/api/store"
	"github.com/gin-gonic/gin"
	storeDto "lib/data/dto/store"
	"lib/helpers"
	"lib/network/http"
)

type Router struct {
	storeApi *store.Client
}

func NewStoreRouter(storeApi *store.Client) *Router {
	return &Router{storeApi}
}

func (r *Router) GetRoutes(router *gin.RouterGroup) {
	router.GET("/list", r.getAllStores)
}

func (r *Router) getAllStores(context *gin.Context) {
	stores, err := r.storeApi.GetAllStores()
	if err != nil {
		context.JSON(500, http.Response[helpers.None]{
			StatusCode: 500,
			Err:        err.Error(),
			Body:       helpers.None{},
		}.GetH())
		return
	}

	context.JSON(200, http.Response[[]storeDto.MetadataDto]{
		StatusCode: 200,
		Body:       stores,
		Err:        "",
	}.GetH())
}

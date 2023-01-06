package auth

import (
	"auth/data/repository"
	"auth/services"
	"github.com/gin-gonic/gin"
	"lib/data/dto/auth"
	"lib/helpers"
	"lib/network/http"
)

type Router struct {
	userRepository *repository.User
}

func NewAuthRouter(userRepository *repository.User) *Router {
	return &Router{userRepository}
}

func (r *Router) GetRoutes(router *gin.RouterGroup) {
	router.POST("/login", r.login)
	router.POST("/register", r.register)
	router.POST("/refresh", r.refresh)
}

func (r *Router) login(context *gin.Context) {
	authDto := auth.LoginRequest{}
	err := context.BindJSON(&authDto)
	if err != nil {
		context.JSON(400, http.Response[helpers.None]{
			StatusCode: 400,
			Err:        err.Error(),
			Body:       helpers.None{},
		}.GetH())
		return
	}

	response := services.HandleLogin(services.NewLoginDetails(&authDto, nil), r.userRepository)
	context.JSON(response.StatusCode, response.GetH())
}

func (r *Router) register(context *gin.Context) {
	authDto := auth.RegisterRequest{}
	err := context.BindJSON(&authDto)
	if err != nil {
		context.JSON(400, http.Response[helpers.None]{
			StatusCode: 400,
			Err:        err.Error(),
			Body:       helpers.None{},
		}.GetH())
		return
	}

	response := services.HandleRegister(authDto, r.userRepository)
	context.JSON(response.StatusCode, response.GetH())
}

func (r *Router) refresh(context *gin.Context) {
	refreshDto := auth.RefreshRequest{}
	err := context.BindJSON(&refreshDto)
	if err != nil {
		context.JSON(400, http.Response[helpers.None]{
			StatusCode: 400,
			Err:        err.Error(),
			Body:       helpers.None{},
		}.GetH())
		return
	}

	response := services.HandleRefresh(refreshDto, r.userRepository)
	context.JSON(response.StatusCode, response.GetH())
}

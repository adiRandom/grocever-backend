package auth

import (
	"api_gateway/services/api/auth"
	"github.com/gin-gonic/gin"
	dto "lib/data/dto/auth"
	"lib/helpers"
	"lib/network/http"
)

type Router struct {
	authApiClient *auth.Client
}

func NewAuthRouter(authApiClient *auth.Client) *Router {
	return &Router{authApiClient}
}

func (r *Router) GetRoutes(router *gin.RouterGroup) {
	router.POST("/login", r.login)
	router.POST("/register", r.register)
	router.POST("/refresh", r.refresh)
}

func (r *Router) login(c *gin.Context) {
	var body dto.LoginRequest
	err := c.BindJSON(&body)
	if err != nil {
		c.JSON(400, http.Response[helpers.None]{
			StatusCode: 400,
			Err:        err.Error(),
			Body:       helpers.None{},
		})
		return
	}

	res, apiErr := r.authApiClient.Login(body)
	if apiErr != nil {
		c.JSON(apiErr.Code, http.Response[helpers.None]{
			StatusCode: apiErr.Code,
			Err:        apiErr.Error(),
			Body:       helpers.None{},
		})
		return
	}

	c.JSON(200, http.Response[dto.AuthResponse]{
		StatusCode: 200,
		Err:        "",
		Body:       res,
	}.GetH())
}

func (r *Router) register(c *gin.Context) {
	var body dto.RegisterRequest
	err := c.BindJSON(&body)
	if err != nil {
		c.JSON(400, http.Response[helpers.None]{
			StatusCode: 400,
			Err:        err.Error(),
			Body:       helpers.None{},
		})
		return
	}

	res, apiErr := r.authApiClient.Register(body)
	if err != nil {
		c.JSON(apiErr.Code, http.Response[helpers.None]{
			StatusCode: apiErr.Code,
			Err:        apiErr.Error(),
			Body:       helpers.None{},
		})
		return
	}

	c.JSON(200, http.Response[dto.AuthResponse]{
		StatusCode: 200,
		Err:        "",
		Body:       res,
	}.GetH())
}

func (r *Router) refresh(c *gin.Context) {
	var body dto.RefreshRequest
	err := c.BindJSON(&body)
	if err != nil {
		c.JSON(400, http.Response[helpers.None]{
			StatusCode: 400,
			Err:        err.Error(),
			Body:       helpers.None{},
		})
	}

	res, apiErr := r.authApiClient.Refresh(body)
	if err != nil {
		c.JSON(apiErr.Code, http.Response[helpers.None]{
			StatusCode: apiErr.Code,
			Err:        apiErr.Error(),
			Body:       helpers.None{},
		})
	}

	c.JSON(200, http.Response[dto.RefreshResponse]{
		StatusCode: 200,
		Err:        "",
		Body:       res,
	}.GetH())
}

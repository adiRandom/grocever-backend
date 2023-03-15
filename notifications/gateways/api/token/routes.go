package token

import (
	"github.com/gin-gonic/gin"
	"lib/data/dto/notifications"
	"lib/helpers"
	"lib/network/http"
	"notifications/data/database/repository"
	"strconv"
)

type Router struct {
	repo *repository.NotificationUserRepository
}

func NewTokenRouter(
	notificationRepository *repository.NotificationUserRepository,
) *Router {
	return &Router{
		repo: notificationRepository,
	}
}

func (r *Router) GetRoutes(router *gin.RouterGroup) {
	router.PUT("/:userId", r.setToken)
}

func (r *Router) setToken(context *gin.Context) {
	var dto notifications.SetUserFcmTokenDto
	err := context.BindJSON(&dto)
	if err != nil {
		context.JSON(400, http.Response[helpers.None]{
			Err:        err.Error(),
			StatusCode: 400,
			Body:       helpers.None{},
		}.GetH())
		return
	}

	userId, err := strconv.Atoi(context.Param("userId"))
	if err != nil {
		context.JSON(400, http.Response[helpers.None]{
			Err:        err.Error(),
			StatusCode: 400,
			Body:       helpers.None{},
		}.GetH())
		return
	}

	err = r.repo.CreateOrUpdate(userId, dto.Token)
	if err != nil {
		context.JSON(500, http.Response[helpers.None]{
			Err:        err.Error(),
			StatusCode: 500,
			Body:       helpers.None{},
		}.GetH())
		return
	}

	context.Status(204)
}

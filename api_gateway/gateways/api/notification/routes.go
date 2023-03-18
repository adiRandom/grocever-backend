package notification

import (
	"api_gateway/services/api/notification"
	"github.com/gin-gonic/gin"
	"lib/api/middleware"
	"lib/data/dto/notifications"
	"lib/helpers"
	"lib/network/http"
)

type Router struct {
	notificationClient *notification.Client
}

func NewNotificationRouter(notificationClient *notification.Client) *Router {
	return &Router{notificationClient}
}

func (r *Router) GetRoutes(router *gin.RouterGroup) {
	router.PUT("/token", r.sendFcmToken)
}

func (r *Router) sendFcmToken(context *gin.Context) {
	var dto notifications.SetUserFcmTokenDto
	err := context.BindJSON(&dto)
	if err != nil {
		context.JSON(400, http.Response[helpers.None]{
			Err:        err.Error(),
			StatusCode: 500,
			Body:       helpers.None{},
		}.GetH())
		return
	}

	apiErr := r.notificationClient.SendFcmToken(dto, context.GetInt(middleware.UserIdKey))
	if apiErr != nil {
		context.JSON(500, http.Response[helpers.None]{
			Err:        apiErr.Error(),
			StatusCode: 500,
			Body:       helpers.None{},
		}.GetH())
		return
	}

	context.Status(204)
}

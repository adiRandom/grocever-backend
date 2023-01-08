package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"lib/data/dto/ocr"
	"lib/events/rabbitmq"
	"lib/helpers"
	"lib/network/http"
	"mime/multipart"
	"ocr/gateways/events"
)

type Router struct {
	engine *gin.Engine
	broker *rabbitmq.JsonBroker[ocr.UploadDto]
}

const imageParam = "image"
const userIdParam = "userId"

var router *Router = nil

func GetRouter() *Router {
	if router == nil {
		router = &Router{
			engine: gin.Default(),
			broker: events.GetRabbitMqBroker(),
		}

		router.initEndpoints()
	}
	return router
}

func (c *Router) initEndpoints() {
	c.engine.POST("/ocr", c.processImage)
}

func (c *Router) processImage(ctx *gin.Context) {
	image, err := ctx.FormFile(imageParam)
	if err != nil {
		ctx.JSON(500, http.Response[helpers.None]{
			Err:        err.Error(),
			StatusCode: 500,
			Body:       helpers.None{},
		}.GetH())
		return
	}

	userId, exists := ctx.GetPostForm(userIdParam)
	if !exists {
		ctx.JSON(400, http.Response[helpers.None]{
			Err:        "userId is required",
			StatusCode: 400,
			Body:       helpers.None{},
		}.GetH())
		return
	}

	// Read the image into a byte array
	imageFile, err := image.Open()
	if err != nil {
		ctx.JSON(500, http.Response[helpers.None]{
			Err:        err.Error(),
			StatusCode: 500,
			Body:       helpers.None{},
		}.GetH())
		return
	}
	defer func(imageFile multipart.File) {
		err := imageFile.Close()
		if err != nil {
			fmt.Println(err)
		}
	}(imageFile)

	// Read the image into a byte array
	imageBytes := make([]byte, image.Size)
	_, err = imageFile.Read(imageBytes)
	if err != nil {
		ctx.JSON(500, http.Response[helpers.None]{
			Err:        err.Error(),
			StatusCode: 500,
			Body:       helpers.None{},
		}.GetH())
		return
	}

	// Send the image to the OCR service
	c.broker.SendInput(ocr.UploadDto{
		Bytes:  imageBytes,
		Size:   image.Size,
		UserId: userId,
	})

	ctx.JSON(200, http.Response[helpers.None]{
		Err:        "",
		StatusCode: 200,
		Body:       helpers.None{},
	})
}

func (c *Router) Run(port string) {
	err := c.engine.Run(port)
	if err != nil {
		panic(err)
	}
}

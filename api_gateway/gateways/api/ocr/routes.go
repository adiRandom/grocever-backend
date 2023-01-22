package auth

import (
	"api_gateway/services/api/ocr"
	"fmt"
	"github.com/gin-gonic/gin"
	ocrDto "lib/data/dto/ocr"
	"lib/helpers"
	"lib/network/http"
	"mime/multipart"
	"strconv"
)

type Router struct {
	ocrApiClient *ocr.Client
}

func NewAuthRouter(ocrApiClient *ocr.Client) *Router {
	return &Router{ocrApiClient}
}

func (r *Router) GetRoutes(router *gin.RouterGroup) {
	router.POST("/", r.uploadImage)
}

func (r *Router) uploadImage(ctx *gin.Context) {
	image, err := ctx.FormFile(ocrDto.UploadImageParam)
	if err != nil {
		ctx.JSON(500, http.Response[helpers.None]{
			Err:        err.Error(),
			StatusCode: 500,
			Body:       helpers.None{},
		}.GetH())
		return
	}

	userId, exists := ctx.GetPostForm(ocrDto.UploadImageUserIdParam)
	if !exists {
		ctx.JSON(400, http.Response[helpers.None]{
			Err:        "userId is required",
			StatusCode: 400,
			Body:       helpers.None{},
		}.GetH())
		return
	}

	userIdInt, err := strconv.Atoi(userId)
	if err != nil {
		ctx.JSON(401, http.Response[helpers.None]{
			Err:        "Wrong format of userId",
			StatusCode: 401,
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

	err = r.ocrApiClient.UploadImage(*ocrDto.NewUploadImageRequest(&imageFile, userIdInt))
	if err != nil {
		ctx.JSON(500, http.Response[helpers.None]{
			Err:        err.Error(),
			StatusCode: 500,
			Body:       helpers.None{},
		}.GetH())
		return
	} else {
		ctx.JSON(200, http.Response[helpers.None]{
			Err:        "",
			StatusCode: 200,
			Body:       helpers.None{},
		}.GetH())
		return
	}

}

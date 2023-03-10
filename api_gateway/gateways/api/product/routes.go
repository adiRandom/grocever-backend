package product

import (
	"api_gateway/services/api/products"
	"github.com/gin-gonic/gin"
	"lib/api/middleware"
	"lib/data/dto/product"
	"lib/events/rabbitmq"
	"lib/helpers"
	"lib/network/amqp"
	"lib/network/http"
)

type Router struct {
	productApiClient *products.Client
}

func NewProductRouter(productApiClient *products.Client) *Router {
	return &Router{productApiClient}
}

func (r *Router) GetRoutes(router *gin.RouterGroup) {
	router.GET("/list", r.getAllUserProducts)
	router.POST("", r.createPurchaseInstalmentNoOcr)
	router.POST("/report", r.reportMissLink)
	router.GET("/report/list", r.getUserReports)
}

func (r *Router) getAllUserProducts(context *gin.Context) {
	userId, exists := context.Get(middleware.UserIdKey)
	if !exists {
		context.JSON(401, http.Response[helpers.None]{
			StatusCode: 401,
			Err:        "Unauthorized",
			Body:       helpers.None{},
		})
		return
	}

	productList, apiError := r.productApiClient.GetProductList(userId.(int))
	if apiError != nil {
		context.JSON(apiError.Code, http.Response[helpers.None]{
			StatusCode: apiError.Code,
			Err:        apiError.Error(),
			Body:       helpers.None{},
		}.GetH())
		return
	}
	context.JSON(200, http.Response[product.UserProductListDto]{
		StatusCode: 200,
		Body: product.UserProductListDto{
			Products: productList,
		},
		Err: "",
	}.GetH())
}

func (r *Router) createPurchaseInstalmentNoOcr(context *gin.Context) {
	userId, exists := context.Get(middleware.UserIdKey)
	if !exists {
		context.JSON(401, http.Response[helpers.None]{
			StatusCode: 401,
			Err:        "Unauthorized",
			Body:       helpers.None{},
		})
		return
	}

	var dto product.CreatePurchaseInstalmentNoOcrDto
	err := context.BindJSON(&dto)
	if err != nil {
		context.JSON(400, http.Response[helpers.None]{
			Err:        err.Error(),
			StatusCode: 400,
			Body:       helpers.None{},
		}.GetH())
		return
	}

	dtoWithUserId := product.CreatePurchaseInstalmentNoOcrWithUserDto{
		CreatePurchaseInstalmentNoOcrDto: dto, UserId: uint(userId.(int)),
	}

	resDto, apiErr := r.productApiClient.CreatePurchaseInstalmentNoOcr(userId.(int), dtoWithUserId)
	if apiErr != nil {
		context.JSON(500, http.Response[helpers.None]{
			Err:        apiErr.Error(),
			StatusCode: 500,
			Body:       helpers.None{},
		}.GetH())
		return
	}

	searchDto := product.PurchaseInstalmentWithUserDto{
		PurchaseInstalmentDto: *resDto,
		UserId:                userId.(int),
	}

	err = rabbitmq.PushToQueue[product.PurchaseInstalmentWithUserDto](amqp.SearchQueue, searchDto)
	if err != nil {
		context.JSON(500, http.Response[helpers.None]{
			Err:        "Failed to push to queue",
			StatusCode: 500,
			Body:       helpers.None{},
		}.GetH())
		return
	}

	context.JSON(200, http.Response[product.PurchaseInstalmentDto]{
		Err:        "",
		Body:       *resDto,
		StatusCode: 200,
	}.GetH())
}

func (r *Router) reportMissLink(context *gin.Context) {
	userId, exists := context.Get(middleware.UserIdKey)
	if !exists {
		context.JSON(401, http.Response[helpers.None]{
			StatusCode: 401,
			Err:        "Unauthorized",
			Body:       helpers.None{},
		})
		return
	}

	var dto product.ReportDto
	err := context.BindJSON(&dto)
	if err != nil {
		context.JSON(400, http.Response[helpers.None]{
			Err:        err.Error(),
			StatusCode: 400,
			Body:       helpers.None{},
		}.GetH())
		return
	}

	dtoWithUserId := product.NewReportWithUserIdDto(dto.ProductId, dto.OcrProductName, userId.(int))
	apiErr := r.productApiClient.ReportMissLink(*dtoWithUserId)
	if apiErr != nil {
		context.JSON(500, http.Response[helpers.None]{
			Err:        apiErr.Error(),
			StatusCode: 500,
			Body:       helpers.None{},
		}.GetH())
		return
	}

	context.JSON(200, http.Response[helpers.None]{
		Err:        "",
		Body:       helpers.None{},
		StatusCode: 200,
	}.GetH())
}

func (r *Router) getUserReports(context *gin.Context) {
	userId, exists := context.Get(middleware.UserIdKey)
	if !exists {
		context.JSON(401, http.Response[helpers.None]{
			StatusCode: 401,
			Err:        "Unauthorized",
			Body:       helpers.None{},
		})
		return
	}

	reports, apiErr := r.productApiClient.GetReportsByUser(userId.(int))
	if apiErr != nil {
		context.JSON(500, http.Response[helpers.None]{
			Err:        apiErr.Error(),
			StatusCode: 500,
			Body:       helpers.None{},
		}.GetH())
		return
	}

	context.JSON(200, http.Response[[]product.ReportDto]{
		Err:        "",
		Body:       *reports,
		StatusCode: 200,
	}.GetH())
}

package products

import (
	"fmt"
	"lib/data/dto/product"
	"lib/network/http"
	"os"
)

type Client struct {
	baseUrl string
}

var client *Client = nil

func GetClient() *Client {
	if client == nil {
		client = &Client{
			baseUrl: os.Getenv("PRODUCT_PROCESSING_API_BASE_URL"),
		}
	}
	return client
}

func (s *Client) GetProductList(userId int) ([]product.UserProductDto, *http.Error) {
	res, err := http.UnwrapHttpResponse(
		http.GetSync[http.Response[product.UserProductListDto]](fmt.Sprintf(s.baseUrl+"/product/%d/list", userId)),
	)
	if err != nil {
		return nil, err
	}
	return res.Products, nil
}

func (s *Client) CreatePurchaseInstalmentNoOcr(userId int, purchaseInstalmentNoOcrDto product.CreatePurchaseInstalmentNoOcrWithUserDto) (*product.PurchaseInstalmentDto, *http.Error) {
	res, err := http.UnwrapHttpResponse(
		http.PostSync[http.Response[product.PurchaseInstalmentDto]](fmt.Sprintf(s.baseUrl+"/product/%d", userId), purchaseInstalmentNoOcrDto),
	)
	if err != nil {
		return nil, err
	}
	return res, nil
}

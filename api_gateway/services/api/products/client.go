package products

import (
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

func (s *Client) GetProductList() ([]product.PurchaseInstalmentDto, *http.Error) {
	res, err := http.ParseHttpResponse(
		http.GetSync[http.Response[product.UserProductListDto]](s.baseUrl + "/product/list"),
	)
	if err != nil {
		return nil, err
	}
	return res.Products, nil
}

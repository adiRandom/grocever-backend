package store

import (
	"lib/data/dto/product_processing"
	"lib/network/http"
	"os"
)

type Client struct {
	baseUrl string
}

var client *Client

func GetClient() *Client {
	if client == nil {
		client = &Client{
			baseUrl: os.Getenv("PRODUCT_PROCESSING_API_BASE_URL"),
		}
	}
	return client
}

func (s *Client) GetProductList() []product_processing.UserProductDto {
	res, err := http.GetSync[http.Response[product_processing.UserProductListDto]](s.baseUrl + "/product/list")
	if err != nil {
		return nil
	}
	return res.Body.Products
}

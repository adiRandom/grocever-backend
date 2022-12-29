package store

import (
	"lib/data/dto"
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
			baseUrl: os.Getenv("STORES_API_BASE_URL"),
		}
	}
	return client
}

func (s *Client) GetAllStores() []dto.StoreMetadata {
	res, err := http.ParseHttpResponse(
		http.GetSync[http.Response[[]dto.StoreMetadata]](s.baseUrl + "/store/list"),
	)
	if err != nil {
		return nil
	}
	return *res
}

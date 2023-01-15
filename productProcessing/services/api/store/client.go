package store

import (
	"lib/data/dto/store"
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
			baseUrl: os.Getenv("STORE_API_BASE_URL"),
		}
	}
	return client
}

func (s *Client) GetAllStores() ([]store.MetadataDto, error) {
	res, err := http.ParseHttpResponse(
		http.GetSync[http.Response[[]store.MetadataDto]](s.baseUrl + "/store/list"),
	)
	if err != nil {
		return nil, err
	}
	return *res, nil
}

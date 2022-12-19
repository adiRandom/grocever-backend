package store

import (
	"lib/data/dto"
	"lib/functional"
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

func (s *Client) GetStoreMetadataForName(name string) (dto.StoreMetadata, error) {
	res, err := http.GetSync[dto.StoreMetadata](s.baseUrl + "/store/" + name)
	if err != nil {
		return dto.StoreMetadata{}, err
	}

	return *res, nil
}

func (s *Client) GetAllStoreNames() []string {
	res, err := http.GetSync[[]dto.StoreMetadata](s.baseUrl + "/store/list")
	if err != nil {
		return nil
	}

	return functional.Map[dto.StoreMetadata, string](*res, func(store dto.StoreMetadata) string {
		return store.Name
	})
}

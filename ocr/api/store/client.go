package store

import (
	"lib/data/dto/store"
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
			baseUrl: os.Getenv("STORE_API_BASE_URL"),
		}
	}
	return client
}

func (s *Client) GetStoreMetadataForName(name string) (store.MetadataDto, error) {
	res, err := http.UnwrapHttpResponse(http.GetSync[http.Response[store.MetadataDto]](s.baseUrl + "/store/" + name))
	if err != nil {
		return store.MetadataDto{}, err
	}

	return *res, nil
}

func (s *Client) GetAllStoreNames() []string {
	res, err := http.UnwrapHttpResponse(http.GetSync[http.Response[[]store.MetadataDto]](s.baseUrl + "/store/list"))
	if err != nil {
		return nil
	}

	return functional.Map[store.MetadataDto, string](*res, func(store store.MetadataDto) string {
		return store.Name
	})
}

package product_processing

import (
	dto "lib/data/dto/product_processing/ocr"
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

func (s *Client) OcrProductsExists(ocrNames []string) ([]bool, error) {
	res, err := http.ParseHttpResponse(http.PostSync[http.Response[dto.OcrProductExistsResponse]](
		s.baseUrl+"/product/exists",
		dto.OcrProductExists{
			OcrNames: ocrNames,
		},
	))

	if err != nil {
		return nil, err
	}

	return res.Exists, nil
}

package nlp

import (
	"lib/network/http"
	"os"
	"productProcessing/data/dto"
)

type Client struct {
	baseUrl string
	headers map[string]string
}

var client *Client

func GetClient() *Client {
	if client == nil {
		headers := map[string]string{
			"Content-Type": "application/json",
			"X-Api-Key":    os.Getenv("NLP_API_KEY"),
		}
		client = &Client{
			baseUrl: os.Getenv("NLP_API_BASE_URL"),
			headers: headers,
		}
	}
	return client
}

func (c *Client) GetSimilarity(ocrProductName, productName string) float64 {
	res, err := http.PostWithHeadersSync[dto.TextSimilarityResultDto](c.baseUrl+"textsimilarity", dto.TextSimilarityDto{
		Text1: ocrProductName,
		Text2: productName,
	}, c.headers)
	if err != nil {
		return 0
	}

	return res.Similarity
}

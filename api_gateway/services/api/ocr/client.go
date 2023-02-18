package ocr

import (
	"lib/data/dto/ocr"
	"lib/helpers"
	"lib/network/http"
	"os"
	"strconv"
)

type Client struct {
	baseUrl string
}

var client *Client = nil

func GetClient() *Client {
	if client == nil {
		client = &Client{
			baseUrl: os.Getenv("OCR_API_BASE_URL"),
		}
	}
	return client
}

func (s *Client) UploadImage(body ocr.UploadImageRequest) *http.Error {
	reqBody := make(map[string]http.PostFormValue)
	reqBody[ocr.UploadImageParam] = body.Image
	reqBody[ocr.UploadImageUserIdParam] = strconv.Itoa(body.UserId)
	_, err := http.UnwrapHttpResponse[helpers.None](http.PostFormSync[http.Response[helpers.None]](s.baseUrl+"/ocr", reqBody))
	return err
}

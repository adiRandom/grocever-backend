package notification

import (
	"fmt"
	"lib/data/dto/notifications"
	"lib/helpers"
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
			baseUrl: os.Getenv("NOTIFICATION_API_BASE_URL"),
		}
	}
	return client
}

func (s *Client) SendFcmToken(body notifications.SetUserFcmTokenDto, userId int) *http.Error {
	_, err := http.UnwrapHttpResponse[helpers.None](http.PutSync[http.Response[helpers.None]](
		fmt.Sprintf(s.baseUrl+"/token/%d", userId), body),
	)
	if err != nil {
		return err
	}
	return nil
}

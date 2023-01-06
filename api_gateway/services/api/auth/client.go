package auth

import (
	"lib/data/dto/auth"
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
			baseUrl: os.Getenv("AUTH_API_BASE_URL"),
		}
	}
	return client
}

func (s *Client) Login(body auth.LoginRequest) (auth.AuthResponse, *http.Error) {
	res, err := http.ParseHttpResponse[auth.AuthResponse](http.PostSync[http.Response[auth.AuthResponse]](s.baseUrl+"/auth/login", body))
	if err != nil {
		return auth.AuthResponse{}, err
	}
	return *res, nil
}

func (s *Client) Register(body auth.RegisterRequest) (auth.AuthResponse, *http.Error) {
	res, err := http.ParseHttpResponse[auth.AuthResponse](http.PostSync[http.Response[auth.AuthResponse]](s.baseUrl+"/auth/register", body))
	if err != nil {
		return auth.AuthResponse{}, err
	}
	return *res, nil
}

func (s *Client) Refresh(body auth.RefreshRequest) (auth.RefreshResponse, *http.Error) {
	res, err := http.ParseHttpResponse[auth.RefreshResponse](http.PostSync[http.Response[auth.RefreshResponse]](s.baseUrl+"/auth/refresh", body))
	if err != nil {
		return auth.RefreshResponse{}, err
	}
	return *res, nil
}

func (s *Client) Validate(body auth.ValidateRequest) (auth.ValidateResponse, *http.Error) {
	res, err := http.ParseHttpResponse[auth.ValidateResponse](http.PostSync[http.Response[auth.ValidateResponse]](s.baseUrl+"/auth/validate", body))
	if err != nil {
		return auth.ValidateResponse{}, err
	}
	return *res, nil
}

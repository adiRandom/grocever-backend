package http

import (
	"lib/helpers"
)

func ParseHttpResponse[T any](response *Response[T], err error) (*T, error) {
	if response.Err != "" {
		return nil, helpers.HttpError{Msg: response.Err, Code: response.StatusCode}
	}
	return &response.Body, nil
}

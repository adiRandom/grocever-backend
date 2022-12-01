package http

import (
	"github.com/gin-gonic/gin"
	"reflect"
)

type Response[T any] struct {
	StatusCode int    `json:"statusCode"`
	Body       T      `json:"body"`
	Err        string `json:"err"`
}

func (r Response[T]) GetH() gin.H {

	if reflect.TypeOf(r.Body).Name() == "None" {
		return gin.H{
			"statusCode": r.StatusCode,
			"err":        r.Err,
		}
	}

	return gin.H{
		"statusCode": r.StatusCode,
		"body":       r.Body,
		"err":        r.Err,
	}
}

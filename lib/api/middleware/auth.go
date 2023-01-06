package middleware

import (
	"github.com/gin-gonic/gin"
	"lib/helpers"
	"lib/network/http"
)

// AuthHandler
// A function that validates the token and returns the user id
// If the token is invalid returns an error
type AuthHandler = func(access string) (int, error)

const UserIdKey = "user_id"

func AuthMiddleware(handler AuthHandler) gin.HandlerFunc {
	return func(c *gin.Context) {
		access := c.GetHeader("Authorization")
		if access == "" {
			c.AbortWithStatusJSON(401, http.Response[helpers.None]{
				StatusCode: 401,
				Err:        "Unauthorized",
				Body:       helpers.None{},
			})
			return
		}

		token := access[7:]
		userId, err := handler(token)
		if err != nil {
			c.AbortWithStatusJSON(401, http.Response[helpers.None]{
				StatusCode: 401,
				Err:        "Unauthorized",
				Body:       helpers.None{},
			})
		}

		c.Set(UserIdKey, userId)
		c.Next()
	}
}

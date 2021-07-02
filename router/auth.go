package router

import (
	"auth/serializer"
	"auth/util/jwt"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

// AuthRequired 需要登录
func AuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		uid, err := jwt.Get(c)
		if err == nil {
			c.Set("uid", uid)
			c.Next()
			return
		}

		c.JSON(200, serializer.Err(http.StatusUnauthorized, errors.New("请先登录")))
		c.Abort()
	}
}

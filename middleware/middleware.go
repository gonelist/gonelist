package middleware

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"GOIndex/mg_auth"
)

func CheckLogin() gin.HandlerFunc {
	return func(c *gin.Context) {
		if mg_auth.GetClient() == nil {
			// 没有 Client 则重定向到登陆
			c.Redirect(http.StatusFound, "/login")
			c.Abort()
		}
	}
}

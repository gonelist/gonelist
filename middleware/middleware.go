package middleware

import (
	"GOIndex/mg_auth"
	"GOIndex/pkg/app"
	"GOIndex/pkg/e"
	"github.com/gin-gonic/gin"
	"net/http"
)

func CheckLogin() gin.HandlerFunc {
	return func(c *gin.Context) {
		if mg_auth.GetClient() == nil {
			// 没有 Client 则重定向到登陆
			app.Response(c, http.StatusOK, e.REDIRECT, "需要重定向")
			//c.Redirect(http.StatusFound, "/login")
			c.Abort()
		}
	}
}

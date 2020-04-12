package middleware

import (
	"gonelist/mg_auth"
	"gonelist/pkg/app"
	"gonelist/pkg/e"
	"github.com/gin-gonic/gin"
	"net/http"
)

func CheckLogin() gin.HandlerFunc {
	return func(c *gin.Context) {
		if mg_auth.GetClient() == nil {
			// 没有 Client 则重定向到登陆
			app.Response(c, http.StatusOK, e.REDIRECT_LOGIN, "需要重定向到登陆")
			//c.Redirect(http.StatusFound, "/login")
			c.Abort()
		}
	}
}

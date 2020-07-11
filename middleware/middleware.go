package middleware

import (
	"github.com/gin-gonic/gin"
	"gonelist/onedrive"
	"gonelist/pkg/app"
	"gonelist/pkg/e"
	"net/http"
)

func CheckLogin() gin.HandlerFunc {
	return func(c *gin.Context) {
		if onedrive.GetClient() == nil {
			// 没有 Client 则重定向到登陆
			app.Response(c, http.StatusOK, e.REDIRECT_LOGIN, e.GetMsg(e.REDIRECT_LOGIN))
			//c.Redirect(http.StatusFound, "/login")
			c.Abort()
		}
	}
}

func CheckOnedriveInit() gin.HandlerFunc {
	return func(c *gin.Context) {
		if !onedrive.FileTree.IsLogin() {
			// 判断是否初始化完成
			app.Response(c, http.StatusOK, e.LOAD_NOT_READY, e.GetMsg(e.LOAD_NOT_READY))
			//c.Redirect(http.StatusFound, "/login")
			c.Abort()
		}
	}
}

package middleware

import (
	"github.com/gin-gonic/gin"
	"gonelist/mg_auth"
	"gonelist/onedrive"
	"gonelist/pkg/app"
	"gonelist/pkg/e"
	"net/http"
	"sync"
)

var cacheGoOnce sync.Once

func CheckLogin() gin.HandlerFunc {
	return func(c *gin.Context) {
		if mg_auth.GetClient() == nil {
			// 没有 Client 则重定向到登陆
			app.Response(c, http.StatusOK, e.REDIRECT_LOGIN, "需要重定向到登陆")
			//c.Redirect(http.StatusFound, "/login")
			c.Abort()
		} else {
			cacheGoOnce.Do(func() {
				onedrive.GetAllFiles()
				mg_auth.IsLogin = true
				// 如果首页有 README.md 则下载到本地
				onedrive.DownloadREADME()
				go onedrive.SetAutoRefresh()
			})
		}
	}
}

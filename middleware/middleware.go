package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"gonelist/conf"
	"gonelist/pkg/app"
	"gonelist/pkg/e"
	"gonelist/service/onedrive"
	"gonelist/service/onedrive/cache"
)

// 判断 onedrive 是否 login
func CheckLogin() gin.HandlerFunc {
	return func(c *gin.Context) {
		if onedrive.GetClient() == nil {
			// 没有 Client 则重定向到登陆
			app.Response(c, http.StatusOK, e.REDIRECT_LOGIN, nil)
			//c.Redirect(http.StatusFound, "/login")
			c.Abort()
		}
	}
}

// TODO 登陆之后有一个等待初始化的时间
// 如果没有初始化完成需要返回给前端做判断
func CheckOnedriveInit() gin.HandlerFunc {
	return func(c *gin.Context) {
		if !onedrive.FileTree.IsLogin() {
			// 判断是否初始化完成
			app.Response(c, http.StatusOK, e.LOAD_NOT_READY, nil)
			//c.Redirect(http.StatusFound, "/login")
			c.Abort()
		}
	}
}

// 判断文件夹密码是否正确
func CheckFolderPass() gin.HandlerFunc {
	return func(c *gin.Context) {
		p := c.Query("path")
		if conf.UserSet.Local.Enable && strings.HasPrefix(p, "/"+conf.UserSet.Local.Name) {
			c.Next()
			return
		}
		pass := c.GetHeader("pass")
		// 判断 config.yml 中的密码
		if !onedrive.CheckPassCorrect(p, pass) {
			// 如果密码错误，则返回
			app.Response(c, http.StatusOK, e.PASS_ERROR, nil)
			c.Abort()
		}
		// 判断路径下是否有 .password 文件
		node, ok := cache.Cache.Get(p)
		if !ok || node == nil {
			app.Response(c, http.StatusOK, e.PASS_ERROR, nil)
			c.Abort()
		}

		if node.Password == "" || node.Password == pass {
			c.Next()
		} else {
			app.Response(c, http.StatusOK, e.PASS_ERROR, nil)
			c.Abort()
		}
	}
}

// CheckSecret
/**
 * @Description: 用户权限鉴权，需要鉴权的api有，上传文件，创建文件夹，删除文件
 * @return gin.HandlerFunc
 */
func CheckSecret() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		secret := ctx.Query("secret")
		if secret != conf.UserSet.Admin.Secret {
			app.Response(ctx, http.StatusOK, e.SECRET_ERROR, nil)
			ctx.Abort()
			return
		}
		ctx.Next()
	}
}

package api

import (
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"gonelist/conf"
	"gonelist/onedrive"
	"gonelist/pkg/app"
	"gonelist/pkg/e"
	"net/http"
)

// 通过监听一个地址，跳转打开 login
func Login(c *gin.Context) {
	conf.UserSet.Server.SiteUrl = c.GetHeader("Host")
	// 判断是否登录
	if onedrive.FileTree.IsLogin() == true {
		// 有 Client 则重定向到首页
		c.Redirect(http.StatusFound, "/onedrive/getpath?path=/")
	} else {
		c.Redirect(http.StatusFound, "/loginmg")
	}
	c.Abort()
}

// 跳转到网页登录
func LoginMG(c *gin.Context) {
	onedrive.RedirectLoginMG(c)
	c.Abort()
}

// 接受 code
func GetCode(c *gin.Context) {
	var err error
	code := &onedrive.ReceiveCode{
		Code: c.Query("code"),
		//SessionState: c.Query("session_state"), // 有的账号好像没有
		State: c.Query("state"),
	}
	err = c.ShouldBind(code)
	if err != nil {
		log.Warn(err)
		app.Response(c, http.StatusOK, e.INVALID_PARAMS, nil)
		return
	}
	// 获取 AccessToken
	err = onedrive.GetAccessToken(*code)
	if err != nil {
		app.Response(c, http.StatusOK, e.GetErrorCode(err), "登陆失败，请重新登陆")
	} else {
		// 初始化 onedrive 的连接
		onedrive.InitOnedive()
		// 登陆成功跳转到网站首页
		c.Redirect(http.StatusTemporaryRedirect, conf.UserSet.Server.SiteUrl)
		//app.Response(c, http.StatusOK, e.SUCCESS, "登陆成功")
	}
}

// 注销登陆
func CancelLogin(c *gin.Context) {
	//mg_auth.ClearCLient()
	onedrive.FileTree.SetLogin(false)
	app.Response(c, http.StatusOK, e.SUCCESS, "注销成功")
}

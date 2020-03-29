package api

import (
	"GOIndex/mg_auth"
	"GOIndex/onedrive"
	"GOIndex/pkg/app"
	"GOIndex/pkg/e"
	"github.com/gin-gonic/gin"
	"net/http"
)

// 通过监听一个地址，跳转打开 login
func Login(c *gin.Context) {
	// 判断是否登录
	if mg_auth.GetClient() != nil {
		// 有 Client 则重定向到首页
		c.Redirect(http.StatusFound, "/onedrive/getallfiles")
	} else {
		c.Redirect(http.StatusFound, "/loginmg")
	}
	c.Abort()
}

// 跳转到网页登录
func LoginMG(c *gin.Context) {
	mg_auth.RedirectLoginMG(c)
	c.Abort()
}

// 接受 code
func GetCode(c *gin.Context) {
	var err error
	code := &mg_auth.ReceiveCode{
		Code:         c.Query("code"),
		SessionState: c.Query("session_state"),
		State:        c.Query("state"),
	}
	err = c.ShouldBind(code)
	if err != nil {
		app.Response(c, http.StatusOK, e.INVALID_PARAMS, nil)
		return
	}
	// 获取 AccessToken
	err = mg_auth.GetAccessToken(*code)
	if err != nil {
		app.Response(c, http.StatusOK, e.GetErrorCode(err), "登陆失败，请重新登陆")
	} else {
		// 初始化 onedrive 的连接，读取内容
		go onedrive.SetAutoRefresh()
		app.Response(c, http.StatusOK, e.SUCCESS, "登陆成功")
	}
}

// 注销登陆
func CancelLogin(c *gin.Context) {
	mg_auth.ClearCLient()
	
	app.Response(c, http.StatusOK, e.SUCCESS, "注销成功")
}

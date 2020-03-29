package mg_auth

import (
	"GOIndex/conf"
	"context"
	"errors"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"golang.org/x/oauth2"
	"net/http"
)

// 使用代码流的方式获取授权，文档
// https://docs.microsoft.com/zh-cn/onedrive/developer/rest-api/getting-started/graph-oauth?view=odsp-graph-online#code-flow

// 获取授权代码
// https://login.microsoftonline.com/common/oauth2/authorize?response_type=code&client_id=${client_id}&redirect_uri=${redirect_uri}

var oauthConfig oauth2.Config
var oauthStateString string
var client *http.Client

func SetUserInfo(user conf.UserInfo) {
	oauthConfig = oauth2.Config{
		Endpoint: oauth2.Endpoint{
			AuthURL:  "https://login.microsoftonline.com/common/oauth2/v2.0/authorize",
			TokenURL: "https://login.microsoftonline.com/common/oauth2/v2.0/token",
		},
		Scopes:       []string{"offline_access", "files.readwrite.all"},
		ClientID:     user.ClientID,
		ClientSecret: user.ClientSecret,
		RedirectURL:  user.RedirectURL,
	}
	oauthStateString = user.State
	client = nil
}

// redirect to microsoft login
func RedirectLoginMG(c *gin.Context) {
	url := oauthConfig.AuthCodeURL(oauthStateString)
	c.Redirect(http.StatusFound, url)
}

type ReceiveCode struct {
	Code         string `binding:"required"`
	SessionState string `binding:"required"`
	State        string `binding:"required"`
}

// receive code ,get access_token
func GetAccessToken(code ReceiveCode) error {
	ctx := context.Background()
	if code.State != oauthStateString {
		return errors.New("state 字符串与设置的不一致，请检查设置")
	}
	// 获取 AccessToken
	tok, err := oauthConfig.Exchange(ctx, code.Code)
	if err != nil {
		log.WithFields(log.Fields{
			"token": tok,
			"error": err,
		}).Info("获取 AccessToken 错误")
		return errors.New("获取 AccessToken 错误")
	}
	// 如果登陆成功返回成功，前端去跳转
	client = oauthConfig.Client(ctx, tok)
	log.WithField("token", tok).Info("获取 AccessToken 成功")
	return nil
}

func GetClient() *http.Client {
	return client
}

func ClearCLient() {
	client = nil
}

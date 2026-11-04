package auth

import (
	"context"
	"errors"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"golang.org/x/oauth2"

	"gonelist/conf"
)

// 使用代码流的方式获取授权，文档
// https://docs.microsoft.com/zh-cn/onedrive/developer/rest-api/getting-started/graph-oauth?view=odsp-graph-online#code-flow

// 获取授权代码
// https://login.microsoftonline.com/common/oauth2/authorize?response_type=code&client_id=${client_id}&redirect_uri=${redirect_uri}

// 跟 oauth2 有关的内容
var (
	clientID         string
	clientSecret     string
	oauthConfig      Config
	oauthStateString string
	client           *http.Client
	CacheOne         sync.Once
)

func SetOnedriveInfo(conf *conf.AllSet, initFunc func()) {
	user := conf.Onedrive
	clientID = user.ClientID
	clientSecret = user.ClientSecret

	endPoint := user.RemoteConf.EndPoint
	// 只申请读权限，避免应用程序进行修改，但使用 config.yml 给的默认 id 还是不太安全
	var scopes = []string{"offline_access", "files.read", "https://graph.microsoft.com/Sites.Read.All"}
	// 如果允许上传，则申请读写权限
	if conf.Admin.EnableWrite {
		scopes = append(scopes, "https://graph.microsoft.com/Files.ReadWrite.All")
	}

	// 初始化 oauth 的 config
	oauthConfig = Config{
		Config: &oauth2.Config{
			Endpoint:     endPoint,
			Scopes:       scopes,
			ClientID:     clientID,
			ClientSecret: clientSecret,
			RedirectURL:  user.RedirectURL,
		},
		Storage: &FileStorage{Path: user.TokenPath},
	}
	oauthStateString = user.State
	ctx := context.WithValue(context.Background(), oauth2.HTTPClient, &http.Client{Timeout: 60 * time.Second})
	tok, err := oauthConfig.Storage.GetToken()
	if err == nil {
		client = oauthConfig.Client(ctx, tok)
		log.Info("从文件读取 OAuth 令牌成功")
		// 初始化 onedrive 的内容
		initFunc()
		return
	}
	client = nil
}

// redirect to microsoft login
func RedirectLoginMG(c *gin.Context) {
	url := oauthConfig.AuthCodeURL(oauthStateString)
	c.Redirect(http.StatusFound, url)
}

type ReceiveCode struct {
	Code string `binding:"required"`
	//SessionState string `binding:"required"`
	State string `binding:"required"`
}

// receive code ,get access_token
func GetAccessToken(code ReceiveCode) error {
	ctx := context.WithValue(context.Background(), oauth2.HTTPClient, &http.Client{Timeout: 60 * time.Second})
	if code.State != oauthStateString {
		return errors.New("state 字符串与设置的不一致，请检查设置")
	}

	tok, err := GetToken(ctx, oauthConfig, code.Code)
	if err != nil {
		log.Warn("获取 AccessToken 错误")
		return errors.New("获取 AccessToken 错误")
	}
	// 如果登陆成功返回成功，前端去跳转
	client = oauthConfig.Client(ctx, tok)
	log.Info("获取 AccessToken 成功")
	return nil
}

func GetClient() *http.Client {
	return client
}

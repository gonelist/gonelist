package onedrive

import (
	"context"
	"errors"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"golang.org/x/oauth2"
	"gonelist/conf"
	"net/http"
	"sync"
)

// 使用代码流的方式获取授权，文档
// https://docs.microsoft.com/zh-cn/onedrive/developer/rest-api/getting-started/graph-oauth?view=odsp-graph-online#code-flow

// 获取授权代码
// https://login.microsoftonline.com/common/oauth2/authorize?response_type=code&client_id=${client_id}&redirect_uri=${redirect_uri}

var oauthConfig Config
var oauthStateString string
var client *http.Client
var cacheGoOnce sync.Once

func SetUserInfo(user *conf.UserSetting) {
	var endPoint oauth2.Endpoint
	// 设置 ChinaCloud 相关
	if user.ChinaCloud.Enable == true {
		endPoint = oauth2.Endpoint{
			AuthURL:  "https://login.chinacloudapi.cn/common/oauth2/v2.0/authorize",
			TokenURL: "https://login.chinacloudapi.cn/common/oauth2/v2.0/token",
		}
	} else {
		endPoint = oauth2.Endpoint{
			AuthURL:  "https://login.microsoftonline.com/common/oauth2/v2.0/authorize",
			TokenURL: "https://login.microsoftonline.com/common/oauth2/v2.0/token",
		}
	}
	SetROOTUrl(conf.UserSet.ChinaCloud.Enable)

	// 初始化 oauth 的 config
	oauthConfig = Config{
		Config: &oauth2.Config{
			Endpoint:     endPoint,
			Scopes:       []string{"offline_access", "files.read"}, // 只申请读权限，避免应用程序进行修改，但使用 config.json 给的默认 id 还是不太安全
			ClientID:     user.ClientID,
			ClientSecret: user.ClientSecret,
			RedirectURL:  user.RedirectURL,
		},
		Storage: &FileStorage{Path: user.TokenPath},
	}
	oauthStateString = user.State
	ctx := context.Background()
	tok, err := oauthConfig.Storage.GetToken()
	if err == nil {
		client = oauthConfig.Client(ctx, tok)
		log.WithField("refresh_token", tok.RefreshToken).Info("从文件读取refresh_token成功")
		if _, err := GetAllFiles(); err != nil {
			log.Fatal(err)
		}
		SetLogin(true)
		// 如果首页有 README.md 则下载到本地
		DownloadREADME()
		cacheGoOnce.Do(func() {
			go SetAutoRefresh()
		})

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
	ctx := context.Background()
	if code.State != oauthStateString {
		return errors.New("state 字符串与设置的不一致，请检查设置")
	}

	tok, err := GetToken(ctx, oauthConfig, code.Code)
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

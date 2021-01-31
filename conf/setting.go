package conf

import (
	"encoding/json"
	"errors"
	"fmt"
	log "github.com/sirupsen/logrus"
	"gonelist/pkg/file"
	"os"
	"path"
	"strings"
	"time"
)

// 服务器设置
type Server struct {
	Port         int `json:"port"`
	RefreshTime  int `json:"refresh_time"` //单位为分钟
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	DistPATH     string `json:"dist_path"` // 静态文件目录
	BindGlobal   bool   `json:"bind_global"`
	SiteUrl      string `json:"site_url"`   // 网站网址，如 https://gonelist.cugxuan.cn
	FolderSub    string `json:"folder_sub"` // onedrive 的子文件夹
	Gzip         bool   `json:"gzip"`       // 是否打开 Gzip 加速
}

var defaultServerSetting = &Server{
	Port:         8000,
	RefreshTime:  10,
	ReadTimeout:  60,
	WriteTimeout: 60,
	DistPATH:     "dist",
	BindGlobal:   true,
	SiteUrl:      "https://gonelist.cugxuan.cn",
	FolderSub:    "/",
	Gzip:         true,
}

type ChinaCloud struct {
	Enable       bool   `json:"enable"`
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
}

var defaultChinaCloudSetting = &ChinaCloud{
	Enable:       false,
	ClientID:     "",
	ClientSecret: "",
}

// 用户信息设置
type UserSetting struct {
	// 获取授权代码
	ResponseType string `json:"-"` // 值为 code
	ClientID     string `json:"client_id"`
	RedirectURL  string `json:"redirect_url"`
	State        string `json:"state"` // 用户设置的标识
	// 获取 access_token
	ClientSecret           string `json:"client_secret"`
	Code                   string `json:"-"`                        // 服务器收到的中间内容
	GrantType              string `json:"-"`                        // 值为 authorization_code
	Scope                  string `json:"-"`                        // 值为 offline_access files.readwrite.all
	AccessToken            string `json:"-"`                        // 令牌
	RefreshToken           string `json:"-"`                        // 刷新令牌
	TokenPath              string `json:"token_path"`               // token 文件位置
	DownloadRedirectPrefix string `json:"download_redirect_prefix"` // 下载重定向前缀
	// 世纪互联
	ChinaCloud *ChinaCloud `json:"china_cloud"`
	// 用户设置
	Server *Server `json:"server"`
	// 目录密码
	PassList []*Pass `json:"pass_list"`
}

var UserSet = &UserSetting{}

func LoadUserConfig(configPath string) error {
	var content []byte
	var err error

	if len(configPath) == 0 {
		return errors.New("配置文件名不能为空")
	}
	envValue := os.Getenv("CONF_PATH")

	if envValue != "" {
		configPath = envValue
	}

	log.Infof("当前使用的配置文件为:%s", configPath)

	content, _ = file.ReadFromFile(configPath)
	err = json.Unmarshal(content, &UserSet)
	if err != nil {
		return fmt.Errorf("导入用户配置出现错误: %w", err)
	}
	// Server 的设置
	if UserSet.Server == nil {
		UserSet.Server = defaultServerSetting
	}
	// PassList 设置
	if UserSet.PassList == nil {
		UserSet.PassList = defaultPassListSetting
	}
	if UserSet.Server.FolderSub == "" {
		UserSet.Server.FolderSub = "/"
	}
	// ChinaCloud 设置
	if UserSet.ChinaCloud == nil {
		UserSet.ChinaCloud = defaultChinaCloudSetting
	}
	// TokenPath 不为 ""，token 保存在用户设置的目录
	// 否则 token 将保存在用户 config.json 所在的目录
	if UserSet.TokenPath == "" {
		UserSet.TokenPath = GetTokenPath(configPath)
	} else {
		//用户一般写目录，此处转成文件
		if !strings.HasSuffix(UserSet.TokenPath, ".token") {
			UserSet.TokenPath = path.Join(UserSet.TokenPath, ".token")
		}
	}
	log.Info("成功导入用户配置")
	return nil
}

// return the refresh time from the settings
func GetRefreshTime() time.Duration {
	return time.Duration(UserSet.Server.RefreshTime) * time.Minute
}

func GetBindAddr(bind bool, port int) string {
	var prefix string
	if bind == false {
		prefix = "127.0.0.1"
	}
	return fmt.Sprintf("%s:%d", prefix, port)
}

func GetDistPATH() string {
	return UserSet.Server.DistPATH
}

func GetTokenPath(configPath string) string {
	lastIndex := strings.LastIndex(configPath, string(os.PathSeparator))
	return configPath[:lastIndex+1] + ".token"
}

type Pass struct {
	Path string `json:"path"`
	Pass string `json:"pass"`
}

var defaultPassListSetting = []*Pass{
	{
		Path: "",
		Pass: "",
	},
}

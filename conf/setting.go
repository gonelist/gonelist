package conf

import (
	"GOIndex/pkg/file"
	"encoding/json"
	"errors"
	"fmt"
	log "github.com/sirupsen/logrus"
	"time"
)

// 服务器设置
type Server struct {
	WebPort      int    `json:"web_port"`
	BackPort     int    `json:"back_port"`
	RefreshTime  int    `json:"refresh_time"` //单位为分钟
	ReadTimeout  time.Duration
	WriteTimeout time.Duration

	WebBindGlobal  bool   `json:"web_bind_global"`
	BackBindGlobal bool   `json:"back_bind_global"`
	SiteUrl        string `json:"site_url"`   // 网站网址，如 https://goindex.cugxuan.cn
	FolderSub      string `json:"folder_sub"` // onedrive 的子文件夹
}

var defaultServerSetting = &Server{
	WebPort:      8001,
	BackPort:     8000,
	RefreshTime:  10,
	ReadTimeout:  60,
	WriteTimeout: 60,

	WebBindGlobal:  true,
	BackBindGlobal: false,
	SiteUrl:        "https://goindex.cugxuan.cn",
	FolderSub:      "/",
}

// 用户信息设置
type UserSetting struct {
	// 获取授权代码
	ResponseType string `json:"-"`             // 值为 code
	ClientID     string `json:"client_id"`
	RedirectURL  string `json:"redirect_url"`
	State        string `json:"state"`         // 用户设置的标识
	// 获取 access_token
	ClientSecret string `json:"client_secret"`
	Code         string `json:"-"`             // 服务器收到的中间内容
	GrantType    string `json:"-"`             // 值为 authorization_code
	Scope        string `json:"-"`             // 值为 offline_access files.readwrite.all
	AccessToken  string `json:"-"`             // 令牌
	RefreshToken string `json:"-"`             // 刷新令牌
	// 用户设置
	Server *Server `json:"server"`
}

var UserSet UserSetting

func LoadUserConfig(filePath string) error {
	var content string
	var err error

	if len(filePath) == 0 {
		return errors.New("配置文件名不能为空")
	}

	log.Infof("当前使用的配置文件为:%s", filePath)

	content = file.ReadFromFile(filePath)
	err = json.Unmarshal([]byte(content), &UserSet)
	if err != nil {
		return fmt.Errorf("导入用户配置出现错误: %w", err)
	}
	if UserSet.Server == nil {
		UserSet.Server = defaultServerSetting
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

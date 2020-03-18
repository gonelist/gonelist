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
	RunMode      string  `json:"run_mode"`
	HttpPort     int     `json:"http_port"`
	refreshTime  int  `json:"refresh_time"` //单位为分钟
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}


var defaultServerSetting = &Server{
	RunMode:      "release",
	HttpPort:     8000,
	refreshTime:  10,
	ReadTimeout:  60,
	WriteTimeout: 60,
}

// 用户信息设置
type UserInfo struct {
	// 获取授权代码
	ResponseType string `json:"-"` // 值为 code
	ClientID     string `json:"client_id"`
	RedirectURL  string `json:"redirect_url"`
	State        string `json:"state"` // 用户设置的标识
	// 获取 access_token
	ClientSecret string `json:"client_secret"`
	Code         string `json:"-"` // 服务器收到的中间内容
	GrantType    string `json:"-"` // 值为 authorization_code
	Scope        string `json:"-"` // 值为 offline_access files.readwrite.all
	AccessToken  string `json:"-"` // 令牌
	RefreshToken string `json:"refresh_token"` //刷新令牌
	Server      *Server `json:"server"`
}

var UserSetting UserInfo

func LoadUserConfig(filePath string) error {
	var content string
	var err error

	if len(filePath) == 0 {
		return errors.New("配置文件名不能为空")
	}

	log.Infof("当前使用的配置文件为:%s", filePath)

	content = file.ReadFromFile(filePath)
	err = json.Unmarshal([]byte(content), &UserSetting)
	if err != nil {
		return fmt.Errorf("导入用户配置出现错误: %w", err)
	}
	if UserSetting.Server == nil {
		UserSetting.Server = defaultServerSetting
	}
	log.Info("成功导入用户配置")
	return nil
}

// return the refresh time from the settings
func GetRefreshTime() time.Duration {
	return time.Duration(UserSetting.Server.refreshTime) * time.Minute
}
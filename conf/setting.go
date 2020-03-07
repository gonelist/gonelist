package conf

import (
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	"GOIndex/pkg/file"
	"time"
)

// 服务器设置
type Server struct {
	RunMode      string
	HttpPort     int
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}

var ServerSetting = &Server{
	RunMode:      "release",
	HttpPort:     8000,
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
	RefreshToken string `json:"-"` //刷新令牌
}

var UserSetting UserInfo

func LoadUserConfig(filePath string) error {
	var content string
	var err error

	fmt.Println("导入", filePath, "配置")
	content = file.ReadFromFile(filePath)
	err = json.Unmarshal([]byte(content), &UserSetting)
	if err != nil {
		fmt.Println("导入用户配置出现错误")
		log.Warn(err)
		return err
	}
	fmt.Println("成功导入用户配置")
	return nil
}

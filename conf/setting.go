package conf

import (
	"errors"
	"fmt"
	log "github.com/sirupsen/logrus"
	"gonelist/pkg/file"
	"gopkg.in/yaml.v2"
	"os"
	"path"
	"strings"
	"time"
)

// 服务器设置
type Server struct {
	Port         int           `json:"port" yaml:"port"`
	ReadTimeout  time.Duration `yaml:"read_timeout"`
	WriteTimeout time.Duration `yaml:"write_timeout"`
	DistPATH     string        `json:"dist_path" yaml:"dist_path"` // 静态文件目录
	BindGlobal   bool          `json:"bind_global" yaml:"bind_global"`
	SiteUrl      string        `json:"site_url" yaml:"site_url"` // 网站网址，如 https://gonelist.cugxuan.cn
	Gzip         bool          `json:"gzip" yaml:"gzip"`         // 是否打开 Gzip 加速
}

var defaultServerSetting = &Server{
	Port:         8000,
	ReadTimeout:  60,
	WriteTimeout: 60,
	DistPATH:     "./dist/",
	BindGlobal:   true,
	SiteUrl:      "https://gonelist.cugxuan.cn",
	Gzip:         true,
}

type Onedrive struct {
	// Remote to load RemoteConf
	Remote     string `json:"remote" yaml:"remote"`
	RemoteConf Remote `json:"-" yaml:"-"`
	Level      int    `json:"level" yaml:"level"`
	// 自动刷新时间，单位为分钟
	RefreshTime int `json:"refresh_time" yaml:"refresh_time"`
	// 获取授权代码
	ResponseType string `json:"-" yaml:"-"` // 值为 code
	ClientID     string `json:"client_id" yaml:"client_id"`
	RedirectURL  string `json:"redirect_url" yaml:"redirect_url"`
	State        string `json:"state" yaml:"state"` // 用户设置的标识
	// 获取 access_token
	ClientSecret string `json:"client_secret" yaml:"client_secret"`
	Code         string `json:"-" yaml:"-"`                   // 服务器收到的中间内容
	GrantType    string `json:"-" yaml:"-"`                   // 值为 authorization_code
	Scope        string `json:"-" yaml:"-"`                   // 值为 offline_access files.readwrite.all
	AccessToken  string `json:"-" yaml:"-"`                   // 令牌
	RefreshToken string `json:"-" yaml:"-"`                   // 刷新令牌
	TokenPath    string `json:"token_path" yaml:"token_path"` // token 文件位置
	// 其他设置
	FolderSub              string `json:"folder_sub" yaml:"folder_sub"`                             // onedrive 的子文件夹
	DownloadRedirectPrefix string `json:"download_redirect_prefix" yaml:"download_redirect_prefix"` // 下载重定向前缀
	// 目录密码
	PassList []*Pass `json:"pass_list" yaml:"pass_list"`
}

// 用户信息设置
type AllSet struct {
	Name    string `json:"name" yaml:"name"`
	Version string `json:"version" yaml:"version"`
	// Server 配置，用于定义服务的特性
	Server *Server `json:"server" yaml:"server"`
	// 网盘挂载类型
	ListType string `json:"list_type" yaml:"list_type"`
	// Onedrive
	Onedrive *Onedrive `json:"onedrive" yaml:"onedrive"`
}

var UserSet = &AllSet{}

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

	if content, err = file.ReadFromFile(configPath); err != nil {
		return fmt.Errorf("read config err,path: %s", configPath)
	}
	err = yaml.Unmarshal(content, &UserSet)
	if err != nil {
		return fmt.Errorf("导入用户配置出现错误: %w", err)
	}
	// Server 的设置
	if UserSet.Server == nil {
		return fmt.Errorf("Server 设置读取出现错误")
	}
	switch UserSet.ListType {
	case "onedrive":
		// 处理 Remote 地址
		switch UserSet.Onedrive.Remote {
		case "onedrive":
			UserSet.Onedrive.RemoteConf = OneDrive
		case "chinacloud":
			UserSet.Onedrive.RemoteConf = ChinaCloud
		}
		// PassList 设置
		if UserSet.Onedrive.FolderSub == "" {
			UserSet.Onedrive.FolderSub = "/"
		}
		if UserSet.Onedrive.PassList == nil {
			UserSet.Onedrive.PassList = defaultPassListSetting
		}
		// TokenPath 不为 ""，token 保存在用户设置的目录
		// 否则 token 将保存在用户 config.yml 所在的目录
		if UserSet.Onedrive.TokenPath == "" {
			UserSet.Onedrive.TokenPath = GetTokenPath(configPath)
		} else {
			//用户一般写目录，此处转成文件
			if !strings.HasSuffix(UserSet.Onedrive.TokenPath, ".token") {
				UserSet.Onedrive.TokenPath = path.Join(UserSet.Onedrive.TokenPath, ".token")
			}
		}
	default:
		return fmt.Errorf("不支持的网盘挂载类型")
	}

	log.Info("成功导入用户配置")
	return nil
}

// return the refresh time from the settings
func GetRefreshTime() time.Duration {
	return time.Duration(UserSet.Onedrive.RefreshTime) * time.Minute
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

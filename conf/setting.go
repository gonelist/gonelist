package conf

import (
	"crypto/sha256"
	"errors"
	"fmt"
	"os"
	"path"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"

	"gonelist/pkg/file"
)

// 服务器设置
type Server struct {
	ReadTimeout  time.Duration `yaml:"read_timeout"`
	WriteTimeout time.Duration `yaml:"write_timeout"`
	BindGlobal   bool          `json:"bind_global" yaml:"bind_global"` // 是否绑定到0.0.0.0
	DistPATH     string        `json:"dist_path" yaml:"dist_path"`     // 静态文件目录
	Gzip         bool          `json:"gzip" yaml:"gzip"`               // 是否打开 Gzip 加速
	Port         int           `json:"port" yaml:"port"`               // 绑定端口
	SiteUrl      string        `json:"site_url" yaml:"site_url"`       // 网站网址，如 https://gonelist.cugxuan.cn
}

var defaultServerSetting = &Server{
	ReadTimeout:  60,
	WriteTimeout: 60,
	BindGlobal:   true,
	DistPATH:     "./dist/",
	Port:         8000,
	Gzip:         true,
	SiteUrl:      "https://gonelist.cugxuan.cn",
}

// WebDav WebDav配置
type WebDav struct {
	Enable   bool   `json:"enable" yaml:"enable"`     //是否开启webdav
	Account  string `json:"account" yaml:"account"`   // webdav账号
	Password string `json:"password" yaml:"password"` // webdav密码
}

type Admin struct {
	EnableWrite     bool   `json:"enable_write" yaml:"enable_write"`           // 是否允许客户端写入文件到onedrive服务端，写入包括创建文件夹，上传文件，删除文件
	Secret          string `json:"secret" yaml:"secret"`                       // 写入权限的secret，前端升级权限时需要，建议更改默认secret
	UploadSliceSize int    `json:"upload_slice_size" yaml:"upload_slice_size"` // 大文件分片上传时得分片大小，默认为32MB,数字为1表示320kb
}

// Local
// 本地目录挂载配置
type Local struct {
	Enable bool   `json:"enable" yaml:"enable"` // 是否开启本地挂载
	Name   string `json:"name" yaml:"name"`     // 本地挂载的目录在gonelist中显示的名称
	Path   string `json:"path" yaml:"path"`     // 本地挂载目录的路径
}

type Onedrive struct {
	// Remote to load RemoteConf
	Remote     string `json:"remote" yaml:"remote"`
	RemoteConf Remote `json:"-" yaml:"-"`
	// 刷新模式
	Model       string `json:"model" yaml:"model"`               // 刷新模式
	Level       int    `json:"level" yaml:"level"`               // 刷新层级
	RefreshTime int    `json:"refresh_time" yaml:"refresh_time"` // 自动刷新时间，单位为分钟
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
	Name      string `json:"name" yaml:"name"`
	PageTitle string `json:"page_title" yaml:"page_title"`
	Version   string `json:"version" yaml:"version"`
	// Server 配置，用于定义服务的特性
	Server *Server `json:"server" yaml:"server"`
	// 网盘挂载类型
	ListType string `json:"list_type" yaml:"list_type"`
	// Onedrive
	Onedrive *Onedrive `json:"onedrive" yaml:"onedrive"`
	// 权限管理
	Admin *Admin `json:"admin" yaml:"admin"`

	// 本地目录挂载
	Local *Local `json:"local" yaml:"local"`

	WebDav *WebDav `json:"web_dav" yaml:"web_dav"`
}

var UserSet = &AllSet{}

func LoadUserConfig(configPath string) error {
	var (
		content []byte
		err     error
	)

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
	loaded := &AllSet{}
	err = yaml.Unmarshal(content, loaded)
	if err != nil {
		return errors.New("导入用户配置出现错误，请检查 YAML 格式和字段类型")
	}
	if loaded.Server == nil || loaded.Onedrive == nil {
		return errors.New("server and onedrive configuration are required")
	}
	if loaded.Admin == nil {
		loaded.Admin = &Admin{}
	}
	if loaded.Local == nil {
		loaded.Local = &Local{}
	}
	if loaded.WebDav == nil {
		loaded.WebDav = &WebDav{}
	}
	if loaded.Admin.EnableWrite && (len(loaded.Admin.Secret) < 16 || loaded.Admin.Secret == "CHANGE_ME_ADMIN_SECRET") {
		return errors.New("admin.secret must be a private random secret of at least 16 characters when writes are enabled")
	}
	if loaded.WebDav.Enable && (loaded.WebDav.Account == "" || len(loaded.WebDav.Password) < 12 || loaded.WebDav.Password == loaded.WebDav.Account) {
		return errors.New("WebDAV requires an account and a private password of at least 12 characters different from the account")
	}
	if loaded.Local.Enable && (loaded.Local.Path == "" || loaded.Local.Name == "" || loaded.Local.Name == "." || loaded.Local.Name == ".." || strings.ContainsAny(loaded.Local.Name, "/\\\x00")) {
		return errors.New("local mount requires a path and a single directory name")
	}
	if loaded.Onedrive.Remote != "local" && (loaded.Onedrive.ClientID == "" || loaded.Onedrive.ClientSecret == "") {
		return errors.New("configure your own onedrive.client_id and client_secret before starting")
	}
	// These fingerprints identify the shared credentials published in old templates.
	digest := fmt.Sprintf("%x", sha256.Sum256([]byte(loaded.Onedrive.ClientSecret)))
	if digest == "821093316eb629192fa9972c25f752ec65a6d3f8b4b703c66bde465db0d8a5a1" || digest == "ba8fb5201df33fa76068f8424afd8506c2bce91be0c66954d38aa89c55c12d40" {
		return errors.New("the legacy shared client secret is public; create your own application and rotate its credentials")
	}
	switch loaded.ListType {
	case "onedrive":
		// 处理 Remote 地址
		switch loaded.Onedrive.Remote {
		case "onedrive":
			loaded.Onedrive.RemoteConf = OneDrive
		case "chinacloud":
			loaded.Onedrive.RemoteConf = ChinaCloud
		case "local":
		default:
			return errors.New("unsupported onedrive.remote")
		}
		// PassList 设置
		if loaded.Onedrive.FolderSub == "" {
			loaded.Onedrive.FolderSub = "/"
		}
		if loaded.Onedrive.PassList == nil {
			loaded.Onedrive.PassList = defaultPassListSetting
		}
		// TokenPath 不为 ""，token 保存在用户设置的目录
		// 否则 token 将保存在用户 config.yml 所在的目录
		if loaded.Onedrive.TokenPath == "" {
			loaded.Onedrive.TokenPath = GetTokenPath(configPath)
		} else {
			// 用户一般写目录，此处转成文件
			if !strings.HasSuffix(loaded.Onedrive.TokenPath, ".token") {
				loaded.Onedrive.TokenPath = path.Join(loaded.Onedrive.TokenPath, ".token")
			}
		}
		// 设置大文件上传的分片
		if loaded.Admin.UploadSliceSize == 0 {
			loaded.Admin.UploadSliceSize = 100
		}
	default:
		return fmt.Errorf("不支持的网盘挂载类型")
	}

	log.Infof("gonelist 监听端口:%v,成功导入用户配置", loaded.Server.Port)
	UserSet = loaded
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

package routers

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-contrib/gzip"
	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
	cors "github.com/rs/cors/wrapper/gin"
	log "github.com/sirupsen/logrus"

	"gonelist/conf"
	"gonelist/middleware"
	"gonelist/routers/api"
	"gonelist/routers/webdav"
)

func InitRouter() *gin.Engine {

	r := gin.New()
	// 对于 router 中多个斜杠（slash）进行兼容
	// 如 /ping //ping 是同一个接口
	r.RemoveExtraSlash = true
	// Do not log queries: legacy clients send admin secrets and OAuth codes there.
	r.Use(gin.LoggerWithFormatter(func(p gin.LogFormatterParams) string {
		return fmt.Sprintf("%s %s %d %s\n", p.TimeStamp.Format(time.RFC3339), p.Method, p.StatusCode, p.Request.URL.EscapedPath())
	}))
	r.Use(func(c *gin.Context) {
		defer func() {
			if recover() != nil {
				log.Error("request panic recovered")
				c.AbortWithStatus(http.StatusInternalServerError)
			}
		}()
		c.Next()
	})
	r.Use(cors.New(cors.Options{
		AllowedHeaders: []string{"pass"}, // 允许 header
	}))

	if conf.UserSet.Server.Gzip {
		r.Use(gzip.Gzip(gzip.DefaultCompression))
	}

	// Gin may clean dot segments before assigning route parameters. Reject the
	// original ambiguous URL before any handler can use the cleaned parameters.
	r.Use(func(c *gin.Context) {
		for _, part := range strings.Split(c.Request.URL.Path, "/") {
			if part == ".." || strings.ContainsAny(part, "\\\x00") {
				c.AbortWithStatus(http.StatusBadRequest)
				return
			}
		}
		c.Next()
	})

	r.Use(static.Serve("/", static.LocalFile(conf.GetDistPATH(), false)))

	// 测试接口
	r.GET("/testapi", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "hello world",
		})
	})

	// 检测webdav是否开启
	if conf.UserSet.WebDav.Enable {
		// 初始化webdav处理器
		handler := webdav.DavInit()
		r.Use(func(ctx *gin.Context) {
			// 挂载目录
			if ctx.Request.URL.Path == "/webdav" || strings.HasPrefix(ctx.Request.URL.Path, "/webdav/") {
				// 获取用户名/密码
				username, password, ok := ctx.Request.BasicAuth()
				if !ok {
					ctx.Writer.Header().Set("WWW-Authenticate", `Basic realm="Restricted"`)
					ctx.AbortWithStatus(http.StatusUnauthorized)
					return
				}
				// 验证用户名/密码
				if username != conf.UserSet.WebDav.Account || password != conf.UserSet.WebDav.Password {
					ctx.Header("WWW-Authenticate", `Basic realm="Restricted"`)
					ctx.AbortWithStatus(http.StatusUnauthorized)
					return
				}
				if !conf.UserSet.Admin.EnableWrite && ctx.Request.Method != http.MethodGet && ctx.Request.Method != http.MethodHead && ctx.Request.Method != "PROPFIND" && ctx.Request.Method != http.MethodOptions {
					ctx.AbortWithStatus(http.StatusForbidden)
					return
				}
				// 修改path
				ctx.Request.URL.Path = strings.TrimPrefix(ctx.Request.URL.Path, "/webdav")
				if !strings.HasPrefix(ctx.Request.URL.Path, "/") {
					ctx.Request.URL.Path = "/" + ctx.Request.URL.Path
				}
				handler.ServeHTTP(ctx.Writer, ctx.Request)

				ctx.Abort()
				return
			}
		})
	}

	r.GET("/info", api.Info)

	r.GET("/login", api.Login)
	r.GET("/loginmg", api.LoginMG)
	r.GET("/auth", api.GetCode)
	//r.GET("/cancelLogin", api.CancelLogin)

	r.GET("/update_permission", api.UpdatePermission())

	// 直接下载接口
	root := r.Group("/")
	root.Use(middleware.CheckLogin())
	{
		root.GET("/d/*path", api.LocalDownload(), api.Download)
		root.GET("/README", middleware.CheckFolderPass(), api.GetREADME)
	}

	onedrive := r.Group("/onedrive")
	onedrive.Use(middleware.CheckLogin())
	{
		onedrive.GET("/mkdir", middleware.CheckSecret(), api.LocalMkdir(), api.MkDir())
		// 上传文件，仅支持大小4MB
		onedrive.POST("/upload", middleware.CheckSecret(), api.LocalUpload(), api.Upload())
		// 删除文件
		onedrive.GET("/delete_file", middleware.CheckSecret(), api.LocalDelete(), api.DeleteFile())
		// 主动获取所有文件，返回整个树的目录
		onedrive.GET("/getallfiles", api.MGGetFileTree)
		// 根据路径获取对应数据
		onedrive.GET("/getpath", middleware.CheckFolderPass(), api.CacheGetPath)
		onedrive.GET("/search", api.Search)
	}

	return r
}

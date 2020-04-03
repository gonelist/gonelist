package routers

import (
	"GOIndex/api"
	"GOIndex/conf"
	"GOIndex/middleware"
	"github.com/gin-gonic/gin"
	cors "github.com/rs/cors/wrapper/gin"
	"net/http"
)

func InitRouter() *gin.Engine {
	r := gin.New()

	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	r.Use(cors.Default())

	// 测试接口
	r.GET("/testapi", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "hello world",
		})
	})

	r.GET("/login", api.Login)
	r.GET("/loginmg", api.LoginMG)
	r.GET("/auth", api.GetCode)
	r.GET("/cancelLogin", api.CancelLogin)
	onedrive := r.Group("/onedrive")
	// 中间件判断是否已经登录
	onedrive.Use(middleware.CheckLogin())
	{
		// 主动获取所有文件，返回整个树的目录
		onedrive.GET("/getallfiles", api.MGGetFileTree)
		// 根据路径获取对应数据
		onedrive.GET("/getpath", api.CacheGetPath)
	}
	// 前端内容
	r.StaticFS("/"+conf.UserSetting.SubPath, http.Dir("dist"))

	return r
}

package main

import (
	"fmt"
	"GOIndex/conf"
	"GOIndex/mg_auth"
	"GOIndex/routers"
)

func main() {
	// 加载用户配置
	conf.LoadUserConfig("conf/config.json")
	mg_auth.SetUserInfo(conf.UserSetting)
	fmt.Println(conf.UserSetting)
	// 启动服务器
	r := routers.InitRouter()

	r.Run(fmt.Sprintf(":%d", conf.ServerSetting.HttpPort))
}

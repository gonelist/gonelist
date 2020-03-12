package GOIndex

import (
	"GOIndex/conf"
	"GOIndex/mg_auth"
	"GOIndex/routers"
	"flag"
	"fmt"
	log "github.com/sirupsen/logrus"

)

func main() {

	confPath := flag.String("conf", "conf/config.json", "指定配置文件路径")
	flag.Parse()

	log.SetFormatter(&log.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02 15:04:05",
	})
	// 加载用户配置
	if err := conf.LoadUserConfig(*confPath); err != nil {
		log.Fatal(err)
	}

	mg_auth.SetUserInfo(conf.UserSetting)

	log.Printf("%v", conf.UserSetting)
	// 启动服务器
	r := routers.InitRouter()

	panic(r.Run(fmt.Sprintf(":%d", conf.ServerSetting.HttpPort)))
}

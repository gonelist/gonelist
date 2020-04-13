package main

import (
	"flag"
	log "github.com/sirupsen/logrus"
	"golang.org/x/sync/errgroup"
	"gonelist/conf"
	"gonelist/mg_auth"
	"gonelist/routers"
	"net/http"
	"time"
)

var (
	g errgroup.Group
)

func main() {

	confPath := flag.String("conf", "config.json", "指定配置文件路径")
	flag.Parse()

	log.SetFormatter(&log.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02 15:04:05",
	})
	// 加载用户配置
	if err := conf.LoadUserConfig(*confPath); err != nil {
		log.Fatal(err)
	}
	mg_auth.SetUserInfo(conf.UserSet)

	// 处理端口绑定
	Addr := conf.GetBindAddr(conf.UserSet.Server.BindGlobal, conf.UserSet.Server.Port)

	// 启动服务器
	server := &http.Server{
		Addr:         Addr,
		Handler:      routers.InitRouter(),
		ReadTimeout:  60 * time.Second,
		WriteTimeout: 60 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	panic(server.ListenAndServe())
}

package main

import (
	"GOIndex/conf"
	"GOIndex/mg_auth"
	"GOIndex/routers"
	"flag"
	log "github.com/sirupsen/logrus"
	"golang.org/x/sync/errgroup"
	"net/http"
	"time"
)

var (
	g errgroup.Group
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
	mg_auth.SetUserInfo(conf.UserSet)

	// 处理端口绑定
	backAddr := conf.GetBindAddr(conf.UserSet.Server.BackBindGlobal, conf.UserSet.Server.BackPort)
	webAddr := conf.GetBindAddr(conf.UserSet.Server.WebBindGlobal, conf.UserSet.Server.WebPort)

	// 启动服务器
	serverBack := &http.Server{
		Addr:         backAddr,
		Handler:      routers.InitRouter(),
		ReadTimeout:  60 * time.Second,
		WriteTimeout: 60 * time.Second,
	}

	serverFront := &http.Server{
		Addr:         webAddr,
		Handler:      routers.InitWeb(),
		ReadTimeout:  60 * time.Second,
		WriteTimeout: 60 * time.Second,
	}

	g.Go(func() error {
		err := serverBack.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			log.Fatal(err)
		}
		return err
	})

	g.Go(func() error {
		err := serverFront.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			log.Fatal(err)
		}
		return err
	})

	if err := g.Wait(); err != nil {
		log.Fatal(err)
	}
}

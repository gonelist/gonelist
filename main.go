package main

import (
	"flag"
	"fmt"
	log "github.com/sirupsen/logrus"
	"gonelist/conf"
	"gonelist/onedrive"
	"gonelist/routers"
	"net/http"
	"runtime"
	"time"
)

func main() {

	confPath := flag.String("conf", "config.json", "指定配置文件路径")
	versionB := flag.Bool("version", false, "Show current version of gonelist.")
	debugB := flag.Bool("debug", false, "debug log level")
	flag.Parse()

	log.SetFormatter(&log.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02 15:04:05",
	})
	if *versionB {
		versionPrint()
		return
	}
	if *debugB {
		log.SetLevel(log.DebugLevel)
	}
	// 加载用户配置
	if err := conf.LoadUserConfig(*confPath); err != nil {
		log.Fatal(err)
	}

	onedrive.SetUserInfo(conf.UserSet)
	onedrive.SetROOTUrl(conf.UserSet.ChinaCloud.Enable)

	// 处理端口绑定
	Addr := conf.GetBindAddr(conf.UserSet.Server.BindGlobal, conf.UserSet.Server.Port)

	// 启动服务器
	server := &http.Server{
		Addr:           Addr,
		Handler:        routers.InitRouter(),
		ReadTimeout:    60 * time.Second,
		WriteTimeout:   60 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	panic(server.ListenAndServe())
}

var (
	Version      string
	gitCommit    string
	gitTreeState = ""                     // state of git tree, either "clean" or "dirty"
	buildDate    = "1970-01-01T00:00:00Z" // build date, output of $(date +'%Y-%m-%dT%H:%M:%S')
)

func versionPrint() {
	fmt.Printf(`Name: gonelist
Version: %s
CommitID: %s
GitTreeState: %s
BuildDate: %s
GoVersion: %s
Compiler: %s
Platform: %s/%s
`, Version, gitCommit, gitTreeState, buildDate, runtime.Version(), runtime.Compiler, runtime.GOOS, runtime.GOARCH)
}

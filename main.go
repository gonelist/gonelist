package main

import (
	"flag"
	"fmt"
	"net/http"
	"runtime"
	"time"

	log "github.com/sirupsen/logrus"

	"gonelist/conf"
	"gonelist/routers"
	"gonelist/service/onedrive"
)

func main() {
	confPath := flag.String("conf", "config.yml", "指定配置文件路径")
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
	// 加载 config.yml
	if err := conf.LoadUserConfig(*confPath); err != nil {
		log.Fatal(err)
	}
	// 处理用户定义的 passList
	onedrive.InitPass(conf.UserSet)

	// 设置 onedrive 的相关配置，如果有 .token 那么会在这儿进行处理初始化
	// 否则在端口绑定之后通过接口登陆之后初始化
	onedrive.SetOnedriveInfo(conf.UserSet)

	// 设置 version
	conf.UserSet.Version = Version
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

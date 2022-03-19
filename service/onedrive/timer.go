package onedrive

import (
	"io"
	"os"
	"time"

	log "github.com/sirupsen/logrus"

	"gonelist/conf"
)

// 定时刷新缓存
func SetAutoRefresh() {
	d := conf.GetRefreshTime()
	go timer(AutoRefresh, d)
}

// 刷新 onedrive 文件列表的缓存
func AutoRefresh() {
	RefreshOnedriveAll()
}
func setToken(token string) {
	err := os.WriteFile(".file_token", []byte(token), 0666)
	if err != nil {
		log.Errorln(err.Error())
		return
	}
}
func getToken() string {
	file, _ := os.OpenFile(".file_token", os.O_CREATE|os.O_RDWR, 0666)
	data, _ := io.ReadAll(file)
	return string(data)
}

// 定时刷新函数，等待上一轮刷新执行结束，再过 duration 分钟
func timer(AutoCallFunction func(), duration time.Duration) {
	for {
		select {
		case <-time.After(duration):
			if FileTree.IsLogin() == false {
				log.Info("停止刷新缓存")
				return
			} else {
				AutoCallFunction()
			}
		}
	}
}

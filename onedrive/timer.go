package onedrive

import (
	log "github.com/sirupsen/logrus"
	"gonelist/conf"
	"time"
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

func timer(AutoCallFunction func(), duration time.Duration) {
	ticker := time.NewTicker(duration)
	for {
		select {
		case <-ticker.C:
			if FileTree.IsLogin() == false {
				log.Info("停止刷新缓存")
				return
			} else {
				AutoCallFunction()
			}
		}
	}
}

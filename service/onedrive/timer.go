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

// 使用时间间隔，等待上一轮刷新玩之后再过 duration 分钟
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

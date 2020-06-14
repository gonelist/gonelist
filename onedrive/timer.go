package onedrive

import (
	log "github.com/sirupsen/logrus"
	"gonelist/conf"
	"time"
)

// 定时刷新缓存
func SetAutoRefresh() {
	timer(AutoRefresh)
}

func AutoRefresh() {
	log.Info("开始自动刷新所有缓存")
	if _, err := GetAllFiles(); err != nil { // 获取所有文件并且刷新树结构
		log.Error(err)
	}
	log.Debug(FileTree)
}

func timer(AutoRefresh func()) {
	ticker := time.NewTicker(conf.GetRefreshTime())
	for {
		select {
		case <-ticker.C:
			if IsLogin() == false {
				log.Info("停止刷新缓存")
				return
			} else {
				AutoRefresh()
			}
		}
	}
}

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
	log.Info("开始刷新文件缓存")
	if _, err := GetAllFiles(); err != nil { // 获取所有文件并且刷新树结构
		log.WithField("err", err).Error("刷新文件缓存遇到错误")
		return
	}
	log.Infof("结束刷新文件缓存")
	log.Debug(FileTree)
	log.Info("开始刷新 README 缓存")
	if err := RefreshREADME(); err != nil {
		log.WithField("err", err).Error("刷新 README 缓存遇到错误")
		return
	}
	log.Info("结束刷新 README 缓存")
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

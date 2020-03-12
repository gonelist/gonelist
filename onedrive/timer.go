package onedrive

import (
	log "github.com/sirupsen/logrus"
	"time"
)

// 定时刷新缓存
func SetAutoRefresh() {
	timer(AutoRefresh)
}

func AutoRefresh() {
	log.WithFields(log.Fields{
		"time": time.Now(),
	}).Info("自动刷新所有缓存")
	FileTree = GetAllFiles()
	log.Debug(FileTree)
}

func timer(timer func()) {
	ticker := time.NewTicker(1 * time.Minute)
	for {
		select {
		case <-ticker.C:
			timer()
		}
	}
}

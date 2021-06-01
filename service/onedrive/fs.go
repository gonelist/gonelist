package onedrive

import (
	"fmt"
	gocache "github.com/patrickmn/go-cache"
	log "github.com/sirupsen/logrus"
	"time"
)

// 设置缓存的默认时间为 2 天，每 2 天清空已经失效的缓存
var reCache = gocache.New(DefaultTime, DefaultTime)

// 在缓存中 key 的形式是 README_path
// eg. README_/, README_/exampleFolder
const (
	README      = "README_"
	DefaultTime = time.Hour * 24
	FS          = "FS_"
)


// 真实的下载 URL 通过 Cache 进行存储
func GetPathInCache(p string) ([]byte, error) {
	ans, ok := reCache.Get(FS + p)
	if !ok {
		log.WithFields(log.Fields{
			"path": p,
		}).Info("FS not in cache")
		return nil, fmt.Errorf("FS not in cache")
	}

	return ans.([]byte), nil
}

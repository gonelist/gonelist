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
	READEME     = "README_"
	DefaultTime = time.Hour * 24
)

// 刷新每个文件夹的 README
func RefreshREADME() error {
	// 获取根节点开始
	root := FileTree.GetRoot()
	if err := GetCurrentAndChildrenREADME(root); err != nil {
		return err
	}
	return nil
}

// 递归所有节点，下载 README
func GetCurrentAndChildrenREADME(current *FileNode) error {
	if current == nil {
		return fmt.Errorf("GetCurrentAndChildrenREADME get a nil pointer")
	}

	// 当前节点有 READMEURL，就下载存到 cache
	if current.READMEUrl != "" {
		if readmeBytes, err := RequestOneUrl(current.READMEUrl); err != nil {
			log.WithFields(log.Fields{
				"path": current.Path,
				"url":  current.READMEUrl,
			}).Infof("download readme file to cache error")
		} else {
			p := GetReplacePath(current.Path)
			reCache.Set(READEME+p, readmeBytes, DefaultTime)
		}
	}

	for i := range current.Children {
		if err := GetCurrentAndChildrenREADME(current.Children[i]); err != nil {
			return err
		}
	}
	return nil
}

func GetREADMEInCache(p string) ([]byte, error) {
	ans, ok := reCache.Get(READEME + p)
	if !ok {
		log.WithFields(log.Fields{
			"path": p,
		}).Info("README not in cache")
		return nil, fmt.Errorf("README not in cache")
	}

	return ans.([]byte), nil
}

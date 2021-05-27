package onedrive

import (
	"fmt"
	gocache "github.com/patrickmn/go-cache"
	log "github.com/sirupsen/logrus"
	"gonelist/pkg/markdown"
	"strings"
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

// 刷新每个文件夹的 README 和 Password
func RefreshREADME() error {
	// 获取根节点开始
	root := FileTree.GetRoot()
	if err := GetAllREADMEAndPass(root); err != nil {
		return err
	}
	return nil
}

// 递归所有节点，下载 README
func GetAllREADMEAndPass(current *FileNode) error {
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
			// 转化成 HTML 的结果
			finalBytes := markdown.MarkdownToHTMLByBytes(readmeBytes)
			reCache.Set(READEME+p, finalBytes, DefaultTime)
		}
	}

	// 当前节点有 .password，下载并且赋值
	if current.PasswordUrl != "" {
		if readmeBytes, err := RequestOneUrl(current.PasswordUrl); err != nil {
			log.WithFields(log.Fields{
				"path": current.Path,
				"url":  current.PasswordUrl,
			}).Infof("download password file error")
		} else {
			current.Password = strings.TrimSpace(string(readmeBytes))
		}
	}

	for i := range current.Children {
		if err := GetAllREADMEAndPass(current.Children[i]); err != nil {
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

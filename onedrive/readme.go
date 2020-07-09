package onedrive

import (
	gocache "github.com/patrickmn/go-cache"
	log "github.com/sirupsen/logrus"
	"gonelist/pkg/file"
	"time"
)

// 设置缓存的默认时间为 2 天，每 2 天清空已经失效的缓存
var reCache = gocache.New(time.Hour*48, time.Hour*48)

// 刷新每个文件夹的 README
func RefreshREADME() error {
	// 获取根节点开始
	root := FileTree.GetRoot()
	if err := GetCurrentAndChildrenREADME(root); err != nil {
		return err
	}
	return nil
}

func GetCurrentAndChildrenREADME(current *FileNode) error {
	// 当前节点有 READMEURL，就下载存到 cache

	//for _, item := range current {
	//
	//}
	return nil
}

// 遍历获取所有的 README 文件
func DownloadREADME() {
	README := "README.md"
	log.Info("下载", README)

	// 判断是否有 README.md 这个文件
	if file.IsExistFile(README) {
		log.Info("已有 README.md 不进行下载")
	}

	if err := DownloadRootPathFile(README); err != nil {
		log.Warn("下载 README.md 失败")
	}

}

// 传入 filePath 来下载对应文件，暂时保存到根目录
func DownloadRootPathFile(filePath string) error {
	root := FileTree.GetRoot()
	for _, item := range root.Children {
		if item.Name == filePath {
			err := file.DownloadFile(item.DownloadUrl, filePath)
			if err != nil {
				return err
			}
			break
		}
	}
	return nil
}

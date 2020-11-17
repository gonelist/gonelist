package onedrive

import (
	"errors"
	"gonelist/conf"
	"strings"

	log "github.com/sirupsen/logrus"
)

// 初始化登陆状态
// 如果初始化时获取失败直接退出
// 如果在自动刷新时失败给出 error 警告，见 onedrive/timer.go
func InitOnedive() {
	// 获取文件内容和初始化 README 缓存
	err := RefreshOnedriveAll()
	if err != nil {
		log.WithField("err", err).Fatal("InitOnedrive 出现错误")
	}
	// 设置 onedrive 登陆状态
	FileTree.SetLogin(true)
	cacheGoOnce.Do(func() {
		go SetAutoRefresh()
	})
}

// 刷新所有 onedrive 的内容
// 包括 文件列表，README，password，搜索索引
func RefreshOnedriveAll() error {
	log.Info("开始刷新文件缓存")
	if _, err := GetAllFiles(); err != nil { // 获取所有文件并且刷新树结构
		log.WithField("err", err).Error("刷新文件缓存遇到错误")
		return err
	}
	log.Infof("结束刷新文件缓存")
	log.Debug(FileTree)

	log.Info("开始刷新 README 缓存")
	if err := RefreshREADME(); err != nil {
		log.WithField("err", err).Error("刷新 README 缓存遇到错误")
		return err
	}
	log.Info("结束刷新 README 缓存")
	// 构建搜索
	return nil
}

// 从缓存获取某个路径下的所有内容
func CacheGetPathList(oPath string) (*FileNode, error) {
	var (
		root    *FileNode
		isFound bool
	)

	root = FileTree.GetRoot()
	pArray := strings.Split(oPath, "/")

	if oPath == "" || oPath == "/" || len(pArray) < 2 {
		return ConvertReturnNode(root), nil
	}

	for i := 1; i < len(pArray); i++ {
		isFound = false
		for _, item := range root.Children {
			if pArray[i] == item.Name {
				root = item
				isFound = true
			}
		}
		if isFound == false {
			log.WithFields(log.Fields{
				"oPath":    oPath,
				"pArray":   pArray,
				"position": pArray[i],
			})
			return nil, errors.New("未找到该路径")
		}
	}

	// 只返回当前层的内容
	reNode := ConvertReturnNode(root)
	return reNode, nil
}

func ConvertReturnNode(node *FileNode) *FileNode {
	if node == nil {
		return nil
	}

	reNode := CopyFileNode(node)
	for key := range node.Children {
		if node.Children[key].Name == ".password" {
			continue
		}
		tmpNode := node.Children[key]
		reNode.Children = append(reNode.Children, CopyFileNode(tmpNode))
	}
	return reNode
}

func CopyFileNode(node *FileNode) *FileNode {
	if node == nil {
		return nil
	}
	//path := GetReplacePath(node.Path)

	return &FileNode{
		Name:           node.Name,
		Path:           node.Path,
		IsFolder:       node.IsFolder,
		DownloadUrl:    node.DownloadUrl,
		LastModifyTime: node.LastModifyTime,
		Size:           node.Size,
		Children:       nil,
		Password:       node.Password,
	}
}

func GetDownloadUrl(filePath string) (string, error) {
	var (
		fileInfo    *FileNode
		err         error
		downloadUrl string
	)

	if fileInfo, err = CacheGetPathList(filePath); err != nil || fileInfo == nil || fileInfo.IsFolder == true {
		log.WithFields(log.Fields{
			"filePath": filePath,
			"err":      err,
		}).Info("请求的文件未找到")
		return "", err
	}

	// 如果有重定向前缀，就加上
	downloadUrl = conf.UserSet.DownloadRedirectPrefix + fileInfo.DownloadUrl

	return downloadUrl, nil
}

// 替换路径
// 如果设置了 folderSub 为 /public
// 那么 /public 替换为 /, /public/test 替换为 /test
func GetReplacePath(pSrc string) string {
	if conf.UserSet.Server.FolderSub != "/" {
		pSrc = strings.Replace(pSrc, conf.UserSet.Server.FolderSub, "", 1)
	}
	return pSrc
}

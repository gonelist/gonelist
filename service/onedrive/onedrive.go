package onedrive

import (
	"errors"
	"net/url"
	"strings"
	"sync"

	"gonelist/conf"
	"gonelist/pkg/file"

	log "github.com/sirupsen/logrus"
)

const (
	REFRESH_NONE    = 0 // 未在刷新状态
	REFRESH_RUNNING = 1 // 正在刷新
)

var (
	gRefreshStatus = REFRESH_NONE
	gRefreshLock   sync.Mutex
)

// 初始化登陆状态, 如果初始化时获取失败直接退出程序
// 如果在自动刷新时失败给出 error 警告，见 onedrive/timer.go
func InitOnedrive() {
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
	if gRefreshStatus == REFRESH_RUNNING {
		log.Info("当前正在刷新中")
		return nil
	}

	// TODO，修改为 TryLock，现在有概率阻塞协程
	gRefreshLock.Lock()
	defer gRefreshLock.Unlock()
	gRefreshStatus = REFRESH_RUNNING

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

	gRefreshStatus = REFRESH_NONE

	// 构建搜索
	return nil
}

// TODO
// 刷新逻辑，level 表示刷新文件的层数
// level = -1 时刷新全部文件
func RefreshOnedriveByLevel() {

}

func CacheGetPathList(oPath string) ([]*FileNode, error) {
	root, err := GetNode(oPath)
	if err != nil {
		return []*FileNode{}, err
	}
	return ReturnNode(root), nil
}

// 获取树的某个节点，不论是不是叶子节点
func GetNode(oPath string) (*FileNode, error) {
	var (
		root    *FileNode
		isFound bool
	)

	root = FileTree.GetRoot()
	oPath = strings.TrimRight(oPath, "/")
	pArray := strings.Split(oPath, "/")
	if oPath == "" || oPath == "/" {
		//return ConvertReturnNode(root), nil
		//return root.Children, nil
		return root, nil
	}

	for i := 1; i < len(pArray); i++ {
		isFound = false
		for _, item := range root.Children {
			if pArray[i] == item.Name {
				root = item
				isFound = true
				break
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

	return root, nil
	// 只返回当前层的内容
	//reNode := ConvertReturnNode(root)
	//return reNode, nil
}

func ReturnNode(node *FileNode) []*FileNode {
	var reNode []*FileNode
	if node == nil {
		return reNode
	}
	if node.Children == nil {
		return make([]*FileNode, 0)
	}

	return node.Children
}

// 旧版返回逻辑
//func ConvertReturnNode(node *FileNode) *FileNode {
//	if node == nil {
//		return nil
//	}
//
//	reNode := CopyFileNode(node)
//	for key := range node.Children {
//		if node.Children[key].Name == ".password" {
//			continue
//		}
//		tmpNode := node.Children[key]
//		reNode.Children = append(reNode.Children, CopyFileNode(tmpNode))
//	}
//	return reNode
//}
//
//func CopyFileNode(node *FileNode) *FileNode {
//	if node == nil {
//		return nil
//	}
//	//path := GetReplacePath(node.Path)
//
//	return &FileNode{
//		Name:           node.Name,
//		Path:           node.Path,
//		IsFolder:       node.IsFolder,
//		DownloadUrl:    node.DownloadUrl,
//		LastModifyTime: node.LastModifyTime,
//		Size:           node.Size,
//		Children:       nil,
//		Password:       node.Password,
//	}
//}
//func GetDownloadUrl(filePath string) (string, error) {
//	var (
//		fileInfo    *FileNode
//		err         error
//		downloadUrl string
//	)
//
//	if fileInfo, err = GetNode(filePath); err != nil || fileInfo == nil || fileInfo.DownloadUrl == "" {
//		log.WithFields(log.Fields{
//			"filePath": filePath,
//			"err":      err,
//		}).Info("请求的文件未找到")
//		return "", err
//	}
//
//	// 如果有重定向前缀，就加上
//	downloadUrl = conf.UserSet.Onedrive.DownloadRedirectPrefix + fileInfo.DownloadUrl
//
//	return downloadUrl, nil
//}

func GetDownloadUrl(filePath string) (string, error) {

	var (
		fileInfo    *FileNode
		err         error
		downloadUrl string
	)

	// 判断同级目录下是否存在.password
	if fileInfo, err = GetNode(file.FatherPath(filePath)); err != nil || fileInfo.PasswordUrl == ".password" {
		log.WithFields(log.Fields{
			"filePath": filePath,
			"err":      err,
		}).Info("请求的文件未找到")
		return "", err
	}
	apath := strings.Split(filePath, "/")
	fileName := apath[len(apath)-1]
	for _, item := range fileInfo.Children {
		if item.Name == fileName {
			if conf.UserSet.Onedrive.DownloadRedirectPrefix == "" {
				downloadUrl = item.DownloadUrl
			} else {
				downloadUrl = conf.UserSet.Onedrive.DownloadRedirectPrefix + url.QueryEscape(item.DownloadUrl)
			}
			return downloadUrl, nil
		}
	}

	return "", err

}

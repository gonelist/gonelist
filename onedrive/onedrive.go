package onedrive

import (
	"errors"
	log "github.com/sirupsen/logrus"
	"gonelist/conf"
	"strings"
)

// 初始化登陆状态
func InitOnedive() {
	// 获取文件内容和初始化 README 缓存
	if _, err := GetAllFiles(); err != nil {
		log.Fatal(err)
	}
	if err := RefreshREADME(); err != nil {
		log.Fatal(err)
	}
	// 设置 onedrive 登陆状态
	FileTree.SetLogin(true)
	cacheGoOnce.Do(func() {
		go SetAutoRefresh()
	})
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
		tmpNode := node.Children[key]
		reNode.Children = append(reNode.Children, CopyFileNode(tmpNode))
	}
	return reNode
}

func CopyFileNode(node *FileNode) *FileNode {
	if node == nil {
		return nil
	}
	var (
		folderSub string
		path      string
	)

	// 替换相对路径
	if folderSub = conf.UserSet.Server.FolderSub; folderSub != "/" {
		path = strings.Replace(node.Path, conf.UserSet.Server.FolderSub, "", 1)
	} else {
		path = node.Path
	}

	return &FileNode{
		Name:           node.Name,
		Path:           path,
		IsFolder:       node.IsFolder,
		DownloadUrl:    node.DownloadUrl,
		LastModifyTime: node.LastModifyTime,
		Size:           node.Size,
		Children:       nil,
	}
}

func GetDownloadUrl(filePath string) (string, error) {
	var (
		fileInfo *FileNode
		err      error
	)

	if fileInfo, err = CacheGetPathList(filePath); err != nil || fileInfo == nil || fileInfo.IsFolder == true {
		log.WithFields(log.Fields{
			"filePath": filePath,
			"err":      err,
		}).Info("请求的文件未找到")
		return "", err
	}

	return fileInfo.DownloadUrl, nil
}

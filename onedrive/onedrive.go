package onedrive

import (
	"errors"
	log "github.com/sirupsen/logrus"
	"gonelist/conf"
	"gonelist/pkg/file"
	"strings"
)

// 从缓存获取
func CacheGetPathList(oPath string) (*FileNode, error) {
	var (
		root    *FileNode
		isFound bool
	)

	root = FileTree
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
	root := FileTree
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

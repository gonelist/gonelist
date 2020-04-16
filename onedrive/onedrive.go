package onedrive

import (
	"encoding/json"
	"errors"
	log "github.com/sirupsen/logrus"
	"gonelist/conf"
	"gonelist/mg_auth"
	"io/ioutil"
	"strings"
	"time"
)

const (
	ROOTUrl = "https://graph.microsoft.com/v1.0/me/drive/root/children"
)

// 获取某个路径的内容，如果 token 失效或没有正常结果返回 err
func GetUrlToAns(relativePath string) (Answer, error) {
	var (
		url = ROOTUrl
		ans Answer
	)

	if relativePath != "" {
		// eg. /test
		url = "https://graph.microsoft.com/v1.0/me/drive/root:" + relativePath + ":/children"
	}

	client := mg_auth.GetClient()
	resp, err := client.Get(url)
	if err != nil {
		log.WithFields(log.Fields{
			"url":  url,
			"resp": resp,
			"err":  err,
		}).Info("请求 graph.microsoft.com 失败")
		return ans, err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.WithField("err", err).Info("读取 graph.microsoft.com 返回内容失败")
		return ans, err
	}

	// 解析内容
	json.Unmarshal(body, &ans)
	err = CheckAnswerValid(ans, relativePath)
	if err != nil { //如果获取内容
		return ans, err
	}
	return ans, nil
}

// 获取所有文件的树
func GetAllFiles() (*FileNode, error) {
	var (
		err    error
		prefix string
		root   *FileNode
	)

	root = &FileNode{
		Name:           "root",
		Path:           "/",
		IsFolder:       false,
		LastModifyTime: time.Now(),
		Children:       nil,
	}

	if conf.UserSet.Server.FolderSub == "/" {
		prefix = ""
	} else {
		prefix = conf.UserSet.Server.FolderSub
	}
	list, err := GetTreeFileNode(prefix, "")
	if err != nil {
		log.Info(err)
		return nil, err
	} else {
		root.Children = list
		if root.Children != nil {
			root.IsFolder = true
		}
	}
	// 更新树结构
	FileTree = root
	return root, nil
}

func GetTreeFileNode(prefix, relativePath string) (list []*FileNode, err error) {
	var (
		ans   Answer
		oPath = prefix + relativePath
	)

	ans, err = GetUrlToAns(oPath)
	if err != nil {
		log.WithFields(log.Fields{
			"ans": ans,
			"err": err,
		}).Info("请求 graph.microsoft.com 出现错误")
		return nil, err
	}

	// 解析对应 list
	list = ConvertAnsToFileNodes(oPath, ans)
	for i, _ := range list {
		// 如果下一层是文件夹则继续
		if list[i].IsFolder == true {
			tmpList, err := GetTreeFileNode(list[i].Path, "")
			if err == nil {
				list[i].Children = tmpList
			}
		}
	}
	return list, nil
}

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
	for key, _ := range node.Children {
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
		file *FileNode
		err  error
	)

	if file, err = CacheGetPathList(filePath); err != nil || file == nil || file.IsFolder == true {
		log.WithFields(log.Fields{
			"filePath": filePath,
			"err":      err,
		}).Info("请求的文件未找到")
		return "", err
	}

	return file.DownloadUrl, nil
}

package onedrive

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"gonelist/pkg/markdown"
	"gonelist/service/onedrive/cache"
	"gonelist/service/onedrive/model"

	log "github.com/sirupsen/logrus"
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

	// 构建搜索
	go RefreshFiles()
	go RefreshReadme()
	go RefreshPassword()

	return nil
}

func RefreshPassword() {
	defer func() {
		err := recover()
		if err != nil {
			log.Errorln("刷新缓存过程中出现了不可预料的错误")
			log.Errorln(err)
		}
	}()
	log.Infoln("开始刷新password缓存")
	err := GetPasswordNode()
	if err != nil {
		log.Errorln(err.Error())
		return
	}
	log.Infoln("password缓存刷新结束")
}

func RefreshReadme() {
	defer func() {
		err := recover()
		if err != nil {
			log.Errorln("刷新缓存过程中出现了不可预料的错误")
			log.Errorln(err)
		}
	}()
	log.Infoln("开始刷新README.md缓存")
	err := GetReadMeNodes()
	if err != nil {
		log.Errorln(err.Error())
		return
	}
	log.Infoln("README.md缓存刷新结束")
}

func RefreshFiles() {
	defer func() {
		err := recover()
		if err != nil {
			log.Errorln("刷新缓存过程中出现了不可预料的错误")
			log.Errorln(err)
		}
	}()
	log.Info("开始刷新文件缓存")
	start := time.Now()
	_, token, _ := Delta(getToken())
	setToken(token)
	duration := time.Since(start)
	log.Infoln(fmt.Sprintf("共用时%.2f分钟", duration.Minutes()))
	log.Infoln("刷新文件缓存结束")
}

func GetPasswordNode() error {
	passNodes, err := model.FindByName(".password")
	if err != nil {
		log.Errorln("查询password节点出现错误")
		return err
	}
	for _, node := range passNodes {
		var downloadUrl string
		if node.DownloadURL != "" {
			downloadUrl = node.DownloadURL
		} else {
			url, err := getDownloadUrl(node)
			if err != nil {
				return err
			}
			downloadUrl = url
		}
		resp, err := GetData(http.MethodGet, downloadUrl, map[string]string{}, nil)
		if err != nil {
			return err
		}
		parentNode, err := model.Find(node.ParentID)
		if err != nil {
			return err
		}
		log.Debugln(string(resp))
		parentNode.Password = string(resp)
		parentNode.PasswordURL = downloadUrl
		_ = model.UpdateFile(parentNode)
	}
	return err
}

// GetReadMeNodes
/**
 * @Description: 解析文件中所有的READEME.md文件
 * @return error
 */
func GetReadMeNodes() error {
	readmeNodes, err := model.FindByName("README.md")
	if err != nil {
		log.Errorln("查询readme节点出现错误")
		return err
	}
	for _, node := range readmeNodes {
		var downloadUrl string
		if node.DownloadURL != "" {
			downloadUrl = node.DownloadURL
		} else {
			url, err := getDownloadUrl(node)
			if err != nil {
				return err
			}
			downloadUrl = url
		}
		resp, err := GetData(http.MethodGet, downloadUrl, map[string]string{}, nil)
		if err != nil {
			return err
		}
		parentNode, err := model.Find(node.ParentID)
		if err != nil {
			return err
		}
		data := markdown.MarkdownToHTMLByBytes(resp)
		reCache.Set(README+parentNode.Path, data, DefaultTime)
	}
	return err
}

// TODO
// 刷新逻辑，level 表示刷新文件的层数
// level = -1 时刷新全部文件
func RefreshOnedriveByLevel() {

}

func CacheGetPathList(oPath string) ([]*model.FileNode, error) {
	//root, err := GetNode(oPath)
	//if err != nil {
	//	return []*model.FileNode{}, err
	//}
	//return ReturnNode(root), nil
	node, ok := cache.Cache.Get(oPath)
	if node == nil || ok == false {
		return nil, errors.New("file not found")
	}
	nodes, err := model.GetChildrenByID(node.ID)
	if err != nil {
		return []*model.FileNode{}, err
	}
	return nodes, err
}

// 获取树的某个节点，不论是不是叶子节点
func GetNode(oPath string) (*model.FileNode, error) {
	var (
		root    *model.FileNode
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

func ReturnNode(node *model.FileNode) []*model.FileNode {
	var reNode []*model.FileNode
	if node == nil {
		return reNode
	}
	if node.Children == nil {
		return make([]*model.FileNode, 0)
	}

	return node.Children
}

// 旧版返回逻辑
//func ConvertReturnNode(node *model.FileNode) *FileNode {
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

func getDownloadUrl(node *model.FileNode) (string, error) {
	baseURl := "https://graph.microsoft.com/v1.0/me/drive/items/" + node.ID
	resp, err := GetData(http.MethodGet, baseURl, map[string]string{}, nil)
	if err != nil {
		return "", err
	}
	v := new(Value)
	err = json.Unmarshal(resp, v)
	if err != nil {
		return "", err
	}
	return v.MicrosoftGraphDownloadURL, err
}

func GetDownloadUrl(filePath string) (string, error) {

	node, ok := cache.Cache.Get(filePath)
	if !ok {
		return "", errors.New("file not found")
	}
	return getDownloadUrl(node)
	//var (
	//	fileInfo    *model.FileNode
	//	err         error
	//	downloadUrl string
	//)
	//
	//// 判断同级目录下是否存在.password
	//if fileInfo, err = GetNode(file.FatherPath(filePath)); err != nil || fileInfo.PasswordURL == ".password" {
	//	log.WithFields(log.Fields{
	//		"filePath": filePath,
	//		"err":      err,
	//	}).Info("请求的文件未找到")
	//	return "", err
	//}
	//apath := strings.Split(filePath, "/")
	//fileName := apath[len(apath)-1]
	//for _, item := range fileInfo.Children {
	//	if item.Name == fileName {
	//		if conf.UserSet.Onedrive.DownloadRedirectPrefix == "" {
	//			downloadUrl = item.DownloadURL
	//		} else {
	//			downloadUrl = conf.UserSet.Onedrive.DownloadRedirectPrefix + url.QueryEscape(item.DownloadURL)
	//		}
	//		return downloadUrl, nil
	//	}
	//}
	//
	//return "", err

}

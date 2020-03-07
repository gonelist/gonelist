package onedrive

import (
	"encoding/json"
	"errors"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"GOIndex/mg_auth"
	"strings"
	"time"
)

// 获取某个路径的内容，如果 token 失效或没有正常结果返回 err
func GetUrlToAns(relativePath string) (Answer, error) {
	var url string
	var ans Answer

	if relativePath == "" {
		url = "https://graph.microsoft.com/v1.0/me/drive/root/children"
	} else {
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
func GetAllFiles() *FileNode {
	var err error
	root := &FileNode{
		Name:           "root",
		Path:           "/",
		LastModifyTime: time.Now(),
		Children:       nil,
	}

	list, err := GetTreeFileNode("", "")
	if err != nil {
		log.Info(err)
	} else {
		root.Children = list
	}
	return root
}

func GetTreeFileNode(prefix, relativePath string) (list []*FileNode, err error) {
	var ans Answer
	oPath := prefix + relativePath
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

func CacheGetPathList(oPath string) ([]*FileNode, error) {
	root := FileTree
	pArray := strings.Split(oPath, "/")

	if oPath == "" || oPath == "/" || len(pArray) < 2 {
		return root.Children, nil
	}

	for i := 1; i < len(pArray); i++ {
		flag := false
		for _, item := range root.Children {
			if pArray[i] == item.Name {
				root = item
				flag = true
			}
		}
		if flag == false {
			log.WithFields(log.Fields{
				"oPath":    oPath,
				"pArray":   pArray,
				"position": pArray[i],
			})
			return nil, errors.New("未找到该路径")
		}
	}

	return root.Children, nil
}

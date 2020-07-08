package onedrive

import (
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	"gonelist/conf"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

var (
	// 默认 https://graph.microsoft.com/v1.0/me/drive/root/children
	// ChinaCloud https://microsoftgraph.chinacloudapi.cn/v1.0/me/drive/root/children
	ROOTUrl  string
	UrlBegin string
	UrlEnd   string
)

func SetROOTUrl(chinaCloud bool) {
	if chinaCloud == false {
		ROOTUrl = "https://graph.microsoft.com/v1.0/me/drive/root/children"
		UrlBegin = "https://graph.microsoft.com/v1.0/me/drive/root:"
	} else {
		ROOTUrl = "https://microsoftgraph.chinacloudapi.cn/v1.0/me/drive/root/children"
		UrlBegin = "https://microsoftgraph.chinacloudapi.cn/v1.0/me/drive/root:"
	}
	UrlEnd = ":/children"
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

// 获取树的一个节点
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
		}).Errorf("请求 graph.microsoft.com 出现错误: prefix:%v, relativePath:%v", prefix, relativePath)
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

// 获取某个路径的内容，如果 token 失效或没有正常结果返回 err
func GetUrlToAns(relativePath string) (Answer, error) {
	// 默认一次获取 3000 个文件
	var (
		url    = ROOTUrl + "?$top=3000"
		ans    Answer
		tmpAns Answer
		err    error
	)

	if relativePath != "" {
		// 每次获取 3000 个文件
		// eg. /test -> https://graph.microsoft.com/v1.0/me/drive/root:/test:/children
		url = UrlBegin + relativePath + UrlEnd + "?$top=3000"
	}

	for {
		tmpAns, err = RequestAnswer(url, relativePath)
		// 判断是否合并两次请求
		if len(ans.Value) == 0 {
			ans = tmpAns
		} else {
			ans.Value = append(ans.Value, tmpAns.Value...)
		}
		// 判断是否继续请求下一个链接
		if err != nil {
			return ans, err
		} else if tmpAns.OdataNextLink == "" {
			break
		}
		url = ans.OdataNextLink
	}

	return ans, nil
}

func RequestAnswer(url string, relativePath string) (Answer, error) {
	var (
		ans Answer
	)
	body, err := RequestOneUrl(url)
	if err != nil {
		return ans, err
	}
	// 解析内容
	if err := json.Unmarshal(body, &ans); err != nil {
		return ans, err
	}
	log.Debugf("url:%s relativePath:%s | body:%s", url, relativePath, string(body))
	err = CheckAnswerValid(ans, relativePath)
	if err != nil { //如果获取内容
		return ans, err
	}
	return ans, nil
}

func RequestOneUrl(url string) (body []byte, err error) {

	var (
		client *http.Client // 获取全局的 client 来请求接口
		resp   *http.Response
	)
	if client = GetClient(); client == nil {
		log.Errorln("cannot get client to start request.")
		return nil, fmt.Errorf("cannot get client")
	}

	// 如果超时，重试两次
	for retryCount := 2; retryCount > 0; retryCount-- {
		if resp, err = client.Get(url); err != nil && strings.Contains(err.Error(), "timeout") {
			<-time.After(time.Second)
		} else {
			break
		}
	}

	if err != nil {
		log.WithFields(log.Fields{
			"url":  url,
			"resp": resp,
			"err":  err,
		}).Info("请求 graph.microsoft.com 失败")
		return body, err
	}

	if body, err = ioutil.ReadAll(resp.Body); err != nil {
		log.WithField("err", err).Info("读取 graph.microsoft.com 返回内容失败")
		return body, err
	}
	return body, nil
}

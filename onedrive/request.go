package onedrive

import (
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	"gonelist/conf"
	"io/ioutil"
	"net/http"
	"net/url"
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
	if chinaCloud {
		ROOTUrl = "https://microsoftgraph.chinacloudapi.cn/v1.0/me/drive/root/children"
		UrlBegin = "https://microsoftgraph.chinacloudapi.cn/v1.0/me/drive/root:"
	} else {
		ROOTUrl = "https://graph.microsoft.com/v1.0/me/drive/root/children"
		UrlBegin = "https://graph.microsoft.com/v1.0/me/drive/root:"
	}
	UrlEnd = ":/children"
}

// 获取所有文件的树
func GetAllFiles() (*FileNode, error) {
	var (
		err  error
		root *FileNode
	)

	root = &FileNode{
		Name:           "root",
		Path:           "/",
		IsFolder:       false,
		LastModifyTime: time.Now(),
		Children:       nil,
	}

	list, readmeUrl, passUrl, err := GetTreeFileNode("/")
	if err != nil {
		log.Info(err)
		return nil, err
	}
	root.Children = list
	root.READMEUrl = readmeUrl
	root.PasswordUrl = passUrl
	if root.Children != nil {
		root.IsFolder = true
	}

	// 更新树结构
	FileTree.SetRoot(root)
	// 更新索引结构
	FileTree.SetIndex()
	return root, nil
}

// 获取树的一个节点
// list 返回当前文件夹中的所有文件夹和文件
// readmeUrl 这个是当前文件夹 readme 文件的下载链接
// err 返回错误
func GetTreeFileNode(relativePath string) (list []*FileNode, readmeUrl, passUrl string, err error) {
	var (
		ans Answer
	)

	ans, err = GetUrlToAns(relativePath)
	if err != nil {
		log.WithFields(log.Fields{
			"ans": ans,
			"err": err,
		}).Errorf("请求 graph.microsoft.com 出现错误: relativePath:%v", relativePath)
		return nil, "", "", err
	}

	// 解析对应 list
	list = ConvertAnsToFileNodes(relativePath, ans)
	for i := range list {
		// 如果下一层是文件夹则继续
		if list[i].IsFolder == true {
			tmpList, tmpReadmeUrl, tmpPassUrl, err := GetTreeFileNode(list[i].Path)
			if err == nil {
				list[i].Children = tmpList
				list[i].READMEUrl = tmpReadmeUrl
				list[i].PasswordUrl = tmpPassUrl
			}
		} else if list[i].Name == "README.md" {
			readmeUrl = list[i].DownloadUrl
		} else if list[i].Name == ".password" {
			passUrl = list[i].DownloadUrl
			// 隐藏 .password 文件的 url 和 size
			list[i].DownloadUrl = ""
			list[i].Size = 0
		}
	}
	return list, readmeUrl, passUrl, nil
}

// 获取某个路径的内容，如果 token 失效或没有正常结果返回 err
func GetUrlToAns(relativePath string) (Answer, error) {
	// 默认一次获取 3000 个文件
	var (
		baseURL string
		ans     Answer
		tmpAns  Answer
		err     error
	)

	// 每次获取 3000 个文件
	if relativePath == "/" && conf.UserSet.Server.FolderSub == "/" {
		// https://graph.microsoft.com/v1.0/me/drive/root/children
		baseURL = ROOTUrl + "?$top=3000"
	} else if relativePath == "/" {
		// eg. /test -> https://graph.microsoft.com/v1.0/me/drive/root:/test:/children
		// UrlBegin: https://graph.microsoft.com/v1.0/me/drive/root:
		// conf.UserSet.Server.FolderSub: /public
		// UrlEnd: :/children
		baseURL = UrlBegin + conf.UserSet.Server.FolderSub + UrlEnd + "?$top=3000"
	} else {
		// TODO ，好像会出现多个 / 的情况，但是暂时不影响使用
		//baseURL = UrlBegin + conf.UserSet.Server.FolderSub + relativePath + UrlEnd + "?$top=3000"

		baseURL = UrlBegin + url.QueryEscape(conf.UserSet.Server.FolderSub+relativePath) + UrlEnd + "?$top=3000"
		baseURL = strings.Replace(baseURL, "+", "%20", -1)
	}

	for {
		tmpAns, err = RequestAnswer(baseURL, relativePath)
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
		baseURL = ans.OdataNextLink
	}

	return ans, nil
}

func RequestAnswer(urlstr string, relativePath string) (Answer, error) {
	var (
		ans Answer
	)
	if strings.Contains(urlstr, "%") {
		log.Debugf("123")
	}
	// 首先对非英文字符进行 encode，兼容出现 %20 和 特殊字符 等情况
	//encodeURL := url.QueryEscape(urlstr)
	//body, err := RequestOneUrl(encodeURL)
	//m, err := url.Parse(urlstr)
	//if err != nil {
	//	log.Infof("url 出现问题，url:%s", urlstr)
	//	return ans, fmt.Errorf("url 出现问题")
	//}
	//encodeURL := m.String()
	//body, err := RequestOneUrl(encodeURL)
	body, err := RequestOneUrl(urlstr)
	if err != nil {
		return ans, err
	}

	// 解析内容
	if err := json.Unmarshal(body, &ans); err != nil {
		return ans, err
	}
	log.Debugf("url:%s relativePath:%s | body:%s", urlstr, relativePath, string(body))
	err = CheckAnswerValid(ans, relativePath)

	//如果获取内容不正常，则返回
	if err != nil {
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
		return nil, fmt.Errorf("RequestOneURL cannot get client")
	}

	// 如果超时，重试两次
	for retryCount := 3; retryCount > 0; retryCount-- {
		if resp, err = client.Get(url); err != nil && strings.Contains(err.Error(), "timeout") {
			log.WithFields(log.Fields{
				"url":  url,
				"resp": resp,
				"err":  err,
			}).Info("RequestOneUrl 出现错误，开始重试")
			<-time.After(time.Second / 3)
		} else {
			break
		}
	}

	if err != nil {
		log.WithFields(log.Fields{
			"url":  url,
			"resp": resp,
			"err":  err,
		}).Info("请求 graph.microsoft.com 失败, request timeout")
		return body, err
	}

	if body, err = ioutil.ReadAll(resp.Body); err != nil {
		log.WithField("err", err).Info("读取 graph.microsoft.com 返回内容失败")
		return body, err
	}
	return body, nil
}

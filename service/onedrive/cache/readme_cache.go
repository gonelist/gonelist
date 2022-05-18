package cache

import (
	"encoding/json"
	"net/http"
	"strings"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"

	"gonelist/conf"
	"gonelist/pkg/markdown"
	"gonelist/service/onedrive/model"
	"gonelist/service/onedrive/pojo"
	"gonelist/service/onedrive/utils"
)

var (
	ReadmeCache *READMECache
)

type ReadmeData struct {
	ExpiredTime int64 //数据的过期时间戳
	Data        []byte
}

type READMECache struct {
	cap         int                    // 缓存容量
	datas       map[string]*ReadmeData // map结构，存储数据
	list        *DoubleList            // 双向链表
	expiredTime int                    // Lru设置的过期时间
	rwLock      sync.RWMutex
}

func init() {
	InitREADMECache(10, 30)
}

// InitREADMECache 初始化LRU缓存的容量大小，过期时间
func InitREADMECache(cap, expiredTime int) {
	ReadmeCache = new(READMECache)
	ReadmeCache.cap = cap
	ReadmeCache.expiredTime = expiredTime
	ReadmeCache.list = NewDoubleList()
	ReadmeCache.datas = make(map[string]*ReadmeData)
}

// PutREADME
/**
 * @Description:
 * @receiver cache
 * @param node
 * @return error
 */
func (c *READMECache) PutREADME(path string, data []byte) error {
	c.rwLock.Lock()
	defer c.rwLock.Unlock()

	if _, ok := c.datas[path]; ok {
		// 数据已经存在缓存中但还在插入，说明缓存中数据已经过期，重新更新数据
		c.list.MoveToHead(path)
		c.datas[path] = &ReadmeData{
			ExpiredTime: time.Now().Add(time.Minute * time.Duration(c.expiredTime)).Unix(),
			Data:        data,
		}
		return nil
	}

	// 缓存中容量超过上限
	if c.list.Size >= c.cap {
		// 删除链表尾部节点
		k := c.list.RemoveOneNodeByTail()
		delete(c.datas, k)
	}
	// 正常插入数据
	c.list.InsertList(path)
	c.datas[path] = &ReadmeData{
		ExpiredTime: time.Now().Add(time.Minute * time.Duration(c.expiredTime)).Unix(),
		Data:        data,
	}
	return nil
}

// GetREADME
/**
 * @Description: 从缓存中获取数据
 * @receiver c
 * @param key
 * @return *model.FileNode
 * @return bool
 */
func (c *READMECache) GetREADME(key string) ([]byte, bool) {
	//c.rwLock.Lock()
	//defer c.rwLock.Unlock()
	if conf.UserSet.Onedrive.FolderSub != "/" {
		key = conf.UserSet.Onedrive.FolderSub + key
		key = strings.TrimRight(key, "/")
	} else if key == "/" { // 在不设置folder_path下，根目录会出现不能访问，原因是key=/
		key = ""
	}
	// 先从缓存中取出数据，如果数据存在并且数据没用过去则直接返回值，数据过期直接走插入流程
	data, ok := c.datas[key]
	if ok && data.ExpiredTime > time.Now().Unix() {
		c.list.MoveToHead(key)
		return data.Data, true
	}
	// 先去从数据库中获取node节点，不从缓存获取是因为已经加上了folder_sub，从缓存获取会出错
	node, err := model.FindByPath(key + "/README.md")
	if err != nil {
		return nil, false
	}
	// 获取节点文件的下载链接
	url, err := getDownloadUrl(node)
	if err != nil {
		log.Errorln("获取节点下载链接错误" + err.Error())
		return nil, false
	}
	content, err := utils.GetData(http.MethodGet, url, map[string]string{}, nil)
	if err != nil {
		return nil, false
	}
	d := markdown.MarkdownToHTMLByBytes(content)
	_ = c.PutREADME(key, d)
	//node, err := model.FindByPath(key)
	//if err != nil {
	//	return nil, false
	//}
	//_ = c.Put(node)
	//return node, true
	return d, true
}

func getDownloadUrl(node *model.FileNode) (string, error) {
	baseURl := "https://graph.microsoft.com/v1.0/me/drive/items/" + node.ID
	resp, err := utils.GetData(http.MethodGet, baseURl, map[string]string{}, nil)
	if err != nil {
		return "", err
	}
	v := new(pojo.Value)
	err = json.Unmarshal(resp, v)
	if err != nil {
		return "", err
	}
	return v.MicrosoftGraphDownloadURL, err
}

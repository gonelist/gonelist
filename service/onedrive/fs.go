package onedrive

import (
	"fmt"
	"sort"
	"sync"
	"time"

	gocache "github.com/patrickmn/go-cache"
	log "github.com/sirupsen/logrus"
)

// 设置缓存的默认时间为 2 天，每 2 天清空已经失效的缓存
var reCache = gocache.New(DefaultTime, DefaultTime)

// 在缓存中 key 的形式是 README_path
// eg. README_/, README_/exampleFolder
const (
	README      = "README_"
	DefaultTime = time.Hour * 24
	FS          = "FS_"
)

// 存储的目录结构
// RefreshTime 表示最近一次的刷新时间
type FileNode struct {
	ID             string      `json:"id"`
	Name           string      `json:"name"`
	Path           string      `json:"path"`
	READMEUrl      string      `json:"-"`
	IsFolder       bool        `json:"is_folder"`
	DownloadUrl    string      `json:"-"`
	LastModifyTime time.Time   `json:"last_modify_time"`
	Size           int64       `json:"size"`
	Children       []*FileNode `json:"-"`
	RefreshTime    time.Time   `json:"refresh_time"`
	// .password 内容
	Password    string `json:"-"`
	PasswordUrl string `json:"-"`
}

// Answer 是请求接口返回内容，里面包含的 Value 是一个列表
func ConvertAnsToFileNodes(oldPath string, ans Answer) []*FileNode {
	var (
		list []*FileNode
		path string
	)

	for _, item := range ans.Value {
		if oldPath == "/" {
			path = oldPath + item.Name
		} else {
			path = oldPath + "/" + item.Name
		}
		node := &FileNode{
			ID:             item.ID,
			Name:           item.Name,
			Path:           path,
			LastModifyTime: item.FileSystemInfo.LastModifiedDateTime,
			DownloadUrl:    item.MicrosoftGraphDownloadURL,
			IsFolder:       false,
			Size:           item.Size,
			Children:       nil,
		}
		// 如果是文件夹，就设置状态，并且对文件夹设置当前的刷新时间
		if item.Folder.ChildCount != -1 {
			node.IsFolder = true
			node.RefreshTime = time.Now()
		}
		list = append(list, node)
	}

	// 对 list 进行排序
	sort.Slice(list, func(i, j int) bool {
		return list[i].Name < list[j].Name
	})

	return list
}

type Tree struct {
	sync.Mutex
	root       *FileNode
	isLogin    bool
	FirstReady int
}

var FileTree = &Tree{}

func (t *Tree) SetLogin(status bool) {
	t.Lock()
	defer t.Unlock()

	t.isLogin = status
}

func (t *Tree) IsLogin() bool {
	return t.isLogin
}

func (t *Tree) GetRoot() *FileNode {
	t.Lock()
	defer t.Unlock()

	return t.root
}

func (t *Tree) SetRoot(root *FileNode) {
	t.Lock()
	defer t.Unlock()

	t.root = root
}

func (t *Tree) dfsIndexTree(f *FileNode) {

	// 如果不是文件夹或者是一个加密的文件夹
	// 那么直接退出
	if !f.IsFolder || f.PasswordUrl != "" {
		return
	}

	//if strings.Contains(t.Name, "加密") {
	//	log.Info(t.Name)
	//}

	for _, child := range f.Children {
		// 如果该文件夹加密，则直接退出
		if child.Name == ".password" {
			return
		}
		t.dfsIndexTree(child)
	}
}

// TODO
// 真实的下载 URL 通过 Cache 进行存储
func GetPathInCache(p string) ([]byte, error) {
	ans, ok := reCache.Get(FS + p)
	if !ok {
		log.WithFields(log.Fields{
			"path": p,
		}).Info("FS not in cache")
		return nil, fmt.Errorf("FS not in cache")
	}

	return ans.([]byte), nil
}

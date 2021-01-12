package onedrive

import (
	"encoding/json"
	"errors"
	"gonelist/onedrive/internal"
	"gonelist/onedrive/normal_index"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"
)

type ErrJson struct {
	Code       string `json:"code"`
	Message    string `json:"message"`
	InnerError struct {
		RequestID string `json:"request-id"`
		Date      string `json:"date"`
	} `json:"innerError"`
}

type Folder struct {
	ChildCount int `json:"childCount"`
}

type Value struct {
	CreatedDateTime      time.Time `json:"createdDateTime"` // 创建时间
	ETag                 string    `json:"eTag"`
	ID                   string    `json:"id"`
	LastModifiedDateTime time.Time `json:"lastModifiedDateTime"`
	Name                 string    `json:"name"`
	WebURL               string    `json:"webUrl"`
	CTag                 string    `json:"cTag"`
	Size                 int64     `json:"size"`
	CreatedBy            struct {
		User struct {
			Email       string `json:"email"`
			ID          string `json:"id"`
			DisplayName string `json:"displayName"`
		} `json:"user"`
	} `json:"createdBy,omitempty"`
	LastModifiedBy struct {
		User struct {
			Email       string `json:"email"`
			ID          string `json:"id"`
			DisplayName string `json:"displayName"`
		} `json:"user"`
	} `json:"lastModifiedBy,omitempty"`
	ParentReference struct {
		DriveID   string `json:"driveId"`
		DriveType string `json:"driveType"`
		ID        string `json:"id"`
		Path      string `json:"path"`
	} `json:"parentReference"`
	FileSystemInfo struct {
		CreatedDateTime      time.Time `json:"createdDateTime"`
		LastModifiedDateTime time.Time `json:"lastModifiedDateTime"`
	} `json:"fileSystemInfo"`
	Folder        Folder `json:"folder,omitempty"`
	SpecialFolder struct {
		Name string `json:"name"`
	} `json:"specialFolder,omitempty"`
	MicrosoftGraphDownloadURL string `json:"@microsoft.graph.downloadUrl,omitempty"`
	File                      struct {
		MimeType string `json:"mimeType"`
		Hashes   struct {
			QuickXorHash string `json:"quickXorHash"`
		} `json:"hashes"`
	} `json:"file,omitempty"`
	Shared struct {
		Scope string `json:"scope"`
	} `json:"shared,omitempty"`
	Image struct {
		Height int `json:"height"`
		Width  int `json:"width"`
	} `json:"image,omitempty"`
}

type Answer struct {
	OdataContext  string  `json:"@odata.context"`
	OdataNextLink string  `json:"@odata.nextLink"`
	Value         []Value `json:"value"`
	Error         ErrJson `json:"error,omitempty"`
}

// 修改 Folder 的默认值，为 -1 时不是文件夹
func (v *Value) UnmarshalJSON(b []byte) error {
	type xvalue Value

	xf := &xvalue{Folder: Folder{ChildCount: -1}}
	if err := json.Unmarshal(b, xf); err != nil {
		return err
	}
	*v = Value(*xf)
	return nil
}

// 判断收到的 Answer 是否正常
func CheckAnswerValid(ans Answer, relativePath string) error {
	if ans.Error.Code != "" {
		log.WithFields(log.Fields{
			"Answer": ans,
			"Path":   relativePath,
		}).Info("获取的 Answer 不正确")
		return errors.New("获取的 Answer 不正确")
	}
	return nil
}

// 存储的目录结构
type FileNode struct {
	Name           string      `json:"name"`
	Path           string      `json:"path"`
	READMEUrl      string      `json:"readme_url"`
	IsFolder       bool        `json:"is_folder"`
	DownloadUrl    string      `json:"download_url"`
	LastModifyTime time.Time   `json:"last_modify_time"`
	Size           int64       `json:"size"`
	Children       []*FileNode `json:"children"`
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
			Name:           item.Name,
			Path:           path,
			LastModifyTime: item.FileSystemInfo.LastModifiedDateTime,
			DownloadUrl:    item.MicrosoftGraphDownloadURL,
			IsFolder:       false,
			Size:           item.Size,
			Children:       nil,
		}
		if item.Folder.ChildCount != -1 {
			node.IsFolder = true
		}
		list = append(list, node)
	}
	return list
}

type Tree struct {
	sync.Mutex
	root       *FileNode
	isLogin    bool
	FirstReady int
	Index      internal.Index
	NewIndex   internal.Index
}

var FileTree = &Tree{
	//Index: normal_index.NewNIndex(),
}

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

// 搜索索引相关
func (t *Tree) SetIndex() {
	t.NewIndex = normal_index.NewNIndex()
	t.dfsIndexTree(t.root)

	log.Debug(t.NewIndex)
	t.Index = t.NewIndex
	t.NewIndex = nil
}

func (t *Tree) dfsIndexTree(f *FileNode) {
	t.NewIndex.Insert(f.Name, internal.Item{
		Path:     f.Path,
		IsFolder: f.IsFolder,
	})

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

package onedrive

import (
	"encoding/json"
	"errors"
	log "github.com/sirupsen/logrus"
	"time"
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
	IsFolder       bool        `json:"is_folder"`
	DownloadUrl    string      `json:"download_url"`
	LastModifyTime time.Time   `json:"last_modify_time"`
	Size           int64       `json:"size"`
	Children       []*FileNode `json:"children"`
}

var IsLogin bool
var FileTree *FileNode

// Answer 是请求接口返回内容，里面包含的 Value 是一个列表
func ConvertAnsToFileNodes(oldPath string, ans Answer) []*FileNode {
	var list []*FileNode
	for _, item := range ans.Value {
		node := &FileNode{
			Name:           item.Name,
			Path:           oldPath + "/" + item.Name,
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

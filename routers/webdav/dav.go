package webdav

import (
	"errors"
	"io"
	"net/http"

	"github.com/emersion/go-webdav"
	log "github.com/sirupsen/logrus"

	"gonelist/service/onedrive"
	"gonelist/service/onedrive/cache"
	"gonelist/service/onedrive/model"
)

type Dav struct {
}

func (d *Dav) Open(name string) (io.ReadCloser, error) {
	log.Debugln("[webdav] " + "打开文件" + name)
	url, err := onedrive.GetDownloadUrl(name)
	if err != nil {
		log.Errorln("获取下载链接错误" + err.Error())
		return nil, err
	}
	response, err := http.Get(url)
	if err != nil {
		log.Errorln("[webdav] " + "下载文件错误" + err.Error())
		return nil, err
	}
	return response.Body, nil
}

func (d *Dav) Stat(name string) (*webdav.FileInfo, error) {
	log.Debugln("[webdav] " + "获取文件状态 ==> " + name)
	node, b := cache.Cache.Get(name)
	if !b {
		log.Errorln("获取文件状态错误")
		return nil, errors.New("file not found")
	}
	return &webdav.FileInfo{
		Path:     name,
		Size:     node.Size,
		ModTime:  node.LastModifyTime,
		IsDir:    node.IsFolder,
		MIMEType: "",
		ETag:     "",
	}, nil
}

func (d *Dav) Readdir(name string, recursive bool) ([]webdav.FileInfo, error) {
	log.Infoln("[webdav] " + "开始读取文件目录 ==> " + name)
	node, b := cache.Cache.Get(name)
	if !b {
		log.Errorln("获取文件目录错误")
		return nil, errors.New("file not found")
	}
	nodes, err := model.GetChildrenByID(node.ID)
	if err != nil {
		log.Errorln("获取文件目录子目录错误" + err.Error())
		return nil, err
	}
	files := make([]webdav.FileInfo, len(nodes))
	for _, fileNode := range nodes {
		files = append(files, webdav.FileInfo{
			Path:     fileNode.Path,
			Size:     fileNode.Size,
			ModTime:  fileNode.LastModifyTime,
			IsDir:    fileNode.IsFolder,
			MIMEType: "",
			ETag:     "",
		})
	}
	return files, err
}

func (d *Dav) Create(name string) (io.WriteCloser, error) {
	//TODO implement me
	panic("implement me")
}

func (d *Dav) RemoveAll(name string) error {
	//TODO implement me
	panic("implement me")
}

func (d *Dav) Mkdir(name string) error {
	//TODO implement me
	panic("implement me")
}

func (d *Dav) Copy(name, dest string, recursive, overwrite bool) (created bool, err error) {
	//TODO implement me
	panic("implement me")
}

func (d *Dav) MoveAll(name, dest string, overwrite bool) (created bool, err error) {
	//TODO implement me
	panic("implement me")
}

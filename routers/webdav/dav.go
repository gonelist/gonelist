package webdav

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/emersion/go-webdav"
	log "github.com/sirupsen/logrus"

	"gonelist/service/onedrive"
	"gonelist/service/onedrive/cache"
	"gonelist/service/onedrive/model"
)

// DavInit
// 初始化webdav
func DavInit() *webdav.Handler {
	defer func() {
		err := recover()
		log.Errorln("[webdav] webdav初始化异常")
		log.Errorln(err)
	}()
	han := &webdav.Handler{FileSystem: &Dav{}}
	return han
	//panic(http.ListenAndServe(fmt.Sprintf("%v:%v", conf.UserSet.WebDav.Host, conf.UserSet.WebDav.Port), nil))
}

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
	if name != "/" {
		name = strings.TrimRight(name, "/")
	}
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
	if name != "/" {
		name = strings.TrimRight(name, "/")
	}
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

	log.Debugln("[webdav] " + "开始创建文件，路径==》" + strings.TrimRight(name, filepath.Base(name)) + ", 文件名 ==》" + filepath.Base(name))
	u := &uploader{
		path: strings.TrimRight(name, filepath.Base(name)),
		name: filepath.Base(name),
		data: []byte{},
	}
	return u, nil
}

func (d *Dav) RemoveAll(name string) error {
	log.Debugln(fmt.Sprintf("[webdav] 开始删除文件 ==》 %v", name))
	name = strings.TrimRight(name, "/")
	node, b := cache.Cache.Get(name)
	if !b {
		log.Errorln("获取文件目录错误")
		return errors.New("file not found")
	}
	return onedrive.DeleteFile(node.ID)
}

func (d *Dav) Mkdir(name string) error {
	log.Debugln("[webdav] " + "开始创建文件夹，路径==》" + strings.TrimRight(name, filepath.Base(name)+"/") + ", 文件名 ==》" + filepath.Base(name))
	return onedrive.Mkdir(strings.TrimRight(name, filepath.Base(name)+"/"), filepath.Base(name))
}

func (d *Dav) Copy(name, dest string, recursive, overwrite bool) (created bool, err error) {

	// 删除/webdav前缀
	dest = trimLeft(dest)

	log.Debugln(fmt.Sprintf("[webdav] 开始从%v复制文件%v", name, strings.TrimRight(dest, filepath.Base(dest))))
	err = onedrive.Copy(name, strings.TrimRight(dest, filepath.Base(dest)))
	if err != nil {
		log.Errorln("[webdav] 复制文件出现错误 " + err.Error())
		return false, err
	}
	onedrive.RefreshFiles()
	return true, err
}

func (d *Dav) MoveAll(name, dest string, overwrite bool) (created bool, err error) {

	// 删除/webdav前缀
	dest = trimLeft(dest)

	log.Debugln(fmt.Sprintf("[webdav] 开始从%v移动文件%v", name, strings.TrimRight(dest, filepath.Base(dest))))
	err = onedrive.Move(name, strings.TrimRight(dest, filepath.Base(dest)))
	if err != nil {
		log.Errorln("[webdav] 移动文件出现错误 " + err.Error())
		return false, err
	}
	onedrive.RefreshFiles()
	return true, err
}

// 小文件上传器，仅支持4Mb以内
type uploader struct {
	path string
	name string
	data []byte
}

func (u *uploader) Write(p []byte) (n int, err error) {
	if len(u.data) > 4194304 {
		return -1, errors.New("the file over 4 Mb,the app not support")
	}
	u.data = append(u.data, p...)
	return len(p), nil
}

func (u *uploader) Close() error {
	return onedrive.Upload(u.path, u.name, u.data)
}

func trimLeft(path string) string {
	return strings.TrimPrefix(path, "/webdav")
}

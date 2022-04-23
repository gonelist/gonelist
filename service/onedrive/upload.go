package onedrive

import (
	"errors"
	"fmt"
	"net/http"

	log "github.com/sirupsen/logrus"
	"github.com/tidwall/gjson"

	"gonelist/service/onedrive/cache"
)

// Uploader
/*
 * @Description: 文件上传器
 */
type Uploader struct {
	sessionURL   string
	fileSize     int64
	currentWrite int64
}

func (u *Uploader) Close() error {
	_, err := GetData(http.MethodDelete, u.sessionURL, map[string]string{}, nil)
	if err != nil {
		return err
	}
	return err
}

func NewUploader() *Uploader {
	return &Uploader{}
}

// Write
/**
 * @Description: 实现Writer接口
 * @receiver u
 * @param p
 * @return n
 * @return err
 */
func (u *Uploader) Write(p []byte) (n int, err error) {
	m := map[string]string{"Content-Range": fmt.Sprintf("bytes %d-%d/%d",
		u.currentWrite, u.currentWrite+int64(len(p))-1, u.fileSize)}
	log.Debugln(fmt.Sprintf("文件上传中==》%v", (u.currentWrite/u.fileSize)*100))
	resp, err := GetData(http.MethodPut, u.sessionURL, m, p)
	if err != nil {
		return 0, err
	}
	u.currentWrite += int64(len(p))
	log.Debugln(gjson.GetBytes(resp, "@this|@pretty"))
	return len(p), err
}

// CreateSession
/**
 * @Description: 创建一个uploadSession
 * @receiver u
 * @param path
 * @param fileName
 * @param fileSize
 * @return error
 */
func (u *Uploader) CreateSession(path, fileName string, fileSize int64) error {
	node, b := cache.Cache.Get(path)
	if !b {
		return errors.New("parent folder not found")
	}
	sessionURL := fmt.Sprintf(ROOTUrl+"/items/%s:/%s:/createUploadSession",
		node.ID, fileName)
	data, err := GetData(http.MethodPost, sessionURL, map[string]string{}, nil)
	if err != nil {
		return err
	}
	uploadURL := gjson.GetBytes(data, "uploadUrl")
	if !uploadURL.Exists() {
		log.Errorln(gjson.GetBytes(data, "@this|@pretty"))
		return errors.New("the uploadUrl not found")
	}
	u.fileSize = fileSize
	u.sessionURL = uploadURL.String()
	log.Infoln(uploadURL.String())
	return err
}

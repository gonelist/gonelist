package file

import (
	"errors"
	"fmt"
	log "github.com/sirupsen/logrus"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

// 从某个文件读取内容
func ReadFromFile(filename string) ([]byte, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		log.WithFields(log.Fields{
			"filename": filename,
			"err":      err,
		}).Warn("ReadFromFile 出错")
		return nil, err
	}
	return data, nil
}

// 将内容写入到文件
func WriteToFile(filename string, data []byte) error {
	err := ioutil.WriteFile(filename, data, os.ModeType)
	if err != nil {
		log.WithFields(log.Fields{
			"filename": filename,
			"err":      err,
		}).Warn("WriteToFile 出错")
		return err
	}
	return nil
}

// 从 url 下载到文件
func DownloadFile(url, filename string) error {
	res, err := (&http.Client{Timeout: 60 * time.Second}).Get(url)
	if err != nil {
		log.Warn("下载文件失败")
		return errors.New("download request failed")
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("download returned HTTP %d", res.StatusCode)
	}
	f, err := os.Create(filename)
	if err != nil {
		log.WithFields(log.Fields{
			"filePath": filename,
			"err":      err,
		}).Warn("创建", filename, "失败")
		return err
	}
	defer f.Close()
	_, err = io.Copy(f, res.Body)
	if err != nil {
		log.WithFields(log.Fields{
			"filePath": filename,
			"err":      err,
		}).Warn("写入", filename, "失败")
		return err
	}
	return nil
}

// 判读是否存在文件
func IsExistFile(p string) bool {
	_, err := os.Lstat(p)
	return !os.IsNotExist(err)
}

package file

import (
	log "github.com/sirupsen/logrus"
	"io"
	"io/ioutil"
	"net/http"
	"os"
)

// 从某个文件读取内容
func ReadFromFile(filename string) ([]byte, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		log.WithFields(log.Fields{
			"filename": filename,
			"data":     data,
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
	res, err := http.Get(url)
	if err != nil {
		log.WithFields(log.Fields{
			"filePath": filename,
			"err":      err,
		}).Warn("下载", filename, "失败")
		return err
	}
	f, err := os.Create(filename)
	if err != nil {
		log.WithFields(log.Fields{
			"filePath": filename,
			"err":      err,
		}).Warn("创建", filename, "失败")
		return err
	}
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

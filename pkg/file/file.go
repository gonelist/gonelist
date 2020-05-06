package file

import (
	log "github.com/sirupsen/logrus"
	"io/ioutil"
)

// 从某个文件读取内容
func ReadFromFile(filename string) ([]byte, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		log.WithFields(log.Fields{
			"filename": filename,
			"data":     data,
			"err":      err,
		}).Info("ReadFromFile 出错")
		return nil, err
	}
	return data, nil
}

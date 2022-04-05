// Package local
// 挂载本地目录
package local

import (
	"os"
	"strings"
	"time"

	"gonelist/conf"
	"gonelist/service/onedrive/model"
)

func HandlePath(path string) string {
	p1 := strings.TrimPrefix(path, "/"+conf.UserSet.Local.Name)
	p := conf.UserSet.Local.Path + p1
	return p
}

func GetPath(path string) ([]*model.FileNode, error) {
	dir := HandlePath(path)
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	var nodes []*model.FileNode
	for _, file := range entries {
		name := file.Name()
		info, _ := file.Info()
		nodes = append(nodes, &model.FileNode{
			ID:             "",
			Name:           name,
			Path:           path + "/" + name,
			READMEUrl:      "",
			IsFolder:       file.IsDir(),
			DownloadURL:    "",
			LastModifyTime: info.ModTime(),
			Size:           0,
			Children:       nil,
			RefreshTime:    time.Time{},
			Password:       "",
			PasswordURL:    "",
			ParentID:       "",
		})
	}
	return nodes, err
}

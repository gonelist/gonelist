package model

import (
	"testing"
	"time"

	log "github.com/sirupsen/logrus"
)

var (
	file = &FileNode{
		ID:             "122",
		Name:           "test.txt",
		Path:           "/test",
		READMEUrl:      "123",
		IsFolder:       true,
		DownloadURL:    "https://test.com",
		LastModifyTime: time.Now(),
		Size:           100,
		Children:       nil,
		RefreshTime:    time.Now(),
		Password:       "123",
		PasswordURL:    "123",
		ParentID:       "122111111",
	}
)

func TestInsert(t *testing.T) {
	err := InsetFile(file)
	if err != nil {
		panic(err.Error())
	}
}

func TestBatchInsertFile(t *testing.T) {
	err := BatchInsertFile([]*FileNode{file})
	if err != nil {
		panic(err.Error())
	}
}

func TestDeleteFile(t *testing.T) {
	err := DeleteFile(file.ID)
	if err != nil {
		panic(err)
	}
}

func TestUpdateFile(t *testing.T) {
	err := UpdateFile(file)
	if err != nil {
		panic(err.Error())
	}
}

func TestFind(t *testing.T) {
	node, err := Find("122111111")
	if err != nil {
		panic(err.Error())
	}
	log.Infoln(node.Name)
}

func TestGetChildrenByID(t *testing.T) {
	nodes, err := GetChildrenByID("122111111")
	if err != nil {
		panic(err.Error())
	}
	for _, node := range nodes {
		log.Infoln(node.LastModifyTime)
	}
}

func TestFindByPath(t *testing.T) {
	node, err := FindByPath("/test")
	if err != nil {
		panic(err.Error())
	}
	log.Infoln(node.Name)
}

func TestFindByName(t *testing.T) {
	nodes, err := FindByName("test.txt")
	if err != nil {
		panic(err.Error())
	}
	for _, node := range nodes {
		log.Infoln(node.Name)
	}
}

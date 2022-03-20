package model

import (
	"database/sql"
	"time"

	log "github.com/sirupsen/logrus"

	//  sqlite驱动
	_ "modernc.org/sqlite"
)

var (
	db *sql.DB
)

func init() {
	var err error
	// 连接数据库

	db, err = sql.Open("sqlite", "data.db")
	if err != nil {
		log.Fatalf("打开数据库失败+ %v", err.Error())
		return
	}
	// 设置数据库连接数为1，防止出现database is locked的情况
	db.SetMaxOpenConns(1)
	// 创建表
	initTable()
}

// FileNode 存储的目录结构
// RefreshTime 表示最近一次的刷新时间
type FileNode struct {
	ID             string      `json:"id"`
	Name           string      `json:"name"`
	Path           string      `json:"path"`
	READMEUrl      string      `json:"-"`
	IsFolder       bool        `json:"is_folder"`
	DownloadURL    string      `json:"-"`
	LastModifyTime time.Time   `json:"last_modify_time"`
	Size           int64       `json:"size"`
	Children       []*FileNode `json:"-"`
	RefreshTime    time.Time   `json:"refresh_time"`
	// .password 内容
	Password    string `json:"-"`
	PasswordURL string `json:"-"`
	ParentID    string `json:"parent_id"`
}

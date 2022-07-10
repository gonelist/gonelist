package model

import (
	"database/sql"
	"time"

	log "github.com/sirupsen/logrus"
)

var (
	nodeChan chan *FileNode
)

func init() {
	nodeChan = make(chan *FileNode, 100)
	go insertFile()
}

// initTable
/**
 * @Description: 初始化表格
 */
func initTable() {
	_, err := db.Exec(`create table if not exists file
(
    id               TEXT not null
        constraint file_pk
            primary key,
    name             TEXT,
    path             TEXT,
    readme_url       TEXT,
    is_folder        integer default 0,
    download_url     TEXT,
    last_modify_time INTEGER,
    size             integer,
    password         TEXT,
    password_url     TEXT,
    parent_id        TEXT
);`)
	if err != nil {
		log.Errorln(err.Error())
		return
	}
	// 为path创建唯一索引
	_, err = db.Exec(`create unique index  path_index on file (path)`)
	if err != nil {
		log.Errorln(err.Error())
		return
	}
}

func insertFile() {
	for {
		node := <-nodeChan
		_, err := db.Exec(`insert into file  values (?,?,?,?,?,?,?,?,?,?,?);`,
			node.ID,
			node.Name,
			node.Path,
			node.READMEUrl,
			node.IsFolder,
			node.DownloadURL,
			node.LastModifyTime.Unix(),
			node.Size,
			node.Password,
			node.PasswordURL,
			node.ParentID)
		if err != nil {
			log.Errorln("数据库插入失败 " + err.Error())
		}
	}
}

var i = 0

// InsetFile
/**
 * @Description: 单条插入数据
 * @param node
 * @return error
 */
func InsetFile(node *FileNode) error {
	i++
	log.Infoln("添加了一个新的文件 ==》" + node.Path)
	nodeChan <- node
	return nil
}

// BatchInsertFile
/**
 * @Description: 批量插入数据
 * @param nodes
 * @return error
 */
func BatchInsertFile(nodes []*FileNode) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	// 方法退出时提交事务
	defer func(tx *sql.Tx) {
		err = tx.Commit()
		if err != nil {
			log.Errorln("提交事务失败" + err.Error())
		}
	}(tx)
	for _, node := range nodes {
		_, err := tx.Exec(`insert into file  values (?,?,?,?,?,?,?,?,?,?,?);`,
			node.ID,
			node.Name,
			node.Path,
			node.READMEUrl,
			node.IsFolder,
			node.DownloadURL,
			node.LastModifyTime.Unix(),
			node.Size,
			node.Password,
			node.PasswordURL,
			node.ParentID)
		if err != nil {
			log.Errorln("数据库插入失败 " + err.Error())
			// 插入错误时回滚数据
			err := tx.Rollback()
			if err != nil {
				return err
			}
			return err
		}
	}
	return err
}

// DeleteFile
/**
 * @Description: 删除一个文件
 * @param id
 * @return error
 */
func DeleteFile(id string) error {
	log.Errorln("检测到文件删除 ==》" + id)
	_, err := db.Exec(`delete from main.file where id = ?;`, id)
	if err != nil {
		return err
	}
	return err
}

// UpdateFile
/**
 * @Description: 更新文件数据
 * @param node
 * @return error
 */
func UpdateFile(node *FileNode) error {
	_, err := db.Exec(`update file set 
                name = ?,
            	path=?,
                readme_url=?,
                is_folder=?,
                download_url=?,
                last_modify_time=?,
                size=?,
                password=?,
                password_url=?,
                parent_id=? where id=?;`,
		node.Name, node.Path, node.READMEUrl, node.IsFolder, node.DownloadURL, node.LastModifyTime.Unix(),
		node.Size, node.Password, node.PasswordURL, node.ParentID, node.ID)
	if err != nil {
		log.Errorln("数据更新失败 " + err.Error())
		return err
	}
	return err
}

// Find
/**
 * @Description: 根据id查询单个文件
 * @param id
 * @return *FileNode
 * @return error
 */
func Find(id string) (*FileNode, error) {
	node := new(FileNode)
	var t int64
	err := db.QueryRow(`select * from file where id=?`, id).
		Scan(&node.ID, &node.Name, &node.Path, &node.READMEUrl, &node.IsFolder,
			&node.DownloadURL, &t, &node.Size, &node.Password,
			&node.PasswordURL, &node.ParentID)
	if err != nil {
		// log.Errorln("数据查找出现错误 " + err.Error())
		return nil, err
	}
	node.LastModifyTime = time.Unix(t, 0)
	return node, err
}

// GetChildrenByID
/**
 * @Description: 根据id获取该item下的子文件夹
 * @param id
 * @return []*FileNode
 * @return error
 */
func GetChildrenByID(id string) ([]*FileNode, error) {
	var (
		nodes []*FileNode
		t     int64
	)
	rows, err := db.Query("select * from file where parent_id=?", id)
	if err != nil {
		return nil, err
	}
	defer func(rows *sql.Rows) {
		_ = rows.Close()
	}(rows)
	for rows.Next() {
		node := new(FileNode)
		err := rows.Scan(&node.ID, &node.Name, &node.Path, &node.READMEUrl, &node.IsFolder,
			&node.DownloadURL, &t, &node.Size, &node.Password,
			&node.PasswordURL, &node.ParentID)
		if err != nil {
			log.Errorln("查询children出现错误 " + err.Error())
			return nil, err
		}
		node.LastModifyTime = time.Unix(t, 0)
		nodes = append(nodes, node)
	}
	return nodes, err
}

// FindByPath
/**
 * @Description: 根据path查找文件
 * @param path
 * @return *FileNode
 * @return error
 */
func FindByPath(path string) (*FileNode, error) {
	node := new(FileNode)
	var t int64
	err := db.QueryRow(`select * from file where path=?`, path).
		Scan(&node.ID, &node.Name, &node.Path, &node.READMEUrl, &node.IsFolder,
			&node.DownloadURL, &t, &node.Size, &node.Password,
			&node.PasswordURL, &node.ParentID)
	if err != nil {
		return nil, err
	}
	node.LastModifyTime = time.Unix(t, 0)
	return node, err
}

// FindByName
/**
 * @Description: 根据文件名查找信息
 * @param name
 * @return []*FileNode
 * @return error
 */
func FindByName(name string) ([]*FileNode, error) {
	var (
		nodes []*FileNode
		t     int64
	)
	rows, err := db.Query("select * from file where name=?", name)
	if err != nil {
		return nil, err
	}
	defer func(rows *sql.Rows) {
		_ = rows.Close()
	}(rows)
	for rows.Next() {
		node := new(FileNode)
		err := rows.Scan(&node.ID, &node.Name, &node.Path, &node.READMEUrl, &node.IsFolder,
			&node.DownloadURL, &t, &node.Size, &node.Password,
			&node.PasswordURL, &node.ParentID)
		if err != nil {
			log.Errorln("查询children出现错误 " + err.Error())
			return nil, err
		}
		node.LastModifyTime = time.Unix(t, 0)
		nodes = append(nodes, node)
	}
	return nodes, err
}

func Search(key string, path string) ([]*FileNode, error) {
	var (
		nodes []*FileNode
		t     int64
		rows  *sql.Rows
		err   error
	)
	if path == "" {
		rows, err = db.Query("select * from file where name like ?;", "%"+key+"%")
	} else {
		node, err := FindByPath(path)
		if err != nil {
			return nil, err
		}
		rows, err = db.Query("select * from file where name like ? and parent_id=?", "%"+key+"%", node.ID)
		if err != nil {
			return nil, err
		}
	}

	defer func(rows *sql.Rows) {
		_ = rows.Close()
	}(rows)
	for rows.Next() {
		node := new(FileNode)
		err := rows.Scan(&node.ID, &node.Name, &node.Path, &node.READMEUrl, &node.IsFolder,
			&node.DownloadURL, &t, &node.Size, &node.Password,
			&node.PasswordURL, &node.ParentID)
		if err != nil {
			log.Errorln("查询children出现错误 " + err.Error())
			return nil, err
		}
		node.LastModifyTime = time.Unix(t, 0)
		nodes = append(nodes, node)
	}
	return nodes, err
}

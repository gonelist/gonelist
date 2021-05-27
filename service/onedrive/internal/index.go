package internal

// Index 是索引的接口，用来实现全局搜索的功能
// 里面应该存放每个文件和目录的路径
type Index interface {
	// 直接设置所有内容
	SetData(data map[string][]Item)
	// 插入一条内容
	Insert(name string, item Item)
	// 插入多条数据
	InsertDatas(name string, item []Item)
	// 插入一组内容
	InsertDataMap(data map[string][]Item)
	// 清空内容
	Clear()
	// 搜索
	Search(key string) []Item
}

type Item struct {
	Path     string `json:"path"`
	IsFolder bool   `json:"is_folder"`
}

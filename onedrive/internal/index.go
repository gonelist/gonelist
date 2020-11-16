package internal

// Index 是索引的接口，用来实现全局搜索的功能
// 里面应该存放每个文件和目录的路径
type Index interface {
	// 直接设置所有内容
	SetData(data map[string]string)
	// 插入一条内容
	Insert(name, path string)
	// 插入一组内容
	InsertDataMap(data map[string]string)
	// 清空内容
	Clear()
	// 搜索
	Search(key string) []string
}

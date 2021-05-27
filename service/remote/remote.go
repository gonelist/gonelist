package remote

type Remote interface {
	// 初始化网盘内的索引数据
	Init()
	// 刷新网盘内索引数据
	Refresh()
	// 获取路径信息
	GetPath()
}

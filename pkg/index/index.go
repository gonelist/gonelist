package index

type Index interface {
	SetData(data []string)
	// 插入内容
	Insert(singleInfo string)
	InsertArray(singleInfoList []string)
	// 清空内容
	Clear()
	// 搜索
	Search(key string) []string
}

var IndexImpl Index

func init() {
	IndexImpl = &kmpIndex{}
}

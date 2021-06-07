package onedrive

import "gonelist/service/onedrive/internal"

func (t *Tree) SetData(data map[string][]internal.Item) {
}

func (t *Tree) Insert(name string, item internal.Item) {
}

func (t *Tree) InsertDatas(name string, item []internal.Item) {
}

func (t *Tree) InsertDataMap(data map[string][]internal.Item) {
}

func (t *Tree) Clear() {
}

// TODO
// 直接递归文件树查找
func (t *Tree) Search(key string) []internal.Item {
	var list []internal.Item
	root := t.root
	list = append(list, internal.Item{
		Path:     root.Path,
		IsFolder: true,
	})
	return list
}

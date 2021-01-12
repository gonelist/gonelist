package normal_index

import (
	"gonelist/onedrive/internal"
	"strings"
	"sync"
)

// TODO 可能出现重复文件名的情况
// data map[string][]string
type NIndex struct {
	sync.Mutex
	data map[string][]internal.Item
}

// 创建一个搜索的索引
func NewNIndex() *NIndex {
	index := &NIndex{}
	index.data = make(map[string][]internal.Item)
	return index
}

func (i *NIndex) SetData(data map[string][]internal.Item) {
	i.Lock()
	defer i.Unlock()

	if data == nil {
		i.data = make(map[string][]internal.Item)
	} else {
		i.data = data
	}
}

func (i *NIndex) Insert(name string, item internal.Item) {
	i.Lock()
	defer i.Unlock()

	i.data[name] = append(i.data[name], item)
}

func (i *NIndex) InsertDatas(name string, item []internal.Item) {
	i.Lock()
	defer i.Unlock()

	i.data[name] = append(i.data[name], item...)
}

func (i *NIndex) InsertDataMap(data map[string][]internal.Item) {
	i.Lock()
	defer i.Unlock()

	for name, items := range data {
		i.data[name] = append(i.data[name], items...)
	}
}

func (i *NIndex) Clear() {
	i.Lock()
	defer i.Unlock()

	i.data = make(map[string][]internal.Item)
}

func (i *NIndex) Search(key string) []internal.Item {
	ans := []internal.Item{}
	for name, items := range i.data {
		// 只返回 100 个结果
		if len(ans) >= 100 {
			break
		}
		if strings.Contains(name, key) {
			ans = append(ans, items...)
		}
	}
	return ans
}

package normal_index

import (
	"strings"
	"sync"
)

type NIndex struct {
	sync.Mutex
	data map[string]string
}

func (i *NIndex) SetData(data map[string]string) {
	i.Lock()
	defer i.Unlock()

	if data == nil {
		i.data = make(map[string]string)
	} else {
		i.data = data
	}
}

func (i *NIndex) Insert(name, path string) {
	i.Lock()
	defer i.Unlock()

	i.data[name] = path
}

func (i *NIndex) InsertDataMap(data map[string]string) {
	i.Lock()
	defer i.Unlock()

	for name, path := range data {
		i.data[name] = path
	}
}

func (i *NIndex) Clear() {
	i.Lock()
	defer i.Unlock()

	i.data = make(map[string]string)
}

func (i *NIndex) Search(key string) []string {
	ans := []string{}
	for name, path := range i.data {
		if strings.Contains(name, key) {
			ans = append(ans, path)
		}
	}
	return ans
}

package index

import (
	"strings"
	"sync"
)

type kmpIndex struct {
	sync.Mutex
	data []string
}

func (k *kmpIndex) SetData(data []string) {
	k.Lock()
	k.Unlock()

	k.data = data
}

func (k *kmpIndex) Insert(singleInfo string) {
	k.Lock()
	defer k.Unlock()

	k.data = append(k.data, singleInfo)
}

func (k *kmpIndex) InsertArray(singleInfoList []string) {
	k.Lock()
	defer k.Unlock()

	k.data = append(k.data, singleInfoList...)
}

func (k *kmpIndex) Clear() {
	k.Lock()
	defer k.Unlock()

	k.data = []string{}
}

func (k *kmpIndex) Search(key string) []string {
	ans := []string{}
	data := k.data
	for _, str := range data {
		if strings.Contains(str, key) {
			ans = append(ans, str)
		}
	}
	return ans
}

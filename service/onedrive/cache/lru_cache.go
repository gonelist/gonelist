package cache

import (
	"strings"
	"sync"

	"gonelist/conf"
	"gonelist/service/onedrive/model"
)

type LRUCache struct {
	cap    int                        // 缓存容量
	datas  map[string]*model.FileNode // map结构，存储数据
	list   *DoubleList                // 双向链表
	rwLock sync.RWMutex
}

var (
	Cache *LRUCache
)

func init() {
	InitCache(100)
}

func InitCache(cap int) {
	Cache = new(LRUCache)
	Cache.cap = cap
	Cache.list = NewDoubleList()
	Cache.datas = make(map[string]*model.FileNode)
}

// Put
/**
 * @Description:
 * @receiver cache
 * @param node
 * @return error
 */
func (cache *LRUCache) Put(node *model.FileNode) error {
	cache.rwLock.Lock()
	defer cache.rwLock.Unlock()

	if fileNode, ok := cache.datas[node.Path]; ok {
		cache.list.MoveToHead(fileNode.Path)
		return nil
	}

	// 缓存中容量超过上限
	if cache.list.Size >= cache.cap {
		// 删除链表尾部节点
		k := cache.list.RemoveOneNodeByTail()
		delete(cache.datas, k)
	}
	cache.list.InsertList(node.ID)
	cache.datas[node.ID] = node
	return nil
}

// Get
/**
 * @Description: 从缓存中获取数据
 * @receiver cache
 * @param key
 * @return *model.FileNode
 * @return bool
 */
func (cache *LRUCache) Get(key string) (*model.FileNode, bool) {
	// cache.rwLock.Lock()
	// defer cache.rwLock.Unlock()
	if conf.UserSet.Onedrive.FolderSub != "/" {
		key = conf.UserSet.Onedrive.FolderSub + key
		key = strings.TrimRight(key, "/")
	}
	fileNode, ok := cache.datas[key]
	if ok {
		cache.list.MoveToHead(fileNode.Path)
		return fileNode, true
	}
	node, err := model.FindByPath(key)
	if err != nil {
		return nil, false
	}
	_ = cache.Put(node)
	return node, true
}

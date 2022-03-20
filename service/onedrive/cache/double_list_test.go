package cache

import (
	"log"
	"testing"
)

func TestDoubleList_InsertList(t *testing.T) {
	list := NewDoubleList()
	list.InsertList("1")
	list.InsertList("2")
	list.InsertList("3")
	list.InsertList("4")
	log.Println(list)
}

func TestDoubleList_RemoveOneNodeByTail(t *testing.T) {
	list := NewDoubleList()
	list.InsertList("1")
	list.InsertList("2")
	list.InsertList("3")
	list.InsertList("4")

	list.RemoveOneNodeByTail()
	list.RemoveOneNodeByTail()
}

func TestDoubleList_RemoveNode(t *testing.T) {
	list := NewDoubleList()
	list.InsertList("1")
	list.InsertList("2")
	list.InsertList("3")
	list.InsertList("4")

	list.RemoveNode("2")
	list.RemoveNode("1")
	list.RemoveNode("4")
}

func TestDoubleList_MoveToHead(t *testing.T) {
	list := NewDoubleList()
	list.InsertList("1")
	list.InsertList("2")
	list.InsertList("3")
	list.InsertList("4")

	list.MoveToHead("2")
	list.MoveToHead("1")
}

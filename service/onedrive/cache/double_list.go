package cache

type Node struct {
	prev *Node
	next *Node
	key  string
}

type DoubleList struct {
	Head *Node
	Tail *Node
	Size int
}

// NewDoubleList
/**
 * @Description: 初始化链表
 * @return *DoubleList
 */
func NewDoubleList() *DoubleList {
	d := new(DoubleList)
	d.Head = nil
	d.Tail = nil
	d.Size = 0
	return d
}

func (l *DoubleList) MoveToHead(data string) {
	if l.Head.key == data {
		return
	}
	l.RemoveNode(data)
	l.InsertList(data)
}

// RemoveNode
/**
 * @Description: 根据值删除节点
 * @receiver l
 * @param data
 */
func (l *DoubleList) RemoveNode(data string) {
	if l.Size == 0 {
		return
	}
	temp := l.Head
	flag := false
	for {
		if temp.next == nil {
			if temp.key == data {
				flag = true
			}
			break
		} else {
			if temp.key == data {
				flag = true
				break
			}
		}
		temp = temp.next
	}
	if flag {
		if temp.next != nil {
			temp.next.prev = temp.prev
		} else {
			l.Tail = temp.prev
		}
		if temp.prev != nil {
			temp.prev.next = temp.next
		} else {
			l.Head = temp.next
		}
		l.Size--
	}
}

// InsertList
/**
 * @Description: 从链表头部插入数据
 * @receiver l
 * @param data
 */
func (l *DoubleList) InsertList(data string) {
	n := new(Node)
	n.key = data
	// 链表为空时
	if l.Size < 1 {
		l.Head = n
		l.Tail = n
	} else {
		// 链表的长度>=1时
		firstNode := l.Head
		n.next = firstNode
		firstNode.prev = n
		n.prev = nil
		l.Head = n
	}
	l.Size++
}

// RemoveOneNodeByTail
/**
 * @Description: 删除链表尾部节点
 * @receiver l
 */
func (l *DoubleList) RemoveOneNodeByTail() string {
	result := l.Tail.key
	tempNode := l.Tail.prev
	tempNode.next = nil
	l.Tail = tempNode
	l.Size--
	return result

}

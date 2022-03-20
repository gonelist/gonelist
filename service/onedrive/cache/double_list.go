package cache

type Node struct {
	prev *Node
	next *Node
	key  string
}

type DoubleList struct {
	head *Node
	tail *Node
	size int
}

// NewDoubleList
/**
 * @Description: 初始化链表
 * @return *DoubleList
 */
func NewDoubleList() *DoubleList {
	d := new(DoubleList)
	d.head = nil
	d.tail = nil
	d.size = 0
	return d
}

func (l *DoubleList) MoveToHead(data string) {
	if l.head.key == data {
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
	if l.size == 0 {
		return
	}
	temp := l.head
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
			l.tail = temp.prev
		}
		if temp.prev != nil {
			temp.prev.next = temp.next
		} else {
			l.head = temp.next
		}
		l.size--
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
	if l.size < 1 {
		l.head = n
		l.tail = n
	} else {
		// 链表的长度>=1时
		firstNode := l.head
		n.next = firstNode
		firstNode.prev = n
		n.prev = nil
		l.head = n
	}
	l.size++
}

// RemoveOneNodeByTail
/**
 * @Description: 删除链表尾部节点
 * @receiver l
 */
func (l *DoubleList) RemoveOneNodeByTail() string {
	result := l.tail.key
	tempNode := l.tail.prev
	tempNode.next = nil
	l.tail = tempNode
	l.size--
	return result

}

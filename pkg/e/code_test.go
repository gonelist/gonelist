package e

import (
	"fmt"
	"sort"
	"testing"
)

func TestGetMsg(t *testing.T) {
	params := MsgFlags
	// 按字母顺序遍历map
	keys := make([]int, 0)
	for k, _ := range params {
		keys = append(keys, k)
	}
	sort.Ints(keys)
	for _, k := range keys {
		fmt.Println("|", k, "|", params[k], "|")
		//handler(k, params[k])
	}
}

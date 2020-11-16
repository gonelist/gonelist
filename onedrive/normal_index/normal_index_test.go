package normal_index

import (
	"gonelist/onedrive/internal"
	"gopkg.in/go-playground/assert.v1"
	"testing"
)

func TestNIndex_InsertArray(t *testing.T) {
	var (
		index internal.Index
		ans   []string
	)
	index = &NIndex{}
	example := make(map[string]string)
	example["这是一个目录"] = "/测试/这是一个目录"
	example["123"] = "/public/123"
	example["a012,我们去大草原的湖边"] = "/a012,我们去大草原的湖边"
	index.SetData(example)

	ans = index.Search("我")
	assert.Equal(t, ans, []string{"/a012,我们去大草原的湖边"})

	ans = index.Search("/")
	assert.Equal(t, ans, []string{})

	ans = index.Search("1")
	t.Log(ans)
	assert.Equal(t, ans, []string{"/public/123", "/a012,我们去大草原的湖边"})

	ans = index.Search("123")
	assert.Equal(t, ans, []string{"/public/123"})

}

// TODO test multi goroutine concurrence
func TestNIndex_SetData(t *testing.T) {
	//var (
	//	wg sync.WaitGroup
	//	k  NIndex
	//)
	//
	//k.SetData([]string{"1", "2", "3"})
	//go func() {
	//	k.Search()
	//}()
}

package index

import (
	"gopkg.in/go-playground/assert.v1"
	"testing"
)

func TestKmpIndex_InsertArray(t *testing.T) {
	var (
		index   Index
		ans     []string
		example []string
	)
	index = &kmpIndex{}
	example = []string{"/测试/这是一个目录", "/public/123", "/a012,我们去大草原的湖边"}
	index.InsertArray(example)

	ans = index.Search("我")
	assert.Equal(t, ans, []string{example[2]})

	ans = index.Search("/")
	assert.Equal(t, ans, example)

	ans = index.Search("1")
	assert.Equal(t, ans, []string{example[1], example[2]})

	ans = index.Search("123")
	assert.Equal(t, ans, []string{example[1]})

}

// TODO test multi goroutine concurrence
func TestKmpIndex_SetData(t *testing.T) {
	//var (
	//	wg sync.WaitGroup
	//	k  kmpIndex
	//)
	//
	//k.SetData([]string{"1", "2", "3"})
	//go func() {
	//	k.Search()
	//}()
}

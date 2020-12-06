package normal_index

import (
	internal "gonelist/onedrive/internal"
	"gopkg.in/go-playground/assert.v1"
	"testing"
)

func TestNIndex_InsertArray(t *testing.T) {
	var (
		index *NIndex
		ans   []internal.Item
	)
	index = NewNIndex()
	example := make(map[string][]internal.Item)
	example["这是一个目录"] = []internal.Item{
		{
			Path:     "/测试/这是一个目录",
			IsFolder: true,
		},
	}
	example["123"] = []internal.Item{
		{
			Path:     "/public/123",
			IsFolder: false,
		},
	}
	example["a012,我们去大草原的湖边"] = []internal.Item{
		{
			Path:     "/a012,我们去大草原的湖边",
			IsFolder: false,
		},
	}
	index.SetData(example)

	ans = index.Search("我")
	assert.Equal(t, ans, []internal.Item{
		{
			Path:     "/a012,我们去大草原的湖边",
			IsFolder: false,
		},
	})

	ans = index.Search("/")
	assert.Equal(t, ans, []internal.Item{})

	ans = index.Search("1")
	t.Log(ans)
	assert.Equal(t, ans, []internal.Item{
		{
			Path:     "/public/123",
			IsFolder: false,
		},
		{
			Path:     "/a012,我们去大草原的湖边",
			IsFolder: false,
		},
	})

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

package onedrive

import (
	"encoding/json"
	"fmt"
	"GOIndex/pkg/file"
	"testing"
)

func TestOnedriveGetPath(t *testing.T) {
	data := file.ReadFromFile("../example/root.json")
	var ans Answer
	json.Unmarshal([]byte(data), &ans)
	fmt.Println(ans)
}

func TestGetAnsErr(t *testing.T) {
	data := file.ReadFromFile("../example/InvalidToken.json")
	var ans Answer
	json.Unmarshal([]byte(data), &ans)
	fmt.Println(ans)
}

func TestCheckAnswerValid(t *testing.T) {
	// 如果都用 ans，第二个的 Error 会是第一个的内容
	var data string
	var ans1 Answer
	var ans2 Answer
	var valid error

	data = file.ReadFromFile("../example/InvalidToken.json")
	json.Unmarshal([]byte(data), &ans1)
	valid = CheckAnswerValid(ans1, "/example/InvalidToken.json")
	fmt.Println(valid)
	data = file.ReadFromFile("../example/root.json")
	json.Unmarshal([]byte(data), &ans2)
	valid = CheckAnswerValid(ans2, "/example/root.json")
	fmt.Println(valid)
}

func TestCacheGetPathList(t *testing.T) {
	var data string
	var filetree *FileNode

	data = file.ReadFromFile("../example/filetree.json")
	json.Unmarshal([]byte(data), &filetree)

	FileTree = filetree

	child, err := CacheGetPathList("/test")
	if err != nil {
		fmt.Println(err)
	}
	for _, item := range child {
		fmt.Println(item)
	}
}

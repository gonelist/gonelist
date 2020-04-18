package onedrive

import (
	"encoding/json"
	"fmt"
	"gonelist/conf"
	"gonelist/mg_auth"
	"gonelist/pkg/file"
	"log"
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

	root, err := CacheGetPathList("/ttt")
	if err != nil {
		fmt.Println(err)
		return
	}
	if root.IsFolder == false {
		// 是一个文件
		b, _ := json.Marshal(root)
		fmt.Println(string(b))
	} else {
		for _, item := range root.Children {
			fmt.Println(item)
		}
	}
}

func TestConvertReturnNode(t *testing.T) {
	var data string
	var filetree *FileNode

	if err := conf.LoadUserConfig("../example/config.json"); err != nil {
		log.Fatal(err)
	}

	mg_auth.SetUserInfo(conf.UserSet)

	data = file.ReadFromFile("../example/filetree.json")
	json.Unmarshal([]byte(data), &filetree)

	reNode := ConvertReturnNode(filetree)
	fmt.Println(reNode)
}

// 测试一个目录下，上千文件的情况
func TestThousand(t *testing.T) {
	data := file.ReadFromFile("../example/pdf.json")
	var ans Answer
	var ans1 Answer
	json.Unmarshal([]byte(data), &ans)
	fmt.Println(len(ans.Value))
	data1 := file.ReadFromFile("../example/pdf1.json")
	json.Unmarshal([]byte(data1), &ans1)
	fmt.Println(len(ans1.Value))
}

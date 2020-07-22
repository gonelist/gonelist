package onedrive

import (
	"encoding/json"
	"gonelist/conf"
	"gonelist/pkg/file"
	"log"
	"testing"
)

const example = "../example/"

func TestOnedriveGetPath(t *testing.T) {
	data, _ := file.ReadFromFile(example + "root.json")
	var ans Answer
	json.Unmarshal(data, &ans)
	t.Log(ans)
}

func TestGetAnsErr(t *testing.T) {
	data, _ := file.ReadFromFile(example + "InvalidToken.json")
	var ans Answer
	json.Unmarshal(data, &ans)
	t.Log(ans)
}

func TestCheckAnswerValid(t *testing.T) {
	// 如果都用 ans，第二个的 Error 会是第一个的内容
	var data []byte
	var ans1 Answer
	var ans2 Answer
	var valid error

	data, _ = file.ReadFromFile(example + "InvalidToken.json")
	json.Unmarshal(data, &ans1)
	valid = CheckAnswerValid(ans1, "/example/InvalidToken.json")
	t.Log(valid)
	data, _ = file.ReadFromFile(example + "root.json")
	json.Unmarshal(data, &ans2)
	valid = CheckAnswerValid(ans2, "/example/root.json")
	t.Log(valid)
}

func TestCacheGetPathList(t *testing.T) {
	if err := conf.LoadUserConfig("../config.json"); err != nil {
		log.Fatal(err)
	}

	var data []byte
	var filetree *FileNode

	data, _ = file.ReadFromFile(example + "filetree.json")
	json.Unmarshal(data, &filetree)

	FileTree.SetRoot(filetree)
	FileTree.SetLogin(true)

	root, err := CacheGetPathList("/ttt")
	if err != nil {
		t.Log(err)
		return
	}
	if root.IsFolder == false {
		// 是一个文件
		b, _ := json.Marshal(root)
		t.Log(string(b))
	} else {
		for _, item := range root.Children {
			t.Log(item)
		}
	}
}

func TestConvertReturnNode(t *testing.T) {
	var data []byte
	var filetree *FileNode

	if err := conf.LoadUserConfig(example + "config.json"); err != nil {
		log.Fatal(err)
	}

	SetOnedriveInfo(conf.UserSet)

	data, _ = file.ReadFromFile(example + "filetree.json")
	json.Unmarshal(data, &filetree)

	reNode := ConvertReturnNode(filetree)
	t.Log(reNode)
}

// 测试一个目录下，上千文件的情况
func TestThousand(t *testing.T) {
	data, _ := file.ReadFromFile(example + "pdf.json")
	var ans Answer
	var ans1 Answer
	json.Unmarshal(data, &ans)
	t.Log(len(ans.Value))
	data1, _ := file.ReadFromFile(example + "pdf1.json")
	json.Unmarshal([]byte(data1), &ans1)
	t.Log(len(ans1.Value))
}

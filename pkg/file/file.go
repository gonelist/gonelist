package file

import (
	"fmt"
	"io/ioutil"
)

// 从某个文件读取内容
func ReadFromFile(filename string) string {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Println(err)
		//return
	}
	return string(data)
}
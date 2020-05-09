package markdown

import (
	"errors"
	"github.com/russross/blackfriday"
	"gonelist/pkg/file"
)

// 输入文件路径，给出对应的内容
func MarkdownToHTML(filePath string) ([]byte, error) {
	var err error

	if !file.IsExistFile(filePath) {
		return nil, errors.New("文件不存在")
	}
	input, err := file.ReadFromFile(filePath)
	if err != nil {
		return nil, err
	}
	output := blackfriday.MarkdownBasic(input)
	return output, nil
}

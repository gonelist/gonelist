package markdown

import (
	"github.com/russross/blackfriday"
	"gonelist/pkg/file"
)

// 输入文件路径，给出对应的内容
func MarkdownToHTML(filePath string) ([]byte, error) {
	input, err := file.ReadFromFile(filePath)
	if err != nil {
		return nil, err
	}
	output := blackfriday.MarkdownBasic(input)
	return output, nil
}

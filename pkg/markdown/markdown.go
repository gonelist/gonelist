package markdown

import (
	"github.com/russross/blackfriday"
	"gonelist/pkg/file"
)

// 输入文件路径，给出对应的内容
func MarkdownToHTML(filePath string) []byte {
	input := file.ReadFromFile(filePath)
	output := blackfriday.MarkdownBasic(input)
	return output
}

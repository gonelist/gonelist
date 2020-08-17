package markdown

import (
	"errors"
	"github.com/russross/blackfriday/v2"
	"gonelist/pkg/file"
)

// 输入 []byte，得到结果
func MarkdownToHTMLByBytes(input []byte) []byte {
	// blackfriday.HardLineBreak，将 markdown 中的换行直接替换为 <br />
	output := blackfriday.Run(input,
		blackfriday.WithExtensions(blackfriday.HardLineBreak|blackfriday.Autolink|blackfriday.FencedCode))
	return output
}

// 输入文件路径，给出对应的内容
func MarkdownToHTMLByFile(filePath string) ([]byte, error) {
	if !file.IsExistFile(filePath) {
		return nil, errors.New("文件不存在")
	}

	input, err := file.ReadFromFile(filePath)
	if err != nil {
		return nil, err
	}
	output := MarkdownToHTMLByBytes(input)
	return output, nil
}

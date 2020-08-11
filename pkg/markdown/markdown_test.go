package markdown

import (
	"testing"
)

const example = "../../example/"

func TestMarkdownToHTML(t *testing.T) {
	output, err := MarkdownToHTMLByFile(example + "README.md")
	if err != nil {
		t.Log(err)
	}
	t.Log(string(output))
}

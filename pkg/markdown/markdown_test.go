package markdown

import (
	"fmt"
	"testing"
)

const example = "../../example/"

func TestMarkdownToHTML(t *testing.T) {
	output := MarkdownToHTML(example + "README.md")
	fmt.Println(string(output))
}

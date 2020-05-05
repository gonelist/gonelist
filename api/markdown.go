package api

import (
	"github.com/gin-gonic/gin"
	"gonelist/pkg/markdown"
	"net/http"
)

// 返回 README 的内容
func GetREADME(c *gin.Context) {
	filePath := "./README.md"
	output := markdown.MarkdownToHTML(filePath)

	c.String(http.StatusOK, string(output))
}

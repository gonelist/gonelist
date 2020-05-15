package api

import (
	"github.com/gin-gonic/gin"
	"gonelist/pkg/app"
	"gonelist/pkg/e"
	"gonelist/pkg/markdown"
	"net/http"
)

// 返回 README 的内容
func GetREADME(c *gin.Context) {
	filePath := "./README.md"

	output, err := markdown.MarkdownToHTML(filePath)
	if err != nil {
		app.Response(c, http.StatusOK, e.ITEM_NOT_FOUND, nil)
	} else {
		app.Response(c, http.StatusOK, e.SUCCESS, string(output))
	}
}

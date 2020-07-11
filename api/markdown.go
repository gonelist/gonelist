package api

import (
	"github.com/gin-gonic/gin"
	"gonelist/onedrive"
	"gonelist/pkg/app"
	"gonelist/pkg/e"
	"gonelist/pkg/markdown"
	"net/http"
)

// 返回 README 的内容
func GetREADME(c *gin.Context) {
	path := c.Query("path")

	// 从 cache 中获取对应的 input
	input, err := onedrive.GetREADMEInCache(path)
	if err != nil {
		app.Response(c, http.StatusOK, e.CACHE_NOT_FIND, nil)
		return
	}

	output, err := markdown.MarkdownToHTMLByBytes(input)
	if err != nil {
		app.Response(c, http.StatusOK, e.ITEM_NOT_FOUND, nil)
	} else {
		app.Response(c, http.StatusOK, e.SUCCESS, string(output))
	}
}

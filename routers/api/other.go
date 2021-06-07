package api

import (
	"github.com/gin-gonic/gin"
	"gonelist/conf"
	"gonelist/pkg/app"
	"gonelist/pkg/e"
	"gonelist/service/onedrive"
	"net/http"
)

// 返回 README 的内容
func GetREADME(c *gin.Context) {
	path := c.Query("path")

	if path == "" {
		path = "/"
	}
	// 从 cache 中获取对应的 readme 的内容
	readmeBytes, err := onedrive.GetREADMEInCache(path)
	if err != nil {
		app.Response(c, http.StatusOK, e.CACHE_NOT_FIND, nil)
	} else {
		app.Response(c, http.StatusOK, e.SUCCESS, string(readmeBytes))
	}
}

// 搜索功能
func Search(c *gin.Context) {
	key := c.Query("key")

	if key == "" {
		app.Response(c, http.StatusOK, e.ITEM_NOT_FOUND, nil)
		return
	}

	ans := onedrive.FileTree.Search(key)
	app.Response(c, http.StatusOK, e.SUCCESS, ans)
}

// 版本和网盘信息
func Info(c *gin.Context) {
	ans := make(map[string]string)
	ans["name"] = conf.UserSet.Name
	ans["version"] = conf.UserSet.Version
	app.Response(c, http.StatusOK, e.SUCCESS, ans)
}

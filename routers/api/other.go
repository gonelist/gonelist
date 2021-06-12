package api

import (
	"github.com/gin-gonic/gin"
	"gonelist/conf"
	"gonelist/pkg/app"
	"gonelist/pkg/e"
	"gonelist/service/onedrive"
	"net/http"
)

// swagger:operation GET /README?path= readme
// ---
// summary: 获取 README 的内容
// description: 获取 README 的内容
// parameters:
// - name: path
//   in: path
//   description: 地址id
//   type: string
//   required: true
// responses:
//   200: repoResp
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

// swagger:operation GET /info info
// ---
// summary: 获取网盘信息
// description: 获取网盘信息
// parameters:
//
// responses:
//   200: repoResp
func Info(c *gin.Context) {
	ans := make(map[string]string)
	ans["name"] = conf.UserSet.Name
	ans["version"] = conf.UserSet.Version
	app.Response(c, http.StatusOK, e.SUCCESS, ans)
}

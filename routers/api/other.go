package api

import (
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"

	"gonelist/conf"
	"gonelist/pkg/app"
	"gonelist/pkg/e"
	"gonelist/service/onedrive"
	"gonelist/service/onedrive/model"

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

	//ans := onedrive.FileTree.Search(key)
	nodes, err := model.Search(key)
	if err != nil {
		log.Errorln(err.Error())
		return
	}
	app.Response(c, http.StatusOK, e.SUCCESS, nodes)
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
	ans := make(map[string]interface{})
	ans["name"] = conf.UserSet.Name
	ans["page_title"] = conf.UserSet.PageTitle
	ans["upload"] = conf.UserSet.Admin.EnableWrite
	ans["version"] = conf.UserSet.Version
	app.Response(c, http.StatusOK, e.SUCCESS, ans)
}

// UpdatePermission
/**
 * @Description: 提升客户端的权限为管理员权限
 * @return gin.HandlerFunc
 */
func UpdatePermission() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if conf.UserSet.Admin.EnableWrite {
			if ctx.Query("secret") == conf.UserSet.Admin.Secret {
				app.Response(ctx, http.StatusOK, e.SUCCESS, nil)
			} else {
				app.Response(ctx, http.StatusOK, e.ACCESS_TOKEN_ERROR, nil)
			}
		} else {
			app.Response(ctx, http.StatusOK, e.ACCESS_TOKEN_ERROR, nil)
		}
	}
}

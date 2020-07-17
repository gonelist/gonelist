package api

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"gonelist/onedrive"
	"gonelist/pkg/app"
	"gonelist/pkg/e"
	"net/http"
	"strings"
)

// 测试接口，从 MG 获取整个树结构
func MGGetFileTree(c *gin.Context) {
	root, err := onedrive.GetAllFiles()
	if err != nil {
		log.Warn("请求 graph.microsoft.com 错误")
		app.Response(c, http.StatusOK, e.MG_ERROR, e.GetMsg(e.MG_ERROR))
		return
	}

	str, _ := json.Marshal(root)
	log.Debug("*root", string(str))

	app.Response(c, http.StatusOK, e.SUCCESS, root)
}

// 获取对应路径的文件
func CacheGetPath(c *gin.Context) {
	oPath := c.Query("path")
	pass := c.GetHeader("pass")

	root, err := onedrive.CacheGetPathList(oPath)
	if err != nil {
		app.Response(c, http.StatusOK, e.ITEM_NOT_FOUND, e.GetMsg(e.ITEM_NOT_FOUND))
	} else if root.Password != "" && pass != root.Password {
		app.Response(c, http.StatusOK, e.PASS_ERROR, e.GetMsg(e.PASS_ERROR))
	} else {
		app.Response(c, http.StatusOK, e.SUCCESS, root)
	}
}

// 分享文件下载链接
func Download(c *gin.Context) {
	filePath := c.Param("path")

	// 屏蔽 .password 文件的下载
	list := strings.Split(filePath, "/")
	if list[len(list)-1] == ".password" {
		app.Response(c, http.StatusOK, e.PASSWORD_FORBIT_DOWNLOAD, e.GetMsg(e.PASSWORD_FORBIT_DOWNLOAD))
		c.Abort()
	}

	downloadURL, err := onedrive.GetDownloadUrl(filePath)
	if err != nil {
		app.Response(c, http.StatusOK, e.ITEM_NOT_FOUND, e.GetMsg(e.ITEM_NOT_FOUND))
	} else {
		c.Redirect(http.StatusFound, downloadURL)
		c.Abort()
	}
}

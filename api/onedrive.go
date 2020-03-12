package api

import (
	"GOIndex/onedrive"
	"GOIndex/pkg/app"
	"GOIndex/pkg/e"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

// 测试接口，从 MG 获取树结构
func MGGetFileTree(c *gin.Context) {
	root := onedrive.GetAllFiles()

	str, _ := json.Marshal(root)
	fmt.Println("*root", string(str))

	app.Response(c, http.StatusOK, e.SUCCESS, root)
}

// 获取对应路径的文件
func CacheGetPath(c *gin.Context) {
	oPath := c.Query("path")

	root, err := onedrive.CacheGetPathList(oPath)
	if err != nil {
		app.Response(c, http.StatusOK, e.ITEM_NOT_FOUND, nil)
	} else {
		app.Response(c, http.StatusOK, e.SUCCESS, root)
	}
}

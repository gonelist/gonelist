package api

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"

	"gonelist/conf"
	"gonelist/pkg/app"
	"gonelist/pkg/e"
	"gonelist/service/onedrive"
)

// 测试接口，从 MG 获取整个树结构
func MGGetFileTree(c *gin.Context) {
	root, err := onedrive.GetAllFiles()
	if err != nil {
		log.Warn("请求 graph.microsoft.com 错误")
		app.Response(c, http.StatusOK, e.MG_ERROR, nil)
		return
	}
	onedrive.RefreshREADME()

	str, _ := json.Marshal(root)
	log.Debug("*root", string(str))

	app.Response(c, http.StatusOK, e.SUCCESS, "已刷新缓存")
}

// 获取对应路径的文件
func CacheGetPath(c *gin.Context) {
	oPath := c.Query("path")

	root, err := onedrive.CacheGetPathList(oPath)
	// 如果没有找到文件则返回 404
	if err != nil {
		app.Response(c, http.StatusNotFound, e.ITEM_NOT_FOUND, nil)
	} else if root == nil {
		app.Response(c, http.StatusNotFound, e.LOAD_NOT_READY, nil)
	} else {
		app.Response(c, http.StatusOK, e.SUCCESS, root)
	}
}

// CheckUploadSecret
/**
 * @Description: 用于检查文件上传时的token
 * @return gin.HandlerFunc
 */
func CheckUploadSecret() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		secret := ctx.Query("secret")
		if secret != conf.UserSet.Onedrive.UploadSecret {
			app.Response(ctx, http.StatusOK, e.SECRET_ERROR, nil)
			ctx.Abort()
			return
		}
		ctx.Next()
	}
}

// Upload
/**
 * @Description: 上传文件
 * @return gin.HandlerFunc
 */
func Upload() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if !conf.UserSet.Server.EnableUpload {
			app.Response(ctx, 403, 1403, "the api not open")
			return
		}
		path := ctx.Query("path")
		file, err := ctx.FormFile("file")
		if err != nil {
			app.Response(ctx, http.StatusBadGateway, e.INVALID_PARAMS, e.MsgFlags[e.INVALID_PARAMS])
			return
		}
		// 检查文件大小
		if file.Size > 4194304 {
			app.Response(ctx, http.StatusBadGateway, e.INVALID_PARAMS, "文件大小大于4MB")
		}
		f, err := file.Open()
		if err != nil {
			return
		}
		content, err := io.ReadAll(f)
		if err != nil {
			return
		}
		err = f.Close()
		if err != nil {
			return
		}
		log.Infoln("开始上传文件：文件名==》"+file.Filename+" 文件大小===》", file.Size/1024/1024, "Mb", " 上传路径==》"+path)
		err = onedrive.Upload(path, file.Filename, content)
		if err != nil {
			app.Response(ctx, http.StatusOK, e.SUCCESS, nil)
			return
		}
		ctx.JSON(200, nil)
	}
}

// 分享文件下载链接
func Download(c *gin.Context) {
	filePath := c.Param("path")
	// 屏蔽 .password 文件的下载
	list := strings.Split(filePath, "/")
	if list[len(list)-1] == ".password" {
		app.Response(c, http.StatusOK, e.PASSWORD_FORBIT_DOWNLOAD, nil)
		c.Abort()
	}

	downloadURL, err := onedrive.GetDownloadUrl(filePath)
	log.Info(downloadURL)
	if err != nil || downloadURL == "" {
		app.Response(c, http.StatusOK, e.ITEM_NOT_FOUND, nil)
	} else {
		c.Redirect(http.StatusFound, downloadURL)
		c.Abort()
	}
}

package api

import (
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
	err := onedrive.RefreshOnedriveAll()
	if err != nil {
		return
	}

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

// 创建文件夹
func MkDir() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		path := ctx.Query("path")
		folderName := ctx.Query("folder_name")
		err := onedrive.Mkdir(path, folderName)
		if err != nil {
			app.Response(ctx, 403, 403, "")
			return
		}
		app.Response(ctx, http.StatusOK, e.SUCCESS, "")
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

// swagger:operation POST /onedrive/upload info
// ---
// summary: 上传文件
// description: 上传一个文件，目前仅支持单文件
// parameters:
// 		- name: path
//   	in: path
//   	description: 地址id
//   	type: string
//   	required: true
//
// responses:
//   200: repoResp
//   403： the api not open
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
		// 小文件直接上传
		if file.Size < 4194304 {
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
			app.Response(ctx, http.StatusOK, e.SUCCESS, nil)
			// app.Response(ctx, http.StatusBadGateway, e.INVALID_PARAMS, "文件大小大于4MB")
		} else {
			// 大文件通过session进行分片上传
			uploader := onedrive.NewUploader()
			// 创建一个上传session,获取到上传url
			err := uploader.CreateSession(path, file.Filename, file.Size)
			if err != nil {
				return
			}
			data, err := file.Open()
			if err != nil {
				return
			}
			temp := make([]byte, conf.UserSet.Onedrive.UploadSliceSize*327680)
			_, err = io.CopyBuffer(uploader, data, temp)
			if err != nil {
				return
			}
			err = data.Close()
			if err != nil {
				return
			}
			app.Response(ctx, http.StatusOK, e.SUCCESS, nil)
		}
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
	if err != nil || downloadURL == "" {
		app.Response(c, http.StatusOK, e.ITEM_NOT_FOUND, nil)
	} else {
		c.Redirect(http.StatusFound, downloadURL)
		c.Abort()
	}
}

package api

import (
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"

	"gonelist/conf"
	"gonelist/pkg/app"
	"gonelist/pkg/e"
	"gonelist/service/local"
)

// LocalDownload
/**
 * @Description: 本地文件下载
 * @return gin.HandlerFunc
 */
func LocalDownload() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if conf.UserSet.Local.Enable {
			path := ctx.Param("path")
			if strings.HasPrefix(path, "/"+conf.UserSet.Local.Name) {
				ctx.Header("Content-Type", "application/octet-stream")
				ctx.Header("Content-Disposition", "attachment; filename="+strings.Split(path, "/")[len(strings.Split(path, "/"))-1])
				ctx.Header("Content-Transfer-Encoding", "binary")
				ctx.File(local.HandlePath(path))
				ctx.Abort()
				return
			}
		}
		ctx.Next()
	}
}

// LocalMkdir
/**
 * @Description: 本地文件夹创建
 * @return gin.HandlerFunc
 */
func LocalMkdir() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if conf.UserSet.Local.Enable {
			path := ctx.Query("path")
			folderName := ctx.Query("folder_name")
			if strings.HasPrefix(path, "/"+conf.UserSet.Local.Name) {
				err := os.Mkdir(local.HandlePath(path)+"/"+folderName, 0666)
				if err != nil {
					app.Response(ctx, http.StatusBadGateway, e.SUCCESS, err)
					ctx.Abort()
					return
				}
				app.Response(ctx, http.StatusOK, e.SUCCESS, nil)
				ctx.Abort()
				return
			}
		}
		ctx.Next()
	}
}

// LocalUpload
/**
 * @Description: 本地文件上传
 * @return gin.HandlerFunc
 */
func LocalUpload() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if conf.UserSet.Local.Enable {
			path := ctx.Query("path")
			file, err := ctx.FormFile("file")
			if err != nil {
				ctx.Abort()
				return
			}
			if strings.HasPrefix(path, "/"+conf.UserSet.Local.Name) {
				err := ctx.SaveUploadedFile(file, local.HandlePath(path)+"/"+file.Filename)
				if err != nil {
					app.Response(ctx, http.StatusBadGateway, e.SUCCESS, err)
					ctx.Abort()
					return
				}
				app.Response(ctx, http.StatusOK, e.SUCCESS, nil)
				ctx.Abort()
				return
			}
		}
		ctx.Next()
	}
}

func LocalDelete() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if conf.UserSet.Local.Enable {
			path := ctx.Query("path")
			if strings.HasPrefix(path, "/"+conf.UserSet.Local.Name) {
				err := os.Remove(local.HandlePath(path))
				if err != nil {
					app.Response(ctx, http.StatusBadGateway, e.SUCCESS, err)
					ctx.Abort()
					return
				}
				app.Response(ctx, http.StatusOK, e.SUCCESS, nil)
				ctx.Abort()
				return
			}
		}
		ctx.Next()
	}
}

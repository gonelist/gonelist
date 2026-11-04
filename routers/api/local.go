package api

import (
	"mime"
	"net/http"

	"github.com/gin-gonic/gin"
	"gonelist/conf"
	"gonelist/pkg/app"
	"gonelist/pkg/e"
	"gonelist/service/local"
)

func localError(ctx *gin.Context) {
	app.Response(ctx, http.StatusForbidden, e.ERROR, "local file operation denied or failed")
	ctx.Abort()
}

func LocalDownload() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		p := ctx.Param("path")
		if !conf.UserSet.Local.Enable || !local.IsPath(p) {
			ctx.Next()
			return
		}
		f, err := local.Open(p)
		if err != nil {
			localError(ctx)
			return
		}
		defer f.Close()
		info, err := f.Stat()
		if err != nil || !info.Mode().IsRegular() {
			localError(ctx)
			return
		}
		ctx.Header("Content-Type", "application/octet-stream")
		ctx.Header("Content-Disposition", mime.FormatMediaType("attachment", map[string]string{"filename": info.Name()}))
		http.ServeContent(ctx.Writer, ctx.Request, info.Name(), info.ModTime(), f)
		ctx.Abort()
	}
}

func LocalMkdir() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		p := ctx.Query("path")
		if !conf.UserSet.Local.Enable || !local.IsPath(p) {
			ctx.Next()
			return
		}
		if err := local.Mkdir(p, ctx.Query("folder_name")); err != nil {
			localError(ctx)
			return
		}
		app.Response(ctx, http.StatusOK, e.SUCCESS, nil)
		ctx.Abort()
	}
}

func LocalUpload() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		p := ctx.Query("path")
		if !conf.UserSet.Local.Enable || !local.IsPath(p) {
			ctx.Next()
			return
		}
		file, err := ctx.FormFile("file")
		if err != nil {
			localError(ctx)
			return
		}
		src, err := file.Open()
		if err != nil {
			localError(ctx)
			return
		}
		defer src.Close()
		if err := local.Upload(p, file.Filename, src); err != nil {
			localError(ctx)
			return
		}
		app.Response(ctx, http.StatusOK, e.SUCCESS, nil)
		ctx.Abort()
	}
}

func LocalDelete() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		p := ctx.Query("path")
		if !conf.UserSet.Local.Enable || !local.IsPath(p) {
			ctx.Next()
			return
		}
		if err := local.Delete(p); err != nil {
			localError(ctx)
			return
		}
		app.Response(ctx, http.StatusOK, e.SUCCESS, nil)
		ctx.Abort()
	}
}

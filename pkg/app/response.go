package app

import (
	"github.com/gin-gonic/gin"
	"gonelist/pkg/e"
)

type Res struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data,omitempty"`
}

// Response setting gin.JSON
func Response(c *gin.Context, httpCode, errCode int, data interface{}) {
	res := Res{
		Code: errCode,
		Msg:  e.GetMsg(errCode),
		Data: data,
	}
	c.JSON(httpCode, res)
}

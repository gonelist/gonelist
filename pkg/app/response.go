package app

import (
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"GOIndex/pkg/e"
)

type Res struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

// Response setting gin.JSON
func Response(c *gin.Context, httpCode, errCode int, data interface{}) {
	res := Res{
		Code: errCode,
		Msg:  e.GetMsg(errCode),
		Data: data,
	}
	log.WithFields(log.Fields{
		"header":   c.Request.Header,
		"response": res,
	}).Info("返回内容")
	c.JSON(httpCode, res)
	return
}

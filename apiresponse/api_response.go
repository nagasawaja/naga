package apiresponse

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

type respBody struct {
	Code int64       `json:"code"`
	Data interface{} `json:"data"`
	Msg  string      `json:"msg"`
}

const (
	OK         = 0
	ParamsErr  = 1
	Timeout    = 2
	UnknownErr = 3
	WxidBusy   = 4
)

var errCodeMsg = map[int64]string{
	ParamsErr:  "错误参数",
	Timeout:    "查询超时",
	UnknownErr: "未知错误",
	WxidBusy:   "帐号业务繁忙",
}

func Ok(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, respBody{
		Code: OK,
		Data: data,
		Msg:  "",
	})
}

// arg params first param as msg, second param as data
func Fail(c *gin.Context, code int64, arg ...interface{}) {
	var msg = ""
	var data interface{}
	var argLen = len(arg)
	if argLen > 0 {
		// get arg first param as msg
		msg = fmt.Sprint(arg[0])
		if argLen > 1 {
			// get arg second param as data
			data = arg[1]
		}
	} else {
		// no arg param, so get default codeMsg
		v, ok := errCodeMsg[code]
		if ok == true {
			msg = v
		}
	}
	c.JSON(http.StatusOK, respBody{
		Code: code,
		Data: data,
		Msg:  msg,
	})
}

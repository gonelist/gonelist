package e

const (
	SUCCESS        = 200
	INVALID_PARAMS = 300
	REDIRECT_LOGIN = 400
	ERROR          = 500
	LOAD_NOT_READY = 600
	// 具体返回内容
	INVALID_STATE      = 10001
	ITEM_NOT_FOUND     = 10002
	ACCESS_TOKEN_ERROR = 10003
	MG_ERROR           = 10004
	CACHE_NOT_FIND     = 10005
	CONVERT_MD_ERROR   = 10006
	PASS_ERROR         = 10007
)

var MsgFlags = map[int]string{
	SUCCESS:        "ok",
	INVALID_PARAMS: "请求参数错误",
	REDIRECT_LOGIN: "需要重定向到登陆",
	ERROR:          "fail",
	LOAD_NOT_READY: "文件未加载完成",
	// 具体返回内容
	INVALID_STATE:      "state 字符串与设置的不一致，请检查设置",
	ITEM_NOT_FOUND:     "未找到对应项目",
	ACCESS_TOKEN_ERROR: "获取 AccessToken 错误",
	MG_ERROR:           "请求graph.microsoft.com错误",
	CACHE_NOT_FIND:     "缓存中未找到请求内容",
	CONVERT_MD_ERROR:   "转化markdown出现错误",
	PASS_ERROR:         "输入密码错误",
}

func GetMsg(code int) string {
	msg, ok := MsgFlags[code]
	if ok {
		return msg
	}

	return MsgFlags[ERROR]
}

func GetErrorCode(err error) int {
	for key, val := range MsgFlags {
		if err.Error() == val {
			return key
		}
	}
	return 500 //未找到具体错误
}

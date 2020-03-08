package e

const (
	SUCCESS            = 200
	INVALID_PARAMS     = 300
	ERROR              = 500
	// 具体返回内容
	INVALID_STATE  = 10001
	ITEM_NOT_FOUND = 10002
)

var MsgFlags = map[int]string{
	SUCCESS:            "ok",
	INVALID_PARAMS:     "请求参数错误",
	ERROR:              "fail",
	INVALID_STATE:      "state 字符串与设置的不一致，请检查设置",
	ITEM_NOT_FOUND:     "未找到对应项目",
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
	return 500; //未找到具体错误
}

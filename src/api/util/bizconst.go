package util

//空对象定义
type EmptyObject struct{}

var EMPTY_OBJECT = EmptyObject{}

//业务常量
const (
	SUCCESS     = 0
	SUCCESS_MSG = "success"

	PLATFORM_ANDROID = "android"
	PLATFORM_IOS     = "ios"

	PARAM_INVALID_ERRCODE = 10001
	PARAM_INVALID_ERRMSG  = "参数错误"

	REQUEST_INVALID_ERRCODE = 10000
	REQUEST_INVALID_ERRMSG  = "非法请求"

	GET_DB_ERRCODE = 10002
	GET_DB_ERRMSG  = "DB连接异常"

	REQUEST_TIMEOUT_ERRCODE = 10003
	REQUEST_TIMEOUT_ERRMSG  = "请求超时"
)

//配置常量
const (
	//请求时长
	MAX_REQUEST_INTERVAL = 10
)

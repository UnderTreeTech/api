package v2

import (
	"api/config"
	"api/util"
	"github.com/astaxie/beego"
	"sort"
	"strings"
)

//在beego.Controller上的进一步封装，主要进行一些验证
type BaseController struct {
	beego.Controller
	requestParams map[string]string
}

// Prepare implemented Prepare() method for baseController.
// It's used for request param valid and check
func (b *BaseController) Prepare() {
	b.requestParams = b.params()
	b.validParam()
	b.validRequest()
}

//获取请求的所有参数,以供业务controller使用
//规定参数不能为空
func (b *BaseController) params() (params map[string]string) {
	params = make(map[string]string)
	urlValues := b.Input()
	for paramKey, _ := range urlValues {
		paramVal := b.Ctx.Input.Query(paramKey)
		if "" == paramVal {
			b.RespResult(util.PARAM_INVALID_ERRCODE, util.PARAM_INVALID_ERRMSG, nil)
		}
		params[paramKey] = paramVal
	}

	return params
}

//校验必须参数
func (b *BaseController) validParam() {
	params := b.Input()
	requireParams := []string{"requestId", "signature", "timestamp", "appVersion", "platform"}

	for _, key := range requireParams {
		if _, ok := params[key]; !ok {
			b.RespResult(util.PARAM_INVALID_ERRCODE, util.PARAM_INVALID_ERRMSG, nil)
		}
	}
}

//校验请求签名是否一致
func (b *BaseController) validRequest() {
	platform := b.requestParams["platform"]
	signature := b.requestParams["signature"]
	timestamp := b.requestParams["timestamp"]
	appVersion := b.requestParams["appVersion"]
	requestId := b.requestParams["requestId"]

	appKey := b.getAppKey(platform, appVersion)

	dict := []string{appKey, timestamp, requestId}
	sort.Strings(dict)
	dictStr := strings.Join(dict, "")

	encryptStr := util.SHA1Encrypt(dictStr)

	if encryptStr != signature {
		b.RespResult(util.REQUEST_INVALID_ERRCODE, util.REQUEST_INVALID_ERRMSG, nil)
	}
}

//获取预埋appKey
//ios及android可以根据不同的版本设置不同的appkey
func (b *BaseController) getAppKey(platform string, appVersion string) (appKey string) {
	ok := true
	switch platform {
	case util.PLATFORM_ANDROID:
		if appKey, ok = config.AppkeyAndroid[appVersion]; !ok {
			b.RespResult(util.PARAM_INVALID_ERRCODE, util.PARAM_INVALID_ERRMSG, nil)
		}

	case util.PLATFORM_IOS:
		if appKey, ok = config.AppkeyIOS[appVersion]; !ok {
			b.RespResult(util.PARAM_INVALID_ERRCODE, util.PARAM_INVALID_ERRMSG, nil)
		}
	default:
		b.RespResult(util.PARAM_INVALID_ERRCODE, util.PARAM_INVALID_ERRMSG, nil)
	}

	return appKey
}

//请求回复
func (b *BaseController) RespResult(errCode int, errMsg string, data interface{}) {
	result := make(map[string]interface{})

	result["errCode"] = errCode
	result["errMsg"] = errMsg
	if nil == data {
		result["data"] = util.EMPTY_OBJECT
	} else {
		result["data"] = data
	}

	b.Data["json"] = result
	b.ServeJSON()
}

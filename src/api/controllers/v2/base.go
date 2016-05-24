package v2

import (
	"api/config"
	"api/dao"
	"api/util"
	"api/util/xss"
	"github.com/astaxie/beego"
	"sort"
	"strconv"
	"strings"
	"time"
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
	b.cleanXss()
	beego.Error(b.requestParams)
	b.validParam()
	b.validRequest()
	b.validToken()
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

//clean xss attack
func (b *BaseController) cleanXss() {
	for paramName, paramVal := range b.requestParams {
		cleanedVal := xss.GetXssHandler().Sanitize(paramVal)
		b.requestParams[paramName] = cleanedVal
	}
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

	//如果一个请求从发出到服务端接收到，超过10秒，丢弃请求
	//这一步是用来防请求重放攻击，同一请求不能无限次被抓包使用
	//注意，要做这一步必须做app与服务端的时间戳同步
	//另外，默认超时是10秒，有可能会误杀真正的慢请求
	now := time.Now().Unix()
	reqeustTime, err := strconv.ParseInt(timestamp, 10, 64)
	if err != nil {
		b.RespResult(util.PARAM_INVALID_ERRCODE, util.PARAM_INVALID_ERRMSG, nil)
	}

	if (now - reqeustTime) > util.MAX_REQUEST_INTERVAL {
		b.RespResult(util.REQUEST_TIMEOUT_ERRCODE, util.REQUEST_TIMEOUT_ERRMSG, nil)
	}

	appKey := b.getAppKey(platform, appVersion)

	dict := []string{appKey, timestamp, requestId}
	sort.Strings(dict)
	dictStr := strings.Join(dict, "")

	encryptStr := util.SHA1Encrypt(dictStr)

	if encryptStr != signature {
		b.RespResult(util.REQUEST_INVALID_ERRCODE, util.REQUEST_INVALID_ERRMSG, nil)
	}
}

//token是代表用户的唯一标志，前后端的往来只使用token，不能出现userid
//token是有时效性的，默认一个月后过期
//并不是所有接口都需要token参数，所以对请求做差异化处理
//对有token的接口，对其进行token有效性验证
func (b *BaseController) validToken() {
	params := b.Input()
	if token, ok := params["token"]; ok {
		//如果有token，不能为空
		if "" == token[0] {
			b.RespResult(util.PARAM_INVALID_ERRCODE, util.PARAM_INVALID_ERRMSG, nil)
		}

		//校验登录下的重放攻击
		b.validRequestId(token[0], b.requestParams["signature"])
		b.checkTokenInfo(token[0])
	}
}

//token与userid是绑定关系，token对外代表用户
//这一步暂时时间实现，新项目开始忙了，说下思路
//根据token获取token详情，验证token是否存在及过期
//如不存在，非法请求
//如没过期，拿到userid，从而查到userinfo，设置到缓存及requestParams["userinfo"]，方便后续调用
//如果过期，报登录态失败错误，提示app跳转登录
func (b *BaseController) checkTokenInfo(token string) {
	// cacheName := util.USER_INFO + token
	// userinfo, _ := dao.GetCache(cacheName)
	// if userinfo != nil {
	// 	userToken =
	// }
}

//对于有token的请求，证明用户在登录态下
//对请求可再做一次粗略的重放校验，以保证用户数据的安全
func (b *BaseController) validRequestId(token string, signature string) {
	cacheName := util.REQUEST_HISTORY + token
	lastSignature, _ := dao.GetCache(cacheName)
	if (lastSignature != nil) && (lastSignature == signature) {
		b.RespResult(util.REQUEST_INVALID_ERRCODE, util.REQUEST_INVALID_ERRMSG, nil)
	}
	//此处可设置cache的缓存时间，与token的生存周期一致
	dao.SetCache(cacheName, signature)
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

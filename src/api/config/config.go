package config

//一些基础配置
import (
	"github.com/astaxie/beego"
)

//app版本与接口版本的映射关系
var ApiVersionAndroid = map[string]string{
	"3.1.0" : "v1",
}

var ApiVersionIOS = map[string]string{
	"3.1.0" : "v1",
}

//app版本与appkey的映射关系
var AppkeyAndroid = make(map[string]string)

var AppkeyIOS = make(map[string]string)

//初始化版本与appkey的映射值
func init() {
	if runmode := beego.AppConfig.String("runmode"); runmode == "dev" {
		AppkeyAndroid["3.1.0"] = "testCcmsIam500QiangA"
		AppkeyIOS["3.1.0"] = "testCcmsIam500QiangA"
	}else{
		AppkeyAndroid["3.1.0"] = "xxxx"
		AppkeyIOS["3.1.0"] = "xxxx"
	}
}
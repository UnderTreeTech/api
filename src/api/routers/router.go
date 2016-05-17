// @APIVersion 1.0.0
// @Title beego Test API
// @Description beego has a very cool tools to autogenerate documents for your API
// @Contact astaxie@gmail.com
// @TermsOfServiceUrl http://beego.me/
// @License Apache 2.0
// @LicenseUrl http://www.apache.org/licenses/LICENSE-2.0.html
package routers

import (
	"api/controllers/v1"
	"api/controllers/v2"
	"github.com/astaxie/beego"
)

//配置路由规则
func init() {
	v1NS := beego.NewNamespace("/v1",
		beego.NSRouter("/person/list", &v1.PersonController{}, "get:List"),
	)

	v2NS := beego.NewNamespace("/v2", 
		beego.NSNamespace("/person", 
			beego.NSRouter("/get", &v2.PersonController{}, "get:Get"),
		),
		beego.NSRouter("/person/list", &v2.PersonController{}, "get:List"),
		
	)

	beego.AddNamespace(v1NS,v2NS)
}

package v2

import (
	// "api/dao"
	"api/models/v2"
	"api/util"
	// "github.com/astaxie/beego"
)

type PersonController struct {
	BaseController
}

func (p *PersonController) List() {
	persons := v2.List()
	p.Data["json"] = persons
	p.ServeJSON()
}

func (p *PersonController) Get() {
	requestParams := p.params()
	userId := requestParams["userId"]

	person := v2.Get(userId)

	//验证redis是否只初始化了一次
	// beego.Info(dao.GetRedisPool() == dao.GetRedisPool())

	p.RespResult(util.SUCCESS, util.SUCCESS_MSG, person)
}

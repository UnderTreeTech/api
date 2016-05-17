package v2

import (
	// "api/dao"
	"api/models/v2"
	"api/util"
	"github.com/astaxie/beego"
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
	beego.Info(p.requestParams)

	userId := p.requestParams["userId"]
	person := v2.Get(userId)

	p.RespResult(util.SUCCESS, util.SUCCESS_MSG, person)
}

package v1

import (
	"github.com/astaxie/beego"
	"api/models/v1"
)


type PersonController struct {
	beego.Controller
}

func (p *PersonController) List() {
	persons :=  v1.List()
	p.Data["json"] = persons
	p.ServeJSON()
}
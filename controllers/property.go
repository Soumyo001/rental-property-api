package controllers

import (
	"github.com/beego/beego/v2/core/logs"
	beego "github.com/beego/beego/v2/server/web"
)

type PropertyController struct {
	beego.Controller
}

// @Title Dummy GET property
// @Description Get the dummy JSON response data
// @Success 200 {object} string "ok"
// @Failure 400 {object} string "internal error"
// @Failure 500 {object} string "internal error"
// @router / [get]
func (p *PropertyController) GetDummy() {
	p.Data["json"] = map[string]string{"message": "Hello world"}
	if err := p.ServeJSON(); err != nil {
		logs.Error("Failed to send json response: %v", err)
	}
}

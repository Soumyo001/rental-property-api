package controllers

import (
	"rental-property-api/models"

	"github.com/beego/beego/v2/core/logs"
	beego "github.com/beego/beego/v2/server/web"
)

type BaseController struct {
	beego.Controller
}

func (b *BaseController) writeJSON(status int, payload interface{}) {
	b.Ctx.Output.SetStatus(status)
	b.Data["json"] = payload

	if err := b.ServeJSON(); err != nil {
		logs.Error("failed to send JSON response: %v", err)
	}
}

func (b *BaseController) writeError(status int, message string) {
	b.writeJSON(status, models.ErrorResponse{Error: message})
}

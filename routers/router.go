// @APIVersion 1.0.0
// @Title Rental Property API
// @Description Main route of rental property API
// @Contact admin@w3engineers.com
package routers

import (
	"rental-property-api/controllers"

	beego "github.com/beego/beego/v2/server/web"
)

func init() {
	ns := beego.NewNamespace("/v1",
		beego.NSNamespace("/properties",
			beego.NSInclude(&controllers.PropertyController{}),
		),
	)
	beego.AddNamespace(ns)
}

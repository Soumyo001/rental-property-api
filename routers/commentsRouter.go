package routers

import (
	beego "github.com/beego/beego/v2/server/web"
	"github.com/beego/beego/v2/server/web/context/param"
)

func init() {

    beego.GlobalControllerRouter["rental-property-api/controllers:PropertyController"] = append(beego.GlobalControllerRouter["rental-property-api/controllers:PropertyController"],
        beego.ControllerComments{
            Method: "GetDummy",
            Router: `/`,
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

}

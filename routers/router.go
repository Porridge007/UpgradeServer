// @APIVersion 1.0.0
// @Title beego Test API
// @Description beego has a very cool tools to autogenerate documents for your API
// @Contact astaxie@gmail.com
// @TermsOfServiceUrl http://beego.me/
// @License Apache 2.0
// @LicenseUrl http://www.apache.org/licenses/LICENSE-2.0.html
package routers

import (
	"UpgraderServer/controllers"

	"github.com/astaxie/beego"
)

func init() {
	beego.Router("/query_latest", &controllers.QueryLatestController{})
	beego.Router("/",&controllers.IndexController{})
	ns := beego.NewNamespace("/upgrade",
		beego.NSNamespace("/device",
			beego.NSInclude(
				&controllers.DeviceController{},
			),
		),

		beego.NSNamespace("/package",
			beego.NSInclude(
				&controllers.PackageController{},
			),
		),

		beego.NSNamespace("/user",
			beego.NSInclude(
				&controllers.UserController{},
			),
		),
	)
	beego.AddNamespace(ns)
}

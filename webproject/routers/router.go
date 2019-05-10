package routers

import (
	"webproject/controllers"
	"github.com/astaxie/beego"
)

func init() {
    beego.Router("/", &controllers.MainController{})
    beego.Router("/login", &controllers.UserController{}, "get:ShowLogin")
    beego.Router("/register", &controllers.UserController{}, "get:ShowReg;post:HandleReg")
	beego.Router("/sendVerCode",&controllers.UserController{},"post:HandleVerCode")
	beego.Router("/active",&controllers.UserController{},"get:ActiveUser")
}

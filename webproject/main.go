package main

import (
	_ "webproject/routers"
	_ "webproject/models"
	"github.com/astaxie/beego"
)

func main() {
	beego.Run()
}


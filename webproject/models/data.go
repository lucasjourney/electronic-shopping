package models

import (
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
)

type  User struct {
	Id int `orm:"pk;auto"`
	Name string `orm:"unique;size(40)"`
	PassWord string `orm:"size(40)"`
	Email string `orm:"unique;size(40)"`
	Active bool `orm:"default(false)"`
}

func init()  {
	//注册数据库
	orm.RegisterDataBase("default","mysql","root:123456@tcp(192.168.70.150:3306)/webProject")
	//注册表
	orm.RegisterModel(new(User))
	//跑起来
	orm.RunSyncdb("default",false,true)
}

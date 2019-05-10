package main

import (
	"fmt"
	"github.com/astaxie/beego/utils"
	"strconv"
)

type  User struct {
	Id int `orm:"pk;auto"`
	Name string `orm:"unique;size(40)"`
	PassWord string `orm:"size(40)"`
	Email string `orm:"unique;size(40)"`
	Active bool `orm:"default(false)"`
}

func main() {
	//配置用户
	user := User{}
	user.Id = 1
	user.Email = "312415754@qq.com"

	//发送邮件
	config := `{"username":"1510271838@qq.com","password":"ynojniemjvbnigch","host":"smtp.qq.com","port":587}`
	temail := utils.NewEMail(config)
	temail.To = []string{user.Email}
	temail.From = "1510271838@qq.com"
	temail.Subject = "真吃货用户激活"

	temail.HTML = "复制该连接到浏览器中激活：127.0.0.1:8080/active?id=" + strconv.Itoa(user.Id)

	err := temail.Send()
	if err != nil {
		//this.Data["errmsg"] = "发送激活邮件失败，请重新注册！"
		//this.TplName = "register.html"
		fmt.Println("发送邮件失败")
		return
	}
	fmt.Println("发送邮件成功")
	//this.Ctx.WriteString("注册成功，请前往邮箱激活!")

}

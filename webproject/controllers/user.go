package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"github.com/astaxie/beego/utils"
	"math/rand"
	"regexp"
	"strconv"
	"time"
	"webproject/models"
)

type UserController struct {
	beego.Controller
}
//展示登录页面
func (this *UserController) ShowLogin()  {
	this.TplName = "login.html"
}
//展示注册页面
func (this *UserController) ShowReg() {
	this.TplName = "register.html"
}

//处理验证码
func (this *UserController) HandleVerCode() {
	//接收ajax的值 phone
	phone := this.GetString("phone")
	//定义返回值字段
	resp := make(map[string]interface{})

	//传值
	defer RespFunc(this, resp)
	//判断phone是否合法
	if phone == "" {
		resp["errno"] = 1
		resp["errmsg"] = "电话号码不能为空,请重新输入"
		return
	}

	//初始化客户端  需要accessKey  需要开通申请
	client, err := sdk.NewClientWithAccessKey("default", "LTAIvFXQmq69AXbm", "z4bbiK2XAx8HunIBBXWPY3JXjBI3A0")
	if err != nil {
		resp["errno"] = 2
		resp["errmsg"] = "阿里云客户端初始化失败"
		fmt.Println("阿里云客户端初始化失败")
		return
	}
	//获取6位数随机码
	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))
	vcode := fmt.Sprintf("%06d", rnd.Int31n(1000000))

	//初始化请求对象
	request := requests.NewCommonRequest()
	request.Method = "POST"//设置请求方法
	request.Scheme = "https" // https | http   //设置请求协议
	request.Domain = "dysmsapi.aliyuncs.com"  //域名
	request.Version = "2017-05-25"			//版本号
	request.ApiName = "SendSms"				//api名称
	request.QueryParams["PhoneNumbers"] = phone  //需要发送的电话号码
	request.QueryParams["SignName"] = "真吃货之家"    //签名名称   需要申请
	request.QueryParams["TemplateCode"] = "SMS_165115755"   //模板号   需要申请
	request.QueryParams["TemplateParam"] = `{"code":`+vcode+`}`   //发送短信验证码

	response, err := client.ProcessCommonRequest(request)  //发送短信
	if err != nil {
		resp["errno"] = 3
		resp["errmsg"] = "发送短信失败"
		return
	}
	//fmt.Println(string(response.GetHttpContentBytes()))
	//fmt.Println("结束")
	var msg models.MSG

	json.Unmarshal(response.GetHttpContentBytes(),&msg)  //解析发送结果
	if msg.Message != "OK"{
		resp["errno"] = 4
		resp["errmsg"] = "短信发送失败"
		return
	}
	fmt.Println(msg)
	resp["errno"] = 5
	resp["errmsg"] = "短信发送成功"
	resp["verCode"] = vcode
}

//向ajax返回json值
func RespFunc(this *UserController, resp map[string]interface{})  {
	//ajax前后端通信
	this.Data["json"] = resp
	//fmt.Println("resp的值是：", resp)
	this.ServeJSON()
	return
}

//处理注册业务
func (this *UserController)HandleReg() {

	//1.获取数据
	userName := this.GetString("phone")
	pwd := this.GetString("password")
	cpwd := this.GetString("repassword")
	email := this.GetString("email")

	//对接收到的数据进行非空校验
	if userName == "" || pwd == "" || cpwd == "" || email == ""{
		this.Data["errmsg"] = "输入信息不完整,请重新输入！"
		this.TplName = "register.html"
		return
	}
	//判断“密码”和”确认密码“是否一致。
	if pwd != cpwd{
		this.Data["errmsg"] = "两次输入密码不一致，请重新输入！"
		this.TplName = "register.html"
		return
	}
	//对邮箱格式的校验：注意这里使用的是regexp下的Compile方法。
	reg,_ := regexp.Compile("^[A-Za-z0-9\u4e00-\u9fa5]+@[a-zA-Z0-9_-]+(\\.[a-zA-Z0-9_-]+)+$")
	res := reg.FindString(email)
	if res == ""{
		this.Data["errmsg"] = "邮箱格式不正确"
		this.TplName = "register.html"
		return
	}

	o := orm.NewOrm()
	user := models.User{}
	user.Name = userName
	user.PassWord = pwd
	user.Email = email
	//检验重复的手机号
	err := o.Read(&user,"Name")
	if err != orm.ErrNoRows{
		this.Data["errmsg"] = "用户以存在，请重新注册！"
		this.TplName = "register.html"
		return
	}
	//插入数据
	_,err = o.Insert(&user)
	if err != nil {
		this.Data["errmsg"] = "插入失败，请重新注册！"
		this.TplName = "register.html"
		return
	}

	//发送激活邮件部分
	//发送邮件
	config := `{"username":"1510271838@qq.com","password":"ynojniemjvbnigch","host":"smtp.qq.com","port":587}`
	temail := utils.NewEMail(config)
	temail.To = []string{user.Email}
	temail.From = "1510271838@qq.com"
	temail.Subject = "真吃货用户激活"

	temail.HTML = "复制该连接到浏览器中激活：127.0.0.1:8080/active?id=" + strconv.Itoa(user.Id)

	err = temail.Send()
	if err != nil {
		//this.Data["errmsg"] = "发送激活邮件失败，请重新注册！"
		//this.TplName = "register.html"
		fmt.Println("发送邮件失败")
		return
	}
	beego.Info("发送邮件成功")
	this.Ctx.WriteString("注册成功，请前往邮箱激活!")
	//this.TplName = "login.html"
}

//处理激活邮箱业务
func (this *UserController)ActiveUser() {
	//获取数据
	id, err := this.GetInt("id")
	if err != nil {
		this.Data["errmsg"] = "激活路径不正确"
		this.TplName = "login.html"
		return
	}
	//查询表
	o := orm.NewOrm()
	var user models.User
	user.Id = id
	o.Read(&user)
	if err != nil {
		this.Data["errmsg"] = "激活路径不正确"
		this.TplName = "login.html"
		return
	}
	//更新表
	user.Active = true
	_, err = o.Update(&user)
	if err != nil {
		this.Data["errmsg"] = "激活失败"
		this.TplName = "login.html"
		return
	}
	this.Redirect("/login", 302)
}
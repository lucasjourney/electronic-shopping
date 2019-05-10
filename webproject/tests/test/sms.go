package main

import (
	"encoding/json"
	"fmt"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"math/rand"
	"time"
)

func main()  {
	//定义返回值字段
	resp := make(map[string]interface{})
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
	request.QueryParams["PhoneNumbers"] = "13810549573"  //需要发送的电话号码
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
	var msg MSG

	json.Unmarshal(response.GetHttpContentBytes(),&msg)  //解析发送结果
	if msg.Message != "OK"{
		resp["errno"] = 4
		resp["errmsg"] = "短信发送失败"

		return
	}
	fmt.Println(msg)
}

type MSG struct {
	Message string `json:"Message"`
	RequestId string `json:"RequestId"`
	BizId string `json:"BizId"`
	Code string `json:"Code"`
}

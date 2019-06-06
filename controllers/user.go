package controllers

import (
	"BeeProject/models"
	"encoding/base64"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
)

type UserController struct {
	beego.Controller
}

//注册函数
func (self  *UserController)ShowRegister(){
	self.TplName = "register.html"
}

//提交注册数据

func (self *UserController)PostData()  {
	//获取数据
	name := self.GetString("userName")
	fmt.Println("kkkkkk",name)
	password := self.GetString("password")
	//校验数据
	logs.Info(name,password)
	if name == "" || password == ""{
		self.Data["error"] = "数据不能为空"
		logs.Info("注册数据不能为空")
		self.TplName = "register.html"
		return
	}
	//新建orm
	o := orm.NewOrm()
	//操作数据
	var user models.User
	user.UserName = name
	user.PassWord = password
	o.Insert(&user)
	//返回页面
	self.Redirect("/login",302)
}


//登陆页面
func (self *UserController)ShowLogin(){
	cookie := self.Ctx.GetCookie("userName")
	logs.Info("cookies",cookie)
	if cookie == ""{
		self.Data["userName"] = ""
		self.Data["checked"] = ""
	}else {
		self.Data["userName"] = cookie
		self.Data["checked"] = "checked"
	}
	self.TplName = "login.html"
}

func (self *UserController)PostLogin(){
	//获取数据
	name := self.GetString("userName")
	password := self.GetString("password")
	if name == "" || password == ""{
		self.Data["error"] = "表格不能为空"
		logs.Info("注册数据不能为空")
		self.TplName = "login.html"
		return
	}
	o := orm.NewOrm()
	//操作数据
	var user models.User
	user.UserName = name


	err := o.Read(&user,"UserName")
	if err != nil{
		logs.Error("用户名错误,请重新输入")
		self.TplName = "login.html"
		return
	}
	if user.PassWord !=password{
		logs.Error("密码错误")
		self.TplName = "login.html"
		return
	}

	//获取记住用户状态

	data := self.GetString("remember")
	if data == "on"{
		temp := base64.StdEncoding.EncodeToString([]byte(name))
		//设置cookie
		self.Ctx.SetCookie("userName",temp,100)
	}else{
		self.Ctx.SetCookie("userName",name,-1)
	}

	self.SetSession("userName",name)
	self.Redirect("/article/showArticleList",302)
}

func (self *UserController)ShowLogout(){
	//删除session
	self.DelSession("userName")
	self.Redirect("/login",302)
}
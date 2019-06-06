package controllers

import (
	"BeeProject/models"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"


)

type MainController struct {
	beego.Controller
}

func (c *MainController) Get() {
	c.Data["Name"] = "1111"
	c.Data["Password"] = "522222"
	c.TplName = "form.html"
}


func (c *MainController) Post() {
	c.Data["Name"] = "admin"
	c.Data["Password"] = "admin"
	c.TplName = "form.html"
}


func (c *MainController) ShowGet() {
	//orm操作

	o := orm.NewOrm()
	//插入操作
	var user models.User
	user.UserName = "hihaha"
	user.PassWord = "wowow"
	count,err := o.Insert(&user)
	if err != nil{
		logs.Error("insert data error,err:%v",err)

	}
	logs.Info("count is:%v",count)
	//c.Data["Name"] = "admin"
	//c.Data["Password"] = "admin"
	c.TplName = "form.html"

	//查询对象
	var userData models.User
	userData.Id = 1
	err = o.Read(&userData,"id")
	if err != nil{
		logs.Error("read data error from table,err:%v",err)
		return
	}
	logs.Info(userData)

	//更新操作
	var data models.User
	data.Id = 1
	err = o.Read(&data)
	if err != nil{
		logs.Error("要更新的数据不存在")

	}
	data.UserName = "wudi"
	count,err = o.Update(&user)
	if err != nil{
		logs.Error("更新失败")
	}
	logs.Info(count)

	//删除操作
	var deleteData models.User
	deleteData.Id = 2
	num,err := o.Delete(&deleteData)
	if err != nil{
		logs.Error("删除失败")
	}
	logs.Info(num)

}

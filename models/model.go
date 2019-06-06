package models

import (
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
	"time"
)
/*
数据库表的设计
 */
//func init(){
//	//操作数据库
//
//	conn,err := sql.Open("mysql","root:root@tcp(127.0.0.1:3306)/wudi?charset=utf8")
//	if err != nil{
//		logs.Error("connect mysql error,err:%v",err)
//		return
//	}
//	defer conn.Close()
//	logs.Debug("init sql now")
//	//res,err := conn.Exec("create table  user (username VARCHAR (255),pwd varchar (255));")
//	//if err != nil{
//	//	logs.Error("sql table create error,err:%v\n,result is :%v",err,res)
//	//	return
//	//}
//
//	logs.Debug("table create now")
//	//conn.Exec("INSERT user(username,pwd)value(?,?)","wudi2","wudi2")
//	res,err := conn.Query("SELECT username from user")
//	if err != nil{
//		logs.Error("select data from table error,err:%v",err)
//		return
//	}
//	var username string
//	for res.Next(){
//		res.Scan(&username)
//		logs.Info(username)
//
//	}
//
//}

type User struct {
	Id int
	UserName string
	PassWord string

	Article []*Article `orm:"reverse(many)"`

}

type Article struct{
	Id int `orm:"pk;auto"`
	ArtiName string `orm:"size(50)"`
	Atime time.Time `orm:"auto_now"`
	Acount int `orm:"default(0);null"`
	Acontent string `orm:"size(500)"`
	Aimg string  `orm:"size(100)"`

	ArticleType *ArticleType `orm:"rel(fk)"`

	Users []*User `orm:"rel(m2m)"`

}

type ArticleType struct{
	Id int
	TypeName string `orm:"size(20)"`
	Articles []*Article `orm:"reverse(many)"`
}
func init(){
	// orm操作,双下滑线有特殊含义
	//获取连接对象
	orm.RegisterDataBase("default","mysql","root:root@tcp(127.0.0.1:3306)/test?charset=utf8")
	//创建表
	orm.RegisterModel(new(User),new(Article),new(ArticleType))
	//生成表
	orm.RunSyncdb("default",false,true)
	//操作表

}
package routers

import (
	"BeeProject/controllers"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"
)

func init() {
    //beego.Router("/", &controllers.MainController{},"get:ShowGet;post:Post")
    //过滤器函数
	beego.InsertFilter("/article/*",beego.BeforeRouter,filter)
	beego.Router("/register", &controllers.UserController{},"get:ShowRegister;post:PostData")
	beego.Router("/login", &controllers.UserController{},"get:ShowLogin;post:PostLogin")
	//实现文章列表页访问
	beego.Router("/article/showArticleList",&controllers.ArticleController{},"get:ShowArticleList")
	//添加文章
	beego.Router("/article/addArticle",&controllers.ArticleController{},"get:ShowAddArticle;post:HandleArticle")
	//显示文章详情
	beego.Router("/article/showArticleDetail",&controllers.ArticleController{},"get:ShowArticleDetail")
	//编辑文章页面
	beego.Router("/article/updateArticle",&controllers.ArticleController{},"get:ShowUpdateArticle;post:HandleUpdateArticle")
	//删除文章
	beego.Router("/article/deleteArticle",&controllers.ArticleController{},"get:ShowDeleteArticle")
	//添加分类
	beego.Router("/article/addArticleType",&controllers.ArticleController{},"get:ShowAddType;post:HandleAddType")
   //删除分类
	beego.Router("/article/deleteArticleType",&controllers.ArticleController{},"get:ShowDeleteType")
    //退出
	beego.Router("/article/logout",&controllers.UserController{},"get:ShowLogout")
	//一个请求指定一个方法
	//beego.Router("/login",&controllers.LoginController{},"get:Login;post:ShowPost")
	////多个请求指定一个方法
	//beego.Router("/index",&controllers.IndexController{},"get,post:GetPost")
    ////所有请求指定一个方法
	//beego.Router("/handle",&controllers.HandleController{},"*:HandleFunc")
}


var filter = func(ctx *context.Context){
	//获取session
	userName :=ctx.Input.Session("userName")
	if userName == nil{
		ctx.Redirect(302,"/login")
		return
	}
}
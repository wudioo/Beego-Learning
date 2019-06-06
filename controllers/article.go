package controllers

import (
	"BeeProject/models"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
	"math"
	"path"
	"time"
)

type ArticleController struct {
	beego.Controller
}

func (self *ArticleController)ShowArticleList(){

	//session判断
	userName := self.GetSession("userName")
	if userName == nil{
		self.Redirect("/login",302)
		return
	}
	//高级查询
	o := orm.NewOrm()
	qs := o.QueryTable("article")
	var articles []models.Article

	//count,err :=qs.All(&articles)
	//if err != nil{
	//	logs.Error("查询数据错误")
	//
	//}
	//查询总记录数
	//获取下拉框数据类型
	typeName := self.GetString("select")

	var count int64
	pageSize := 3
	//天花板函数,向上取整,获取总页数


	//获取页码
	pageIndex,err := self.GetInt("pageIndex")
	if err != nil{
		logs.Error("获取数据类型错误")
		pageIndex = 1
	}
	start := (pageIndex - 1) * pageSize
	if typeName == ""{
		count,err = qs.Count()
	}else {
		count,err = qs.Limit(pageSize,start).RelatedSel("ArticleType").Filter("ArticleType__TypeName",typeName).All(&articles)
	}
	pageNum := math.Ceil(float64(count) / float64(pageSize))
	if err != nil{
		logs.Error("查询数据错误")
	}
	//第一个参数获取几条,第二个参数表示从第几条数据开始
	//qs.Limit(pageSize,start).RelatedSel("ArticleType").All(&articles)
	var types []models.ArticleType
	o.QueryTable("ArticleType").All(&types)
	self.Data["types"] = types
	logs.Info("类型",types)


	logs.Info("下拉框获取",typeName)
	if typeName == ""{
		qs.Limit(pageSize,start).RelatedSel("ArticleType").All(&articles)
	}else {
		//RelatedSel关联的是字段名
		qs.Limit(pageSize,start).RelatedSel("ArticleType").Filter("ArticleType__TypeName",typeName).All(&articles)
	}




	self.Data["typeName"] = typeName
	self.Data["count"] = count
	self.Data["articles"] = articles
	self.Data["pageNum"] = int(pageNum)
	self.Data["pageIndex"] = pageIndex
	LayoutName := self.GetSession("userName")
	self.Data["LayoutName"] = LayoutName.(string)
	//指定视图布局
	self.Layout = "layout.html"
	self.TplName = "index.html"
}

//添加文章
func (self *ArticleController)ShowAddArticle(){
	o := orm.NewOrm()
	var types []models.ArticleType
	o.QueryTable("ArticleType").All(&types)
	self.Data["Types"] = types
	LayoutName := self.GetSession("userName")
	self.Data["LayoutName"] = LayoutName.(string)
	self.Layout = "layout.html"
	self.TplName = "add.html"
}

//获取添加文章内容
func (self *ArticleController)HandleArticle(){
	//1.获取数据
	articleName := self.GetString("articleName")
	content := self.GetString("content")

	//2.校验数据
	if articleName == "" || content == ""{
		self.Data["errmsg"] = "添加数据不能为空"
		self.TplName = "add.html"
		return
	}
	logs.Info(articleName,content)
	filePath :=uploadFile(&self.Controller,"uploadname","add.html")
	//3.操作数据
	//插入操作
	o := orm.NewOrm()
	var article models.Article
	article.ArtiName = articleName
	article.Acontent = content
	//不能有.
	article.Aimg = filePath
	//文章类型
	typeName := self.GetString("select")

	//根据名称查询类型
	var articleType models.ArticleType
	articleType.TypeName = typeName
	o.Read(&articleType,"TypeName")

	article.ArticleType = &articleType

	o.Insert(&article)
	//4.返回页面
	self.Redirect("/showArticleList",302)


}

func (self  *ArticleController)ShowArticleDetail(){
	//1.获取数据
	articleId,err := self.GetInt("articleId")

	if err !=nil{
		logs.Error("传递数据错误")
	}

	o := orm.NewOrm()
	var article models.Article
	article.Id = articleId

	//o.Read(&article)
	o.QueryTable("Article").RelatedSel("ArticleType").Filter("Id",articleId).One(&article)

	//修改阅读量
	article.Acount += 1
	o.Update(&article)


	//多对多插入
	m2m:= o.QueryM2M(&article,"Users")
	userName := self.GetSession("userName")
	logs.Info("session Name :",userName)
	if userName == nil{
		self.Redirect("/login",302)
		return

	}
	var users models.User
	var usersAll []models.User
	users.UserName = userName.(string)
	o.Read(&users,"UserName")
	logs.Info(users)

	m2m.Add(users)
	//o.LoadRelated(&article,"Users")
	o.QueryTable("User").Filter("Article__Article__Id",articleId).Distinct().All(&usersAll)

	self.Data["article"] = article
	self.Data["users"] = usersAll
	LayoutName := self.GetSession("userName")
	self.Data["LayoutName"] = LayoutName.(string)
	self.Layout = "layout.html"
	self.TplName = "content.html"

}

//编辑文章页面

func (self *ArticleController)ShowUpdateArticle(){
	articleId,err := self.GetInt("articleId")

	if err !=nil{
		logs.Error("传递数据错误")
	}

	o := orm.NewOrm()
	var article models.Article
	article.Id = articleId

	o.Read(&article)
	self.Data["article"] = article
	LayoutName := self.GetSession("userName")
	self.Data["LayoutName"] = LayoutName.(string)
	self.Layout = "layout.html"

	self.TplName = "update.html"
}
func uploadFile(self *beego.Controller,fieldText string,tplName string)string{
	//处理文件上传
	file,head,err := self.GetFile(fieldText)
	if head.Filename == ""{
		return "Noimg"

	}
	if err != nil{
		self.Data["errmsg"] = "文件上传失败"
		self.TplName = tplName
		return ""
	}
	defer file.Close()

	//1.文件大小
	if head.Size > 5000000{
		self.Data["errmsg"] = "文件大小不得超过5M"
		self.TplName = tplName
		return ""
	}
	//文件格式,后缀
	ext := path.Ext(head.Filename)
	if ext != ".jpg" && ext != ".png" && ext != ".jpeg"{
		self.Data["errmsg"] = "文件格式不支持"
		self.TplName = tplName
		return ""
	}
	fileName := time.Now().Format("2006-01-02 15-04-05") + ext
	logs.Info(fileName,head.Size)
	//必须要有.
	self.SaveToFile(fieldText,"./static/img/"+ fileName)
	return "/static/img/"+ fileName
}

//编辑界面
func (self *ArticleController)HandleUpdateArticle(){
	articleId,err := self.GetInt("articleId")
	if err != nil{
		logs.Error("更新页面数据类型错误")
		return
	}
	artiName := self.GetString("articleName")
	content := self.GetString("content")
	filePath :=uploadFile(&self.Controller,"uploadname","update.html")


	if err != nil || artiName == "" || content == "" ||filePath == ""{
		logs.Error("请求错误")
		return
	}


	o := orm.NewOrm()
	var article models.Article
	article.Id = articleId
	err = o.Read(&article)
	if err !=nil{
		logs.Error("更新数据不存在")
		return
	}
	article.ArtiName = artiName
	article.Acontent = content
	if filePath != "Noimg"{
		article.Aimg = filePath
	}

	o.Update(&article)

	//返回页面
	self.Redirect("/showArticleList",302)


}

//删除文章处理
func (self *ArticleController)ShowDeleteArticle(){
	//1.获取数据
	articleId,err := self.GetInt("articleId")
	if err != nil{
		logs.Error("更新页面数据类型错误")
		return
	}
	o := orm.NewOrm()
	var article models.Article
	article.Id = articleId
	o.Delete(&article)

	//返回页面
	self.Redirect("/showArticleList",302)

}

//添加分类
func (self *ArticleController)ShowAddType(){
	o := orm.NewOrm()
	var articleType []models.ArticleType
	o.QueryTable("ArticleType").All(&articleType)
	self.Data["ArticleType"] = articleType
	LayoutName := self.GetSession("userName")
	self.Data["LayoutName"] = LayoutName.(string)
	self.Layout = "layout.html"
	self.TplName = "addType.html"

}

//删除类型
func (self *ArticleController)ShowDeleteType(){
	//获取数据
	Id,err := self.GetInt("Id")
	if err != nil{
		logs.Info("获取类型Id错误",err)
		return
	}
	o := orm.NewOrm()
	var articleType models.ArticleType
	articleType.Id = Id
	o.Delete(&articleType)
	self.Redirect("/article/addArticleType",302)

}

func (self *ArticleController)HandleAddType(){
	//1.获取数据
	typeName := self.GetString("typeName")

	//校验数据
	if typeName == ""{
		logs.Error("类型不能为空")
		self.Data["errmsg"] = "类型不能为空"
		self.TplName = "addType.html"
		return
	}
	//处理数据
	o := orm.NewOrm()
	var articleType models.ArticleType
	articleType.TypeName = typeName
	o.Insert(&articleType)

	self.Redirect("/article/addArticleType",302)

}
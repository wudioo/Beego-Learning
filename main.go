package main

import (
	_ "BeeProject/routers"
	"github.com/astaxie/beego"
	_"BeeProject/models"
)

func main() {
	beego.AddFuncMap("prepage",ShowPrePage)
	beego.AddFuncMap("nextpage",ShowNextPage)
	beego.Run()
}


//后台定义函数
func ShowPrePage(pageIndex int)int{
	if pageIndex == 1{
		return  pageIndex
	}
	return pageIndex -1

}


func ShowNextPage(pageIndex int,pageNum int)int{
	if pageIndex == pageNum{
		return pageIndex
	}
	return pageIndex + 1

}

/*
处理视图中简单业务逻辑
1.创建后台函数
2.在视图中定义函数名
3.在beego.run之前关联函数

 */

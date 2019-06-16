package main

import (
	//_ "Article-Manage/models"
	_ "Article-Manage/routers"
	"github.com/astaxie/beego"
)

func main() {
	beego.AddFuncMap("ShowPrePage", ShowPrePage)
	beego.AddFuncMap("ShowNextPage", ShowNextPage)
	beego.Run()
}

//视图函数，获取上一页页码

/*
1.在视图中定义视图函数函数名    | funcName

2.一般在main.go里面实现视图函数

3.在main函数里面把实现的函数和视图函关联起来   beego.AddFuncMap()
*/
func ShowPrePage(pageindex int) (preIndex int) {
	if pageindex == 1 {
		return 1
	} else {
		preIndex = pageindex - 1
	}
	return
}

func ShowNextPage(pageindex int) (nextIndex int) {
	nextIndex = pageindex + 1
	return
}

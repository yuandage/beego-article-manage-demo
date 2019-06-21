package routers

import (
	"Article-Manage/controllers"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"
)

func init() {

	var Filter = func(ctx *context.Context) {
		//判断是否登录
		userName := ctx.Input.Session("userName")
		if userName == nil {
			ctx.Redirect(302, "/login")
			return
		}
	}

	beego.InsertFilter("/index", beego.BeforeRouter, Filter)
	beego.InsertFilter("/addArticle", beego.BeforeRouter, Filter)
	beego.InsertFilter("/content", beego.BeforeRouter, Filter)
	beego.InsertFilter("/update", beego.BeforeRouter, Filter)
	beego.InsertFilter("/delete", beego.BeforeRouter, Filter)
	beego.InsertFilter("/addType", beego.BeforeRouter, Filter)

	beego.Router("/", &controllers.MainController{})
	beego.Router("/abc", &controllers.MainController{})

	beego.Router("/register", &controllers.MainController{}, "get:ShowRegister")
	//注意：当实现了自定义的请求方法，请求将不会访问默认方法
	beego.Router("/login", &controllers.MainController{}, "get:ShowLogin;post:HandleLogin")
	beego.Router("/logout", &controllers.MainController{}, "get:Logout")
	beego.Router("/index", &controllers.MainController{}, "get:ShowIndex")

	beego.Router("/addArticle", &controllers.MainController{}, "get:ShowAdd;post:HandleAdd")

	beego.Router("/content", &controllers.MainController{}, "get:ShowContent")

	beego.Router("/update", &controllers.MainController{}, "get:ShowUpdate;post:HandleUpdate")

	beego.Router("/delete", &controllers.MainController{}, "get:HandleDelete")

	beego.Router("/addType", &controllers.MainController{}, "get:ShowAddType;post:HandleAddType")
	beego.Router("/deleteType", &controllers.MainController{}, "get:HandleDeleteType")
	beego.Router("/updateType", &controllers.MainController{}, "post:HandleUpdateType")
}

package controllers

import (
	"Article-Manage/models"
	"github.com/astaxie/beego"
	"math"

	//"github.com/astaxie/beego/orm"
	//"class/models"
	"github.com/astaxie/beego/orm"
	"path"
	"time"
)

type MainController struct {
	beego.Controller
}

func (c *MainController) Get() {
	o := orm.NewOrm()
	var articles []models.Article //结构体数组
	_, err := o.QueryTable("Article").All(&articles)
	if err != nil {
		beego.Info("查询所有文章信息出错")
		return
	}
	c.Data["articles"] = articles
	c.TplName = "home.html"
}

func (c *MainController) ShowRegister() {
	c.TplName = "register.html"
}

//注册
func (c *MainController) Post() {

	//1.拿到数据
	userName := c.GetString("userName")
	pwd := c.GetString("pwd")
	//2.对数据进行校验
	if userName == "" || pwd == "" {
		beego.Info("数据不能为空")
		c.Redirect("/register", 302)
		return
	}
	//3.插入数据库
	o := orm.NewOrm()

	user := models.User{}
	user.Name = userName
	user.Pwd = pwd
	_, err := o.Insert(&user)
	if err != nil {
		beego.Info("插入数据失败")
		c.Redirect("/register", 302)
		return
	}
	//4.返回登陆界面
	c.Redirect("/login", 302)
}

func (c *MainController) ShowLogin() {
	userName := c.Ctx.GetCookie("userName")
	if userName != "" {
		c.Data["userName"] = userName
		c.Data["check"] = "checked"
	}
	c.TplName = "login.html"
}

//登陆业务处理
func (c *MainController) HandleLogin() {
	//c.Ctx.WriteString("这是登陆的POST请求")
	//1.拿到数据
	userName := c.GetString("userName")
	pwd := c.GetString("pwd")
	//2.判断数据是否合法
	if userName == "" || pwd == "" {
		beego.Info("输入数据不合法")
		c.TplName = "login.html"
		return
	}
	//3.查询账号密码是否正确
	o := orm.NewOrm()
	user := models.User{}

	user.Name = userName
	err := o.Read(&user, "Name")
	if err != nil {
		beego.Info("查询失败")
		c.TplName = "login.html"
		return
	}
	//记住用户名
	remember := c.GetString("remember")
	if remember == "on" {
		c.Ctx.SetCookie("userName", userName, time.Hour)
	} else {
		c.Ctx.SetCookie("userName", "", -1)
	}
	//登陆状态保存
	c.SetSession("userName", userName)
	//4.跳转
	c.Redirect("/index", 302)
}

//注销
func (c *MainController) Logout() {
	//1.删除登陆状态
	c.DelSession("userName")
	//2.跳转登陆页面
	c.Redirect("/", 302)
}

//显示列表页面内容
func (c *MainController) ShowIndex() {

	userName := c.GetSession("userName")

	o := orm.NewOrm()
	id, _ := c.GetInt("select")

	//beego.Info("id=", id)
	var articles []models.Article //结构体数组
	_, err := o.QueryTable("Article").All(&articles)
	if err != nil {
		beego.Info("查询所有文章信息出错")
		return
	}

	pageIndex, err := c.GetInt("pageIndex")
	if err != nil {
		pageIndex = 1
	}

	//获取类型数据
	var artiTypes []models.ArticleType
	artiTypesCount, err := o.QueryTable("ArticleType").All(&artiTypes)
	if err != nil {
		beego.Info("获取类型错误")
		return
	}
	if id > int(artiTypesCount) {
		id = 1
	}
	var count int64
	if id == 0 { //显示全部的数据
		count, err = o.QueryTable("Article").Count() //总记录条数
		//beego.Info("count=", count)
	} else {
		count, err = o.QueryTable("Article").RelatedSel("ArticleType").Filter("ArticleType__Id", id).Count() //总记录条数
	}
	pageSize := 2                                              //一页的记录条数
	pageCount := math.Ceil(float64(count) / float64(pageSize)) //Ceil向上取整,Floor()向下取整

	if pageIndex <= 0 {
		pageIndex = 1
	}
	if pageIndex > int(pageCount) {
		pageIndex = int(pageCount)
	}
	start := pageSize * (pageIndex - 1)
	if id == 0 {
		o.QueryTable("Article").Limit(pageSize, start).RelatedSel("ArticleType").All(&articles)
	} else {
		o.QueryTable("Article").Limit(pageSize, start).RelatedSel("ArticleType").Filter("ArticleType__Id", id).All(&articles)
	}

	if err != nil {
		beego.Info("查询错误")
		return
	}

	c.Data["pageCount"] = pageCount
	c.Data["count"] = count
	c.Data["pageIndex"] = pageIndex
	c.Data["articles"] = articles
	c.Data["articleType"] = artiTypes
	c.Data["typeid"] = id //文章类型ID
	c.Data["userName"] = userName

	c.Layout = "layout.html"
	c.TplName = "index.html"
}

//显示添加文章界面
func (c *MainController) ShowAdd() {
	o := orm.NewOrm()
	var artiTypes []models.ArticleType
	_, err := o.QueryTable("ArticleType").All(&artiTypes)
	if err != nil {
		beego.Info("获取类型错误")
		return
	}

	c.Data["articleType"] = artiTypes
	c.Layout = "layout.html"
	c.TplName = "add.html"
}

//处理添加文章界面数据
func (c *MainController) HandleAdd() {
	//1.拿到数据
	artiName := c.GetString("articleName")
	artiContent := c.GetString("content")
	id, err := c.GetInt("select")
	if err != nil {
		beego.Info("获取类型错误")
		return
	}
	f, h, err := c.GetFile("uploadname")
	defer f.Close()

	//1.要限定格式
	fileext := path.Ext(h.Filename) //文件后缀
	if fileext != ".jpg" && fileext != ".png" {
		beego.Info("上传文件格式错误")
		return
	}
	//2.限制大小
	if h.Size > 50000000 {
		beego.Info("上传文件过大")
		return
	}

	//3.需要对文件重命名，防止文件名重复
	filename := time.Now().Format("2006-01-02 15-04-05") + fileext //6-1-2 3:4:5 文件名称不能有:

	if err != nil {
		beego.Info("上传文件失败")
		return
	} else {
		c.SaveToFile("uploadname", "./static/img/"+filename)
	}

	//2.判断数据是否合法
	if artiContent == "" || artiName == "" {
		beego.Info("添加文章数据错误")
		return
	}
	//3.插入数据
	o := orm.NewOrm()
	arti := models.Article{}
	arti.ArtiName = artiName
	arti.Acontent = artiContent
	arti.Aimg = "./static/img/" + filename
	//查找type对象
	artiType := models.ArticleType{Id: id}
	o.Read(&artiType)

	arti.ArticleType = &artiType

	_, err = o.Insert(&arti)
	if err != nil {
		beego.Info("插入数据库错误")
		return
	}

	//4.返回文章界面
	c.Redirect("/index", 302)
}

//显示内容详情页面
func (c *MainController) ShowContent() {
	//1.获取文章ID
	id, err := c.GetInt("id")
	//beego.Info("id is ",id)
	if err != nil {
		beego.Info("获取文章ID错误", err)
		return
	}
	//2.查询数据库获取数据
	o := orm.NewOrm()
	arti := models.Article{Id: id}
	err = o.Read(&arti)
	if err != nil {
		beego.Info("查询错误", err)
		return
	}
	arti.Acount += 1

	m2m := o.QueryM2M(&arti, "User")
	user := models.User{}
	userName := c.GetSession("userName")
	user.Name = userName.(string)
	o.Read(&user, "Name")
	//多对多插入
	_, err = m2m.Add(&user)
	if err != nil {
		beego.Info("插入失败")
		return
	}
	o.Update(&arti)
	//显示多对多查询
	o.LoadRelated(&arti, "User") //第一种方法
	//o.QueryTable("Article").Filter("User__User__Name",userName.(string)).Distinct().All(&arti)	//第二种方法

	//3.传递数据给试图
	c.Data["article"] = arti
	c.Layout = "layout.html"
	c.TplName = "content.html"

}

//显示编辑界面
func (c *MainController) ShowUpdate() {
	//1.获取文章ID
	id, err := c.GetInt("id")
	//beego.Info("id is ",id)
	if err != nil {
		beego.Info("获取文章ID错误", err)
		return
	}
	//2.查询数据库获取数据
	o := orm.NewOrm()
	arti := models.Article{Id: id}
	err = o.Read(&arti)
	if err != nil {
		beego.Info("查询错误", err)
		return
	}
	//3.传递数据给试图
	c.Data["article"] = arti
	c.Layout = "layout.html"
	c.TplName = "update.html"
}

//处理更新业务数据
func (c *MainController) HandleUpdate() {
	//1.拿到数据
	id, _ := c.GetInt("id")
	artiName := c.GetString("articleName")
	content := c.GetString("content")
	f, h, err := c.GetFile("uploadname")
	var filename string
	if err != nil {
		beego.Info("上传文件失败")
		return
	} else {
		defer f.Close()

		//1.要限定格式
		fileext := path.Ext(h.Filename)
		if fileext != ".jpg" && fileext != "png" {
			beego.Info("上传文件格式错误")
			return
		}
		//2.限制大小
		if h.Size > 50000000 {
			beego.Info("上传文件过大")
			return
		}

		//3.需要对文件重命名，防止文件名重复
		filename = time.Now().Format("2006-01-02 15-04-05") + fileext //6-1-2 3:4:5
		c.SaveToFile("uploadname", "./static/img/"+filename)
	}

	//2.对数据进行一个处理
	if artiName == "" || content == "" {
		beego.Info("更新数据获取失败")
		return
	}

	//3.更新操作
	o := orm.NewOrm()
	arti := models.Article{Id: id}
	err = o.Read(&arti)
	if err != nil {
		beego.Info("查询数据错误")
		return
	}
	arti.ArtiName = artiName
	arti.Acontent = content
	arti.Aimg = "./static/img/" + filename

	_, err = o.Update(&arti)
	if err != nil {
		beego.Info("更新数据显示错误", err)
		return
	}
	//4.返回列表页面
	c.Redirect("/index", 302)
}

//删除操作
func (c *MainController) HandleDelete() {
	//1.拿到数据
	id, err := c.GetInt("id")
	if err != nil {
		beego.Info("获取id数据错误")
		return
	}

	//2.执行删除操作
	o := orm.NewOrm()
	arti := models.Article{Id: id}
	err = o.Read(&arti)
	if err != nil {
		beego.Info("查询错误")
		return
	}
	o.Delete(&arti)

	//3.返回列表页面
	c.Redirect("/index", 302)
}

//显示添加类型界面
func (c *MainController) ShowAddType() {
	o := orm.NewOrm()
	var artiTypes []models.ArticleType
	_, err := o.QueryTable("ArticleType").All(&artiTypes)
	if err != nil {
		beego.Info("没有获取到类型数据")
	}

	c.Data["articleType"] = artiTypes
	c.Layout = "layout.html"
	c.TplName = "addType.html"
}

//处理添加类型传输的信息
func (c *MainController) HandleAddType() {
	//1.获取内容
	typeName := c.GetString("typeName")
	//2.判断数据是否合法
	if typeName == "" {
		beego.Info("类型名为空")
		return
	}
	//3.写入数据
	o := orm.NewOrm()
	artiType := models.ArticleType{}
	artiType.Tname = typeName
	_, err := o.Insert(&artiType)
	if err != nil {
		beego.Info("插入类型错误")
		return
	}
	//4.返回界面
	c.Redirect("/addType", 302)
}

//删除文章类型
func (c *MainController) HandleDeleteType() {
	//1.拿到数据
	id, err := c.GetInt("id")
	if err != nil {
		beego.Info("获取id数据错误")
		return
	}

	//2.执行删除操作
	o := orm.NewOrm()
	artiType := models.ArticleType{Id: id}
	err = o.Read(&artiType)
	if err != nil {
		beego.Info("查询错误")
		return
	}
	o.Delete(&artiType)

	//3.返回列表页面
	c.Redirect("/addType", 302)
}

//处理文章类型更新
func (c *MainController) HandleUpdateType() {
	//1.拿到数据
	id, _ := c.GetInt("id")
	typeName := c.GetString("typeName")

	if typeName == "" {
		beego.Info("类型名称为空")
		c.Redirect("/addType", 302)
		return
	}
	//3.更新操作
	o := orm.NewOrm()
	artiType := models.ArticleType{Id: id, Tname: typeName}

	_, err := o.Update(&artiType)
	if err != nil {
		beego.Info("更新数据错误")
		c.Redirect("/addType", 302)
		return
	}
	//4.返回列表页面
	c.Redirect("/addType", 302)
}

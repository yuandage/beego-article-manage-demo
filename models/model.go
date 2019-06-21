package models

import (
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
	"time"
)

//表的设计
type User struct {
	Id      int
	Name    string     `orm:"unique"` //用户名
	Pwd     string     //密码
	Article []*Article `orm:"rel(m2m)"`
}

//文章结构体
type Article struct {
	Id       int       `orm:"pk;auto"`         //文章Id,主键,自增
	ArtiName string    `orm:"size(50)"`        //文章名称
	Atime    time.Time `orm:"auto_now"`        //文章时间
	Acount   int       `orm:"default(0);null"` //阅读量 允许为空
	Acontent string    //文章内容
	Aimg     string    //文章图片
	Atype    string    //文章类型

	ArticleType *ArticleType `orm:"rel(fk)"`
	User        []*User      `orm:"reverse(many)"`
}

//类型表
type ArticleType struct {
	Id      int
	Tname   string     //文章类型
	Article []*Article `orm:"reverse(many)"`
}

func init() {
	// 设置数据库基本信息
	orm.RegisterDataBase("default", "mysql", "root:165035@tcp(127.0.0.1:3306)/article_manage?charset=utf8")
	// 映射model数据
	orm.RegisterModel(new(User), new(Article), new(ArticleType))
	// 生成表
	//orm.RunSyncdb("default", false, true)
}

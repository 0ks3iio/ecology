package mysql

import (
	"errors"
	"strings"

	"github.com/astaxie/beego/orm"
	"github.com/goecology/ecology/appgo/model/mysql/store"
	"github.com/goecology/ecology/appgo/pkg/utils"
)

var tableCategory = "category"

// 分类
type Category struct {
	Id     int    `gorm:"not null;primary_key;AUTO_INCREMENT"json:"id"` //自增主键
	Pid    int    `gorm:"not null;"json:"pid"`                          //分类id
	Title  string `gorm:"not null;"json:"title"`                        //分类名称
	Intro  string `gorm:"not null;"json:"intro"`                        //介绍
	Icon   string `gorm:"not null;"json:"icon"`                         //分类icon
	Cnt    int    `gorm:"not null;"json:"cnt"`                          //分类下的文档项目统计
	Sort   int    `gorm:"not null;"json:"sort"`                         //排序
	Status bool   `gorm:"not null;"json:"status"`                       //分类状态，true表示显示，否则表示隐藏
	//PrintBookCount int    `orm:"default(0)" json:"print_book_count"`
	//WikiCount      int    `orm:"default(0)" json:"wiki_count"`
	//ArticleCount   int    `orm:"default(0)" json:"article_count"`
}

func (m *Category) TableName() string {
	return "category"
}

func NewCategory() *Category {
	return &Category{}
}

//新增分类
func (this *Category) AddCates(pid int, cates string) (err error) {
	slice := strings.Split(cates, "\n")
	if len(slice) == 0 {
		return
	}

	o := orm.NewOrm()
	for _, item := range slice {
		if item = strings.TrimSpace(item); item != "" {
			var cate = Category{
				Pid:    pid,
				Title:  item,
				Status: true,
			}
			if o.Read(&cate, "title"); cate.Id == 0 {
				_, err = orm.NewOrm().Insert(&cate)
			}
		}
	}
	return
}

//删除分类（如果分类下的文档项目不为0，则不允许删除）
func (this *Category) Del(id int) (err error) {
	var cate = Category{Id: id}

	o := orm.NewOrm()
	if err = o.Read(&cate); cate.Cnt > 0 { //当前分类下文档项目数量不为0，不允许删除
		return errors.New("删除失败，当前分类下的问下项目不为0，不允许删除")
	}

	if _, err = o.Delete(&cate, "id"); err != nil {
		return
	}
	_, err = o.QueryTable(tableCategory).Filter("pid", id).Delete()
	if err != nil { //删除分类图标
		return
	}

	switch utils.StoreType {
	case utils.StoreOss:
		store.ModelStoreOss.DelFromOss(cate.Icon)
	case utils.StoreLocal:
		store.ModelStoreLocal.DelFiles(cate.Icon)
	}
	return
}

//查询单个分类
func (this *Category) Find(id int) (cate Category) {
	cate.Id = id
	orm.NewOrm().Read(&cate)
	return cate
}
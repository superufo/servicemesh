package common

import (
	"errors"

	"github.com/go-mesh/openlogging"
	"github.com/jinzhu/gorm"

	_ "github.com/jinzhu/gorm/dialects/mysql"
)

var DB *gorm.DB

func init() {
	err := errors.New("")
	DB, err = gorm.Open("mysql", "gshop:T4dABtXMbs@tcp(111.229.150.102)/gshop_test?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		openlogging.Error("连接数据库失败", openlogging.WithTags(openlogging.Tags{
			"err": err.Error(),
		}))
		panic("连接数据库失败" + err.Error())
	}

	// 去掉复数  否则在表名后面加s
	DB.SingularTable(true)
}

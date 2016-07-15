package main

import (
	"github.com/astaxie/beego/orm"
	"zhongzi/save/mysql"
	"log"
)

func main() {
	orm.Debug = true

	oorm := orm.NewOrm()
	i_info := &mysql.MetaInfo{
		InfoHash:"1111111111111",
		BtInfo:"333333333333333333",
	}
	oorm.Insert(i_info)

	log.Println(i_info)
}


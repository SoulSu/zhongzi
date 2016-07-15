package mysql

import(
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
	"log"
)




var aliasName = "default"

func init() {
	err:=orm.RegisterDataBase(aliasName, "mysql", "root:@/bt?charset=utf8")
	if err != nil {
		log.Println(err.Error())
	}


	orm.SetMaxIdleConns(aliasName, 30)

	orm.SetMaxOpenConns(aliasName, 30)


	orm.RegisterModel(new(MetaInfo))

	orm.RunSyncdb(aliasName, false, true)
}

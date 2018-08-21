package models

import (
	"fenqi_v1/utils"
	"plugin"
	"zcm_tools/orm"

	_ "github.com/go-sql-driver/mysql"
)

func init() {
	if utils.RunMode == "release" {
		p, err := plugin.Open("../yg_orm/register_orm.so")
		if err != nil {
			panic(err)
		}
		f, err := p.Lookup("RegisterOrmFqV1")
		if err != nil {
			panic(err)
		}
		f.(func())()
	} else {
		orm.RegisterDataBase("default", "mysql", utils.MYSQL_URL)
		orm.RegisterDataBase("read", "mysql", utils.MYSQL_READ_URL)
		orm.RegisterDataBase("fq_log", "mysql", utils.MYSQL_LOG_URL)
		orm.RegisterDataBase("dm_xjfq", "mysql", utils.MYSQL_RPT_URL)
		//orm.Debug = true
	}
	orm.RegisterModel(
		new(OrderCredit),
		new(ConnectRecord),
		new(Salesman),
	)
}

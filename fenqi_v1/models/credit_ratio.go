package models

import (
	"time"

	"zcm_tools/orm"
)

//案件分配比例
type CreditRatio struct {
	Id         int       `orm:"column(id);pk"`
	PartWeight int       `orm:"column(part_weight);null" description:"拆分方案比例"`
	AllWeight  int       `orm:"column(all_weight);null" description:"完整方案比例"`
	Operator   string    `orm:"column(operator);size(20);null" description:"操作人"`
	CreateTime time.Time `orm:"column(create_time);type(datetime);null" description:"创建时间"`
	Remark     string    `orm:"column(remark);size(250);null" description:"修改时间"`
}

//添加案件分配比例
func AddCreditRatio(partWeight, allWeight int, operator, remark string) (err error) {
	o := orm.NewOrm()
	sql := `INSERT INTO credit_ratio(part_weight,all_weight,operator,create_time,remark) VALUES(?,?,?,NOW(),?)`
	_, err = o.Raw(sql, partWeight, allWeight, operator, remark).Exec()
	return
}

//查询所有分配比例
func QueryCreditRatios() (vs []CreditRatio, err error) {
	o := orm.NewOrm()
	sql := `SELECT * FROM credit_ratio ORDER BY create_time DESC `
	_, err = o.Raw(sql).QueryRows(&vs)
	return
}

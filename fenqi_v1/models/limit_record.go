package models

import (
	"time"
	"zcm_tools/orm"
)

type LimitRecord struct {
	Id           int       `orm:"column(id);auto"`
	Uid          int       `orm:"column(uid)" description:"授信审核表ID"`
	Remark       string    `orm:"column(remark);size(255);null" description:"备注"`
	OperatorTime time.Time `orm:"column(operator_time);type(datetime);null" description:"操作时间"`
	BalanceMoney int       `orm:"column(balance_money);null" description:"额度"`
	OperaId      int       `orm:"column(opera_id)" description:"操作人ID"`
	Displayname  string    `orm:"column(displayname)"`
}

//额度历史
func GetLimitRecordByUid(uid, start, pageSize int) (v []LimitRecord, err error) {
	o := orm.NewOrm()
	sql := `SELECT lr.uid,lr.remark,lr.operator_time,lr.balance_money,s.displayname FROM limit_record lr INNER JOIN sys_user AS s ON lr.opera_id = s.id WHERE uid = ? ORDER BY lr.operator_time DESC limit ?,?`
	_, err = o.Raw(sql, uid, start, pageSize).QueryRows(&v)
	return
}

//额度历史count
func GetLimitRecordCount(uid int) (count int, err error) {
	o := orm.NewOrm()
	sql := `SELECT COUNT(1) FROM limit_record lr LEFT JOIN sys_user AS s ON lr.opera_id = s.id WHERE uid = ? `
	err = o.Raw(sql, uid).QueryRow(&count)
	return
}

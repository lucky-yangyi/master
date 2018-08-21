package models

import (
	"time"
	"zcm_tools/orm"
)

//授信审核记录
type CreditAduitRecord struct {
	Id            int
	CreditAduitId int       //授信审核ID
	OperatorTime  time.Time //操作时间
	Content       string    //操作内容
	Remark        string    //备注
}

//增加授信分配记录
func AddCreditAduitRecord(creditAduitId, uid int, content, remark string) (err error) {
	sql := `INSERT into credit_aduit_record (credit_aduit_id,uid,content,remark,operator_time) values (?,?,?,?,now()) `
	_, err = orm.NewOrm().Raw(sql, creditAduitId, uid, content, remark).Exec()
	return err
}

//查询授信审核记录
func QueryCreditAduitRecord(uid int) (list []CreditAduitRecord, err error) {
	sql := `SELECT id,credit_aduit_id,operator_time,content,remark FROM credit_aduit_record WHERE uid = ? ORDER BY operator_time DESC,id DESC`
	o := orm.NewOrm()
	_, err = o.Raw(sql, uid).QueryRows(&list)
	return
}

//更新额度记录
func AddLimitRecord(uid, opera_id, balance_money int, remark string) (err error) {
	sql := `INSERT into limit_record (uid,opera_id,balance_money,remark,operator_time) values (?,?,?,?,now()) `
	_, err = orm.NewOrm().Raw(sql, uid, opera_id, balance_money, remark).Exec()
	return err
}

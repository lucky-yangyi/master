package models

import (
	"zcm_tools/orm"
)

//还款记录
type PaymentRecord struct {
	Id           int
	Uid          int     //用户ID
	CreateTime   string  //还款时间
	ReturnMoney  float64 //还款金额
	Days         int     //还款状态
	PaymentState string  //还款状态值
	Channel      int     //还款类型
	PaymentType  string  //还款类型值
	OrderNumber  string  //订单号
	OperatorId   int     //还款方式:小于0为主动还款
	State        string  //还款状态
	Remark       string  //备注
	Displayname  string  //操作人姓名
}

//查询用户还款记录
func QueryPaymentRecords(loanId int) (list []PaymentRecord, err error) {
	sql := `SELECT
				pr.id,
				pr.uid,
				pr.create_time,
				pr.return_money,
				TIMESTAMPDIFF(DAY, rs.loan_return_date,DATE(pr.create_time)) AS days,
				pr.channel
           FROM payment_record AS pr
		   LEFT JOIN repayment_schedule AS rs
		   ON pr.repayment_schedule_id = rs.id
		   WHERE pr.loan_id = ?
		   ORDER BY pr.create_time DESC`
	o := orm.NewOrm()
	_, err = o.Raw(sql, loanId).QueryRows(&list)
	return
}

//查询用户还款记录（带分页）
func QueryUsersManagePaymentRecords(uid, start, pageSize int) (list []PaymentRecord, err error) {
	sql := `SELECT
				pr.id,
				pr.uid,
				pr.create_time,
				pr.return_money,
				pr.channel,
				pr.order_number,
				pr.operator_id,
				pr.state,
				pr.remark,
				su.displayname
           FROM payment_record AS pr
           LEFT JOIN sys_user AS su
           ON pr.operator_id = su.id
		   WHERE pr.uid = ?
		   ORDER BY pr.create_time DESC
           LIMIT ?,?`
	o := orm.NewOrm()
	_, err = o.Raw(sql, uid, start, pageSize).QueryRows(&list)
	return
}

//查询用户还款记录总数
func QueryUsersManagePaymentRecordCount(uid int) (count int, err error) {
	sql := `SELECT
				COUNT(1)
           FROM payment_record
		   WHERE uid = ?`
	o := orm.NewOrm()
	err = o.Raw(sql, uid).QueryRow(&count)
	return
}

//更新还款记录状态
func UpdatePaymentRecordState(id int, state string) (err error) {
	o := orm.NewOrm()
	sql := `UPDATE payment_record SET state = ? WHERE id = ?`
	_, err = o.Raw(sql, state, id).Exec()
	return
}

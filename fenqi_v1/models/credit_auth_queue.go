package models

import (
	"fenqi_v1/utils"
	"strconv"
	"time"
	"zcm_tools/orm"
)

//获取授信中预约的数量
func GetAuthNumCreditQueueAp() (count int) {
	sql := `SELECT COUNT(1) from sys_user as a,credit_auth_queue as b WHERE a.id = b.operator_id AND b.auth_state = "OUTQUEUE"`
	orm.NewOrm().Raw(sql).QueryRow(&count)
	return
}

//获取预约时间队列id
func GetAuthAppointmentQueueId() (v CreditAduit, err error) {
	sql := `SELECT b.id,b.uid,b.inqueue_type,b.inqueue_time,b.credit_aduit_id FROM credit_auth_queue as b WHERE b.auth_state = "OUTQUEUE" AND b.inqueue_time IS NOT NULL AND b.inqueue_time <= NOW() AND b.inqueue_type > 0 ORDER BY inqueue_time LIMIT 1`
	err = orm.NewOrm().Raw(sql).QueryRow(&v)
	return
}

//获取授信排队中信息id
func GetCreditAuthQueueUpId() (v CreditAduit, err error) {
	sql := `SELECT b.id,b.uid,b.credit_aduit_id FROM credit_auth_queue as b WHERE b.auth_state = "QUEUEING" ORDER BY b.queue_time LIMIT 1`
	err = orm.NewOrm().Raw(sql).QueryRow(&v)
	return
}

//获取排队中数量
func GetCreditAuthQueueUpIdCount() (count int) {
	sql := `SELECT COUNT(1) FROM credit_auth_queue as b WHERE b.auth_state = 'QUEUEING' ORDER BY b.queue_time`
	orm.NewOrm().Raw(sql).QueryRow(&count)
	return
}

//更新授信中记录
func UpdateAuthCueditQueueRemark(remark string, credit_aduit_id int) error {
	sql := `UPDATE credit_auth_queue SET remark = ? WHERE credit_aduit_id =? `
	_, err := orm.NewOrm().Raw(sql, remark, credit_aduit_id).Exec()
	return err
}

//Pass状态下更新授信
func UpdateAuthCueditQueuePassStatus(auth_state, remark string, auth_balance_money, operator_id, credit_aduit_id int) error {
	sql := `UPDATE credit_auth_queue SET auth_state=?,remark = ?,auth_balance_money=?,operator_id=?,handling_time=NOW() WHERE credit_aduit_id=? `
	_, err := orm.NewOrm().Raw(sql, auth_state, remark, auth_balance_money, operator_id, credit_aduit_id).Exec()
	return err
}

//更新退回队列时间与状态
func UpdateAuthCreditOutqueueTime(credit_aduit_id, operator_id int, inqueue_time string) (err error) {
	sql := `UPDATE credit_auth_queue SET auth_state = "OUTQUEUE"`
	if inqueue_time == "" {
		sql += `,inqueue_time = NULL,operator_id=?,inqueue_type=0,handling_time=NOW() WHERE credit_aduit_id =?`
		_, err = orm.NewOrm().Raw(sql, operator_id, credit_aduit_id).Exec()

	} else {
		sql += `,inqueue_time=?,operator_id=?,inqueue_type=1,handling_time=NOW() WHERE credit_aduit_id =?`
		_, err = orm.NewOrm().Raw(sql, inqueue_time, operator_id, credit_aduit_id).Exec()
	}
	return err
}

//更新分配時間
func UpdateAuthCreditAlloctionTime(credit_aduit_id int) error {
	sql := `UPDATE credit_auth_queue SET allocation_time = NOW(),inqueue_time = NULL WHERE credit_aduit_id = ? `
	_, err := orm.NewOrm().Raw(sql, credit_aduit_id).Exec()
	return err
}

//统计45分钟之内
func QueryAuthCreditHandingCountIn(credit_aduit_id int) (count int) {
	sql := `SELECT COUNT(1) FROM credit_auth_queue WHERE credit_aduit_id =? AND allocation_time >= DATE_SUB(NOW(),INTERVAL 45 MINUTE)`
	orm.NewOrm().Raw(sql, credit_aduit_id).QueryRow(&count)
	return
}

//超时清缓存并入队列
func UpdateAuthCreditQueueing(credit_aduit_id int) error {
	sql := `UPDATE credit_auth_queue SET auth_state = "QUEUEING",queue_time = NOW() WHERE credit_aduit_id = ? `
	_, err := orm.NewOrm().Raw(sql, credit_aduit_id).Exec()
	return err
}

//超时分配查询attime时间
func QueryAuthCreditAttime(credit_aduit_id int) (v CreditAduit) {
	sql := "SELECT allocation_time FROM credit_auth_queue WHERE credit_aduit_id=?"
	orm.NewOrm().Raw(sql, credit_aduit_id).QueryRow(&v)
	return
}

//更新授信中状态
func UpdateCueditAuthQueueStatusInqueueTime(auth_state, auth_name string, operator_id, credit_aduit_id int) error {
	sql := `UPDATE credit_auth_queue SET auth_state=?,auth_name=?,inqueue_time=NULL,operator_id=? WHERE credit_aduit_id =? `
	_, err := orm.NewOrm().Raw(sql, auth_state, auth_name, operator_id, credit_aduit_id).Exec()
	return err
}

//更新is_auth
func UpdateCreditIsAuth(id int) error {
	sql := `UPDATE credit_aduit SET is_auth=1 WHERE id =? `
	_, err := orm.NewOrm().Raw(sql, id).Exec()
	return err
}

//更新授信中处理人状态
func UpdateAuthCueditQueueStatusOp(auth_state, remark string, credit_aduit_id, operator_id int) error {
	sql := `UPDATE credit_auth_queue SET auth_state = ?,remark=?,operator_id=?,handling_time=NOW() WHERE credit_aduit_id =? `
	_, err := orm.NewOrm().Raw(sql, auth_state, remark, operator_id, credit_aduit_id).Exec()
	return err
}

//更新state
func UpdateAuthCreditIsState(id int) error {
	sql := `UPDATE credit_auth_queue SET auth_state="OUTQUEUE" WHERE id =? `
	_, err := orm.NewOrm().Raw(sql, id).Exec()
	return err
}

//添加授信中记录
func InsetAuthCueditAdviseRemark(remark, state string, credit_cut_id int) error {
	sql := `INSERT into credit_cut_advise (credit_cut_id,state,remark,operator_time) values (?,?,?,now()) `
	_, err := orm.NewOrm().Raw(sql, credit_cut_id, state, remark).Exec()
	return err
}

//更新退回队列中入队类型(排队,插队)
func UpdateAuthInqueueType(id int, queue_time time.Time) (err error) {
	sql := `UPDATE credit_auth_queue SET auth_state = "QUEUEING",queue_time=? WHERE id =?`
	_, err = orm.NewOrm().Raw(sql, queue_time, id).Exec()
	return
}

//45分钟 关机 睡眠 处理机制
func CreditAuthHandingLogOut() {
	var id []int
	sql := `SELECT id FROM credit_auth_queue WHERE state = "HANDING"`
	orm.NewOrm().Raw(sql).QueryRows(&id)
	for _, v := range id {
		if !utils.Rc.IsExist("auth:" + utils.CacheKeyCreditMessage + "_" + strconv.Itoa(v)) {
			sql := `UPDATE credit_auth_queue SET state = "QUEUEING",queue_time = NOW() WHERE id = ? `
			orm.NewOrm().Raw(sql, v).Exec()
		}

	}
}

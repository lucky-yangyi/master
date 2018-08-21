package models

import (
	"time"
	"zcm_tools/orm"
)

type RepaymentSchedule struct {
	Id             int
	Uid            int
	Amount         float64   //合计
	LoanReturnDate time.Time //还款日期
	CapitalAmount  float64   //本金
	TaxAmount      float64   //利息
	Shouldmoney    float64   //应还金额
	Term_no        int       //期数
}

//还款计划
func GetRepaymentScheduleByRepayId(loanId int) (v []RepaymentSchedule, err error) {
	o := orm.NewOrm()
	sql := `SELECT loan_return_date,capital_amount,tax_amount,(capital_amount + tax_amount) AS amount FROM repayment_schedule WHERE loan_id  = ? ORDER BY loan_return_date`
	_, err = o.Raw(sql, loanId).QueryRows(&v)
	return
}

//更新还款计划延期处理天数
func UpdateAdjournHandlingDay(adjournHandlingDay, id int) (err error) {
	o := orm.NewOrm()
	sql := `UPDATE repayment_schedule SET adjourn_handling_day=? WHERE id = ?`
	_, err = o.Raw(sql, adjournHandlingDay, id).Exec()
	return
}

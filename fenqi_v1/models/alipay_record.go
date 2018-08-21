package models

import (
	"time"
	"zcm_tools/orm"
)

type AlipayRecord struct {
	Id             int
	OidPaybill     string
	ReturnTime     time.Time
	Remark         string
	AmountIncome   float64
	AmountOutlay   float64
	IsDeal         int
	Account        string
	CreateTime     time.Time
	Remark1        string
	Remark2        string
	StateString    string
	ExtraRepayment float64
	OperatorIpDeal string //处理的操作人IP
}

func GetAlipayRecordMoreThanZ(condition string, isLimit bool, begin, count int, pars ...string) (a []AlipayRecord, err error) {
	o := orm.NewOrm()
	sql := `SELECT
				id,
				oid_paybill,
				return_time,
				amount_income,
				remark,
				remark1,
				remark2,
				extra_repayment,
				operator_ip_deal,
				is_deal
			FROM alipay_record WHERE amount_income > 0 `
	if condition != "" {
		sql += condition
	}
	if isLimit {
		sql += " ORDER BY return_time DESC  LIMIT ?, ?"
		_, err = o.Raw(sql, pars, begin, count).QueryRows(&a)
	} else {
		sql += " ORDER BY return_time DESC"
		_, err = o.Raw(sql, pars).QueryRows(&a)
	}
	return
}
func GetAlipayRecordMoreThanZ2(condition string, pars ...string) (a []AlipayRecord, err error) {
	o := orm.NewOrm()
	sql := `SELECT
				id,
				oid_paybill,
				return_time,
				amount_income,
				remark,
				remark1,
				remark2,
				extra_repayment,
				operator_ip_deal,
				is_deal
			FROM alipay_record WHERE amount_income > 0 `
	if condition != "" {
		sql += condition
	}
	sql += " ORDER BY return_time DESC "
	_, err = o.Raw(sql, pars).QueryRows(&a)
	return
}
func GetAlipayRecordMoreThanZSumReturnMoney(condition string, pars ...string) (sumRetrunMoney float64) {
	o := orm.NewOrm()
	sql := `SELECT sum(amount_income) from alipay_record where amount_income > 0 `
	if condition != "" {
		sql += condition
	}
	o.Raw(sql, pars).QueryRow(&sumRetrunMoney)
	return
}

func GetAlipayRecordCount(condition string, pars ...string) int {
	o := orm.NewOrm()
	sql := `SELECT count(1) from alipay_record where amount_income > 0 `
	if condition != "" {
		sql += condition
	}
	var count int
	o.Raw(sql, pars).QueryRow(&count)
	return count
}
func GetAlipayRecordRemark2(id int) string {
	o := orm.NewOrm()
	sql := `SELECT remark2 from alipay_record where id = ?`
	var remark2 string
	o.Raw(sql, id).QueryRow(&remark2)
	return remark2
}

func UpdateAlipayRecord(isDeal, id, operatorIdRepay int, remark, operatorIp string) error {
	o := orm.NewOrm()
	sql := `UPDATE alipay_record  SET  is_deal = ? ,remark1 = ?,operator_id_repay = ?,operator_time_repay = NOW(),operator_ip_repay = ? where id = ?`
	_, err := o.Raw(sql, isDeal, remark, operatorIdRepay, operatorIp, id).Exec()
	return err
}

func UpdateAlipayRecord2(isDeal, id int, remark string) error {
	o := orm.NewOrm()
	sql := `UPDATE alipay_record  SET  is_deal = ? ,remark1 = ? where id = ?`
	_, err := o.Raw(sql, isDeal, remark, id).Exec()
	return err
}

func UpdateAlipayRecordRemark2(remark2 string, extra float64, operatorIp string, id, operatorIdDeal int) error {
	o := orm.NewOrm()
	sql := `UPDATE alipay_record  SET remark2 = ? ,extra_repayment=?,operator_id_deal = ?,operator_time_deal = NOW(),operator_ip_deal = ? where id = ?`
	_, err := o.Raw(sql, remark2, extra, operatorIdDeal, operatorIp, id).Exec()
	return err
}

func UpdateAlipayIsDeal(isDeal, id int) error {
	o := orm.NewOrm()
	sql := `UPDATE alipay_record  SET is_deal = ? where id = ?`
	_, err := o.Raw(sql, isDeal, id).Exec()
	return err
}

func GetAlipayRecordMoreThanZSumExtraRepayment(condition string, pars ...string) (sumExtraRepayment float64) {
	o := orm.NewOrm()
	sql := `SELECT sum(extra_repayment) from alipay_record where extra_repayment > 0 `
	if condition != "" {
		sql += condition
	}
	o.Raw(sql, pars).QueryRow(&sumExtraRepayment)
	return
}

//更新处理结果
func UpdateAlipayRecordRemark(id, operatorIdMotifyResult int, newRemark, operatorIp string, extra_repayment float64) (err error) {
	sql := ``
	if extra_repayment == 0 {
		sql += `UPDATE alipay_record SET remark2 = ?,operator_id_motify_result = ?,operator_time_motify_result = NOW(),operator_ip_motify_result = ? WHERE id = ?`
		_, err = orm.NewOrm().Raw(sql, newRemark, operatorIdMotifyResult, operatorIp, id).Exec()
	} else {
		sql += `UPDATE alipay_record SET remark2 = ?,operator_id_motify_result = ?,operator_time_motify_result = NOW(),operator_ip_motify_result = ?,extra_repayment=? WHERE id = ?`
		_, err = orm.NewOrm().Raw(sql, newRemark, operatorIdMotifyResult, operatorIp, extra_repayment, id).Exec()
	}
	return
}

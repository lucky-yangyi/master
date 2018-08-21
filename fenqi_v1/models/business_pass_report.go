package models

import (
	"time"
	"zcm_tools/orm"
)

type PassTotalList struct {
	Createtime      time.Time //日期
	NumRegister     int       //注册用户
	MunCertSucess   int       //认证通过用户
	MunCreditSucess int       //授信通过用户
	NumLoanSucess   int       //借款成功用户
	OrderLoanSucess int       //借款成功订单
	MoneyLoanSucess float64   //借款成功金额
}

type PassTotalData struct {
	TotalNumRegister int     //累计注册用户
	TotalNumCert     int     //累计认证用户
	TotalMunCredit   int     //累计授信用户
	TotalNumLoan     int     //累计借款用户
	TotalOrderLoan   int     //累计借款成功订单
	TotalMoneyLoan   float64 //累计借款成功金额
}

//获取通过汇总累计数据
func GetTotalData(condition string, pars ...interface{}) (data PassTotalData, err error) {
	o := orm.NewOrm()
	o.Using("dm_xjfq")
	sql := `SELECT
				SUM(total_num_register) AS total_num_register,
				SUM(total_num_cert) AS total_num_cert,
				SUM(total_mun_credit) AS total_mun_credit,
				SUM(total_num_loan) AS total_num_loan,
				SUM(total_order_loan) AS total_order_loan,
				SUM(total_money_loan) AS total_money_loan
			FROM
				HT_summary_pass_rate_hourly
            WHERE 1=1`
	if condition != "" {
		sql += condition
	}
	sql += ` GROUP BY
				createtime
			ORDER BY
				createtime DESC
			LIMIT 1 `
	err = o.Raw(sql, pars).QueryRow(&data)
	return
}

//获取通过汇总列表
func GetPassTotalList(start, pageSize int, isLimit bool, condition string, pars ...interface{}) (list []PassTotalList, err error) {
	o := orm.NewOrm()
	o.Using("dm_xjfq")
	sql := ` SELECT
				createtime,
				SUM(num_register) AS num_register,
				SUM(mun_cert_sucess) AS mun_cert_sucess,
				SUM(mun_credit_sucess) AS mun_credit_sucess,
				SUM(num_loan_sucess) AS num_loan_sucess,
				SUM(order_loan_sucess) AS order_loan_sucess,
				SUM(money_loan_sucess) AS money_loan_sucess
			FROM
				HT_summary_pass_rate_hourly
			WHERE
				1 = 1 `
	if condition != "" {
		sql += condition
	}
	sql += ` GROUP BY createtime ORDER BY createtime DESC`
	if isLimit {
		sql += ` LIMIT ? ,?`
		_, err = o.Raw(sql, pars, start, pageSize).QueryRows(&list)
	} else {
		_, err = o.Raw(sql, pars).QueryRows(&list)
	}
	return
}

//获取通过汇总列表总数
func PassTotalListCount(condition string, pars ...interface{}) (count int, err error) {
	o := orm.NewOrm()
	o.Using("dm_xjfq")
	sql := ``
	sql += `SELECT COUNT(DISTINCT createtime) as count FROM HT_summary_pass_rate_hourly WHERE 1=1 `
	if condition != "" {
		sql += condition
	}
	err = o.Raw(sql, pars).QueryRow(&count)
	return
}

type BusinessData struct {
	Createtime string
	Count      float64
}

//通过汇总折线图数据
func GetPassTotalLineData(chooseType int, condition string, pars ...interface{}) (list []BusinessData, err error) {
	o := orm.NewOrm()
	o.Using("dm_xjfq")
	sql := `SELECT createtime `
	switch chooseType {
	case 0:
		sql += ` ,SUM(num_register) as count`
	case 1:
		sql += ` ,SUM(mun_cert_sucess) as count`
	case 2:
		sql += ` ,SUM(mun_credit_sucess) as count`
	case 3:
		sql += ` ,SUM(num_loan_sucess) as count`
	case 4:
		sql += ` ,SUM(order_loan_sucess) as count`
	case 5:
		sql += ` ,SUM(money_loan_sucess) as count`
	}
	sql += ` FROM HT_summary_pass_rate_hourly WHERE 1=1 `
	if condition != "" {
		sql += condition
	}
	sql += ` GROUP BY createtime `
	_, err = o.Raw(sql, pars).QueryRows(&list)
	return
}

type ProcessConversionList struct {
	Createtime           time.Time //日期
	NumRegister          int       //注册_人数
	NumCertApply         int       //认证申请_人数
	NumCertSucess        int       //认证通过_人数
	NumLivingSucess      int       //活体通过_人数
	NumUsersBaseInfo     int       //个人信息补充_人数
	NumZmAuth            int       //芝麻信用_人数
	NumLinkMan           int       //常用联系人_人数
	NumMobileOperatorsMx int       //运营商认证_人数
	NumBindCard          int       //收款银行卡_人数
	NumGjj               int       //公积金_人数
	NumAliPay            int       //支付宝_人数
	NumCreditApply       int       //授信申请_人数
	NumCreditSucess      int       //授信通过_人数
	NumFirstLoanApply    int       //首次借款申请_人数
	NumFirstLoanSucess   int       //首次借款通过_人数
}

//流程转化率数据
func GetProcessConversionList(changeCycle, start, pageSize int, isLimit bool, condition string, pars ...interface{}) (list []ProcessConversionList, err error) {
	o := orm.NewOrm()
	o.Using("dm_xjfq")
	sql := `SELECT createtime,SUM(num_register) AS num_register, `
	if changeCycle == 1 {
		sql += `  SUM(num_cert_apply_one) AS num_cert_apply,
				SUM(num_cert_sucess_one) AS num_cert_sucess,
				SUM(num_living_sucess_one) AS num_living_sucess,
				SUM(num_users_base_info_one) AS num_users_base_info,
				SUM(num_zm_auth_one) AS num_zm_auth,
				SUM(num_link_man_one) AS num_link_man,
				SUM(num_mobile_operators_mx_one) AS num_mobile_operators_mx,
				SUM(num_bind_card_one) AS num_bind_card,
				SUM(num_gjj_one) AS num_gjj,
				SUM(num_ali_pay_one) AS num_ali_pay,
				SUM(num_credit_apply_one) AS num_credit_apply,
				SUM(num_credit_sucess_one) AS num_credit_sucess,
				SUM(num_first_loan_apply_one) AS num_first_loan_apply,
				SUM(num_first_loan_sucess_one) AS num_first_loan_sucess
				FROM HT_users_conver1_daily `
	}
	if changeCycle == 3 {
		sql += ` SUM(num_cert_apply_three) AS num_cert_apply,
				SUM(num_cert_sucess_three) AS num_cert_sucess,
				SUM(num_living_sucess_three) AS num_living_sucess,
				SUM(num_users_base_info_three) AS num_users_base_info,
				SUM(num_zm_auth_three) AS num_zm_auth,
				SUM(num_link_man_three) AS num_link_man,
				SUM(num_mobile_operators_mx_three) AS num_mobile_operators_mx,
				SUM(num_bind_card_three) AS num_bind_card,
				SUM(num_gjj_three) AS num_gjj,
				SUM(num_ali_pay_three) AS num_ali_pay,
				SUM(num_credit_apply_three) AS num_credit_apply,
				SUM(num_credit_sucess_three) AS num_credit_sucess,
				SUM(num_first_loan_apply_three) AS num_first_loan_apply,
				SUM(num_first_loan_sucess_three) AS num_first_loan_sucess
				FROM HT_users_conver3_daily `
	}
	if changeCycle == 7 {
		sql += ` SUM(num_cert_apply_seven) as num_cert_apply,
				 SUM(num_cert_sucess_seven) as num_cert_sucess,
				 SUM(num_living_sucess_seven) as num_living_sucess,
				 SUM(num_users_base_info_seven) as num_users_base_info,
				 SUM(num_zm_auth_seven) as num_zm_auth,
				 SUM(num_link_man_seven) as num_link_man,
				 SUM(num_mobile_operators_mx_seven) as num_mobile_operators_mx,
				 SUM(num_bind_card_seven) as num_bind_card,
				 SUM(num_gjj_seven) as num_gjj,
				 SUM(num_ali_pay_seven) as num_ali_pay,
				 SUM(num_credit_apply_seven) as num_credit_apply,
				 SUM(num_credit_sucess_seven) as num_credit_sucess,
				 SUM(num_first_loan_apply_seven) as num_first_loan_apply,
				 SUM(num_first_loan_sucess_seven) as num_first_loan_sucess
				 FROM HT_users_conver7_daily`
	}
	if changeCycle == 30 {
		sql += ` SUM(num_cert_apply_thirty) as num_cert_apply,
				 SUM(num_cert_sucess_thirty) as num_cert_sucess,
				 SUM(num_living_sucess_thirty) as num_living_sucess,
				 SUM(num_users_base_info_thirty) as num_users_base_info,
				 SUM(num_zm_auth_thirty) as num_zm_auth,
				 SUM(num_link_man_thirty) as num_link_man,
				 SUM(num_mobile_operators_mx_thirty) as num_mobile_operators_mx,
				 SUM(num_bind_card_thirty) as num_bind_card,
				 SUM(num_gjj_thirty) as num_gjj,
				 SUM(num_ali_pay_thirty) as num_ali_pay,
				 SUM(num_credit_apply_thirty) as num_credit_apply,
				 SUM(num_credit_sucess_thirty) as num_credit_sucess,
				 SUM(num_first_loan_apply_thirty) as num_first_loan_apply,
				 SUM(num_first_loan_sucess_thirty) as num_first_loan_sucess
				 FROM HT_users_conver30_daily`
	}
	sql += ` WHERE 1=1 `
	if condition != "" {
		sql += condition
	}
	sql += ` GROUP BY createtime ORDER BY createtime DESC`
	if isLimit {
		sql += ` LIMIT ? ,?`
		_, err = o.Raw(sql, pars, start, pageSize).QueryRows(&list)
	} else {
		_, err = o.Raw(sql, pars).QueryRows(&list)
	}
	return
}

//获取流程转化率总数
func GetProcessConversionCount(changeCycle int, condition string, pars ...interface{}) (count int, err error) {
	o := orm.NewOrm()
	o.Using("dm_xjfq")
	sql := `SELECT COUNT(DISTINCT createtime) as count FROM `
	if changeCycle == 1 {
		sql += ` HT_users_conver1_daily `
	}
	if changeCycle == 3 {
		sql += ` HT_users_conver3_daily `
	}
	if changeCycle == 7 {
		sql += ` HT_users_conver7_daily `
	}
	if changeCycle == 30 {
		sql += ` HT_users_conver30_daily `
	}
	sql += ` WHERE 1=1 `
	if condition != "" {
		sql += condition
	}
	err = o.Raw(sql, pars).QueryRow(&count)
	return
}

//获取流程转化率折线图
func GetProcessConversionLineData(chooseType, changeCycle int, dataType, condition string, pars ...interface{}) (list []BusinessData, err error) {
	o := orm.NewOrm()
	o.Using("dm_xjfq")
	sql := `SELECT createtime, `
	if changeCycle == 1 {
		if dataType == "用户人数" {
			switch chooseType {
			case 0:
				sql += `SUM(num_register) as count `
			case 1:
				sql += `SUM(num_cert_apply_one) as count `
			case 2:
				sql += `SUM(num_cert_sucess_one) as count `
			case 3:
				sql += `SUM(num_living_sucess_one) as count `
			case 4:
				sql += `SUM(num_users_base_info_one) as count `
			case 5:
				sql += `SUM(num_zm_auth_one) as count `
			case 6:
				sql += `SUM(num_link_man_one) as count `
			case 7:
				sql += `SUM(num_mobile_operators_mx_one) as count `
			case 8:
				sql += `SUM(num_bind_card_one) as count `
			case 9:
				sql += `SUM(num_gjj_one) as count `
			case 10:
				sql += `SUM(num_ali_pay_one) as count `
			case 11:
				sql += `SUM(num_credit_apply_one) as count `
			case 12:
				sql += `SUM(num_credit_sucess_one) as count `
			case 13:
				sql += `SUM(num_first_loan_apply_one) as count `
			case 14:
				sql += `SUM(num_first_loan_sucess_one) as count `
			}
		} else {
			switch chooseType {
			case 0:
				sql += `SUM(num_register) as count `
			case 1:
				sql += `FORMAT(SUM(num_cert_apply_one)/SUM(num_register)*100,2) as count `
			case 2:
				sql += `FORMAT(SUM(num_cert_sucess_one)/SUM(num_register)*100,2) as count `
			case 3:
				sql += `FORMAT(SUM(num_living_sucess_one)/SUM(num_register)*100,2) as count `
			case 4:
				sql += `FORMAT(SUM(num_users_base_info_one)/SUM(num_register)*100,2) as count `
			case 5:
				sql += `FORMAT(SUM(num_zm_auth_one)/SUM(num_register)*100,2) as count `
			case 6:
				sql += `FORMAT(SUM(num_link_man_one)/SUM(num_register)*100,2) as count `
			case 7:
				sql += `FORMAT(SUM(num_mobile_operators_mx_one)/SUM(num_register)*100,2) as count `
			case 8:
				sql += `FORMAT(SUM(num_bind_card_one)/SUM(num_register)*100,2) as count `
			case 9:
				sql += `FORMAT(SUM(num_gjj_one)/SUM(num_register)*100,2) as count `
			case 10:
				sql += `FORMAT(SUM(num_ali_pay_one)/SUM(num_register)*100,2) as count `
			case 11:
				sql += `FORMAT(SUM(num_credit_apply_one)/SUM(num_register)*100,2) as count `
			case 12:
				sql += `FORMAT(SUM(num_credit_sucess_one)/SUM(num_register)*100,2) as count `
			case 13:
				sql += `FORMAT(SUM(num_first_loan_apply_one)/SUM(num_register)*100,2) as count `
			case 14:
				sql += `FORMAT(SUM(num_first_loan_sucess_one)/SUM(num_register)*100,2) as count `
			}
		}
		sql += ` FROM HT_users_conver1_daily`
	}
	if changeCycle == 3 {
		if dataType == "用户人数" {
			switch chooseType {
			case 0:
				sql += `SUM(num_register) as count `
			case 1:
				sql += `SUM(num_cert_apply_three) as count `
			case 2:
				sql += `SUM(num_cert_sucess_three) as count `
			case 3:
				sql += `SUM(num_living_sucess_three) as count `
			case 4:
				sql += `SUM(num_users_base_info_three) as count `
			case 5:
				sql += `SUM(num_zm_auth_three) as count `
			case 6:
				sql += `SUM(num_link_man_three) as count `
			case 7:
				sql += `SUM(num_mobile_operators_mx_three) as count `
			case 8:
				sql += `SUM(num_bind_card_three) as count `
			case 9:
				sql += `SUM(num_gjj_three) as count `
			case 10:
				sql += `SUM(num_ali_pay_three) as count `
			case 11:
				sql += `SUM(num_credit_apply_three) as count `
			case 12:
				sql += `SUM(num_credit_sucess_three) as count `
			case 13:
				sql += `SUM(num_first_loan_apply_three) as count `
			case 14:
				sql += `SUM(num_first_loan_sucess_three) as count `
			}
		} else {
			switch chooseType {
			case 0:
				sql += `SUM(num_register) as count `
			case 1:
				sql += `FORMAT(SUM(num_cert_apply_three)/SUM(num_register)*100,2) as count `
			case 2:
				sql += `FORMAT(SUM(num_cert_sucess_three)/SUM(num_register)*100,2) as count `
			case 3:
				sql += `FORMAT(SUM(num_living_sucess_three)/SUM(num_register)*100,2) as count `
			case 4:
				sql += `FORMAT(SUM(num_users_base_info_three)/SUM(num_register)*100,2) as count `
			case 5:
				sql += `FORMAT(SUM(num_zm_auth_three)/SUM(num_register)*100,2) as count `
			case 6:
				sql += `FORMAT(SUM(num_link_man_three)/SUM(num_register)*100,2) as count `
			case 7:
				sql += `FORMAT(SUM(num_mobile_operators_mx_three)/SUM(num_register)*100,2) as count `
			case 8:
				sql += `FORMAT(SUM(num_bind_card_three)/SUM(num_register)*100,2) as count `
			case 9:
				sql += `FORMAT(SUM(num_gjj_three)/SUM(num_register)*100,2) as count `
			case 10:
				sql += `FORMAT(SUM(num_ali_pay_three)/SUM(num_register)*100,2) as count `
			case 11:
				sql += `FORMAT(SUM(num_credit_apply_three)/SUM(num_register)*100,2) as count `
			case 12:
				sql += `FORMAT(SUM(num_credit_sucess_three)/SUM(num_register)*100,2) as count `
			case 13:
				sql += `FORMAT(SUM(num_first_loan_apply_three)/SUM(num_register)*100,2) as count `
			case 14:
				sql += `FORMAT(SUM(num_first_loan_sucess_three)/SUM(num_register)*100,2) as count `
			}
		}
		sql += ` FROM HT_users_conver3_daily`
	}
	if changeCycle == 7 {
		if dataType == "用户人数" {
			switch chooseType {
			case 0:
				sql += `SUM(num_register) as count `
			case 1:
				sql += `SUM(num_cert_apply_seven) as count `
			case 2:
				sql += `SUM(num_cert_sucess_seven) as count `
			case 3:
				sql += `SUM(num_living_sucess_seven) as count `
			case 4:
				sql += `SUM(num_users_base_info_seven) as count `
			case 5:
				sql += `SUM(num_zm_auth_seven) as count `
			case 6:
				sql += `SUM(num_link_man_seven) as count `
			case 7:
				sql += `SUM(num_mobile_operators_mx_seven) as count `
			case 8:
				sql += `SUM(num_bind_card_seven) as count `
			case 9:
				sql += `SUM(num_gjj_seven) as count `
			case 10:
				sql += `SUM(num_ali_pay_seven) as count `
			case 11:
				sql += `SUM(num_credit_apply_seven) as count `
			case 12:
				sql += `SUM(num_credit_sucess_seven) as count `
			case 13:
				sql += `SUM(num_first_loan_apply_seven) as count `
			case 14:
				sql += `SUM(num_first_loan_sucess_seven) as count `
			}
		} else {
			switch chooseType {
			case 0:
				sql += `SUM(num_register) as count `
			case 1:
				sql += `FORMAT(SUM(num_cert_apply_seven)/SUM(num_register)*100,2) as count `
			case 2:
				sql += `FORMAT(SUM(num_cert_sucess_seven)/SUM(num_register)*100,2) as count `
			case 3:
				sql += `FORMAT(SUM(num_living_sucess_seven)/SUM(num_register)*100,2) as count `
			case 4:
				sql += `FORMAT(SUM(num_users_base_info_seven)/SUM(num_register)*100,2) as count `
			case 5:
				sql += `FORMAT(SUM(num_zm_auth_seven)/SUM(num_register)*100,2) as count `
			case 6:
				sql += `FORMAT(SUM(num_link_man_seven)/SUM(num_register)*100,2) as count `
			case 7:
				sql += `FORMAT(SUM(num_mobile_operators_mx_seven)/SUM(num_register)*100,2) as count `
			case 8:
				sql += `FORMAT(SUM(num_bind_card_seven)/SUM(num_register)*100,2) as count `
			case 9:
				sql += `FORMAT(SUM(num_gjj_seven)/SUM(num_register)*100,2) as count `
			case 10:
				sql += `FORMAT(SUM(num_ali_pay_seven)/SUM(num_register)*100,2) as count `
			case 11:
				sql += `FORMAT(SUM(num_credit_apply_seven)/SUM(num_register)*100,2) as count `
			case 12:
				sql += `FORMAT(SUM(num_credit_sucess_seven)/SUM(num_register)*100,2) as count `
			case 13:
				sql += `FORMAT(SUM(num_first_loan_apply_seven)/SUM(num_register)*100,2) as count `
			case 14:
				sql += `FORMAT(SUM(num_first_loan_sucess_seven)/SUM(num_register)*100,2) as count `
			}
		}
		sql += ` FROM HT_users_conver7_daily`
	}
	if changeCycle == 30 {
		if dataType == "用户人数" {
			switch chooseType {
			case 0:
				sql += `SUM(num_register) as count `
			case 1:
				sql += `SUM(num_cert_apply_thirty) as count `
			case 2:
				sql += `SUM(num_cert_sucess_thirty) as count `
			case 3:
				sql += `SUM(num_living_sucess_thirty) as count `
			case 4:
				sql += `SUM(num_users_base_info_thirty) as count `
			case 5:
				sql += `SUM(num_zm_auth_thirty) as count `
			case 6:
				sql += `SUM(num_link_man_thirty) as count `
			case 7:
				sql += `SUM(num_mobile_operators_mx_thirty) as count `
			case 8:
				sql += `SUM(num_bind_card_thirty) as count `
			case 9:
				sql += `SUM(num_gjj_thirty) as count `
			case 10:
				sql += `SUM(num_ali_pay_thirty) as count `
			case 11:
				sql += `SUM(num_credit_apply_thirty) as count `
			case 12:
				sql += `SUM(num_credit_sucess_thirty) as count `
			case 13:
				sql += `SUM(num_first_loan_apply_thirty) as count `
			case 14:
				sql += `SUM(num_first_loan_sucess_thirty) as count `
			}
		} else {
			switch chooseType {
			case 0:
				sql += `SUM(num_register) as count `
			case 1:
				sql += `FORMAT(SUM(num_cert_apply_thirty)/SUM(num_register)*100,2) as count `
			case 2:
				sql += `FORMAT(SUM(num_cert_sucess_thirty)/SUM(num_register)*100,2) as count `
			case 3:
				sql += `FORMAT(SUM(num_living_sucess_thirty)/SUM(num_register)*100,2) as count `
			case 4:
				sql += `FORMAT(SUM(num_users_base_info_thirty)/SUM(num_register)*100,2) as count `
			case 5:
				sql += `FORMAT(SUM(num_zm_auth_thirty)/SUM(num_register)*100,2) as count `
			case 6:
				sql += `FORMAT(SUM(num_link_man_thirty)/SUM(num_register)*100,2) as count `
			case 7:
				sql += `FORMAT(SUM(num_mobile_operators_mx_thirty)/SUM(num_register)*100,2) as count `
			case 8:
				sql += `FORMAT(SUM(num_bind_card_thirty)/SUM(num_register)*100,2) as count `
			case 9:
				sql += `FORMAT(SUM(num_gjj_thirty)/SUM(num_register)*100,2) as count `
			case 10:
				sql += `FORMAT(SUM(num_ali_pay_thirty)/SUM(num_register)*100,2) as count `
			case 11:
				sql += `FORMAT(SUM(num_credit_apply_thirty)/SUM(num_register)*100,2) as count `
			case 12:
				sql += `FORMAT(SUM(num_credit_sucess_thirty)/SUM(num_register)*100,2) as count `
			case 13:
				sql += `FORMAT(SUM(num_first_loan_apply_thirty)/SUM(num_register)*100,2) as count `
			case 14:
				sql += `FORMAT(SUM(num_first_loan_sucess_thirty)/SUM(num_register)*100,2) as count `
			}
		}
		sql += ` FROM HT_users_conver30_daily`
	}
	sql += ` WHERE 1=1 `
	if condition != "" {
		sql += condition
	}
	sql += ` GROUP BY createtime `
	_, err = o.Raw(sql, pars).QueryRows(&list)
	return
}

type ControlPassList struct {
	Createtime              time.Time
	AllCreditApply          int //申请
	AllCreditSuccess        int //通过_数量
	AllCreditRejected       int //驳回_数量
	AllCreditCloseDays      int //关闭30天_数量
	AllCreditClosePermanent int //永久关闭_数量
	AllCrediting            int //流程中_数量
}

//获取风控通过率数据
func GetControlPassList(countDimension string, start, pageSize int, isLimit bool, condition string, pars ...interface{}) (list []ControlPassList, err error) {
	o := orm.NewOrm()
	o.Using("dm_xjfq")
	sql := `SELECT createtime,`
	if countDimension == "次数" {
		sql += ` SUM(all_munc_credit_apply) as all_credit_apply,
				 SUM(all_munc_credit_sucess) as all_credit_success,
				 SUM(all_munc_credit_rejected) as all_credit_rejected,
				 SUM(all_munc_credit_close_30days) as all_credit_close_days,
				 SUM(all_munc_credit_close_permanent) as all_credit_close_permanent,
				 SUM(all_munc_crediting) as all_crediting `
	}
	if countDimension == "人数" {
		sql += ` SUM(all_mun_credit_apply) as all_credit_apply,
				 SUM(all_mun_credit_sucess) as all_credit_success,
				 SUM(all_mun_credit_rejected) as all_credit_rejected,
				 SUM(all_mun_credit_close_30days) as all_credit_close_days,
				 SUM(all_mun_credit_close_permanent) as all_credit_close_permanent,
				 SUM(all_mun_crediting) as all_crediting `
	}
	sql += ` FROM HT_credit_pass_rate_hourly WHERE 1=1 `
	if condition != "" {
		sql += condition
	}
	sql += ` GROUP BY createtime ORDER BY createtime DESC`
	if isLimit {
		sql += ` LIMIT ? ,?`
		_, err = o.Raw(sql, pars, start, pageSize).QueryRows(&list)
	} else {
		_, err = o.Raw(sql, pars).QueryRows(&list)
	}
	return
}

//获取风控通过率折线图数据
func GetControlPassRateLineData(chooseType int, countDimension, dataType, condition string, pars ...interface{}) (list []BusinessData, err error) {
	o := orm.NewOrm()
	o.Using("dm_xjfq")
	sql := `SELECT createtime,`
	if dataType == "数量" {
		if countDimension == "次数" {
			switch chooseType {
			case 0:
				sql += `SUM(all_munc_credit_apply) as count`
			case 1:
				sql += `SUM(all_munc_credit_sucess) as count`
			case 2:
				sql += `SUM(all_munc_credit_rejected) as count`
			case 3:
				sql += `SUM(all_munc_credit_close_30days) as count`
			case 4:
				sql += `SUM(all_munc_credit_close_permanent) as count`
			case 5:
				sql += `SUM(all_munc_crediting) as count`
			}
		}
		if countDimension == "人数" {
			switch chooseType {
			case 0:
				sql += `SUM(all_mun_credit_apply) as count`
			case 1:
				sql += `SUM(all_mun_credit_sucess) as count`
			case 2:
				sql += `SUM(all_mun_credit_rejected) as count`
			case 3:
				sql += `SUM(all_mun_credit_close_30days) as count`
			case 4:
				sql += `SUM(all_mun_credit_close_permanent) as count`
			case 5:
				sql += `SUM(all_mun_crediting) as count`
			}
		}
	} else {
		if countDimension == "次数" {
			switch chooseType {
			case 0:
				sql += `SUM(all_munc_credit_apply) as count`
			case 1:
				sql += `FORMAT(SUM(all_munc_credit_sucess)/SUM(all_munc_credit_apply)*100,2) as count`
			case 2:
				sql += `FORMAT(SUM(all_munc_credit_rejected)/SUM(all_munc_credit_apply)*100,2) as count`
			case 3:
				sql += `FORMAT(SUM(all_munc_credit_close_30days)/SUM(all_munc_credit_apply)*100,2) as count`
			case 4:
				sql += `FORMAT(SUM(all_munc_credit_close_permanent)/SUM(all_munc_credit_apply)*100,2) as count`
			case 5:
				sql += `FORMAT(SUM(all_munc_crediting)/SUM(all_munc_credit_apply)*100,2) as count`
			}
		}
		if countDimension == "人数" {
			switch chooseType {
			case 0:
				sql += `SUM(all_mun_credit_apply) as count`
			case 1:
				sql += `FORMAT(SUM(all_mun_credit_sucess)/SUM(all_mun_credit_apply)*100,2) as count`
			case 2:
				sql += `FORMAT(SUM(all_mun_credit_rejected)/SUM(all_mun_credit_apply)*100,2) as count`
			case 3:
				sql += `FORMAT(SUM(all_mun_credit_close_30days)/SUM(all_mun_credit_apply)*100,2) as count`
			case 4:
				sql += `FORMAT(SUM(all_mun_credit_close_permanent)/SUM(all_mun_credit_apply)*100,2) as count`
			case 5:
				sql += `FORMAT(SUM(all_mun_crediting)/SUM(all_mun_credit_apply)*100,2) as count`
			}
		}
	}
	sql += `  FROM HT_credit_pass_rate_hourly WHERE 1=1 `
	if condition != "" {
		sql += condition
	}
	sql += ` GROUP BY createtime `
	_, err = o.Raw(sql, pars).QueryRows(&list)
	return

}

//系统通过率
type SystemPassList struct {
	Createtime                 time.Time //日期
	AllCreditApply             int       //申请(次数/人数)
	SystemCreditSucess         int       //通过(次数/人数)
	SystemCreditGjjRejected    int       //公积金授权关闭30天(次数/人数)即公积金锁定30天
	SystemCreditZfbRejected    int       //支付宝授权关闭30天(次数/人数)即支付宝锁定30天
	SystemCreditMobileRejected int       //运营商授权不满足准入(次数/人数)即运营商驳回

	SystemCreditGjjzfbRejected       int //公积金和支付宝(次数/人数)
	SystemCreditGjjmobileRejected    int //公积金和运营商(次数/人数)
	SystemCreditZfbmobileRejected    int //支付宝和运营商(次数/人数)
	SystemCreditGjjzfbmobileRejected int //公积金和支付宝和运营商(次数/人数)

	SystemCreditIntoXS         int //审核(次数/人数)
	SystemCreditCloseDays      int //关闭30天(次数/人数)
	SystemCreditClosePermanent int //永久关闭(次数/人数)
	SystemCrediting            int //流程中(次数/人数)
}

//获取系统通过率数据
func GetSystemPassList(countDimension string, start, pageSize int, isLimit bool, condition string, pars ...interface{}) (list []SystemPassList, err error) {
	o := orm.NewOrm()
	o.Using("dm_xjfq")
	sql := `SELECT createtime,`
	if countDimension == "次数" {
		sql += ` SUM(all_munc_credit_apply) as all_credit_apply,
				 SUM(system_munc_credit_sucess) as system_credit_sucess,
				 SUM(system_munc_credit_gjj_rejected) as system_credit_gjj_rejected,
                 SUM(system_munc_credit_zfb_rejected) as system_credit_zfb_rejected,
				 SUM(system_munc_credit_mobile_rejected) as system_credit_mobile_rejected,
				 SUM(system_munc_credit_gjjzfb_rejected) as system_credit_gjjzfb_rejected,	
 				 SUM(system_munc_credit_gjjmobile_rejected) as system_credit_gjjmobile_rejected,
				 SUM(system_munc_credit_zfbmobile_rejected) as system_credit_zfbmobile_rejected,
				 SUM(system_munc_credit_gjjzfbmobile_rejected) as system_credit_gjjzfbmobile_rejected,
				 SUM(system_munc_credit_intoXS) as system_credit_into_x_s,
				 SUM(system_munc_credit_close_30days) as system_credit_close_days,
				 SUM(system_munc_credit_close_permanent) as system_credit_close_permanent,
				 SUM(system_munc_crediting) as system_crediting `
	}
	if countDimension == "人数" {
		sql += ` SUM(all_mun_credit_apply) as all_credit_apply,
				 SUM(system_mun_credit_sucess) as system_credit_sucess,
				 SUM(system_mun_credit_gjj_rejected) as system_credit_gjj_rejected,
                 SUM(system_mun_credit_zfb_rejected) as system_credit_zfb_rejected,
				 SUM(system_mun_credit_mobile_rejected) as system_credit_mobile_rejected,
				 SUM(system_mun_credit_gjjzfb_rejected) as system_credit_gjjzfb_rejected,	
 				 SUM(system_mun_credit_gjjmobile_rejected) as system_credit_gjjmobile_rejected,
				 SUM(system_mun_credit_zfbmobile_rejected) as system_credit_zfbmobile_rejected,
				 SUM(system_mun_credit_gjjzfbmobile_rejected) as system_credit_gjjzfbmobile_rejected,
				 SUM(system_mun_credit_intoXS) as system_credit_into_x_s,
				 SUM(system_mun_credit_close_30days) as system_credit_close_days,
				 SUM(system_mun_credit_close_permanent) as system_credit_close_permanent,
				 SUM(system_mun_crediting) as system_crediting `
	}
	sql += ` FROM HT_credit_pass_rate_hourly WHERE 1=1 `
	if condition != "" {
		sql += condition
	}
	sql += ` GROUP BY createtime ORDER BY createtime DESC`
	if isLimit {
		sql += ` LIMIT ? ,?`
		_, err = o.Raw(sql, pars, start, pageSize).QueryRows(&list)
	} else {
		_, err = o.Raw(sql, pars).QueryRows(&list)
	}
	return
}

//获取系统通过率折线图数据
func GetSystemPassRateLineData(chooseType int, countDimension, dataType, condition string, pars ...interface{}) (list []BusinessData, err error) {
	o := orm.NewOrm()
	o.Using("dm_xjfq")
	sql := `SELECT createtime,`
	if dataType == "数量" {
		if countDimension == "次数" {
			switch chooseType {
			case 0:
				sql += `SUM(all_munc_credit_apply) as count`
			case 1:
				sql += `SUM(system_munc_credit_sucess) as count`
			case 2:
				sql += `SUM(system_munc_credit_gjj_rejected) as count`
			case 3:
				sql += `SUM(system_munc_credit_zfb_rejected) as count`
			case 4:
				sql += `SUM(system_munc_credit_mobile_rejected) as count`
			case 5:
				sql += `SUM(system_munc_credit_gjjzfb_rejected) as count`
			case 6:
				sql += `SUM(system_munc_credit_gjjmobile_rejected) as count`
			case 7:
				sql += `SUM(system_munc_credit_zfbmobile_rejected) as count`
			case 8:
				sql += `SUM(system_munc_credit_gjjzfbmobile_rejected) as count`
			case 9:
				sql += `SUM(system_munc_credit_intoXS) as count`
			case 10:
				sql += `SUM(system_munc_credit_close_30days) as count`
			case 11:
				sql += `SUM(system_munc_credit_close_permanent) as count`
			case 12:
				sql += `SUM(system_munc_crediting) as count`
			}
		}
		if countDimension == "人数" {
			switch chooseType {
			case 0:
				sql += `SUM(all_mun_credit_apply) as count`
			case 1:
				sql += `SUM(system_mun_credit_sucess) as count`
			case 2:
				sql += `SUM(system_mun_credit_gjj_rejected) as count`
			case 3:
				sql += `SUM(system_mun_credit_zfb_rejected) as count`
			case 4:
				sql += `SUM(system_mun_credit_mobile_rejected) as count`
			case 5:
				sql += `SUM(system_mun_credit_gjjzfb_rejected) as count`
			case 6:
				sql += `SUM(system_mun_credit_gjjmobile_rejected) as count`
			case 7:
				sql += `SUM(system_mun_credit_zfbmobile_rejected) as count`
			case 8:
				sql += `SUM(system_mun_credit_gjjzfbmobile_rejected) as count`
			case 9:
				sql += `SUM(system_mun_credit_intoXS) as count`
			case 10:
				sql += `SUM(system_mun_credit_close_30days) as count`
			case 11:
				sql += `SUM(system_mun_credit_close_permanent) as count`
			case 12:
				sql += `SUM(system_mun_crediting) as count`
			}
		}
	} else {
		if countDimension == "次数" {
			switch chooseType {
			case 0:
				sql += `SUM(all_munc_credit_apply) as count`
			case 1:
				sql += `FORMAT(SUM(system_munc_credit_sucess)/SUM(all_munc_credit_apply)*100,2) as count`
			case 2:
				sql += `FORMAT(SUM(system_munc_credit_gjj_rejected)/SUM(all_munc_credit_apply)*100,2) as count`
			case 3:
				sql += `FORMAT(SUM(system_munc_credit_zfb_rejected)/SUM(all_munc_credit_apply)*100,2) as count`
			case 4:
				sql += `FORMAT(SUM(system_munc_credit_mobile_rejected)/SUM(all_munc_credit_apply)*100,2) as count`
			case 5:
				sql += `FORMAT(SUM(system_munc_credit_gjjzfb_rejected)/SUM(all_munc_credit_apply)*100,2) as count`
			case 6:
				sql += `FORMAT(SUM(system_munc_credit_gjjmobile_rejected)/SUM(all_munc_credit_apply)*100,2) as count`
			case 7:
				sql += `FORMAT(SUM(system_munc_credit_zfbmobile_rejected)/SUM(all_munc_credit_apply)*100,2) as count`
			case 8:
				sql += `FORMAT(SUM(system_munc_credit_gjjzfbmobile_rejected)/SUM(all_munc_credit_apply)*100,2) as count`
			case 9:
				sql += `FORMAT(SUM(system_munc_credit_intoXS)/SUM(all_munc_credit_apply)*100,2) as count`
			case 10:
				sql += `FORMAT(SUM(system_munc_credit_close_30days)/SUM(all_munc_credit_apply)*100,2) as count`
			case 11:
				sql += `FORMAT(SUM(system_munc_credit_close_permanent)/SUM(all_munc_credit_apply)*100,2) as count`
			case 12:
				sql += `FORMAT(SUM(system_munc_crediting)/SUM(all_munc_credit_apply)*100,2) as count`
			}
		}
		if countDimension == "人数" {
			switch chooseType {
			case 0:
				sql += `SUM(all_mun_credit_apply) as count`
			case 1:
				sql += `FORMAT(SUM(system_mun_credit_sucess)/SUM(all_mun_credit_apply)*100,2) as count`
			case 2:
				sql += `FORMAT(SUM(system_mun_credit_gjj_rejected)/SUM(all_mun_credit_apply)*100,2) as count`
			case 3:
				sql += `FORMAT(SUM(system_mun_credit_zfb_rejected)/SUM(all_mun_credit_apply)*100,2) as count`
			case 4:
				sql += `FORMAT(SUM(system_mun_credit_mobile_rejected)/SUM(all_mun_credit_apply)*100,2) as count`
			case 5:
				sql += `FORMAT(SUM(system_mun_credit_gjjzfb_rejected)/SUM(all_mun_credit_apply)*100,2) as count`
			case 6:
				sql += `FORMAT(SUM(system_mun_credit_gjjmobile_rejected)/SUM(all_mun_credit_apply)*100,2) as count`
			case 7:
				sql += `FORMAT(SUM(system_mun_credit_zfbmobile_rejected)/SUM(all_mun_credit_apply)*100,2) as count`
			case 8:
				sql += `FORMAT(SUM(system_mun_credit_gjjzfbmobile_rejected)/SUM(all_mun_credit_apply)*100,2) as count`
			case 9:
				sql += `FORMAT(SUM(system_mun_credit_intoXS)/SUM(all_munc_credit_apply)*100,2) as count`
			case 10:
				sql += `FORMAT(SUM(system_mun_credit_close_30days)/SUM(all_munc_credit_apply)*100,2) as count`
			case 11:
				sql += `FORMAT(SUM(system_mun_credit_close_permanent)/SUM(all_munc_credit_apply)*100,2) as count`
			case 12:
				sql += `FORMAT(SUM(system_mun_crediting)/SUM(all_munc_credit_apply)*100,2) as count`
			}
		}
	}
	sql += `  FROM HT_credit_pass_rate_hourly WHERE 1=1 `
	if condition != "" {
		sql += condition
	}
	sql += ` GROUP BY createtime `
	_, err = o.Raw(sql, pars).QueryRows(&list)
	return

}

type CreditPassList struct {
	Createtime             time.Time
	SystemCreditIntoXS     int `orm:"column(system_credit_intoXS)"`      //审核
	XSCreditSuccess        int `orm:"column(XS_credit_sucess)"`          //通过_数量
	XSCreditRejected       int `orm:"column(XS_credit_rejected)"`        //驳回_数量
	XSCreditCloseDays      int `orm:"column(XS_credit_close_days)"`      //关闭30天_数量
	XSCreditClosePermanent int `orm:"column(XS_credit_close_permanent)"` //永久关闭_数量
	XSCrediting            int `orm:"column(XS_crediting)"`              //流程中_数量
}

//获取信审通过率数据
func GetCreditPassList(countDimension string, start, pageSize int, isLimit bool, condition string, pars ...interface{}) (list []CreditPassList, err error) {
	o := orm.NewOrm()
	o.Using("dm_xjfq")
	sql := `SELECT createtime,`
	if countDimension == "次数" {
		sql += ` SUM(system_munc_credit_intoXS) AS system_credit_intoXS,
				 SUM(XS_munc_credit_sucess) AS XS_credit_sucess,
				 SUM(XS_munc_credit_rejected) AS XS_credit_rejected,
				 SUM(XS_munc_credit_close_30days) AS XS_credit_close_days,
				 SUM(XS_munc_credit_close_permanent) AS XS_credit_close_permanent,
				 SUM(XS_munc_crediting) AS XS_crediting `
	}
	if countDimension == "人数" {
		sql += ` SUM(system_mun_credit_intoXS) AS system_credit_intoXS,
				 SUM(XS_mun_credit_sucess) AS XS_credit_sucess,
				 SUM(XS_mun_credit_rejected) AS XS_credit_rejected,
				 SUM(XS_mun_credit_close_30days) AS XS_credit_close_days,
				 SUM(XS_mun_credit_close_permanent) AS XS_credit_close_permanent,
				 SUM(XS_mun_crediting) AS XS_crediting `
	}
	sql += ` FROM HT_credit_pass_rate_hourly WHERE 1=1 `
	if condition != "" {
		sql += condition
	}
	sql += ` GROUP BY createtime ORDER BY createtime DESC`
	if isLimit {
		sql += ` LIMIT ? ,?`
		_, err = o.Raw(sql, pars, start, pageSize).QueryRows(&list)
	} else {
		_, err = o.Raw(sql, pars).QueryRows(&list)
	}
	return
}

//获取系统通过率折线图数据
func GetCreditPassRateLineData(chooseType int, countDimension, dataType, condition string, pars ...interface{}) (list []BusinessData, err error) {
	o := orm.NewOrm()
	o.Using("dm_xjfq")
	sql := `SELECT createtime,`
	if dataType == "数量" {
		if countDimension == "次数" {
			switch chooseType {
			case 0:
				sql += `SUM(system_munc_credit_intoXS) as count`
			case 1:
				sql += `SUM(XS_munc_credit_sucess) as count`
			case 2:
				sql += `SUM(XS_munc_credit_rejected) as count`
			case 3:
				sql += `SUM(XS_munc_credit_close_30days) as count`
			case 4:
				sql += `SUM(XS_munc_credit_close_permanent) as count`
			case 5:
				sql += `SUM(XS_munc_crediting) as count`
			}
		}
		if countDimension == "人数" {
			switch chooseType {
			case 0:
				sql += `SUM(system_mun_credit_intoXS) as count`
			case 1:
				sql += `SUM(XS_mun_credit_sucess) as count`
			case 2:
				sql += `SUM(XS_mun_credit_rejected) as count`
			case 3:
				sql += `SUM(XS_mun_credit_close_30days) as count`
			case 4:
				sql += `SUM(XS_mun_credit_close_permanent) as count`
			case 5:
				sql += `SUM(XS_mun_crediting) as count`
			}
		}
	} else {
		if countDimension == "次数" {
			switch chooseType {
			case 0:
				sql += `SUM(system_munc_credit_intoXS) as count`
			case 1:
				sql += `FORMAT(SUM(XS_munc_credit_sucess)/SUM(system_munc_credit_intoXS)*100,2) as count`
			case 2:
				sql += `FORMAT(SUM(XS_munc_credit_rejected)/SUM(system_munc_credit_intoXS)*100,2) as count`
			case 3:
				sql += `FORMAT(SUM(XS_munc_credit_close_30days)/SUM(system_munc_credit_intoXS)*100,2) as count`
			case 4:
				sql += `FORMAT(SUM(XS_munc_credit_close_permanent)/SUM(system_munc_credit_intoXS)*100,2) as count`
			case 5:
				sql += `FORMAT(SUM(XS_munc_crediting)/SUM(system_munc_credit_intoXS)*100,2) as count`
			}
		}
		if countDimension == "人数" {
			switch chooseType {
			case 0:
				sql += `SUM(system_mun_credit_intoXS) as count`
			case 1:
				sql += `FORMAT(SUM(XS_mun_credit_sucess)/SUM(system_mun_credit_intoXS)*100,2) as count`
			case 2:
				sql += `FORMAT(SUM(XS_mun_credit_rejected)/SUM(system_mun_credit_intoXS)*100,2) as count`
			case 3:
				sql += `FORMAT(SUM(XS_mun_credit_close_30days)/SUM(system_mun_credit_intoXS)*100,2) as count`
			case 4:
				sql += `FORMAT(SUM(XS_mun_credit_close_permanent)/SUM(system_mun_credit_intoXS)*100,2) as count`
			case 5:
				sql += `FORMAT(SUM(XS_mun_crediting)/SUM(system_mun_credit_intoXS)*100,2) as count`
			}
		}
	}
	sql += `  FROM HT_credit_pass_rate_hourly WHERE 1=1 `
	if condition != "" {
		sql += condition
	}
	sql += ` GROUP BY createtime `
	_, err = o.Raw(sql, pars).QueryRows(&list)
	return

}

//获取授信通过率数据总数(风控通过率、系统通过率、信审通过率)
func CreditPassListCount(condition string, pars ...interface{}) (count int, err error) {
	o := orm.NewOrm()
	o.Using("dm_xjfq")
	sql := `SELECT COUNT(DISTINCT createtime) FROM HT_credit_pass_rate_hourly WHERE 1=1`
	if condition != "" {
		sql += condition
	}
	err = o.Raw(sql, pars).QueryRow(&count)
	return
}

package models

import (
	"time"
	"zcm_tools/orm"
)

//借款记录
type LoanRecord struct {
	Id                  int
	Uid                 int       //用户ID
	LoanDate            time.Time //申请时间
	ContractCode        string    //合同编号
	Money               float64   //借款金额
	LoanTax             float64   //月利息
	TotalLoanTax        float64   //总利息 = 月利息*期数
	AuthFee             float64   //信审数据费
	RealMoney           float64   //到账金额
	LoanTermCount       int       //借款期限
	ReplyTime           time.Time //批复时间
	EndDate             time.Time //应还时间
	FinishDate          time.Time //结清日期
	State               string    //借款状态
	OperatorName        string    //信审操作人
	AuditType           int       //1-机器审核 2-人工审核
	OrderState          string    //订单审核状态
	LoanAgreementUrl    string    //借款协议
	ServiceAgreementUrl string    //服务协议
	Mark                string    //备注
	RepaymentScheduleId int       //还款计划id
}
type LoanDisplay struct {
	Uid            int
	Account        string
	Verifyrealname string
	DtOrder        time.Time
	DtOrderStr     string
	Money          float64
	LoanTermCount  int
	Displayname    string
	OrderNumber    string
	State          string
}

//查询用户借款记录
func QueryUsersLoanRecords(uid int) (list []LoanRecord, err error) {
	sql := `SELECT
				l.id,
				l.uid,
				l.loan_date,
				l.money,
				l.auth_fee,
				l.real_money,
				l.loan_term_count,
				l.reply_time,
				l.end_date,
				l.finish_date,
				l.mark,
				l.state
			FROM loan AS l
			WHERE l.uid = ?
			ORDER BY l.create_time DESC`
	o := orm.NewOrm()
	_, err = o.Raw(sql, uid).QueryRows(&list)
	return
}

//查询用户借款记录(带分页)
func QueryUsersManageLoanRecords(uid, start, pageSize int) (list []LoanRecord, err error) {
	sql := `SELECT
				l.id,
				l.uid,
				l.loan_date,
				l.contract_code,
				l.money,
				l.loan_tax,
				l.loan_tax AS total_loan_tax,
				l.auth_fee,
				l.real_money,
				l.loan_term_count,
				l.reply_time,
				l.end_date,
				l.finish_date,
				l.state,
				l.credit_operator AS operator_name,
				l.audit_type,
				l.order_state,
				l.mark,
				l.repayment_schedule_id
			FROM loan AS l
			WHERE l.uid = ?
			ORDER BY l.loan_date DESC
			LIMIT ?,?`
	o := orm.NewOrm()
	_, err = o.Raw(sql, uid, start, pageSize).QueryRows(&list)
	return
}

//查询用户借款记录总数
func QueryUsersManageLoanRecordCount(uid int) (count int, err error) {
	sql := `SELECT
				COUNT(1)
			FROM loan AS l
			WHERE l.uid = ?`
	o := orm.NewOrm()
	err = o.Raw(sql, uid).QueryRow(&count)
	return
}

//查询用户借款几笔
func QueryLoanRecordOverDueCount(uid int) (count int, err error) {
	sql := `SELECT
			COUNT(1)
			FROM
				(
					SELECT
						COUNT(1) AS hk
					FROM
						repayment_schedule AS b
					WHERE
						b.state = "BACKING"
					AND b.overdue_days > 0
					AND uid = ?
					GROUP BY
						loan_id
				) t`
	o := orm.NewOrm()
	err = o.Raw(sql, uid).QueryRow(&count)
	return
}

//还款计划
type PaymentSchedules struct {
	Id                             int
	Money                          float64 //借款本金
	TermNo                         int     //期数
	Uid                            int     //用户ID
	LoanReturnDate                 string  //应还日期
	CapitalAmount                  float64 //应还本金
	TaxAmount                      float64 //应还利息
	OverdueBreachOfAmount          float64 //应还违约金
	OverdueMoneyAmount             float64 //应还滞纳金
	AheadOfTimeClearedAmount       float64 //应收提前结清违约金
	DataServiceFee                 float64 //应收信审数据费
	ReturnCapitalAmount            float64 //已还本金
	ReturnTaxAmount                float64 //已还利息
	ReturnOverdueBreachOfAmount    float64 //已还违约金
	ReturnOverdueMoneyAmount       float64 //已还滞纳金
	ReturnAheadOfTimeClearedAmount float64 //实还提前结清违约金
	ReturnDataServiceFee           float64 //实还信审数据费
	OverdueDays                    int     //逾期天数
	State                          string  //状态
	RemainMoneyChargeUpAmount      float64 //挂账金额
}

//查询用户还款计划
func QueryPaymentSchedules(loanId int) (list []PaymentSchedules, err error) {
	sql := `SELECT
				id,
				uid,
				term_no,
				loan_return_date,
				capital_amount,
				tax_amount,
				overdue_breach_of_amount,
				overdue_money_amount,
				ahead_of_time_cleared_amount,
				data_service_fee,
				return_capital_amount,
				return_tax_amount,
				return_overdue_breach_of_amount,
				return_overdue_money_amount,
				overdue_days,
				return_ahead_of_time_cleared_amount,
				return_data_service_fee,
				state,
				remain_money_charge_up_amount
           FROM repayment_schedule
		   WHERE loan_id = ?
		   ORDER BY loan_return_date`
	o := orm.NewOrm()
	_, err = o.Raw(sql, loanId).QueryRows(&list)
	return
}

//获取用户成功借款笔数
func QueryUserLoanSuccessCount(uid int) (count int, err error) {
	sql := `SELECT loan_success_count FROM users_loan_detail WHERE uid=?`
	err = orm.NewOrm().Raw(sql, uid).QueryRow(&count)
	return
}

//获取用户当前未结清的借款
func QueryUserBackingLoan(uid int) (rs []PaymentSchedules, err error) {
	o := orm.NewOrm()
	sql := `SELECT
				l.money,
				rs.capital_amount,
				rs.tax_amount,
				rs.overdue_money_amount,
				rs.overdue_breach_of_amount,
				rs.remain_money_charge_up_amount,
				rs.term_no,
				rs.loan_return_date,
				rs.overdue_days,
				rs.state
			FROM loan AS l
			INNER JOIN repayment_schedule AS rs
			ON l.id = rs.loan_id
			WHERE l.state = 'BACKING'
			AND l.uid = ?
			ORDER BY rs.loan_id DESC,rs.term_no ASC`
	_, err = o.Raw(sql, uid).QueryRows(&rs)
	return
}

//获取借款列表
func LoanList(start, pageSize int, isLimit bool, condition string, paras ...interface{}) (rs []*LoanDisplay, err error) {
	o := orm.NewOrm()
	sql := `SELECT
				l.uid,
				l.account,
				um.verifyrealname,
				tr.dt_order,
				l.money,
				l.loan_term_count,
				su.displayname,
                tr.order_number,
                tr.state
			FROM loan AS l 
			INNER JOIN users_metadata um ON l.uid = um.uid
 			LEFT JOIN trade_record tr on l.id=tr.loan_id 
			LEFT JOIN sys_user su on l.operator_id = su.id
			WHERE tr.dt_order is NOT NULL `
	if condition != "" {
		sql += condition
	}
	sql += ` ORDER BY tr.dt_order DESC  `
	if isLimit {
		sql += ` LIMIT ?,? `
		_, err = o.Raw(sql, paras, start, pageSize).QueryRows(&rs)

	} else {
		_, err = o.Raw(sql, paras).QueryRows(&rs)
	}
	return
}
func QueryListCount(condition string, paras ...interface{}) (count int, err error) {
	o := orm.NewOrm()
	sql := ` SELECT
				count(1)
			FROM loan l 
			INNER JOIN users_metadata um ON l.uid = um.uid
 			LEFT JOIN trade_record tr on l.id=tr.loan_id 
			LEFT JOIN sys_user su on l.operator_id = su.id
			WHERE tr.dt_order is NOT NULL`
	if condition != "" {
		sql += condition
	}
	err = o.Raw(sql, paras).QueryRow(&count)
	return
}
func GetCreditUsersList() (names []string, err error) {
	o := orm.NewOrm()
	sql := `SELECT
				displayname from sys_user
			WHERE role_id in (53,54,2,3)`
	_, err = o.Raw(sql).QueryRows(&names)
	return
}

func VerfyLoan(account string) (repayment_schedule_id int) {
	o := orm.NewOrm()
	sql := `SELECT
				a.repayment_schedule_id
			FROM
				loan AS a,
				repayment_schedule AS b
			WHERE
				a.account = ?
			AND a.state = "BACKING"
			AND b.id = a.repayment_schedule_id
			ORDER BY b.loan_return_date ASC,b.capital_amount DESC 
			LIMIT 1`
	o.Raw(sql, account).QueryRow(&repayment_schedule_id)
	return
}

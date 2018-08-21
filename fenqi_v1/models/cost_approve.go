package models

import (
	"zcm_tools/orm"
	// "jqb_v1/services"
	"fmt"
	//"jqb_v1/utils"
	"zcm_tools/email"
)

//费用减免审批
type CostReliefApprove struct {
	Id              int     `json:"id"                    description:"编号"`
	Uid             int     `json:"uid"                   description:"用户编号"`
	Phone           string  `json:"phone"                 description:"手机号"`
	State           string  `json:"state"                 description:"审批状态(NONE 待审核,REFUSE 拒绝,AGREE 同意,DONE 已完成)"`
	RepSchId        int     `json:"repSchId"              description:"还款计划id"`
	LoanId          int     `json:"loanId"                description:"借款id"`
	Money           float64 `json:"money"                 description:"减免金额"`
	Reason          string  `json:"reason"                description:"减免原因"`
	ApproveResult   string  `json:"approveResult"         description:"审批结果"`
	RequestSysId    int     `json:"requestSysId"          description:"提交申请sys_userid"`
	RequestSysName  string  `json:"requestSysName"        description:"审批结果"`
	ApproveSysId    int     `json:"approveSysId"          description:"审批人sys_userid"`
	DisposeSysId    int     `json:"disposeSysId"          description:"处理人sys_userid"`
	DisposeResult   string  `json:"disposeResult"         description:"处理结果"`
	CreateTime      string  `json:"createTime"            description:"创建时间"`
	Displayname     string  `json:"displayname"           description:"审批人姓名"`
	ImagesUrl       string  `json:"imagesUrl"           	 description:"审批图片url"`
	ApproveMoney    float64 `json:"approve_money"         description:"审批通过所填金额"`
	IsApprove       int     `json:"is_approve"            description:"是否审批通过：1-是"`
	SubmitMoney     float64 `json:"submit_money"          description:"提交上级所填的金额"`
	IsSubmit        int     `json:"is_submit"           	 description:"是否提交上级：1-是"`
	CollectionPhase string  // 催收阶段
	Company         string  //所属公司
}

//借款所需参数
type LoanResult struct {
	Id               int     `json:"id"                    description:"编号"`
	LoanMoney        float64 `json:"loanMoney"             description:"借款金额"`
	LoanTime         string  `json:"loanTime"              description:"借款时间"`
	RepayTime        string  `json:"repayTime"             description:"到期时间"`
	OverDay          int     `json:"overDay"               description:"逾期天数"`
	OverMoney        float64 `json:"overMoney"             description:"逾期费用"`
	SumRepayMoney    float64 `json:"sumRepayMoney"         description:"应还总金额"`
	RepayMoney       float64 `json:"repayMoney"            description:"已还金额"`
	WaitMoney        float64 `json:"waitMoney"             description:"待还金额"`
	ProductId        int     `json:"productId"             description:"产品 id"`
	ShouldAlsoAmount float64 `json:"shouldAlsoAmount"      description:"分期应还金额"`
	TermNo           int     `json:"termNo"                description:"期号"`
	TaxAmount        float64 `json:"taxAmount"             description:"利息"`
	CapitalAmount    float64 `json:"capitalAmount"         description:"本金"`
	StageFee         float64 `json:"stageFee"              description:"分期费"`
}

//审批处理记录
type CostReliefApproveRecord struct {
	Id             int    `json:"id"                    description:"编号"`
	CraId          int    `json:"craId"                 description:"费用减免id"`
	OldApproveId   int    `json:"oldApproveId"          description:"原审批人id"`
	NewApproveId   int    `json:"newApproveId"          description:"新审批人id"`
	NewApproveName string `json:"newApproveName"        description:"新审批人姓名"`
	Source         string `json:"source"                description:"来源于审批(APPROVE)还是处理(DISPOSE)"`
	Opinion        string `json:"opinion"               description:"处理审批意见"`
	CreateTime     string `json:"createTime"            description:"创建时间"`
	Displayname    string `json:"displayname"           description:"审批人姓名"`
	State          string `json:"state"           		description:"当前状态"`
}

func FindByIdLoan(loanId int) (loan LoanResult) {
	o := orm.NewOrm()
	sql := `select
					l.create_time as loan_time,l.money as loan_money,rs.loan_return_date as repay_time,rs.overdue_days as over_day,
					sum(ifnull(rs.tax_amount,0)+ifnull(rs.overdue_money_amount,0)+ifnull(rs.overdue_breach_of_amount,0)
					+ifnull(rs.ahead_of_time_cleared_amount,0))as over_money,
					sum(ifnull(rs.return_tax_amount,0)+ifnull(rs.return_overdue_money_amount,0)+ifnull(rs.return_overdue_breach_of_amount,0)
					+ifnull(rs.return_ahead_of_time_cleared_amount,0)+ifnull(rs.return_stage_fee,0)+ifnull(rs.return_capital_amount,0))as repay_money,l.product_id,rs.tax_amount
					from loan l inner join repayment_schedule rs  on  l.id=rs.loan_id where  l.id=?`
	o.Raw(sql, loanId).QueryRow(&loan)
	return
}

func FindByIdLoanFenQi(loanId int) (list []LoanResult, err error) {
	o := orm.NewOrm()
	sql := `SELECT
			rs.loan_start_date AS loan_time,
			l.money AS loan_money,
			rs.loan_return_date AS repay_time,
			rs.overdue_days AS over_day,
			SUM(IFNULL(rs.tax_amount,0)+ IFNULL(rs.overdue_money_amount,0)+ IFNULL(rs.overdue_breach_of_amount,0)+ IFNULL(rs.ahead_of_time_cleared_amount,0)+ IFNULL(rs.data_service_fee,0)) AS over_money, SUM(IFNULL(rs.tax_amount,0)+ IFNULL(rs.overdue_money_amount,0)+ IFNULL(rs.overdue_breach_of_amount,0)+ IFNULL(rs.ahead_of_time_cleared_amount,0)+ IFNULL(rs.data_service_fee,0)+ IFNULL(rs.capital_amount,0)) AS should_also_amount,
			SUM(IFNULL(rs.return_tax_amount,0)+ IFNULL(rs.return_overdue_money_amount,0)+ IFNULL(rs.return_overdue_breach_of_amount,0)+ IFNULL(rs.return_ahead_of_time_cleared_amount,0)+ IFNULL(rs.return_data_service_fee,0)+ IFNULL(rs.return_capital_amount,0)) AS repay_money,rs.term_no,rs.tax_amount,rs.capital_amount,rs.data_service_fee
			FROM loan l
			INNER JOIN repayment_schedule rs ON l.id=rs.loan_id
			WHERE l.id=? group by rs.term_no`
	_, err = o.Raw(sql, loanId).QueryRows(&list)
	return
}

//提交费用减免审批
func AddCostReliefApprove(uid, repSchId, loanId, requestSysId, approveSysId int,
	requestSysName, phone, reason, imagesUrl, collectionPhase, company string, money float64) (err error) {
	o := orm.NewOrm()
	_, err = o.Raw(`insert cost_relief_approve (uid,rep_sch_id,loan_id,request_sys_id,approve_sys_id,request_sys_name,phone,reason,money,state,create_time,images_url,collection_phase,company)
		values(?,?,?,?,?,?,?,?,?,'NONE',now(),?,?,?)`, uid, repSchId, loanId, requestSysId, approveSysId, requestSysName, phone, reason, money, imagesUrl, collectionPhase, company).Exec()
	if err == nil {
		sysUser, err := FindByIdSysUser(approveSysId)
		if err == nil && sysUser.Email != "" {
			email.Send("费用减免通知", "您有一条新的费用减免申请等待审批，请及时处理", sysUser.Email, "huawuyou")
		}
	}
	return
}

//查询该笔借款状态
func GetLoanState(loanId int) (state string, err error) {
	o := orm.NewOrm()
	sql := "SELECT state FROM loan WHERE id = ?"
	err = o.Raw(sql, loanId).QueryRow(&state)
	return
}

//查询未完结数量
func QueryCostUnfinishCount(uid, loanId int) (count int, err error) {
	o := orm.NewOrm()
	sql := "select count(1) from cost_relief_approve where uid=? and loan_id=? and  state in ('NONE','AGREE')"
	err = o.Raw(sql, uid, loanId).QueryRow(&count)
	return
}

//查询所有费用减免
func QueryCostReliefApprove(pageNo, pageSize int, where string, params ...string) (list []CostReliefApprove, err error) {
	o := orm.NewOrm()
	sql := `SELECT
			crf.id,
			crf.uid,
			crf.phone,
			crf.state,
			crf.rep_sch_id,
			crf.loan_id,
			crf.money,
			crf.reason,
			crf.approve_result,
			crf.request_sys_id,
			crf.request_sys_name,
			crf.approve_sys_id,
			crf.dispose_sys_id,
			crf.create_time,
			s.displayname,
			crf.dispose_result,
			crf.approve_money,
			crf.is_approve,
			crf.submit_money,
			crf.is_submit,
			crf.collection_phase,
			crf.company
		FROM
			cost_relief_approve crf
		LEFT JOIN sys_user s ON crf.approve_sys_id = s.id
		WHERE	1 = 1 `
	sql += where
	sql += ` ORDER BY crf.id DESC limit ?,?`

	_, err = o.Raw(sql, params, pageNo, pageSize).QueryRows(&list)
	return
}

//查询所有费用减免(不分页)
func QueryCostReliefApproveNoPage(where string, params ...string) (list []CostReliefApprove, err error) {
	o := orm.NewOrm()
	sql := `SELECT
				crf.id,
				crf.request_sys_name,
				crf.create_time,
				crf.phone,
				crf.money,
				crf.state,
				crf.approve_money,
				crf.is_approve,
				crf.submit_money,
				crf.is_submit,
				s.displayname
			FROM cost_relief_approve crf
			LEFT JOIN sys_user s ON crf.approve_sys_id=s.id
			WHERE 1=1  `
	if where != "" {
		sql += where
	}
	sql += ` ORDER BY crf.id DESC`
	_, err = o.Raw(sql, params).QueryRows(&list)
	return
}

//查询条数
func QueryCount(where string, pars ...string) int {
	o := orm.NewOrm()
	sql := `select count(1)
	        from cost_relief_approve crf left join sys_user s on crf.approve_sys_id=s.id where 1=1 `
	sql += where
	var count int
	o.Raw(sql, pars).QueryRow(&count)
	return count
}

//查询单个审批记录
func FindByApprove(id int) (approve CostReliefApprove, err error) {
	if id < 1 {
		return
	}
	o := orm.NewOrm()
	sql := `select id,uid,phone,state,rep_sch_id,loan_id,money,reason,approve_result,dispose_result,
	               request_sys_id,request_sys_name,approve_sys_id,dispose_sys_id,create_time,images_url,
				   approve_money,is_approve,submit_money,is_submit
		        from cost_relief_approve where id=?`
	err = o.Raw(sql, id).QueryRow(&approve)
	return
}

//查询管理员姓名
func FindBySysName(sysId int) (approve CostReliefApprove) {
	o := orm.NewOrm()
	o.Raw("select displayname from sys_user where id=?", sysId).QueryRow(&approve)
	return
}

//更新费用减免表
func UpdateCostReliefApprove(id int, money float64, isTag string) (err error) {
	sql := ``
	if isTag == "Approve" { //审批通过
		sql = `UPDATE cost_relief_approve SET approve_money = ?,is_approve = 1 WHERE id = ?`
	} else if isTag == "Submit" { //提交上级
		sql = `UPDATE cost_relief_approve SET submit_money = ?,is_submit = 1 WHERE id = ?`
	}
	_, err = orm.NewOrm().Raw(sql, money, id).Exec()
	return
}

//减免审批
func ApproveRelie(id, sysId int, state, reason string) (has bool, msg string) {
	o := orm.NewOrm()
	sql := ``
	var err error
	if state == "DONE" || state == "REFUSEDONE" {
		sql += `update cost_relief_approve set state=?,dispose_result=?,dispose_sys_id=? where id=? `
		_, err = o.Raw(sql, state, reason, sysId, id).Exec()
		if state == "REFUSEDONE" && err == nil {
			AddCostReliefApproveRecord(0, id, sysId, reason, "DISPOSE", state)
		}
	} else {
		sql += `update cost_relief_approve set state=?,approve_result=? where id=? `
		_, err = o.Raw(sql, state, reason, id).Exec()
		if err == nil {
			AddCostReliefApproveRecord(0, id, sysId, reason, "APPROVE", state)
		}
	}
	if err != nil {
		msg = err.Error()
		return
	}
	approve, err := FindByApprove(id)
	requestSysId := approve.RequestSysId
	if err == nil {
		if state == "AGREE" {
			// sysUser, err := FindByIdSysUser(98)
			var list []SysUser
			var str string
			sql := `select dispose_id from cost_config where id =1`
			err := o.Raw(sql).QueryRow(&str)
			if err != nil {
				fmt.Println(err)
			}
			sysUserSql := `select id,email from sys_user where id in(` + str + `)`
			num, err := o.Raw(sysUserSql).QueryRows(&list)
			if err != nil {
				fmt.Println(err)
			}
			if num > 0 {

				// o.Raw("select email from sys_user where id in(97,98,114)").QueryRows(&list)
				// if len(list) > 0 {
				for _, v := range list {
					if v.Id == 20 {
						continue
					}
					if v.Email != "" {
						email.Send("费用减免通知", "您有一条新的费用减免申请等待处理，请及时处理", v.Email, "huawuyou")
					}
				}
			}

		} else {
			sysUser, err := FindByIdSysUser(requestSysId)
			if err == nil && sysUser.Email != "" {
				if state == "REFUSE" || state == "REFUSEDONE" {
					email.Send("费用减免通知", "您发起的手机号为"+approve.Phone+"的费用减免申请被退回，请及时查看并处理", sysUser.Email, "huawuyou")
				}
				if state == "DONE" {
					email.Send("费用减免通知", "您发起的手机号为"+approve.Phone+"的费用减免申请已处理完成", sysUser.Email, "huawuyou")
				}
			}
		}
	}
	has = true
	msg = "success"
	return
}

func AddCostReliefApproveRecord(id, craId, newApproveId int, opinion, source, state string) (err error) {
	var oldApproveId int
	oldApproveId = newApproveId
	o := orm.NewOrm()
	if id == 1 {
		approve, _ := FindByApprove(craId)
		oldApproveId = approve.ApproveSysId
		o.Raw("update cost_relief_approve set approve_sys_id=? where id=?", newApproveId, craId).Exec()
		approveSysId := approve.ApproveSysId
		if approveSysId > 0 {
			sysUser, err := FindByIdSysUser(approveSysId)
			if err == nil && sysUser.Email != "" {
				email.Send("费用减免通知", "您有一条新的费用减免申请等待审批，请及时处理", sysUser.Email, "huawuyou")
			}
		}
	}
	_, err = o.Raw("insert into cost_relief_approve_record(cra_id,old_approve_id,new_approve_id,opinion,source,create_time,state)values(?,?,?,?,?,now(),?)",
		craId, oldApproveId, newApproveId, opinion, source, state).Exec()
	return
}

func FindByCraIdApproveRecord(craId int) (list []CostReliefApproveRecord) {
	o := orm.NewOrm()
	o.Raw(`select c.id,c.cra_id,c.old_approve_id,c.new_approve_id,c.source,c.opinion,c.create_time,c.state,s.displayname
   	         from cost_relief_approve_record c left join sys_user s on c.old_approve_id=s.id where c.cra_id=?`,
		craId).QueryRows(&list)
	return
}

func FindByCraIdRecordCount(craId int) (count int) {
	o := orm.NewOrm()
	o.Raw("select count(1) from cost_relief_approve_record where cra_id=? and state='NONE'", craId).QueryRow(&count)
	return
}

func GetCollectionMoneyByloanid2(loanid string) (money float64, err error) {
	sql := `SELECT SUM(capital_amount+tax_amount+overdue_breach_of_amount+overdue_money_amount+data_service_fee-remain_money_charge_up_amount)
			FROM repayment_schedule
			WHERE loan_id=? AND state='backing'
			GROUP BY loan_id `
	o := orm.NewOrm()
	o.Using("read")
	err = o.Raw(sql, loanid).QueryRow(&money)
	return
}

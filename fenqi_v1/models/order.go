package models

import (
	"fenqi_v1/utils"
	"time"
	"zcm_tools/orm"

	"github.com/astaxie/beego"
)

type OrderCredit struct {
	Id                  int
	Money               float64 //额度
	LoanTermCount       int     //期数
	CardNumber          string  //银行卡
	BankName            string
	BankMobile          string    //银行卡预留手机号
	Uid                 int       //用户Id
	VerifyRealName      string    //姓名
	Account             string    //手机号
	IdCard              string    //身份证
	OrderState          string    //订单状态
	InqueueType         int       //入列状态
	CreateTime          time.Time //提交时间
	InqueueTime         time.Time //入列时间
	HandleTime          time.Time //处理时间
	CreditOperator      string    // 处理人
	Mark                string    //备注
	Content             string    //操作内容
	State               string    //打款状态
	Source              string    //渠道
	RepaymentScheduleId int       `orm:"column(repayment_schedule_id)"` //还款计划id
	Atime               time.Time `orm:"column(allocation_time)"`       //分配时间
	IsQueue             int       `orm:"column(is_queue)"`              //是否入队列
	CurrentLoanState    int       `orm:"column(current_loan_state)"`
}

//获取借款退回列表总数
func GetOutQueueCount(condition string, paras ...interface{}) (count int, err error) {
	o := orm.NewOrm()
	sql := ` SELECT
				COUNT(1)
			FROM loan AS l
			INNER JOIN users_metadata AS um
			ON l.uid = um.uid
			WHERE 1=1
			AND l.is_now = 1
			AND((l.inqueue_time IS NOT NULL
			AND l.inqueue_time > NOW())
			OR (l.inqueue_type = 0))`
	if condition != "" {
		sql += condition
	}
	sql += ` AND l.order_state = "OUTQUEUE" `
	err = o.Raw(sql, paras).QueryRow(&count)
	return
}

//获取借款退回列表
func GetOutQueueList(start, pageSize int, condition string, paras ...interface{}) (list []OrderCredit, err error) {
	o := orm.NewOrm()
	sql := `SELECT
				l.id,
				l.money,
				l.loan_term_count,
				l.order_state,
				l.handle_time,
				l.credit_operator,
				l.inqueue_time,
				um.uid,
				um.verifyrealname as verify_real_name,
				um.account,
				l.mark,
				l.create_time
			FROM loan AS l
			LEFT JOIN users_metadata AS um
			ON l.uid = um.uid
			WHERE 1=1
			AND l.is_now = 1
			AND((l.inqueue_time IS NOT NULL
			AND l.inqueue_time > NOW())
			OR (l.inqueue_type = 0))`
	if condition != "" {
		sql += condition
	}
	sql += ` AND l.order_state = "OUTQUEUE" ORDER BY l.create_time DESC LIMIT ?,?`
	_, err = o.Raw(sql, paras, start, pageSize).QueryRows(&list)
	return
}

//获取订单管理列表
func GetOrderCreditList(start, pageSize int, condition string, paras ...interface{}) (list []OrderCredit, err error) {
	o := orm.NewOrm()
	sql := `SELECT
				l.id,
				l.money,
				l.loan_term_count,
				l.order_state,
				l.handle_time,
				l.credit_operator,
				l.allocation_time,
				l.create_time,
				l.inqueue_time,
				l.current_loan_state,
				um.uid,
				um.verifyrealname as verify_real_name,
				um.account
			FROM loan AS l
			LEFT JOIN users_metadata AS um
			ON l.uid = um.uid
			WHERE 1=1  AND is_now = 1`
	if condition != "" {
		sql += condition
	}
	sql += ` ORDER BY l.create_time DESC LIMIT ?,?`
	_, err = o.Raw(sql, paras, start, pageSize).QueryRows(&list)
	return
}

//获取订单管理列表总数
func GetOrderCreditCount(condition string, paras ...interface{}) (count int, err error) {
	o := orm.NewOrm()
	sql := ` SELECT
				COUNT(1)
			FROM loan AS l
			LEFT JOIN users_metadata AS um
			ON l.uid = um.uid
			WHERE 1=1 AND is_now = 1`
	if condition != "" {
		sql += condition
	}
	err = o.Raw(sql, paras).QueryRow(&count)
	return
}

//==============================================================================

//获取预约时间的订单
func GetOrdering() (orderCredit OrderCredit, err error) {
	sql := `SELECT l.id,l.uid,l.inqueue_type,l.order_state,l.handle_time,l.credit_operator
			FROM loan AS l	
			WHERE
				l.inqueue_time IS NOT NULL
			AND 
				l.inqueue_time <= NOW()
			AND 
				l.inqueue_type >  0
			AND
				l.order_state = 'OUTQUEUE'
			AND 
				l.is_queue = 1
			ORDER BY l.inqueue_time
			LIMIT 1`
	err = orm.NewOrm().Raw(sql).QueryRow(&orderCredit)
	return
}

//获取排队中的订单
func GetQueueing() (orderCredit OrderCredit, err error) {
	sql := `SELECT l.id,l.uid,l.inqueue_type,l.order_state,l.handle_time,l.credit_operator
			FROM loan AS l
			WHERE
				l.order_state = 'QUEUEING'
			AND l.is_queue = 1
			ORDER BY l.queue_time 
			LIMIT 1`
	err = orm.NewOrm().Raw(sql).QueryRow(&orderCredit)
	return
}

//排队 (更改queue_time)
func UpdateQueueTime(loanId int, inqueueTime time.Time) error {
	sql := `UPDATE loan SET order_state = "QUEUEING" ,queue_time= ?  WHERE id = ? `
	_, err := orm.NewOrm().Raw(sql, inqueueTime, loanId).Exec()

	return err
}

//超时清缓存并入队列
func UpdateQueueing(loanId int) error {
	sql := `UPDATE loan SET order_state = "QUEUEING",queue_time = NOW() WHERE id = ? `
	_, err := orm.NewOrm().Raw(sql, loanId).Exec()
	return err
}

//选择出列订单更新入列状态
func UpadateInQueue(loading Loading) (err error) {
	o := orm.NewOrm()
	err = o.Begin()
	if err != nil {
		return err
	}
	sql := `UPDATE loan SET inqueue_type = ? ,`
	if loading.InqueueTime == "" {
		sql += `inqueue_time= NULL `
	} else {
		sql += `inqueue_time= ? `
	}
	sql += ` WHERE id = ?`
	p, err := o.Raw(sql).Prepare()
	defer p.Close()
	for _, v := range loading.LoanIds {
		if loading.InqueueTime == "" {
			_, err = p.Exec(loading.InqueueType, v)
		} else {
			_, err = p.Exec(loading.InqueueType, loading.InqueueTime, v)
		}
		if err != nil {
			o.Rollback()
			return err
		}
	}
	o.Commit()
	return nil
}

//插入订单审核记录
func InsertLoanAduitRecord(uid, loanId int, content, remark string) (err error) {
	sql := `INSERT INTO loan_aduit_record (uid,loan_id,content,remark,operator_time) VALUES (?,?,?,?,NOW()) `
	o := orm.NewOrm()
	_, err = o.Raw(sql, uid, loanId, content, remark).Exec()
	return
}

//订单“退回”操作事务
func UpdateStateByOutqueue(uid, loanId, inqueueType, result int, inqueueTime, operator, remark, orderState string) (err error) {
	o := orm.NewOrm()
	o.Begin()
	defer func() {
		if err := recover(); err != nil {
			o.Rollback()
		} else {
			o.Commit()
		}
	}()
	err = UpdateLoanInfoByOutqueue(loanId, inqueueType, inqueueTime, operator, remark, orderState) //更新loan表状态
	if err != nil {
		return err
	}
	content := utils.GenerateContent(orderState, operator)
	err = InsertLoanAduitRecord(uid, loanId, content, remark) //插入订单审核记录
	if err != nil {
		return err
	}
	if orderState == "OUTQUEUE" {
		err = AddLoanManmageLog(uid, loanId, result) //插入人工审核结果日志记录
		if err != nil {
			return err
		}
	}
	return nil
}

//人工审核结果记录
func AddLoanManmageLog(uid, loanId, result int) error {
	sql := `INSERT INTO loan_manage_log  (uid,loan_id,create_time,result) VALUES (?,?,NOW(),?)`
	_, err := orm.NewOrm().Raw(sql, uid, loanId, result).Exec()
	return err
}

//订单操作事务(包括通过、正常关闭、关闭30天、永久关闭)
func UpdateStateByShutWithPass(uid, loanId int, operator, remark, orderState string) (err error) {
	o := orm.NewOrm()
	o.Begin()
	defer func() {
		if err := recover(); err != nil {
			o.Rollback()
		} else {
			o.Commit()
		}
	}()
	err = UpdateLoanInfoByOperation(loanId, orderState, operator, remark) //更新loan表状态
	if err != nil {
		return err
	}
	content := utils.GenerateContent(orderState, operator)
	err = InsertLoanAduitRecord(uid, loanId, content, remark) //插入订单审核记录
	if err != nil {
		return err
	}
	return nil
}

//订单处理
func UpdateHandleMessages(oc OrderCredit, inqueueTime string, is_req int) error {
	paras := []interface{}{}
	o := orm.NewOrm()
	err := o.Begin()
	if err != nil {
		return err
	}
	sql := ""
	if oc.OrderState == "PAUSE" || oc.OrderState == "CLOSE" || oc.OrderState == "PASS" || oc.OrderState == "CANCEL" {
		loanState := utils.StateToUsersMetadate(oc.OrderState)
		sql = `UPDATE users_metadata SET loan_state = ?`
		if loanState == 8 {
			sql += ` ,loan_state_valid_time = DATE_ADD(NOW(), INTERVAL 30 DAY)`
		}
		if loanState == 1 && oc.OrderState == "CANCEL" {
			sql += ` ,loan_state_valid_time= NOW()`
		}
		sql += ` WHERE uid = ?`
		_, err = o.Raw(sql, loanState, oc.Uid).Exec()
		if err != nil {
			o.Rollback()
			beego.Info("更新users_metadata发生错误")
			return err
		}
	}
	paras = append(paras, oc.OrderState, oc.Mark, oc.CreditOperator, oc.InqueueType)
	sql = `UPDATE loan SET order_state = ?, mark = ?, credit_operator=?,allocation_time = NULL ,handle_time = NoW(),inqueue_type=? `
	if inqueueTime == "" {
		sql += ` ,inqueue_time = NULL`
	} else {
		sql += ` ,inqueue_time = ?`
		paras = append(paras, inqueueTime)
	}
	sql += ` WHERE id = ?`
	paras = append(paras, oc.Id)
	_, err = o.Raw(sql, paras).Exec()
	if err != nil {
		o.Rollback()
		beego.Info("更新loan发生错误")
		return err
	}
	if is_req == 0 {
		sql = `INSERT INTO loan_aduit_record (loan_id,content,remark,operator_time) VALUES (?,?,?,NOW()) `
		_, err = orm.NewOrm().Raw(sql, oc.Id, oc.Content, oc.Mark).Exec()
		if err != nil {
			o.Rollback()
			beego.Info("更新loan_aduit_record发生错误")
			return err
		}
	}
	o.Commit()
	return nil
}

//订单信息
func QueryOrderByUid(loanId int) (oc OrderCredit, err error) {
	sql := `SELECT l.id , l.uid , l.money, l.loan_term_count, l.card_number,l.repayment_schedule_id, l.allocation_time, l.state, 
					ub.bank_mobile,ub.bank_name
			FROM loan AS l
			LEFT JOIN users_bankcards AS ub ON l.uid = ub.uid
			WHERE  l.id = ? `

	err = orm.NewOrm().Raw(sql, loanId).QueryRow(&oc)
	return
}

//Query loan HANDING count
func QueryHandingCount(loanId int) (count int, err error) {
	sql := `SELECT COUNT(1) FROM loan WHERE id = ? AND order_state = "HANDING" AND allocation_time <= DATE_SUB(NOW(), INTERVAL 45 MINUTE) AND is_queue =1 `
	err = orm.NewOrm().Raw(sql, loanId).QueryRow(&count)
	return
}

//update allocation_time
func UpdateAlloctionTime(loanId int) error {
	sql := `UPDATE loan SET allocation_time = NOW(),inqueue_time = NULL WHERE id = ? `
	_, err := orm.NewOrm().Raw(sql, loanId).Exec()
	return err
}

//更改订单状态(处理)
func UpdateOrderState(loanId, id int, creditOperator, orderState string) error {
	sql := `UPDATE loan SET order_state = ?,operator_id=?,credit_operator=? WHERE id = ? `
	_, err := orm.NewOrm().Raw(sql, orderState, id, creditOperator, loanId).Exec()
	return err
}

type LoanAduitRecord struct {
	Id           int
	LoanId       int
	Content      string
	Remark       string
	OperatorTime time.Time
}

//查询订单审核记录
func QueryLoanAduitRecord(uid int) (list []LoanAduitRecord, err error) {
	sql := `SELECT id,loan_id,content,remark,operator_time FROM loan_aduit_record WHERE uid = ? ORDER BY operator_time DESC,id DESC`
	o := orm.NewOrm()
	_, err = o.Raw(sql, uid).QueryRows(&list)
	return
}

//添加分配历史
func AddLoanAduitRecord(uid, loanId int, content string) error {
	sql := `INSERT INTO loan_aduit_record (uid,loan_id,content,operator_time) VALUES (?,?,?,NOW()) `
	_, err := orm.NewOrm().Raw(sql, uid, loanId, content).Exec()
	return err
}

//获取产品期数
func QueryProductTermCount() (termCount []int, err error) {
	sql := `SELECT term_count FROM product ORDER BY term_count`
	o := orm.NewOrm()
	_, err = o.Raw(sql).QueryRows(&termCount)
	return
}

//更新订单状态(api失败要出列的)
func UpdateLoanState(loanId, inqueueType int, state string) (err error) {
	sql := `UPDATE loan SET order_state = ?,inqueue_type = ?,inqueue_time = NULL WHERE id = ?`
	o := orm.NewOrm()
	_, err = o.Raw(sql, state, inqueueType, loanId).Exec()
	return
}

//更新订单状态(退回操作的)
func UpdateLoanInfoByOutqueue(loanId, inqueueType int, inqueueTime, operator, remark, orderState string) (err error) {
	sql := `UPDATE loan SET order_state = ?,credit_operator = ?,mark = ?,handle_time = NOW()`
	if inqueueTime == "" {
		sql += `,inqueue_time = NULL`
	} else {
		sql += `,inqueue_time = ?`
	}
	sql += `,inqueue_type = ? WHERE id = ?`
	o := orm.NewOrm()
	if inqueueTime == "" {
		_, err = o.Raw(sql, orderState, operator, remark, inqueueType, loanId).Exec()
	} else {
		_, err = o.Raw(sql, orderState, operator, remark, inqueueTime, inqueueType, loanId).Exec()
	}
	return
}

//更新订单信息
func UpdateLoanInfoByOperation(loanId int, orderState, operator, remark string) (err error) {
	sql := `UPDATE loan SET order_state = ?,credit_operator = ?,mark = ?,handle_time = NOW(),inqueue_time = NULL WHERE id = ?`
	o := orm.NewOrm()
	_, err = o.Raw(sql, orderState, operator, remark, loanId).Exec()
	return
}

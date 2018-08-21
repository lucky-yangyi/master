package models

import (
	"time"
	"zcm_tools/orm"
)

type RepayRecord struct {
	Uid            string    `orm:"column(uid)"`            //用户id
	VerifyRealName string    `orm:"column(verifyrealname)"` //姓名
	Account        string    `orm:"column(account)"`        //手机号
	CreateTime     time.Time `orm:"column(create_time)"`    //还款时间
	Channel        int       `orm:"column(channel)"`        //渠道
	OidPayBill     string    `orm:"column(oid_paybill)"`    //第三方订单号
	OperatorId     int       `orm:"column(operator_id)"`    //还款方式
	State          string    `orm:"column(state)"`          //还款结果
	ReturnMoney    string    `orm:"column(return_money)"`   //还款金額
	OrderNumber    string    `orm:"column(order_number)"`   //订单号

}

//获取还款管理列表
func RepayRecordList(start, pageSize int, isLimit bool, condition string, paras ...interface{}) (list []RepayRecord, err error) {
	o := orm.NewOrm()
	sql := ` SELECT
				ca.verifyrealname,
				ca.account,
				su.channel,
				su.create_time,
				su.oid_paybill,
				su.operator_id,
				su.state,
				su.return_money,
				su.uid,
				su.order_number
			FROM
				users_metadata AS ca,
				payment_record AS su
			WHERE
				ca.uid = su.uid`
	if condition != "" {
		sql += condition
	}
	if isLimit {
		sql += ` ORDER BY create_time DESC LIMIT ?,?`
		_, err = o.Raw(sql, paras, start, pageSize).QueryRows(&list)
	} else {
		sql += ` ORDER BY create_time DESC`
		_, err = o.Raw(sql, paras).QueryRows(&list)
	}
	return
}

//获取还款管理列表总数
func RepayRecordListCount(condition string, paras ...interface{}) (count int, err error) {
	o := orm.NewOrm()
	sql := `SELECT
				COUNT(1)
			FROM
				users_metadata AS ca,
				payment_record AS su
			WHERE
				ca.uid = su.uid`
	if condition != "" {
		sql += condition
	}
	err = o.Raw(sql, paras).QueryRow(&count)
	return
}

type ReturnRecord struct {
	Uid            int       `orm:"column(uid)"`             //用户id
	Account        string    `orm:"column(account)"`         //手机号
	VerifyRealName string    `orm:"column(verifyrealname)"`  //姓名
	CreateTime     time.Time `orm:"column(create_time)"`     //还款时间
	Money          float64   `orm:"column(money)"`           //借款金额
	LoanTermCount  int       `orm:"column(loan_term_count)"` //借款期限
	Rmoney         float64   `orm:"column(rmoney)"`          //到账金额
	DisplayName    string    `orm:"column(displayname)"`     //操作人
	OperatorId     int       `orm:"column(operator_id)"`     //放款方式
}

//获取回款管理列表
func ReturnRecordList(start, pageSize int, isLimit bool, condition string, paras ...interface{}) (list []ReturnRecord, err error) {
	o := orm.NewOrm()
	sql := `SELECT
				su.uid,
				ca.verifyrealname,
				ca.account,
				su.create_time,
				su.money,
				su.loan_term_count,
				su.money - su.auth_fee AS rmoney,
				sy.displayname,
				su.operator_id
			FROM
				loan AS su
			INNER JOIN users_metadata AS ca ON ca.uid = su.uid
			LEFT JOIN sys_user AS sy ON sy.id = su.operator_id
			WHERE
				su.state = "BACKING"`
	if condition != "" {
		sql += condition
	}
	if isLimit {
		sql += ` ORDER BY create_time DESC LIMIT ?,?`
		_, err = o.Raw(sql, paras, start, pageSize).QueryRows(&list)
	} else {
		sql += ` ORDER BY create_time DESC`
		_, err = o.Raw(sql, paras).QueryRows(&list)
	}
	return
}

//获取回款管理总数
func ReturnRecordListCount(condition string, paras ...interface{}) (count int, err error) {
	o := orm.NewOrm()
	sql := `SELECT
				COUNT(1)
			FROM
				loan AS su
			INNER JOIN users_metadata AS ca ON ca.uid = su.uid
			LEFT JOIN sys_user AS sy ON sy.id = su.operator_id
			WHERE
				su.state = "BACKING"`
	if condition != "" {
		sql += condition
	}
	err = o.Raw(sql, paras).QueryRow(&count)
	return
}

type RepaymentRecord struct {
	Id                   int
	Repaymen_schedule_id int
	Order_number         string
	Oid_paybill          string
	Return_money         float64
	Channel              int
	Remark               string
	Create_time          time.Time
	Return_date          time.Time
	Contract_code        string
	Uid                  int
	// Account              string
	State       string
	Loan_id     int
	Operator_id int
	Operator    string
	Displayname string
	Days        int
	StationId   int
}

// 还款情况
func GetRepaymentRecordByOidPaybill(oidpaybill string) (info *RepaymentRecord, err error) {
	sql := `select id,repayment_schedule_id,order_number,oid_paybill,return_money,channel,remark,create_time,return_date,contract_code,uid,state,loan_id,operator_id 
			from payment_record 
			where oid_paybill=?`
	err = orm.NewOrm().Raw(sql, oidpaybill).QueryRow(&info)
	if err != nil {
		return nil, err
	}
	return info, err
}

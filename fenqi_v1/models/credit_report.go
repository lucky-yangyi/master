package models

import (
	"time"
	"zcm_tools/orm"
	//	"fmt"
)

type UsersCreditReport struct {
	Id             int
	Uid            int
	State          string
	UselessDate    time.Time
	CreateTime     time.Time
	OrderNumber    string
	OidPaybill     string
	BankCardNumber string
	Remark         string
	DtOrder        time.Time
	CreditReportId int
	CreditCardName string
	Phase          int
	Price          float64
	LoanId         int
	PayMethod      int
	PayPrice       float64
	Account        string
	VerifyRealName string `orm:"column(verifyrealname)"`
	UseTime        string
	StateStr       string
}

//购买记录数量
func GetUsersCreditReportCount(condition string, paras ...interface{}) (count int, err error) {
	o := orm.NewOrm()
	sql := ` SELECT
				count(1)
			FROM
				users_metadata AS a,
				users_credit_report AS b
			WHERE
				a.uid = b.uid`
	if condition != "" {
		sql += condition
	}
	err = o.Raw(sql, paras).QueryRow(&count)
	return
}

//购买列表
func GetCreditReport(start, pageSize int, isLimit bool, condition string, paras ...interface{}) (list []UsersCreditReport, err error) {
	o := orm.NewOrm()
	sql := ` SELECT
				b.*, a.account,a.verifyrealname
			FROM
				users_metadata AS a,
				users_credit_report AS b
			WHERE
				a.uid = b.uid`
	if condition != "" {
		sql += condition
	}
	if isLimit {
		sql += ` ORDER BY b.create_time DESC LIMIT ?,?`
		_, err = o.Raw(sql, paras, start, pageSize).QueryRows(&list)
	} else {
		sql += ` ORDER BY b.create_time DESC`
		_, err = o.Raw(sql, paras).QueryRows(&list)
	}
	return
}

//查看个人购买记录
func GetPerCreditReport(start, pageSize, uid int) (list []UsersCreditReport, err error) {
	o := orm.NewOrm()
	sql := ` SELECT
				b.*, a.account
			FROM
				users_metadata AS a,
				users_credit_report AS b
			WHERE
				a.uid= b.uid AND b.uid = ?`
	sql += ` ORDER BY b.create_time DESC LIMIT ?,?`
	_, err = o.Raw(sql, uid, start, pageSize).QueryRows(&list)
	return
}

//查看个人购买记录数量
func GetPerUsersCreditReportCount(uid int) (count int, err error) {
	o := orm.NewOrm()
	sql := `SELECT
				count(1)
			FROM
				users_metadata AS a,
				users_credit_report AS b
			WHERE
				a.uid= b.uid AND b.uid = ?`
	err = o.Raw(sql, uid).QueryRow(&count)
	return
}

type CreditReportDate struct {
	State  string
	Remark string
}

//查看订单状态
func GeUidStateFromUsersCreditReport(id int) (repdata CreditReportDate, err error) {
	o := orm.NewOrm()
	sql := "SELECT state,remark FROM users_credit_report WHERE id =?"
	err = o.Raw(sql, id).QueryRow(&repdata)
	return
}

//更新状态
func UpdateUidStateFromUsersCreditReport(id int, state string) (err error) {
	o := orm.NewOrm()
	sql := "UPDATE users_credit_report SET state =? WHERE id =?"
	_, err = o.Raw(sql, state, id).Exec()
	return
}

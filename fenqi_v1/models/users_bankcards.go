package models

import (
	"time"
	"zcm_tools/orm"
)

//银行卡
type UsersBankcards struct {
	Id             int       `orm:"column(id);auto"`
	Uid            int       `orm:"column(uid);null"`
	BankId         int       `orm:"column(bank_id);null"`
	BankName       string    `orm:"column(bank_name);size(50);null"`   //银行卡名称
	CardNumber     string    `orm:"column(card_number);size(30);null"` //卡号
	State          string    `orm:"column(state);null"`
	Channel        int8      `orm:"column(channel);null"`
	CreateTime     time.Time `orm:"column(create_time);type(timestamp);null;auto_now_add"`
	Account        string    `orm:"column(account);size(11);null"`
	RealPayMoney   float64   `orm:"column(real_pay_money);null;digits(10);decimals(3)"`
	EncashMoney    float64   `orm:"column(encash_money);null;digits(10);decimals(3)"`
	BankMobile     string    `orm:"column(bank_mobile);size(11);null"`
	NoAgree        string    `orm:"column(no_agree);size(20);null"`
	RetCode        string    `orm:"column(ret_code);size(10);null"`
	Deadbeat       string    `orm:"column(deadbeat);size(2);null"`
	LogoImage      string    //银行icon
	Verifyrealname string    //用户姓名
}

func QueryBankcardByUid(uid int) (u UsersBankcards, err error) {
	o := orm.NewOrm()
	sql := `SELECT
			card_number
			FROM users_bankcards
			WHERE uid = ?`
	err = o.Raw(sql, uid).QueryRow(&u)
	return
}

//根据uid获取用户银行卡
func QueryUsersBankcards(uid, start, pageSize int) (list []UsersBankcards, err error) {
	o := orm.NewOrm()
	sql := `SELECT
				ub.uid,
				ub.bank_name,
				ub.bank_mobile,
				b.logo_image,
				ub.card_number,
				um.verifyrealname
			FROM users_bankcards AS ub
			INNER JOIN bankcard AS b
			ON ub.bank_id = b.id
			LEFT JOIN users_metadata AS um
			ON ub.uid = um.uid
			WHERE ub.uid = ?
			ORDER BY ub.create_time DESC
			LIMIT ?,?`
	_, err = o.Raw(sql, uid, start, pageSize).QueryRows(&list)
	return
}

//根据uid获取用户银行卡总数
func QueryUsersBankcardsCount(uid int) (count int, err error) {
	o := orm.NewOrm()
	sql := `SELECT
				COUNT(1)
			FROM users_bankcards AS ub
			INNER JOIN bankcard AS b
			ON ub.bank_id = b.id
			LEFT JOIN users_metadata AS um
			ON ub.uid = um.uid
			WHERE ub.uid = ?`
	err = o.Raw(sql, uid).QueryRow(&count)
	return
}

//获取用户绑定的卡
func QueryUsersBankcardByUid(uid int) (userCard UsersBankcards, err error) {
	o := orm.NewOrm()
	sql := `SELECT
			id,
			uid,
			card_number
		FROM users_bankcards
		WHERE uid = ?
		AND state = "USING"
		ORDER BY create_time DESC
		LIMIT 1`
	err = o.Raw(sql, uid).QueryRow(&userCard)
	return
}

//更新用户绑定的卡
func UpdateUsersBankcard(uid int) (err error) {
	o := orm.NewOrm()
	sql := `UPDATE users_bankcards SET state = "ABANDON" WHERE uid = ? AND state = "USING"`
	_, err = o.Raw(sql, uid).Exec()
	return
}

//添加银行卡解绑日志
func AddUsersBankcardLog(uid, bankcardId, operatorId int, cardNumber, operator string) (err error) {
	o := orm.NewOrm()
	o.Using("fq_log")
	sql := `INSERT INTO users_bankcard_log (uid,bankcard_id,operator_id,card_number,operator,operator_time) VALUES(?,?,?,?,?,NOW())`
	_, err = o.Raw(sql, uid, bankcardId, operatorId, cardNumber, operator).Exec()
	return
}

package models

import (
	"zcm_tools/orm"
)

type UsersCoupon struct {
	Id           int    `orm:"column(id);auto"`
	Uid          int    `orm:"column(uid)"`                            //用户id
	Account      string `orm:"column(account);size(11);null"`          //用户账号
	CouponType   string `orm:"column(coupon_type);size(11)"`           //优惠券类型
	IsUsed       string `orm:"column(is_used);null"`                   //使用情况
	BeginTime    string `orm:"column(begin_time);type(datetime);null"` //开始时间
	EndTime      string `orm:"column(end_time);type(datetime);null"`   //结束时间
	UseTime      string `orm:"column(use_time);type(datetime);null"`   //使用时间
	GetTime      string `orm:"column(get_time);type(datetime);null"`   //获得时间
	Type         string `orm:"column(type);null"`
	LoanId       string `orm:"column(loan_id);size(20);null"`
	FavorMoney   int    `orm:"column(favor_money);null"`           //优惠金额
	Monetary     int    `orm:"column(monetary);null"`              //满足金额
	CouponName   string `orm:"column(coupon_name);size(100);null"` //优惠券名称
	PMaxDay      int    `orm:"column(p_max_day);null"`
	IsRead       int    `orm:"column(is_read);null"` //是否已读:0-已读 1-未读
	ApplyExplain string //获得原因
}

//根据uid获取用户优惠券
func QueryUsersCoupons(uid, start, pageSize int) (list []UsersCoupon, err error) {
	o := orm.NewOrm()
	sql := `SELECT
				uc.uid,
				DATE_FORMAT(uc.begin_time,"%Y-%m-%d") AS begin_time,
				DATE_FORMAT(uc.end_time,"%Y-%m-%d") AS end_time,
				uc.get_time,
				c.coupon_name,
				c.favor_money,
				c.monetary,
				c.apply_explain,
				uc.is_read,
				uc.is_used
			FROM users_coupon AS uc
			INNER JOIN coupon AS c
			ON uc.coupon_type = c.id
			WHERE uc.uid = ?
			ORDER BY uc.get_time DESC LIMIT ?,?`
	_, err = o.Raw(sql, uid, start, pageSize).QueryRows(&list)
	return
}

//根据uid获取用户优惠券总数
func QueryUsersCouponsCount(uid int) (count int, err error) {
	o := orm.NewOrm()
	sql := `SELECT
				COUNT(1)
			FROM users_coupon AS uc
			INNER JOIN coupon AS c
			ON uc.coupon_type = c.id
			WHERE uc.uid = ?`
	err = o.Raw(sql, uid).QueryRow(&count)
	return
}

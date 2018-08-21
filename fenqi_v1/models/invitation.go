package models

import (
	"strconv"
	"time"
	"zcm_tools/orm"
)

type CouponInfo struct {
	Id         int
	CouponName string
	FavorMoney int
	Monetary   int
	CType      string `orm:"column(type)"`
	TermLimit  int
}

func AddInvitationReward(uid, money int, iType int) (err error) {
	var oldUid int
	o := orm.NewOrm()
	if err != nil {
		return
	}
	//查询邀请人uid
	sql := `SELECT old_uid FROM invitation WHERE new_uid = ?`
	err = o.Raw(sql, uid).QueryRow(&oldUid)
	if err != nil {
		return nil
	}
	if oldUid == 0 {
		return
	}

	var count int
	//查询是否有邀请奖励发放记录
	sql = `SELECT COUNT(1) FROM invitation_reward WHERE new_uid = ? AND old_uid = ? AND reward_type = ` + strconv.Itoa(iType)
	err = o.Raw(sql, uid, oldUid).QueryRow(&count)
	if err != nil || count > 0 {
		return nil
	}
	//查询手机号
	var account int
	sql = `SELECT account FROM users WHERE id = ?`
	err = o.Raw(sql, uid).QueryRow(&account)
	if err != nil || count > 0 {
		return nil
	}
	//查询优惠券信息
	var couponInfo CouponInfo
	sql = `SELECT id,coupon_name,favor_money,monetary,type,term_limit FROM coupon WHERE id =?`
	err = o.Raw(sql, 7).QueryRow(&couponInfo)
	if err != nil {
		return
	}
	err = o.Begin()
	if err != nil {
		return
	}
	//更新用户邀请记录状态
	sql = `UPDATE invitation SET state = "AUTHEN" WHERE new_uid = ?`
	_, err = o.Raw(sql, uid).Exec()
	if err != nil {
		o.Rollback()
		return
	}
	//插入用户邀请奖励记录
	sql = `INSERT INTO invitation_reward SET old_uid=?,new_uid=?,money=?,reward_type=?,createtime=NOW()`
	_, err = o.Raw(sql, oldUid, uid, money, iType).Exec()
	if err != nil {
		o.Rollback()
		return
	}
	//用户添加减免券
	endTimeStr := time.Now().AddDate(0, 0, 30).Format("2006-01-02") + " 23:59:59"
	sql = `INSERT INTO users_coupon SET uid=?,account=?,coupon_type=?,begin_time=?,type=?,favor_money=?,monetary=?,coupon_name=?,end_time=?,get_time=?`
	_, err = o.Raw(sql, oldUid, account, couponInfo.Id, time.Now().AddDate(0, 0, 7), couponInfo.CType, couponInfo.FavorMoney, couponInfo.Monetary, couponInfo.CouponName, endTimeStr, time.Now()).Exec()
	if err != nil {
		o.Rollback()
		return
	}
	o.Commit()
	return
}

//邀请记录
type InvitationRecord struct {
	OldUid     int       //邀请人id
	NewUid     int       //被邀请人id
	NewAccount string    //被邀请人手机号
	CreateTime time.Time //邀请时间
	RewardType int       //奖励类型:0-话费奖励 1-抵用券奖励
}

//获取用户邀请记录
func QueryInvitationRecords(uid, start, pageSize int) (list []InvitationRecord, err error) {
	o := orm.NewOrm()
	sql := `SELECT
				i.old_uid,
				i.new_uid,
				i.new_account,
				i.create_time,
				ir.reward_type
			FROM invitation AS i
			INNER JOIN invitation_reward AS ir
			ON i.old_uid = ir.old_uid
			AND i.new_uid = ir.new_uid
			WHERE i.old_uid = ?
			ORDER BY i.create_time DESC
			LIMIT ?,?`
	_, err = o.Raw(sql, uid, start, pageSize).QueryRows(&list)
	return
}

//获取用户邀请记录总数
func QueryInvitationRecordCount(uid int) (count int, err error) {
	o := orm.NewOrm()
	sql := `SELECT
				COUNT(1)
			FROM invitation AS i
			INNER JOIN invitation_reward AS ir
			ON i.old_uid = ir.old_uid
			AND i.new_uid = ir.new_uid
			WHERE i.old_uid = ?`
	err = o.Raw(sql, uid).QueryRow(&count)
	return
}

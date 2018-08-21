package models

import (
	"time"
	"zcm_tools/orm"
)

//查询用户的认证状态
type UserAuthState struct {
	CurrentAuthState int `description:"1:未认证2:认证完成(在有效时间内)3:认证完成(已过期)4：授信关闭30天5：授信永久关闭6：授信驳回"`
}

//更新授信状态
func UpadteUserAuthState(current_auth_state, uid int) error {
	sql := `UPDATE users_auth SET current_auth_state=? WHERE uid=?`
	_, err := orm.NewOrm().Raw(sql, current_auth_state, uid).Exec()
	return err
}

//用户认证
type UsersAuth struct {
	Uid            int       //用户id
	RealNameTime   time.Time //认证时间
	Account        string    //手机号
	Verifyrealname string    //姓名
	IdCard         string    //身份证
	Balance        float64   //可借额度
	UseBalance     float64   //已使用额度
	RemaindBalance float64   //剩余额度
	Source         string    //注册渠道
	State          string    //最终授信状态
	AuthState      int       //users_metadata的授信状态
	SalesmanId     int
	IsSalemanType  string
}

//查询用户认证次数
func QueryUserCreditCount(uid int) (count int) {
	o := orm.NewOrm()
	sql := `SELECT COUNT(1) FROM users_auth WHERE uid = ?`
	o.Raw(sql, uid).QueryRow(&count)
	return
}

//查询用户最后一次授信
func QueryUserLastCredit(uid int) (state int, err error) {
	o := orm.NewOrm()
	sql := `SELECT current_auth_state
			FROM users_auth
			WHERE uid = ?
			AND is_valid = 2 AND is_real_name = 2
			ORDER BY submit_time DESC LIMIT 1`
	err = o.Raw(sql, uid).QueryRow(&state)
	return
}

//查询用户认证列表
func QueryUsersAuthList(start, pageSize int, condition string, paras ...interface{}) (list []*UsersAuth, err error) {
	o := orm.NewOrm()
	sql := `SELECT 
				um.uid,
				u.account,
				um.verifyrealname,
				um.id_card,
				ua.real_name_time,
				um.balance,
				um.use_balance,
				(um.balance-um.use_balance) AS remaind_balance,
				ua.current_auth_state AS auth_state,
				u.source,
				u.salesman_id,
				ca.state
			FROM users_metadata um 
			INNER JOIN users u ON um.uid = u.id
			INNER JOIN users_auth ua ON um.uid = ua.uid AND ua.is_valid = 1 AND ua.is_real_name = 2
			LEFT JOIN credit_aduit ca ON um.uid = ca.uid AND ca.is_now = 1
			WHERE 1=1 `
	if condition != "" {
		sql += condition
	}
	sql += ` ORDER BY ua.real_name_time DESC LIMIT ?,?`
	_, err = o.Raw(sql, paras, start, pageSize).QueryRows(&list)
	return
}

//查询用户认证总数
func QueryUsersAuthCount(condition string, paras ...interface{}) (count int, err error) {
	o := orm.NewOrm()
	sql := `SELECT 
				COUNT(1)
			FROM users_metadata um 
			INNER JOIN users u ON um.uid = u.id
			INNER JOIN users_auth ua ON um.uid = ua.uid AND ua.is_valid = 1 AND ua.is_real_name = 2
			LEFT JOIN credit_aduit ca ON um.uid = ca.uid AND ca.is_now = 1
			WHERE 1=1 `
	if condition != "" {
		sql += condition
	}
	err = o.Raw(sql, paras).QueryRow(&count)
	return
}

func GetUsersAuthCount(uid int) (count int, err error) {
	o := orm.NewOrm()
	sql := `SELECT
				COUNT(1)
			FROM users_auth
			WHERE is_valid = 1
			AND is_real_name = 2
			AND uid = ?`
	err = o.Raw(sql, uid).QueryRow(&count)
	return
}

//获取用户银行卡绑定情况
func QueryUsersBankcardBind(uid int) (isBindCard int, err error) {
	o := orm.NewOrm()
	sql := `SELECT is_bind_card FROM users_auth WHERE is_valid = 1 AND uid = ?`
	err = o.Raw(sql, uid).QueryRow(&isBindCard)
	return
}

//修改用户绑卡状态
func UpdateUsersBankcardBind(uid, isBindCard int) (err error) {
	o := orm.NewOrm()
	sql := `UPDATE users_auth SET is_bind_card = ? WHERE uid = ? AND is_valid = 1`
	_, err = o.Raw(sql, isBindCard, uid).Exec()
	return
}

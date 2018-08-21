package models

import (
	"zcm_tools/orm"
)

type BlackSignSave struct {
	Uid            int    //用户id
	Account        string //账号
	SubmitManId    int    // 提交人Id
	SubmitTime     string // 提交时间
	State          int    //状态 ：0：待审核  :1：已标记 :2： 已退回
	IsTag          int    //入黑类型 ：0 未标记欺诈 1 家人代偿,2 本人还款意愿差,3 其他
	BlackTime      string //审核时间
	CheckManId     int    //审核人id
	Remark         string //备注
	CollectionType string // 催收阶段   0：M0,  1 :S0,  2:S1,  3:S2  4:委外催收

}
type BlackSign struct {
	Uid            int    //用户id
	Account        string //账号
	SubmitMan      string // 提交人姓名
	SubmitTime     string // 提交时间
	State          int    //状态 ：0：待审核  :1：已标记 :2： 已退回
	IsTag          int    //入黑类型 ：0 未标记欺诈 1 家人代偿,2 本人还款意愿差,3 其他
	CheckTime      string //审核时间
	CheckManId     int    //审核人id
	Remark         string //备注
	CollectionType string //催收阶段   0：M0,  1 :S0,  2:S1,  3:S2  4:委外催收
	CheckMan       string //审核人
	ActionType     string //行动分类
}
type LoanOverdueDate struct {
	OverdueDays       int // 逾期天数
	DistributionCount int //分配次数
}

// save
func SaveBlackSign(b BlackSignSave) (err error) {
	o := orm.NewOrm()
	sql := ` INSERT INTO black_sign (uid,account,submit_man_id,submit_time,state,is_tag,remark,collection_type) VALUES (?,?,?,?,?,?,?,?) `
	_, err = o.Raw(sql, b.Uid, b.Account, b.SubmitManId, b.SubmitTime, b.State, b.IsTag, b.Remark, b.CollectionType).Exec()
	return
}

// 入黑标记list
func GetBlackSignList(condition string, params []interface{}, begin, size int) (bs []BlackSign, err error) {
	o := orm.NewOrm()
	sql := `SELECT bs.uid, bs.account,su.displayname submit_man,bs.submit_time,bs.is_tag,bs.remark,bs.state,bs.collection_type,bs.check_time,s1.displayname as check_man,bs.action_type
			FROM  black_sign bs
			LEFT JOIN sys_user su ON bs.submit_man_id=su.id
			LEFT JOIN sys_user AS s1 ON s1.id=bs.check_man_id WHERE 1=1 `
	if condition != "" {
		sql += condition
	}
	sql += ` ORDER BY bs.submit_time DESC  LIMIT ?,? `
	_, err = o.Raw(sql, params, begin, size).QueryRows(&bs)
	return
}

// 获取总页数
func GetBlackSignPageCount(condition string, params []interface{}) (pageCount int, err error) {
	o := orm.NewOrm()
	sql := ` SELECT COUNT(1) FROM black_sign bs WHERE 1=1 `
	if condition != "" {
		sql += condition
	}
	err = o.Raw(sql, params).QueryRow(&pageCount)
	return
}

// 更改状态
func HandlerBlackSignUserState(state, checkManId, blackUserId, is_tag int, tagtype string) (err error) {
	o := orm.NewOrm()
	sql := ` UPDATE black_sign bs SET bs.state = ?,bs.check_man_id = ?,check_time=now() WHERE bs.uid = ? `
	_, err = o.Raw(sql, state, checkManId, blackUserId).Exec()
	if state == 2 && err == nil {
		var is_black = 2
		sql2 := ` UPDATE users SET is_tag =?, is_black =?,tag_type=?,black_time= NOW()  WHERE id = ? `
		_, err = o.Raw(sql2, is_tag, is_black, tagtype, blackUserId).Exec()
	}
	return
}

//更新入黑标记行动分类
func UpdateBlackSignActionType(uid, submitId int, actionType string) (err error) {
	sql := `UPDATE black_sign SET action_type = ?,submit_man_id = ?,submit_time = NOW() WHERE uid = ?`
	_, err = orm.NewOrm().Raw(sql, actionType, submitId, uid).Exec()
	return
}

//批量更改状态
func BatchUpdateBlackSignState(state, checkManId int, isTags []int, tagTypes, uids []string) (err error) {
	o := orm.NewOrm()
	sql := ` UPDATE black_sign bs SET bs.state = ?,bs.check_man_id = ?,check_time=NOW() WHERE bs.uid = ? `
	update, err := o.Raw(sql).Prepare()
	defer update.Close()
	if err == nil {
		o.Begin()
		for k, uid := range uids {
			_, err = update.Exec(state, checkManId, uid)
			if state == 2 && err == nil {
				sql2 := ` UPDATE users SET is_tag =?, is_black = 2,tag_type=?,black_time= NOW()  WHERE id = ? `
				_, err = o.Raw(sql2, isTags[k], tagTypes[k], uid).Exec()
			}
		}
		if err != nil {
			o.Rollback()
			return err
		}
		o.Commit()
	}
	return
}

//
func QueryBlackSign(uid int) (state int, err error) {
	o := orm.NewOrm()
	sql := ` SELECT state FROM black_sign WHERE  uid = ? `
	err = o.Raw(sql, uid).QueryRow(&state)
	return
}

// 根据id 获取入黑类型
func GetIsTag(uid int) (istag int, err error) {
	o := orm.NewOrm()
	sql := ` SELECT is_tag FROM black_sign WHERE uid = ? `
	err = o.Raw(sql, uid).QueryRow(&istag)
	return
}

// 更新已退回用户
func UpdateReturnUser(bs BlackSignSave) (err error) {
	o := orm.NewOrm()
	sql := ` UPDATE black_sign SET submit_man_id=?,submit_time=?,state=?,is_tag=?,check_time=null,check_man_id=null,remark=?,collection_type=? WHERE uid= ? `
	_, err = o.Raw(sql, bs.SubmitManId, bs.SubmitTime, bs.State, bs.IsTag, bs.Remark, bs.CollectionType, bs.Uid).Exec()
	return
}

// 重新标记判断是否存在
func UpdateQueryBlackSign(uid int) (count int, err error) {
	o := orm.NewOrm()
	sql := ` SELECT count(1) count FROM black_sign WHERE   uid = ? `
	err = o.Raw(sql, uid).QueryRow(&count)
	return
}

// 根据 loanId 查询未还最大逾期天数
func GetOverdueDaysByLoanId(loanId int) (data LoanOverdueDate, err error) {
	o := orm.NewOrm()
	sql := `SELECT MAX(rs.overdue_days)
			FROM repayment_schedule rs
			WHERE rs.loan_id =? AND rs.state="backing" `
	err = o.Raw(sql, loanId).QueryRow(&data)
	return
}

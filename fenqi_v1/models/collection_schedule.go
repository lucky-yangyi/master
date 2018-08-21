package models

import (
	"zcm_tools/orm"
)

//催收排班
type CollectionSchedule struct {
	Id        int
	SysUid    int    //催收人ID
	Name      string //催收人姓名
	Type      int    //催收人所属分组
	PlanYear  int    //排班日期（年）
	PlanMonth int    //排班日期（月）
	PlanDay   int    //排班日期（日）
	State     int    //催收人排班状态（0-班/1-休）
}
type CollectionScheduleRecord struct {
	Name        string
	SysUid      int
	Type        int
	Days        [31]int
	AccountType int //是否为机器人:1-是
}

//获取催收排班详情
func GetCollectionSchedule(where string, pars []interface{}, group string) (collectionSchedules []CollectionSchedule, err error) {
	sql := `SELECT u.id AS sys_uid,u.displayname AS name,cs.id,cs.type,cs.state,cs.plan_year,cs.plan_month,cs.plan_day
			FROM sys_user AS u
			INNER JOIN sys_station AS s ON u.station_id = s.id
			INNER JOIN sys_station_type AS t ON s.id = t.station_id
			LEFT JOIN collection_schedule AS cs ON u.id = cs.sys_uid
			WHERE cs.state = 1
			AND u.accountstatus = '启用'`
	if where != "" {
		sql += where
	}
	if group != "" {
		sql += ` AND t.type in(` + group + `)`
	} else {
		sql += ` AND t.type in (11,12,13,14,15,16)`
	}
	sql += " ORDER BY sys_uid"
	_, err = orm.NewOrm().Raw(sql, pars).QueryRows(&collectionSchedules)
	return
}

func GetCollectionScheduleGroup(group_type, displayName string) (collectionScheduleRecords []CollectionScheduleRecord, err error) {
	sql := `SELECT u.id as sys_uid,s.org_id,u.displayname AS name,t.type,u.account_type
			FROM sys_user AS u
			INNER JOIN sys_station AS s
			ON u.station_id=s.id
			INNER JOIN sys_station_type AS t
			ON s.id=t.station_id
			WHERE u.accountstatus='启用'`
	if group_type != "" {
		sql += " AND t.type in (" + group_type + ")"
	} else {
		sql += " AND t.type  in (11,12,13,14,15,16)"
	}
	if displayName != "" {
		sql += ` AND u.displayname LIKE '%` + displayName + `%'`
	}
	sql += " ORDER BY sys_uid "
	_, err = orm.NewOrm().Raw(sql).QueryRows(&collectionScheduleRecords)
	return
}

//检查该条记录是否存在
func CheckCollectionScheduleIsExist(sysUid, types, year, month, day int) bool {
	sql := `SELECT COUNT(1) FROM collection_schedule WHERE sys_uid = ? AND type = ? AND plan_year = ? AND plan_month = ? AND plan_day = ?`
	var count int
	err := orm.NewOrm().Raw(sql, sysUid, types, year, month, day).QueryRow(&count)
	if err != nil {
		return false
	}
	if count > 0 {
		return true
	}
	return false
}

//更新排班催收分组
func UpdateCollectionScheduleType(sysUid, types int) (err error) {
	sql := `UPDATE collection_schedule SET type = ? WHERE sys_uid = ?`
	_, err = orm.NewOrm().Raw(sql, types, sysUid).Exec()
	return
}

//更新排班信息
func UpdateMultilCollectionSchedules(sysUid, state, year, month, day int, date string, periodTypes []int) (err error) {
	sql := `UPDATE collection_schedule SET state = ?,date = ? WHERE sys_uid = ? AND type = ? AND plan_year = ? AND plan_month = ? AND plan_day = ?`
	o := orm.NewOrm()
	o.Begin()
	update, err := o.Raw(sql).Prepare()
	defer update.Close()
	if err == nil {
		for _, v := range periodTypes {
			_, err = update.Exec(state, date, sysUid, v, year, month, day)
			if err != nil {
				o.Rollback()
				return
			}
		}
	}
	o.Commit()
	return
}

//插入排班信息
func InsertMultilCollectionSchedules(sysUid, state, year, month, day int, date string, periodTypes []int) (err error) {
	sql := `INSERT INTO collection_schedule(sys_uid,type,state,plan_year,plan_month,plan_day,date)  VALUES(?,?,?,?,?,?,?)`
	o := orm.NewOrm()
	o.Begin()
	insert, err := o.Raw(sql).Prepare()
	defer insert.Close()
	if err == nil {
		for _, v := range periodTypes {
			_, err = insert.Exec(sysUid, v, state, year, month, day, date)
			if err != nil {
				o.Rollback()
				return
			}
		}
	}
	o.Commit()
	return
}

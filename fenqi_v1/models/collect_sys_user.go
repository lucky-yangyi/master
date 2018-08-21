package models

import "zcm_tools/orm"

// 催收系统用户列表
func CollectSysUserList(condition string, pars []string, begin, count int) (list []SysUserMini, err error) {
	sql := `SELECT su.id,su.name,su.displayname,su.email,su.accountstatus, sr.displayname role_name,ss.name as station_name,su.last_operation_time
			FROM sys_user su
			LEFT JOIN sys_role sr ON su.role_id=sr.id
			LEFT JOIN sys_station ss ON su.station_id=ss.id
			WHERE 1=1 AND su.is_collect_account = 1 `
	sql += condition
	sql += " ORDER BY su.name LIMIT ?, ?"
	_, err = orm.NewOrm().Raw(sql, pars, begin, count).QueryRows(&list)
	if err != nil {
		return nil, err
	}
	return list, nil
}

func CollectSysUserCounts(name string) (count int, err error) {
	sql := `SELECT COUNT(1) FROM sys_user WHERE name = ?`
	o := orm.NewOrm()
	o.Using("read")
	err = o.Raw(sql, name).QueryRow(&count)
	return
}

func CollectSysUserCount(condition string, pars []string) int {
	sql := `SELECT count(1)
			FROM sys_user su
			LEFT JOIN sys_role sr ON su.role_id=sr.id
			LEFT JOIN sys_station ss ON su.station_id=ss.id
			WHERE 1=1 AND  su.is_collect_account = 1`
	sql += condition
	var count int
	orm.NewOrm().Raw(sql, pars).QueryRow(&count)
	return count
}

func CollectSysUserDetail(uid int) (user *SysUserMini, err error) {
	sql := `SELECT su.id,su.name,su.password,su.displayname,su.email,su.role_id,
			su.accountstatus,su.station_id,su.qn_account,su.qn_password,su.place_role,
			su.tape_id,su.create_time,su.account_type,su.last_operation_time, sr.displayname role_name
			FROM sys_user su
			LEFT JOIN sys_role sr ON su.role_id=sr.id
			WHERE su.id=?`
	o := orm.NewOrm()
	o.Using("read")
	err = o.Raw(sql, uid).QueryRow(&user)
	return
}

func (u *SysUserMini) CollectInsert() error {
	sql := `INSERT INTO sys_user (name, password,  displayname, email, accountstatus, role_id, create_time,station_id,qn_account,qn_password,place_role,account_type,is_collect_account)
			values(?, ?, ?, ?, ?, ?, now(),?,?,?,?,?,?)`
	_, err := orm.NewOrm().Raw(sql, u.Name, u.Password, u.Displayname, u.Email, u.Accountstatus, u.Role_id, u.Station_id, u.Qn_account, u.Qn_password, u.Place_role, u.AccountType, u.IsCollectAccount).Exec()
	return err
}

func CollectDeleteUser(uid int) error {
	sql := `DELETE FROM sys_user WHERE id=?`
	_, err := orm.NewOrm().Raw(sql, uid).Exec()
	return err
}

func GetCollectUserType(stationId int) (userType int, err error) {
	sql := `SELECT t.type
			FROM sys_user AS u
			INNER JOIN sys_station AS s
			ON u.station_id=s.id
			INNER JOIN sys_station_type AS t
			ON s.id=t.station_id
			WHERE u.station_id = ?`
	err = orm.NewOrm().Raw(sql, stationId).QueryRow(&userType)
	return
}

//根据岗位ID查找角色ID
func QueryCollectRoleIdByStationId(stationId int) (roleId int, err error) {
	sql := `SELECT role_id FROM sys_station WHERE id = ?`
	err = orm.NewOrm().Raw(sql, stationId).QueryRow(&roleId)
	return
}

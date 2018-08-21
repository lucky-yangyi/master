package models

import (
	"time"
	"zcm_tools/orm"
)

type SysUser struct {
	Id                int
	Code              string `orm:"column(code);null"`
	Name              string `orm:"column(name);null"`
	Password          string `orm:"column(password);null"`
	Password1         string `orm:"column(password1);null"`
	DisplayName       string `orm:"column(displayname);null"`
	Sex               string `orm:"column(sex);null"`
	Phone             string `orm:"column(phone);null"`
	Email             string `orm:"column(email);null"`
	AccountStatus     string `orm:"column(accountstatus);null"`
	Remark            string `orm:"column(remark);null"`
	RoleId            int    `orm:"column(role_id);null"`
	StationId         int    `orm:"column(station_id);null"`
	QnAccount         string `description:"青牛账号"`
	QnPassword        string `description:"青牛密码"`
	PlaceRole         string `description:"座席岗位"`
	TapeId            int    `description:"坐席最新通话记录id"`
	AccountType       int    `json:"account_type" description:"账户类型，0外部账户，1内部账户"`
	LastOperationTime time.Time
	LoginState        string `orm:"column(login_state)"` //登录状态

}

// type SysUserExpand struct {
// 	Id            int
// 	Code          string `orm:"column(code);null"`
// 	Name          string `orm:"column(name);null"`
// 	Password      string `orm:"column(password);null"`
// 	Password1     string `orm:"column(password1);null"`
// 	DisplayName   string `orm:"column(displayname);null"`
// 	Sex           string `orm:"column(sex);null"`
// 	Phone         string `orm:"column(phone);null"`
// 	Email         string `orm:"column(email);null"`
// 	AccountStatus string `orm:"column(accountstatus);null"`
// 	Remark        string `orm:"column(remark);null"`
// 	RoleId        string `orm:"column(role_id);null"`
// }

type SysUserMini struct {
	Id                int
	Name              string
	Password          string
	Displayname       string
	Email             string
	Role_id           int
	RoleName          string
	Accountstatus     string
	Secret            string // 验证码密钥
	AuthURL           string
	Station_id        int    //岗位ID
	StationName       string //岗位名称
	Qn_account        string `description:"青牛账号"`
	Qn_password       string `description:"青牛密码"`
	Place_role        string `description:"座席岗位"`
	TapeId            int    `description:"坐席最新通话记录id"`
	CreateTime        string `description:"创建时间"`
	CallType          string `description:"呼叫类型"`
	AccountType       int    `json:"account_type" description:"账户类型，0外部账户，1内部账户"`
	LastOperationTime string //最后操作时间
	IsCollectAccount  int    `description:"是不是催收账号：0-不是 1-是"`
}

func FindByIdSysUser(sysId int) (sysUser SysUser, err error) {
	o := orm.NewOrm()
	err = o.Raw("select email from sys_user where id=?", sysId).QueryRow(&sysUser)
	return
}

// 系统用户列表
func SysCollectionUserList(condition string, pars []string, begin, count int) (list []SysUserMini, err error) {
	sql := `SELECT su.*, sr.displayname role_name,ifnull(tr.call_type,"挂断") as call_type
			FROM sys_user su 
			LEFT JOIN sys_role sr ON su.role_id=sr.id
			LEFT JOIN tape_record tr ON su.tape_id=tr.id
			WHERE 1=1`
	sql += condition
	sql += " ORDER BY create_time DESC LIMIT ?, ?"
	o := orm.NewOrm()
	o.Using("read")
	_, err = o.Raw(sql, pars, begin, count).QueryRows(&list)
	if err != nil {
		return nil, err
	}
	return list, nil
}

func SysCollectionUserCount(condition string, pars []string) int {
	sql := `SELECT count(1) 
			FROM sys_user su 
			LEFT JOIN sys_role sr ON su.role_id=sr.id
			LEFT JOIN tape_record tr ON su.tape_id=tr.id
			WHERE 1=1`
	sql += condition
	var count int
	o := orm.NewOrm()
	o.Using("read")
	o.Raw(sql, pars).QueryRow(&count)
	return count
}

func Login(name, password string) (v *SysUser, err error) {
	o := orm.NewOrm()
	sql := `SELECT * FROM sys_user WHERE name= ? and password=? and accountStatus='启用'  `
	// fmt.Println(sql, name, password)
	err = o.Raw(sql, name, password).QueryRow(&v)
	return
}

// 系统用户列表
func SysUserList(condition string, pars []string, begin, count int) (list []SysUserMini, err error) {
	sql := `SELECT su.*, sr.displayname role_name,ss.name as station_name,su.last_operation_time
			FROM sys_user su
			LEFT JOIN sys_role sr ON su.role_id=sr.id
			LEFT JOIN sys_station ss ON su.station_id=ss.id
			WHERE 1=1`
	sql += condition
	sql += " ORDER BY su.name LIMIT ?, ?"
	_, err = orm.NewOrm().Raw(sql, pars, begin, count).QueryRows(&list)
	if err != nil {
		return nil, err
	}
	return list, nil
}

func SysUserCount(condition string, pars []string) int {
	sql := `SELECT count(1)
			FROM sys_user su
			LEFT JOIN sys_role sr ON su.role_id=sr.id
			LEFT JOIN sys_station ss ON su.station_id=ss.id
			WHERE 1=1`
	sql += condition
	var count int
	orm.NewOrm().Raw(sql, pars).QueryRow(&count)
	return count
}

func SysUserDetail(uid int) (user *SysUserMini, err error) {
	sql := `SELECT su.* , sr.displayname role_name
			FROM sys_user su
			LEFT JOIN sys_role sr ON su.role_id=sr.id
			WHERE su.id=?`
	o := orm.NewOrm()
	o.Using("read")
	err = o.Raw(sql, uid).QueryRow(&user)
	return
}

func (u *SysUserMini) Insert() error {
	sql := `INSERT INTO sys_user (name, password,  displayname, email, accountstatus, role_id, create_time,station_id,qn_account,qn_password,place_role,account_type)
			values(?, ?, ?, ?, ?, ?, now(),?,?,?,?,?)`
	_, err := orm.NewOrm().Raw(sql, u.Name, u.Password, u.Displayname, u.Email, u.Accountstatus, u.Role_id, u.Station_id, u.Qn_account, u.Qn_password, u.Place_role, u.AccountType).Exec()
	return err
}

func (u *SysUserMini) Update() error {
	sql := `UPDATE sys_user SET password=?, displayname=?, email=?, accountstatus=?, role_id=?,station_id=?,qn_account=?,qn_password=?,place_role=?,account_type=?
			WHERE id=?`
	_, err := orm.NewOrm().Raw(sql, u.Password, u.Displayname, u.Email, u.Accountstatus, u.Role_id, u.Station_id, u.Qn_account, u.Qn_password, u.Place_role, u.AccountType, u.Id).Exec()
	return err
}

func DeleteUser(uid int) error {
	sql := `DELETE FROM sys_user WHERE id=?`
	_, err := orm.NewOrm().Raw(sql, uid).Exec()
	return err
}

//修改用户密码
func UpdatePassword(id int, password, password1 string) (bool, error) {
	o := orm.NewOrm()
	sql := `update sys_user set password=? where id=?`
	_, err := o.Raw(sql, password, id).Exec()
	if err != nil {
		return false, err
	} else {
		return true, nil
	}
}

func UpdateLastOperationTime(id int, last time.Time) (err error) {
	sql := `UPDATE sys_user SET last_operation_time=?  where id=?`
	_, err = orm.NewOrm().Raw(sql, last, id).Exec()
	return
}

func GetUserType(stationId int) (userType int, err error) {
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
func GetIsUseByIp(ip string) (flag bool) {
	sql := `SELECT count(1)
           FROM ip_white_list
           WHERE ip_address=? AND is_use =0`
	var count int
	orm.NewOrm().Raw(sql, ip).QueryRow(&count)
	if count > 0 {
		return true
	}
	return false
}

//更新系统用户状态
func UpdateLoginStateById(id int, loginState string) error {
	_, err := orm.NewOrm().Raw(`UPDATE sys_user SET login_state = ? WHERE id = ?`, loginState, id).Exec()
	return err
}

//根据系统用户id查找用户信息
func FindSysUserById(id int) (sysUser SysUser, err error) {
	sql := `SELECT * FROM sys_user WHERE id = ?`
	err = orm.NewOrm().Raw(sql, id).QueryRow(&sysUser)
	return
}

//根据岗位id得到org_id
func GetOrgIdByStationId(sid int) int {
	o := orm.NewOrm()
	sql := ` SELECT org_id
			FROM sys_station ss
			WHERE  ss.id=?  `
	var orgId int
	o.Raw(sql, sid).QueryRow(&orgId)
	return orgId
}

//根据岗位ID查找角色ID
func QueryRoleIdByStationId(stationId int) (roleId int, err error) {
	sql := `SELECT role_id FROM sys_station WHERE id = ?`
	err = orm.NewOrm().Raw(sql, stationId).QueryRow(&roleId)
	return
}

//根据stationId更新系统用户roleId
func UpdateSysUserRoleId(stationId, roleId int) (err error) {
	o := orm.NewOrm()
	sql := `UPDATE sys_user SET role_id = ? WHERE station_id = ?`
	_, err = o.Raw(sql, roleId, stationId).Exec()
	return
}

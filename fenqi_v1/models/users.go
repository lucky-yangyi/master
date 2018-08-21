package models

import (
	"time"
	"zcm_tools/orm"
)

type Users struct {
	Id             int       `json "Uid" description:"用户ID"`
	Account        string    `description:"账号"`
	App            int       `description:"来自哪个平台 1 ios,2 android,3 wx,4pc，"`
	AppVersion     string    `description:"app版本"`
	MobileType     string    `description:"手机型号"`
	MobileVersion  string    `description:"手机版本号"`
	Token          string    `description:"登录成功后返回一个token,作为下次post 请求的"`
	IsEmulator     bool      `description:"是否是虚拟器"`
	DeviceUniqueID string    `description:"手机设备号"`
	FingerKey      string    `description:"白骑士设备指纹"`
	PkgType        int       `description:"分包"`
	DeviceSingle   string    `description:"设备唯一标示"`
	CodeType       string    `description:"登录注册验证码类型"`
	CreateTime     time.Time `description:"注册时间"`
	Source         string    `description:"渠道"`
	State          int       `description:"用户状态:0-正常 1-冻结"`
	RepSchId       int       `description:"本期应还的还款计划"`
	LoanId         int       `description:"借款ID"`
	Mtype          int       `description:"催收阶段"`
	OperationTime  time.Time `description:"操作时间"`
	SalesmanId     int
	IsSalemanType  string
}

func QueryUserByUid(uid int) (users Users, err error) {
	sql := `SELECT * FROM users WHERE id = ? `
	err = orm.NewOrm().Raw(sql, uid).QueryRow(&users)
	return
}

//查询用户注册时的信息
func QueryUsersInfo(uid int) (u *Users, err error) {
	o := orm.NewOrm()
	sql := `SELECT id,account,create_time,source,mobile_type,state FROM users WHERE id = ?`
	err = o.Raw(sql, uid).QueryRow(&u)
	return
}

//查询注册用户列表
func QueryUsersList(start, pageSize int, condition string, paras ...interface{}) (list []*Users, err error) {
	o := orm.NewOrm()
	sql := ` SELECT
				u.id,
				u.account,
				u.create_time,
				u.pkg_type,
				u.source,
				u.salesman_id
			FROM users AS u
			WHERE 1=1`
	if condition != "" {
		sql += condition
	}
	sql += ` ORDER BY u.create_time DESC LIMIT ?,?`
	_, err = o.Raw(sql, paras, start, pageSize).QueryRows(&list)
	return
}

//查询注册用户总数
func QueryUsersCount(condition string, paras ...interface{}) (count int, err error) {
	o := orm.NewOrm()
	sql := ` SELECT
				COUNT(1)
			FROM users AS u
			WHERE 1=1`
	if condition != "" {
		sql += condition
	}
	err = o.Raw(sql, paras).QueryRow(&count)
	return
}

func FindByAccountUser(phone string) (user Users, err error) {
	o := orm.NewOrm()
	sql := `select b.name as company,a.id as loan_id ,a.repayment_schedule_id as rep_sch_id,a.uid as id,datediff(now(),a.end_date) as overdue_days,a.org_id,a.mtype from loan as a left join sys_organization as b on a.org_id=b.id
             where a.account=? and a.state='BACKING' order by a.id desc limit 1`
	err = o.Raw(sql, phone).QueryRow(&user)
	return
}

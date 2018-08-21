package models

import (
	"time"
	"zcm_tools/orm"
)

type Advises struct {
	Id            int       `orm:"column(id);auto"`
	Uid           int       `orm:"column(uid);null"`                   //用户id
	Account       string    `orm:"column(account);size(11);null"`      //手机号
	AdviseType    string    `orm:"column(advise_type);size(20);null"`  //
	Content       string    `orm:"column(content);null"`               //反馈内容
	CreateTime    time.Time `orm:"column(create_time);type(datetime)"` //反馈时间
	IsChecked     int8      `orm:"column(is_checked);null"`
	Remarks       string    `orm:"column(remarks);null"`               //备注
	Dealpeople    string    `orm:"column(dealpeople);size(40);null"`   //处理人
	AppVersion    string    `orm:"column(app_version);size(255);null"` //App版本号
	App           int8      `orm:"column(app);null"`
	MobileVersion string    `orm:"column(mobile_version);size(100);null"`   //手机版本号
	MobileType    string    `orm:"column(mobile_type);size(100);null"`      //手机机型
	ModifyTime    time.Time `orm:"column(modify_time);type(datetime);null"` //修改时间
	Screenshot    string    `orm:"column(screenshot);size(1024);null"`      //反馈截图
	Status        bool
}

//获取意见反馈
func QueryAdviseList(start, pageSize int, condition string, paras ...interface{}) (list []Advises, err error) {
	o := orm.NewOrm()
	sql := `SELECT
				id,
				uid,
				account,
				create_time,
				is_checked,
				mobile_type,
				app_version,
				content,
				remarks,
				dealpeople,
				screenshot
			FROM advises
			WHERE 1=1`
	if condition != "" {
		sql += condition
	}
	sql += ` ORDER BY create_time DESC LIMIT ?,?`
	_, err = o.Raw(sql, paras, start, pageSize).QueryRows(&list)
	return
}

//获取意见反馈总数
func QueryAdviseCount(condition string, paras ...interface{}) (count int, err error) {
	o := orm.NewOrm()
	sql := `SELECT
				COUNT(1)
			FROM advises
			WHERE 1=1`
	if condition != "" {
		sql += condition
	}
	err = o.Raw(sql, paras).QueryRow(&count)
	return
}

//根据id获取一条反馈
func QueryAdviseById(id int) (advise Advises, err error) {
	o := orm.NewOrm()
	sql := `SELECT
				id,
				is_checked,
				remarks
			FROM advises
			WHERE id = ?`
	err = o.Raw(sql, id).QueryRow(&advise)
	return
}

//更新意见反馈备注
func UpdateAdviseRemark(remark, name string, isChecked, id int) error {
	sql := `UPDATE advises SET remarks=? ,dealpeople = ?,is_checked=?,modify_time = NOW() WHERE id=?`
	_, err := orm.NewOrm().Raw(sql, remark, name, isChecked, id).Exec()
	return err
}

type SalesAdvise struct {
	Id          int
	Account     string    `orm:"column(account);size(45)"`          //业务手机号
	Saleman     string    `orm:"column(saleman);size(45)"`          //业务员名字
	InviteCode  string    `orm:"column(invite_code);size(45);null"` //业务员邀请码
	CreateTime  time.Time `orm:"column(creat_time);null"`           //反馈时间
	Remark      string    `orm:"column(remark);size(256)"`          //备注
	Advise      string    `orm:"column(advise);size(256)"`          //备注
	Displayname string    `orm:"column(displayname)"`               //处理人
}

//获取业务员意见反馈
func QuerySalmanAdviseList(start, pageSize int, condition string, paras ...interface{}) (list []SalesAdvise, err error) {
	o := orm.NewOrm()
	sql := `SELECT
				a.id,
				a.advise,
				a.remark,
				a.displayname,
				b.saleman,
				a.creat_time,
				b.account,
				b.invite_code
			FROM
				salesman_advise AS a,
				salesman AS b
			WHERE
				a.account = b.account`
	if condition != "" {
		sql += condition
	}
	sql += ` ORDER BY a.creat_time DESC LIMIT ?,?`
	_, err = o.Raw(sql, paras, start, pageSize).QueryRows(&list)
	return
}

//获取业务员意见反馈
func QuerySalemanAdviseCount(condition string, paras ...interface{}) (count int, err error) {
	o := orm.NewOrm()
	sql := `SELECT
				COUNT(1)
			FROM
				salesman_advise AS a,
				salesman AS b
			WHERE
				a.account = b.account`
	if condition != "" {
		sql += condition
	}
	sql += ` ORDER BY a.creat_time`
	err = o.Raw(sql, paras).QueryRow(&count)
	return
}

//更新意见反馈备注
func UpdateSalemanAdviseRemark(remark, displayname string, id int) error {
	sql := `UPDATE salesman_advise SET remark=? ,displayname = ?,state=1 WHERE id=?`
	_, err := orm.NewOrm().Raw(sql, remark, displayname, id).Exec()
	return err
}

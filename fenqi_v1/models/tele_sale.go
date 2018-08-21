package models

import (
	"time"
	"zcm_tools/orm"
	"fmt"
)

//电销列表
type TeleSale struct {
	Id          int       `orm:"column(id)"`               //主键
	Phone       string    `orm:"column(phone)"`            //电话号码
	Name        string    `orm:"column(name)"`             //用户姓名
	PhoneStatus string    `orm:"column(phone_status)"`     //号码状态
	Remark      string    `orm:"column(remark)"`           //备注
	CreateTime  time.Time `orm:"column(create_time)"`      //创建时间
	EndCallTime time.Time `orm:"column(end_call_time)"`    //最后拨打时间
	CallMan     string    `orm:"column(call_man)"`         //电销员
	QnAccount   string    `orm:"column(qn_account)"`       //青牛账号
	QnPassword  string    `orm:"column(qn_password)"`      //青牛密码
	TeleId      int       `orm:"column(tele_id)"`          //电销员
	State       int       `orm:"column(final_auth_state)"` //授信状态
	Uid         int       `orm:"column(uid)"`              //用户id
	CreditSate  string    `orm:"column(state)"`            //排队中
	PlatFrom  string    `orm:"column(platfrom)"`            //排队中
}

type TeleMan struct {
	Phone string
	Name  string
}

//获取电销员列表
func TeleSaleList(start, pageSize int, condition string, paras ...interface{}) (v []TeleSale, err error) {
	sql := `SELECT
				b.*,c.final_auth_state,
				a.id AS uid,
				d.qn_account,
				d.qn_password,
				f.state
			FROM
				tele_sale AS b
			LEFT JOIN users AS a ON a.account = b.phone
			LEFT JOIN users_auth AS c ON a.id = c.uid AND c.is_valid = 1
			LEFT JOIN credit_aduit as f ON a.id= f.uid AND f.is_now = 1
			INNER JOIN sys_user AS d ON d.id = b.tele_id WHERE 1=1 `
	if condition != "" {
		sql += condition
	}
	sql += ` ORDER BY b.modify_time LIMIT ?,?`
	fmt.Println(sql)
	_, err = orm.NewOrm().Raw(sql, paras, start, pageSize).QueryRows(&v)
	//fmt.Println(start,paras)
	return
}

//获取电销列表总数
func TeleSaleListCount(condition string, paras ...interface{}) (count int, err error) {
	sql := `SELECT
				COUNT(1)
			FROM
				tele_sale AS b
			LEFT JOIN users AS a ON a.account = b.phone
			LEFT JOIN users_auth AS c ON a.id = c.uid AND c.is_valid = 1
			LEFT JOIN credit_aduit as f ON a.id= f.uid AND f.is_now = 1
			INNER JOIN sys_user AS d ON d.id = b.tele_id WHERE 1=1 `
	if condition != "" {
		sql += condition
	}
	err = orm.NewOrm().Raw(sql, paras).QueryRow(&count)
	return
}

//分配权限
func TeleSaleAllomatPower(tele_id int) (count int, err error) {
	sql := `SELECT
			COUNT(1)
		FROM
			sys_user AS a,
			sys_role_menu AS b
		WHERE
			a.role_id = b.role_id
		AND a.id = ?
		AND b.menu_id IN (87, 88)`
	err = orm.NewOrm().Raw(sql, tele_id).QueryRow(&count)
	return
}

type AllomatUser struct{
	Id          int            `orm:"column(id)"`
	Displayname string  `orm:"column(displayname)"`
}

//获取电销人员列表
func TeleSaleListDisplayname() (v []AllomatUser, err error) {
	o := orm.NewOrm()
	sql := `SELECT
				a.displayname,
				a.id
			FROM
				sys_user AS a,
				sys_station_type AS b
			WHERE
				a.station_id = b.station_id
			AND b.type = 17`
	_, err = o.Raw(sql).QueryRows(&v)
	return
}

//更新状态
func UpdateTeleSaleStatus(phone_status, remark string, id int) error {
	sql := `UPDATE tele_sale SET phone_status=?,remark=?,modify_time=NOW() WHERE id =? `
	_, err := orm.NewOrm().Raw(sql, phone_status, remark, id).Exec()
	return err
}

//查询未分配数量
func NoAllomentTeleSaleListCount() (count int, err error) {
	o := orm.NewOrm()
	sql := `SELECT
				COUNT(1)
			FROM
				tele_sale AS b
			WHERE
			b.phone_status="no_allotment"`
	err = o.Raw(sql).QueryRow(&count)
	return
}

//更新分配
func UpdateNoAllomentTeleSale(number, tele_id int, call_man string) (err error) {
	var id []int
	sql := `SELECT
				id
			FROM
				tele_sale AS b
			WHERE
				b.phone_status = "no_allotment"
			ORDER BY
				b.create_time DESC
			LIMIT ?`
	_, err = orm.NewOrm().Raw(sql, number).QueryRows(&id)
	if err != nil {
		return err
	}
	if len(id) > 0 {
		for _, v := range id {
			sql := `UPDATE tele_sale
						SET phone_status = "no_call",
						 	call_man =?,
						 	tele_id =?,
							create_time=NOW()
						WHERE
							id =? `
			_, err := orm.NewOrm().Raw(sql, call_man, tele_id, v).Exec()
			if err != nil {
				return err
			}
		}
	}
	return
}

//添加电销客户
func AddTeleSaleFromExcel(teleman []TeleMan, tele_id int) (err error) {
	sql_t := `INSERT INTO tele_sale (phone,name,tele_id,phone_status) VALUES(?,?,?,"")`
	p, err := orm.NewOrm().Raw(sql_t).Prepare()

	sql_t1 := `INSERT INTO tele_sale (phone,name,tele_id) VALUES(?,?,?)`
	q, err := orm.NewOrm().Raw(sql_t1).Prepare()

	for _, v := range teleman {
		if v.Phone != "" {
			var count int
			sql_t := `SELECT COUNT(1) FROM users WHERE account=?`
			orm.NewOrm().Raw(sql_t,v.Phone).QueryRow(&count)
			if count >0{
				_, err = p.Exec(v.Phone, v.Name, tele_id)
			}else{
				_, err = q.Exec(v.Phone, v.Name, tele_id)
			}
		}
	}
	p.Close()
	q.Close()
	return
}

//更新状态
func UpdateTeleEndCallTime(phone string, id int) error {
	sql := `UPDATE tele_sale SET end_call_time=NOW() WHERE phone =? AND id=?`
	_, err := orm.NewOrm().Raw(sql, phone, id).Exec()
	return err
}

// //更新状态
// func UpdateTeleEndCallTime(id int) error {
// 	sql := `UPDATE tele_sale SET end_call_time=NOW() WHERE id =? `
// 	_, err := orm.NewOrm().Raw(sql, id).Exec()
// 	return err
// }

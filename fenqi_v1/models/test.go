package models

import (
	"time"
	"zcm_tools/orm"
)

//测试管理
type Test struct {
	Id          int       `orm:"column(id)"`               //主键ID
	Phone       string    `orm:"column(phone)"`            //电话号码
	Name        string    `orm:"column(name)"`             //用户姓名
	PhoneStatus string    `orm:"column(phone_status)"`      //号码状态
	Remark      string    `orm:"column(remark)"`            //备注
	CreateTime  time.Time  `orm:"column(create_time)"`      //创建时间
	EndCallTime time.Time  `orm:"column(end_call_time)"`    //最后拨打时间
}
 func TestList(start, pageSize int, condition string, paras ...interface{})(list []Test,err error){
	 o := orm.NewOrm()
 	sql := `SELECT * FROM tele_sale WHERE 1=1`
	 if condition != "" {
		sql += condition
	 }
	 sql += ` ORDER BY modify_time DESC LIMIT ?,?`
	 _, err = o.Raw(sql,paras,start,pageSize).QueryRows(&list)
	 //fmt.Println(sql)
	 return
 }
 //获取测试管理总数
 func QueryTestCount(condition string,paras ...interface{})(count int,err error){
 	o := orm.NewOrm()
 	sql := `SELECT
                COUNT(1)
            FROM tele_sale
            WHERE 1=1`
            if condition !=""{
            	sql +=condition
			}
			err = o.Raw(sql, paras).QueryRow(&count)
			return
 	}
 	//更新状态
 	func UpdateStatus(phone_status,remark string, id int) error {
 		   o :=orm.NewOrm()
            sql := `UPDATE tele_sale SET phone_status=?,remark=?,modify_time=NOW() WHERE id=?`
            _,err := o.Raw(sql,phone_status,remark,id).Exec()
            return err
            }

	//添加电销员
	func AddTestExcel(teleman []TeleMan,tele_id int)(err error){
		//fmt.Println("===========tale_id--",tale_id)
		o :=orm.NewOrm()
		sql_t := `INSERT INTO tele_sale (phone,name,tele_id,phone_status) VALUES(?,?,?,"")`
		p,err := o.Raw(sql_t).Prepare()
		sql_t1 := `INSERT INTO tele_sale (phone,name,tele_id) VALUES(?,?,?)`
		q,err := o.Raw(sql_t1).Prepare()
		//fmt.Println("--err4--",err)
		for _, v := range teleman{
			if v.Phone != "" {
				var count int
				sql_t := `SELECT COUNT(1) FROM users WHERE account=?`
				orm.NewOrm().Raw(sql_t,v.Phone).QueryRow(&count)
				if count > 0{
					_, err = p.Exec(v.Phone,v.Name,tele_id)
				}else{
					_, err = q.Exec(v.Phone,v.Name,tele_id)
				}
			}
		}
		p.Close()
		q.Close()
		//fmt.Println("--err3--",err)
		return
	}




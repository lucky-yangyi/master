package models

import (
	"time"
	"zcm_tools/orm"
)

type ConnectRecord struct {
	Id          int       `orm:"column(id);pk" `
	Uid         int       `orm:"column(uid)" description:"用户ID"`
	CreateTime  time.Time `orm:"column(create_time);type(datetime);null" description:"创建时间"`
	Operator    string    `orm:"column(operator)" description:" 操作人"`
	Context     string    `orm:"column(context);size(255);null" description:"联系内容"`
	ConnectType string    `orm:"column(connect_type);null" description:"联系类型  COMPLAIN:投诉处理  ENTER:  客服登记   TOUCH : 信审联系"`
	ConnectObj  string    `orm:"column(connect_obj);null" description:"联系对象  客户本人  紧急联系人1  紧急联系人2  其他(自填)"`
}

// 添加联系历史
func AddConnectRecord(m *ConnectRecord) (id int64, err error) {
	o := orm.NewOrm()
	id, err = o.Insert(m)
	return
}

// 获取联系历史根据uid
func GetConnectRecordByUid(uid int, connectType string, start, pageSize int) (v []ConnectRecord, err error) {
	o := orm.NewOrm()
	sql := `SELECT uid,create_time,operator,context,connect_type,connect_obj FROM connect_record WHERE uid = ? AND connect_type = ? ORDER BY create_time DESC LIMIT ? ,?`
	_, err = o.Raw(sql, uid, connectType, start, pageSize).QueryRows(&v)
	return v, err
}

//	联系历史count
func GetConnectRecordCountByUid(uid int, connectType string) (count int, err error) {
	o := orm.NewOrm()
	sql := `SELECT COUNT(1) FROM connect_record WHERE uid = ? AND connect_type = ?`
	err = o.Raw(sql, uid, connectType).QueryRow(&count)
	return
}

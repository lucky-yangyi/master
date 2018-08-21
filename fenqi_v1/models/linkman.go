package models

import (
	"time"

	"zcm_tools/orm"

	"gopkg.in/mgo.v2/bson"
)

//通讯录
type MailList struct {
	Id         bson.ObjectId `bson:"_id"`
	Uid        int
	Contact    []Contact
	LinkType   int
	CreateTime time.Time
}

type Contact struct {
	ContactName        string   `description:"联系姓名"`
	ContactPhoneNumber []string `description:"联系号码"`
	SignState          int      `description:"标记状态"`
	SignRemark         string   `description:"标记备注"`
	SignId             int      `description:"标记ID"`
}

type ContactData struct {
	ContactName        string `description:"联系姓名"`
	ContactPhoneNumber string `description:"联系号码"`
}

type MailList2 struct {
	Id         bson.ObjectId `bson:"_id"`
	Uid        int           `bson:"uid"`
	Contact    []Contact     `bson:"contact"`
	LinkType   int
	CreateTime time.Time
}

//用户紧急联系人
type UsersLinkman struct {
	Id                 int    `orm:"column(id);auto"`
	Uid                int    `description:"用户ID"`
	Relation           string `description:"关系"`
	LinkmanName        string `description:"联系人姓名"`
	ContactPhoneNumber string `description:"紧急联系人电话"`
	IsNormal           int8   `description:"紧急联系人电话是否异常:1:正常，2:异常"`
	AbnormalResult     string `description:"异常原因"`
	CreateTime         time.Time
	Source             int `description:"1伊顿"`
	SignRemark         string
}

//用户通讯录联系人---注意区分和紧急联系人
type UsersSimpleLinkman struct {
	Id          int
	Uid         int
	ContactName string //通讯录联系人姓名
	CreateTime  time.Time
}

//用户通讯录号码
type LinkmanPhone struct {
	Id          int
	LinkmanId   int
	PhoneNumber string
}

type SmsRecordsOneHour struct {
	Uid          int
	PlatformName string
	ProductId    int
	SmsRecords   []SmsRecordsItem
}

type SmsRecordsItem struct {
	PhoneNumber string
	SmsContent  string
	Date        string
	Type        int
}

type Thesaurus struct {
	Name string
	Id   int
}

func QueryUsersLinkmanList(uid int) (list []UsersLinkman, err error) {
	o := orm.NewOrm()
	sql := `select * from users_linkman where uid=? order by create_time desc limit 2`
	_, err = o.Raw(sql, uid).QueryRows(&list)
	return
}

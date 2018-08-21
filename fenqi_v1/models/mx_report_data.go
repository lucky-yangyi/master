package models

import "gopkg.in/mgo.v2/bson"

type Mxreportdata struct {
	Id bson.ObjectId `bson:"_id"`
	Rt Report        `bson:"report"`
}

type Report struct {
	Ccd []CallContactDetail `bson:"call_contact_detail"`
}

type CallContactDetail struct {
	Call_Cnt_1w int    `bson:"call_cnt_1w"`
	Call_Cnt_1m int    `bson:"call_cnt_1m"`
	Call_Cnt_6m int    `bson:"call_cnt_6m"`
	Peer_num    string `bson:"peer_num"`
	City        string `bson:"city"`
	SignState   int    `description:"标记状态"`
	SignRemark  string `description:"标记备注"`
	SignId      int    `description:"标记ID"`
}

//近3个月通话记录统计
type MxrThreeReportData struct {
	Id bson.ObjectId `bson:"_id"`
	Rt ReportData    `bson:"report"`
}
type ReportData struct {
	Ccd []Call_Contact_Detail `bson:"call_contact_detail"`
}

type Call_Contact_Detail struct {
	Peer_num      string `bson:"peer_num"`      //手机号
	City          string `bson:"city"`          //号码归属地
	Dial_Cnt_3m   int    `bson:"dial_cnt_3m"`   //主叫次数
	Dialed_Cnt_3m int    `bson:"dialed_cnt_3m"` //被叫次数
	Call_Time_3m  int    `bson:"call_time_3m"`  //通话时长
	Total_Cnt     int    //联系次数
	ContactName   string //联系人
}

type MatchCall struct {
	ContactName        string //联系人
	ContactPhoneNumber string //手机号
	ConnectCount       int    //联系次数
	MonType            string //运营商
	ConnectFlag        bool   //联系人记录flag
	MobileFlag         bool   //通讯录记录flag
}

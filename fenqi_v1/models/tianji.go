package models

import (
	"gopkg.in/mgo.v2/bson"
)

type Tijireport struct {
	Id     bson.ObjectId                           `bson:"_id"`
	Tianji Tianji_api_tianjireport_detail_response `bson:"tianji_api_tianjireport_detail_response"`
}

type Tianji_api_tianjireport_detail_response struct {
	CallLog []CallLogInfo `bson:"call_log"` //通讯记录
}

type CallLogInfo struct {
	Phone                 string `bson:"phone"`                 //号码	string
	Phone_location        string `bson:"phone_location"`        //号码归属地	string
	Phone_info            string `bson:"phone_info"`            //互联网标识	string
	Phone_label           string `bson:"phone_label"`           //类别标签	string
	First_contact_date    string `bson:"first_contact_date"`    //首次联系时间	string
	Last_contact_date     string `bson:"last_contact_date"`     //最后联系时间	string
	Talk_seconds          int    `bson:"talk_seconds"`          //通话时长	int
	Talk_cnt              int    `bson:"talk_cnt"`              //通话次数	int
	Call_seconds          int    `bson:"call_seconds"`          //主叫时长	int
	Call_cnt              int    `bson:"call_cnt"`              //主叫次数	int
	Called_seconds        int    `bson:"called_seconds"`        //被叫时长	int
	Called_cnt            int    `bson:"called_cnt"`            //被叫次数	int
	Msg_cnt               int    `bson:"msg_cnt"`               //短信总数	int
	Send_cnt              int    `bson:"send_cnt"`              //发送短信数	int
	Receive_cnt           int    `bson:"receive_cnt"`           //接收短信数	int
	Contact_1w            int    `bson:"contact_1w"`            //近一周联系次数	int
	Contact_1m            int    `bson:"contact_1m"`            //近一个月联系次数	int
	Contact_3m            int    `bson:"contact_3m"`            //近三个月联系次数	int
	Contact_early_morning int    `bson:"contact_early_morning"` //凌晨联系次数	int
	Contact_morning       int    `bson:"contact_morning"`       //早晨联系次数	int
	Contact_noon          int    `bson:"contact_noon"`          //上午联系次数	int
	Contact_afternoon     int    `bson:"contact_afternoon"`     //下午联系次数	int
	Contact_evening       int    `bson:"contact_evening"`       //夜晚联系次数	int
	Contact_night         int    `bson:"contact_night"`         //深夜联系次数	int
	Contact_weekday       int    `bson:"contact_weekday"`       //工作日联系次数	int
	Contact_weekend       int    `bson:"contact_weekend"`       //周末联系次数	int
	ContactName           string //联系人
}

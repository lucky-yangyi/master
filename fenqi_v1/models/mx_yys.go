package models

import "time"

//运营商数据
type MdbMxYYSData struct {
	Uid        int       `bson:"uid"`
	MXYYSData  MXYYSData `bson:"mxyysdata"`
	CreateTime time.Time `bson:"createtime"`
}

type MXYYSData struct {
	Calls []Calls `bson:"calls"`
}

type Calls struct {
	BillMonth string  `bson:"bill_month"`
	Items     []Items `bson:"items"`
}

type Items struct {
	PeerNumber string  `bson:"peer_number"`
	Duration   float64 `bson:"duration"`
	DialType   string  `bson:"dial_type"`
	Time       string  `bson:"time"`
}

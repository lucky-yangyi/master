package models

import (
	"fenqi_v1/utils"

	"strconv"

	"github.com/astaxie/beego"
)

type RiskAdvise struct {
	Uid         string      `bson:"uId"`
	OrderId     string      `bson:"orderId"`
	AuditAdvice string      `bson:"auditAdvice"`
	EndTime     string      `bson:"endTime"`
	LimitMoney  int         `bson:"limitMoney"`
	RiseItems   []*RiskItem `bson:"riskItems"`
}
type RiskItem struct {
	RiskName   string `bson:"riskName"`
	RiskAdvise string `bson:"riskAdvice"`
	State      int
}

func GetRiskAdvises(uid, businessType int) (v *RiskAdvise, err error) {
	session := utils.GetYeadunSession()
	paramMap := make(map[string]interface{})
	paramMap["uId"] = strconv.Itoa(uid)
	paramMap["businessType"] = businessType
	defer session.Close()
	err = session.DB(utils.MGO_YEADUN_DB).C("result_data").Find(paramMap).Sort("-startTime").One(&v)
	beego.Info(v)
	return
}

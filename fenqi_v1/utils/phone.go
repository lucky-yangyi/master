package utils

import (
	"encoding/json"
	"github.com/astaxie/beego/httplib"
)

type Phone struct {
	Ret   int    `json:"ret"`
	Msg   string `json:"msg"`
	LogID string `json:"log_id"`
	Data  struct {
		Types    string `json:"types"`
		Lng      string `json:"lng"`
		City     string `json:"city"`
		Num      int    `json:"num"`
		Isp      string `json:"isp"`
		AreaCode string `json:"area_code"`
		CityCode string `json:"city_code"`
		Prov     string `json:"prov"`
		ZipCode  string `json:"zip_code"`
		Lat      string `json:"lat"`
	} `json:"data"`
}

//手机号归属地查询
func QueryLocating(phone string) (phone_info string) {
	var m Phone
	//缓存
	if Rc.IsExist("xjfq_phone:" + PhoneInfo + phone) {
		if data, err := Rc.RedisBytes("xjfq_phone:" + PhoneInfo + phone); err == nil {
			json.Unmarshal(data, &m)
			phone_info = m.Data.Prov + " " + m.Data.City + " " + m.Data.Types
			return phone_info
		}
	} else {
		//第三方链接
		res, err := httplib.Get("http://api04.aliyun.venuscn.com/mobile?mobile="+phone).Header("Authorization", "APPCODE 89cec72da6c049308689da99ed33d0d5").Bytes()
		if err != nil {
			return
		}
		if err := json.Unmarshal(res, &m); err == nil {
			phone_info = m.Data.Prov + " " + m.Data.City + " " + m.Data.Types
			Rc.Put("xjfq_phone:"+PhoneInfo+phone, res, RedisCacheTime_Year)
		}
		//, _ := json.Marshal(res)
		return phone_info
	}
	return
}

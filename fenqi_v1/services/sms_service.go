package services

import (
	"encoding/json"
	"fenqi_v1/utils"
	"fmt"
	"net/url"
	"zcm_tools/http"
)

type SmsResult struct {
	Ret int
	Msg string
}

func SendSms(account, content, source, sendtype, ip, category string) bool {

	if len(account) != 11 {
		return false
	}
	params := url.Values{}
	params.Set("merchant", utils.Merchant)
	params.Set("account", account)
	params.Set("content", content)
	params.Set("source", source)
	params.Set("ip", ip)
	params.Set("category", category)
	b, err := http.Post(utils.SMSURL, params.Encode()) //http://127.0.0.1:8104/v1/sms/sendsms  http://xjser.zcmlc.com/v1/sms/sendsms
	if err != nil {
		fmt.Println("2", err.Error())
		return false
	} else {
		var m SmsResult
		if err := json.Unmarshal(b, &m); err == nil && m.Ret == 200 {
			return true
		} else {
			return false
		}
	}
	return true
}

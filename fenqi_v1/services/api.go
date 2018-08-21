package services

import (
	"encoding/json"
	"fenqi_v1/models"
	"fenqi_v1/utils"
	"strconv"
	"strings"
	"zcm_tools/http"

	"github.com/astaxie/beego"
)

//post 方式请求api
func PostApi(method string, params map[string]interface{}) ([]byte, error) {
	params["MobileType"] = "YGFQ_SYS"
	content, err := json.Marshal(params)
	if err != nil {
		return nil, err
	}
	// data := utils.DesBase64Encrypt(content)
	b, err := http.Post(utils.FQ_API_URL+method, string(content))
	if err != nil {
		return nil, err
	}
	return b, nil
}

func GetApi(method string) ([]byte, error) {
	b, err := http.Get(utils.FQ_API_URL + method)
	if err != nil {
		return nil, err
	}
	return b, nil
}

func ShareEqual(ids []int, num int, list []*models.SalesmanAllotment) (shares []models.Share) {
	if len(ids) > 0 && num > 0 {
		a := num % len(ids)
		b := (num - a) / len(ids)
		var share models.Share
		for k, _ := range ids {
			share.Id = ids[k]
			if a > 0 {
				share.New = b + 1
				a = a - 1
			} else {
				share.New = b
			}
			i := 1
			beego.Info(share)
			e := ""
			for _, v := range list {
				if i <= share.New {
					beego.Info(v)
					if v.State == 0 {
						e += "," + strconv.Itoa(v.Id)
						v.State = 1
						i++
					}
				}
			}
			share.AllotmentId = strings.Trim(e, ",")
			sale, _ := models.GetSalemanById(ids[k])
			share.All = sale.AllAllotment + share.New
			shares = append(shares, share)

			beego.Info(shares)
		}
	}
	return
}

package models

import (
	"encoding/json"
	"fenqi_v1/utils"
	"time"
	"zcm_tools/orm"
)

type Config struct {
	Id          int
	ConfigKey   string `description:"键名"`
	ConfigValue string `description:"键值"`
	ConfigDesc  string `description:"内容"`
	Remark      string `description:"备注"`
	ConfigUrl   string `description:"键连接"`
	UrlParam    int    `description:"url是否加参数,0:不带参数,1带参数"`
	Title       string `description:"跳转页面的标题"`
	ShowVersion string `description:"更新版本显示"`
}

func GetConfigCache(code string) (cf *Config, err error) {
	o := orm.NewOrm()
	sql := `SELECT
				id,
				config_key,
				config_value,
				config_desc,
				remark,
				config_url,
				url_param,
				title,
				show_version
			FROM config
			WHERE config_key = ? `
	err = o.Raw(sql, code).QueryRow(&cf)
	if data, err2 := json.Marshal(cf); err2 == nil && utils.Re == nil {
		utils.Rc.Put(utils.FQ_CACHE_KEY_CONFIG+code, data, 24*time.Hour)
	}
	return cf, err
}

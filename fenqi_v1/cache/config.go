package cache

import (
	"encoding/json"
	"fenqi_v1/models"
	"fenqi_v1/utils"
)

//获取配置信息
func GetConfigByCKeyCache(code string) (m *models.Config, err error) {
	if utils.Re == nil && code != "" && utils.Rc.IsExist(utils.FQ_CACHE_KEY_CONFIG+code) {
		if data, err1 := utils.Rc.RedisBytes(utils.FQ_CACHE_KEY_CONFIG + code); err1 == nil {
			err = json.Unmarshal(data, &m)
			if m != nil {
				return
			}
		}
	}
	return models.GetConfigCache(code)
}

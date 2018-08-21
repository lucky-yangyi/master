package cache

import (
	"encoding/json"
	"fenqi_v1/models"
	"fenqi_v1/utils"
	"strconv"
)

//岗位所拥有的菜单权限-树结构
func GetSysMenuTreeByRoleId(role_id int) (m models.SysMenuList, err error) {
	if utils.Re == nil && utils.Rc.IsExist(utils.CacheKeyRoleMenuTreePrefix+strconv.Itoa(role_id)) {
		if data, err := utils.Rc.RedisBytes(utils.CacheKeyRoleMenuTreePrefix + strconv.Itoa(role_id)); err == nil {
			err = json.Unmarshal(data, &m)
			if m != nil {
				return m, err
			}
		}
	}
	return
}

//岗位所拥有的菜单权限-树结构
func GetSysMenuTreeMapByRoleId(role_id int) (m map[int]int, err error) {
	if utils.Re == nil && utils.Rc.IsExist(utils.CacheKeyRoleMenuMapTreePrefix+strconv.Itoa(role_id)) {
		if data, err := utils.Rc.RedisBytes(utils.CacheKeyRoleMenuMapTreePrefix + strconv.Itoa(role_id)); err == nil {
			err = json.Unmarshal(data, &m)
			if m != nil {
				return m, err
			}
		}
	}
	return
}

//组织架构信息
func GetSysOrganization() ([]models.SysOrganization, error) {
	if utils.Re == nil && utils.Rc.IsExist(utils.CacheKeySystemOrganization) {
		if data, err := utils.Rc.RedisBytes(utils.CacheKeySystemOrganization); err == nil {
			var m []models.SysOrganization
			err = json.Unmarshal(data, &m)
			if m != nil {
				return m, err
			}
		}
	}
	return models.GetSysOrganization()
}

//获取菜单信息
func GetCacheSysMenu() (m map[string]models.SysMenu, err error) {
	if utils.Re == nil && utils.Rc.IsExist(utils.CacheKeySystemMenu) {
		if data, err1 := utils.Rc.RedisBytes(utils.CacheKeySystemMenu); err1 == nil {
			err = json.Unmarshal(data, &m)
			if m != nil {
				return
			}
		}
	}
	m, err = models.GetSysMenu()
	return
}

//权限
func GetCacheDataByStation(stationId int) (str string, err error) {
	if utils.Re == nil && utils.Rc.IsExist(utils.CacheKeySysStationData+strconv.Itoa(stationId)) {
		if data, err := utils.Rc.RedisBytes(utils.CacheKeySysStationData + strconv.Itoa(stationId)); err == nil {
			str = string(data)
			if str != "" {
				return str, err
			}
		}
	}
	return models.GetDataByStation(stationId)
}

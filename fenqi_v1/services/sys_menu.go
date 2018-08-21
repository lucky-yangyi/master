package services

import (
	"encoding/json"
	"fenqi_v1/cache"
	"fenqi_v1/models"
	"fenqi_v1/utils"
	"sort"
	"strconv"
)

//获取用户菜单
func GetSysMenuTreeByRoleId(role_id int) (models.SysMenuList, error) {
	//t, _ := cache.GetSysMenuTreeByRoleId(role_id)
	// if t != nil {
	// 	return t, nil
	// }
	//m, err := models.GetSysMenuTreeByRoleId(role_id)
	m, err := models.GetSysMenuTreeByStationId(role_id) //根据岗位ID,获取菜单信息
	l := len(m)
	var menu models.SysMenuList
	for i := 0; i < l; i++ {
		if m[i].ParentId == 0 {
			for j := 0; j < l; j++ {
				if m[j].ParentId == m[i].Id {
					m[i].ChildMenu = append(m[i].ChildMenu, &m[j])
				}
			}
			sort.Sort(m[i].ChildMenu)
			menu = append(menu, &m[i])
		}
	}
	sort.Sort(menu)
	// if data, err2 := json.Marshal(menu); err == nil && err2 == nil && utils.Re == nil {
	// 	utils.Rc.Put(utils.CacheKeyRoleMenuTreePrefix+strconv.Itoa(role_id), data, utils.RedisCacheTime_Role)
	// }\
	return menu, err
}
func GetSysMenuByRoleId(role_id int) (map[int]int, error) {

	t, _ := cache.GetSysMenuTreeMapByRoleId(role_id)
	if t != nil {
		return t, nil
	}
	var menu map[int]int = map[int]int{}
	m1, err := models.GetSysMenuTreeByStationId(role_id) //根据岗位ID,获取菜单信息
	for _, k := range m1 {
		menu[k.Id] = 1
	}
	if data, err2 := json.Marshal(menu); err == nil && err2 == nil && utils.Re == nil {
		utils.Rc.Put(utils.CacheKeyRoleMenuMapTreePrefix+strconv.Itoa(role_id), data, utils.RedisCacheTime_User)
	}
	return menu, err
}

func GetSysMenuZTree(list []models.SysMenu) []map[string]interface{} {
	var menu []map[string]interface{}
	l := len(list)
	if l == 0 {
		return []map[string]interface{}{}
	}
	for i := 0; i < l; i++ {
		menu = append(menu, map[string]interface{}{"id": list[i].Id, "pId": list[i].ParentId, "open": true, "name": list[i].DisplayName})
	}
	return menu
}

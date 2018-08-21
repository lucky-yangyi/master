package services

import (
	"fenqi_v1/models"
)

//返回ztree格式组织架构数据
func GetCollectSysOrganizationZTree() ([]map[string]interface{}, error) {
	list, err := models.GetCollectSysOrganization()
	l := len(list)
	var org []map[string]interface{}
	for i := 0; i < l; i++ {
		org = append(org, map[string]interface{}{"id": list[i].Id, "pId": list[i].ParentId, "open": true, "name": list[i].Name, "remark": list[i].Remark})
	}
	//beego.Info(utils.ToString(list))
	return org, err
}

//返回ztree格式菜单数据
func GetCollectAllSysMenuZTree() ([]map[string]interface{}, error) {
	list, err := models.GetCollectSysMenuTreeAll()
	var menu []map[string]interface{}
	l := len(list)
	for i := 0; i < l; i++ {
		menu = append(menu, map[string]interface{}{"id": list[i].Id, "pId": list[i].ParentId, "open": true, "name": list[i].DisplayName})
	}
	return menu, err
}

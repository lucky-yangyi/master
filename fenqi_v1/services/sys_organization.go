package services

import (
	"fenqi_v1/cache"
	"fenqi_v1/models"
	"strconv"
)

//返回ztree格式组织架构数据
func GetSysOrganizationZTree() ([]map[string]interface{}, error) {
	list, err := cache.GetSysOrganization()
	l := len(list)
	var org []map[string]interface{}
	for i := 0; i < l; i++ {
		org = append(org, map[string]interface{}{"id": list[i].Id, "pId": list[i].ParentId, "open": true, "name": list[i].Name, "remark": list[i].Remark})
	}
	return org, err
}

//返回ztree格式菜单数据
func GetAllSysMenuZTree() ([]map[string]interface{}, error) {
	list, err := models.GetSysMenuTreeAll()
	var menu []map[string]interface{}
	l := len(list)
	for i := 0; i < l; i++ {
		menu = append(menu, map[string]interface{}{"id": list[i].Id, "pId": list[i].ParentId, "open": true, "name": list[i].DisplayName})
	}
	return menu, err
}

//获取部门下面的所有子部门
func GetSubSysOrganizationCode(orgcode string) (string, error) {
	// orgcode = "22de5463-80d4-43b2-9350-c9c5ea1e7f25"
	list, err := cache.GetSysOrganization()
	l := len(list)
	orgs := "'" + orgcode + "'"
	if l > 0 {
		for i := 0; i < l; i++ {
			if strconv.Itoa(list[i].ParentId) == orgcode {
				sub := GetChildOrganizationCode(list, strconv.Itoa(list[i].Id)) //, ch
				orgs = orgs + ",'" + strconv.Itoa(list[i].Id) + "'" + sub
			}
		}
	}
	//beego.Info(utils.ToString(org))
	return orgs, err
}

func GetChildOrganizationCode(list []models.SysOrganization, parentcode string) string {
	l := len(list)
	orgs := ""
	for i := 0; i < l; i++ {
		if strconv.Itoa(list[i].ParentId) == parentcode {
			sub := GetChildOrganizationCode(list, strconv.Itoa(list[i].Id)) //, ch
			orgs = orgs + ",'" + strconv.Itoa(list[i].Id) + "'" + sub
		}
	}
	return orgs
}

//查看当前单子所在营业部
//返回 上级部门名称，营业部名称，营业部code
func FindOrgTypeDepartment(org map[string]models.SysOrganization, orgcode string) (string, string, string, string) {
	if o, ok := org[orgcode]; ok {
		//判断是否营业部  organizationtype_group
		d := o.Name
		c := ""
		e := ""
		dcode := strconv.Itoa(o.Id)
		if p, ok1 := org[strconv.Itoa(o.ParentId)]; ok1 {
			c = p.Name
			if q, ok2 := org[strconv.Itoa(p.ParentId)]; ok2 {
				e = q.Name
			}
		}
		return c, d, dcode, e

	}
	return "", "", "", ""
}

//返回 上级部门名称，营业部名称，营业部code
func FindOrgTypePlace(org map[string]models.SysOrganization, orgcode string) (string, string, string) {
	if o, ok := org[orgcode]; ok {
		//判断是否营业部  organizationtype_group
		d := o.Name
		c := ""
		dcode := strconv.Itoa(o.Id)
		if p, ok1 := org[strconv.Itoa(o.ParentId)]; ok1 {
			c = p.Name

		}
		return c, d, dcode

	}
	return "", "", ""
}

func QueryOrgTypeDepartment(org map[string]models.SysOrganization, dep_id, place_id, region_id int) (large_area_name, place_name, bus_dep_name string) {
	//dep_id 营业部id place_id 省级id region_id区运营中心id  营业部 有的省份无营业部
	//大区
	if large_area, ok := org[strconv.Itoa(region_id)]; ok {
		large_area_name = large_area.Name
		//省
		if place, ok1 := org[strconv.Itoa(place_id)]; ok1 {
			place_name = place.Name
			//营业部
			if bus_dep, ok := org[strconv.Itoa(dep_id)]; ok {
				bus_dep_name = bus_dep.Name
			}
		}
	}
	return large_area_name, place_name, bus_dep_name
}

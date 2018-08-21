package models

import (
	"encoding/json"
	"fenqi_v1/utils"
	"time"
	"zcm_tools/orm"
)

type SysMenuList []*SysMenu

func (list SysMenuList) Len() int {
	return len(list)
}

func (list SysMenuList) Less(i, j int) bool {
	return list[i].SortIndex < list[j].SortIndex
}

func (list SysMenuList) Swap(i, j int) {
	list[i], list[j] = list[j], list[i]
}

//系统菜单
type SysMenu struct {
	Id          int         `orm:"column(id);pk"`
	DisplayName string      `orm:"column(displayname);null"`
	ControlUrl  string      `orm:"column(controlurl);null"`
	HomeUrl     string      `orm:"column(homeurl);null"`
	ParentId    int         `orm:"column(parent_id);null"`
	SortIndex   int         `orm:"column(sortindex);null"`
	ChildMenu   SysMenuList `orm:"-"`
}

//按钮权限
type SysBtn struct {
	Id         int    `orm:"column(id);pk"`
	ControlUrl string `orm:"column(controlurl);null"`
	IsShow     bool
}

func GetSysMenuTreeByRoleId(role_id int) ([]SysMenu, error) {
	o := orm.NewOrm()
	sql := "SELECT DISTINCT d.* from  sys_role_menu c INNER JOIN sys_menu d on c.menu_id=d.id WHERE d.isvisible=1 and c.role_id =? order by sortindex "
	res := []SysMenu{}
	_, err := o.Raw(sql, role_id).QueryRows(&res)
	return res, err
}

//根据岗位ID获取菜单信息
func GetSysMenuTreeByStationId(stationId int) ([]SysMenu, error) {
	o := orm.NewOrm()
	sql := `SELECT me.* FROM sys_station AS m
	INNER JOIN sys_role_menu AS r
	ON m.role_id=r.role_id
	INNER JOIN sys_menu AS me
	ON r.menu_id=me.id
	WHERE m.id=?
	ORDER BY sortindex`
	res := []SysMenu{}
	_, err := o.Raw(sql, stationId).QueryRows(&res)
	return res, err
}

func GetSysMenuTreeAll() (list []SysMenu, err error) {
	sql := `SELECT * FROM sys_menu WHERE isvisible=1 order by sortindex`
	_, err = orm.NewOrm().Raw(sql).QueryRows(&list)
	return
}

//获取所有组织架构
func GetSysMenu() (map[string]SysMenu, error) {
	o := orm.NewOrm()
	sql := "SELECT * FROM sys_menu where parent_id>0"
	var list []SysMenu
	_, err := o.Raw(sql).QueryRows(&list)
	m := map[string]SysMenu{}
	if err == nil && len(list) > 0 {
		for _, k := range list {
			m[k.ControlUrl] = k
		}
	}
	if data, err2 := json.Marshal(m); err == nil && err2 == nil && utils.Re == nil {
		utils.Rc.Put(utils.CacheKeySystemMenu, data, 5*time.Minute)
	}
	return m, err
}

//根据岗位ID获取菜单信息
func GetBtnPermissions(btnIds string, roleId int) ([]SysBtn, error) {
	o := orm.NewOrm()
	sql := `SELECT sm.controlurl,sm.id FROM sys_menu AS sm 
			INNER JOIN sys_role_menu AS rm ON sm.id=rm.menu_id
			WHERE sm.controlurl IN(` + btnIds + `)
			AND rm.role_id=?`
	res := []SysBtn{}
	_, err := o.Raw(sql, roleId).QueryRows(&res)
	return res, err
}

//校验该按钮是否有权限
func IsBtnPermission(btnId string, roleId int) bool {
	o := orm.NewOrm()
	sql := `SELECT
				COUNT(1)
			FROM sys_menu AS sm
			INNER JOIN sys_role_menu AS rm
			ON sm.id = rm.menu_id
			WHERE sm.controlurl = ?
			AND rm .role_id = ?`
	var count int
	err := o.Raw(sql, btnId, roleId).QueryRow(&count)
	if err != nil {
		return false
	}
	if count > 0 {
		return true
	}
	return false
}

func IsDataPermissionByStationId(stationId, stationType int) bool {
	o := orm.NewOrm()
	sql := `SELECT COUNT(1) FROM sys_station_type WHERE station_id = ? AND type = ?`
	var count int
	err := o.Raw(sql, stationId, stationType).QueryRow(&count)
	if err != nil {
		return false
	}
	if count > 0 {
		return true
	}
	return false
}

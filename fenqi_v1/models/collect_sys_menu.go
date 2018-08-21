package models

import (
	"zcm_tools/orm"
)

func GetCollectSysMenuTreeByRoleId(role_id int) ([]SysMenu, error) {
	o := orm.NewOrm()
	sql := "SELECT DISTINCT d.id,d.parent_id,d.displayname from  sys_role_menu c INNER JOIN sys_menu d on c.menu_id=d.id WHERE d.isvisible=1 AND is_collect_menu = 1 and c.role_id =? order by sortindex "
	res := []SysMenu{}
	_, err := o.Raw(sql, role_id).QueryRows(&res)
	return res, err
}
func GetCollectSysMenuTreeAll() (list []SysMenu, err error) {
	sql := `SELECT id,displayname,controlurl,homeurl,parent_id,sortindex FROM sys_menu WHERE isvisible=1 AND is_collect_menu = 1 order by sortindex`
	_, err = orm.NewOrm().Raw(sql).QueryRows(&list)
	return
}

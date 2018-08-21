package models

import (
	"zcm_tools/orm"
)

func CollectSysRoleList() (list []SysRole, err error) {
	sql := `SELECT id,displayname,remark,org_id FROM sys_role where is_collect_role = 1 ORDER BY CONVERT(displayname USING gbk) COLLATE gbk_chinese_ci ASC`
	_, err = orm.NewOrm().Raw(sql).QueryRows(&list)
	return
}

func CollectSelectRoleList(condition string) (list []SysRole, err error) {
	sql := `SELECT id,displayname FROM sys_role WHERE 1=1 AND is_collect_role = 1 `
	if condition != "" {
		sql += condition
	}
	sql += ` ORDER BY CONVERT(displayname USING gbk) COLLATE gbk_chinese_ci ASC`
	_, err = orm.NewOrm().Raw(sql).QueryRows(&list)
	return
}

func (sr *SysRole) CollectInsert(menu_ids []string) error {
	sql := `INSERT INTO sys_role (displayname, remark, org_id,is_collect_role)
			values(?, ?, ?,?)`
	o := orm.NewOrm()
	o.Begin()
	res, err := o.Raw(sql, sr.Displayname, sr.Remark, sr.Org_id, 1).Exec()
	if err != nil {
		o.Rollback()
		return err
	}

	if len(menu_ids) > 0 {
		rid, err := res.LastInsertId()
		if err != nil {
			o.Rollback()
			return err
		}

		sql = ` INSERT INTO sys_role_menu (role_id, menu_id) values(?, ?)`
		for i := 0; i < len(menu_ids); i++ {
			_, err = o.Raw(sql, rid, menu_ids[i]).Exec()
			if err != nil {
				break
			}
		}
		if err != nil {
			o.Rollback()
			return err
		}
	}
	o.Commit()
	return nil
}

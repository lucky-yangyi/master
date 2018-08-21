package utils

import (
	"strings"
	"zcm_tools/orm"
)

func FlagToTable(rejectType string) (table string, flag, sort_no int) {
	T := map[string]string{"身份证": "ocr_info", "活体验证": "living_info", "居住地": "users_base_info", "联系人1": "users_linkman", "联系人2": "users_linkman"}
	t, ok := T[rejectType]
	if ok {
		table = t
		if strings.Contains(table, "users_linkman") {
			flag = 2
		} else {
			flag = 1
		}
	}
	if strings.Contains(rejectType, "1") {
		sort_no = 1
	}
	if strings.Contains(rejectType, "2") {
		sort_no = 2
	}
	return
}

func UpdateRecjectTable(table, reject_reason string, sort_no, uid int) error {
	o := orm.NewOrm()
	sql := `UPDATE ` + table + ` SET reject_reason = ? WHERE uid = ?`
	if sort_no != 0 {
		sql += ` AND sort_no = ?`
	}
	sql += ` ORDER BY id DESC LIMIT 1`
	p, err := o.Raw(sql).Prepare()
	if sort_no != 0 {
		_, err = p.Exec(reject_reason, uid, sort_no)
	} else {
		_, err = p.Exec(reject_reason, uid)
	}
	return err
}

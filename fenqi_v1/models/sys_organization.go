package models

import (
	"encoding/json"
	"fenqi_v1/utils"
	"strconv"
	"zcm_tools/orm"
)

//组织架构返回数据
type SysOrganization struct {
	Id                   int    `orm:"column(id);pk"`          //主键
	ParentId             int    `orm:"column(parent_id);null"` //父级ID
	Name                 string `orm:"column(name);null"`      //组织架构名称
	Remark               string `orm:"column(remark);null"`    //备注
	NoCheck              bool   `orm:"column(nocheck);null"`   //选择标识
	Type                 int
	InvitationCodePrefix string `orm:"column(invitation_code_prefix);null"`
	IsCollectOrg         int    `orm:"column(is_collect_org)"` //是否是催收组织：0-不是 1-是
}

//结清阶段
type OutSourceType struct {
	OutSourceType string
	OutSourceVal  int
}

//获取所有组织架构
func GetSysOrganization() ([]SysOrganization, error) {
	o := orm.NewOrm()
	sql := "SELECT a.id,a.parent_id,a.name,a.remark from  sys_organization a"
	res := []SysOrganization{}
	_, err := o.Raw(sql).QueryRows(&res)
	if data, err2 := json.Marshal(res); err == nil && err2 == nil && utils.Re == nil {
		utils.Rc.Put(utils.CacheKeySystemOrganization, data, utils.RedisCacheTime_Organization)
	}
	return res, err
}

//根据ID获取组织架构
func GetSysOrganizationById(organizationId int) (organization *SysOrganization, err error) {
	sql := `SELECT DISTINCT a.id,a.parent_id,a.name,a.remark from sys_organization a WHERE a.id=?`
	err = orm.NewOrm().Raw(sql, organizationId).QueryRow(&organization)
	return
}

//根据名称获取组织架构
func GetSysOrganizationByName(organizationName string) (organization *SysOrganization, err error) {
	sql := `SELECT DISTINCT a.id,a.parent_id,a.name,a.remark from sys_organization a WHERE a.name=?`
	err = orm.NewOrm().Raw(sql, organizationName).QueryRow(&organization)
	return
}

//新增组织架构
func (so *SysOrganization) Insert() error {
	sql := `INSERT INTO sys_organization (parent_id, name, remark)
			values(?, ?, ?)`
	_, err := orm.NewOrm().Raw(sql, so.ParentId, so.Name, so.Remark).Exec()
	err = utils.Rc.Delete(utils.CacheKeySystemOrganization)
	return err
}

//修改组织架构名称
func UpdateOrganizationName(name string, id int) (err error) {
	sql := `UPDATE sys_organization SET name=? WHERE id =?`
	_, err = orm.NewOrm().Raw(sql, name, id).Exec()
	return
}

//获取组织架构岗位信息
func GetOrganizationStations() ([]map[string]interface{}, error) {
	sql := `SELECT id,parent_id,name,1 nocheck  FROM sys_organization
			UNION ALL
			SELECT id+100000 id,org_id AS parent_id,name,0 nocheck FROM sys_station`
	list := []*SysOrganization{}

	_, err := orm.NewOrm().Raw(sql).QueryRows(&list)
	l := len(list)
	var org []map[string]interface{}
	for i := 0; i < l; i++ {
		org = append(org, map[string]interface{}{"id": list[i].Id, "pId": list[i].ParentId, "open": false, "name": list[i].Name, "nocheck": list[i].NoCheck})
	}
	return org, err
}

func QueryDisplayQn() (sList []SysStation) {
	var sys_1 []SysOrganization
	o := orm.NewOrm()
	inStr := ""
	sql := `select id,name,parent_id from sys_organization 
	        where parent_id in (select id from sys_organization where name ='电销中心')`
	o.Raw(sql).QueryRows(&sys_1)
	if len(sys_1) > 0 {
		for i := 0; i < len(sys_1); i++ {
			inStr += strconv.Itoa(sys_1[i].Id) + ","
		}
		inStr2 := inStr[:len(inStr)-1]
		sql2 := `select id,name,parent_id from sys_organization 
	        where parent_id in(` + inStr2 + `)`
		var sys_2 []SysOrganization
		o.Raw(sql2).QueryRows(&sys_2)
		if len(sys_2) > 0 {
			inStr3 := ""
			for i := 0; i < len(sys_2); i++ {
				inStr += strconv.Itoa(sys_2[i].Id) + ","
				inStr3 += strconv.Itoa(sys_2[i].Id) + ","
			}
			inStr3 = inStr3[:len(inStr3)-1]
			sql3 := `select id,name,parent_id from sys_organization 
	        where parent_id in(` + inStr3 + `)`
			var sys_3 []SysOrganization
			o.Raw(sql3).QueryRows(&sys_3)
			if len(sys_3) > 0 {
				inStr4 := ""
				for i := 0; i < len(sys_3); i++ {
					inStr += strconv.Itoa(sys_3[i].Id) + ","
					inStr4 += strconv.Itoa(sys_3[i].Id) + ","
				}
				inStr4 = inStr4[:len(inStr4)-1]
				sql4 := `select id,name,parent_id from sys_organization
				         where parent_id in(` + inStr4 + `)`
				var sys_4 []SysOrganization
				o.Raw(sql4).QueryRows(&sys_4)
				if len(sys_4) > 0 {
					// inStr5:=""
					for i := 0; i < len(sys_4); i++ {
						inStr += strconv.Itoa(sys_4[i].Id) + ","
						// inStr5 += strconv.Itoa(sys_4[i].Id) + ","
					}
					// inStr5 = inStr5[:len(inStr4)-1]
					// sql5 := `select id,name,parent_id from sys_organization
					//         where parent_id in(` + inStr5 + `)`
					// var sys_5 []SysOrganization
					// o.Raw(sql5).QueryRows(&sys_5)
					// for i := 0; i < len(sys_5); i++ {
					// 	inStr += strconv.Itoa(sys_5[i].Id) + ","
					// }
				}
			}
		}
	}
	var sys SysOrganization
	o.Raw(`select id from sys_organization where name ='电销中心'`).QueryRow(&sys)
	inStr += strconv.Itoa(sys.Id)
	stationSql := `select * from sys_station where org_id in (` + inStr + `)`
	o.Raw(stationSql).QueryRows(&sList)
	return
}

//根据ID获取子组织架构
func GetParentSysOrganizationById(organizationId int) ([]SysOrganization, error) {
	var sos []SysOrganization
	sql := `SELECT a.id,a.parent_id,a.name,a.remark from sys_organization a WHERE a.parent_id=?`
	_, err := orm.NewOrm().Raw(sql, organizationId).QueryRows(&sos)
	return sos, err
}

//根据区省营业部和岗位id判断权限
func QuerySysRrganizationIdByRPO(region, province, office, stationId int) (is bool, err error) {
	if region == 0 {
		region = 28
	}
	sql := `SELECT IF ((SELECT id FROM sys_station WHERE org_id = ? AND id =?)IS NULL ,(SELECT IF ((SELECT id FROM sys_station WHERE org_id = ? AND id =?)IS NULL ,(SELECT IF ((SELECT id FROM sys_station WHERE org_id = ? AND id =?)IS NULL ,"false ","true")),"true")),"true")`
	err = orm.NewOrm().Raw(sql, office, stationId, province, stationId, region, stationId).QueryRow(&is)
	return
}

//根据stationId获取组织id
func QueryOrgIdByStationId(stationId int) (orgId int, err error) {
	sql := `SELECT so.id FROM sys_organization so LEFT JOIN sys_station ss ON so.id = ss.org_id WHERE ss.id = ?`
	err = orm.NewOrm().Raw(sql, stationId).QueryRow(&orgId)
	return
}

func QueryParentIdByOrgId(orgId int) (parentId int, err error) {
	sql := `SELECT so.parent_id FROM sys_organization so WHERE so.id = ?`
	err = orm.NewOrm().Raw(sql, orgId).QueryRow(&parentId)
	return
}

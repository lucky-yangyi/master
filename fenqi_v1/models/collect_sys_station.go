package models

import (
	"strconv"
	"zcm_tools/orm"
)

//添加岗位信息
func (station *SysStation) CollectInsert(stationType, stationData []string) error {
	o := orm.NewOrm()
	sql := `INSERT INTO sys_station(name,role_id,org_id,is_collect_station)values(?,?,?,?)`
	o.Begin()
	res, err := o.Raw(sql, station.Name, station.RoleId, station.OrgId, station.IsCollectStation).Exec()
	if err != nil {
		o.Rollback()
		return err
	}
	//获取岗位ID
	stationId, err := res.LastInsertId()
	if err != nil {
		o.Rollback()
		return err
	}
	//添加岗位类型
	if len(stationType) > 0 {
		sql = `INSERT INTO sys_station_type(station_id,type)values(?,?)`
		typeSql, err := o.Raw(sql).Prepare()
		for i := 0; i < len(stationType); i++ {
			typeInt, _ := strconv.Atoi(stationType[i])
			if typeInt > 0 {
				_, err = typeSql.Exec(stationId, typeInt)
				if err != nil {
					break
				}
			}
		}
		typeSql.Close()
		if err != nil {
			o.Rollback()
			return err
		}
	}
	//添加岗位权限
	if len(stationData) > 0 {
		sql = `INSERT INTO sys_station_data(station_id,org_id)values(?,?)`
		dataSql, err := o.Raw(sql).Prepare()
		defer dataSql.Close()
		for i := 0; i < len(stationData); i++ {
			dataInt, _ := strconv.Atoi(stationData[i])
			if dataInt > 0 {
				_, err = dataSql.Exec(stationId, dataInt)
				if err != nil {
					break
				}
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

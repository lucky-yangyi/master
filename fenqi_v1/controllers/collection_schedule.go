package controllers

import (
	"fenqi_v1/cache"
	"fenqi_v1/models"
	"fenqi_v1/utils"
	"strconv"
	"time"
)

//催收排班控制器
type CollectionScheduleController struct {
	BaseController
}

//@router /getcollectionscheduledetail [get]
func (this *CollectionScheduleController) GetCollectionScheduleDetail() {
	this.IsNeedTemplate()
	where := ""
	var pars []interface{}
	var currentTime, nextTime time.Time
	var monthSlice []string
	where += " AND cs.plan_year = ? AND cs.plan_month = ?"
	//获取排班时间(年月)
	if date := this.GetString("beginTime"); date != "" {
		dates, _ := time.Parse("2006-01-02", date+"-01")
		pars = append(pars, dates.Year(), int(dates.Month()))
		currentTime, _ = time.Parse("2006-01", date)
	} else {
		pars = append(pars, time.Now().Year(), int(time.Now().Month()))
		currentTime, _ = time.Parse("2006-01", time.Now().Format("2006-01"))
	}
	nextTime = currentTime.AddDate(0, 1, -1)
	monthSlice = utils.GetSeriesMonth(currentTime, nextTime)
	//催收所属分组
	group := this.GetString("type")
	if group == "" {
		group = "11"
	}
	this.Data["group"] = group
	collectionName := this.GetString("collectionName") //催收人员姓名
	if collectionName != "" {
		where += " AND u.displayname LIKE '%" + collectionName + "%'"
	}
	collectionSchedules, err := models.GetCollectionSchedule(where, pars, group)
	collectionScheduleRecords, err := models.GetCollectionScheduleGroup(group, collectionName)
	if err != nil {
		cache.RecordLogs(this.User.Id, 0, this.User.Name, this.User.DisplayName, "查询催收排班列表失败", "催收排班/GetCollectionScheduleDetail", err.Error(), this.Ctx.Input)
	}
	for k, v1 := range collectionScheduleRecords {
		for key, v := range collectionSchedules {
			if v.SysUid == v1.SysUid && v.Type == v1.Type {
				collectionScheduleRecords[k].Days[v.PlanDay-1] = v.State
				if key != len(collectionSchedules)-1 {
					currentSysUid := collectionSchedules[key].SysUid
					currentType := strconv.Itoa(collectionSchedules[key].Type)
					currentYear := collectionSchedules[key].PlanYear
					currentMonth := collectionSchedules[key].PlanMonth
					nextSysUid := collectionSchedules[key+1].SysUid
					nextType := strconv.Itoa(collectionSchedules[key+1].Type)
					nextYear := collectionSchedules[key+1].PlanYear
					nextMonth := collectionSchedules[key+1].PlanMonth
					if currentSysUid != nextSysUid && currentType != nextType && currentYear != nextYear && currentMonth != nextMonth {
						break
					}
				}
			}
		}
		collectionScheduleRecords[k].Name = v1.Name + utils.Mtype(v1.Type)
	}
	this.Data["nowTime"] = time.Now().Unix() * 1000
	this.Data["collectionScheduleRecords"] = collectionScheduleRecords
	this.Data["oneMonth"] = monthSlice
	this.Data["monthLength"] = len(monthSlice) - 1
	this.TplName = "collection/calendar_list.html"
}

//更新催收排班表
//@router /updatecollectionschedule [post]
func (this *CollectionScheduleController) UpdateCollectionSchedule() {
	resultMap := make(map[string]interface{})
	resultMap["ret"] = 304
	defer func() {
		this.Data["json"] = resultMap
		this.ServeJSON()
	}()
	sysUid, err := this.GetInt("sysUid") //催收人员Id
	types, err := this.GetInt("type")    //催收人员所属分组
	date := this.GetString("date")       //排班日期
	state, err := this.GetInt("state")   //排班状态
	if err != nil {
		resultMap["err"] = "参数有误"
		return
	}
	currentDate := time.Now().Format("2006-01-02")
	dates, err := time.Parse("2006-01-02", date)
	if err != nil {
		resultMap["err"] = "格式化时间出错!"
		return
	}
	if currentDate >= date {
		resultMap["err"] = "时间参数有误，该日期排班不能修改!"
		return
	}
	planYear, planMonth, planDay := dates.Date()
	var periodTypes []int
	periodTypes = append(periodTypes, types)
	isExist := models.CheckCollectionScheduleIsExist(sysUid, types, planYear, int(planMonth), planDay)
	if isExist {
		err = models.UpdateMultilCollectionSchedules(sysUid, state, planYear, int(planMonth), planDay, date, periodTypes)
		if err != nil {
			resultMap["err"] = "保存排班信息失败"
			return
		}
	} else {
		err = models.InsertMultilCollectionSchedules(sysUid, state, planYear, int(planMonth), planDay, date, periodTypes)
		if err != nil {
			resultMap["err"] = "保存排班信息失败"
			return
		}
	}
	resultMap["ret"] = 200
	resultMap["msg"] = "保存排班成功"
}

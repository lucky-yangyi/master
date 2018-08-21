package controllers

import (
	"fenqi_v1/cache"
	"fenqi_v1/models"
	"fenqi_v1/utils"
	"strconv"
	"strings"
	"time"
)

//  入黑标记
type BlackController struct {
	BaseController
}

// 获取列表
func (this *BlackController) IntoBlackSign() {
	this.IsNeedTemplate()
	condition := ""
	params := []interface{}{}
	//获取当前页
	currentPage, _ := this.GetInt("page")
	if currentPage <= 0 {
		currentPage = 1
	}
	pageSize, _ := this.GetInt("pageSize")
	if pageSize < 1 {
		pageSize = 30
	}
	startTime := this.GetString("startTime")
	if startTime != "" {
		condition += "AND bs.check_time >= ?"
		params = append(params, startTime)
	}
	endTime := this.GetString("endTime")
	if endTime != "" {
		condition += " AND bs.check_time <= ? "
		params = append(params, endTime)

	}
	account := this.GetString("account")
	if account != "" {
		condition += "  AND bs.account = ? "
		params = append(params, account)
	}

	// 催收阶段
	collectionType := this.GetString("collectionType")
	if collectionType != "" {
		condition += " AND bs.collection_type = ? "
		params = append(params, collectionType)
	}

	//行动分类
	actionType := this.GetString("actionType")
	if actionType != "" {
		condition += " AND bs.action_type = ? "
		params = append(params, actionType)
	}

	isTag, _ := this.GetInt("isTag")
	if isTag > 0 {
		condition += "  AND bs.is_tag = ? "
		params = append(params, isTag)
	}

	state, _ := this.GetInt("state")
	if state > 0 {
		condition += " AND bs.state = ? "
		params = append(params, state)
	}
	//list
	blackSignList, err := models.GetBlackSignList(condition, params, utils.StartIndex(currentPage, pageSize), pageSize)
	if err != nil {
		cache.RecordLogs(this.User.Id, 0, this.User.Name, this.User.DisplayName, "获取入黑标记数据失败", "催收管理/入黑标记", err.Error(), this.Ctx.Input)
	}
	//pagecount
	pageCount, err := models.GetBlackSignPageCount(condition, params)
	if err != nil {
		cache.RecordLogs(this.User.Id, 0, this.User.Name, this.User.DisplayName, "获取入黑标记总记录数失败", "催收管理/入黑标记", err.Error(), this.Ctx.Input)
	}
	this.Data["blackSignList"] = blackSignList
	this.Data["pageCount"] = pageCount
	this.Data["currentPage"] = currentPage
	this.Data["pageSize"] = pageSize
	this.Data["pageNum"] = utils.PageCount(pageCount, pageSize)
	this.TplName = "collection/collection_shielding_tag.html"

}

//批量入黑标记处理操作
func (this *BlackController) BatchHandlerBlackSign() {
	resultMap := make(map[string]interface{})
	resultMap["ret"] = 304
	defer func() {
		this.Data["json"] = resultMap
		this.ServeJSON()
	}()
	blackUids := this.GetString("blackUids")
	if blackUids == "" {
		cache.RecordLogs(this.User.Id, 0, this.User.Name, this.User.DisplayName, "uid不能为空", "入黑标记/BatchHandlerBlackSign", "uid不能为空", this.Ctx.Input)
		resultMap["err"] = "uid不能为空"
		return
	}
	uids := strings.Split(blackUids, ",")
	checkManId := this.User.Id
	state, _ := this.GetInt("state") //状态id
	// 查询标记类型
	var isTags []int
	for _, uid := range uids {
		if uid != "" {
			id, err := strconv.Atoi(uid)
			isTage, err := models.GetIsTag(id)
			if err == nil {
				isTags = append(isTags, isTage)
			}
		}
	}
	// 审核 ,同时更新users表
	var tagTypes []string
	for _, v := range isTags {
		var tagtype string
		switch v {
		case 0:
			tagtype = "未标记欺诈"
		case 1:
			tagtype = "家人代偿"
		case 2:
			tagtype = "本人还款意愿差"
		case 3:
			tagtype = "其他"
		}
		tagTypes = append(tagTypes, tagtype)
	}
	err := models.BatchUpdateBlackSignState(state, checkManId, isTags, tagTypes, uids)
	if err != nil {
		cache.RecordLogs(this.User.Id, 0, this.User.Name, this.User.DisplayName, "入黑标记处理错误！", "入黑标记/BatchHandlerBlackSign", err.Error(), this.Ctx.Input)
		resultMap["err"] = "入黑标记批量处理错误"
		return
	}
	resultMap["ret"] = 200
	resultMap["msg"] = "批量入黑标记成功"
}

//入黑标记处理操作
func (this *BlackController) HandlerBlackSign() {
	defer this.ServeJSON()
	blackUserId, _ := this.GetInt("BlackUserId") //入黑用户id
	checkManId := this.User.Id
	state, _ := this.GetInt("State") //状态id
	// 查询标记类型
	istage, err := models.GetIsTag(blackUserId)
	if err != nil {
		this.Data["json"] = map[string]interface{}{"ret": 403, "err": "获取入黑类型失败"}
		return
	}
	// 审核 ,同时更新users表
	var tagtype string
	switch istage {
	case 0:
		tagtype = "未标记欺诈"
	case 1:
		tagtype = "家人代偿"
	case 2:
		tagtype = "本人还款意愿差"
	case 3:
		tagtype = "其他"
	}
	err = models.HandlerBlackSignUserState(state, checkManId, blackUserId, istage, tagtype)
	if err != nil {
		cache.RecordLogs(this.User.Id, 0, this.User.Name, this.User.DisplayName, "入黑标记处理错误！", "催收管理/入黑标记", err.Error(), this.Ctx.Input)
		this.Data["json"] = map[string]interface{}{"ret": 403, "err": "入黑标记处理错误！"}
		return
	}
	this.Data["json"] = map[string]interface{}{"ret": 200}
}

// 提交入黑
func (this *BlackController) SaveBlackSign() {
	defer this.ServeJSON()
	var blacksign models.BlackSignSave
	blacksign.Uid, _ = this.GetInt("uid")
	if blacksign.Uid <= 0 {
		this.Data["json"] = map[string]interface{}{"ret": 403, "msg": "获取用户id失败"}
		return
	}
	blacksign.Account = this.GetString("account")
	blacksign.SubmitManId = this.User.Id
	blacksign.SubmitTime = time.Now().Format("2006-01-02 15:04:05")
	blacksign.State = 1
	blacksign.IsTag, _ = this.GetInt("flagType")

	blacksign.Remark = this.GetString("cheat_content")
	loanId, _ := this.GetInt("loanId")
	mtype, _ := this.GetInt("mtype")
	blacksign.CollectionType = utils.Mtype(mtype)
	// 如果获取的催收阶段为空
	if blacksign.CollectionType == "" {
		// 根据 loanId 查询逾期
		Od, err := models.GetOverdueDaysByLoanId(loanId)
		if err != nil {
			cache.RecordLogs(this.User.Id, 0, this.User.Name, this.User.DisplayName, "根据loanId查询逾期天数异常！", "催收管理/入黑标记", err.Error(), this.Ctx.Input)
		}
		blacksign.CollectionType = utils.OverudDays(Od.OverdueDays)

	}
	// 判断用户是否存在,存在并已退回的话，更新数据
	num, err := models.UpdateQueryBlackSign(blacksign.Uid)
	if err != nil {
		cache.RecordLogs(this.User.Id, 0, this.User.Name, this.User.DisplayName, "查询标记入黑异常！", "催收管理/入黑标记", err.Error(), this.Ctx.Input)
		this.Data["json"] = map[string]interface{}{"ret": 403, "msg": "查询标记入黑异常！"}
		return
	}
	if num == 0 {
		// 保存
		err = models.SaveBlackSign(blacksign)
		if err != nil {
			cache.RecordLogs(this.User.Id, 0, this.User.Name, this.User.DisplayName, "保存入黑标记失败", "催收管理/入黑标记", err.Error(), this.Ctx.Input)
		}
		this.Data["json"] = map[string]interface{}{"ret": 200, "msg": "标记入黑保存成功"}
		return
	}
	if num == 1 {
		// 更新
		err = models.UpdateReturnUser(blacksign)
		if err != nil {
			cache.RecordLogs(this.User.Id, 0, this.User.Name, this.User.DisplayName, "已退回重新标记异常", "催收管理/入黑标记", err.Error(), this.Ctx.Input)
		}
		this.Data["json"] = map[string]interface{}{"ret": 200, "msg": "标记入黑更新成功"}
		return
	}

}

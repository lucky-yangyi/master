package controllers

import (
	"fenqi_v1/cache"
	"fenqi_v1/models"
	"fenqi_v1/utils"
	"strings"
	"time"
)

//意见反馈
type AdviseController struct {
	BaseController
}

//意见反馈列表
func (this *AdviseController) List() {
	this.IsNeedTemplate()
	pars := []interface{}{}
	condition := ""
	pageNum, _ := this.GetInt("page", 1) //分页信息（第几页）
	if pageNum < 1 {
		pageNum = 1
	}
	pageSize, _ := this.GetInt("pageSize")
	if pageSize < 1 {
		pageSize = 30
	}
	account := this.GetString("account") //手机号
	if account != "" {
		condition += " AND account = ?"
		pars = append(pars, account)
	}
	feedbackTime := this.GetString("feedback_time") //反馈时间
	var startFeedbackTimeTime, endFeedbackTimeTime string
	if feedbackTime != "" {
		feedbackTimes := strings.Split(feedbackTime, "~")
		startFeedbackTimeTime = feedbackTimes[0] + " 00:00:00"
		endFeedbackTimeTime = feedbackTimes[1] + " 23:59:59"
	} else { //默认提交时间7天内
		startFeedbackTimeTime = time.Now().AddDate(0, 0, -6).Format("2006-01-02") + " 00:00:00"
		endFeedbackTimeTime = time.Now().Format("2006-01-02") + " 23:59:59"
	}
	condition += ` AND create_time >= ? AND create_time <= ?`
	pars = append(pars, startFeedbackTimeTime)
	pars = append(pars, endFeedbackTimeTime)
	//状态
	if isChecked := this.GetString("is_checked"); isChecked != "" {
		if isChecked == "0" {
			condition += " AND is_checked = 0 "
		}
		if isChecked == "1" {
			condition += " AND is_checked > 0 "
		}
	}

	adviseList, err := models.QueryAdviseList(utils.StartIndex(pageNum, pageSize), pageSize, condition, pars...)
	if err != nil {
		cache.RecordLogs(this.User.Id, 0, this.User.Name, this.User.DisplayName, "查询意见反馈列表失败", "意见反馈/List", err.Error(), this.Ctx.Input)
		this.Abort("查询意见反馈列表失败")
		return
	}
	count, err := models.QueryAdviseCount(condition, pars...)
	if err != nil {
		cache.RecordLogs(this.User.Id, 0, this.User.Name, this.User.DisplayName, "获取总数异常", "意见反馈/List", err.Error(), this.Ctx.Input)
		this.Abort("获取意见反馈总数异常")
		return
	}
	pageCount := utils.PageCount(count, pageSize)
	this.Data["list"] = adviseList
	this.Data["currPage"] = pageNum
	this.Data["count"] = count
	this.Data["pageSize"] = pageSize
	this.Data["pageCount"] = pageCount
	this.TplName = "advise.html"
}

//检查该用户是否为认证用户
func (this *AdviseController) CheckUsersIsAuth() {
	resultMap := make(map[string]interface{})
	resultMap["ret"] = 403
	defer func() {
		this.Data["json"] = resultMap
		this.ServeJSON()
	}()
	uid, _ := this.GetInt("uid")
	if uid <= 0 {
		cache.RecordLogs(this.User.Id, 0, this.User.Name, this.User.DisplayName, "uid参数传递错误", "意见反馈/CheckUsersIsAuth", "", this.Ctx.Input)
		resultMap["err"] = "uid参数传递错误"
		return
	}
	count, err := models.GetUsersAuthCount(uid)
	if err != nil {
		cache.RecordLogs(this.User.Id, 0, this.User.Name, this.User.DisplayName, "查询认证用户数据错误", "意见反馈/CheckUsersIsAuth", err.Error(), this.Ctx.Input)
		resultMap["err"] = "查询认证用户数据错误" + err.Error()
		return
	}
	resultMap["ret"] = 200
	resultMap["count"] = count
}

//获得一条反馈
func (this *AdviseController) GetAdvise() {
	resultMap := make(map[string]interface{})
	resultMap["ret"] = 403
	defer func() {
		this.Data["json"] = resultMap
		this.ServeJSON()
	}()
	adviseId, _ := this.GetInt("advise_id")
	if adviseId <= 0 {
		cache.RecordLogs(this.User.Id, 0, this.User.Name, this.User.DisplayName, "advise_id参数传递错误", "意见反馈/GetAdvise", "", this.Ctx.Input)
		resultMap["err"] = "advise_id参数传递错误"
		return
	}
	advise, err := models.QueryAdviseById(adviseId)
	if err != nil && err.Error() != utils.ErrNoRow() {
		cache.RecordLogs(this.User.Id, 0, this.User.Name, this.User.DisplayName, "查询反馈数据错误", "意见反馈/GetAdvise", err.Error(), this.Ctx.Input)
		resultMap["err"] = "查询反馈数据错误" + err.Error()
		return
	}
	resultMap["ret"] = 200
	resultMap["advise"] = advise
}

//更新意见反馈备注
func (this *AdviseController) UpdateAdvise() {
	resultMap := make(map[string]interface{})
	resultMap["ret"] = 403
	defer func() {
		this.Data["json"] = resultMap
		this.ServeJSON()
	}()
	adviseId, _ := this.GetInt("advise_id")
	if adviseId <= 0 {
		cache.RecordLogs(this.User.Id, 0, this.User.Name, this.User.DisplayName, "advise_id参数传递错误", "意见反馈/UpdateAdvise", "", this.Ctx.Input)
		resultMap["err"] = "advise_id参数传递错误"
		return
	}
	remark := this.GetString("remark")
	changeStatus, _ := this.GetInt("change_status")
	if changeStatus == 0 {
		changeStatus = 9
	}
	err := models.UpdateAdviseRemark(remark, this.User.DisplayName, changeStatus, adviseId)
	if err != nil {
		cache.RecordLogs(this.User.Id, 0, this.User.Name, this.User.DisplayName, "意见反馈修改备注失败", "意见反馈/UpdateAdvise", err.Error(), this.Ctx.Input)
		resultMap["err"] = "意见反馈修改备注失败" + err.Error()
		return
	}
	resultMap["ret"] = 200
	resultMap["msg"] = "更改状态成功"
}

//业务员意见反馈
func (this *AdviseController) SalesmanList() {
	this.IsNeedTemplate()
	pars := []interface{}{}
	condition := ""
	pageNum, _ := this.GetInt("page", 1) //分页信息（第几页）
	if pageNum < 1 {
		pageNum = 1
	}
	pageSize, _ := this.GetInt("pageSize")
	if pageSize < 1 {
		pageSize = 30
	}
	startSubmitTime := this.GetString("start_register_time") //提交时间
	if startSubmitTime != "" {
		condition += ` AND a.creat_time >= ?`
		pars = append(pars, startSubmitTime)
	}
	endSubmitTime := this.GetString("end_register_time") //提交时间
	if endSubmitTime != "" {
		condition += ` AND a.creat_time <= ?`
		pars = append(pars, endSubmitTime)
	}
	saleman := this.GetString("saleman") //业务员
	if saleman != "" {
		condition += " AND b.saleman = ?"
		pars = append(pars, saleman)
	}
	account := this.GetString("account") //手机号
	if account != "" {
		condition += " AND b.account = ?"
		pars = append(pars, account)
	}
	isChecked := this.GetString("is_checked") //审批状态
	if isChecked == "ok_app" {
		condition += " AND a.state=?"
		pars = append(pars, 1)
	} else if isChecked == "no_app" {
		condition += " AND a.state=?"
		pars = append(pars, 0)
	}
	invite_code := this.GetString("invite_code") //邀请码
	if invite_code != "" {
		condition += " AND b.invite_code = ?"
		pars = append(pars, invite_code)
	}
	//if condition == "" { //默认提交时间7天内
	//	condition += ` AND b.create_time >= ? AND b.create_time <= ?`
	//	startSubmitTime = time.Now().AddDate(0, 0, -6).Format("2006-01-02") + " 00:00:00"
	//	endSubmitTime = time.Now().Format("2006-01-02") + " 23:59:59"
	//	pars = append(pars, startSubmitTime)
	//	pars = append(pars, endSubmitTime)
	//}
	salemnan_advise_list, err := models.QuerySalmanAdviseList(utils.StartIndex(pageNum, pageSize), pageSize, condition, pars...)
	if err != nil {
		cache.RecordLogs(this.User.Id, 0, this.User.Name, this.User.DisplayName, "查询意见反馈列表失败", "意见反馈/SalesmanList", err.Error(), this.Ctx.Input)
		this.Abort("查询意见反馈列表失败")
		return
	}
	count, err := models.QuerySalemanAdviseCount(condition, pars...)
	if err != nil {
		cache.RecordLogs(this.User.Id, 0, this.User.Name, this.User.DisplayName, "获取总数异常", "意见反馈/SalesmanList", err.Error(), this.Ctx.Input)
		this.Abort("获取意见反馈总数异常")
		return
	}
	pageCount := utils.PageCount(count, pageSize)
	this.Data["list"] = salemnan_advise_list
	this.Data["currPage"] = pageNum
	this.Data["count"] = count
	this.Data["pageSize"] = pageSize
	this.Data["pageCount"] = pageCount
	this.TplName = "clerk/clerk_advice.html"
}

//更新意见反馈备注
func (this *AdviseController) UpdateSalemanAdvise() {
	resultMap := make(map[string]interface{})
	resultMap["ret"] = 403
	defer func() {
		this.Data["json"] = resultMap
		this.ServeJSON()
	}()
	id, _ := this.GetInt("id")
	if id <= 0 {
		cache.RecordLogs(this.User.Id, 0, this.User.Name, this.User.DisplayName, "advise_id参数传递错误", "意见反馈/UpdateSalemanAdvise", "", this.Ctx.Input)
		resultMap["err"] = "advise_id参数传递错误"
		return
	}
	remark := this.GetString("remark")
	err := models.UpdateSalemanAdviseRemark(remark, this.User.DisplayName, id)
	if err != nil {
		cache.RecordLogs(this.User.Id, 0, this.User.Name, this.User.DisplayName, "意见反馈修改备注失败", "意见反馈/UpdateSalemanAdvise", err.Error(), this.Ctx.Input)
		resultMap["err"] = "意见反馈修改备注失败" + err.Error()
		return
	}
	resultMap["ret"] = 200
	resultMap["msg"] = "更改状态成功"
}

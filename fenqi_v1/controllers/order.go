package controllers

import (
	"encoding/json"
	"fenqi_v1/cache"
	"fenqi_v1/models"
	"fenqi_v1/services"
	"fenqi_v1/utils"
	"github.com/astaxie/beego"
	"strconv"
	"strings"
	"time"
)

//订单管理
func (this *WorkspaceController) OrderList() {
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
	account := strings.TrimSpace(this.GetString("account")) //手机号
	if account != "" {
		condition += " AND um.account = ?"
		pars = append(pars, account)
	}
	userName := strings.TrimSpace(this.GetString("user_name")) //姓名
	if userName != "" {
		condition += " AND um.verifyrealname = ?"
		pars = append(pars, userName)
	}
	loanTermCount := this.GetString("loan_term_count") //期限
	if loanTermCount != "" {
		condition += " AND l.loan_term_count = ?"
		pars = append(pars, loanTermCount)
	}
	operatorName := this.GetString("operator_name") //处理人
	if operatorName != "" {
		condition += " AND l.credit_operator = ?"
		pars = append(pars, operatorName)
	}
	handleState := this.GetString("handle_state") //状态
	if handleState != "" {
		condition += ` AND l.order_state = ?`
		pars = append(pars, handleState)
	}
	startMoney := strings.TrimSpace(this.GetString("start_money")) //金额
	if startMoney != "" {
		condition += ` AND l.money >= ?`
		pars = append(pars, startMoney)
	}
	endMoney := strings.TrimSpace(this.GetString("end_money"))
	if endMoney != "" {
		condition += ` AND l.money <= ?`
		pars = append(pars, endMoney)
	}
	submitTime := this.GetString("submit_time") //提交时间
	var startSubmitTime, endSubmitTime string
	if submitTime != "" {
		submitTimes := strings.Split(submitTime, "~")
		startSubmitTime = submitTimes[0] + " 00:00:00"
		endSubmitTime = submitTimes[1] + " 23:59:59"
	} else { //默认提交时间7天内
		startSubmitTime = time.Now().AddDate(0, 0, -6).Format("2006-01-02") + " 00:00:00"
		endSubmitTime = time.Now().Format("2006-01-02 15:04:05")
	}
	condition += ` AND l.create_time >= ? AND l.create_time <= ?`
	pars = append(pars, startSubmitTime)
	pars = append(pars, endSubmitTime)
	dealTime := this.GetString("deal_time") //处理时间
	if dealTime != "" {
		dealTimes := strings.Split(dealTime, "~")
		condition += ` AND l.handle_time >= ? AND l.handle_time <= ?`
		pars = append(pars, dealTimes[0]+" 00:00:00")
		pars = append(pars, dealTimes[1]+" 23:59:59")
	}
	orderCreditList, err := models.GetOrderCreditList(utils.StartIndex(pageNum, pageSize), pageSize, condition, pars...)
	for k, v := range orderCreditList {
		if v.InqueueTime.Format(utils.FormatDate) != "0001-01-01" && v.InqueueTime.Before(time.Now()) {
			orderCreditList[k].OrderState = "QUEUEING"
			orderCreditList[k].CreditOperator = ""
			t, _ := time.Parse(utils.FormatDateTime, "0001-01-01")
			orderCreditList[k].Atime = t
			orderCreditList[k].HandleTime = t
		}
	}
	if err != nil {
		cache.RecordLogs(this.User.Id, 0, this.User.Name, this.User.DisplayName, "查询订单管理失败", "workspace/OrderList", err.Error(), this.Ctx.Input)
		this.Abort("查询订单管理失败")
	}
	count, err := models.GetOrderCreditCount(condition, pars...)
	if err != nil {
		cache.RecordLogs(this.User.Id, 0, this.User.Name, this.User.DisplayName, "获取总数异常", "workspace/OrderList", err.Error(), this.Ctx.Input)
		this.Abort("获取订单管理总数异常")
	}
	creditOperators, err := models.GetCreditOperators()
	if err != nil {
		cache.RecordLogs(this.User.Id, 0, this.User.Name, this.User.DisplayName, "获取处理员异常", "workspace/OrderList", err.Error(), this.Ctx.Input)
		this.Abort("获取处理员异常")
	}
	termCount, err := models.QueryProductTermCount()
	if err != nil {
		cache.RecordLogs(this.User.Id, 0, this.User.Name, this.User.DisplayName, "获取产品期数异常", "workspace/OrderList", err.Error(), this.Ctx.Input)
		this.Abort("获取产品期数异常")
	}
	pageCount := utils.PageCount(count, pageSize)
	this.Data["operators"] = creditOperators
	this.Data["termCount"] = termCount
	this.Data["list"] = orderCreditList
	this.Data["currPage"] = pageNum
	this.Data["count"] = count
	this.Data["pageSize"] = pageSize
	this.Data["pageCount"] = pageCount
	this.TplName = "workspace/workspace_order_list.html"
}

//订单退回管理
func (this *WorkspaceController) OrderOutqueue() {
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
	account := strings.TrimSpace(this.GetString("account")) //手机号
	if account != "" {
		condition += " AND um.account = ?"
		pars = append(pars, account)
	}
	userName := strings.TrimSpace(this.GetString("user_name")) //姓名
	if userName != "" {
		condition += " AND um.verifyrealname = ?"
		pars = append(pars, userName)
	}

	operatorName := this.GetString("operator_name") //处理人
	if operatorName != "" {
		condition += " AND l.credit_operator = ?"
		pars = append(pars, operatorName)
	}
	submitTime := this.GetString("submit_time") //提交时间
	var startSubmitTime, endSubmitTime string
	if submitTime != "" {
		submitTimes := strings.Split(submitTime, "~")
		startSubmitTime = submitTimes[0] + " 00:00:00"
		endSubmitTime = submitTimes[1] + " 23:59:59"
	} else { //默认提交时间近一个月
		startSubmitTime = time.Now().AddDate(0, -1, 0).Format("2006-01-02") + " 00:00:00"
		endSubmitTime = time.Now().Format("2006-01-02 15:04:05")
	}
	condition += ` AND l.create_time >= ? AND l.create_time <= ?`
	pars = append(pars, startSubmitTime)
	pars = append(pars, endSubmitTime)
	dealTime := this.GetString("deal_time") //处理时间
	if dealTime != "" {
		dealTimes := strings.Split(dealTime, "~")
		condition += ` AND l.handle_time >= ? AND l.handle_time <= ?`
		pars = append(pars, dealTimes[0]+" 00:00:00")
		pars = append(pars, dealTimes[1]+" 23:59:59")
	}
	appointTime := this.GetString("appoint_time") //出列预约时间
	if appointTime != "" {
		appointTimes := strings.Split(appointTime, "~")
		condition += ` AND l.inqueue_time >= ? AND l.inqueue_time <= ?`
		pars = append(pars, appointTimes[0]+" 00:00:00")
		pars = append(pars, appointTimes[1]+" 23:59:59")
	}

	outQueueList, err := models.GetOutQueueList(utils.StartIndex(pageNum, pageSize), pageSize, condition, pars...)
	if err != nil {
		cache.RecordLogs(this.User.Id, 0, this.User.Name, this.User.DisplayName, "查询订单退回列表失败", "workspace/OrderOutqueue", err.Error(), this.Ctx.Input)
		this.Abort("查询订单退回列表失败")
	}
	count, err := models.GetOutQueueCount(condition, pars...)
	if err != nil {
		cache.RecordLogs(this.User.Id, 0, this.User.Name, this.User.DisplayName, "获取总数异常", "workspace/OrderOutqueue", err.Error(), this.Ctx.Input)
		this.Abort("获取订单退回列表总数异常")
	}
	creditOperators, err := models.GetCreditOperators()
	pageCount := utils.PageCount(count, utils.PageSize10)
	this.Data["operators"] = creditOperators
	this.Data["list"] = outQueueList
	this.Data["currPage"] = pageNum
	this.Data["count"] = count
	this.Data["pageSize"] = pageSize
	this.Data["pageCount"] = pageCount
	this.TplName = "workspace/workspace_order_outqueue_list.html"
}

//===========================================================================

//出列状态的可以手动选择入列时间和类型
func (this *WorkspaceController) InQueue() {
	defer this.ServeJSON()
	var oc models.Loading
	err := json.Unmarshal(this.Ctx.Input.RequestBody, &oc)
	if err != nil {
		this.Data["json"] = map[string]interface{}{"ret": 304, "err": err.Error(), "msg": "参数解析失败"}
		cache.RecordLogs(this.User.Id, 0, this.User.Name, this.User.DisplayName, "参数解析失败", "workspace/InQueue", err.Error(), this.Ctx.Input)
		return
	}
	err = models.UpadateInQueue(oc)
	if err != nil {
		this.Data["json"] = map[string]interface{}{"ret": 304, "err": err.Error(), "msg": "更新预约状态和时间失败"}
		cache.RecordLogs(this.User.Id, 0, this.User.Name, this.User.DisplayName, "更新预约状态和时间失败", "workspace/InQueue", err.Error(), this.Ctx.Input)
		return
	}
	this.Data["json"] = map[string]interface{}{"ret": 200, "msg": "入列成功"}
}

//=========================================================================
//处理(订单信审)
func (this *WorkspaceController) OrderHanding() {
	this.IsNeedTemplate()
	if ok := models.IsDataPermissionByStationId(this.User.StationId, 3); !ok {
		this.Abort("您没有该权限!")
	}
	this.TplName = "workspace/workspace_order_handing.html"
	var (
		is_order bool
		oc       models.OrderCredit
		err      error
		uid      int
		loanId   int
		timeDiff float64
		num      = 1
	)
	id := this.User.Id
	loginState := this.User.LoginState
	displayName := this.User.DisplayName
	defer func() {
		if err != nil && err.Error() != utils.ErrNoRow() {
			return
		}
		if loanId <= 0 {
			return
		}
		ubi, err := models.QueryUsersBaseInfo(uid) //用户基本信息
		if err != nil && err.Error() != utils.ErrNoRow() {
			cache.RecordLogs(this.User.Id, 0, this.User.Name, this.User.DisplayName, "查询用户基本信息出错", "workspace/OrderHanding", err.Error(), this.Ctx.Input)
			this.Abort("查询用户基本信息出错" + err.Error())
		}
		if ubi == nil {
			//用户数据异常
			this.Abort("用户信息不存在,请检查数据")
		}
		ubi.IdCard = utils.IdCardFilter(ubi.IdCard)
		oc, err = models.QueryOrderByUid(loanId)
		if err != nil {
			cache.RecordLogs(this.User.Id, 0, this.User.Name, this.User.DisplayName, "订单信息失败", "workspace/OrderHanding", err.Error(), this.Ctx.Input)
			this.Abort("订单信息失败")
		}
		list, err := models.QueryLoanAduitRecord(uid)
		if err != nil {
			cache.RecordLogs(this.User.Id, 0, this.User.Name, this.User.DisplayName, "获取日志信息失败", "workspace/OrderHanding", err.Error(), this.Ctx.Input)
			this.Abort("获取日志信息失败")
		}
		this.Data["phone_info"] = utils.QueryLocating(ubi.Account)
		this.Data["timeDiff"] = timeDiff
		this.Data["is_order"] = is_order
		this.Data["loginState"] = loginState
		this.Data["oc"] = oc
		this.Data["ubi"] = ubi
		this.Data["list"] = list
	}()
	timeDiff, loanId, uid, err = services.CheckHanding(id, displayName)
	if err != nil && err.Error() != utils.ErrNoRow() {
		cache.RecordLogs(this.User.Id, 0, this.User.Name, this.User.DisplayName, "校验缓存订单失败", "workspace/OrderHanding", err.Error(), this.Ctx.Input)
		this.Abort("校验缓存订单失败")

	}
	if uid > 0 && loanId > 0 {
		is_order = true
		return
	}
	if loginState == "OFFLINE" {
		this.Data["is_order"] = false
		this.Data["loginState"] = loginState
		return
	}
	oc, err = services.Order(num)
	if err != nil && err.Error() != utils.ErrNoRow() {
		cache.RecordLogs(this.User.Id, 0, this.User.Name, this.User.DisplayName, "获取预约订单失败", "workspace/OrderHanding", err.Error(), this.Ctx.Input)
		this.Abort("获取预约订单失败")
	}
	if oc.Id > 0 && oc.Uid > 0 {
		uid = oc.Uid
		loanId = oc.Id
		err = services.InsertOrderCache(id, loanId, uid, this.User.DisplayName, "HANDING")
		if err != nil {
			cache.RecordLogs(this.User.Id, 0, this.User.Name, this.User.DisplayName, "订单缓存失败", "workspace/OrderHanding", err.Error(), this.Ctx.Input)
			this.Abort("订单缓存失败")
		}
		err = models.UpdateAlloctionTime(loanId)
		if err != nil {
			cache.RecordLogs(this.User.Id, 0, this.User.Name, this.User.DisplayName, "分配时间失败", "workspace/OrderHanding", err.Error(), this.Ctx.Input)
			this.Abort("分配时间失败")
		}
		timeDiff = 45 * 60
		is_order = true
		content := "订单分配给" + "【" + this.User.DisplayName + "】"
		err = models.AddLoanAduitRecord(uid, loanId, content)
		return
	} else {
		this.Data["timeDiff"] = timeDiff
		this.Data["is_order"] = false
		this.Data["loginState"] = loginState
		return
	}
}

//============================================处理操作===================================================
//信审人员进行订单处理(1.通过 2.关闭 3.退回 )
func (this *WorkspaceController) OrderOpera() {
	var err error
	var oc models.OrderCredit
	var resp models.BaseResponse
	var req models.LoanApprovalRequest
	req.Ret = 200
	var inqueueTime string
	oc.Id, err = this.GetInt("loanId")
	if err != nil {
		this.Data["json"] = map[string]interface{}{"ret": 304, "err": err.Error(), "msg": "LoanId获取失败~"}
		cache.RecordLogs(this.User.Id, 0, this.User.Name, this.User.DisplayName, "LoanId获取失败", "workspace/OrderOpera", err.Error(), this.Ctx.Input)
		return
	}
	defer func() {
		this.Data["json"] = resp
		this.ServeJSON()
	}()
	orderState := this.GetString("orderState")
	oc.Mark = this.GetString("remark")
	uid, err := this.GetInt("userId")
	oc.Uid = uid
	if err != nil {
		this.Data["json"] = map[string]interface{}{"ret": 304, "err": err.Error(), "msg": "uid获取失败~"}
		cache.RecordLogs(this.User.Id, 0, this.User.Name, this.User.DisplayName, "uid获取失败", "workspace/OrderOpera", err.Error(), this.Ctx.Input)
		return
	}
	flag := cache.GetCacheLoanOften(oc.Id)
	if flag {
		resp.Ret = 304
		resp.Msg = "请勿频繁提交"
		return
	}
	if orderState == "OUTQUEUE" {
		inqueueTime = this.GetString("inqueueTime")
		oc.InqueueType, err = this.GetInt("inqueueType")
		oc.OrderState = orderState
		if err != nil {
			resp.Ret = 304
			resp.Err = err.Error()
			resp.Msg = "inqueueType获取失败~"
			cache.RecordLogs(this.User.Id, 0, this.User.Name, this.User.DisplayName, "inqueueType获取失败", "workspace/OrderOpera", err.Error(), this.Ctx.Input)
			return
		}
	}
	orderCredit, err := models.QueryOrderByUid(oc.Id)
	if err != nil {
		resp.Ret = 304
		resp.Err = err.Error()
		resp.Msg = "更新状态失败"
		cache.RecordLogs(this.User.Id, 0, this.User.Name, this.User.DisplayName, "更新状态失败", "workspace/OrderOpera", err.Error(), this.Ctx.Input)
		return
	}
	timeDiff := time.Now().Sub(orderCredit.Atime).Minutes()
	timeOut := 45.00
	if timeDiff >= timeOut {
		//delete redis(处理后清楚缓存)
		err = cache.DeleteOrderRedis(uid, this.User.Id)
		if err != nil {
			resp.Ret = 304
			resp.Err = err.Error()
			resp.Msg = "用户清除处理缓存失败"
			cache.RecordLogs(this.User.Id, 0, this.User.Name, this.User.DisplayName, "用户清除处理缓存失败", "workspace/OrderOpera", err.Error(), this.Ctx.Input)
			return
		}
		this.Data["json"] = map[string]interface{}{"ret": 302, "msg": "处理超时"}
		return
	}
	user, err := models.QueryUserByUid(uid)
	if err != nil {
		resp.Ret = 304
		resp.Err = err.Error()
		resp.Msg = "用户获取失败"
		cache.RecordLogs(this.User.Id, 0, this.User.Name, this.User.DisplayName, "用户获取失败", "workspace/OrderOpera", err.Error(), this.Ctx.Input)
		return
	}
	if orderState == "PASS" {
		oc.OrderState = orderState
		params := map[string]interface{}{
			"Uid":        uid,
			"Account":    user.Account,
			"LoanId":     oc.Id,
			"Operator":   "accept",
			"OperatorId": this.User.Id,
			"Remark":     oc.Mark,
		}
		b, err := services.PostApi("/loan/approval", params)
		err = json.Unmarshal(b, &req)
		if err != nil {
			resp.Ret = 304
			resp.Msg = "通过api请求失败"
			resp.Err = err.Error()
			cache.RecordLogs(this.User.Id, 0, this.User.Name, this.User.DisplayName, "通过api请求失败", "workspace/OrderOpera", err.Error(), this.Ctx.Input)
			return
		}
	}
	if orderState == "PAUSE" || orderState == "CLOSE" || orderState == "CANCEL" {
		oc.OrderState = orderState
		params := map[string]interface{}{
			"Uid":        uid,
			"Account":    user.Account,
			"LoanId":     oc.Id,
			"Operator":   "refuse",
			"OperatorId": this.User.Id,
			"Remark":     oc.Mark,
		}
		b, err := services.PostApi("/loan/approval", params)
		err = json.Unmarshal(b, &req)
		if err != nil {
			resp.Ret = 304
			resp.Err = err.Error()
			resp.Msg = "拒绝api请求失败"
			cache.RecordLogs(this.User.Id, 0, this.User.Name, this.User.DisplayName, "拒绝api请求失败", "workspace/OrderOpera", err.Error(), this.Ctx.Input)
			return
		}
	}
	is_req := 0
	if req.Ret != 200 {
		oc.OrderState = "OUTQUEUE"
		oc.InqueueType = 0
		is_req = 1
	}
	oc.CreditOperator = this.User.DisplayName
	oc.Content = utils.AddRemark(oc.OrderState, oc.CreditOperator)
	err = models.UpdateHandleMessages(oc, inqueueTime, is_req)
	if err != nil {
		resp.Ret = 304
		resp.Err = err.Error()
		resp.Msg = "更新数据库失败"
		cache.RecordLogs(this.User.Id, 0, this.User.Name, this.User.DisplayName, "更新数据库失败", "workspace/OrderOpera", err.Error(), this.Ctx.Input)
		return
	}
	//delete redis(处理后清除缓存)
	err = cache.DeleteOrderRedis(uid, this.User.Id)
	if err != nil {
		resp.Ret = 304
		resp.Err = err.Error()
		resp.Msg = "用户清除处理缓存失败"
		cache.RecordLogs(this.User.Id, 0, this.User.Name, this.User.DisplayName, "用户清除处理缓存失败", "workspace/OrderOpera", err.Error(), this.Ctx.Input)
		return
	}
	if orderState == "OUTQUEUE" {
		resp.Ret = 200
		resp.Msg = "处理成功"

	} else {
		resp.Ret = req.Ret
		resp.Err = req.Err
		resp.Msg = req.Msg
	}
	err = utils.Rc.Put(utils.CacheKeyLoanOften+"_"+strconv.Itoa(oc.Id), resp, utils.RedisCacheTime_LoanOften)
	if err != nil {
		resp.Ret = 304
		resp.Msg = "缓存提交信息失败"
		resp.Err = err.Error()
		return
	}
}

//订单处理操作
func (this *WorkspaceController) LoanOperation() {
	var resp models.BaseResponse
	loanId, err := this.GetInt("loanId")
	defer func() {
		this.Data["json"] = resp
		this.ServeJSON()
	}()
	if loanId <= 0 {
		resp.Ret = 304
		resp.Msg = "loanId传参失败"
		cache.RecordLogs(this.User.Id, 0, this.User.Name, this.User.DisplayName, "loanId传参失败", "workspace/LoanOperation", err.Error(), this.Ctx.Input)
		return
	}
	orderState := this.GetString("orderState")
	remark := this.GetString("remark")
	uid, err := this.GetInt("userId")
	flag := cache.GetCacheLoanOften(loanId)
	if flag {
		resp.Ret = 304
		resp.Msg = "请勿频繁提交"
		return
	}
	if ok := cache.GetCacheLOCKRecheck(uid, loanId); ok {
		resp.Ret = 304
		resp.Msg = "已重检"
		return
	}
	operationName := ""
	oper := ""
	var result int
	switch orderState {
	case "PASS":
		operationName = "'通过'"
		oper = "accept"
		result = 1
	case "CANCEL":
		operationName = "'正常关闭'"
		oper = "refuse"
		result = 2
	case "PAUSE":
		operationName = "'关闭30天'"
		oper = "refuse30"
		result = 3
	case "CLOSE":
		operationName = "'永久关闭'"
		oper = "refuse_forever"
		result = 4
	case "OUTQUEUE":
		operationName = "'退回'"
	}
	orderCredit, err := models.QueryOrderByUid(loanId)
	if err != nil {
		resp.Ret = 304
		resp.Err = err.Error()
		resp.Msg = "查询订单状态失败"
		cache.RecordLogs(this.User.Id, 0, this.User.Name, this.User.DisplayName, "查询订单状态失败", "workspace/LoanOperation", err.Error(), this.Ctx.Input)
		return
	}

	if orderCredit.State != "CONFIRM" {
		resp.Ret = 304
		resp.Msg = "该订单已处理"
	}
	timeDiff := time.Now().Sub(orderCredit.Atime).Minutes()
	if timeDiff >= 45.00 {
		//delete redis(处理后清除缓存)
		err = cache.DeleteOrderRedis(uid, this.User.Id)
		if err != nil {
			resp.Ret = 304
			resp.Err = err.Error()
			resp.Msg = "清除处理缓存失败"
			cache.RecordLogs(this.User.Id, 0, this.User.Name, this.User.DisplayName, "清除处理缓存失败", "workspace/LoanOperation", err.Error(), this.Ctx.Input)
			return
		}
		resp.Ret = 304
		resp.Msg = "处理超时"
		return
	}
	user, err := models.QueryUserByUid(uid)
	if err != nil {
		resp.Ret = 304
		resp.Err = err.Error()
		resp.Msg = "用户信息获取失败"
		cache.RecordLogs(this.User.Id, 0, this.User.Name, this.User.DisplayName, "用户信息获取失败", "workspace/LoanOperation", err.Error(), this.Ctx.Input)
		return
	}
	operator := this.User.DisplayName
	if orderState == "PASS" || orderState == "PAUSE" || orderState == "CLOSE" || orderState == "CANCEL" {
		params := map[string]interface{}{
			"Uid":        uid,
			"Account":    user.Account,
			"LoanId":     loanId,
			"Operator":   oper,
			"OperatorId": this.User.Id,
			"Remark":     remark,
		}
		//请求api
		b, err := services.PostApi(utils.Loan_Approval, params)
		beego.Info(string(b), err)
		if err != nil {
			resp.Ret = 304
			resp.Msg = operationName + "api请求失败"
			resp.Err = err.Error()
			cache.RecordLogs(this.User.Id, 0, this.User.Name, this.User.DisplayName, operationName+"api请求失败", "workspace/LoanOperation", err.Error(), this.Ctx.Input)
			return
		}
		err = json.Unmarshal(b, &resp)
		if err != nil {
			resp.Ret = 304
			resp.Msg = "解析" + operationName + "api请求返回数据失败"
			resp.Err = err.Error()
			cache.RecordLogs(this.User.Id, 0, this.User.Name, this.User.DisplayName, "解析"+operationName+"api请求返回数据失败", "workspace/LoanOperation", err.Error(), this.Ctx.Input)
			return
		}
		if resp.Ret != 200 { //api返回失败，置为出列需手动入列
			models.UpdateLoanState(loanId, 0, "OUTQUEUE")
		}
		err = models.UpdateStateByShutWithPass(uid, loanId, operator, remark, orderState)
		if err != nil {
			resp.Ret = 304
			resp.Err = err.Error()
			resp.Msg = "更新数据库失败"
			cache.RecordLogs(this.User.Id, 0, this.User.Name, this.User.DisplayName, "更新数据库失败", "workspace/LoanOperation", err.Error(), this.Ctx.Input)
			return
		}
	} else if orderState == "OUTQUEUE" {
		inqueueTime := this.GetString("inqueueTime")
		inqueueType, err := this.GetInt("inqueueType")
		err = models.UpdateStateByOutqueue(uid, loanId, inqueueType, result, inqueueTime, operator, remark, orderState)
		if err != nil {
			resp.Ret = 304
			resp.Err = err.Error()
			resp.Msg = operationName + "操作失败"
			cache.RecordLogs(this.User.Id, 0, this.User.Name, this.User.DisplayName, operationName+"操作失败", "workspace/LoanOperation", err.Error(), this.Ctx.Input)
			return
		}
	}
	//delete redis(处理后清除缓存)
	err = cache.DeleteOrderRedis(uid, this.User.Id)
	if err != nil {
		resp.Ret = 304
		resp.Err = err.Error()
		resp.Msg = "清除缓存失败"
		cache.RecordLogs(this.User.Id, 0, this.User.Name, this.User.DisplayName, "清除缓存失败", "workspace/LoanOperation", err.Error(), this.Ctx.Input)
		return
	}
	err = utils.Rc.Put(utils.CacheKeyLoanOften+"_"+strconv.Itoa(loanId), resp, utils.RedisCacheTime_LoanOften)
	if err != nil {
		resp.Ret = 304
		resp.Msg = "删除控制频繁提交缓存失败"
		resp.Err = err.Error()
		return
	}
	resp.Ret = 200
	resp.Msg = "处理成功"
}

// @Title 还款计划
func (this *WorkspaceController) QueryLoanPlan() {
	defer this.ServeJSON()
	loanId, err := this.GetInt("loanId")
	if err != nil {
		this.Data["json"] = map[string]interface{}{"ret": 304, "err": err.Error(), "msg": "repayId获取失败~"}
		cache.RecordLogs(this.User.Id, 0, this.User.Name, this.User.DisplayName, "repayId获取失败", "workspace/QueryLoanPlan", err.Error(), this.Ctx.Input)
		return
	}
	list, err := models.GetRepaymentScheduleByRepayId(loanId)
	if err != nil {
		this.Data["json"] = map[string]interface{}{"ret": 304, "err": err.Error(), "msg": "查询还款计划失败"}
		cache.RecordLogs(this.User.Id, 0, this.User.Name, this.User.DisplayName, "增加历史记录失败", "workspace/QueryLoanPlan", err.Error(), this.Ctx.Input)
	} else {
		this.Data["json"] = map[string]interface{}{"ret": 200, "data": list, "msg": "查询还款计划成功"}
	}
}

//订单查看
func (this *WorkspaceController) OrderLook() {
	this.IsNeedTemplate()
	uid, err := this.GetInt("uid")
	if err != nil {
		cache.RecordLogs(this.User.Id, 0, this.User.Name, this.User.DisplayName, "uid获取失败", "workspace/OrderLook", err.Error(), this.Ctx.Input)
		this.Abort("uid获取失败")
		return
	}
	loanId, err := this.GetInt("loanId")
	if err != nil {
		cache.RecordLogs(this.User.Id, 0, this.User.Name, this.User.DisplayName, "loanId获取失败", "workspace/OrderLook", err.Error(), this.Ctx.Input)
		this.Abort("loanId获取失败")
		return
	}
	ubi, err := models.QueryUsersBaseInfo(uid) //用户基本信息
	if err != nil && err.Error() != utils.ErrNoRow() {
		cache.RecordLogs(this.User.Id, 0, this.User.Name, this.User.DisplayName, "查询用户基本信息出错", "workspace/OrderLook", err.Error(), this.Ctx.Input)
		this.Abort("查询用户基本信息出错" + err.Error())
	}
	oc, err := models.QueryOrderByUid(loanId)
	if err != nil {
		cache.RecordLogs(this.User.Id, 0, this.User.Name, this.User.DisplayName, "订单信息失败", "workspace/OrderLook", err.Error(), this.Ctx.Input)
		this.Abort("订单信息失败")
	}
	if ubi == nil {
		this.Abort("用户信息不存在,请检查数据")
	}
	list, err := models.QueryLoanAduitRecord(uid)
	if err != nil {
		cache.RecordLogs(this.User.Id, 0, this.User.Name, this.User.DisplayName, "获取日志信息失败", "workspace/OrderLook", err.Error(), this.Ctx.Input)
		this.Abort("获取日志信息失败")
	}
	index, _ := this.GetInt("index", 3)
	this.Data["index"] = index
	this.Data["oc"] = oc
	this.Data["ubi"] = ubi
	this.Data["list"] = list
	this.TplName = "workspace/workspace_order_look.html"
}

package controllers

import (
	"encoding/json"
	"fenqi_v1/cache"
	"fenqi_v1/models"
	"fenqi_v1/services"
	"fenqi_v1/utils"
	nethttp "net/http"
	"os"
	"strconv"
	"strings"
	"time"
	"fmt"
)

type FinanceDataController struct {
	BaseController
}

//业务数据
func (this *FinanceDataController) BusinessData() {
	this.IsNeedTemplate()
	this.TplName = "financeData/business_data.html"
}

//业务资金结算
func (this *FinanceDataController) BusinessFundSettlement() {
	this.IsNeedTemplate()
	this.TplName = "financeData/business_fund_settlement.html"
}

//业务收入表
func (this *FinanceDataController) BusinessIncomeStatement() {
	this.IsNeedTemplate()
	this.TplName = "financeData/business_income_statement.html"
}

//运营成本分析
func (this *FinanceDataController) OperationalCostAnalysis() {
	this.IsNeedTemplate()
	this.TplName = "financeData/operational_cost_analysis.html"
}

//还款数据
func (this *FinanceDataController) PaymentData() {
	this.IsNeedTemplate()
	this.TplName = "financeData/payment_data.html"
}

//逾期统计
func (this *FinanceDataController) OverdueStatistical() {
	this.IsNeedTemplate()
	this.TplName = "financeData/overdue_statistical.html"
}

//催收数据统计
func (this *FinanceDataController) CollectDataStatistics() {
	this.IsNeedTemplate()
	this.TplName = "financeData/collect_data_statistics.html"
}

//累计逾期催收数据
func (this *FinanceDataController) AccumulatedOverdueCollectionData() {
	this.IsNeedTemplate()
	this.TplName = "financeData/accumulated_overdue_collection_data.html"
}

//周分析数据
func (this *FinanceDataController) WeekAnalyzeData() {
	this.IsNeedTemplate()
	this.TplName = "financeData/week_analyze_data.html"
}

//回款管理
func (this *FinanceDataController) RepayManagement() {
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
		condition += " AND ca.account = ?"
		pars = append(pars, account)
	}
	userName := this.GetString("user_name") //姓名
	if userName != "" {
		condition += " AND ca.verifyrealname = ?"
		pars = append(pars, userName)
	}
	create_time := this.GetString("submit_time") //还款时间
	if create_time != "" {
		create_times := strings.Split(create_time, "~")
		startloantime := create_times[0] + " 00:00:00"
		endloantime := create_times[1] + " 23:59:59"
		condition += " AND su.create_time >= ? AND su.create_time <= ?"
		pars = append(pars, startloantime)
		pars = append(pars, endloantime)
	}

	ReturnList, err := models.ReturnRecordList(utils.StartIndex(pageNum, pageSize), pageSize, true, condition, pars...)
	if err != nil {
		cache.RecordLogs(this.User.Id, 0, this.User.Name, this.User.DisplayName, "财务管理", "财务管理/回款管理list", err.Error(), this.Ctx.Input)
		this.Abort("查询回款管理失败")
		return
	}
	count, err := models.ReturnRecordListCount(condition, pars...)
	if err != nil {
		cache.RecordLogs(this.User.Id, 0, this.User.Name, this.User.DisplayName, "财务管理", "财务管理/回款管理list", err.Error(), this.Ctx.Input)
		this.Abort("获取/回款管理总数异常")
		return
	}
	pageCount := utils.PageCount(count, pageSize)
	this.Data["list"] = ReturnList
	this.Data["currPage"] = pageNum
	this.Data["count"] = count
	this.Data["pageSize"] = pageSize
	this.Data["pageCount"] = pageCount
	this.TplName = "financeData/repaymanagement.html"
}

//还款管理
func (this *FinanceDataController) ReturnManagement() {
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
		condition += " AND ca.account = ?"
		pars = append(pars, account)
	}
	userName := this.GetString("user_name") //姓名
	if userName != "" {
		condition += " AND ca.verifyRealName = ?"
		pars = append(pars, userName)
	}
	order_number := this.GetString("order_number") //订单号
	if order_number != "" {
		condition += " AND su.order_number = ?"
		pars = append(pars, order_number)
	}
	create_time := this.GetString("submit_time") //还款时间
	if create_time != "" {
		create_times := strings.Split(create_time, "~")
		startloantime := create_times[0] + " 00:00:00"
		endloantime := create_times[1] + " 23:59:59"
		condition += " AND su.create_time >= ? AND su.create_time <= ?"
		pars = append(pars, startloantime)
		pars = append(pars, endloantime)
	}
	oid_paybill, _ := this.GetInt("channel") //还款渠道
	if oid_paybill != 0 {
		condition += " AND su.channel = ?"
		pars = append(pars, oid_paybill)
	}
	operator_id, _ := this.GetInt("operator", -1) //还款方式
	if operator_id != -1 {
		if operator_id == 0 { //代扣
			condition += " AND su.operator_id NOT IN(-1,1)"
		} else { //主动还款
			condition += " AND su.operator_id IN(-1,1)"
		}
	}
	state := this.GetString("ret") //还款结果
	if state != "" {
		condition += " AND su.state = ?"
		pars = append(pars, state)
	}
	ReturnList, err := models.RepayRecordList(utils.StartIndex(pageNum, pageSize), pageSize,true, condition, pars...)
	if err != nil {
		cache.RecordLogs(this.User.Id, 0, this.User.Name, this.User.DisplayName, "财务管理", "财务管理/还款管理list", err.Error(), this.Ctx.Input)
		this.Abort("查询还款管理失败")
		return
	}
	count, err := models.RepayRecordListCount(condition, pars...)
	if err != nil {
		cache.RecordLogs(this.User.Id, 0, this.User.Name, this.User.DisplayName, "财务管理", "财务管理/还款管理list", err.Error(), this.Ctx.Input)
		this.Abort("获取还款管理总数异常")
		return
	}
	pageCount := utils.PageCount(count, pageSize)
	fmt.Println(pageCount)
	this.Data["list"] = ReturnList
	this.Data["currPage"] = pageNum
	this.Data["count"] = count
	this.Data["pageSize"] = pageSize
	this.Data["pageCount"] = pageCount
	this.TplName = "financeData/returnmanagement.html"
}

//报告购买记录
func (this *FinanceDataController) BuyRecord() {
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
		condition += " AND a.account = ?"
		pars = append(pars, account)
	}
	userName := this.GetString("user_name") //姓名
	if userName != "" {
		condition += " AND a.verifyrealname = ?"
		pars = append(pars, userName)
	}
	startSubmitTime := this.GetString("start_submit_time") //提交时间
	if startSubmitTime != "" {
		condition += ` AND b.create_time >= ?`
		pars = append(pars, startSubmitTime)
	}
	endSubmitTime := this.GetString("end_submit_time") //提交时间
	if endSubmitTime != "" {
		condition += ` AND b.create_time <= ?`
		pars = append(pars, endSubmitTime)
	}
	credit_card_name := this.GetString("credit_card_name") //报告类型
	if credit_card_name != "" {
		condition += " AND b.credit_card_name = ?"
		pars = append(pars, credit_card_name)
	}
	state := this.GetString("state") //购买结果
	if state != "" {
		if state == "SUCCESS" {
			condition += " AND (b.state = ? OR b.state = ? OR (b.state = ? AND b.loan_id !=0))"
			pars = append(pars, "SUCCESS", "LOCKED", "USELESS")
		} else if state == "FAIL" {
			condition += " AND (b.state = ? OR (b.state = ? AND b.loan_id = 0))"
			pars = append(pars, "FAIL", "USELESS")
		} else {
			condition += " AND b.state = ?"
			pars = append(pars, "CONFIRM")
		}
	}
	buy_state := this.GetString("buy_state")
	if buy_state != "" {
		condition += " AND b.pay_method = ?"
		pars = append(pars, buy_state)
	}
	if condition == "" { //默认提交时间1天内
		condition += ` AND b.create_time >= ? AND b.create_time <= ?`
		now := time.Now().Format("2006-01-02")
		startSubmitTime = now + " 00:00:00"
		endSubmitTime = now + " 23:59:59"
		pars = append(pars, startSubmitTime)
		pars = append(pars, endSubmitTime)
	}
	buy_list, err := models.GetCreditReport(utils.StartIndex(pageNum, pageSize), pageSize, true, condition, pars...)
	if err != nil {
		cache.RecordLogs(this.User.Id, 0, this.User.Name, this.User.DisplayName, "财务管理", "财务管理/报告购买记录", err.Error(), this.Ctx.Input)
		this.Abort("查询还款管理失败")
		return
	}
	count, err := models.GetUsersCreditReportCount(condition, pars...)
	if err != nil {
		cache.RecordLogs(this.User.Id, 0, this.User.Name, this.User.DisplayName, "财务管理", "财务管理/报告购买记录", err.Error(), this.Ctx.Input)
		this.Abort("获取还款管理总数异常")
		return
	}
	for k, v := range buy_list {
		if v.State == "SUCCESS" || v.State == "LOCKED" || (v.State == "USELESS" && v.LoanId != 0) {
			buy_list[k].StateStr = "成功"
		} else if v.State == "FAIL" || (v.State == "USELESS" && v.LoanId == 0) {
			buy_list[k].StateStr = "失败"
		} else if v.State == "CONFIRM" {
			buy_list[k].StateStr = "处理中"
		}
	}
	pageCount := utils.PageCount(count, pageSize)
	this.Data["startSubmitTime"] = startSubmitTime
	this.Data["endSubmitTime"] = endSubmitTime
	this.Data["list"] = buy_list
	this.Data["currPage"] = pageNum
	this.Data["count"] = count
	this.Data["pageSize"] = pageSize
	this.Data["pageCount"] = pageCount
	this.TplName = "financeData/financeData_report_list.html"
}

//关闭个人征信订单
func (this *FinanceDataController) ClosePerCreditReport() {
	defer func() {
		this.ServeJSON()
	}()
	id, _ := this.GetInt("id")
	if id <= 0 {
		cache.RecordLogs(this.User.Id, 0, this.User.Name, this.User.DisplayName, "id参数传递错误", "关闭个人征信订单/ClosePerCreditReport", "", this.Ctx.Input)
		this.Data["json"] = map[string]interface{}{"ret": 200, "msg": "这个订单不存在"}
		return
	}
	repdata, err := models.GeUidStateFromUsersCreditReport(id)
	if err != nil && err.Error() != utils.ErrNoRow() {
		cache.RecordLogs(this.User.Id, 0, this.User.Name, this.User.DisplayName, "查询用户银行卡绑定情况出错", "关闭个人征信订单/ClosePerCreditReport", err.Error(), this.Ctx.Input)
		this.Data["json"] = map[string]interface{}{"ret": 200, "msg": "这个订单不存在"}
		return
	}
	if repdata.State == "CONFIRM" {
		err = models.UpdateUidStateFromUsersCreditReport(id, "FAIL")
		this.Data["json"] = map[string]interface{}{"ret": 200, "msg": repdata.Remark}
		return
	}
	this.Data["json"] = map[string]interface{}{"ret": 403, "msg": "此人没有购买征信报告"}
	return
}

//报告购买记录导出
func (this *FinanceDataController) FinanceDataToExcel() {
	pars := []interface{}{}
	condition := ""
	account := this.GetString("account") //手机号
	if account != "" {
		condition += " AND a.account = ?"
		pars = append(pars, account)
	}
	userName := this.GetString("user_name") //姓名
	if userName != "" {
		condition += " AND a.verifyrealname = ?"
		pars = append(pars, userName)
	}
	startSubmitTime := this.GetString("startDate_submit") //提交时间
	if startSubmitTime != "" {
		condition += ` AND b.create_time >= ?`
		pars = append(pars, startSubmitTime)
	}
	endSubmitTime := this.GetString("endDate_submit") //提交时间
	if endSubmitTime != "" {
		condition += ` AND b.create_time <= ?`
		pars = append(pars, endSubmitTime)
	}
	credit_card_name := this.GetString("credit_card_name") //报告类型
	if credit_card_name != "" {
		condition += " AND b.credit_card_name = ?"
		pars = append(pars, credit_card_name)
	}
	state := this.GetString("state") //购买结果
	if state != "" {
		if state == "SUCCESS" {
			condition += " AND (b.state = ? OR b.state = ? OR (b.state = ? AND b.loan_id !=0))"
			pars = append(pars, "SUCCESS", "LOCKED", "USELESS")
		} else if state == "FAIL" {
			condition += " AND (b.state = ? OR (b.state = ? AND b.loan_id = 0))"
			pars = append(pars, "FAIL", "USELESS")
		} else {
			condition += " AND b.state = ?"
			pars = append(pars, "CONFIRM")
		}
	}
	buy_state := this.GetString("buy_state")
	if buy_state != "" {
		condition += " AND b.pay_method = ?"
		pars = append(pars, buy_state)
	}

	buy_list, err := models.GetCreditReport(0, 0, false, condition, pars...)
	if err != nil {
		cache.RecordLogs(this.User.Id, 0, this.User.Name, this.User.DisplayName, "财务管理", "财务管理/报告购买记录导出", err.Error(), this.Ctx.Input)
		this.Abort("查询报告购买记录失败")
		return
	}
	for k, v := range buy_list {
		if v.State == "SUCCESS" || v.State == "LOCKED" || (v.State == "USELESS" && v.LoanId != 0) {
			buy_list[k].StateStr = "成功"
		} else if v.State == "FAIL" || (v.State == "USELESS" && v.LoanId == 0) {
			buy_list[k].StateStr = "失败"
		} else if v.State == "CONFIRM" {
			buy_list[k].StateStr = "处理中"
		}
	}

	exl := [][]string{{"手机号", "姓名", "订单号", "购买时间", "报告类型", "支付金额", "支付方式", "购买结果"}}
	colWidth := []float64{20.0, 20.0, 35.0, 30.0, 20.0, 20.0, 20.0, 20.0}
	for _, v := range buy_list {
		payMethod := ""
		if v.PayMethod == 1 {
			payMethod = "连连"
		} else {
			payMethod = "支付宝"
		}
		stateStr := ""
		if v.State == "SUCCESS" || v.State == "LOCKED" || (v.State == "USELESS" && v.LoanId != 0) {
			stateStr = "成功"
		} else if v.State == "FAIL" || (v.State == "USELESS" && v.LoanId == 0) {
			stateStr = "失败"
		} else if v.State == "CONFIRM" {
			stateStr = "处理中"
		}
		exp := []string{
			v.Account,
			v.VerifyRealName,
			v.OrderNumber,
			v.CreateTime.Format("2006-01-02 15:04:05"),
			v.CreditCardName,
			strconv.FormatFloat(v.PayPrice, 'f', -1, 64),
			payMethod,
			stateStr + " " + v.Remark,
		}
		exl = append(exl, exp)
	}
	tablename := "报告购买记录"
	filename, err := utils.ExportToExcel(exl, colWidth, tablename)
	if err != nil {
		cache.RecordLogs(this.User.Id, 0, this.User.Name, this.User.DisplayName, "", "保存文件错误", err.Error(), this.Ctx.Input)
		this.Abort("保存文件错误" + err.Error())
		return
	}
	this.Ctx.Output.Header("Accept-Ranges", "bytes")
	this.Ctx.Output.Header("Content-Disposition", "attachment; filename="+filename)
	this.Ctx.Output.Header("Cache-Control", "must-revalidate, post-check=0, pre-check=0")
	this.Ctx.Output.Header("Pragma", "no-cache")
	this.Ctx.Output.Header("Expires", "0")
	nethttp.ServeFile(this.Ctx.ResponseWriter, this.Ctx.Request, filename)
	err = os.Remove(filename)
	if err != nil {
		cache.RecordLogs(this.User.Id, 0, this.User.Name, this.User.DisplayName, "", "删除文件错误", err.Error(), this.Ctx.Input)
	}
	return
}

//还款管理导出
func (this *FinanceDataController) ReturnManagementToExcel() {
	pars := []interface{}{}
	condition := ""
	account := this.GetString("account") //手机号
	if account != "" {
		condition += " AND ca.account = ?"
		pars = append(pars, account)
	}
	userName := this.GetString("user_name") //姓名
	if userName != "" {
		condition += " AND ca.verifyRealName = ?"
		pars = append(pars, userName)
	}
	order_number := this.GetString("order_number") //订单号
	if order_number != "" {
		condition += " AND su.order_number = ?"
		pars = append(pars, order_number)
	}
	create_time := this.GetString("submit_time") //还款时间
	if create_time != "" {
		create_times := strings.Split(create_time, "~")
		startloantime := create_times[0] + " 00:00:00"
		endloantime := create_times[1] + " 23:59:59"
		condition += " AND su.create_time >= ? AND su.create_time <= ?"
		pars = append(pars, startloantime)
		pars = append(pars, endloantime)
	}
	oid_paybill, _ := this.GetInt("channel") //还款渠道
	if oid_paybill != 0 {
		condition += " AND su.channel = ?"
		pars = append(pars, oid_paybill)
	}
	operator_id, _ := this.GetInt("operator", -1) //还款方式
	if operator_id != -1 {
		if operator_id == 0 { //代扣
			condition += " AND su.operator_id NOT IN(-1,1)"
		} else { //主动还款
			condition += " AND su.operator_id IN(-1,1)"
		}
	}
	state := this.GetString("ret") //还款结果
	if state != "" {
		condition += " AND su.state = ?"
		pars = append(pars, state)
	}
	returnList, err := models.RepayRecordList(0, 0, false, condition, pars...)
	if err != nil {
		cache.RecordLogs(this.User.Id, 0, this.User.Name, this.User.DisplayName, "财务管理", "财务管理/还款管理list导出", err.Error(), this.Ctx.Input)
		this.Abort("查询还款管理失败")
		return
	}
	exl := [][]string{{"手机号", "姓名", "还款金额", "还款时间", "还款渠道", "订单号", "第三方订单号", "还款方式", "还款结果"}}
	colWidth := []float64{20.0, 20.0, 20.0, 30.0, 20.0, 45.0, 45.0, 20.0, 20.0}
	for _, v := range returnList {
		createTime := "暂无"
		if v.CreateTime.Format("2006-01-02 15:04:05") != "0001-01-01 00:00:00" {
			createTime = v.CreateTime.Format("2006-01-02 15:04:05")
		}
		channel := ""
		switch v.Channel {
		case 1:
			channel = "连连支付"
			break
		case 2:
			channel = "支付宝支付"
			break
		case 3:
			channel = "合利宝代扣"
			break
		case 5:
			channel = "支付宝转账"
			break
		case 6:
			channel = "银行转账"
			break
		case 9:
			channel = "费用减免"
			break
		case 7:
			channel = "融宝代扣"
			break
		case 8:
			channel = "畅捷代扣"
			break
		case 10:
			channel = "先锋代扣"
			break
		}
		operatorId := ""
		if v.OperatorId == -1 || v.OperatorId == 1 {
			operatorId = "主动还款"
		} else {
			operatorId = "代扣"
		}
		state := ""
		if v.State == "SUCCESS" {
			state = "成功"
		} else if v.State == "FAIL" {
			state = "失败"
		} else if v.State == "CONFIRM" {
			state = "确认中"
		}
		exp := []string{
			v.Account,
			v.VerifyRealName,
			v.ReturnMoney,
			createTime,
			channel,
			v.OrderNumber,
			v.OidPayBill,
			operatorId,
			state,
		}
		exl = append(exl, exp)
	}
	tablename := "还款管理"
	filename, err := utils.ExportToExcel(exl, colWidth, tablename)
	if err != nil {
		cache.RecordLogs(this.User.Id, 0, this.User.Name, this.User.DisplayName, "", "保存文件错误", err.Error(), this.Ctx.Input)
		this.Abort("保存文件错误" + err.Error())
		return
	}
	this.Ctx.Output.Header("Accept-Ranges", "bytes")
	this.Ctx.Output.Header("Content-Disposition", "attachment; filename="+filename)
	this.Ctx.Output.Header("Cache-Control", "must-revalidate, post-check=0, pre-check=0")
	this.Ctx.Output.Header("Pragma", "no-cache")
	this.Ctx.Output.Header("Expires", "0")
	nethttp.ServeFile(this.Ctx.ResponseWriter, this.Ctx.Request, filename)
	err = os.Remove(filename)
	if err != nil {
		cache.RecordLogs(this.User.Id, 0, this.User.Name, this.User.DisplayName, "", "删除文件错误", err.Error(), this.Ctx.Input)
	}
	return
}

//回款管理导出
func (this *FinanceDataController) RepayManagementToExcel() {
	pars := []interface{}{}
	condition := ""
	account := this.GetString("account") //手机号
	if account != "" {
		condition += " AND ca.account = ?"
		pars = append(pars, account)
	}
	userName := this.GetString("user_name") //姓名
	if userName != "" {
		condition += " AND ca.verifyrealname = ?"
		pars = append(pars, userName)
	}
	create_time := this.GetString("create_time") //还款时间
	if create_time != "" {
		create_times := strings.Split(create_time, "~")
		startloantime := create_times[0] + " 00:00:00"
		endloantime := create_times[1] + " 23:59:59"
		condition += " AND su.create_time >= ? AND su.create_time <= ?"
		pars = append(pars, startloantime)
		pars = append(pars, endloantime)
	}

	returnList, err := models.ReturnRecordList(0, 0, false, condition, pars...)
	if err != nil {
		cache.RecordLogs(this.User.Id, 0, this.User.Name, this.User.DisplayName, "财务管理", "财务管理/回款管理list导出", err.Error(), this.Ctx.Input)
		this.Abort("查询回款管理失败")
		return
	}

	exl := [][]string{{"手机号", "姓名", "放款时间", "借款金额", "借款期限", "到账金额", "操作人"}}
	colWidth := []float64{20.0, 20.0, 30.0, 20.0, 20.0, 20.0, 20.0, 20.0}
	for _, v := range returnList {
		createTime := "暂无"
		if v.CreateTime.Format("2006-01-02 15:04:05") != "0001-01-01 00:00:00" {
			createTime = v.CreateTime.Format("2006-01-02 15:04:05")
		}
		exp := []string{
			v.Account,
			v.VerifyRealName,
			createTime,
			utils.Float64ToStrings(v.Money),
			strconv.Itoa(v.LoanTermCount),
			utils.Float64ToStrings(v.Rmoney),
			v.DisplayName,
		}
		exl = append(exl, exp)
	}
	tablename := "回款管理"
	filename, err := utils.ExportToExcel(exl, colWidth, tablename)
	if err != nil {
		cache.RecordLogs(this.User.Id, 0, this.User.Name, this.User.DisplayName, "", "保存文件错误", err.Error(), this.Ctx.Input)
		this.Abort("保存文件错误" + err.Error())
		return
	}
	this.Ctx.Output.Header("Accept-Ranges", "bytes")
	this.Ctx.Output.Header("Content-Disposition", "attachment; filename="+filename)
	this.Ctx.Output.Header("Cache-Control", "must-revalidate, post-check=0, pre-check=0")
	this.Ctx.Output.Header("Pragma", "no-cache")
	this.Ctx.Output.Header("Expires", "0")
	nethttp.ServeFile(this.Ctx.ResponseWriter, this.Ctx.Request, filename)
	err = os.Remove(filename)
	if err != nil {
		cache.RecordLogs(this.User.Id, 0, this.User.Name, this.User.DisplayName, "", "删除文件错误", err.Error(), this.Ctx.Input)
	}
	return
}

//支付宝还款导出EXCEL
func (c *FinanceDataController) AilirepaymentToExcel() {
	pars := []string{}
	condition := ""
	//处理结果
	deal_result := c.GetString("deal_result")
	if deal_result != "" {
		condition += " and remark2=?"
		pars = append(pars, deal_result)
	}
	//日期
	startDate := c.GetString("startDate")
	if startDate != "" {
		startDate = strings.Replace(startDate, "+", " ", 1)
		condition += " and return_time>=?"
		pars = append(pars, startDate)
	}
	endDate := c.GetString("endDate")
	if endDate != "" {
		endDate = strings.Replace(endDate, "+", " ", 1)
		condition += " and return_time<=?"
		pars = append(pars, endDate)
	}
	//流水号
	paybill := c.GetString("paybill")
	if paybill != "" {
		condition += " and oid_paybill=?"
		pars = append(pars, paybill)
	}
	//状态
	input_status := c.GetString("input_status")
	if input_status != "" {
		condition += " and remark1=?"
		pars = append(pars, input_status)
	}

	//var sumAmountIncome string = "0" //总还款金额
	list, err := models.GetAlipayRecordMoreThanZ(condition, false, 0, 0, pars...)
	if err != nil {
		cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "查询支付宝还款记录数据失败", "贷款管理/Ailirepayment", err.Error(), c.Ctx.Input)
	}
	//sumAmountIncome := models.GetAlipayRecordMoreThanZSumReturnMoney(condition, pars...)      //总还款金额
	//sumExtraRepayment := models.GetAlipayRecordMoreThanZSumExtraRepayment(condition, pars...) //总多余还款金额
	exl := [][]string{{"流水号", "转账时间", "转账金额", "备注", "状态", "处理结果", "多余还款"}}
	colWidth := []float64{45.0, 30.0, 20.0, 30.0, 30.0, 30.0, 20.0}
	for _, v := range list {
		exp := []string{
			v.OidPaybill,
			v.ReturnTime.Format("2006-01-02 15:04:05"),
			utils.Float64ToString(v.AmountIncome),
			v.Remark,
			v.Remark1,
			v.Remark2,
			utils.Float64ToString(v.ExtraRepayment),
		}
		exl = append(exl, exp)
	}
	filename, err := utils.ExportToExcel(exl, colWidth, "支付宝还款")
	if err != nil {
		cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "保存文件错误", err.Error(), c.Ctx.Input)
		c.Abort("保存文件错误" + err.Error())
	}
	c.Ctx.Output.Header("Accept-Ranges", "bytes")
	c.Ctx.Output.Header("Content-Disposition", "attachment; filename="+filename)
	c.Ctx.Output.Header("Cache-Control", "must-revalidate, post-check=0, pre-check=0")
	c.Ctx.Output.Header("Pragma", "no-cache")
	c.Ctx.Output.Header("Expires", "0")
	nethttp.ServeFile(c.Ctx.ResponseWriter, c.Ctx.Request, filename)
	err = os.Remove(filename)
	if err != nil {
		cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "删除文件错误", err.Error(), c.Ctx.Input)
	}
	return
}

//支付宝还款
func (c *FinanceDataController) Ailirepayment() {
	c.IsNeedTemplate()
	pageNum, _ := c.GetInt("page", 1) //分页信息（第几页）
	if pageNum < 1 {
		pageNum = 1
	}
	pageSize,_ := c.GetInt("pageSize")
	if pageSize < 1 {
		pageSize = 15
	}

	pars := []string{}
	condition := ""
	//处理结果
	deal_result := c.GetString("deal_result")
	if deal_result != "" {
		condition += " and remark2=?"
		pars = append(pars, deal_result)
	}
	//日期
	startDate := c.GetString("startDate")
	if startDate != "" {
		startDate = strings.Replace(startDate, "+", " ", 1)
		condition += " and return_time>=?"
		pars = append(pars, startDate)
	}
	endDate := c.GetString("endDate")
	if endDate != "" {
		endDate = strings.Replace(endDate, "+", " ", 1)
		condition += " and return_time<=?"
		pars = append(pars, endDate)
	}
	//流水号
	paybill := c.GetString("paybill")
	if paybill != "" {
		condition += " and oid_paybill=?"
		pars = append(pars, paybill)
	}
	//状态
	input_status := c.GetString("input_status")
	if input_status != "" {
		condition += " and remark1=?"
		pars = append(pars, input_status)
	}

	//var sumAmountIncome string = "0" //总还款金额
	list, err := models.GetAlipayRecordMoreThanZ(condition, true, utils.StartIndex(pageNum, pageSize), pageSize, pars...)
	if err != nil {
		cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "查询支付宝还款记录数据失败", "贷款管理/Ailirepayment", err.Error(), c.Ctx.Input)
	}
	sumAmountIncome := models.GetAlipayRecordMoreThanZSumReturnMoney(condition, pars...)      //总还款金额
	sumExtraRepayment := models.GetAlipayRecordMoreThanZSumExtraRepayment(condition, pars...) //总多余还款金额
	count := models.GetAlipayRecordCount(condition, pars...)
	pageCount := utils.PageCount(count, pageSize)
	c.Data["stationId"] = c.User.StationId
	c.Data["sumAmountIncome"] = sumAmountIncome
	c.Data["sumExtraRepayment"] = sumExtraRepayment
	c.Data["count"] = count
	c.Data["list"] = list
	c.Data["pageSize"] = pageSize
	c.Data["pageCount"] = pageCount
	c.Data["currPage"] = pageNum
	c.TplName = "financeData/ali_repay.html"
}

//修改处理结果
func (this *FinanceDataController) ModifyRemark() {
	defer this.ServeJSON()
	alipayRecordId, _ := this.GetInt("id")
	operatorIp := this.Ctx.Input.IP()
	newRemark := this.GetString("newRemark")
	extramoney, _ := this.GetFloat("money")
	err := models.UpdateAlipayRecordRemark(alipayRecordId, this.User.Id, newRemark, operatorIp, extramoney)
	if err != nil {
		cache.RecordLogs(this.User.Id, 0, this.User.Name, this.User.DisplayName, "修改处理结果有误", "支付宝还款/ModifyRemark", err.Error(), this.Ctx.Input)
		this.Data["json"] = map[string]interface{}{"ret": 304, "err": "修改处理结果失败" + err.Error()}
		return
	}
	this.Data["json"] = map[string]interface{}{"ret": 200, "msg": "更改处理结果成功!"}
}

//处理
func (c *FinanceDataController) GetRemark() {
	defer c.ServeJSON()
	results := c.GetString("resultCode")
	extramoney, _ := c.GetFloat("money")
	operatorId := c.User.Id
	id, _ := c.GetInt("id")
	if id < 1 {
		c.Data["json"] = map[string]interface{}{"ret": 304, "msg": "参数id有误"}
		return
	}
	if results == "多余还款" {
		amountIncome, _ := c.GetFloat("amountIncome")
		extramoney = amountIncome
	}
	operatorIp := c.Ctx.Input.IP()
	err := models.UpdateAlipayRecordRemark2(results, extramoney, operatorIp, id, operatorId)
	if err != nil {
		cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "处理失败", "贷款管理/GetRemark", err.Error(), c.Ctx.Input)
		models.UpdateAlipayIsDeal(3, id)
	} else {
		models.UpdateAlipayIsDeal(2, id)
	}
	c.Data["json"] = map[string]interface{}{"ret": 200, "msg": "更新成功"}
}

//还款
func (c *FinanceDataController) AliRepay() {
	c.IsNeedTemplate()
	defer c.ServeJSON()
	account := c.GetString("account")
	AlipayId, _ := c.GetInt("id")
	oidpaybill := c.GetString("oidpaybill")
	amountIncome, _ := c.GetFloat("amountincome")
	alirepay_time := c.GetString("alirepay_time")
	uid, err := models.GetUsersMetadataByaccount(account)
	if err != nil && err.Error() != "<QuerySeter> no row found" {
		cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "根据account查询用户信息失败", "贷款管理/GetRemark", err.Error(), c.Ctx.Input)
	}
	if uid == 0 {
		c.Data["json"] = map[string]interface{}{"ret": 304, "msg": "未找到该手机用户,请核实手机号"}
		return
	}
	repayment_schedule_id := models.VerfyLoan(account)
	if repayment_schedule_id == 0 {
		c.Data["json"] = map[string]interface{}{"ret": 304, "msg": "未找到需要处理的还款计划"}
		models.UpdateAlipayRecord2(3, AlipayId, "未找到需要处理的还款计划")
		return
	}
	//调用还款接口
	params := map[string]interface{}{
		"Uid":                 uid,
		"MobileType":          "YGFQ_SYS",
		"RepaymentScheduleId": repayment_schedule_id,
		"Channel":             5,
		"OperatorId":          c.User.Id,
		"ReturnMoney":         amountIncome,
		"Remark":              "支付宝还款",
		"OidPaybill":          oidpaybill,
		"RepayTime":           alirepay_time}
	b, err := services.PostApi(utils.Loan_Repayment, params)
	if err != nil {
		cache.RecordLogs(c.User.Id, repayment_schedule_id, c.User.Name, c.User.DisplayName, "还款失败，请求接口失败", "还款管理/GetRemark", err.Error(), c.Ctx.Input)
		c.Data["json"] = map[string]interface{}{"ret": 304, "msg": "还款失败，请求接口失败" + err.Error()}
		models.UpdateAlipayIsDeal(3, AlipayId)
		return
	}
	//添加日志记录
	// cache.RecordLogs(c.User.Id, list[0].RepaymentScheduleId, c.User.Name, c.User.DisplayName, "支付宝还款", "还款管理/GetRemark", string(b), c.Ctx.Input)
	var res models.ApiGeneralResponse
	json.Unmarshal(b, &res)
	if res.ErrTag == 1 {
		res.Msg = "还款金额不能大于本期应还金额"
	} else if res.ErrTag == 2 {
		res.Msg = "该笔还款计划正在还款中"
	} else if res.ErrTag == 3 {
		res.Msg = "获取用户详细信息失败"
	} else if res.ErrTag == 4 {
		res.Msg = "时间格式错误"
	} else if res.ErrTag == 6 {
		res.Msg = "获取还款计划信息失败"
	} else if res.ErrTag == 7 {
		res.Msg = "该笔还款计划已结清"
	} else if res.ErrTag == 9 {
		res.Msg = "还款失败"
	} else if res.ErrTag == 0 {
		res.Msg = "还款成功"
	}
	if res.Ret == 200 {
		models.UpdateAlipayRecord(2, AlipayId, c.User.Id, "还款成功", c.Ctx.Input.IP())
	} else {
		models.UpdateAlipayRecord(3, AlipayId, c.User.Id, res.Msg, c.Ctx.Input.IP())
	}
	c.Data["json"] = map[string]interface{}{"ret": res.Ret, "msg": res.Msg}
}

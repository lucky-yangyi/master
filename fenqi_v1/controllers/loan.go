package controllers

import (
	"fenqi_v1/cache"
	"fenqi_v1/models"
	"fenqi_v1/utils"
	"net/http"
	"os"
	"strconv"
	"time"
)

//用户管理
type LoanController struct {
	BaseController
}

//放款列表
func (this *LoanController) LoanList() {
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
		condition += " AND l.account = ?"
		pars = append(pars, account)
	}
	userName := this.GetString("user_name") //姓名
	if userName != "" {
		condition += " AND um.verifyrealname = ?"
		pars = append(pars, userName)
	}
	idCard := this.GetString("operator") //操作人
	if idCard != "" {
		if idCard == "1" {
			condition += " AND l.audit_type = ?"
			pars = append(pars, 1)
		} else {
			condition += " AND su.displayname = ?"
			pars = append(pars, idCard)
		}
	}
	//放款时间
	startTime := this.GetString("startTime")
	if startTime != "" {
		condition += ` AND tr.dt_order>=?`
		pars = append(pars, startTime)
	}

	endTime := this.GetString("endTime")
	if endTime != "" {
		condition += ` AND tr.dt_order<=?`
		pars = append(pars, endTime)
	}
	if condition == "" { //默认提交时间1天内
		condition += ` AND tr.dt_order >= ? AND tr.dt_order <= ?`
		now := time.Now().Format("2006-01-02")
		startTime = now + " 00:00:00"
		endTime = now + " 23:59:59"
		pars = append(pars, startTime)
		pars = append(pars, endTime)
	}
	//状态
	status := this.GetString("state")
	if status == "" {
		status = "CONFIRM"
	}
	condition += ` AND tr.state=?`
	pars = append(pars, status)

	loanList, err := models.LoanList(utils.StartIndex(pageNum, pageSize), pageSize, true, condition, pars...)

	if err != nil {
		cache.RecordLogs(this.User.Id, 0, this.User.Name, this.User.DisplayName, "查询订单列表", "查询订单列表", err.Error(), this.Ctx.Input)
		this.Abort("查询订单列表")
		return
	}
	for _, v := range loanList {
		if v.DtOrder.IsZero() {
			v.DtOrderStr = ""
		} else {
			v.DtOrderStr = v.DtOrder.Format("2006-01-02 15:04:05")
		}
	}
	count, err := models.QueryListCount(condition, pars...)
	if err != nil {
		cache.RecordLogs(this.User.Id, 0, this.User.Name, this.User.DisplayName, "获取订单总数异常", "获取订单总数异常", err.Error(), this.Ctx.Input)
		this.Abort("获取订单总数异常")
		return
	}
	names, _ := models.GetCreditUsersList()
	pageCount := utils.PageCount(count, pageSize)
	this.Data["startTime"] = startTime
	this.Data["endTime"] = endTime
	this.Data["list"] = loanList
	this.Data["currPage"] = pageNum
	this.Data["count"] = count
	this.Data["pageSize"] = pageSize
	this.Data["creditUsers"] = names
	this.Data["pageCount"] = pageCount
	this.TplName = "loan/loan_list.html"
}

//放款列表导出excel
func (this *LoanController) LoanListToExcel() {
	this.IsNeedTemplate()
	pars := []interface{}{}
	condition := ""
	account := this.GetString("account") //手机号
	if account != "" {
		condition += " AND l.account = ?"
		pars = append(pars, account)
	}
	userName := this.GetString("user_name") //姓名
	if userName != "" {
		condition += " AND um.verifyrealname = ?"
		pars = append(pars, userName)
	}
	idCard := this.GetString("operator") //操作人
	if idCard != "" {
		if idCard == "1" {
			condition += " AND l.audit_type = ?"
			pars = append(pars, 1)
		} else {
			condition += " AND su.displayname = ?"
			pars = append(pars, idCard)
		}
	}
	//放款时间
	if startTime := this.GetString("startTime"); startTime != "" {
		condition += ` and tr.dt_order>=?`
		pars = append(pars, startTime)
	}

	if endTime := this.GetString("endTime"); endTime != "" {
		condition += ` and tr.dt_order<=?`
		pars = append(pars, endTime)
	}
	//状态
	status := this.GetString("state")
	if status == "" {
		status = "CONFIRM"
	}
	condition += ` and tr.state=?`
	pars = append(pars, status)

	loanList, err := models.LoanList(0, 0, false, condition, pars...)
	//fmt.Println(loanList)
	if err != nil {
		cache.RecordLogs(this.User.Id, 0, this.User.Name, this.User.DisplayName, "查询订单列表", "查询订单列表", err.Error(), this.Ctx.Input)
		this.Abort("查询订单列表")
		return
	}
	for _, v := range loanList {
		if v.DtOrder.IsZero() {
			v.DtOrderStr = ""
		} else {
			v.DtOrderStr = v.DtOrder.Format("2006-01-02 15:04:05")
		}
	}
	exl := [][]string{{"手机号", "姓名", "放款时间", "订单号", "借款金额", "借款期限", "操作人"}}
	colWidth := []float64{20.0, 20.0, 27.0, 42.0, 20.0, 20.0, 20.0}
	for _, v := range loanList {
		exp := []string{
			v.Account,
			v.Verifyrealname,
			v.DtOrderStr,
			v.OrderNumber,
			utils.Float64ToString(v.Money),
			strconv.Itoa(v.LoanTermCount),
			v.Displayname,
		}

		exl = append(exl, exp)
	}
	filename, err := utils.ExportToExcel(exl, colWidth, "放款管理")
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
	http.ServeFile(this.Ctx.ResponseWriter, this.Ctx.Request, filename)
	err = os.Remove(filename)
	if err != nil {
		cache.RecordLogs(this.User.Id, 0, this.User.Name, this.User.DisplayName, "", "删除文件错误", err.Error(), this.Ctx.Input)

	}
}

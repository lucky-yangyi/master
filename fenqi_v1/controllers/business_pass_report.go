package controllers

import (
	"fenqi_v1/cache"
	"fenqi_v1/models"
	"fenqi_v1/utils"
	"net/http"
	"os"
	"strconv"
	"strings"
)

type BusinessPassReportController struct {
	BaseController
}

//通过汇总
func (this *BusinessPassReportController) PassTotal() {
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
	dealTime := this.GetString("deal_time") //处理时间
	//fmt.Println("dealTime")
	var startDealTime, endDealTime string
	if dealTime != "" {
		dealTimes := strings.Split(dealTime, "~")
		startDealTime = dealTimes[0] + " 00:00:00"
		endDealTime = dealTimes[1] + " 23:59:59"
	   // fmt.Println(startDealTime,endDealTime)
		condition += ` AND createtime >= ? AND createtime <= ?`
		pars = append(pars, startDealTime)
		pars = append(pars, endDealTime)
	}

	businessType := this.GetString("business_type") //业务类型
	if businessType == "" {
		businessType = "现金分期"
	}
	condition += ` AND business_type = ?`
	pars = append(pars, businessType)

	channel := this.GetString("channel")
	if channel != "" {
		condition += ` AND channel = ?`
		pars = append(pars, channel)
	}

	PassTotalData, err := models.GetTotalData(condition, pars...)
	if err != nil && err.Error() != utils.ErrNoRow() {
		cache.RecordLogs(this.User.Id, 0, this.User.Name, this.User.DisplayName, "获取通过汇总累计数据", "业务通过报表/通过汇总list", err.Error(), this.Ctx.Input)
		this.Abort("获取通过汇总累计数据")
		return
	}

	PassTotalList, err := models.GetPassTotalList(utils.StartIndex(pageNum, pageSize), pageSize, true, condition, pars...)
	if err != nil {
		cache.RecordLogs(this.User.Id, 0, this.User.Name, this.User.DisplayName, "获取通过汇总列表失败", "业务通过报表/通过汇总list", err.Error(), this.Ctx.Input)
		this.Abort("获取通过汇总列表失败")
		return
	}
	count, err := models.PassTotalListCount(condition, pars...)
	if err != nil {
		cache.RecordLogs(this.User.Id, 0, this.User.Name, this.User.DisplayName, "获取通过汇总列表总数失败", "业务通过报表/通过汇总count", err.Error(), this.Ctx.Input)
		this.Abort("获取通过汇总列表总数失败")
		return
	}
	pageCount := utils.PageCount(count, pageSize)
	this.Data["totalNumRegister"] = PassTotalData.TotalNumRegister //累计注册用户
	this.Data["totalNumCert"] = PassTotalData.TotalNumCert         //累计认证用户
	this.Data["totalNumCredit"] = PassTotalData.TotalMunCredit     //累计授信用户
	this.Data["totalNumLoan"] = PassTotalData.TotalNumLoan         //累计借款用户
	this.Data["totalOrderLoan"] = PassTotalData.TotalOrderLoan     //累计借款成功订单
	this.Data["totalMoneyLoan"] = PassTotalData.TotalMoneyLoan     //累计借款成功金额

	this.Data["list"] = PassTotalList
	this.Data["currPage"] = pageNum
	this.Data["count"] = count
	this.Data["pageSize"] = pageSize
	this.Data["pageCount"] = pageCount
	this.TplName = "business/business_pass_total.html"
}

//通过汇总导出EXCEL
//@router /passtotaltoexcel [get]
func (this *BusinessPassReportController) PassTotalToExcel() {
	pars := []interface{}{}
	condition := ""
	dealTime := this.GetString("deal_time") //处理时间
	var startDealTime, endDealTime string
	if dealTime != "" {
		dealTimes := strings.Split(dealTime, "~")
		startDealTime = dealTimes[0] + " 00:00:00"
		endDealTime = dealTimes[1] + " 23:59:59"
		condition += ` AND createtime >= ? AND createtime <= ?`
		pars = append(pars, startDealTime)
		pars = append(pars, endDealTime)
	}

	businessType := this.GetString("business_type") //业务类型
	if businessType == "" {
		businessType = "现金分期"
	}
	condition += ` AND business_type = ?`
	pars = append(pars, businessType)

	channel := this.GetString("channel")
	if channel != "" {
		condition += ` AND channel = ?`
		pars = append(pars, channel)
	}
	PassTotalList, err := models.GetPassTotalList(0, 0, false, condition, pars...)
	if err != nil {
		cache.RecordLogs(this.User.Id, 0, this.User.Name, this.User.DisplayName, "获取通过汇总列表失败", "业务通过报表/通过汇总list", err.Error(), this.Ctx.Input)
		this.Abort("获取通过汇总列表失败")
		return
	}
	exl := [][]string{{"日期", "注册用户", "认证通过用户", "授信通过用户", "借款成功用户", "借款成功订单", "借款成功金额"}}
	colWidth := []float64{15.0, 20.0, 20.0, 20.0, 20.0, 20.0, 20.0}
	for _, v := range PassTotalList {
		exp := []string{
			v.Createtime.Format("2006-01-02"),
			strconv.Itoa(v.NumRegister),
			strconv.Itoa(v.MunCertSucess),
			strconv.Itoa(v.MunCreditSucess),
			strconv.Itoa(v.NumLoanSucess),
			strconv.Itoa(v.OrderLoanSucess),
			utils.Float64ToString(v.MoneyLoanSucess),
		}
		exl = append(exl, exp)
	}
	tablename := "通过汇总"
	if channel != "" {
		tablename = "通过汇总(" + channel + ")"
	}
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
	http.ServeFile(this.Ctx.ResponseWriter, this.Ctx.Request, filename)
	err = os.Remove(filename)
	if err != nil {
		cache.RecordLogs(this.User.Id, 0, this.User.Name, this.User.DisplayName, "", "删除文件错误", err.Error(), this.Ctx.Input)
	}
}

//通过汇总折线图
//router /passtotallinechart [post]
func (this *BusinessPassReportController) PassTotalLineChart() {
	resultMap := make(map[string]interface{})
	resultMap["ret"] = 403
	defer func() {
		this.Data["json"] = resultMap
		this.ServeJSON()
	}()
	pars := []interface{}{}
	condition := ""
	dealTime := this.GetString("deal_time") //处理时间
	var startDealTime, endDealTime string
	if dealTime != "" {
		dealTimes := strings.Split(dealTime, "~")
		startDealTime = dealTimes[0] + " 00:00:00"
		endDealTime = dealTimes[1] + " 23:59:59"
		condition += ` AND createtime >= ? AND createtime <= ?`
		pars = append(pars, startDealTime)
		pars = append(pars, endDealTime)
	}

	businessType := this.GetString("business_type") //业务类型
	if businessType == "" {
		businessType = "现金分期"
	}
	condition += ` AND business_type = ?`
	pars = append(pars, businessType)

	chooseType, _ := this.GetInt("choose_type") //选择折线图数据类型

	channel := this.GetString("channel")
	if channel != "" {
		condition += ` AND channel = ?`
		pars = append(pars, channel)
	}

	passTotalLineData, err := models.GetPassTotalLineData(chooseType, condition, pars...)
	if err != nil {
		cache.RecordLogs(this.User.Id, 0, this.User.Name, this.User.DisplayName, "获取通过汇总折线图数据失败", "业务通过报表/通过汇总list", err.Error(), this.Ctx.Input)
		this.Abort("获取通过汇总折线图数据失败")
		return
	}
	resultMap["data"] = passTotalLineData
	resultMap["ret"] = 200
}

//流程转化率
func (this *BusinessPassReportController) ProcessConversionRate() {
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
	dealTime := this.GetString("deal_time") //处理时间
	var startDealTime, endDealTime string
	if dealTime != "" {
		dealTimes := strings.Split(dealTime, "~")
		startDealTime = dealTimes[0] + " 00:00:00"
		endDealTime = dealTimes[1] + " 23:59:59"
		condition += ` AND createtime >= ? AND createtime <= ?`
		pars = append(pars, startDealTime)
		pars = append(pars, endDealTime)
	}
	businessType := this.GetString("business_type") //业务类型
	if businessType == "" {
		businessType = "现金分期"
	}
	condition += ` AND business_type = ?`
	pars = append(pars, businessType)

	dataType := this.GetString("data_type") //数据类型
	if dataType == "" {
		dataType = "用户人数"
	}

	changeCycle, _ := this.GetInt("change_cycle") //转化周期
	if changeCycle == 0 {
		changeCycle = 1
	}

	channel := this.GetString("channel")
	if channel != "" {
		condition += ` AND channel = ?`
		pars = append(pars, channel)
	}

	ProcessConversionList, err := models.GetProcessConversionList(changeCycle, utils.StartIndex(pageNum, pageSize), pageSize, true, condition, pars...)
	if err != nil {
		cache.RecordLogs(this.User.Id, 0, this.User.Name, this.User.DisplayName, "获取流程转化数据失败", "业务通过报表/流程转化率list", err.Error(), this.Ctx.Input)
		this.Abort("获取流程转化数据失败")
		return
	}
	count, err := models.GetProcessConversionCount(changeCycle, condition, pars...)
	if err != nil {
		cache.RecordLogs(this.User.Id, 0, this.User.Name, this.User.DisplayName, "获取流程转化数据总数失败", "业务通过报表/流程转化率count", err.Error(), this.Ctx.Input)
		this.Abort("获取流程转化数据总数失败")
		return
	}
	pageCount := utils.PageCount(count, pageSize)
	this.Data["dataType"] = dataType
	this.Data["list"] = ProcessConversionList
	this.Data["currPage"] = pageNum
	this.Data["count"] = count
	this.Data["pageSize"] = pageSize
	this.Data["pageCount"] = pageCount
	this.TplName = "business/business_process_percent_conversion_.html"
}

//流程转化率导出EXCEL
//@router /processconversiontoexcel [get]
func (this *BusinessPassReportController) ProcessConversionToExcel() {
	pars := []interface{}{}
	condition := ""
	dealTime := this.GetString("deal_time") //处理时间
	var startDealTime, endDealTime string
	if dealTime != "" {
		dealTimes := strings.Split(dealTime, "~")
		startDealTime = dealTimes[0] + " 00:00:00"
		endDealTime = dealTimes[1] + " 23:59:59"
		condition += ` AND createtime >= ? AND createtime <= ?`
		pars = append(pars, startDealTime)
		pars = append(pars, endDealTime)
	}

	businessType := this.GetString("business_type") //业务类型
	if businessType == "" {
		businessType = "现金分期"
	}
	condition += ` AND business_type = ?`
	pars = append(pars, businessType)

	dataType := this.GetString("data_type") //数据类型
	if dataType == "" {
		dataType = "用户人数"
	}

	changeCycle, _ := this.GetInt("change_cycle") //转化周期
	if changeCycle == 0 {
		changeCycle = 1
	}

	channel := this.GetString("channel")
	if channel != "" {
		condition += ` AND channel = ?`
		pars = append(pars, channel)
	}

	ProcessConversionList, err := models.GetProcessConversionList(changeCycle, 0, 0, false, condition, pars...)
	if err != nil {
		cache.RecordLogs(this.User.Id, 0, this.User.Name, this.User.DisplayName, "获取流程转化数据失败", "业务通过报表/流程转化率list", err.Error(), this.Ctx.Input)
		this.Abort("获取流程转化数据失败")
		return
	}
	exl := [][]string{{"日期", "注册", "认证申请", "认证通过", "活体通过", "个人信息补充", "芝麻信用", "常用联系人", "运营商认证", "收款银行卡", "公积金", "支付宝", "授信申请", "授信通过", "首次借款申请", "首次借款通过"}}
	colWidth := []float64{15.0, 15.0, 15.0, 15.0, 15.0, 15.0, 15.0, 15.0, 15.0, 15.0, 15.0, 15.0, 15.0, 15.0, 15.0, 15.0}
	for _, v := range ProcessConversionList {
		if dataType == "用户人数" {
			exp := []string{
				v.Createtime.Format("2006-01-02"),
				strconv.Itoa(v.NumRegister),
				strconv.Itoa(v.NumCertApply),
				strconv.Itoa(v.NumCertSucess),
				strconv.Itoa(v.NumLivingSucess),
				strconv.Itoa(v.NumUsersBaseInfo),
				strconv.Itoa(v.NumZmAuth),
				strconv.Itoa(v.NumLinkMan),
				strconv.Itoa(v.NumMobileOperatorsMx),
				strconv.Itoa(v.NumBindCard),
				strconv.Itoa(v.NumGjj),
				strconv.Itoa(v.NumAliPay),
				strconv.Itoa(v.NumCreditApply),
				strconv.Itoa(v.NumCreditSucess),
				strconv.Itoa(v.NumFirstLoanApply),
				strconv.Itoa(v.NumFirstLoanSucess),
			}
			exl = append(exl, exp)
		} else {
			exp := []string{
				v.Createtime.Format("2006-01-02"),
				strconv.Itoa(v.NumRegister),
				utils.GetDivide(v.NumCertApply, v.NumRegister),
				utils.GetDivide(v.NumCertSucess, v.NumRegister),
				utils.GetDivide(v.NumLivingSucess, v.NumRegister),
				utils.GetDivide(v.NumUsersBaseInfo, v.NumRegister),
				utils.GetDivide(v.NumZmAuth, v.NumRegister),
				utils.GetDivide(v.NumLinkMan, v.NumRegister),
				utils.GetDivide(v.NumMobileOperatorsMx, v.NumRegister),
				utils.GetDivide(v.NumBindCard, v.NumRegister),
				utils.GetDivide(v.NumGjj, v.NumRegister),
				utils.GetDivide(v.NumAliPay, v.NumRegister),
				utils.GetDivide(v.NumCreditApply, v.NumRegister),
				utils.GetDivide(v.NumCreditSucess, v.NumRegister),
				utils.GetDivide(v.NumFirstLoanApply, v.NumRegister),
				utils.GetDivide(v.NumFirstLoanSucess, v.NumRegister),
			}
			exl = append(exl, exp)
		}

	}
	tablename := "流程转化率(" + dataType + ")"
	if channel != "" {
		tablename = "流程转化率(" + dataType + "-" + channel + ")"
	}
	filename, err := utils.ProcessConversionToExcel(exl, colWidth, tablename)
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

//流程转化率折线图
//router /processconversionlinechart [post]
func (this *BusinessPassReportController) ProcessConversionLineChart() {
	resultMap := make(map[string]interface{})
	resultMap["ret"] = 403
	defer func() {
		this.Data["json"] = resultMap
		this.ServeJSON()
	}()
	pars := []interface{}{}
	condition := ""
	dealTime := this.GetString("deal_time") //处理时间
	var startDealTime, endDealTime string
	if dealTime != "" {
		dealTimes := strings.Split(dealTime, "~")
		startDealTime = dealTimes[0] + " 00:00:00"
		endDealTime = dealTimes[1] + " 23:59:59"
		condition += ` AND createtime >= ? AND createtime <= ?`
		pars = append(pars, startDealTime)
		pars = append(pars, endDealTime)
	}

	businessType := this.GetString("business_type") //业务类型
	if businessType == "" {
		businessType = "现金分期"
	}
	condition += ` AND business_type = ?`
	pars = append(pars, businessType)

	dataType := this.GetString("data_type") //数据类型
	if dataType == "" {
		dataType = "用户人数"
	}

	changeCycle, _ := this.GetInt("change_cycle") //转化周期
	if changeCycle == 0 {
		changeCycle = 1
	}

	chooseType, _ := this.GetInt("choose_type") //选择折线图数据类型

	channel := this.GetString("channel")
	if channel != "" {
		condition += ` AND channel = ?`
		pars = append(pars, channel)
	}

	processConversionLineData, err := models.GetProcessConversionLineData(chooseType, changeCycle, dataType, condition, pars...)
	if err != nil {
		cache.RecordLogs(this.User.Id, 0, this.User.Name, this.User.DisplayName, "获取流程转化率折线图数据失败", "业务通过报表/流程转化率list", err.Error(), this.Ctx.Input)
		this.Abort("获取流程转化率折线图数据失败")
		return
	}
	resultMap["data"] = processConversionLineData
	resultMap["ret"] = 200
}

//授信通过率(creditType 1:风控通过率 2:系统通过率 3:信审通过率)
func (this *BusinessPassReportController) CreditPass(creditType int) {
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
	dealTime := this.GetString("deal_time") //处理时间
	var startDealTime, endDealTime string
	if dealTime != "" {
		dealTimes := strings.Split(dealTime, "~")
		startDealTime = dealTimes[0] + " 00:00:00"
		endDealTime = dealTimes[1] + " 23:59:59"
		condition += ` AND createtime >= ? AND createtime <= ?`
		pars = append(pars, startDealTime)
		pars = append(pars, endDealTime)
	}

	businessType := this.GetString("business_type") //业务类型
	if businessType == "" {
		businessType = "现金分期"
	}
	condition += ` AND business_type = ?`
	pars = append(pars, businessType)

	countDimension := this.GetString("count_dimension") //统计维度
	if countDimension == "" {
		countDimension = "次数"
	}

	dataType := this.GetString("data_type") //数据类型
	if dataType == "" {
		dataType = "数量"
	}

	channel := this.GetString("channel")
	if channel != "" {
		condition += ` AND channel = ?`
		pars = append(pars, channel)
	}

	//creditType 1:风控通过率 2:系统通过率 3:信审通过率
	switch creditType {
	case 1:
		controlPassList, err := models.GetControlPassList(countDimension, utils.StartIndex(pageNum, pageSize), pageSize, true, condition, pars...)
		if err != nil {
			cache.RecordLogs(this.User.Id, 0, this.User.Name, this.User.DisplayName, "获取风控通过率数据失败", "业务通过报表/风控通过率list", err.Error(), this.Ctx.Input)
			this.Abort("获取风控通过率数据失败")
			return
		}
		this.Data["list"] = controlPassList
	case 2:
		systemPassList, err := models.GetSystemPassList(countDimension, utils.StartIndex(pageNum, pageSize), pageSize, true, condition, pars...)
		if err != nil {
			cache.RecordLogs(this.User.Id, 0, this.User.Name, this.User.DisplayName, "获取系统通过率数据失败", "业务通过报表/系统通过率list", err.Error(), this.Ctx.Input)
			this.Abort("获取系统通过率数据失败")
			return
		}
		this.Data["list"] = systemPassList
	case 3:
		creditPassList, err := models.GetCreditPassList(countDimension, utils.StartIndex(pageNum, pageSize), pageSize, true, condition, pars...)
		if err != nil {
			cache.RecordLogs(this.User.Id, 0, this.User.Name, this.User.DisplayName, "获取信审通过率数据失败", "业务通过报表/信审通过率list", err.Error(), this.Ctx.Input)
			this.Abort("获取信审通过率数据失败")
			return
		}
		this.Data["list"] = creditPassList
	}
	count, err := models.CreditPassListCount(condition, pars...)
	if err != nil {
		cache.RecordLogs(this.User.Id, 0, this.User.Name, this.User.DisplayName, "获取风控通过数据总数失败", "业务通过报表/风控通过率count", err.Error(), this.Ctx.Input)
		this.Abort("获取风控通过数据总数失败")
		return
	}
	pageCount := utils.PageCount(count, pageSize)
	this.Data["dataType"] = dataType
	this.Data["currPage"] = pageNum
	this.Data["count"] = count
	this.Data["pageSize"] = pageSize
	this.Data["pageCount"] = pageCount
}

//授信通过率折线图(creditType 1:风控通过率 2:系统通过率 3:信审通过率)
func (this *BusinessPassReportController) CreditPassLineChart(creditType int) {
	resultMap := make(map[string]interface{})
	resultMap["ret"] = 403
	defer func() {
		this.Data["json"] = resultMap
		this.ServeJSON()
	}()
	pars := []interface{}{}
	condition := ""
	dealTime := this.GetString("deal_time") //处理时间
	var startDealTime, endDealTime string
	if dealTime != "" {
		dealTimes := strings.Split(dealTime, "~")
		startDealTime = dealTimes[0] + " 00:00:00"
		endDealTime = dealTimes[1] + " 23:59:59"
		condition += ` AND createtime >= ? AND createtime <= ?`
		pars = append(pars, startDealTime)
		pars = append(pars, endDealTime)
	}

	businessType := this.GetString("business_type") //业务类型
	if businessType == "" {
		businessType = "现金分期"
	}
	condition += ` AND business_type = ?`
	pars = append(pars, businessType)

	countDimension := this.GetString("count_dimension") //统计维度
	if countDimension == "" {
		countDimension = "次数"
	}

	dataType := this.GetString("data_type") //数据类型
	if dataType == "" {
		dataType = "数量"
	}

	channel := this.GetString("channel")
	if channel != "" {
		condition += ` AND channel = ?`
		pars = append(pars, channel)
	}

	chooseType, _ := this.GetInt("choose_type") //选择折线图数据类型
	//creditType 1:风控通过率 2:系统通过率 3:信审通过率
	switch creditType {
	case 1:
		controlPassRateLineData, err := models.GetControlPassRateLineData(chooseType, countDimension, dataType, condition, pars...)
		if err != nil {
			cache.RecordLogs(this.User.Id, 0, this.User.Name, this.User.DisplayName, "获取风控通过率折线图数据失败", "业务通过报表/风控通过率list", err.Error(), this.Ctx.Input)
			this.Abort("获取风控通过率折线图数据失败")
			return
		}
		resultMap["data"] = controlPassRateLineData
	case 2:
		systemPassRateLineData, err := models.GetSystemPassRateLineData(chooseType, countDimension, dataType, condition, pars...)
		if err != nil {
			cache.RecordLogs(this.User.Id, 0, this.User.Name, this.User.DisplayName, "获取系统通过率折线图数据失败", "业务通过报表/系统通过率list", err.Error(), this.Ctx.Input)
			this.Abort("获取系统通过率折线图数据失败")
			return
		}
		resultMap["data"] = systemPassRateLineData
	case 3:
		creditPassRateLineData, err := models.GetCreditPassRateLineData(chooseType, countDimension, dataType, condition, pars...)
		if err != nil {
			cache.RecordLogs(this.User.Id, 0, this.User.Name, this.User.DisplayName, "获取信审通过率折线图数据失败", "业务通过报表/信审通过率list", err.Error(), this.Ctx.Input)
			this.Abort("获取信审通过率折线图数据失败")
			return
		}
		resultMap["data"] = creditPassRateLineData
	}

	resultMap["ret"] = 200

}

//风控通过率
func (this *BusinessPassReportController) ControlPassRate() {
	this.IsNeedTemplate()
	this.CreditPass(1)
	this.TplName = "business/business_control_pass.html"
}

//风控通过率导出EXCEL
//@router /controlpassratetoexcel [get]
func (this *BusinessPassReportController) ControlPassRateToExcel() {
	pars := []interface{}{}
	condition := ""
	dealTime := this.GetString("deal_time") //处理时间
	var startDealTime, endDealTime string
	if dealTime != "" {
		dealTimes := strings.Split(dealTime, "~")
		startDealTime = dealTimes[0] + " 00:00:00"
		endDealTime = dealTimes[1] + " 23:59:59"
		condition += ` AND createtime >= ? AND createtime <= ?`
		pars = append(pars, startDealTime)
		pars = append(pars, endDealTime)
	}

	businessType := this.GetString("business_type") //业务类型
	if businessType == "" {
		businessType = "现金分期"
	}
	condition += ` AND business_type = ?`
	pars = append(pars, businessType)

	countDimension := this.GetString("count_dimension") //统计维度
	if countDimension == "" {
		countDimension = "次数"
	}

	dataType := this.GetString("data_type") //数据类型
	if dataType == "" {
		dataType = "数量"
	}

	channel := this.GetString("channel")
	if channel != "" {
		condition += ` AND channel = ?`
		pars = append(pars, channel)
	}

	controlPassList, err := models.GetControlPassList(countDimension, 0, 0, false, condition, pars...)
	if err != nil {
		cache.RecordLogs(this.User.Id, 0, this.User.Name, this.User.DisplayName, "获取风控通过率数据失败", "业务通过报表/风控通过率list", err.Error(), this.Ctx.Input)
		this.Abort("获取风控通过率数据失败")
		return
	}
	exl := [][]string{{"日期", "申请", "通过", "驳回", "关闭30天", "永久关闭", "流程中"}}
	colWidth := []float64{15.0, 20.0, 20.0, 20.0, 20.0, 20.0, 20.0}
	for _, v := range controlPassList {
		if dataType == "数量" {
			exp := []string{
				v.Createtime.Format("2006-01-02"),
				strconv.Itoa(v.AllCreditApply),
				strconv.Itoa(v.AllCreditSuccess),
				strconv.Itoa(v.AllCreditRejected),
				strconv.Itoa(v.AllCreditCloseDays),
				strconv.Itoa(v.AllCreditClosePermanent),
				strconv.Itoa(v.AllCrediting),
			}
			exl = append(exl, exp)
		} else {
			exp := []string{
				v.Createtime.Format("2006-01-02"),
				strconv.Itoa(v.AllCreditApply),
				utils.GetDivide(v.AllCreditSuccess, v.AllCreditApply),
				utils.GetDivide(v.AllCreditRejected, v.AllCreditApply),
				utils.GetDivide(v.AllCreditCloseDays, v.AllCreditApply),
				utils.GetDivide(v.AllCreditClosePermanent, v.AllCreditApply),
				utils.GetDivide(v.AllCrediting, v.AllCreditApply),
			}
			exl = append(exl, exp)
		}

	}
	tablename := "风控通过率(" + dataType + ")"
	if channel != "" {
		tablename = "风控通过率(" + dataType + "-" + channel + ")"
	}
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
	http.ServeFile(this.Ctx.ResponseWriter, this.Ctx.Request, filename)
	err = os.Remove(filename)
	if err != nil {
		cache.RecordLogs(this.User.Id, 0, this.User.Name, this.User.DisplayName, "", "删除文件错误", err.Error(), this.Ctx.Input)
	}
}

//风控通过率折线图
//@router /controlpassratelinechart [post]
func (this *BusinessPassReportController) ControlPassRateLineChart() {
	this.CreditPassLineChart(1)
}

//系统通过率
func (this *BusinessPassReportController) SystemPassRate() {
	this.IsNeedTemplate()
	this.CreditPass(2)
	this.TplName = "business/business_system_pass.html"
}

//系统通过率导出EXCEL
//@router /systempassratetoexcel [get]
func (this *BusinessPassReportController) SystemPassRateToExcel() {
	pars := []interface{}{}
	condition := ""
	dealTime := this.GetString("deal_time") //处理时间
	var startDealTime, endDealTime string
	if dealTime != "" {
		dealTimes := strings.Split(dealTime, "~")
		startDealTime = dealTimes[0] + " 00:00:00"
		endDealTime = dealTimes[1] + " 23:59:59"
		condition += ` AND createtime >= ? AND createtime <= ?`
		pars = append(pars, startDealTime)
		pars = append(pars, endDealTime)
	}

	businessType := this.GetString("business_type") //业务类型
	if businessType == "" {
		businessType = "现金分期"
	}
	condition += ` AND business_type = ?`
	pars = append(pars, businessType)

	countDimension := this.GetString("count_dimension") //统计维度
	if countDimension == "" {
		countDimension = "次数"
	}

	dataType := this.GetString("data_type") //数据类型
	if dataType == "" {
		dataType = "数量"
	}

	channel := this.GetString("channel")
	if channel != "" {
		condition += ` AND channel = ?`
		pars = append(pars, channel)
	}

	systemPassList, err := models.GetSystemPassList(countDimension, 0, 0, false, condition, pars...)
	if err != nil {
		cache.RecordLogs(this.User.Id, 0, this.User.Name, this.User.DisplayName, "获取系统通过率数据失败", "业务通过报表/系统通过率list", err.Error(), this.Ctx.Input)
		this.Abort("获取系统通过率数据失败")
		return
	}
	exl := [][]string{{"日期", "申请", "通过", "公积金锁定30天", "支付宝锁定30天", "运营商驳回", "公积金&支付宝", "公积金&运营商", "支付宝&运营商", "公积金&支付宝&运营商", "审核", "关闭30天", "永久关闭", "流程中"}}
	colWidth := []float64{15.0, 13.0, 13.0, 23.0, 23.0, 23.0, 23.0, 23.0, 23.0, 23.0, 13.0, 13.0, 13.0, 13.0}
	for _, v := range systemPassList {
		if dataType == "数量" {
			exp := []string{
				v.Createtime.Format("2006-01-02"),
				strconv.Itoa(v.AllCreditApply),
				strconv.Itoa(v.SystemCreditSucess),
				strconv.Itoa(v.SystemCreditGjjRejected),
				strconv.Itoa(v.SystemCreditZfbRejected),
				strconv.Itoa(v.SystemCreditMobileRejected),
				strconv.Itoa(v.SystemCreditGjjzfbRejected),
				strconv.Itoa(v.SystemCreditGjjmobileRejected),
				strconv.Itoa(v.SystemCreditZfbmobileRejected),
				strconv.Itoa(v.SystemCreditGjjzfbmobileRejected),
				strconv.Itoa(v.SystemCreditIntoXS),
				strconv.Itoa(v.SystemCreditCloseDays),
				strconv.Itoa(v.SystemCreditClosePermanent),
				strconv.Itoa(v.SystemCrediting),
			}
			exl = append(exl, exp)
		} else {
			exp := []string{
				v.Createtime.Format("2006-01-02"),
				strconv.Itoa(v.AllCreditApply),
				utils.GetDivide(v.SystemCreditSucess, v.AllCreditApply),
				utils.GetDivide(v.SystemCreditGjjRejected, v.AllCreditApply),
				utils.GetDivide(v.SystemCreditZfbRejected, v.AllCreditApply),
				utils.GetDivide(v.SystemCreditMobileRejected, v.AllCreditApply),
				utils.GetDivide(v.SystemCreditGjjzfbRejected, v.AllCreditApply),
				utils.GetDivide(v.SystemCreditGjjmobileRejected, v.AllCreditApply),
				utils.GetDivide(v.SystemCreditZfbmobileRejected, v.AllCreditApply),
				utils.GetDivide(v.SystemCreditGjjzfbmobileRejected, v.AllCreditApply),
				utils.GetDivide(v.SystemCreditIntoXS, v.AllCreditApply),
				utils.GetDivide(v.SystemCreditCloseDays, v.AllCreditApply),
				utils.GetDivide(v.SystemCreditClosePermanent, v.AllCreditApply),
				utils.GetDivide(v.SystemCrediting, v.AllCreditApply),
			}
			exl = append(exl, exp)
		}

	}
	tablename := "系统通过率(" + dataType + ")"
	if channel != "" {
		tablename = "系统通过率(" + dataType + "-" + channel + ")"
	}
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
	http.ServeFile(this.Ctx.ResponseWriter, this.Ctx.Request, filename)
	err = os.Remove(filename)
	if err != nil {
		cache.RecordLogs(this.User.Id, 0, this.User.Name, this.User.DisplayName, "", "删除文件错误", err.Error(), this.Ctx.Input)
	}
}

//系统通过率折线图
//@router /systempassratelinechart [post]
func (this *BusinessPassReportController) SystemPassRateLineChart() {
	this.CreditPassLineChart(2)
}

//信审通过率
func (this *BusinessPassReportController) CreditPassRate() {
	this.IsNeedTemplate()
	this.CreditPass(3)
	this.TplName = "business/business_credit_pass.html"
}

//信审通过率导出EXCEL
//@router /creditpassratetoexcel [get]
func (this *BusinessPassReportController) CreditPassRateToExcel() {
	pars := []interface{}{}
	condition := ""
	dealTime := this.GetString("deal_time") //处理时间
	var startDealTime, endDealTime string
	if dealTime != "" {
		dealTimes := strings.Split(dealTime, "~")
		startDealTime = dealTimes[0] + " 00:00:00"
		endDealTime = dealTimes[1] + " 23:59:59"
		condition += ` AND createtime >= ? AND createtime <= ?`
		pars = append(pars, startDealTime)
		pars = append(pars, endDealTime)
	}

	businessType := this.GetString("business_type") //业务类型
	if businessType == "" {
		businessType = "现金分期"
	}
	condition += ` AND business_type = ?`
	pars = append(pars, businessType)

	countDimension := this.GetString("count_dimension") //统计维度
	if countDimension == "" {
		countDimension = "次数"
	}

	dataType := this.GetString("data_type") //数据类型
	if dataType == "" {
		dataType = "数量"
	}

	channel := this.GetString("channel")
	if channel != "" {
		condition += ` AND channel = ?`
		pars = append(pars, channel)
	}

	creditPassList, err := models.GetCreditPassList(countDimension, 0, 0, false, condition, pars...)
	if err != nil {
		cache.RecordLogs(this.User.Id, 0, this.User.Name, this.User.DisplayName, "获取信审通过率数据失败", "业务通过报表/信审通过率list", err.Error(), this.Ctx.Input)
		this.Abort("获取信审通过率数据失败")
		return
	}
	exl := [][]string{{"日期", "审核", "通过", "信审驳回", "关闭30天", "永久关闭", "流程中"}}
	colWidth := []float64{15.0, 20.0, 20.0, 20.0, 20.0, 20.0, 20.0}
	for _, v := range creditPassList {
		if dataType == "数量" {
			exp := []string{
				v.Createtime.Format("2006-01-02"),
				strconv.Itoa(v.SystemCreditIntoXS),
				strconv.Itoa(v.XSCreditSuccess),
				strconv.Itoa(v.XSCreditRejected),
				strconv.Itoa(v.XSCreditCloseDays),
				strconv.Itoa(v.XSCreditClosePermanent),
				strconv.Itoa(v.XSCrediting),
			}
			exl = append(exl, exp)
		} else {
			exp := []string{
				v.Createtime.Format("2006-01-02"),
				strconv.Itoa(v.SystemCreditIntoXS),
				utils.GetDivide(v.XSCreditSuccess, v.SystemCreditIntoXS),
				utils.GetDivide(v.XSCreditRejected, v.SystemCreditIntoXS),
				utils.GetDivide(v.XSCreditCloseDays, v.SystemCreditIntoXS),
				utils.GetDivide(v.XSCreditClosePermanent, v.SystemCreditIntoXS),
				utils.GetDivide(v.XSCrediting, v.SystemCreditIntoXS),
			}
			exl = append(exl, exp)
		}
	}
	tablename := "信审通过率(" + dataType + ")"
	if channel != "" {
		tablename = "信审通过率(" + dataType + "-" + channel + ")"
	}
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
	http.ServeFile(this.Ctx.ResponseWriter, this.Ctx.Request, filename)
	err = os.Remove(filename)
	if err != nil {
		cache.RecordLogs(this.User.Id, 0, this.User.Name, this.User.DisplayName, "", "删除文件错误", err.Error(), this.Ctx.Input)
	}
}

//信审通过率折线图
//@router /creditpassratelinechart [post]
func (this *BusinessPassReportController) CreditPassRateLineChart() {
	this.CreditPassLineChart(3)
}

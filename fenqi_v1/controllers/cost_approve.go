package controllers

import (
	// "encoding/json"

	"fenqi_v1/cache"
	"fenqi_v1/models"
	"fenqi_v1/services"
	"fenqi_v1/utils"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
	"zcm_tools/file"

	"github.com/astaxie/beego"
)

//费用减免审批
type CostApproveController struct {
	BaseController
}

//@Summary 添加审批记录
func (this *CostApproveController) AddApprove() {
	defer this.ServeJSON()
	// loginSysId := this.User.Id
	loginSysId := this.User.RoleId
	var requestHas bool
	config, err := models.QueryByIdCostConfig(1)
	if err == nil {
		// requestHas = utils.CheckIsSysId(loginSysId, config.RequestId)
		requestHas = utils.CheckIsSysId(loginSysId, config.RequestRoleId)
	}
	if requestHas {
		imgUrl := this.GetString("imgUrl")
		if imgUrl == "" {
			this.Data["json"] = map[string]interface{}{"ret": 304, "errmsg": "图片信息不能为空"}
			return
		}
		imgUrl = strings.Replace(imgUrl, "data:image/jpeg;base64,", "", -1)
		imgUrl = strings.Replace(imgUrl, "data:image/png;base64,", "", -1)
		arr := strings.Split(imgUrl, "||")
		length := len(arr)
		imagesUrl := ""
		deleteFile := []string{}
		if length > 0 {
			timeStr := strconv.FormatInt(time.Now().Unix(), 10)
			for i := 0; i < length; i++ {
				filename := "costApprove" + timeStr + strconv.Itoa(i) + ".jpg"
				filepath := "./static/image/" + filename
				err := file.SaveBase64ToFile(arr[i], filepath)
				if err != nil {
					cache.RecordLogs(this.User.Id, 0, this.User.Name, this.User.DisplayName, "保存本地图片有误", "催收管理/AddApprove", err.Error(), this.Ctx.Input)
					continue
				}
				deleteFile = append(deleteFile, filename)
				err, url := services.UploadAliyun(filename, filepath)
				if err != nil {
					cache.RecordLogs(this.User.Id, 0, this.User.Name, this.User.DisplayName, "上传本地图片到阿里云有误", "催收管理/AddApprove", err.Error(), this.Ctx.Input)
					continue
				}
				if err == nil {
					imagesUrl += url + ","
				}
			}
		}
		if imagesUrl != "" {
			imagesUrl = imagesUrl[:len(imagesUrl)-1]
		}
		phone := this.GetString("phone")
		if phone == "" {
			this.Data["json"] = map[string]interface{}{"ret": 304, "errmsg": "电话号码不能为空"}
			return
		}
		user, err := models.FindByAccountUser(phone)
		if err != nil || user.RepSchId == 0 {
			this.Data["json"] = map[string]interface{}{"ret": 304, "errmsg": "没有此用户的逾期记录，请检查手机号是否正确"}
			return
		}
		money, _ := this.GetFloat("money")
		if money < 0 {
			this.Data["json"] = map[string]interface{}{"ret": 304, "errmsg": "减免金额不能小于0"}
			return
		}
		approveSysId, _ := this.GetInt("approveSysId")
		reason := this.GetString("reason")
		uid := user.Id
		repSchId := user.RepSchId
		loanId := user.LoanId
		collectionPhase := utils.Mtype(user.Mtype)
		var company string
		state, _ := models.GetLoanState(loanId)
		if state != "BACKING" {
			this.Data["json"] = map[string]interface{}{"ret": 304, "errmsg": "订单已结清，不允许发起减免申请。"}
			return
		}
		count, _ := models.QueryCostUnfinishCount(uid, loanId)
		if count > 0 {
			this.Data["json"] = map[string]interface{}{"ret": 304, "errmsg": "该笔借款有尚未完成的减免。"}
			return
		}
		err = models.AddCostReliefApprove(uid, repSchId, loanId, this.User.Id, approveSysId, this.User.DisplayName, phone, reason, imagesUrl, collectionPhase, company, money)
		if err == nil {
			this.Data["json"] = map[string]interface{}{"ret": 200}
		} else {
			cache.RecordLogs(this.User.Id, 0, this.User.Name, this.User.DisplayName, "提交费用减免审批失败", "催收管理/AddApprove", err.Error(), this.Ctx.Input)
			this.Data["json"] = map[string]interface{}{"ret": 304, "errmsg": err.Error()}
		}
		if len(deleteFile) > 0 {
			for i := 0; i < len(deleteFile); i++ {
				os.Remove("./static/image/" + deleteFile[i])
			}
		}
	} else {
		this.Data["json"] = map[string]interface{}{"ret": 304, "errmsg": "权限不足，请联系管理员。"}
	}

}

// @Summary 查询费用减免审批记录
func (this *CostApproveController) QueryApprove() {
	this.IsNeedTemplate()
	this.Data["isDisplayAdd"] = false      //是否显示添加按钮
	this.Data["isDisplayApplove"] = true   //是否显示退回按钮
	this.Data["loginSysId"] = this.User.Id //比较是否显示退回按钮(记录审批人跟当前登录账号 id 相同才显示该条退回)
	var count int
	var requestHas, approveHas, disposeHas, stationHas, displayExeclButton bool
	config, err := models.QueryByIdCostConfig(1)
	if err == nil {
		requestHas = utils.CheckIsSysId(this.User.RoleId, config.RequestRoleId)           //申请人权限（按角色）
		approveHas = utils.CheckIsSysId(this.User.Id, config.ApproveId)                   //审批人权限(按账号 id)
		disposeHas = utils.CheckIsSysId(this.User.Id, config.DisposeId)                   //处理人权限(按账号 id)
		stationHas = utils.CheckIsSysId(this.User.StationId, config.StationId)            //岗位权限
		displayExeclButton = utils.CheckIsSysId(this.User.RoleId, config.CostExeclRoleId) //导出 execl 权限(按角色)
	} else {
		cache.RecordLogs(this.User.Id, 0, this.User.Name, this.User.DisplayName, "查询所有费用减免人员配置信息失败", "催收管理/QueryApprove", err.Error(), this.Ctx.Input)
	}
	if requestHas || approveHas || disposeHas || stationHas {
		pageNo, _ := this.GetInt("page", 1)
		pageSize, _ := this.GetInt("pageSize")
		if pageSize < 1 {
			pageSize = 30
		}
		state := this.GetString("state")
		beginTime := this.GetString("beginTime")
		beginTime = strings.Replace(beginTime, "+", " ", -1)
		endTime := this.GetString("endTime")
		endTime = strings.Replace(endTime, "+", " ", -1)
		account := this.GetString("account")
		collectionPhase := this.GetString("collectionPhase")
		//company := this.GetString("company")
		pars := []string{}
		condition := ""
		if requestHas {
			this.Data["isDisplayAdd"] = true
			condition += " and crf.request_sys_id =?"
			pars = append(pars, strconv.Itoa(this.User.Id))
			this.Data["isDisplayApplove"] = false
		} else {
			if disposeHas {
				this.Data["isDisplayApplove"] = false
			}
			if stationHas && !approveHas {
				this.Data["isDisplayApplove"] = false
			}
		}
		if state != "" {
			if state == "REFUSE" {
				condition += " and crf.state in('REFUSE','REFUSEDONE')"
			} else {
				condition += " and crf.state=?"
				pars = append(pars, state)
			}
		}
		if beginTime != "" {
			condition += " and crf.create_time>=?"
			pars = append(pars, beginTime)
		}
		if endTime != "" {
			condition += " and crf.create_time<=?"
			pars = append(pars, endTime)
		}
		if account != "" {
			condition += " and crf.phone=?"
			pars = append(pars, account)
		}
		if collectionPhase != "" {
			condition += " and crf.collection_phase=?"
			pars = append(pars, collectionPhase)
		}
		/*	if company != "" {
			condition += " and crf.company=?"
			pars = append(pars, company)
		}*/
		count = models.QueryCount(condition, pars...)
		if count > 0 {
			list, err := models.QueryCostReliefApprove(utils.StartIndex(pageNo, pageSize), pageSize, condition, pars...)
			if err != nil {
				cache.RecordLogs(this.User.Id, 0, this.User.Name, this.User.DisplayName, "查询所有费用减免失败", "催收管理/QueryApprove", err.Error(), this.Ctx.Input)
				beego.Debug(err)
			}
			for k, v := range list {
				if v.IsSubmit == 1 {
					list[k].Money = v.SubmitMoney
					continue
				}
				if v.IsApprove == 1 {
					list[k].Money = v.ApproveMoney
					continue
				}
			}
			this.Data["list"] = list
		}
		pagecount := utils.PageCount(count, pageSize)
		this.Data["currpage"] = pageNo
		this.Data["pagecount"] = pagecount
		this.Data["pageSize"] = pageSize
	}
	/*orgStr, err := cache.GetCacheDataByStation(this.User.StationId)
	if err != nil {
		cache.RecordLogs(this.User.Id, 0, this.User.Name, this.User.DisplayName, "根据岗位获取数据权限失败", "新催收管理/QueryApprove", err.Error(), this.Ctx.Input)
	}*/
	//this.Data["hasBtnPower1"] = models.CheckPowerByStaId2(this.User.StationId) //权限设置(超管，总监，经理)
	this.Data["displayExeclButton"] = displayExeclButton
	this.Data["stationId"] = this.User.StationId
	this.Data["count"] = count
	this.TplName = "collection/cost_relief_list.html"
}

// //费用减免页面数据导出Excel
func (this *CostApproveController) CostApproveDataToExcel() {
	pars := []string{}
	condition := ""
	// loginSysId := this.User.Id
	var requestHas bool
	config, err := models.QueryByIdCostConfig(1)
	if err == nil {
		// requestIdStr := config.RequestId
		// requestHas = utils.CheckIsSysId(loginSysId, requestIdStr)
		requestHas = utils.CheckIsSysId(this.User.RoleId, config.RequestRoleId)
	} else {
		cache.RecordLogs(this.User.Id, 0, this.User.Name, this.User.DisplayName, "查询所有费用减免人员配置信息失败", "费用减免/CostApproveDataToExcel", err.Error(), this.Ctx.Input)
	}
	if requestHas {
		condition += " AND crf.request_sys_id = ?"
		pars = append(pars, strconv.Itoa(this.User.Id))
	}
	beginTime := this.GetString("beginTime") //开始日期
	if beginTime != "" {
		beginTime = strings.Replace(beginTime, "+", " ", -1)
		condition += " AND crf.create_time >= ?"
		pars = append(pars, beginTime)
	}
	endTime := this.GetString("endTime") //结束日期
	if endTime != "" {
		endTime = strings.Replace(endTime, "+", " ", -1)
		condition += " AND crf.create_time <= ?"
		pars = append(pars, endTime)
	}
	state := this.GetString("state") //审批状态
	if state != "" {
		if state == "REFUSE" {
			condition += " AND crf.state IN('REFUSE','REFUSEDONE')"
		} else {
			condition += " AND crf.state = ?"
			pars = append(pars, state)
		}
	}
	account := this.GetString("account") //手机号
	if account != "" {
		condition += " AND crf.phone = ?"
		pars = append(pars, account)
	}
	list, err := models.QueryCostReliefApproveNoPage(condition, pars...)
	if err != nil {
		cache.RecordLogs(this.User.Id, 0, this.User.Name, this.User.DisplayName, "查询所有费用减免失败", "费用减免/CostApproveDataToExcel", err.Error(), this.Ctx.Input)
	}
	exl := [][]string{{"发起人", "发起时间", "减免账号", "减免金额", "状态", "审批人"}}
	colWidth := []float64{15.0, 28.0, 18.0, 20.0, 12.0, 15.0}
	for k, v := range list {
		if list[k].IsSubmit == 1 {
			v.Money = v.SubmitMoney
		} else if list[k].IsApprove == 1 {
			v.Money = v.ApproveMoney
		}
		var stateString = ""
		if v.State == "NONE" {
			stateString = "审核中"
		} else if v.State == "AGREE" {
			stateString = "待处理"
		} else if v.State == "REFUSE" {
			stateString = "已退回"
		} else if v.State == "REFUSEDONE" {
			stateString = "拒绝退回"
		} else if v.State == "DONE" {
			stateString = "处理完成"
		} else if v.State == "FAIL" {
			stateString = "处理失败"
		}
		if v.State != "NONE" {
			list[k].Displayname = ""
		}
		exp := []string{
			v.RequestSysName,
			v.CreateTime,
			v.Phone,
			utils.Float64ToString(v.Money),
			stateString,
			list[k].Displayname,
		}
		exl = append(exl, exp)
	}
	filename, err := utils.ExportToExcel(exl, colWidth, "费用减免")
	if err != nil {
		cache.RecordLogs(this.User.Id, 0, this.User.Name, this.User.DisplayName, "", "保存文件错误", err.Error(), this.Ctx.Input)
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

//审批费用减免
func (this *CostApproveController) CostApprove() {
	defer this.ServeJSON()
	id, _ := this.GetInt("id")
	if id < 1 {
		this.Data["json"] = map[string]interface{}{"ret": 304, "msg": "参数不能为空"}
		return
	}
	state := this.GetString("state")
	result := this.GetString("result")
	approve, err := models.FindByApprove(id)
	if err != nil {
		cache.RecordLogs(this.User.Id, 0, this.User.Name, this.User.DisplayName, "查询信息有误", "催收管理/CostApprove", err.Error(), this.Ctx.Input)
		this.Data["json"] = map[string]interface{}{"ret": 304, "msg": "查询信息有误"}
		return
	}
	oldState := approve.State
	var msg string
	bools := true
	if state == "DONE" || oldState == "REFUSEDONE" {
		if oldState == "DONE" || oldState == "REFUSEDONE" {
			bools = false
			msg = "已完成处理,请不要重复操作。"
		} else if oldState != "AGREE" {
			bools = false
			msg = "同意操作后,在操作。"
		}
	}
	if state == "AGREE" || state == "REFUSE" {
		if oldState != "NONE" {
			bools = false
			msg = "已经审批,不要重复操作"
		}
	}
	if !bools {
		this.Data["json"] = map[string]interface{}{"ret": 304, "msg": msg}
		return
	}
	//判断费用减免金额和代还款金额
	if state == "DONE" {
		idAtr := strconv.Itoa(approve.LoanId)
		money, _ := models.GetCollectionMoneyByloanid2(idAtr)
		//money = 5.0
		if money < approve.ApproveMoney {
			this.Data["json"] = map[string]interface{}{"ret": 304, "msg": "减免金额大于待还金额，请核实"}
			return
		}
	}
	isBack, _ := this.GetBool("isBack")     //不是退回操作
	isFinish, _ := this.GetBool("isFinish") //不是处理完成操作
	if !isBack && !isFinish {
		approveMoney, _ := this.GetFloat("approveMoney")
		err = models.UpdateCostReliefApprove(id, approveMoney, "Approve")
		if err != nil {
			cache.RecordLogs(this.User.Id, 0, this.User.Name, this.User.DisplayName, "更新费用减免表失败", "催收管理/CostApprove", err.Error(), this.Ctx.Input)
			this.Data["json"] = map[string]interface{}{"ret": 304, "msg": "更新费用减免表失败" + err.Error()}
			return
		}
	}
	if isFinish { //处理完成操作
		params := map[string]interface{}{"Uid": approve.Uid, "MobileType": "XjdSystem", "RepaymentScheduleId": approve.RepSchId, "Channel": 9, "OperatorId": this.User.Id, "ReturnMoney": approve.ApproveMoney, "Remark": "催收协商减免", "OidPaybill": ""}
		_, err := services.PostApi(utils.Loan_Repayment, params)
		if err != nil {
			models.ApproveRelie(id, this.User.Id, "FAIL", err.Error())
			cache.RecordLogs(this.User.Id, 0, this.User.Name, this.User.DisplayName, "还款失败，请求接口失败", "催收管理/CostApprove", err.Error(), this.Ctx.Input)
			this.Data["json"] = map[string]interface{}{"ret": 304, "msg": "还款失败，请求接口失败" + err.Error()}
			return
		}
	}
	has, msg := models.ApproveRelie(id, this.User.Id, state, result)
	if has {
		this.Data["json"] = map[string]interface{}{"ret": 200}
	} else {
		this.Data["json"] = map[string]interface{}{"ret": 304, "msg": msg}
	}

}

//审批费用减免
func (this *CostApproveController) GetApprove() {
	defer this.ServeJSON()
	id, _ := this.GetInt("id")
	if id > 0 {
		approve, err := models.FindByApprove(id)
		if err == nil {
			ApproveSysName := models.FindBySysName(approve.ApproveSysId).Displayname
			DisposeSysName := models.FindBySysName(approve.DisposeSysId).Displayname
			this.Data["ApproveSysName"] = ApproveSysName
			this.Data["DisposeSysName"] = DisposeSysName
			this.Data["json"] = approve
		}
	}
}

//审批费用减免
func (this *CostApproveController) GetAdd() {
	this.IsNeedTemplate()
	this.TplName = "collection/cost_add.html"
}

// 查看
func (c *CostApproveController) Detail() {
	c.IsNeedTemplate()
	uid, _ := c.GetInt("uid")
	c.Data["uid"] = uid
	id, _ := c.GetInt("id")
	isDisplayApplove, _ := c.GetBool("isDisplayApplove")
	c.Data["isDisplayApplove"] = isDisplayApplove
	c.Data["isCheckDispose"] = false
	sysUid := c.User.Id
	var approveHas, disposeHas bool
	config, err := models.QueryByIdCostConfig(1)
	if err == nil {
		approveHas = utils.CheckIsSysId(c.User.Id, config.ApproveId) //审批权限
		disposeHas = utils.CheckIsSysId(c.User.Id, config.DisposeId) //处理权限
	} else if err.Error() != "<QuerySeter> no row found" {
		cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "查询所有费用减免人员配置信息失败", "催收管理/Detail", err.Error(), c.Ctx.Input)
	}
	sId := c.User.StationId
	if sId == 5 || sId == 6 {
		approveHas = true
		disposeHas = true
	}
	if approveHas || disposeHas {
		c.Data["isDisplayApplove"] = true
	}
	c.Data["id"] = id
	approve, err := models.FindByApprove(id)
	if err == nil {
		if approve.IsSubmit == 1 {
			approve.Money = approve.SubmitMoney
		} else if approve.IsApprove == 1 {
			approve.Money = approve.ApproveMoney
		}
		ApproveSysName := models.FindBySysName(approve.ApproveSysId).Displayname
		clist := models.FindByCraIdApproveRecord(id)
		if len(clist) > 0 {
			c.Data["clist"] = clist
		}
		if disposeSysId := approve.DisposeSysId; disposeSysId > 0 {
			DisposeSysName := models.FindBySysName(disposeSysId).Displayname
			c.Data["DisposeSysName"] = DisposeSysName
		}
		c.Data["ApproveSysName"] = ApproveSysName
		c.Data["approve"] = approve
		oldApprove := approve.ApproveSysId
		if approveHas {
			if oldApprove != sysUid {
				c.Data["isDisplayApplove"] = false
			}
		} else {
			c.Data["isCheckDispose"] = true
		}
	} else if err.Error() != "<QuerySeter> no row found" {
		cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "查询单个审批记录", "催收管理/Detail", err.Error(), c.Ctx.Input)
	}
	loanId, _ := c.GetInt("loanId")
	var productId int
	if loanId > 0 {
		//loan := models.FindByIdLoan(loanId)
		list, _ := models.FindByIdLoanFenQi(loanId)
		for i := 0; i < len(list); i++ {
			list[i].SumRepayMoney = (list[i].ShouldAlsoAmount*100 + list[i].OverMoney*100) / 100
			list[i].WaitMoney = (list[i].ShouldAlsoAmount*100 - list[i].RepayMoney*100) / 100
		}
		c.Data["list"] = list
	}
	c.Data["productId"] = productId
	//count := 0
	um, err := models.GetUsersMetadata(uid)
	/*reditAss, _ := models.GetReditAssessmentByUId(uid) //得到reditAssessment
	if reditAss == nil {
		reditState, err := models.GetLoanStateByUid(uid)
		if err != nil && err.Error() != "<QuerySeter> no row found" {
			cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "根据uid查询loanState列表失败", "催收管理/Detail", err.Error(), c.Ctx.Input)
		}
		if reditState != nil {
			c.Data["creditResult"] = reditState.Antifruad
			c.Data["creditScore"] = reditState.Score
		} else {
			c.Data["creditResult"] = -1
			c.Data["creditScore"] = ""
		}
	} else {
		c.Data["creditResult"] = reditAss.StrongAntifraud + reditAss.WeakAntifraud
		c.Data["creditScore"] = reditAss.Score
	}*/
	if err != nil && err.Error() != "<QuerySeter> no row found" {
		cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "根据uid查询个人信息有误", "催收管理/Detail", err.Error(), c.Ctx.Input)
	} else {
		//count = models.GetUsersDevicecodeCount(uid)
		um.IdCard = utils.IdCradDispose(um.IdCard)
		c.Data["um"] = um
	}
	//用户注册信息
	info, err := models.QueryUserByUid(uid)
	if err != nil && err.Error() != "<QuerySeter> no row found" {
		cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "根据uid获取用户注册信息失败", "催收管理/QueryUserByUid", err.Error(), c.Ctx.Input)
	}
	//loanCount := models.GetLoanFinishCount(uid)
	token := c.GetString("token")
	c.Data["token"] = token
	c.Data["operation_time"] = info.OperationTime
	//c.Data["count"] = count
	//c.Data["loanCount"] = loanCount
	c.TplName = "collection/cost_state.html"
}

//提交上级
func (this *CostApproveController) AddCostApprove() {
	defer this.ServeJSON()
	sysId := this.User.Id
	carId, _ := this.GetInt("id")
	parentId, _ := this.GetInt("parentId")
	opinion := this.GetString("opinion")
	if carId < 1 {
		this.Data["json"] = map[string]interface{}{"ret": 304, "msg": "参数不能为空。"}
		return
	}
	if sysId == parentId {
		this.Data["json"] = map[string]interface{}{"ret": 304, "msg": "不能提交给自己。"}
		return
	}
	// count := models.FindByCraIdRecordCount(carId)
	// if count > 0 {
	// 	this.Data["json"] = map[string]interface{}{"ret": 304, "msg": "该审核只能提交一次上级,已提交过。"}
	// 	return
	// }
	submitMoney, _ := this.GetFloat("submitMoney")
	err := models.UpdateCostReliefApprove(carId, submitMoney, "Submit")
	if err != nil {
		cache.RecordLogs(this.User.Id, 0, this.User.Name, this.User.DisplayName, "更新费用减免表失败", "催收管理/AddCostApprove", err.Error(), this.Ctx.Input)
		this.Data["json"] = map[string]interface{}{"ret": 304, "msg": "更新费用减免表失败" + err.Error()}
		return
	}
	err = models.AddCostReliefApproveRecord(1, carId, parentId, opinion, "APPROVE", "NONE")
	if err != nil {
		cache.RecordLogs(this.User.Id, 0, this.User.Name, this.User.DisplayName, "提交失败,请稍后再试", "催收管理/AddCostApprove", err.Error(), this.Ctx.Input)
		this.Data["json"] = map[string]interface{}{"ret": 304, "msg": "提交失败,请稍后再试"}
	} else {
		cache.RecordLogs(this.User.Id, 0, this.User.Name, this.User.DisplayName, "新增费用减免审批记录成功", "催收管理/AddCostApprove", "", this.Ctx.Input)
		this.Data["json"] = map[string]interface{}{"ret": 200}
	}
}

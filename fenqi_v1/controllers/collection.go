package controllers

import (
	"fenqi_v1/cache"
	"fenqi_v1/models"
	"fenqi_v1/services"
	"fenqi_v1/utils"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

//催收
type CollectionController struct {
	BaseController
}

//M0
func (c *CollectionController) M0List() {
	c.CollectList(11)
}

//S1
func (c *CollectionController) S1List() {
	c.CollectList(12)
}

//S2
func (c *CollectionController) S2List() {
	c.CollectList(13)
}

//S3
func (c *CollectionController) S3List() {
	c.CollectList(14)
}

//M2
func (c *CollectionController) M2List() {
	c.CollectList(15)
}

//M3
func (c *CollectionController) M3List() {
	c.CollectList(16)
}

//分期催收管理列表
func (c *CollectionController) CollectList(mtype int) {
	c.IsNeedTemplate()
	//当前页
	pageNum, _ := c.GetInt("page", 1) //分页信息（第几页）
	if pageNum < 1 {
		pageNum = 1
	}
	pageSize, _ := c.GetInt("pageSize")
	if pageSize < 1 {
		pageSize = 30
	}
	where, pars := Where(1, c)
	orgStr, err := cache.GetCacheDataByStation(c.User.StationId)
	if err != nil {
		cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "根据岗位获取数据权限失败", "催收管理/CollectList", err.Error(), c.Ctx.Input)
	}
	list, err := models.GetList(where, orgStr, utils.StartIndex(pageNum, pageSize), pageSize, mtype, pars...)
	if err != nil {
		cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "查询催收管理列表失败", "催收管理/CollectList", err.Error(), c.Ctx.Input)
	}
	for k, v := range list {
		//到期待还或逾期待还
		list[k].NeedPayment = models.GetNeedPayment(v.LoanId)
		//最大逾期天数
		// list[k].MaxOverdueDays = models.GetMaxOverdueDays(v.LoanId)
		//判断该用户目前是否需要催收
		// count, err := models.GetCaseStateCount(v.LoanId)
		// if err != nil {
		// 	cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "判断该用户目前是否需要催收失败", "催收管理/CollectList", err.Error(), c.Ctx.Input)
		// }
		// if count > 0 {
		// 	list[k].Flag = true
		// }
	}
	countArr := models.GetCount(where, orgStr, mtype, pars...)
	count := len(countArr)
	pageCount := utils.PageCount(count, pageSize)
	collectionUsers := models.GetCollectUser(orgStr, mtype)
	var distributeHas bool //查看权限
	config, err := models.QueryByIdCostConfig(1)
	if err == nil {
		distributeHas = utils.CheckIsSysId(c.User.RoleId, config.CollectionDistributeRoleId)
	}
	c.Data["sumMoney"] = models.GetNewSumMoney(where, orgStr, mtype, pars)
	c.Data["distributeHas"] = distributeHas
	c.Data["stationId"] = c.User.StationId
	c.Data["currPage"] = pageNum
	c.Data["pageCount"] = pageCount
	c.Data["pageSize"] = pageSize
	c.Data["list"] = list
	c.Data["collecionUsers"] = collectionUsers
	c.Data["count"] = count
	if mtype == 11 {
		c.TplName = "collection/m0_list.html"
	} else if mtype == 12 {
		c.TplName = "collection/s1_list.html"
	} else if mtype == 13 {
		c.TplName = "collection/s2_list.html"
	} else if mtype == 14 {
		c.TplName = "collection/s3_list.html"
	} else if mtype == 15 {
		c.TplName = "collection/m2_list.html"
	} else if mtype == 16 {
		c.TplName = "collection/m3_list.html"
	}
}

//M催收处理页面
func (c *CollectionController) Mhandle() {
	c.IsNeedTemplate()
	//当前页
	loanId, _ := c.GetInt("id")
	//判断是否结清
	loan, err := models.GetLoanById(loanId)
	if err != nil {
		cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "根据load查询借款记录失败", "催收管理/Mhandle", err.Error(), c.Ctx.Input)
		c.Abort("校验借款数据错误")
		return
	}
	if loan.State == "FINISH" {
		c.Abort("该案件已经结清")
		return
	}
	//判断该期是否结清
	count, err := models.GetCaseStateCount(loanId)
	if err != nil {
		cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "根据load查询还款计划失败", "催收管理/Handle", err.Error(), c.Ctx.Input)
		c.Abort("查询还款计划失败")
		return
	}
	if count == 0 {
		//结清
		c.Abort("该期已经结清")
		return
	}
	mtype, _ := c.GetInt("mtype")
	isTag := -1
	m1HandleList, err := models.GetHandleListById(loanId)
	if err != nil && err.Error() != "<QuerySeter> no row found" {
		cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "根据借款ID获取催收信息失败", "催收管理/Mhandle", err.Error(), c.Ctx.Input)
	}
	if len(m1HandleList) <= 0 {
		cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "根据借款ID获取催收信息失败,loan_id="+strconv.Itoa(loanId), "催收管理/Mhandle", err.Error(), c.Ctx.Input)
		c.Abort("获取借款信息失败")
		return
	}
	uid := 0
	for k, _ := range m1HandleList {
		if uid <= 0 {
			uid = m1HandleList[k].Uid
		}
	}

	um, err := models.GetUsersMetadata(uid)
	if err != nil && err.Error() == "<QuerySeter> no row found" {
		cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "根据uid查询个人信息失败", "催收管理/Mhandle", err.Error(), c.Ctx.Input)
		c.Abort("用户不存在")
		return
	}
	if err == nil {
		isTag = um.IsTag
		c.Data["um"] = um
		bankCard, err := models.BankcardInfo(uid)
		if err != nil && err.Error() != "<QuerySeter> no row found" {
			cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "根据UID查询银行卡信息失败", "催收管理/Mhandle", err.Error(), c.Ctx.Input)
		}
		if bankCard != nil {
			c.Data["bankCard"] = bankCard
		}
	}
	//从black_sign 查是记录
	blackSignState, err := models.QueryBlackSign(uid)
	if err != nil && err.Error() != "<QuerySeter> no row found" {
		cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "入黑查询错误", "催收管理/Mhandle", err.Error(), c.Ctx.Input)
	}
	var has, collectionExeclHas bool //导出查看权限
	var bs = true                    // 委外管理/处理/除了A3委外专员外，开通全部权限
	config, err := models.QueryByIdCostConfig(1)
	if err == nil {
		has = utils.CheckIsSysId(c.User.RoleId, config.CollectionExeclLetterRoleId)
		bs = utils.CheckIsSysId(c.User.RoleId, config.OutsourceBlackSignRoleId)
		collectionExeclHas = utils.CheckIsSysId(c.User.RoleId, config.CollectionExeclPhoneRoleId)
	}
	//借款逾期总数
	count_overdue, err := models.QueryLoanRecordOverDueCount(uid)
	if err != nil {
		cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "查询用户借款记录总数出错", "认证用户详情/借款记录GetUsersManageLoanRecords", err.Error(), c.Ctx.Input)
		c.Abort("查询用户借款记录总数出错" + err.Error())
		return
	}
	c.Data["submit_man"] = c.User.Id
	c.Data["collectionExeclHas"] = collectionExeclHas
	c.Data["displayExeclButton"] = has
	c.Data["loanId"] = loanId
	c.Data["isTag"] = isTag
	c.Data["uid"] = uid
	c.Data["blackSignState"] = blackSignState
	c.Data["account"] = um.Account
	c.Data["handlelist"] = m1HandleList
	c.Data["list"] = m1HandleList[0]
	c.Data["mtype"] = mtype
	c.Data["nbs"] = !bs
	c.Data["count_overdue"] = count_overdue
	c.TplName = "collection/review.html"
}

//催收处理
func (c *CollectionController) Handle() {
	defer c.ServeJSON()
	var chandle *models.CollectionHandle
	loanId, _ := c.GetInt("loanId")
	if loanId <= 0 {
		c.Data["json"] = map[string]interface{}{"ret": 304, "msg": "借款ID错误"}
		return
	}
	loan, err := models.GetLoanById(loanId)
	if err != nil {
		cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "根据load查询借款记录失败", "催收管理/Handle", err.Error(), c.Ctx.Input)
		c.Data["json"] = map[string]interface{}{"ret": 304, "msg": "校验借款数据错误"}
		return
	}
	if loan == nil {
		c.Data["json"] = map[string]interface{}{"ret": 304, "msg": "借款记录不存在"}
		return
	}
	if loan.State == "FINISH" {
		c.Data["json"] = map[string]interface{}{"ret": 304, "msg": "该案件已经结清"}
		return
	}
	//判断该期是否结清
	count, err := models.GetCaseStateCount(loanId)
	if err != nil {
		cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "根据load查询还款计划失败", "催收管理/Handle", err.Error(), c.Ctx.Input)
		c.Data["json"] = map[string]interface{}{"ret": 304, "msg": "校验借款数据错误"}
		return
	}
	if count == 0 {
		//结清
		c.Data["json"] = map[string]interface{}{"ret": 304, "msg": "该期已经结清"}
		return
	}
	config, _ := models.QueryByIdCostConfig(1)
	cjHas := utils.CheckIsSysId(c.User.RoleId, config.CollectionSignRoleId)
	if loan.CollectionUserId != c.User.Id && !cjHas {
		c.Data["json"] = map[string]interface{}{"ret": 304, "msg": "没有权限处理该案件"}
		return
	}
	//行动分类
	actionType := c.GetString("actionType")
	if actionType == "" {
		c.Data["json"] = map[string]interface{}{"ret": 304, "msg": "行动分类为空"}
		return
	}
	// 联系内容
	rowRemark := c.GetString("remark")
	remark := strings.TrimSpace(rowRemark)
	if remark == "" {
		c.Data["json"] = map[string]interface{}{"ret": 304, "msg": "联系内容不能为空"}
		return
	}
	remark = utils.FilterEmoji(remark)
	//复核日期
	checkNum := c.GetString("check_num")
	var promiseMoney float64
	if actionType == "PTP" {
		//承诺金额
		promiseMoney, _ = c.GetFloat("promiseMoney")
		if promiseMoney <= 0 {
			c.Data["json"] = map[string]interface{}{"ret": 304, "msg": "承诺金额错误"}
			return
		}
	}

	//用户ID
	uid, _ := c.GetInt("uid")
	//行动分类为‘盗办’，‘死亡’，’坐牢‘，的用户自动入黑
	if actionType == "盗办" || actionType == "死亡" || actionType == "坐牢" {
		err = models.UpdateBlackByUserId(uid, c.User.Id, actionType)
		if err != nil {
			cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "自动入黑失败", "催收管理/Handle", err.Error(), c.Ctx.Input)
			c.Data["json"] = map[string]interface{}{"ret": 304, "msg": err.Error()}
			return
		}
		err = models.UpdateBlackSignActionType(uid, c.User.Id, actionType)
		if err != nil {
			cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "更新入黑行动分类失败", "催收管理/Handle", err.Error(), c.Ctx.Input)
			c.Data["json"] = map[string]interface{}{"ret": 304, "msg": "更新入黑行动分类失败" + err.Error()}
			return
		}
	}
	chandle = &models.CollectionHandle{}
	mtype, _ := c.GetInt("type")
	chandle.LoanId = loanId
	chandle.HandleUserId = c.User.Id
	chandle.ActionType = actionType
	chandle.HandleTime = time.Now()
	chandle.Remark = remark
	chandle.PromiseMoney = promiseMoney
	chandle.Type = mtype
	chandle.CompositeDate = checkNum
	var rcd *models.Conn_record
	rcd = &models.Conn_record{}
	rcd.Uid = uid
	rcd.Conn_type = "COLLECTION"
	rcd.Created_by = c.User.Id
	rcd.Create_time = time.Now()
	rcd.Modify_by = c.User.Id
	rcd.Modify_time = time.Now()
	rcd.Content = remark
	err = rcd.Insert() //添加联系历史信息
	if err != nil {
		cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "新增联系历史失败", "催收管理/Handle", err.Error(), c.Ctx.Input)
	}
	chandle.ConnRecordId = rcd.Id
	err = chandle.HandleInsert() //添加处理信息
	if err != nil {
		cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "新增催收处理信息失败", "催收管理/Handle", err.Error(), c.Ctx.Input)
		c.Data["json"] = map[string]interface{}{"ret": 304, "msg": err.Error()}
		return
	}
	err = chandle.UpdateLoan()
	if err != nil {
		cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "更新催收处理信息失败", "催收管理/Handle", err.Error(), c.Ctx.Input)
		c.Data["json"] = map[string]interface{}{"ret": 304, "msg": err.Error()}
		return
	}
	c.Data["json"] = map[string]interface{}{"ret": 200, "msg": "添加成功!"}
}

//获取查询条件
func Where(ctype int, c *CollectionController) (where string, pars []string) {
	//手机号码
	if account := c.GetString("account"); account != "" {
		where += " AND a.account=?"
		pars = append(pars, account)
	}
	//用户姓名
	if verifyRealName := c.GetString("verifyRealName"); verifyRealName != "" {
		where += " AND a.verifyrealname=?"
		pars = append(pars, verifyRealName)
	}
	//催收员
	collectionUserId, _ := c.GetInt("collectionUserId")
	if collectionUserId > 0 {
		where += " AND a.collection_user_id=?"
		strCollectionUserId := strconv.Itoa(collectionUserId)
		pars = append(pars, strCollectionUserId)
	}
	if collectionUserId < 0 { //无催收员
		where += " AND a.collection_user_id=0 "
	}
	//行动分类
	actiontype := c.GetString("action_type")
	if actiontype != "" {
		if actiontype == "未处理" {
			where += " AND a.collection_handle_id=0"
		} else {
			where += " AND a.action_type=?"
			pars = append(pars, actiontype)
		}
	}
	//到期时间
	registerTime := c.GetString("register_time")
	if registerTime != "" {
		if strings.Contains(registerTime, "~") {
			dayArr := strings.Split(registerTime, "~")
			if len(dayArr) == 2 {
				start := dayArr[0]
				end := dayArr[1]
				where += " AND a.expiration_time>= ? AND a.expiration_time<= ? "
				pars = append(pars, start, end)
			}
		}
	}
	//逾期天数
	overdue_num := c.GetString("overdue_num")
	if overdue_num != "" {
		if strings.Contains(overdue_num, ",") {
			overdue := strings.Split(overdue_num, ",")
			if len(overdue) == 2 {
				start := overdue[0]
				end := overdue[1]
				where += " AND a.max_overdue_days>= ? AND a.max_overdue_days<= ? "
				pars = append(pars, start, end)
			}
		}
	}
	//处理情况
	handleType := c.GetString("handleType")
	if handleType != "" {
		if handleType == "1" {
			where += " AND a.handle_time >= ?"
			pars = append(pars, time.Now().Format(utils.FormatDate)+" 00:00:00")
		} else {
			where += " AND ( a.handle_time < ? OR a.handle_time IS NULL)"
			pars = append(pars, time.Now().Format(utils.FormatDate)+" 00:00:00")
		}
	}
	if ctype > 0 {
		uidStr := strconv.Itoa(c.User.Id)
		pars = append(pars, uidStr)
	}
	return where, pars
}

// 回收案件
func (c *CollectionController) GetRecycleCase() {
	c.IsNeedTemplate()
	page, _ := c.GetInt("page")
	if page < 1 {
		page = 1
	}
	pageSize, _ := c.GetInt("pageSize")
	if pageSize < 1 {
		pageSize = 30
	}
	orgStr, err := cache.GetCacheDataByStation(c.User.StationId)
	if err != nil {
		cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "回收案件根据岗位获取数据权限失败", "新催收管理/NewRecycleCase", err.Error(), c.Ctx.Input)
	}
	var where string
	var pars []string
	//手机号码
	if account := c.GetString("account"); account != "" {
		where += " AND a.account=?"
		pars = append(pars, account)
	}
	//用户姓名
	if verifyRealName := c.GetString("verifyRealName"); verifyRealName != "" {
		where += " AND a.verifyrealname=?"
		pars = append(pars, verifyRealName)
	}
	//到期时间
	registerTime := c.GetString("register_time")
	if registerTime != "" {
		if strings.Contains(registerTime, "~") {
			dayArr := strings.Split(registerTime, "~")
			if len(dayArr) == 2 {
				start := dayArr[0]
				end := dayArr[1]
				where += " AND b.loan_return_date>= ? AND b.loan_return_date<= ? "
				pars = append(pars, start, end)
			}
		}
	} else {
		start := time.Now().Format(utils.FormatDate) + " 00:00:00"
		end := time.Now().Format(utils.FormatDate) + " 59:59:59"
		where += " AND b.return_date>= ? AND b.return_date<= ? "
		pars = append(pars, start, end)
	}
	//逾期天数
	overdue_num := c.GetString("overdue_num")
	if overdue_num != "" {
		if strings.Contains(overdue_num, ",") {
			overdue := strings.Split(overdue_num, ",")
			if len(overdue) == 2 {
				start := overdue[0]
				end := overdue[1]
				where += " AND b.overdue_days>= ? AND b.overdue_days<= ? "
				pars = append(pars, start, end)
			}
		}
	}
	//结清阶段
	handleType := c.GetString("handleType")
	//催收员
	collectionUserId, _ := c.GetInt("collectUserId")
	cleaningMtype := map[int]string{11: "M0", 12: "S1", 13: "S2", 14: "S3", 15: "M2", 16: "M3"}
	collectionUsers := models.GetAllCollectUser(orgStr)
	if handleType != "" {
		where += " AND b.mtype = ? "
		pars = append(pars, handleType)
		if collectionUserId <= 0 {
			//根据结清阶段得到催收员
			mtype, _ := strconv.Atoi(handleType)
			collectionUsers = models.GetCollectUser(orgStr, mtype)
		}
	}
	if collectionUserId > 0 {
		where += " AND b.collection_user_id=?"
		strCollectionUserId := strconv.Itoa(collectionUserId)
		pars = append(pars, strCollectionUserId)
		if handleType == "" {
			//根据催收员得到结清阶段
			mtype := models.GetInfoByCollector(collectionUserId)
			ctype := utils.Mtype(mtype)
			cleaningMtype = map[int]string{mtype: ctype}
		}
	}
	if collectionUserId < 0 { //无催收员
		where += " AND b.collection_user_id=0  "
	}
	uidStr := strconv.Itoa(c.User.Id)
	pars = append(pars, uidStr)

	list, err := models.GetRecycleCase(utils.StartIndex(page, pageSize), pageSize, where, orgStr, pars...)
	if err != nil {
		cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "查询回收案件列表失败", "新催收管理/NewRecycleCase", err.Error(), c.Ctx.Input)
	}
	count := models.GetRecycleCount(where, orgStr, pars...)
	pageCount := utils.PageCount(count, pageSize)
	money := models.GetNewRecycleMoney(where, orgStr, pars)
	c.Data["sumloanmoney"] = money.Sumloanmoney
	c.Data["realrepayment"] = money.Finishmoney
	c.Data["costrelieffee"] = money.CostFee
	c.Data["list"] = list
	c.Data["count"] = count
	c.Data["pageCount"] = pageCount
	c.Data["currPage"] = page
	c.Data["pageSize"] = pageSize
	c.Data["collecionUsers"] = collectionUsers
	c.Data["cleaningMtype"] = cleaningMtype
	c.TplName = "collection/recycle_case.html"
}

//手动分配订单给相应的催收人员
func (c *CollectionController) ManualDistribution() {
	defer c.ServeJSON()
	//获取分发的订单ID
	loanIds := strings.TrimSpace(c.GetString("loanIds"))
	if loanIds == "" {
		c.Data["json"] = map[string]interface{}{"ret": 304, "msg": "请选择要分发的订单"}
		return
	}
	loanIdArr := strings.Split(loanIds, ",")
	loanIdArrNew := utils.RemoveDuplicatesAndEmpty(loanIdArr) //loan_id
	//获取分发给的催收员ID
	collectionUserId, _ := c.GetInt("uid")
	//获取类型
	ctype, _ := c.GetInt("ctype")
	//获取组织机构
	org, err := models.GetOrganizationByUserId(collectionUserId)
	if err != nil {
		cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "获取组织机构失败", "催收管理/ManualDistribution", err.Error(), c.Ctx.Input)
		c.Data["json"] = map[string]interface{}{"ret": 304, "msg": "获取组织机构失败!"}
		return
	}
	if org == nil {
		c.Data["json"] = map[string]interface{}{"ret": 304, "msg": "用户组织架构不存在!"}
		return
	}
	//进行手动分配
	err = models.InManualDistributionModels(loanIdArrNew, ctype, c.User.Id, collectionUserId, org.Id)
	if err != nil {
		cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "重新分配失败", "催收管理/ManualDistribution", err.Error(), c.Ctx.Input)
		c.Data["json"] = map[string]interface{}{"ret": 304, "msg": "重新分配失败"}
		return
	}
	c.Data["json"] = map[string]interface{}{"ret": 200, "msg": "重新分配成功"}
}

//贷后催收
func (c *CollectionController) ConnHistory() {
	uid, _ := c.GetInt("uid")
	c.Data["uid"] = uid
	page, _ := c.GetInt("page")
	if page < 1 {
		page = 1
	}
	pageSize, _ := c.GetInt("pageSize")
	if pageSize < 1 {
		pageSize = 30
	}
	connectType := c.GetString("connectType")
	condition := ""
	pars := []string{}
	connRcds, err := models.ConnRcdsList(uid, condition, pars, utils.StartIndex(page, pageSize), pageSize)

	if err != nil {
		cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "查询联系历史失败", "个人信息/联系历史ConnHistory", err.Error(), c.Ctx.Input)
	}
	count := models.ConnRcdsCount(uid, condition, pars)
	pagecount := utils.PageCount(count, pageSize)
	c.Data["ConnectType"] = connectType
	c.Data["currPage"] = page
	c.Data["pageCount"] = pagecount
	c.Data["pageSize"] = pageSize
	c.Data["count"] = count
	c.Data["connRcds"] = connRcds
	c.TplName = "user/user_connect_history.html"
}

//投诉处理列表
func (c *CollectionController) ComplaintHandlingList() {
	uid, _ := c.GetInt("uid")
	c.Data["uid"] = uid
	page, _ := c.GetInt("page")
	if page < 1 {
		page = 1
	}
	pageSize, _ := c.GetInt("pageSize")
	if pageSize < 1 {
		pageSize = 30
	}
	connectType := c.GetString("connectType")
	list, err := models.ComplaintHandlingList(uid, utils.StartIndex(page, pageSize), pageSize)
	if err != nil {
		cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "查询投诉管理列表失败", "个人信息/查询投诉管理列表", err.Error(), c.Ctx.Input)
	}
	count, err := models.ComplaintHandlingListCount(uid)
	pagecount := utils.PageCount(count, pageSize)
	c.Data["ConnectType"] = connectType
	c.Data["currPage"] = page
	c.Data["pageCount"] = pagecount
	c.Data["pageSize"] = pageSize
	c.Data["count"] = count
	c.Data["list"] = list
	c.TplName = "user/user_connect_history.html"
}

// @router /installcreatecontractpdf [get]
func (this *CollectionController) InstalCreateContractPdf() {
	pdfType, _ := this.GetInt("pdf_type")
	loanId, _ := this.GetInt("id")
	m1HandleInfoList, err := models.GetINstallHandleListById(loanId)
	if err != nil {
		cache.RecordLogs(this.User.Id, 0, this.User.Name, this.User.DisplayName, "获取借款信息失败", "导出缴款通知书和债务律师催告函pdf/CreateContractPdf", err.Error(), this.Ctx.Input)
	}
	m1HandleInfo := models.NewCollectMList{}
	if len(m1HandleInfoList) > 0 {
		for k, _ := range m1HandleInfoList {
			m1HandleInfo.CapitalAmount += m1HandleInfoList[k].CapitalAmount
			m1HandleInfo.TaxAmount += m1HandleInfoList[k].TaxAmount
			m1HandleInfo.OverdueMoneyAmount += m1HandleInfoList[k].OverdueMoneyAmount
			m1HandleInfo.OverdueBreachOfAmount += m1HandleInfoList[k].OverdueBreachOfAmount
			m1HandleInfo.RemainMoneyChargeUpAmount += m1HandleInfoList[k].RemainMoneyChargeUpAmount // 本次挂账余额
			m1HandleInfo.DataServiceFee += m1HandleInfoList[k].DataServiceFee                       // 信审数据费
			m1HandleInfo.VerifyRealName = m1HandleInfoList[k].VerifyRealName
			m1HandleInfo.IdCard = m1HandleInfoList[k].IdCard
			m1HandleInfo.LoanDate = m1HandleInfoList[k].LoanDate
			m1HandleInfo.Displayname = m1HandleInfoList[k].Displayname
			m1HandleInfo.ContactPhone = m1HandleInfoList[k].ContactPhone
			m1HandleInfo.ContractCode = m1HandleInfoList[k].ContractCode
		}
	}
	totalLoanAmount, err := strconv.ParseFloat(utils.SumMoney(6, m1HandleInfo.CapitalAmount, m1HandleInfo.TaxAmount, m1HandleInfo.OverdueMoneyAmount, m1HandleInfo.OverdueBreachOfAmount, m1HandleInfo.DataServiceFee, m1HandleInfo.RemainMoneyChargeUpAmount), 64)
	if err != nil {
		cache.RecordLogs(this.User.Id, 0, this.User.Name, this.User.DisplayName, "转换待还总金额有误", "导出缴款通知书和债务律师催告函pdf/CreateContractPdf", err.Error(), this.Ctx.Input)
	}
	data := &services.PdfBaseParameter{
		BorrowerName:    m1HandleInfo.VerifyRealName,
		BorrowerIdCard:  m1HandleInfo.IdCard,
		LoanDate:        m1HandleInfo.LoanDate,
		TotalLoanAmount: totalLoanAmount,
		HandleName:      m1HandleInfo.Displayname,
		HandlePhone:     m1HandleInfo.ContactPhone,
		ContractCode:    m1HandleInfo.ContractCode,
	}
	var xjdpdf *utils.XJDPdf
	var fileName string = ""
	if pdfType == 1 { //缴款通知书
		xjdpdf = services.GeneratePDFDemandNote(data)
		fileName = "缴款通知书.pdf"
	} else if pdfType == 2 { //债务律师催告函
		xjdpdf = services.GeneratePDFAttorneyLetter(data)
		fileName = "债务律师催告函.pdf"
	}
	if xjdpdf != nil {
		xjdpdf.Out(fileName)
	}
	this.Ctx.Output.Header("Accept-Ranges", "bytes")
	this.Ctx.Output.Header("Content-Disposition", "attachment; filename="+fileName)
	this.Ctx.Output.Header("Cache-Control", "must-revalidate, post-check=0, pre-check=0")
	this.Ctx.Output.Header("Pragma", "no-cache")
	this.Ctx.Output.Header("Expires", "0")
	http.ServeFile(this.Ctx.ResponseWriter, this.Ctx.Request, fileName)
	err = os.Remove(fileName)
	if err != nil {
		cache.RecordLogs(this.User.Id, 0, this.User.Name, this.User.DisplayName, "", "删除文件错误", err.Error(), this.Ctx.Input)
	}
}

//紧急联系人
func (this *CollectionController) PhoneLinkman() {
	defer this.ServeJSON()
	pageNo, _ := this.GetInt("page")
	if pageNo <= 0 {
		pageNo = 1
	}
	uid, _ := this.GetInt("uid")
	if uid <= 0 {
		this.Data["json"] = map[string]interface{}{"ret": 403, "msg": "uid参数异常！"}
		return
	}
	phoneNum := this.GetString("phoneNum")
	var telDerictorySort int = 2
	var linkmanSort int
	var txlist = []models.MailList{}
	var txlist2 = []models.MailList2{}
	var contact = []models.Contact{}
	var count int
	var err2 error
	//count, _ = models.QueryTelephonCount(uid, phoneNum)
	//uid = 500000024
	txlist, err2 = models.QueryTelephon(uid, phoneNum)
	if err2 != nil {
		cache.RecordLogs(this.User.Id, 0, this.User.Name, this.User.DisplayName, "获取通讯录信息异常", "个人信息/手机联系人PhoneLinkman", err2.Error(), this.Ctx.Input)
	}
	if len(txlist) == 0 {
		cache.RecordLogs(this.User.Id, 0, this.User.Name, this.User.DisplayName, "获取通讯录信息异常", "个人信息/手机联系人PhoneLinkman", "", this.Ctx.Input)
	} else {
		txlist2 = services.SortAddBook(txlist)
		for k, v := range txlist2[0].Contact {
			if len(v.ContactPhoneNumber) == 0 {
				continue
			}
			contact = append(contact, txlist2[0].Contact[k])
		}

	}
	count = len(contact)
	var totalPage int
	if count > 0 {
		totalPage = utils.PageCount(count, utils.PageSize10)
		start := utils.StartIndex(pageNo, utils.PageSize10)
		num := utils.StartIndex(pageNo, utils.PageSize10) + utils.PageSize10
		if count >= num {
			contact = contact[start:num]
		} else {
			contact = contact[start:]
		}
	}
	//根据uid,name,phone 得到标记信息
	for k, v := range contact {
		con, err := models.GetInfoByUidAndPhone(uid, v.ContactPhoneNumber[0], v.ContactName)
		if err != nil {
			contact[k].SignId = 0
			contact[k].SignRemark = ""
			contact[k].SignState = 0
		} else {
			contact[k].SignId = con.SignId
			contact[k].SignRemark = con.SignRemark
			contact[k].SignState = con.SignState
		}
	}
	//  查询紧急联系人
	link, err3 := models.QueryInstancyLinkman(uid)

	if err3 != nil {
		this.Data["json"] = map[string]interface{}{"ret": 403, "msg": "查询紧急联系人异常"}
		return
	}
	linkmanSort = 2
	this.Data["json"] = map[string]interface{}{"link": link, "page": pageNo, "uid": uid, "txlist": contact, "totalPage": totalPage, "count": count, "sort": telDerictorySort, "linkmanSort": linkmanSort}
}

//标记紧急联系人
// @router /sign [post]
func (c *CollectionController) SignContacts() {
	defer c.ServeJSON()
	//标记状态
	makeState := c.GetString("makeState")
	if makeState == "" {
		c.Data["json"] = map[string]interface{}{"ret": 304, "msg": "请选择状态"}
		return
	}
	makeContent := c.GetString("makeContent")
	if makeContent == "" {
		c.Data["json"] = map[string]interface{}{"ret": 304, "msg": "请选择备注信息"}
		return
	}
	id, _ := c.GetInt("signId")
	if id < 1 {
		c.Data["json"] = map[string]interface{}{"ret": 304, "msg": "参数为空"}
		return
	}
	signState, err := strconv.Atoi(makeState)
	if err != nil {
		c.Data["json"] = map[string]interface{}{"ret": 304, "msg": "id类型转换失败"}
		return
	}
	if makeState == "1" {
		makeContent = "正常:" + makeContent
	} else if makeState == "2" {
		makeContent = "虚假:" + makeContent
	} else if makeState == "3" {
		makeContent = "异常:" + makeContent
	}
	//标记的从哪里查的数据
	//1-users_modify_linkman 2-users_linkman)
	sort, _ := c.GetInt("sort")
	/*	if sort == 1 {
		num, err2 := models.UpdateCountUsersModifyLinkmanById(id, signState, makeContent)
		if err2 != nil || num != 1 {
			c.Data["json"] = map[string]interface{}{"ret": 403, "msg": "标记/重新标记紧急联系人异常"}
			return
		}
		c.Data["json"] = map[string]interface{}{"ret": 200, "msg": "标记成功"}
	} else*/
	if sort == 2 {
		num, err2 := models.SignUsersTelephoneDirectory(id, signState, makeContent)
		if err2 != nil || num != 1 {
			c.Data["json"] = map[string]interface{}{"ret": 403, "msg": "标记/重新标记紧急联系人异常"}
			return
		}
		c.Data["json"] = map[string]interface{}{"ret": 200, "msg": "标记成功"}
	}
}

//标记通讯录
// @router /frequent [post]
func (c *CollectionController) FrequentContacts() {
	defer c.ServeJSON()
	//标记状态
	makeState := c.GetString("makeState")
	if makeState == "" {
		c.Data["json"] = map[string]interface{}{"ret": 304, "msg": "请选择状态"}
		return
	}
	makeContent := c.GetString("makeContent")
	if makeContent == "" {
		c.Data["json"] = map[string]interface{}{"ret": 304, "msg": "请输入备注内容"}
		return
	}
	signState, err := strconv.Atoi(makeState)
	if err != nil {
		c.Data["json"] = map[string]interface{}{"ret": 304, "msg": "id类型转换失败"}
		return
	}
	uid, _ := c.GetInt("uid")
	//uid = 500000024
	//标记的手机通讯录号码
	mobile := c.GetString("mobile")
	name := c.GetString("phoneName")
	stype := c.GetString("singType")
	if makeState == "1" {
		makeContent = "正常:" + makeContent
	} else if makeState == "2" {
		makeContent = "虚假:" + makeContent
	} else if makeState == "3" {
		makeContent = "异常:" + makeContent
	}
	//根据uid和手机号码 查询是否存在记录
	count := models.GetSignCountByUidAndPhone(uid, name, mobile, stype)
	if count == 0 {
		//无记录 插入记录
		err = models.InsertMoblieDirectory(uid, c.User.Id, signState, name, makeContent, mobile, stype)
		if err != nil {
			c.Data["json"] = map[string]interface{}{"ret": 304, "msg": "标记联系人异常"}
			return
		}

	} else {
		//有记录，更新记录
		num, err := models.SignMoblieDirectory(uid, c.User.Id, signState, name, makeContent, mobile)
		if err != nil || num < 1 {
			c.Data["json"] = map[string]interface{}{"ret": 304, "msg": "重新标记联系人异常"}
			return
		}
	}
	c.Data["json"] = map[string]interface{}{"ret": 200, "msg": "标记成功"}
}

//常用联系人标记
// @router /common [post]
func (c *CollectionController) CommonContactTags() {
	defer c.ServeJSON()
	//标记状态
	makeState := c.GetString("makeState")
	if makeState == "" {
		c.Data["json"] = map[string]interface{}{"ret": 304, "msg": "请选择状态"}
		return
	}
	makeContent := c.GetString("makeContent")
	if makeContent == "" {
		c.Data["json"] = map[string]interface{}{"ret": 304, "msg": "请输入备注内容"}
		return
	}
	signState, err := strconv.Atoi(makeState)
	if err != nil {
		c.Data["json"] = map[string]interface{}{"ret": 304, "msg": "id类型转换失败"}
		return
	}
	if makeState == "1" {
		makeContent = "正常:" + makeContent
	} else if makeState == "2" {
		makeContent = "虚假:" + makeContent
	} else if makeState == "3" {
		makeContent = "异常:" + makeContent
	}
	uid, _ := c.GetInt("uid")
	//uid = 1000949793
	//标记的手机通讯录号码
	mobile := c.GetString("mobile")
	stype := c.GetString("singType")
	//根据uid和手机号码 查询是否存在记录
	count := models.GetSignCountByUidAndPhone(uid, "", mobile, stype)
	if count == 0 {
		//无记录 插入记录
		err = models.InsertMoblieDirectory(uid, c.User.Id, signState, "", makeContent, mobile, stype)
		if err != nil {
			c.Data["json"] = map[string]interface{}{"ret": 304, "msg": "标记联系人异常"}
			return
		}

	} else {
		//有记录，更新记录
		num, err := models.SignMoblieDirectory(uid, c.User.Id, signState, "", makeContent, mobile)
		if err != nil || num < 1 {
			c.Data["json"] = map[string]interface{}{"ret": 304, "msg": "重新标记联系人异常"}
			return
		}
	}
	c.Data["json"] = map[string]interface{}{"ret": 200, "msg": "标记成功"}
}

//手机通讯录导出
func (c *CollectionController) PhoneLinkManToExcel() {
	uid, _ := c.GetInt("uid")
	//uid = 500000025
	accounts := strings.TrimSpace(c.GetString("account"))
	var accountArr = []string{}
	if accounts != "" {
		accountArr = strings.Split(accounts, ",")
	}
	exl := [][]string{{"姓名", "手机号码"}}
	colWidth := []float64{10.0, 18.0}
	var contact = []models.Contact{}
	txlist, err := models.QueryTelephon(uid, "")
	if err != nil {
		cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "获取通讯录信息异常", "手机通讯录导出", err.Error(), c.Ctx.Input)
	}
	if len(txlist) == 0 {
		cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "获取通讯录信息异常", "手机通讯录导出", "", c.Ctx.Input)
	} else {
		txlist2 := services.SortAddBook(txlist)
		for k, v := range txlist2[0].Contact {
			if len(v.ContactPhoneNumber) == 0 {
				continue
			}
			contact = append(contact, txlist2[0].Contact[k])
		}
	}
	exp := []string{}
	if len(contact) > 0 {
		for _, v := range contact {
			for _, av := range accountArr {
				if v.ContactPhoneNumber[0] == av {
					exp = []string{v.ContactName, av}
					exl = append(exl, exp)
				}
			}
		}
	}
	filename, err := utils.ExportToExcel(exl, colWidth, "手机通讯录号码")
	if err != nil {
		cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "保存文件错误", err.Error(), c.Ctx.Input)
	}
	c.Ctx.Output.Header("Accept-Ranges", "bytes")
	c.Ctx.Output.Header("Content-Disposition", "attachment; filename="+filename)
	c.Ctx.Output.Header("Cache-Control", "must-revalidate, post-check=0, pre-check=0")
	c.Ctx.Output.Header("Pragma", "no-cache")
	c.Ctx.Output.Header("Expires", "0")
	http.ServeFile(c.Ctx.ResponseWriter, c.Ctx.Request, filename)
	err = os.Remove(filename)
	if err != nil {
		cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "删除文件错误", err.Error(), c.Ctx.Input)
	}
}

// 常用联系人
func (c *CollectionController) OperatorData() {
	defer c.ServeJSON()
	session := utils.GetSession()
	defer session.Close()
	page, _ := c.GetInt("page")
	if page < 1 {
		page = 1
	}
	uid, err := c.GetInt("uid")
	mobileAuthType, err := models.GetMobileAuthTypeByUserId(uid)
	if err != nil {
		cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "根据用户ID获取运营商授权类型失败", "个人信息/Post", err.Error(), c.Ctx.Input)
	}
	//常用联系人电话
	var telDerictorySort int
	var mx_data []models.CallContactDetail
	var tj_data []models.CallLogInfo
	if mobileAuthType == 1 { //魔蝎
		var userMxData models.Mxreportdata
		session := utils.GetSession()
		defer session.Close()
		//	uid = 1000949793
		err := session.DB(utils.MGO_DB).C("mxreportdata").Find(&models.MonGoQuery{Uid: uid}).One(&userMxData)
		if err != nil {
			cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "mgodb数据获取异常", "个人信息/运营商数据OperatorData", err.Error(), c.Ctx.Input)
		}
		if len(userMxData.Rt.Ccd) == 0 {
			cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "mgodb数据获取为空", "个人信息/运营商数据OperatorData", err.Error(), c.Ctx.Input)
		}
		mx_data = userMxData.Rt.Ccd
		mxStart := utils.StartIndex(page, utils.PageSize10)
		mxCount := utils.StartIndex(page, utils.PageSize10) + utils.PageSize10
		count := len(mx_data)
		if count >= mxCount {
			mx_data = mx_data[mxStart:mxCount]
		} else {
			mx_data = mx_data[mxStart:]
		}
		//根据uid,name,phone 得到标记信息
		for k, v := range mx_data {
			mx, err := models.GetMxInfo(uid, v.Peer_num)
			if err != nil {
				mx_data[k].SignId = 0
				mx_data[k].SignRemark = ""
				mx_data[k].SignState = 0
			} else {
				mx_data[k].SignId = mx.SignId
				mx_data[k].SignRemark = mx.SignRemark
				mx_data[k].SignState = mx.SignState
			}
		}
		telDerictorySort = 1
		pageCount := utils.PageCount(count, utils.PageSize10)
		c.Data["json"] = map[string]interface{}{"list": mx_data, "uid": uid, "mobileAuthType": mobileAuthType, "sort": telDerictorySort, "currpage": page, "count": count, "pagecount": pageCount}
		return
	} else if mobileAuthType == 2 {
		var tjr models.Tijireport
		session := utils.GetSession()
		defer session.Close()
		//uid = 1000949856
		err = session.DB(utils.MGO_DB).C("tianjireport").Find(&models.MonGoQuery{Uid: uid}).One(&tjr)
		if err != nil {
			cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "mgodb数据获取异常", "个人信息/运营商数据OperatorData", err.Error(), c.Ctx.Input)
		}
		if len(tjr.Tianji.CallLog) == 0 {
			cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "mgodb数据获取为空", "个人信息/运营商数据OperatorData", err.Error(), c.Ctx.Input)
		}
		tj_data = tjr.Tianji.CallLog
		tjStart := utils.StartIndex(page, utils.PageSize10)
		tjCount := utils.StartIndex(page, utils.PageSize10) + utils.PageSize10
		count := len(tj_data)
		if count >= tjCount {
			tj_data = tj_data[tjStart:tjCount]
		} else {
			tj_data = tj_data[tjStart:]
		}
		telDerictorySort = 1
		pageCount := utils.PageCount(count, utils.PageSize10)
		c.Data["json"] = map[string]interface{}{"list": tj_data, "uid": uid, "mobileAuthType": mobileAuthType, "sort": telDerictorySort, "currpage": page, "count": count, "pagecount": pageCount}
		return
	}
}

//催收查询
func (c *CollectionController) Msearch() {
	c.IsNeedTemplate()
	var where string
	var pars []string
	//手机号码
	if account := c.GetString("account"); account != "" {
		where += " AND a.account=?"
		pars = append(pars, account)
	}
	//用户姓名
	if verifyRealName := c.GetString("verifyRealName"); verifyRealName != "" {
		where += " AND a.verifyrealname=?"
		pars = append(pars, verifyRealName)
	}
	//用户ID
	if userId, _ := c.GetInt("userId"); userId != 0 {
		where += " AND a.uid = ?"
		pars = append(pars, strconv.Itoa(userId))
	}
	if where != "" {
		var list []*models.NewCollectMList
		//cstartDate, cendDate := getOverdueDaysDate2("0")
		list, err := models.GetMSearch(where, pars)
		if err != nil {
			cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "获得催收查询列表失败", "催收管理/Msearch", err.Error(), c.Ctx.Input)
		}
		for k, v := range list {
			list[k].ReturnMoney = models.GetNeedPayment(v.LoanId)
		}
		c.Data["list"] = list
	}
	// cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "获得催收查询列表成功", "催收管理/Msearch", "", c.Ctx.Input)
	c.TplName = "collection/search.html"
}

// @router /complanit [get,post]
func (c *CollectionController) Complanit() {
	c.IsNeedTemplate()
	status, _ := c.GetInt("handleType") //状态
	mobile := c.GetString("account")
	startTime := c.GetString("start_time")
	endTime := c.GetString("end_time")
	page, _ := c.GetInt("page")
	if page < 1 {
		page = 1
	}
	pageSize, _ := c.GetInt("pageSize")
	if pageSize < 1 {
		pageSize = 30
	}
	complanit, err := models.GetComplanitInfo(startTime, endTime, mobile, status, utils.StartIndex(page, pageSize), pageSize)
	//投诉内容
	if err != nil {
		cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "获取投诉处理信息失败", "催收管理/投诉处理complanit", err.Error(), c.Ctx.Input)
	}
	count, err1 := models.GetComplanitCount(startTime, endTime, mobile, status) //投诉次数
	if err1 != nil {
		cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "获取投诉处理信息总数失败", "催收管理/投诉处理complanit", err.Error(), c.Ctx.Input)
	}
	c.Data["connRcds"] = complanit
	c.Data["count"] = count
	c.Data["currpage"] = page
	c.Data["pagecount"] = utils.PageCount(count, pageSize)
	c.Data["pagesize"] = pageSize
	c.TplName = "collection/complaint_handling.html"
}

func (this *CollectionController) CheckUsersMetadata() {
	defer this.ServeJSON()
	uid, _ := this.GetInt("uid")
	count, err := models.GetUsersCount(uid)
	if err != nil {
		cache.RecordLogs(this.User.Id, 0, this.User.Name, this.User.DisplayName, "查询认证用户数据错误", "意见反馈CheckUsersMetadata", err.Error(), this.Ctx.Input)
		return
	}
	this.Data["json"] = map[string]interface{}{"ret": 200, "count": count}
}

//添加投诉处理结果
// @router /connRcDealUp [post]
func (c *CollectionController) ConnRcDealUp() {
	defer c.ServeJSON()
	id, _ := c.GetInt("id")
	flagId, _ := c.GetInt("flagId")
	uid, _ := c.GetInt("uid")
	if id < 1 || flagId < 1 {
		c.Data["json"] = map[string]interface{}{"ret": 304, "msg": "参数为空"}
		return
	}
	res := c.GetString("manage_content")
	err := models.ConnRcDeal(id, flagId, c.User.Id, res)
	if err != nil {
		cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "添加处理结果失败", "个人信息/联系历史ConnRcDealUp", err.Error(), c.Ctx.Input)
		c.Data["json"] = map[string]interface{}{"ret": 304, "msg": err.Error()}
		return
	}
	err = models.AddComplainOperationLog(uid, id, c.User.Id, flagId, c.User.DisplayName, res)
	if err != nil {
		cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "添加投诉处理记录失败", "个人信息/联系历史ConnRcDealUp", err.Error(), c.Ctx.Input)
	}
	c.Data["json"] = map[string]interface{}{"ret": 200}
}

//获得投诉处理结果
// @router /connRcDealResult [post]
func (c *CollectionController) ConnRcDealResult() {
	defer c.ServeJSON()
	id, _ := c.GetInt("id")
	list, err := models.QueryComplainOperationLog(id)
	if err != nil {
		cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "获取处理结果有误", "个人信息/联系历史ConnRcDealResult", err.Error(), c.Ctx.Input)
		c.Data["json"] = map[string]interface{}{"ret": 304, "msg": err.Error()}
		return
	}
	if list == nil {
		c.Data["json"] = map[string]interface{}{"ret": 304, "msg": "无处理结果。"}
		return
	}
	c.Data["json"] = map[string]interface{}{"ret": 200, "data": list}
}

//坐席组数据
func (c *CollectionController) SeatGroupData() {
	c.IsNeedTemplate()
	condition := ""
	paras := []interface{}{}
	month := time.Now().Format("2006-01")
	day := time.Now().AddDate(0, 0, -1).Format(utils.FormatDate)
	collectionType := c.GetString("collectionType") //催收阶段
	if collectionType != "" {
		condition += ` AND stage=? `
		paras = append(paras, collectionType)
	}

	timeType := c.GetString("timeType") //查看期限
	staTime := c.GetString("beginTime") //统计时间
	if timeType == "" {
		timeType = "day"
	}
	//按日
	if timeType == "day" {
		if staTime != "" {
			day = staTime
			condition += ` AND createtime=? `
			paras = append(paras, staTime)
		} else {
			condition += ` AND createtime=? `
			paras = append(paras, day)
		}
		//按月统计
	} else {
		startDate, endDate := models.GetSelectDayToMonth(staTime)
		month = startDate
		day = endDate
		condition += ` AND createtime =? `
		paras = append(paras, endDate)
	}
	list, err := models.GetInnerGroupDataSum(condition, paras...)
	if err != nil {
		cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "获取内催组日数据失败", "坐席组数据/SeatGroupData", err.Error(), c.Ctx.Input)
	}
	c.Data["list"] = list
	c.Data["timeType"] = timeType
	c.Data["day"] = day
	c.Data["month"] = month
	c.TplName = "collection_data/collection_seat_data.html"
}

//内催组数据
func (c *CollectionController) InnerGroupDataNew() {
	c.IsNeedTemplate()
	condition := ""
	paras := []interface{}{}
	month := time.Now().Format("2006-01")
	day := time.Now().AddDate(0, 0, -1).Format(utils.FormatDate)
	collectionType := c.GetString("collectionType") //催收阶段
	if collectionType != "" {
		condition += ` AND stage=? `
		paras = append(paras, collectionType)
	}

	timeType := c.GetString("timeType") //查看期限
	staTime := c.GetString("beginTime") //统计时间
	if timeType == "" {
		timeType = "day"
	}
	//按日
	if timeType == "day" {
		if staTime != "" {
			day = staTime
			condition += ` AND createtime=? `
			paras = append(paras, staTime)
		} else {
			condition += ` AND createtime=? `
			paras = append(paras, day)
		}

	} else {
		startDate, endDate := models.GetSelectDayToMonth(staTime)
		month = startDate
		day = endDate
		condition += ` AND createtime =? `
		paras = append(paras, endDate)
	}

	var list []models.GroupData
	var sumList models.GroupData
	var err error
	if collectionType == "M0催收组" {
		list, err = models.GetInnerGroupDataM0(timeType, condition, paras...)
		if err != nil {
			cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "获取内催组M0数据失败", "内催组数据/InnerGroupDataNew", err.Error(), c.Ctx.Input)
		}
		sumList, err = models.GetInnerGroupDataSumM0(condition, paras)
		if err != nil {
			cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "获取内催组M0数据总计失败", "内催组数据/InnerGroupDataNew", err.Error(), c.Ctx.Input)
		}
	} else {
		list, err = models.GetInnerGroupData(timeType, condition, paras)
		if err != nil {
			cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "获取内催组数据失败", "内催组数据/InnerGroupDataNew", err.Error(), c.Ctx.Input)
		}
		sumList, err = models.GetInnerGroupDataSum(condition, paras)
		if err != nil {
			cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "获取内催组数据总计失败", "内催组数据/InnerGroupDataNew", err.Error(), c.Ctx.Input)
		}
	}
	c.Data["timeType"] = timeType
	c.Data["list"] = list
	c.Data["sumList"] = sumList
	c.Data["collectionType"] = collectionType
	c.Data["day"] = day
	c.Data["month"] = month
	c.TplName = "collection_data/collection_inside_data.html"
}

//委外组数据
func (c *CollectionController) InOutGroupDataNew() {
	c.IsNeedTemplate()
	condition := ""
	paras := []interface{}{}
	month := time.Now().Format("2006-01")
	day := time.Now().AddDate(0, 0, -1).Format(utils.FormatDate)
	collectionType := c.GetString("collectionType") //催收阶段
	if collectionType != "" {
		condition += ` AND stage=? `
		paras = append(paras, collectionType)
	}
	flag := false
	timeType := c.GetString("timeType") //查看期限
	staTime := c.GetString("beginTime") //统计时间
	if timeType == "" {
		timeType = "day"
	}
	//按日
	if timeType == "day" {
		if staTime != "" {
			day = staTime
			condition += ` AND createtime=? `
			paras = append(paras, staTime)
		} else {
			condition += ` AND createtime=? `
			paras = append(paras, day)
		}

	} else {
		startDate, endDate := models.GetSelectDayToMonth(staTime)
		month = startDate
		day = endDate
		condition += ` AND createtime =? `
		paras = append(paras, endDate)
	}

	list, err := models.GetOutGroupData(timeType, condition, paras)
	if err != nil {
		cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "获取委外组数据失败", "委外组数据/InOutGroupDataNew", err.Error(), c.Ctx.Input)
	}
	sumList, err := models.GetOutGroupDataSum(condition, paras)
	if err != nil {
		cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "获取委外组数据总计失败", "委外组数据/InOutGroupDataNew", err.Error(), c.Ctx.Input)
	}
	if len(list) > 0 {
		flag = true
	}
	c.Data["timeType"] = timeType
	c.Data["list"] = list
	c.Data["sumList"] = sumList
	c.Data["flag"] = flag
	c.Data["day"] = day
	c.Data["month"] = month
	c.TplName = "collection_data/collection_outside_data.html"
}

//回收数据
func (c *CollectionController) AfterloanGroupDataNew() {
	c.IsNeedTemplate()
	condition := ""
	paras := []interface{}{}
	month := time.Now().Format("2006-01")
	day := time.Now().AddDate(0, 0, -1).Format(utils.FormatDate)
	collectionType := c.GetString("collectionType") //催收阶段
	timeType := c.GetString("timeType")             //查看期限
	staTime := c.GetString("beginTime")             //统计时间
	if timeType == "" {
		timeType = "day"
	}
	//按日
	if timeType == "day" {
		if staTime != "" {
			day = staTime
			condition += ` AND createtime=? `
			paras = append(paras, staTime)
		} else {
			condition += ` AND createtime=? `
			paras = append(paras, day)
		}
		//按月
	} else {
		startDate, endDate := models.GetSelectDayToMonth(staTime)
		month = startDate
		day = endDate
		condition += ` AND createtime =? `
		paras = append(paras, endDate)
	}
	var sumList models.GroupData
	if collectionType == "内催组" {
		condition += ` AND stage in ('S1催收组','S2催收组','S3催收组','M0催收组') `
	} else if collectionType == "委外组" {
		condition += ` AND stage in ('A1催收组','A2催收组','A3催收组','A4催收组') `
	} else if collectionType == "坐席组" {

	} else {
		condition += ` AND stage in ('S1催收组','S2催收组','S3催收组','A1催收组','A2催收组','A3催收组','A4催收组','M0催收组') `
	}
	M0List, err := models.GetInnerGroupDataSumM0(condition, paras)
	if err != nil {
		cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "获取内催组M0数据总计失败", "内催组数据/InnerGroupDataNew", err.Error(), c.Ctx.Input)
	}
	innerList, err := models.GetInnerGroupDataSum(condition, paras)
	if err != nil {
		cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "获取内催组数据总计失败", "回收总数据/AfterloanGroupDataNew", err.Error(), c.Ctx.Input)
	}
	outList, err := models.GetOutGroupDataSum(condition, paras)
	if err != nil {
		cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "获取委外组数据总计失败", "回收总数据/AfterloanGroupDataNew", err.Error(), c.Ctx.Input)
	}
	flag := false
	if outList.AmountOnFile > 0 {
		flag = true
	}

	//总计
	if timeType == "day" {
		sumList.AmountOnFile = innerList.AmountOnFile + outList.AmountOnFile
		sumList.MoneyOnFile = innerList.MoneyOnFile + outList.MoneyOnFile
		sumList.CapitalOnFile = innerList.CapitalOnFile + outList.CapitalOnFile
		sumList.LateFeeOnFile = innerList.LateFeeOnFile + outList.LateFeeOnFile
		sumList.ReturnCapital = innerList.ReturnCapital + outList.ReturnCapital
		sumList.ReturnAmount = innerList.ReturnAmount + outList.ReturnAmount
		sumList.ReturnLateFee = innerList.ReturnLateFee + outList.ReturnLateFee
	}

	if timeType == "month" {
		sumList.AmountOnFileMonth = innerList.AmountOnFileMonth + outList.AmountOnFileMonth
		sumList.MoneyOnFileMonth = innerList.MoneyOnFileMonth + outList.MoneyOnFileMonth
		sumList.CapitalOnFileMonth = innerList.CapitalOnFileMonth + outList.CapitalOnFileMonth
		sumList.LateFeeOnFileMonth = innerList.LateFeeOnFileMonth + outList.LateFeeOnFileMonth
		sumList.ReturnCapitalMonth = innerList.ReturnCapitalMonth + outList.ReturnCapitalMonth
		sumList.ReturnAmountMonth = innerList.ReturnAmountMonth + outList.ReturnAmountMonth
		sumList.ReturnLateFeeMonth = innerList.ReturnLateFeeMonth + outList.ReturnLateFeeMonth
	}

	c.Data["M0List"] = M0List
	c.Data["sumList"] = sumList
	c.Data["innerList"] = innerList
	c.Data["outList"] = outList
	c.Data["collectionType"] = collectionType
	c.Data["timeType"] = timeType
	c.Data["flag"] = flag
	c.Data["month"] = month
	c.Data["day"] = day
	c.TplName = "collection_data/collection_recycle_data.html"
}

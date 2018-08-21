package controllers

import (
	"fenqi_v1/cache"
	"fenqi_v1/models"
	"fenqi_v1/utils"
	// "fmt"
	"github.com/astaxie/beego"
	"strings"
	"time"
)

//用户管理
type UsersManageController struct {
	BaseController
}

//认证用户列表
func (this *UsersManageController) UsersMetadataList() {
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
		condition += " AND u.account = ?"
		pars = append(pars, account)
	}
	userName := this.GetString("user_name") //姓名
	if userName != "" {
		condition += " AND um.verifyrealname = ?"
		pars = append(pars, userName)
	}
	idCard := this.GetString("id_card") //身份证
	if idCard != "" {
		condition += " AND um.id_card = ?"
		pars = append(pars, idCard)
	}
	isCredit, _ := this.GetInt("is_credit") //是否有额度
	if isCredit != 0 {
		if isCredit == 1 {
			condition += " AND um.balance <> 0"
		} else if isCredit == 2 {
			condition += " AND um.balance = 0"
		}
	}
	isInvited, _ := this.GetInt("is_invited") //是否业务员邀请
	if isInvited != 0 {
		if isInvited == 1 {
			condition += " AND u.salesman_id <> 0"
		} else if isInvited == 2 {
			condition += " AND u.salesman_id is null"
		}
	}
	source := this.GetString("source") //渠道
	if source != "" {
		condition += ` AND u.source = ?`
		pars = append(pars, source)
	}
	authTime := this.GetString("auth_time") //认证时间
	var startAuthTime, endAuthTime string
	if authTime != "" {
		authTimes := strings.Split(authTime, "~")
		startAuthTime = authTimes[0] + " 00:00:00"
		endAuthTime = authTimes[1] + " 23:59:59"
		condition += ` AND ua.real_name_time >= ? AND ua.real_name_time <= ?`
		pars = append(pars, startAuthTime)
		pars = append(pars, endAuthTime)
	}
	////默认提交时间7天内
	//if account == "" && userName == "" && idCard == "" && authTime == "" && isCredit == 0 && source == "" {
	//	startAuthTime = time.Now().AddDate(0, 0, -6).Format("2006-01-02") + " 00:00:00"
	//	endAuthTime = time.Now().Format("2006-01-02") + " 23:59:59"
	//	condition += ` AND uas.real_name_time >= ? AND uas.real_name_time <= ?`
	//	pars = append(pars, startAuthTime)
	//	pars = append(pars, endAuthTime)
	//}
	userAuthList, err := models.QueryUsersAuthList(utils.StartIndex(pageNum, pageSize), pageSize, condition, pars...)
	if err != nil {
		cache.RecordLogs(this.User.Id, 0, this.User.Name, this.User.DisplayName, "查询认证用户列表失败", "用户管理/认证用户列表UsersMetadataList", err.Error(), this.Ctx.Input)
		this.Abort("查询认证用户列表失败")
		return
	}
	for k, v := range userAuthList {
		count := models.QueryUserCreditCount(v.Uid)
	   //fmt.Println(count)
		if count > 1 {
			if v.AuthState == 0 {
				userAuthList[k].AuthState, _ = models.QueryUserLastCredit(v.Uid)
			}
		}
		if v.AuthState == 1 {
			userAuthList[k].State = "通过"
		} else if v.AuthState == 0 || v.AuthState == 3 {
			userAuthList[k].State = "未提交"
		} else if v.AuthState == 4 || v.AuthState == 5 || v.AuthState == 6 || v.AuthState == 10 ||
			v.AuthState == 11 || v.AuthState == 12 || v.AuthState == 13 || v.AuthState == 14 {
			userAuthList[k].State = "驳回"
		} else if v.AuthState == 8 {
			userAuthList[k].State = "关闭30天"
		} else if v.AuthState == 15 {
			userAuthList[k].State = "关闭180天"
		} else if v.AuthState == 16 {
			userAuthList[k].State = "关闭365天"
		} else if v.AuthState == 9 {
			userAuthList[k].State = "永久关闭"
		} else if v.AuthState == 2 {
			switch v.State {
			case "QUEUEING":
				userAuthList[k].State = "排队中"
			case "HANDING":
				userAuthList[k].State = "处理中"
			case "PASS":
				userAuthList[k].State = "通过"
			case "REJECT":
				userAuthList[k].State = "驳回"
			case "OUTQUEUE":
				userAuthList[k].State = "出列"
			case "PAUSE":
				userAuthList[k].State = "关闭30天"
			case "CLOSE":
				userAuthList[k].State = "永久关闭"
			}
		}

		if v.SalesmanId > 0 {
			v.IsSalemanType = "是"
		} else {
			v.IsSalemanType = "否"
		}

	}

	beego.Info(userAuthList)
	count, err := models.QueryUsersAuthCount(condition, pars...)
	if err != nil {
		cache.RecordLogs(this.User.Id, 0, this.User.Name, this.User.DisplayName, "获取总数异常", "用户管理/认证用户列表UsersMetadataList", err.Error(), this.Ctx.Input)
		this.Abort("获取认证用户总数异常")
		return
	}
	pageCount := utils.PageCount(count, pageSize)
	this.Data["list"] = userAuthList
	this.Data["currPage"] = pageNum
	this.Data["count"] = count
	this.Data["pageSize"] = pageSize
	this.Data["pageCount"] = pageCount
	this.TplName = "manage/manage_usersmetadata_list.html"
}

//认证用户列表查看
func (this *UsersManageController) UsersMetadataLook() {
	this.IsNeedTemplate()
	uid, _ := this.GetInt("uid")
	if uid <= 0 {
		cache.RecordLogs(this.User.Id, 0, this.User.Name, this.User.DisplayName, "uid参数传递错误", "认证用户查看/UsersMetadataLook", "", this.Ctx.Input)
		this.Abort("uid参数传递错误")
		return
	}
	return_money, _ := this.GetInt("return_money") //回款
	ubi, err := models.QueryUsersBaseInfo(uid)     //用户基本信息
	if err != nil && err.Error() != utils.ErrNoRow() {
		cache.RecordLogs(this.User.Id, 0, this.User.Name, this.User.DisplayName, "查询用户基本信息出错", "认证用户查看/UsersMetadataLook", err.Error(), this.Ctx.Input)
		this.Abort("查询用户基本信息出错" + err.Error())
	}
	if ubi == nil {
		this.Abort("用户信息不存在,请检查数据")
	}
	loanCount, err := models.QueryUserLoanSuccessCount(uid)
	if err != nil && err.Error() != utils.ErrNoRow() {
		cache.RecordLogs(this.User.Id, 0, this.User.Name, this.User.DisplayName, "查询用户成功借款笔数出错", "认证用户查看/UsersMetadataLook", err.Error(), this.Ctx.Input)
		this.Abort("查询用户成功借款笔数出错" + err.Error())
	}
	currentBackingLoans, err := models.QueryUserBackingLoan(uid)
	if err != nil {
		cache.RecordLogs(this.User.Id, 0, this.User.Name, this.User.DisplayName, "查询用户当前借款出错", "认证用户查看/UsersMetadataLook", err.Error(), this.Ctx.Input)
		this.Abort("查询用户当前借款出错" + err.Error())
	}
	index, _ := this.GetInt("index", 0)
	inviteCode, err := models.QueryUserInviteCodeByUid(uid)
	if err != nil && err.Error() != utils.ErrNoRow() {
		cache.RecordLogs(this.User.Id, 0, this.User.Name, this.User.DisplayName, "查询用户邀请码出错", "认证用户查看/UsersMetadataLook", err.Error(), this.Ctx.Input)
		this.Abort("查询用户邀请码出错" + err.Error())
	}
	//借款逾期总数
	count_overdue, err := models.QueryLoanRecordOverDueCount(uid)
	if err != nil {
		cache.RecordLogs(this.User.Id, 0, this.User.Name, this.User.DisplayName, "查询用户借款记录总数出错", "认证用户详情/借款记录GetUsersManageLoanRecords", err.Error(), this.Ctx.Input)
		this.Abort("查询用户借款记录总数出错" + err.Error())
		return
	}
	//用户当前授信状态
	authState, _ := models.GetUserAuthState(uid)
	authStateStr := ""
	switch authState {
	case 0:
		authStateStr = "未授信"
	case 1:
		authStateStr = "授信成功"
	case 2:
		authStateStr = "授信中"
	case 3:
		authStateStr = "授信资料过期"
	case 4:
		authStateStr = "授信驳回"
	case 5:
		authStateStr = "运营商驳回"
	case 6:
		authStateStr = "公积金驳回30天"
	case 8:
		authStateStr = "授信关闭30天"
	case 9:
		authStateStr = "授信永久关闭"
	case 10:
		authStateStr = "运营商、公积金同时驳回"
	case 11:
		authStateStr = "支付宝驳回30天"
	case 12:
		authStateStr = "公积金支付宝同时驳回30"
	case 13:
		authStateStr = "运营商驳回，支付宝同时驳回30天"
	case 14:
		authStateStr = "支付宝运营商公积金同时驳回"
	case 15:
		authStateStr = "授信关闭180天"
	case 16:
		authStateStr = "授信关闭365天"
	}
	beego.Info(authStateStr)
	this.Data["index"] = index
	this.Data["authStateStr"] = authStateStr
	this.Data["user"] = ubi
	this.Data["loanCount"] = loanCount
	this.Data["loans"] = currentBackingLoans
	this.Data["inviteCode"] = inviteCode
	this.Data["return_money"] = return_money
	this.Data["count_overdue"] = count_overdue
	this.TplName = "manage/manage_usersmetadata_look.html"
}

//注册用户列表
func (this *UsersManageController) UsersList() {
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
		condition += " AND u.account = ?"
		pars = append(pars, account)
	}
	isInvited, _ := this.GetInt("is_invited") //是否业务员邀请
	if isInvited != 0 {
		if isInvited == 1 {
			condition += " AND u.salesman_id <> 0"
		} else if isInvited == 2 {
			condition += " AND u.salesman_id is null"
		}
	}
	registerTime := this.GetString("register_time") //注册时间
	var startRegisterTime, endRegisterTime string
	if registerTime != "" {
		registerTimes := strings.Split(registerTime, "~")
		startRegisterTime = registerTimes[0] + " 00:00:00"
		endRegisterTime = registerTimes[1] + " 23:59:59"
		condition += ` AND u.create_time >= ? AND u.create_time <= ?`

		pars = append(pars, startRegisterTime)
		//fmt.Println(pars)
		pars = append(pars, endRegisterTime)
	}
	pkgType, _ := this.GetInt("pkg_type", -1) //平台来源
	if pkgType != -1 {
		condition += " AND u.pkg_type = ?"
		pars = append(pars, pkgType)
	}
	source := this.GetString("source") //渠道
	if source != "" {
		condition += ` AND u.source = ?`
		pars = append(pars, source)
	}
	//默认提交时间7天内
	if account == "" && registerTime == "" && pkgType == -1 && source == "" {
		startRegisterTime = time.Now().AddDate(0, 0, -6).Format("2006-01-02") + " 00:00:00"
		endRegisterTime = time.Now().Format("2006-01-02") + " 23:59:59"
		condition += ` AND u.create_time >= ? AND u.create_time <= ?`
		pars = append(pars, startRegisterTime)
		pars = append(pars, endRegisterTime)
	}
	usersList, err := models.QueryUsersList(utils.StartIndex(pageNum, pageSize), pageSize, condition, pars...)

	if err != nil {
		beego.Info(err)
		cache.RecordLogs(this.User.Id, 0, this.User.Name, this.User.DisplayName, "查询注册用户列表失败", "用户管理/注册用户列表UsersList", err.Error(), this.Ctx.Input)
		this.Abort("查询注册用户列表失败")
		return
	}
	for _, v := range usersList {
		beego.Info(v.SalesmanId)
		if v.SalesmanId > 0 {
			v.IsSalemanType = "是"
		} else {
			v.IsSalemanType = "否"
		}
	}
	count, err := models.QueryUsersCount(condition, pars...)
	if err != nil {
		cache.RecordLogs(this.User.Id, 0, this.User.Name, this.User.DisplayName, "获取总数异常", "用户管理/注册用户列表UsersList", err.Error(), this.Ctx.Input)
		this.Abort("获取注册用户总数异常")
		return
	}
	pageCount := utils.PageCount(count, pageSize)
	this.Data["registerTime"] = startRegisterTime + "~" + endRegisterTime
	this.Data["list"] = usersList
	this.Data["currPage"] = pageNum
	this.Data["count"] = count
	this.Data["pageSize"] = pageSize
	this.Data["pageCount"] = pageCount
	this.TplName = "manage/manage_users_list.html"
}

//注册用户查看
func (this *UsersManageController) UsersLook() {
	this.IsNeedTemplate()
	uid, _ := this.GetInt("uid")
	if uid <= 0 {
		cache.RecordLogs(this.User.Id, 0, this.User.Name, this.User.DisplayName, "uid参数传递错误", "注册用户查看/UsersLook", "", this.Ctx.Input)
		this.Abort("uid参数传递错误")
		return
	}
	u, err := models.QueryUsersInfo(uid)
	if u == nil {
		cache.RecordLogs(this.User.Id, 0, this.User.Name, this.User.DisplayName, "该用户不存在", "注册用户查看/UsersLook", "", this.Ctx.Input)
		this.Abort("该用户不存在")
		return
	}
	if err == nil {
		this.Data["user"] = u
	}
	inviteCode, err := models.QueryUserInviteCodeByUid(uid)
	if err != nil && err.Error() != utils.ErrNoRow() {
		cache.RecordLogs(this.User.Id, 0, this.User.Name, this.User.DisplayName, "查询用户邀请码出错", "注册用户查看/UsersLook", "", this.Ctx.Input)
		this.Abort("查询用户邀请码出错" + err.Error())
	}
	this.Data["inviteCode"] = inviteCode
	this.TplName = "manage/manage_users_look.html"
}

//冻结和取消冻结用户账户
func (this *UsersManageController) UpdateUsersForzenState() {
	resultMap := make(map[string]interface{})
	//fmt.Println(resultMap)
	resultMap["ret"] = 403
	defer func() {
		this.Data["json"] = resultMap
		this.ServeJSON()
	}()
	uid, _ := this.GetInt("uid")
	if uid <= 0 {
		cache.RecordLogs(this.User.Id, 0, this.User.Name, this.User.DisplayName, "uid参数传递错误", "用户详情查看/UpdateUsersForzenState", "", this.Ctx.Input)
		resultMap["err"] = "uid参数传递错误"
		return
	}
	state, _ := this.GetInt("state")
	var oper = ""
	var opertorType = ""
	if state == 0 {
		oper = "取消冻结"
		opertorType = "NOFROZEN"
	} else if state == 1 {
		oper = "冻结账号"
		opertorType = "FROZEN"
	}
	err := models.UpdateUsersFrozenState(uid, state)
	if err != nil {
		cache.RecordLogs(uid, 0, this.User.Name, this.User.DisplayName, "用户"+oper+"失败", "用户详情查看/UpdateUsersForzenState", err.Error(), this.Ctx.Input)
		resultMap["err"] = "用户" + oper + "失败" + err.Error()
		return
	}
	err = models.InsertUsersFrozenLog(uid, this.User.DisplayName, opertorType)
	if err != nil {
		cache.RecordLogs(uid, 0, this.User.Name, this.User.DisplayName, "新增"+oper+"记录失败", "用户详情查看/UpdateUsersForzenState", err.Error(), this.Ctx.Input)
		resultMap["err"] = "新增" + oper + "记录失败" + err.Error()
		return
	}
	resultMap["ret"] = 200
	resultMap["msg"] = "用户" + oper + "成功"
}

package controllers

import (
	"encoding/json"
	"fenqi_v1/cache"
	"fenqi_v1/models"
	"fenqi_v1/services"
	"fenqi_v1/utils"
	"strconv"
	"strings"
	//	"sync"
	//	"fmt"
	"github.com/astaxie/beego"
	"time"
)

//信审工作台
type WorkspaceController struct {
	BaseController
}

//授信管理
func (this *WorkspaceController) CreditList() {
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
		condition += " AND ca.phone_no = ?"
		pars = append(pars, account)
	}
	userName := this.GetString("user_name") //姓名
	if userName != "" {
		condition += " AND ca.user_name = ?"
		pars = append(pars, userName)
	}
	operatorName := this.GetString("operator_name") //处理人
	if operatorName != "" {
		condition += " AND ca.displayname = ?"
		pars = append(pars, operatorName)
	}
	handleState := this.GetString("handle_state") //状态
	if handleState != "" {
		condition += ` AND ca.state = ?`
		pars = append(pars, handleState)
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
	condition += ` AND ca.create_time >= ? AND ca.create_time <= ?`
	pars = append(pars, startSubmitTime)
	pars = append(pars, endSubmitTime)
	handleTime := this.GetString("deal_time") //处理时间
	if handleTime != "" {
		handleTimes := strings.Split(handleTime, "~")
		condition += ` AND ca.handling_time >= ? AND ca.handling_time <= ?`
		pars = append(pars, handleTimes[0]+" 00:00:00")
		pars = append(pars, handleTimes[1]+" 23:59:59")
	}

	creditList, err := models.GetCreditAduitList(utils.StartIndex(pageNum, pageSize), pageSize, condition, pars...)
	if err != nil {
		cache.RecordLogs(this.User.Id, 0, this.User.Name, this.User.DisplayName, "查询授信管理失败", "信审工作平台/授信管理CreditList", err.Error(), this.Ctx.Input)
		this.Abort("查询授信管理失败")
		return
	}

	for k, v := range creditList {
		if v.InqueueTime.Format(utils.FormatDate) != "0001-01-01" && v.InqueueTime.Before(time.Now()) {
			creditList[k].State = "QUEUEING"
			creditList[k].Displayname = ""
			t, _ := time.Parse(utils.FormatDateTime, "0001-01-01")
			creditList[k].AllocationTime = t
			creditList[k].HandlingTime = t
		}
	}
	count, err := models.GetCreditAduitCount(condition, pars...)
	if err != nil {
		cache.RecordLogs(this.User.Id, 0, this.User.Name, this.User.DisplayName, "获取总数异常", "信审工作平台/授信管理CreditList", err.Error(), this.Ctx.Input)
		this.Abort("获取授信管理总数异常")
		return
	}
	creditOperators, err := models.GetCreditOperators()
	if err != nil {
		cache.RecordLogs(this.User.Id, 0, this.User.Name, this.User.DisplayName, "获取处理员异常", "信审工作平台/授信管理CreditList", err.Error(), this.Ctx.Input)
		this.Abort("获取处理员异常")
		return
	}
	pageCount := utils.PageCount(count, pageSize)
	this.Data["operators"] = creditOperators
	this.Data["list"] = creditList
	this.Data["currPage"] = pageNum
	this.Data["count"] = count
	this.Data["pageSize"] = pageSize
	this.Data["pageCount"] = pageCount
	this.TplName = "workspace/workspace_credit_list.html"
}

//退回队列
func (this *WorkspaceController) ReturnQueue() {
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
		condition += " AND ca.phone_no = ?"
		pars = append(pars, account)
	}
	userName := this.GetString("user_name") //姓名
	if userName != "" {
		condition += " AND ca.user_name = ?"
		pars = append(pars, userName)
	}
	operatorName := this.GetString("operator_name") //处理人
	if operatorName != "" {
		condition += " AND su.displayname = ?"
		pars = append(pars, operatorName)
	}
	submitTime := this.GetString("submit_time") //提交时间
	var startSubmitTime, endSubmitTime string
	if submitTime != "" {
		submitTimes := strings.Split(submitTime, "~")
		startSubmitTime = submitTimes[0] + " 00:00:00"
		endSubmitTime = submitTimes[1] + " 23:59:59"
	} else { //默认提交时间一个月内
		startSubmitTime = time.Now().AddDate(0, -1, 0).Format("2006-01-02") + " 00:00:00"
		endSubmitTime = time.Now().Format("2006-01-02 15:04:05")
	}
	condition += ` AND ca.create_time >= ? AND ca.create_time <= ?`
	pars = append(pars, startSubmitTime)
	pars = append(pars, endSubmitTime)
	handleTime := this.GetString("deal_time") //退回时间
	if handleTime != "" {
		handleTimes := strings.Split(handleTime, "~")
		condition += ` AND ca.handling_time >= ? AND ca.handling_time <= ?`
		pars = append(pars, handleTimes[0]+" 00:00:00")
		pars = append(pars, handleTimes[1]+" 23:59:59")
	}
	appointTime := this.GetString("appoint_time") //预约入列时间
	if appointTime != "" {
		appointTimes := strings.Split(appointTime, "~")
		condition += ` AND ca.inqueue_time >= ? AND ca.inqueue_time <= ?`
		pars = append(pars, appointTimes[0]+" 00:00:00")
		pars = append(pars, appointTimes[1]+" 23:59:59")
	}
	outQueueList, err := models.GetCreditOutQueueList(utils.StartIndex(pageNum, pageSize), pageSize, condition, pars...)
	if err != nil {
		cache.RecordLogs(this.User.Id, 0, this.User.Name, this.User.DisplayName, "查询授信管理失败", "信审工作平台/退回队列ReturnQueue", err.Error(), this.Ctx.Input)
		this.Abort("查询退回管理失败")
		return
	}
	count, err := models.GetCreditOutQueueCount(condition, pars...)
	if err != nil {
		cache.RecordLogs(this.User.Id, 0, this.User.Name, this.User.DisplayName, "获取总数异常", "信审工作平台/退回队列ReturnQueue", err.Error(), this.Ctx.Input)
		this.Abort("获取退回队列总数异常")
		return
	}
	creditOperators, err := models.GetCreditOperators()
	if err != nil {
		cache.RecordLogs(this.User.Id, 0, this.User.Name, this.User.DisplayName, "获取处理员异常", "信审工作平台/退回队列ReturnQueue", err.Error(), this.Ctx.Input)
		this.Abort("获取处理员异常")
		return
	}
	pageCount := utils.PageCount(count, pageSize)
	this.Data["operators"] = creditOperators
	this.Data["list"] = outQueueList
	this.Data["currPage"] = pageNum
	this.Data["count"] = count
	this.Data["pageSize"] = pageSize
	this.Data["pageCount"] = pageCount
	this.TplName = "workspace/workspace_credit_outqueue_list.html"
}

//处理授信
func (this *WorkspaceController) CreditQueueHanding() {
	isPermission := models.IsDataPermissionByStationId(this.User.StationId, 2)
	if !isPermission {
		this.Abort("没有该数据权限!")
	}
	this.IsNeedTemplate()
	var okuser utils.ShareDdata
	defer func() {
		//只有异常关闭 与 正常分配执行  没有订单不执行
		if okuser.Flag {
			record, _ := models.QueryCreditAduitRecord(okuser.Uid)
			baseinfo, _ := models.QueryUsersBaseInfo(okuser.Uid)
			baseinfo.IdCard = utils.IdCardFilter(baseinfo.IdCard)
			risk, _ := models.GetRiskAdvises(okuser.Uid, 1)
			this.Data["phone_info"] = utils.QueryLocating(baseinfo.Account)
			this.Data["timeDiff"] = okuser.TimeDff
			this.Data["ok_order"] = true
			this.Data["risk"] = risk
			this.Data["baseinfo"] = baseinfo
			this.Data["id"] = okuser.Id
			this.Data["uid"] = okuser.Uid
			this.Data["record"] = record
			this.Data["loginstate"] = this.User.LoginState
			this.TplName = "workspace/workspace_credit_handing.html"
			return
		}
	}()

	//处理异常关闭
	okuser.Id, okuser.Uid, okuser.TimeDff, okuser.Flag = check_credit_handing(this)
	//45分钟之内 handing状态
	if okuser.Flag {
		return
	}

	//最后一单离线分配
	if this.User.LoginState == "OFFLINE" && !okuser.Flag {
		this.Data["loginstate"] = this.User.LoginState
		this.Data["ok_order"] = false
		cache.RecordLogs(this.User.Id, 0, this.User.Name, this.User.DisplayName, "没有订单", "workspace/creditqueuehanding", "", this.Ctx.Input)
		this.TplName = "workspace/workspace_credit_handing.html"
		return
	}

	//信号量 共享值 正常关闭
	pt := make(chan utils.ShareDdata, 1)
	pp := make(chan utils.ShareDdata, 1)
	//处理插队数据
	var pt_tmp utils.ShareDdata
	var pp_tmp utils.ShareDdata
	go credit_queue_a(pt, this)
	pt_tmp = <-pt
	if !pt_tmp.Flag {
		go credit_queue_p(pp, this)
		pp_tmp = <-pp
	}
	okuser = utils.ShareFlag(pt_tmp, pp_tmp)
	//没有订单
	if !okuser.Flag {
		this.Data["ok_order"] = false
		this.Data["loginstate"] = this.User.LoginState
		cache.RecordLogs(this.User.Id, 0, this.User.Name, this.User.DisplayName, "没有订单", "workspace/creditqueuehanding", "", this.Ctx.Input)
		this.TplName = "workspace/workspace_credit_handing.html"
		return
	}
}

func credit_queue_a(pt chan utils.ShareDdata, this *WorkspaceController) {

	//统计预约插队数量
	count := models.GetNumCreditQueueAp()
	go func() {
		for {
			ca, _ := models.GetAppointmentQueueId()
			if ca.InqueueType == 1 {
				//lock false 为lock状态 ture 为unlock状态
				_, flag := cache.GetCacheCreditMessage(ca.Id)
				if flag {
					//通道放入值
					pt <- utils.ShareDdata{Id: ca.Id, Uid: ca.Uid, Flag: true, TimeDff: 45 * 60}
					credit_queue_change_state(this, ca.Id, ca.Uid)
					return
				}
			}
			//退回队列中有排队 放入队列中
			if ca.InqueueType == 2 {
				err := models.UpdateInqueueType(ca.Id, ca.InqueueTime)
				if err != nil {
					cache.RecordLogs(this.User.Id, 0, this.User.Name, this.User.DisplayName, "更新状态失败", "workspace/creditqueuehanding", err.Error(), this.Ctx.Input)
					return
				}
			}
			//防止死锁
			if count < 0 {
				pt <- utils.ShareDdata{Uid: 0, Flag: false, Id: 0}
				return
			}
			count = count - 1
		}
	}()
	return
}

func credit_queue_p(pp chan utils.ShareDdata, this *WorkspaceController) {

	//后期考虑使用 select case 加上timeout机制 统计排队中数量
	count := models.GetCreditQueueUpIdCount()
	go func() {
		for {
			cp, _ := models.GetCreditQueueUpId()
			//排队中有满足订单情况
			if cp.Uid != 0 {
				_, flag := cache.GetCacheCreditMessage(cp.Id)
				if flag {
					pp <- utils.ShareDdata{Uid: cp.Uid, Id: cp.Id, Flag: true, TimeDff: 45 * 60}
					credit_queue_change_state(this, cp.Id, cp.Uid)
					return
				}
			}
			//排队中没有满足订单的情况 防止死锁
			if cp.Uid == 0 {
				pp <- utils.ShareDdata{Uid: 0, Flag: false, Id: 0}
				return
			}

			//防止死锁
			if count <= 0 {
				pp <- utils.ShareDdata{Uid: 0, Flag: false, Id: 0}
				return
			}
			count = count - 1
		}
	}()
	return
}

//校验授信有没有正在处理的状态 异常跳走页面
func check_credit_handing(this *WorkspaceController) (id, uid int, timediff float64, flag bool) {
	//	id = cache.GetCacheCreditHandingUids(this.User.Id)
	//	uid, flag = cache.GetCacheCreditMessage(id)
	id, uid = models.SelectSystemIDHanding(this.User.Id)
	if id > 0 && uid > 0 {
		oc := models.QueryCreditAttime(id)
		timediff = oc.AllocationTime.Add(45 * time.Minute).Sub(time.Now()).Seconds()
		if timediff <= 0 {
			credit_time_out_handle(this, id, uid)
			return 0, 0, 0, false
		} else {
			return id, uid, timediff, true
		}
	}
	//	models.CreditHandingLogOut() //分配订单后不做任何处理 45分钟超时
	return 0, 0, 0, false
}

func credit_time_out_handle(this *WorkspaceController, id, uid int) {

	//45分钟之外 超时
	err := models.UpdateCreditQueueing(id)
	if err != nil {
		cache.RecordLogs(this.User.Id, 0, this.User.Name, this.User.DisplayName, "更新credit状态失败", "workspace/creditqueuehanding", err.Error(), this.Ctx.Input)
	}
	err = utils.Rc.Delete("xjfq:" + utils.CacheKeyCreditHandingUids + "_" + strconv.Itoa(this.User.Id))
	if err != nil {
		cache.RecordLogs(this.User.Id, 0, this.User.Name, this.User.DisplayName, "删去redis失败", "workspace/creditqueuehanding", err.Error(), this.Ctx.Input)
	}
	err = utils.Rc.Delete("xjfq:" + utils.CacheKeyCreditMessage + "_" + strconv.Itoa(id))
	if err != nil {
		cache.RecordLogs(this.User.Id, 0, this.User.Name, this.User.DisplayName, "删去redis失败", "workspace/creditqueuehanding", err.Error(), this.Ctx.Input)
	}
	//增加超时记录
	content := "【" + this.User.DisplayName + "】" + "授信处理超时，返回排队中"
	err = models.AddCreditAduitRecord(id, uid, content, "")
	if err != nil {
		cache.RecordLogs(this.User.Id, 0, this.User.Name, this.User.DisplayName, "获取预约数据失败", "workspace/creditqueuehanding", err.Error(), this.Ctx.Input)
	}

}

//授信状态redis改变
func credit_queue_change_state(this *WorkspaceController, id, uid int) {
	//存redis
	if utils.Re == nil {
		err := utils.Rc.Put("xjfq:"+utils.CacheKeyCreditHandingUids+"_"+strconv.Itoa(this.User.Id), id, utils.RedisCacheTime_45MIN)
		if err != nil {
			cache.RecordLogs(this.User.Id, 0, this.User.Name, this.User.DisplayName, "存入redis失败", "workspace/creditqueuehanding", err.Error(), this.Ctx.Input)
			return
		}
		err = utils.Rc.Put("xjfq:"+utils.CacheKeyCreditMessage+"_"+strconv.Itoa(id), uid, utils.RedisCacheTime_45MIN)
		if err != nil {
			cache.RecordLogs(this.User.Id, 0, this.User.Name, this.User.DisplayName, "存入redis失败", "workspace/creditqueuehanding", err.Error(), this.Ctx.Input)
			return
		}
	}
	//分配成功 "handing"
	err := models.UpdateCueditQueueStatusInqueueTime("HANDING", this.User.DisplayName, this.User.Id, id)
	if err != nil {
		cache.RecordLogs(this.User.Id, 0, this.User.Name, this.User.DisplayName, "更新状态失败", "workspace/creditqueuehanding", err.Error(), this.Ctx.Input)
	}

	//增加分配记录
	content := "订单分配给" + "【" + this.User.DisplayName + "】"
	err = models.AddCreditAduitRecord(id, uid, content, "")
	if err != nil {
		cache.RecordLogs(this.User.Id, 0, this.User.Name, this.User.DisplayName, "获取预约数据失败", "workspace/creditqueuehanding", err.Error(), this.Ctx.Input)
	}

	//增加预防超时
	err = models.UpdateCreditAlloctionTime(id)
	if err != nil {
		cache.RecordLogs(this.User.Id, 0, this.User.Name, this.User.DisplayName, "更新分配时间失败", "workspace/creditqueuehanding", err.Error(), this.Ctx.Input)
	}
}

//信审人员处理授信数据(1.通过 2 关闭 3 退回 4驳回)
func (this *WorkspaceController) CreditOp() {
	var err error
	defer func() {
		this.ServeJSON()
	}()
	var flag int
	var resp models.BaseResponse
	var cp models.CreditAduit
	cp.State = this.GetString("state")
	cp.Remark = this.GetString("remark")
	cp.Uid, _ = this.GetInt("uid")
	cp.Id, _ = this.GetInt("id")
	if cp.State == "PASS" {
		cp.BalanceMoney, err = this.GetInt("balancemoney")
		if err != nil {
			this.Data["json"] = map[string]interface{}{"ret": 304, "err": err.Error(), "msg": "获取额度失败"}
			cache.RecordLogs(this.User.Id, 0, this.User.Name, this.User.DisplayName, "获取额度失败", "workspace/creditop", err.Error(), this.Ctx.Input)
			return
		}
	}
	err = updatecreditstatus(cp, this)
	if err != nil {
		cache.RecordLogs(this.User.Id, 0, this.User.Name, this.User.DisplayName, "更新数据库失败", "workspace/OrderOpera", err.Error(), this.Ctx.Input)
		return
	}
	if cp.State == "REJECT" {
		var table string
		var sort_no int
		reject_reason := this.GetString("rejectReason")
		table, flag, sort_no = utils.FlagToTable(this.GetString("rejectType"))
		utils.UpdateRecjectTable(table, reject_reason, sort_no, cp.Uid)
	}

	//delete redis(处理清除缓存)
	if utils.Re == nil {
		//处理授信key
		err := utils.Rc.Delete("xjfq:" + utils.CacheKeyCreditMessage + "_" + strconv.Itoa(cp.Id))
		if err != nil {
			cache.RecordLogs(this.User.Id, 0, this.User.Name, this.User.DisplayName, "delete redis失败", "workspace/creditop", err.Error(), this.Ctx.Input)
			return
		}
		//处理超时key
		err = utils.Rc.Delete("xjfq:" + utils.CacheKeyCreditHandingUids + "_" + strconv.Itoa(this.User.Id))
		if err != nil {
			cache.RecordLogs(this.User.Id, 0, this.User.Name, this.User.DisplayName, "delete redis失败", "workspace/creditop", err.Error(), this.Ctx.Input)
			return
		}
	}
	params := map[string]interface{}{
		"Uid":          cp.Uid,          //用户id
		"Id":           cp.Id,           //主键id
		"BalanceMoney": cp.BalanceMoney, //额度
		"SysId":        this.User.Id,    //系统id
		"State":        cp.State,        //状态
		"flag":         flag,            //个人信息状态修改 1 ，常用联系人状态修改 2 , 不是驳回状态 0
	}

	b, err := services.PostApi("credit/statusop", params)
	err = json.Unmarshal(b, &resp)
	if err != nil {
		resp.Ret = 304
		resp.Msg = "api请求失败"
		resp.Err = err.Error()
		cache.RecordLogs(this.User.Id, 0, this.User.Name, this.User.DisplayName, "通过api���求失败", "workspace/OrderOpera", err.Error(), this.Ctx.Input)
	}
	if cp.State == "PASS" {
		config, err := cache.GetConfigByCKeyCache("invitation_voucher_award")
		if err != nil {
			cache.RecordLogs(this.User.Id, 0, this.User.Name, this.User.DisplayName, "获取邀请好友奖励配置失败", "workspace/creditop", err.Error(), this.Ctx.Input)
		} else {
			money, err := strconv.Atoi(config.ConfigValue)
			if err != nil {
				cache.RecordLogs(this.User.Id, 0, this.User.Name, this.User.DisplayName, "邀请好友奖励金额转换失败", "workspace/creditop", err.Error(), this.Ctx.Input)
			} else {
				err = models.AddInvitationReward(cp.Uid, money, 1)
				if err != nil {
					cache.RecordLogs(this.User.Id, 0, this.User.Name, this.User.DisplayName, "添加用户被邀请认证成功记录失败", "workspace/creditop", err.Error(), this.Ctx.Input)
				}
			}
		}
		//增加通过记录
		content := "授信通过，获得" + strconv.Itoa(cp.BalanceMoney) + "元额度 ——" + this.User.DisplayName
		err = models.AddCreditAduitRecord(cp.Id, cp.Uid, content, cp.Remark)
		if err != nil {
			cache.RecordLogs(this.User.Id, 0, this.User.Name, this.User.DisplayName, "授信通过失败", "workspace/creditop", err.Error(), this.Ctx.Input)
			return
		}
		//增加额度记录
		err = models.AddLimitRecord(cp.Uid, this.User.Id, cp.BalanceMoney, cp.Remark)
		if err != nil {
			cache.RecordLogs(this.User.Id, 0, this.User.Name, this.User.DisplayName, "授信通过插入额度记录失败", "workspace/creditop", err.Error(), this.Ctx.Input)
			return
		}
		this.Data["json"] = map[string]interface{}{"ret": 200, "msg": "授信通过"}
		return
	}

	if cp.State == "CLOSE" {
		//增加关闭记录
		content := "永久关闭 —— " + this.User.DisplayName
		err = models.AddCreditAduitRecord(cp.Id, cp.Uid, content, cp.Remark)
		if err != nil {
			cache.RecordLogs(this.User.Id, 0, this.User.Name, this.User.DisplayName, "获取预约数据失败", "workspace/creditop", err.Error(), this.Ctx.Input)
			return
		}
		this.Data["json"] = map[string]interface{}{"ret": 200, "msg": "关闭成功"}
		return
	}

	if cp.State == "PAUSE" {
		//增加关闭30天记录
		content := "关闭30天 —— " + this.User.DisplayName
		err = models.AddCreditAduitRecord(cp.Id, cp.Uid, content, cp.Remark)
		if err != nil {
			cache.RecordLogs(this.User.Id, 0, this.User.Name, this.User.DisplayName, "关闭30天失败", "workspace/creditop", err.Error(), this.Ctx.Input)
			return
		}
		this.Data["json"] = map[string]interface{}{"ret": 200, "msg": "关闭成功"}
		return
	}

	if cp.State == "OUTQUEUE" {
		//增加授信退回记录
		content := "授信退回 —— " + this.User.DisplayName
		err = models.AddCreditAduitRecord(cp.Id, cp.Uid, content, cp.Remark)
		if err != nil {
			cache.RecordLogs(this.User.Id, 0, this.User.Name, this.User.DisplayName, "授信退回失败", "workspace/creditop", err.Error(), this.Ctx.Input)
			return
		}
		this.Data["json"] = map[string]interface{}{"ret": 200, "msg": "退回成功"}
		return
	}

	if cp.State == "REJECT" {
		cp.Remark = this.GetString("rejectType") + ":" + this.GetString("rejectReason")
		content := "由于" + cp.Remark + "审核未通过，请重新编写后再提交申请."
		params := map[string]interface{}{
			"Uid":     cp.Uid, //用户id
			"Message": content,
			"Title":   "审核未通过通知!",
		}
		b, err := services.PostApi("usermessage/addMessage", params)
		err = json.Unmarshal(b, &resp)
		if err != nil {
			resp.Ret = 304
			resp.Msg = "api请求失败"
			resp.Err = err.Error()
			cache.RecordLogs(this.User.Id, 0, this.User.Name, this.User.DisplayName, "通过api请求失败", "workspace/OrderOpera", err.Error(), this.Ctx.Input)
			return
		}
		//增加驳回记录
		r_content := "授信驳回 —— " + this.User.DisplayName
		err = models.AddCreditAduitRecord(cp.Id, cp.Uid, r_content, cp.Remark)
		if err != nil {
			cache.RecordLogs(this.User.Id, 0, this.User.Name, this.User.DisplayName, "授信驳回失败", "workspace/creditop", err.Error(), this.Ctx.Input)
			return
		}
		this.Data["json"] = map[string]interface{}{"ret": 200, "msg": "驳回成功"}
		return
	}

}

func updatecreditstatus(cp models.CreditAduit, this *WorkspaceController) (err error) {
	if cp.State == "PASS" {
		//更新cuedit_audit
		if err = models.UpdateCueditQueuePassStatus(cp.State, cp.BalanceMoney, this.User.Id, cp.Id); err != nil {
			return err
		}

	} else if cp.State == "OUTQUEUE" {
		//更新cuedit_audit
		InqueueTime := this.GetString("inqueuetime")
		if err = models.UpdateCreditOutqueueTime(cp.Id, this.User.Id, InqueueTime); err != nil {
			return err
		}

	} else {
		//更新cuedit_audit
		if err = models.UpdateCueditQueueStatusOp(cp.State, cp.Id, this.User.Id); err != nil {
			return err
		}
	}

	//驳回
	if cp.State == "REJECT" {
		cp.Remark = this.GetString("rejectType") + ":" + this.GetString("rejectReason")
	}

	//增加授信记录
	if err = models.UpdateCueditQueueRemark(cp.Remark, cp.Id); err != nil {
		return err
	}
	return err
}

//信审人员在退回队列中手动添加排队或者插队时间
func (this *WorkspaceController) UpdateInQueue() {
	defer this.ServeJSON()
	var iq models.CreditQueueLine
	json.Unmarshal(this.Ctx.Input.RequestBody, &iq)
	for _, v := range iq.Ids {
		is_znoe := models.GetCreditIsCut(v)
		//不拆分
		if is_znoe == 0 {
			err := models.UpadatecCredintInQueue(iq.InqueueTime, iq.InqueueType, v)
			if err != nil {
				this.Data["json"] = map[string]interface{}{"ret": 304, "err": err.Error(), "msg": "更新预约状态和时间失败"}
				cache.RecordLogs(this.User.Id, 0, this.User.Name, this.User.DisplayName, "更新预约状态和时间失败", "workspace/UpdateInQueue", err.Error(), this.Ctx.Input)
				return
			}
		} else {
			//拆分
			err := models.UpadatecCutCredintInQueue(iq.InqueueTime, iq.InqueueType, v)
			if err != nil {
				this.Data["json"] = map[string]interface{}{"ret": 304, "err": err.Error(), "msg": "更新预约状态和时间失败"}
				cache.RecordLogs(this.User.Id, 0, this.User.Name, this.User.DisplayName, "更新预约状态和时间失败", "workspace/UpdateInQueue", err.Error(), this.Ctx.Input)
				return
			}
		}
	}
	this.Data["json"] = map[string]interface{}{"ret": 200, "msg": "入列成功"}
	return
}

//授信查看
func (this *WorkspaceController) CreditLook() {
	this.IsNeedTemplate()
	uid, _ := this.GetInt("uid")

	if uid <= 0 {
		cache.RecordLogs(this.User.Id, 0, this.User.Name, this.User.DisplayName, "uid参数���递错误", "授信查看/CreditLook", "", this.Ctx.Input)
		this.Abort("uid参数传递错误")
		return
	}
	creditAduitId, _ := this.GetInt("credit_aduit_id")
	if creditAduitId <= 0 {
		cache.RecordLogs(this.User.Id, 0, this.User.Name, this.User.DisplayName, "credit_aduit_id参数传递错误", "授信查看/CreditLook", "", this.Ctx.Input)
		this.Abort("credit_aduit_id参数传递错误")
		return
	}
	user, err := models.QueryUsersBaseInfo(uid) //用户基本信息
	if err != nil && err.Error() != utils.ErrNoRow() {
		cache.RecordLogs(this.User.Id, 0, this.User.Name, this.User.DisplayName, "查询用户基本信息出错", "授信查看/CreditLook", err.Error(), this.Ctx.Input)
		this.Abort("查询用户基本信息出错" + err.Error())
		return
	}
	if user == nil {
		this.Abort("用户信息不存在")
		return
	}
	creditInfo, err := models.GetCreditInfoByUid(creditAduitId)
	if err != nil && err.Error() != utils.ErrNoRow() {
		cache.RecordLogs(this.User.Id, 0, this.User.Name, this.User.DisplayName, "查询风控信息出错", "授信查看/CreditLook", err.Error(), this.Ctx.Input)
		this.Abort("查询风控信息出错" + err.Error())
		return
	}

	examineInfo, err := models.QueryCreditAduitRecord(uid)
	if err != nil && err.Error() != utils.ErrNoRow() {
		cache.RecordLogs(this.User.Id, 0, this.User.Name, this.User.DisplayName, "查询授信审批历史出错", "信审工作平台/GetCreditAduitList", err.Error(), this.Ctx.Input)
		this.Abort("查询授信审批历史出错" + err.Error())
		return
	}

	//拆分审批意见
	advise, err := models.GetCreditAdvise(creditAduitId)
	services.AdviseSort(advise)
	if err != nil {
		this.Data["json"] = map[string]interface{}{"ret": 304, "err": err.Error(), "msg": "id获取失败~"}
		cache.RecordLogs(this.User.Id, 0, this.User.Name, this.User.DisplayName, "uid获取失败", "workspace/creditop", err.Error(), this.Ctx.Input)
		return
	}

	this.Data["user"] = user
	this.Data["creditInfo"] = creditInfo
	this.Data["examineInfo"] = examineInfo
	this.Data["adviselist"] = advise
	this.TplName = "workspace/workspace_credit_look.html"
}

//案件分配比例管理
func (c *WorkspaceController) DistributeRatioManage() {
	c.IsNeedTemplate()
	list, err := models.QueryCreditRatios()
	if err != nil {
		cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "获取方案比例失败", "方案分配比例管理/DistributeRatioManage", err.Error(), c.Ctx.Input)
		c.Abort("获取方案比例失败")
	}
	c.Data["list"] = list
	c.TplName = "workspace/workspace_distribute_ratio.html"
}

//修改案件分配比例
func (c *WorkspaceController) MotifyCreditRatio() {
	resultMap := make(map[string]interface{})
	resultMap["ret"] = 304
	defer func() {
		c.Data["json"] = resultMap
		c.ServeJSON()
	}()
	partWeight, _ := c.GetInt("part_weight", 0) //拆分分配权重
	allWeight, _ := c.GetInt("all_weight", 0)   //完整分配权重
	if partWeight != 0 && allWeight != 0 {
		comNumber := utils.MaxCommonDivisor(partWeight, allWeight)
		partWeight /= comNumber
		allWeight /= comNumber
	}
	if partWeight == 0 && allWeight == 0 {
		partWeight = 1

		allWeight = 1
	}
	remark := c.GetString("remark") //修改说明
	operator := c.User.DisplayName
	err := models.AddCreditRatio(partWeight, allWeight, operator, remark)
	if err != nil {
		resultMap["err"] = "方案比例插入数据库失败" + err.Error()
		cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "方案比例插入数据库失败", "方案分配比例管理/MotifyCreditRatio", err.Error(), c.Ctx.Input)
		return
	}
	//重新设置比例将计数器重置
	if utils.Re == nil && utils.Rc.IsExist(utils.CacheKeyModeCount) {
		err := utils.Rc.Delete(utils.CacheKeyModeCount)
		if err != nil {
			beego.Info("delete cache mode count fail:", err)
		}
	}
	resultMap["ret"] = 200
	resultMap["msg"] = "方案比例设置成功"
}

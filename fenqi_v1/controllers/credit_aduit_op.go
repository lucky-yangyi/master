package controllers

import (
	"encoding/json"
	"fenqi_v1/cache"
	"fenqi_v1/models"
	"fenqi_v1/services"
	"fenqi_v1/utils"
	"strconv"
	"strings"
	"time"
)

//信审人员处理授信数据(1.通过 2 关闭 3 退回 4驳回)
func (this *WorkspaceController) CreditCutOp() {
	defer func() {
		this.ServeJSON()
	}()
	var cp models.CreditAduit
	cp.State = this.GetString("state")
	cp.Remark = this.GetString("remark")
	check, err := this.GetInt("check")
	cp.Uid, err = this.GetInt("uid")
	cp.Id, err = this.GetInt("id")
	if err != nil {
		this.Data["json"] = map[string]interface{}{"ret": 304, "err": err.Error(), "msg": "id获取失败~"}
		cache.RecordLogs(this.User.Id, 0, this.User.Name, this.User.DisplayName, "uid获取失败", "workspace/creditop", err.Error(), this.Ctx.Input)
		return
	}

	if cp.State == "PASS" {
		cp.BalanceMoney, err = this.GetInt("balancemoney")
		if err != nil {
			this.Data["json"] = map[string]interface{}{"ret": 304, "err": err.Error(), "msg": "获取额度失败"}
			cache.RecordLogs(this.User.Id, 0, this.User.Name, this.User.DisplayName, "获取额度失败", "workspace/creditop", err.Error(), this.Ctx.Input)
			return
		}
	}
	//更新本人
	if check == 1 {
		err = UpdateMyInfoCreditStatus(cp, this)
		if err != nil {
			cache.RecordLogs(this.User.Id, 0, this.User.Name, this.User.DisplayName, "更新数据库失败", "workspace/OrderOpera", err.Error(), this.Ctx.Input)
			return
		}
	}
	//更新其他
	if check == 2 {
		err = UpdateOtherCreditStatus(cp, this)
		if err != nil {
			cache.RecordLogs(this.User.Id, 0, this.User.Name, this.User.DisplayName, "更新数据库失败", "workspace/OrderOpera", err.Error(), this.Ctx.Input)
			return
		}
	}
	//更新联系人
	if check == 3 {
		err = UpdateLinkManCreditStatus(cp, this)
		if err != nil {
			cache.RecordLogs(this.User.Id, 0, this.User.Name, this.User.DisplayName, "更新数据库失败", "workspace/OrderOpera", err.Error(), this.Ctx.Input)
			return
		}
	}
	//更新授权
	if check == 4 {
		err = UpdateAuthCreditStatus(cp, this)
		if err != nil {
			cache.RecordLogs(this.User.Id, 0, this.User.Name, this.User.DisplayName, "更新数据库失败", "workspace/OrderOpera", err.Error(), this.Ctx.Input)
			return
		}
	}

	//全部完成
	is_ok := models.GetCreditIsOk(cp.Id)
	if is_ok == 4 {
		go cut_op_ok(this, cp)
	}

	//未全部完成
	if cp.State == "PASS" {
		this.Data["json"] = map[string]interface{}{"ret": 200, "msg": "授信通过"}
		return
	}
	if cp.State == "CLOSE" {
		this.Data["json"] = map[string]interface{}{"ret": 200, "msg": "关闭成功"}
		return
	}
	if cp.State == "PAUSE" {
		this.Data["json"] = map[string]interface{}{"ret": 200, "msg": "关闭成功"}
		return
	}
	if cp.State == "REJECT" {
		this.Data["json"] = map[string]interface{}{"ret": 200, "msg": "驳回成功"}
		return
	}
	if cp.State == "OUTQUEUE" {
		this.Data["json"] = map[string]interface{}{"ret": 200, "msg": "退回成功"}
		return
	}

}

func cut_op_ok(this *WorkspaceController, cp models.CreditAduit) {
	var resp models.BaseResponse
	var handling_time, at_time time.Time
	var flag, sort_no int
	state, diaplayname := models.GetCreditIsCutStatus(cp.Id)
	if state != "" {
		cp.State = models.GetCreditCutState(state)
		if cp.State == "PASS" {
			if strings.Count(state, "PASS") == 4 {
				cp.BalanceMoney = models.GetCreditCutPassMoney(cp.Id)
			}
		}
	}

	handling_time, at_time, this.User.Id, this.User.DisplayName, cp.Remark = models.HandlingTimeOp(cp.State, cp.Uid, cp.Id)
	if cp.State == "REJECT" {
		var table string
		table, flag, sort_no = utils.FlagToTable(strings.Split(cp.Remark, ":")[0])
		utils.UpdateRecjectTable(table, strings.Split(cp.Remark, ":")[1], sort_no, cp.Uid)
	}

	params := map[string]interface{}{
		"Uid":          cp.Uid,          //用户id
		"Id":           cp.Id,           //主键Id
		"BalanceMoney": cp.BalanceMoney, //额度
		"SysId":        this.User.Id,    //系统id
		"State":        cp.State,        //状态
		"Flag":         flag,            //个人信息状态修改 1 ，常用联系人状态修改 -1 , 不是驳回状态 0
	}

	b, err := services.PostApi("credit/statusop", params)
	err = json.Unmarshal(b, &resp)
	if err != nil {
		resp.Ret = 304
		resp.Msg = "api请求失败"
		resp.Err = err.Error()
		cache.RecordLogs(this.User.Id, 0, this.User.Name, this.User.DisplayName, "通过api请求失败", "workspace/creditcutop", err.Error(), this.Ctx.Input)
		return
	}

	err = UpdateCutCreditStatus(cp, this, diaplayname, handling_time, at_time)
	if err != nil {
		cache.RecordLogs(this.User.Id, 0, this.User.Name, this.User.DisplayName, "更新数据库失败", "workspace/creditcutop", err.Error(), this.Ctx.Input)
		return
	}

	//全部分完记录
	cut_content := "订单分配给" + "【" + diaplayname + "】"
	err = models.AddCreditAduitRecord(cp.Id, cp.Uid, cut_content, "")
	if err != nil {
		cache.RecordLogs(this.User.Id, 0, this.User.Name, this.User.DisplayName, "授信通过失败", "workspace/creditcutop", err.Error(), this.Ctx.Input)
		return
	}

	if cp.State == "PASS" {
		config, err := cache.GetConfigByCKeyCache("invitation_voucher_award")
		if err != nil {
			cache.RecordLogs(this.User.Id, 0, this.User.Name, this.User.DisplayName, "获取邀请好友奖励配置失败", "workspace/creditcutop", err.Error(), this.Ctx.Input)
		} else {
			money, err := strconv.Atoi(config.ConfigValue)
			if err != nil {
				cache.RecordLogs(this.User.Id, 0, this.User.Name, this.User.DisplayName, "邀请好友奖励金额转换失败", "workspace/creditcutop", err.Error(), this.Ctx.Input)
			} else {
				err = models.AddInvitationReward(cp.Uid, money, 1)
				if err != nil {
					cache.RecordLogs(this.User.Id, 0, this.User.Name, this.User.DisplayName, "添加用户被邀请认证成功记录失败", "workspace/creditcutop", err.Error(), this.Ctx.Input)
				}
			}
		}

		//增加通过记录
		content := "授信通过，获得" + strconv.Itoa(cp.BalanceMoney) + "元额度"
		err = models.AddCreditAduitRecord(cp.Id, cp.Uid, content, cp.Remark)
		if err != nil {
			cache.RecordLogs(this.User.Id, 0, this.User.Name, this.User.DisplayName, "授信通过失败", "workspace/creditcutop", err.Error(), this.Ctx.Input)
			return
		}
		//增加额度记录
		err = models.AddLimitRecord(cp.Uid, this.User.Id, cp.BalanceMoney, cp.Remark)
		if err != nil {
			cache.RecordLogs(this.User.Id, 0, this.User.Name, this.User.DisplayName, "授信通过插入额度记录失败", "workspace/creditcutop", err.Error(), this.Ctx.Input)
			return
		}
		//发送系统消息
		umessage := &models.UserMessage{Title: "授信成功", Content: "授信评估通过，获得" + strconv.Itoa(cp.BalanceMoney) + "元额度，借款最快2分钟放款哦", Uid: cp.Uid}
		models.AddUsersMessage(umessage)

		user, err := models.QueryUserByUid(cp.Uid)
		if err != nil {
			cache.RecordLogs(this.User.Id, 0, this.User.Name, this.User.DisplayName, "获取用户信息失败", "workspace/creditcutop", err.Error(), this.Ctx.Input)
		}
		content = "恭喜您，授信评估通过，获得" + strconv.Itoa(cp.BalanceMoney) + "元额度，快来借款吧"
		//发送短信信息
		services.SendSms(user.Account, content, "fenqi_v1", "授信成功", this.Ctx.Input.IP(), "0")
		return
	}

	if cp.State == "CLOSE" {
		//增加关闭记录
		content := "永久关闭"
		err = models.AddCreditAduitRecord(cp.Id, cp.Uid, content, cp.Remark)
		if err != nil {
			cache.RecordLogs(this.User.Id, 0, this.User.Name, this.User.DisplayName, "获取预约数据失败", "workspace/creditcutop", err.Error(), this.Ctx.Input)
			return
		}
		return
	}

	if cp.State == "PAUSE" {
		//增加关闭30天记录
		content := "关闭30天"
		err = models.AddCreditAduitRecord(cp.Id, cp.Uid, content, cp.Remark)
		if err != nil {
			cache.RecordLogs(this.User.Id, 0, this.User.Name, this.User.DisplayName, "关闭30天失败", "workspace/creditcutop", err.Error(), this.Ctx.Input)
			return
		}
		return
	}

	if cp.State == "OUTQUEUE" {
		//增加授信退回记录
		content := "授信退回"
		err = models.AddCreditAduitRecord(cp.Id, cp.Uid, content, cp.Remark)
		if err != nil {
			cache.RecordLogs(this.User.Id, 0, this.User.Name, this.User.DisplayName, "授信退回失败", "workspace/creditcutop", err.Error(), this.Ctx.Input)
			return
		}
		return
	}

	if cp.State == "REJECT" {
		//增加驳回记录
		content := "由于" + cp.Remark + "审核未通过，请重新编写后再提交申请."
		params := map[string]interface{}{
			"Uid":     cp.Uid, //用户id
			"Message": content,
			"Title":   "审核未通过通知！",
		}
		b, err := services.PostApi("usermessage/addMessage", params)
		err = json.Unmarshal(b, &resp)
		if err != nil {
			resp.Ret = 304
			resp.Msg = "api请求失败"
			resp.Err = err.Error()
			cache.RecordLogs(this.User.Id, 0, this.User.Name, this.User.DisplayName, "通过api请求失败", "workspace/creditcutop", err.Error(), this.Ctx.Input)
			return
		}
		r_content := "授信驳回"
		err = models.AddCreditAduitRecord(cp.Id, cp.Uid, r_content, cp.Remark)
		if err != nil {
			cache.RecordLogs(this.User.Id, 0, this.User.Name, this.User.DisplayName, "授信驳回失败", "workspace/creditcutop", err.Error(), this.Ctx.Input)
			return
		}

	}

}

func UpdateMyInfoCreditStatus(cp models.CreditAduit, this *WorkspaceController) (err error) {

	//驳回
	if cp.State == "REJECT" {
		cp.Remark = this.GetString("rejectType") + ":" + this.GetString("rejectReason")
	}

	if cp.State == "PASS" {
		//更新cuedit_audit
		if err = models.UpdateMyInfoCueditQueuePassStatus(cp.State, cp.Remark, cp.BalanceMoney, this.User.Id, cp.Id); err != nil {
			return err
		}
	} else {
		//更新cuedit_audit
		if err = models.UpdateMyInfoCueditQueueStatusOp(cp.State, cp.Remark, cp.Id, this.User.Id); err != nil {
			return err
		}
	}

	//增加授信记录
	if err = models.InsetMyInfoCueditAdviseRemark(cp.Remark, cp.State, cp.Id); err != nil {
		return err
	}

	//更新主表是否完成
	if err = models.UpdateCreditIsMyInfo(cp.Id); err != nil {
		return err
	}

	//清空缓存
	if utils.Re == nil {
		//处理授信key
		err := utils.Rc.Delete("myinfo:" + utils.CacheKeyCreditMessage + "_" + strconv.Itoa(cp.Id))
		if err != nil {
			cache.RecordLogs(this.User.Id, 0, this.User.Name, this.User.DisplayName, "delete redis失败", "workspace/creditcutop", err.Error(), this.Ctx.Input)
			return err
		}
		//处理超时key
		err = utils.Rc.Delete("myinfo:" + utils.CacheKeyCreditHandingUids + "_" + strconv.Itoa(this.User.Id))
		if err != nil {
			cache.RecordLogs(this.User.Id, 0, this.User.Name, this.User.DisplayName, "delete redis失败", "workspace/creditcutop", err.Error(), this.Ctx.Input)
			return err
		}
	}
	return
}

func UpdateOtherCreditStatus(cp models.CreditAduit, this *WorkspaceController) (err error) {
	//驳回
	if cp.State == "REJECT" {
		cp.Remark = this.GetString("rejectType") + ":" + this.GetString("rejectReason")
	}

	if cp.State == "PASS" {
		//更新cuedit_audit
		if err = models.UpdateOhterCueditAuthQueuePassStatus(cp.State, cp.Remark, cp.BalanceMoney, this.User.Id, cp.Id); err != nil {
			return err
		}
	} else {
		//更新cuedit_audit
		if err = models.UpdateOtherCueditQueueStatusOp(cp.State, cp.Remark, cp.Id, this.User.Id); err != nil {
			return err
		}
	}

	//更新主表是否完成
	if err = models.UpdateCreditIsOther(cp.Id); err != nil {
		return err
	}
	//增加授信记录
	if err = models.InsetOtherCueditAdviseRemark(cp.Remark, cp.State, cp.Id); err != nil {
		return err
	}
	//清空缓存
	if utils.Re == nil {
		//处理授信key
		err := utils.Rc.Delete("other:" + utils.CacheKeyCreditMessage + "_" + strconv.Itoa(cp.Id))
		if err != nil {
			cache.RecordLogs(this.User.Id, 0, this.User.Name, this.User.DisplayName, "delete redis失败", "workspace/creditcutop", err.Error(), this.Ctx.Input)
			return err
		}
		//处理超时key
		err = utils.Rc.Delete("other:" + utils.CacheKeyCreditHandingUids + "_" + strconv.Itoa(this.User.Id))
		if err != nil {
			cache.RecordLogs(this.User.Id, 0, this.User.Name, this.User.DisplayName, "delete redis失败", "workspace/creditcutop", err.Error(), this.Ctx.Input)
			return err
		}
	}
	return
}

func UpdateLinkManCreditStatus(cp models.CreditAduit, this *WorkspaceController) (err error) {

	//驳回
	if cp.State == "REJECT" {
		cp.Remark = this.GetString("rejectType") + ":" + this.GetString("rejectReason")
	}

	if cp.State == "PASS" {
		//更新cuedit_audit
		if err = models.UpdateLinkManCueditQueuePassStatus(cp.State, cp.Remark, cp.BalanceMoney, this.User.Id, cp.Id); err != nil {
			return err
		}

	} else {
		//更新cuedit_audit
		if err = models.UpdateLinkManCueditQueueStatusOp(cp.State, cp.Remark, cp.Id, this.User.Id); err != nil {
			return err
		}
	}

	//增加授信记录
	if err = models.InsetLinkManCueditAdviseRemark(cp.Remark, cp.State, cp.Id); err != nil {
		return err
	}

	//更新主表是否完成
	if err = models.UpdateCreditIsLinkMan(cp.Id); err != nil {
		return err
	}

	//清空缓存
	if utils.Re == nil {
		//处理授信key
		err := utils.Rc.Delete("linkman:" + utils.CacheKeyCreditMessage + "_" + strconv.Itoa(cp.Id))
		if err != nil {
			cache.RecordLogs(this.User.Id, 0, this.User.Name, this.User.DisplayName, "delete redis失败", "workspace/creditcutop", err.Error(), this.Ctx.Input)
			return err
		}
		//处理超时key
		err = utils.Rc.Delete("linkman:" + utils.CacheKeyCreditHandingUids + "_" + strconv.Itoa(this.User.Id))
		if err != nil {
			cache.RecordLogs(this.User.Id, 0, this.User.Name, this.User.DisplayName, "delete redis失败", "workspace/creditcutop", err.Error(), this.Ctx.Input)
			return err
		}
	}
	return
}

func UpdateAuthCreditStatus(cp models.CreditAduit, this *WorkspaceController) (err error) {

	//驳回
	if cp.State == "REJECT" {
		cp.Remark = this.GetString("rejectType") + ":" + this.GetString("rejectReason")
	}

	if cp.State == "PASS" {
		//更新cuedit_audit
		if err = models.UpdateAuthCueditQueuePassStatus(cp.State, cp.Remark, cp.BalanceMoney, this.User.Id, cp.Id); err != nil {
			return err
		}

	} else {
		//更新cuedit_audit
		if err = models.UpdateAuthCueditQueueStatusOp(cp.State, cp.Remark, cp.Id, this.User.Id); err != nil {
			return err
		}
	}

	//增加授信记录
	if err = models.InsetAuthCueditAdviseRemark(cp.Remark, cp.State, cp.Id); err != nil {
		return err
	}

	//更新授信主表
	if err = models.UpdateCreditIsAuth(cp.Id); err != nil {
		return err
	}

	//清空缓存
	if utils.Re == nil {
		//处理授信key
		err := utils.Rc.Delete("auth:" + utils.CacheKeyCreditMessage + "_" + strconv.Itoa(cp.Id))
		if err != nil {
			cache.RecordLogs(this.User.Id, 0, this.User.Name, this.User.DisplayName, "delete redis失败", "workspace/creditcutop", err.Error(), this.Ctx.Input)
			return err
		}
		//处理超时key
		err = utils.Rc.Delete("auth:" + utils.CacheKeyCreditHandingUids + "_" + strconv.Itoa(this.User.Id))
		if err != nil {
			cache.RecordLogs(this.User.Id, 0, this.User.Name, this.User.DisplayName, "delete redis失败", "workspace/creditcutop", err.Error(), this.Ctx.Input)
			return err
		}
	}
	return
}

func UpdateCutCreditStatus(cp models.CreditAduit, this *WorkspaceController, diaplayname string, at_time, handling_time time.Time) (err error) {

	if cp.State == "PASS" {
		//更新cuedit_audit
		if err = models.UpdateCutCueditQueuePassStatus(cp.State, diaplayname, cp.BalanceMoney, this.User.Id, cp.Id, handling_time); err != nil {
			return err
		}
	} else {
		//更新cuedit_audit
		if err = models.UpdateCutCueditQueueStatusOp(cp.State, diaplayname, cp.Id, this.User.Id, handling_time); err != nil {
			return err
		}
	}

	//增加授信记录
	if err = models.UpdateCueditQueueRemark(cp.Remark, cp.Id); err != nil {
		return err
	}

	//更新分配时间
	if err = models.UpdateCreditAttime(cp.Id, at_time); err != nil {
		return err
	}
	return
}

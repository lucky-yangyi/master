package controllers

import (
	"fenqi_v1/cache"
	"fenqi_v1/models"
	"fenqi_v1/utils"
	"strconv"
	"time"
	//	"fmt"
	"github.com/astaxie/beego"
)
type Workspace struct {
	beego.Controller
}

//处理授信
func (this *WorkspaceController) MetadataQueue() {
	isPermission := models.IsDataPermissionByStationId(this.User.StationId, 4)
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
			this.Data["timeDiff"] = okuser.TimeDff
			this.Data["ok_order"] = true
			this.Data["baseinfo"] = baseinfo
			this.Data["id"] = okuser.Id
			this.Data["uid"] = okuser.Uid
			this.Data["record"] = record
			this.Data["loginstate"] = this.User.LoginState
			this.TplName = "workspace/credit_myinfo.html"
			return
		}
	}()

	//处理异常关闭
	okuser.Id, okuser.Uid, okuser.TimeDff, okuser.Flag = check_myinfo_credit_handing(this)
	//45分钟之内 handing状态
	if okuser.Flag {
		return
	}
	//最后一单 离线分配订单
	if this.User.LoginState == "OFFLINE" && !okuser.Flag {
		this.Data["loginstate"] = this.User.LoginState
		this.Data["ok_order"] = false
		cache.RecordLogs(this.User.Id, 0, this.User.Name, this.User.DisplayName, "没有订单", "workspace/metadataqueue", "", this.Ctx.Input)
		this.TplName = "workspace/credit_myinfo.html"
		return
	}

	//信号量 共享值 正常关闭
	pt := make(chan utils.ShareDdata, 1)
	pp := make(chan utils.ShareDdata, 1)

	//处理插队数据
	var pt_tmp utils.ShareDdata
	var pp_tmp utils.ShareDdata
	go credit_myinfo_queue_a(pt, this)
	pt_tmp = <-pt
	if !pt_tmp.Flag {
		go credit_myinfo_queue_p(pp, this)
		pp_tmp = <-pp
	}
	okuser = utils.ShareFlag(pt_tmp, pp_tmp)
	//没有订单
	if !okuser.Flag {
		this.Data["ok_order"] = false
		this.Data["loginstate"] = this.User.LoginState
		cache.RecordLogs(this.User.Id, 0, this.User.Name, this.User.DisplayName, "没有订单", "workspace/metadataqueue", "", this.Ctx.Input)
		this.TplName = "workspace/credit_myinfo.html"
		return
	}

}

func credit_myinfo_queue_a(pt chan utils.ShareDdata, this *WorkspaceController) {
	//统计预约插队数量
	count := models.GetMyInfoNumCreditQueueAp()
	go func() {
		for {
			ca, _ := models.GetMyInfoAppointmentQueueId()
			if ca.InqueueType == 1 {
				//lock false 为lock状态 ture 为unlock状态
				_, flag := cache.GetMyInfoCacheCreditMessage(ca.CreditAduitId)
				if flag {
					//通道放入值
					pt <- utils.ShareDdata{Id: ca.CreditAduitId, Uid: ca.Uid, Flag: true, TimeDff: 45 * 60}
					credit_myinfo_queue_change_state(this, ca.CreditAduitId, ca.Uid)
					return
				}
			}
			//退回队列中有排队 放入队列中
			if ca.InqueueType == 2 {
				err := models.UpadateMyInfoInqueueType(ca.Id, ca.InqueueTime)
				if err != nil {
					cache.RecordLogs(this.User.Id, 0, this.User.Name, this.User.DisplayName, "更新状态失败", "workspace/metadataqueue", err.Error(), this.Ctx.Input)
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

func credit_myinfo_queue_p(pp chan utils.ShareDdata, this *WorkspaceController) {
	//后期考虑使用 select case 加上timeout机制 统计排队中数量
	count := models.GetMyInfoCreditQueueUpIdCount()
	go func() {
		for {
			cp, _ := models.GetMyIinfoCreditQueueUpId()

			//排队中有满足订单情况
			if cp.Uid != 0 {
				_, flag := cache.GetMyInfoCacheCreditMessage(cp.CreditAduitId)
				if flag {
					pp <- utils.ShareDdata{Uid: cp.Uid, Id: cp.CreditAduitId, Flag: true, TimeDff: 45 * 60}
					credit_myinfo_queue_change_state(this, cp.CreditAduitId, cp.Uid)
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
func check_myinfo_credit_handing(this *WorkspaceController) (id, uid int, timediff float64, flag bool) {
	models.CreditMyinfoHandingLogOut()
	id = cache.GetMyInfoCacheCreditHandingUids(this.User.Id)
	uid, flag = cache.GetMyInfoCacheCreditMessage(id)
	if !flag {
		oc := models.QueryMyInfoCreditAttime(id)
		timediff := oc.AllocationTime.Add(45 * time.Minute).Sub(time.Now()).Seconds()
		if timediff <= 0 {
			credit_myinfo_time_out_handle(this, id, uid)
			return 0, 0, 0, false
		} else {
			return id, uid, timediff, true
		}
	}
	return 0, 0, 0, false
}

func credit_myinfo_time_out_handle(this *WorkspaceController, credit_aduit_id, uid int) {

	//45分钟之外 超时
	err := models.UpdatMyInfoeCreditQueueing(credit_aduit_id)
	if err != nil {
		cache.RecordLogs(this.User.Id, 0, this.User.Name, this.User.DisplayName, "更新credit状态失败", "workspace/metadataqueue", err.Error(), this.Ctx.Input)
	}
	err = utils.Rc.Delete("myinfo:" + utils.CacheKeyCreditHandingUids + "_" + strconv.Itoa(this.User.Id))
	if err != nil {
		cache.RecordLogs(this.User.Id, 0, this.User.Name, this.User.DisplayName, "删去redis失败", "workspace/metadataqueue", err.Error(), this.Ctx.Input)
	}
	err = utils.Rc.Delete("myinfo:" + utils.CacheKeyCreditMessage + "_" + strconv.Itoa(credit_aduit_id))
	if err != nil {

		cache.RecordLogs(this.User.Id, 0, this.User.Name, this.User.DisplayName, "删去redis失败", "workspace/metadataqueue", err.Error(), this.Ctx.Input)
	}
	//增加超时记录
	content := "【" + this.User.DisplayName + "】" + "授信处理超时，返回排队中"
	err = models.AddCreditAduitRecord(credit_aduit_id, uid, content, "")
	if err != nil {
		cache.RecordLogs(this.User.Id, 0, this.User.Name, this.User.DisplayName, "获取预约数据失败", "workspace/creditqueuehanding", err.Error(), this.Ctx.Input)
	}

}

//授信状态redis改变
func credit_myinfo_queue_change_state(this *WorkspaceController, credit_aduit_id, uid int) {
	//存redis
	if utils.Re == nil {
		err := utils.Rc.Put("myinfo:"+utils.CacheKeyCreditHandingUids+"_"+strconv.Itoa(this.User.Id), credit_aduit_id, utils.RedisCacheTime_Year)
		if err != nil {
			cache.RecordLogs(this.User.Id, 0, this.User.Name, this.User.DisplayName, "存入redis失败", "workspace/metadataqueue", err.Error(), this.Ctx.Input)
			return
		}
		err = utils.Rc.Put("myinfo:"+utils.CacheKeyCreditMessage+"_"+strconv.Itoa(credit_aduit_id), uid, utils.RedisCacheTime_Year)
		if err != nil {
			cache.RecordLogs(this.User.Id, 0, this.User.Name, this.User.DisplayName, "存入redis失败", "workspace/metadataqueue", err.Error(), this.Ctx.Input)
			return
		}
	}
	//分配成功 "handing"
	err := models.UpdateMyInfoCueditQueueStatusInqueueTime("HANDING", this.User.DisplayName, this.User.Id, credit_aduit_id)
	if err != nil {
		cache.RecordLogs(this.User.Id, 0, this.User.Name, this.User.DisplayName, "更新状态失败", "workspace/metadataqueue", err.Error(), this.Ctx.Input)
	}
	//更新主表状态
	displayname := models.GetDisplayNameStatus(credit_aduit_id)
	err = models.UpdateCueditQueueStatusInqueueTime("HANDING", displayname, this.User.Id, credit_aduit_id)
	if err != nil {
		cache.RecordLogs(this.User.Id, 0, this.User.Name, this.User.DisplayName, "更新状态失败", "workspace/metadataqueue", err.Error(), this.Ctx.Input)
	}
	//增加分配记录
	content := "订单分配给" + "【" + this.User.DisplayName + "】"
	err = models.AddCreditAduitRecord(credit_aduit_id, uid, content, "")
	if err != nil {
		cache.RecordLogs(this.User.Id, 0, this.User.Name, this.User.DisplayName, "获取预约数据失败", "workspace/creditqueuehanding", err.Error(), this.Ctx.Input)
	}

	//增加预防超时
	err = models.UpdateMyInfoCreditAlloctionTime(credit_aduit_id)
	if err != nil {
		cache.RecordLogs(this.User.Id, 0, this.User.Name, this.User.DisplayName, "更新分配时间失败", "workspace/metadataqueue", err.Error(), this.Ctx.Input)
	}
}

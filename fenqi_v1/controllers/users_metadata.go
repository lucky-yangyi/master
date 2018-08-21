package controllers

import (
	"encoding/json"
	"fenqi_v1/cache"
	"fenqi_v1/models"
	"fenqi_v1/services"
	"fenqi_v1/utils"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/astaxie/beego"
	"gopkg.in/mgo.v2/bson"
)

//用户信息
type UsersMetadataController struct {
	BaseController
}

//获取用户基本信息
func (this *UsersMetadataController) GetUsersBaseInfo() {
	uid, _ := this.GetInt("uid")
	resultMap := make(map[string]interface{})
	resultMap["ret"] = 403
	defer func() {
		this.Data["json"] = resultMap
		this.ServeJSON()
	}()
	if uid <= 0 {
		cache.RecordLogs(this.User.Id, 0, this.User.Name, this.User.DisplayName, "uid参数传递错误", "信审工作平台/GetUsersBaseInfo", "", this.Ctx.Input)
		resultMap["err"] = "uid参数传递错误"
		return
	}
	user, err := models.QueryUsersBaseInfo(uid)
	if err != nil && err.Error() != utils.ErrNoRow() {
		cache.RecordLogs(this.User.Id, 0, this.User.Name, this.User.DisplayName, "查询用户基本信息出错", "信审工作平台/GetUsersBaseInfo", err.Error(), this.Ctx.Input)
		resultMap["err"] = "查询用户基本信息出错"
		return
	}
	resultMap["ret"] = 200
	resultMap["data"] = user
}

//更新用户基本信息
func (this *UsersMetadataController) UpdateUsersBaseInfo() {
	ubiId, _ := this.GetInt("ubi_id")
	uid, _ := this.GetInt("uid")
	beego.Info(uid)
	id, err := this.GetInt("id")
	beego.Info(id, err)
	resultMap := make(map[string]interface{})
	resultMap["ret"] = 403
	defer func() {
		this.Data["json"] = resultMap
		this.ServeJSON()
	}()
	if ubiId <= 0 {
		cache.RecordLogs(this.User.Id, 0, this.User.Name, this.User.DisplayName, "ubi_id参数传递错误", "信审工作平台/UpdateUsersBaseInfo", "", this.Ctx.Input)
		resultMap["err"] = "ubi_id参数传递错误"
		return
	}
	if flag := cache.GetCacheRecheckOften(uid); flag {
		cache.RecordLogs(this.User.Id, 0, this.User.Name, this.User.DisplayName, "请不要重复重检", "信审工作平台/UpdateUsersBaseInfo", "", this.Ctx.Input)
		resultMap["err"] = "请不要重复重检"
		return
	}
	education := this.GetString("education")                   //学历
	marriage := this.GetString("marriage")                     //婚姻
	province := this.GetString("province")                     //省
	city := this.GetString("city")                             //市
	district := this.GetString("district")                     //县或区
	liveDetailAddress := this.GetString("live_detail_address") //详细地址

	err = models.UpdateUsersBaseInfo(ubiId, education, marriage, province, city, district, liveDetailAddress)
	if err != nil {
		cache.RecordLogs(this.User.Id, 0, this.User.Name, this.User.DisplayName, "更新用户基本信息出错", "信审工作平台/UpdateUsersBaseInfo", err.Error(), this.Ctx.Input)
		resultMap["err"] = "更新用户基本信息出错"
		return
	}
	if uid > 0 {
		mark := "重检———【" + this.User.DisplayName + "】"
		beego.Info(id)
		err = models.UpdateCueditQueueStatusInqueueTime("RECHECK", this.User.DisplayName, this.User.Id, id)
		if err != nil {
			cache.RecordLogs(this.User.Id, 0, this.User.Name, this.User.DisplayName, "重检缓存锁失败", "信审工作平台/UpdateUsersBaseInfo", err.Error(), this.Ctx.Input)
			resultMap["err"] = "重检失败"
			return
		}
		params := map[string]interface{}{
			"Id":             id,
			"Uid":            uid,
			"RiskRecordType": 1,
		}
		_, err := services.PostApi("/check/anewcheck", params)
		beego.Info(err)
		if err != nil {
			err = utils.Rc.Put(utils.CacheKeyRecheck+"_"+strconv.Itoa(uid), params, 5*time.Minute)
			if err != nil {
				cache.RecordLogs(this.User.Id, 0, this.User.Name, this.User.DisplayName, "重检缓存锁失败", "信审工作平台/UpdateUsersBaseInfo", err.Error(), this.Ctx.Input)
				resultMap["err"] = "重检缓存锁失败"
				return
			}
		}

		beego.Info(err)
		err = models.AddCreditAduitRecord(id, uid, mark, "重检")
		beego.Info(err)
		if err != nil {
			cache.RecordLogs(this.User.Id, 0, this.User.Name, this.User.DisplayName, "重检缓存锁失败", "信审工作平台/UpdateUsersBaseInfo", err.Error(), this.Ctx.Input)
			resultMap["err"] = "重检记录失败"
			return
		}

	}
	resultMap["ret"] = 200
	resultMap["msg"] = "更新用户基本信息成功"
}

//重检阻塞
func (c *UsersMetadataController) GetRecheck() {
	defer c.ServeJSON()

	id, err := c.GetInt("uid")
	resultMap := make(map[string]interface{})
	resultMap["ret"] = 403
	if err != nil {
		cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "授信ID获取失败", "信审工作平台/GetRecheck", "", c.Ctx.Input)
		resultMap["err"] = "授信ID获取失败"
		return
	}
	state, err := models.GetCreditState(id)
	if err != nil {
		cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "授信状态获取失败", "信审工作平台/GetRecheck", "", c.Ctx.Input)
		resultMap["err"] = "授信状态获取失败"
		return
	}
	resultMap["ret"] = 200
	resultMap["data"] = state

	c.Data["json"] = resultMap
}

//订单重检
func (this *UsersMetadataController) OrderReview() {
	uid, _ := this.GetInt("ubi_id")
	id, _ := this.GetInt("mainId")
	beego.Info(id, uid)
	resultMap := make(map[string]interface{})
	resultMap["ret"] = 403
	defer func() {
		this.Data["json"] = resultMap
		this.ServeJSON()
	}()
	if flag := cache.GetCacheRecheckOften(uid); flag {
		cache.RecordLogs(this.User.Id, 0, this.User.Name, this.User.DisplayName, "请不要重复重检", "信审工作平台/UpdateUsersBaseInfo", "", this.Ctx.Input)
		resultMap["err"] = "请不要重复重检"
		return
	}

	if uid > 0 {
		params := map[string]interface{}{
			"Id":             id,
			"Uid":            uid,
			"RiskRecordType": 2,
		}
		_, err := services.PostApi("/check/anewcheck", params)
		if err != nil {
			err = utils.Rc.Put(utils.CacheKeyRecheck+"_"+strconv.Itoa(uid), params, 5*time.Minute)
			if err != nil {
				cache.RecordLogs(this.User.Id, 0, this.User.Name, this.User.DisplayName, "重检缓存锁失败", "信审工作平台/OrderReview", err.Error(), this.Ctx.Input)
				resultMap["err"] = "重检缓存锁失败"
				return
			}
		}
		mark := "重检———【" + this.User.DisplayName + "】"
		models.UpdateOrderState(id, this.User.Id, this.User.DisplayName, "RECHECK")
		models.InsertLoanAduitRecord(uid, id, mark, "重检")
		// cache.DeleteOrderRedis(uid, id)
	}
	resultMap["ret"] = 200
	resultMap["msg"] = "更新用户基本信息成功"
}

//获取实名认证信息
func (this *UsersMetadataController) GetRealnameAuthInfo() {
	uid, _ := this.GetInt("uid")
	if uid <= 0 {
		cache.RecordLogs(this.User.Id, 0, this.User.Name, this.User.DisplayName, "uid参数传递错误", "信审工作平台/GetRealnameAuthInfo", "", this.Ctx.Input)
		this.Abort("uid参数传递错误")
		return
	}
	user, err := models.QueryRealnameAuthInfo(uid)
	if err != nil && err.Error() != utils.ErrNoRow() {
		cache.RecordLogs(this.User.Id, 0, this.User.Name, this.User.DisplayName, "查询实名认证信息出错", "信审工作平台/GetRealnameAuthInfo", err.Error(), this.Ctx.Input)
		this.Abort("查询实名认证信息出错" + err.Error())
		return
	}
	ocrInfo, err := models.GetUsersOCRInfo(uid)
	if err != nil {
		cache.RecordLogs(this.User.Id, 0, this.User.Name, this.User.DisplayName, "查询OCR信息出错", "信审工作平台/GetRealnameAuthInfo", err.Error(), this.Ctx.Input)
		this.Abort("查询OCR信息出错" + err.Error())
		return
	}
	identifyInfo, err := models.GetUsersIdentifyInfo(uid)
	if err != nil {
		cache.RecordLogs(this.User.Id, 0, this.User.Name, this.User.DisplayName, "查询人脸比对信息出错", "信审工作平台/GetRealnameAuthInfo", err.Error(), this.Ctx.Input)
		this.Abort("查询人脸比对信息出错" + err.Error())
		return
	}
	this.Data["user"] = user
	this.Data["ocrInfo"] = ocrInfo
	this.Data["identifyInfo"] = identifyInfo
	this.TplName = "user/user_realname_auth.html"
}

//获取芝麻信用信息
func (this *UsersMetadataController) GetZMXYInfo() {
	page, _ := this.GetInt("page")
	if page < 1 {
		page = 1
	}
	uid, _ := this.GetInt("uid")
	if uid <= 0 {
		cache.RecordLogs(this.User.Id, 0, this.User.Name, this.User.DisplayName, "uid参数传递错误", "信审工作平台/GetZMXYInfo", "", this.Ctx.Input)
		this.Abort("uid参数传递错误")
		return
	}
	list, err := models.QueryZmxyInfos(uid, utils.StartIndex(page, utils.PageSize30), utils.PageSize30)
	if err != nil {
		cache.RecordLogs(this.User.Id, 0, this.User.Name, this.User.DisplayName, "查询用户芝麻信用出错", "信审工作平台/GetZMXYInfo", err.Error(), this.Ctx.Input)
		this.Abort("查询用户芝麻信用出错" + err.Error())
		return
	}
	count, err := models.QueryZmxyCount(uid)
	if err != nil {
		cache.RecordLogs(this.User.Id, 0, this.User.Name, this.User.DisplayName, "查询用户芝麻信用总数出错", "信审工作平台/GetZMXYInfo", err.Error(), this.Ctx.Input)
		this.Abort("查询用户芝麻信用总数出错" + err.Error())
		return
	}
	pageCount := utils.PageCount(count, utils.PageSize30)
	this.Data["data"] = list
	this.Data["uid"] = uid
	this.Data["currPage"] = page
	this.Data["count"] = count
	this.Data["pageSize"] = utils.PageSize30
	this.Data["pageCount"] = pageCount
	this.TplName = "user/user_zmxy.html"
}

//获取常用联系人
func (this *UsersMetadataController) GetLinkmanList() {
	uid, _ := this.GetInt("uid")
	riskRecordType, _ := this.GetInt("riskRecordType")
	beego.Info(riskRecordType)
	if uid <= 0 {
		cache.RecordLogs(this.User.Id, 0, this.User.Name, this.User.DisplayName, "uid参数传递错误", "信审工作平台/GetLinkmanList", "", this.Ctx.Input)
		this.Abort("uid参数传递错误")
		return
	}

	isEdit, _ := this.GetBool("is_edit", false)
	linkmans, err := models.QueryLinkmans(uid)
	if err != nil {
		cache.RecordLogs(this.User.Id, 0, this.User.Name, this.User.DisplayName, "查询联系人出错", "信审工作平台/GetLinkmanList", err.Error(), this.Ctx.Input)
		this.Abort("查询联系人出错" + err.Error())
		return
	}
	//从Mongodb上获取
	var mdbMxYYSData models.MdbMxYYSData
	uidMap := make(map[string]interface{})
	uidMap["uid"] = uid
	session := utils.GetSession()
	defer session.Close()
	err = session.DB(utils.MGO_DB).C("mxdata").Find(uidMap).Sort("-createtime").Limit(1).One(&mdbMxYYSData)
	if err != nil && err.Error() != utils.MongdbErrNoRow() {
		cache.RecordLogs(this.User.Id, 0, this.User.Name, this.User.DisplayName, "mgodb数据获取mxdata异常", "信审工作平台/GetLinkmanList", err.Error(), this.Ctx.Input)
		this.Abort("mgodb数据获取mxdata异常" + err.Error())
		return
	}
	var finish sync.WaitGroup
	finish.Add(len(linkmans))
	//获取手机号归属地和近一个月通话记录
	for k, v := range linkmans {
		go services.YYSMonthCallData(mdbMxYYSData.MXYYSData.Calls, mdbMxYYSData.CreateTime, v.ContactPhoneNumber, &linkmans[k], &finish)
	}
	finish.Wait()
	beginTime := mdbMxYYSData.CreateTime
	before30 := beginTime.AddDate(0, 0, -29) //前30天
	this.Data["beginTime"] = beginTime
	this.Data["before30"] = before30
	this.Data["isEdit"] = isEdit
	this.Data["sortLinkmans"] = linkmans
	this.Data["type"] = 1
	this.Data["riskRecordType"] = riskRecordType
	this.TplName = "user/user_linkmans.html"
}

//更新用户联系人
func (this *UsersMetadataController) UpdateLinkmans() {
	resultMap := make(map[string]interface{})
	resultMap["ret"] = 403
	defer func() {
		this.Data["json"] = resultMap
		this.ServeJSON()
	}()
	var linkmans []models.Linkman
	err := json.Unmarshal(this.Ctx.Input.RequestBody, &linkmans)
	beego.Info(linkmans[0].RiskRecordType)

	if err != nil {
		cache.RecordLogs(this.User.Id, 0, this.User.Name, this.User.DisplayName, "参数解析错误", "信审工作平台/UpdateLinkmans", err.Error(), this.Ctx.Input)
		resultMap["err"] = "参数解析错误"
		return
	}
	err = models.UpdateLinkmans(linkmans)
	if err != nil {
		cache.RecordLogs(this.User.Id, 0, this.User.Name, this.User.DisplayName, "更新联系人出错", "信审工作平台/UpdateLinkmans", err.Error(), this.Ctx.Input)
		resultMap["err"] = "更新联系人出错"
		return
	}
	var id int
	beego.Info(linkmans[0].Id)
	if len(linkmans) > 0 {
		id, err = models.GetCreditId(linkmans[0].Uid)
		if flag := cache.GetCacheRecheckOften(linkmans[0].Uid); flag {
			cache.RecordLogs(this.User.Id, 0, this.User.Name, this.User.DisplayName, "请不要重复重检", "信审工作平台/UpdateUsersBaseInfo", "", this.Ctx.Input)
			resultMap["err"] = "请不要重复重检"
			return
		}
		params := map[string]interface{}{
			"Id":             id,
			"Uid":            linkmans[0].Uid,
			"RiskRecordType": linkmans[0].RiskRecordType,
		}
		_, err = services.PostApi("/check/anewcheck", params)
		if err != nil {
			err = utils.Rc.Put(utils.CacheKeyRecheck+"_"+strconv.Itoa(linkmans[0].Uid), params, 5*time.Minute)
			if err != nil {
				cache.RecordLogs(this.User.Id, 0, this.User.Name, this.User.DisplayName, "重检缓存锁失败", "信审工作平台/OrderReview", err.Error(), this.Ctx.Input)
				resultMap["err"] = "重检缓存锁失败"
				return
			}
		}

		if linkmans[0].RiskRecordType == 1 {
			mark := "重检———【" + this.User.DisplayName + "】"
			err := models.UpdateCueditQueueStatusInqueueTime("RECHECK", this.User.DisplayName, this.User.Id, id)
			beego.Info(err)
			err = models.AddCreditAduitRecord(id, linkmans[0].Uid, mark, "重检")
			beego.Info(err)
		} else {
			mark := "重检———【" + this.User.DisplayName + "】"
			err := models.UpdateOrderState(linkmans[0].Id, this.User.Id, this.User.DisplayName, "RECHECK")
			beego.Info(err)
			err = models.InsertLoanAduitRecord(linkmans[0].Uid, linkmans[0].Id, mark, "重检")
			beego.Info(err)
			cache.DeleteOrderRedis(linkmans[0].Uid, linkmans[0].Id)

		}

	}
	resultMap["ret"] = 200
	resultMap["msg"] = "更新联系人成功"
}

//手机通讯录
func (this *UsersMetadataController) GetPhoneBook() {
	page, _ := this.GetInt("page")
	if page <= 0 {
		page = 1
	}
	pageSize, _ := this.GetInt("pageSize")
	if pageSize < 1 {
		pageSize = 15
	}
	uid, _ := this.GetInt("uid")
	if uid <= 0 {
		this.Data["json"] = map[string]interface{}{"ret": 403, "msg": "uid参数异常！"}
		return
	}
	var contactlist = []models.ContactData{}
	var contact = models.ContactData{}
	txlist, err2 := models.QueryTelephonInfo(uid)
	isHas := true
	if err2 != nil {
		cache.RecordLogs(this.User.Id, 0, this.User.Name, this.User.DisplayName, "获取通讯录信息异常", "个人信息/手机联系人PhoneLinkman", err2.Error(), this.Ctx.Input)
	}
	if len(txlist) > 0 {
		txlist2 := services.SortAddBook(txlist)
		if len(txlist2[0].Contact) > 0 {
			for _, v := range txlist2[0].Contact {
				if len(v.ContactPhoneNumber) > 0 {
					contact.ContactPhoneNumber = v.ContactPhoneNumber[0]
					contact.ContactName = v.ContactName
					contactlist = append(contactlist, contact)
				}
			}
		} else {
			this.Abort("请开启获取手机通讯录权限")
		}
		count := len(contactlist)
		if count > 0 {
			if res, err := json.Marshal(contactlist); err == nil {
				utils.Rc.Put("xjfq_phone_info:"+"_"+strconv.Itoa(uid), res, utils.RedisCacheTime_7Day)
			}
		}
		pageCount := utils.PageCount(count, pageSize)
		start := utils.StartIndex(page, pageSize)
		num := utils.StartIndex(page, pageSize) + pageSize
		if count >= num {
			contactlist = contactlist[start:num]
		} else {
			contactlist = contactlist[start:]
		}
		this.Data["list"] = contactlist
		this.Data["currPage"] = page
		this.Data["count"] = count
		this.Data["pageSize"] = pageSize
		this.Data["pageCount"] = pageCount
		this.Data["uid"] = uid
	} else {
		isHas = false
	}
	this.Data["type"] = 2
	this.Data["is_Has"] = isHas
	this.TplName = "user/user_linkmans.html"
}

//通话记录
func (this *UsersMetadataController) GetCallRecord() {
	session := utils.GetSession()
	defer session.Close()
	page, _ := this.GetInt("page")
	if page <= 0 {
		page = 1
	}
	pageSize, _ := this.GetInt("pageSize")
	if pageSize < 1 {
		pageSize = 15
	}
	uid, _ := this.GetInt("uid")
	if uid <= 0 {
		cache.RecordLogs(this.User.Id, 0, this.User.Name, this.User.DisplayName, "uid参数传递错误", "信审工作平台/GetCallRecord", "", this.Ctx.Input)
		this.Abort("uid参数传递错误")
		return
	}
	mobileAuthType, err := models.GetMobileAuthTypeByUserId(uid)
	if err != nil && err.Error() != utils.ErrNoRow() {
		cache.RecordLogs(this.User.Id, 0, this.User.Name, this.User.DisplayName, "根据用户ID获取运营商授权类型失败", "个人信息/GetCallRecord", err.Error(), this.Ctx.Input)
	}
	isHas := true
	if mobileAuthType == 1 { //魔蝎
		var userMxData models.MxrThreeReportData
		session := utils.GetSession()
		defer session.Close()
		err := session.DB(utils.MGO_DB).C("mxreportdata").Find(&models.MonGoQuery{Uid: uid}).Sort("-createtime").One(&userMxData)
		if err != nil {
			cache.RecordLogs(this.User.Id, 0, this.User.Name, this.User.DisplayName, "mgodb数据获取异常", "个人信息/运营商数据MxreportData", err.Error(), this.Ctx.Input)
			this.Abort("请开启获取手机通讯录权限" + err.Error())
			return
		}

		if len(userMxData.Rt.Ccd) == 0 {
			cache.RecordLogs(this.User.Id, 0, this.User.Name, this.User.DisplayName, "mgodb数据获取为空", "个人信息/运营商数据MxreportData", "", this.Ctx.Input)
			this.Abort("请开启获取手机通讯录权限")
			return
		}

		var mx_data_list []models.Call_Contact_Detail
		for i := 0; i < len(userMxData.Rt.Ccd); i++ {
			userMxData.Rt.Ccd[i].Total_Cnt = userMxData.Rt.Ccd[i].Dial_Cnt_3m + userMxData.Rt.Ccd[i].Dialed_Cnt_3m
			mx_data_list = append(mx_data_list, userMxData.Rt.Ccd[i])

		}
		sort_mx_data_list := services.SortMXCntTime(mx_data_list)
		if len(sort_mx_data_list) >= 15 {
			sort_mx_data_list = sort_mx_data_list[:15]
		}
		for j := 0; j < len(sort_mx_data_list); j++ {
			sort_mx_data_list[j].ContactName = get_match_records(uid, sort_mx_data_list[j].Peer_num, this)
		}
		//存缓存
		if len(sort_mx_data_list) > 0 {
			if res, err := json.Marshal(sort_mx_data_list); err == nil {
				utils.Rc.Put("xjfq_call_info:"+"_"+strconv.Itoa(uid), res, utils.RedisCacheTime_7Day)
			}
		}
		this.Data["list"] = sort_mx_data_list
		this.Data["mtype"] = 1
		this.Data["uid"] = uid
		this.Data["currPage"] = page
		this.Data["count"] = len(sort_mx_data_list)
		this.Data["pageCount"] = len(sort_mx_data_list)
	} else if mobileAuthType == 2 { //天玑
		var tjr models.Tijireport
		session := utils.GetSession()
		defer session.Close()
		err := session.DB(utils.MGO_DB).C("tianjireport").Find(&models.MonGoQuery{Uid: uid}).Sort("-createtime").One(&tjr) //900000451
		if err != nil {
			cache.RecordLogs(this.User.Id, 0, this.User.Name, this.User.DisplayName, "mgodb数据获取异常", "个人信息/运营商数据OperatorData", err.Error(), this.Ctx.Input)
			this.Abort("mgodb数据获取tianjireport异常" + err.Error())
			return
		}
		if len(tjr.Tianji.CallLog) == 0 {
			cache.RecordLogs(this.User.Id, 0, this.User.Name, this.User.DisplayName, "mgodb数据获取为空", "个人信息/运营商数据OperatorData", err.Error(), this.Ctx.Input)
			this.Abort("mgodb数据获取tianjireport异常" + err.Error())
			return
		}
		var tianjiList []models.CallLogInfo
		for i := 0; i < len(tjr.Tianji.CallLog); i++ {
			tianjiList = append(tianjiList, tjr.Tianji.CallLog[i])
		}
		sort_tianji := services.SortTradeTime(tianjiList)
		if len(sort_tianji) >= 15 {
			sort_tianji = sort_tianji[:15]
		}
		for j := 0; j < len(sort_tianji); j++ {
			sort_tianji[j].ContactName = get_match_records(uid, sort_tianji[j].Phone, this)
		}
		//存缓存
		if len(sort_tianji) > 0 {
			if res, err := json.Marshal(sort_tianji); err == nil {
				utils.Rc.Put("xjfq_call_info:"+"_"+strconv.Itoa(uid), res, utils.RedisCacheTime_7Day)
			}
		}
		this.Data["list"] = sort_tianji
		this.Data["mtype"] = 2
		this.Data["uid"] = uid
		this.Data["currPage"] = page
		this.Data["count"] = len(sort_tianji)
		this.Data["pageCount"] = len(sort_tianji)
	} else {
		isHas = false
	}
	this.Data["isHas"] = isHas
	this.Data["type"] = 3
	this.TplName = "user/user_linkmans.html"
}

//匹配通讯记录
func get_match_records(uid int, phoneNum string, this *UsersMetadataController) (contact_name string) {
	var contactlist = []models.ContactData{}
	//缓存
	if utils.Rc.IsExist("xjfq_phone_info:" + "_" + strconv.Itoa(uid)) {
		if data, err := utils.Rc.RedisBytes("xjfq_phone_info:" + "_" + strconv.Itoa(uid)); err == nil {
			json.Unmarshal(data, &contactlist)
			//匹配处理
			for _, v1 := range contactlist {
				if v1.ContactPhoneNumber == phoneNum {
					return v1.ContactName
				}
			}
		}
	} else {
		txlist, err2 := models.QueryTelephonInfo(uid)
		if err2 != nil {
			cache.RecordLogs(this.User.Id, 0, this.User.Name, this.User.DisplayName, "获取通讯录信息异常", "个人信息/手机联系人PhoneLinkman", err2.Error(), this.Ctx.Input)
		}
		if len(txlist) > 0 {
			txlist2 := services.SortAddBook(txlist)
			if len(txlist2[0].Contact) > 0 {
				for _, v := range txlist2[0].Contact {
					if len(v.ContactPhoneNumber) > 0 {
						if v.ContactPhoneNumber[0] == phoneNum {
							return v.ContactName
						}
					}
				}
			}
		}
	}
	return
}

//匹配通讯记录
func (this *UsersMetadataController) GetMatchRecord() {
	defer this.ServeJSON()
	uid, _ := this.GetInt("uid")
	if uid <= 0 {
		this.Data["json"] = map[string]interface{}{"ret": 403, "msg": "uid参数异常！"}
		return
	}
	phoneNum := this.GetString("phoneNum")
	txlist, err2 := models.QueryTelephonInfo(uid)
	if err2 != nil {
		cache.RecordLogs(this.User.Id, 0, this.User.Name, this.User.DisplayName, "获取通讯录信息异常", "个人信息/手机联系人PhoneLinkman", err2.Error(), this.Ctx.Input)
	}
	if len(txlist) > 0 {
		txlist2 := services.SortAddBook(txlist)
		if len(txlist2[0].Contact) > 0 {
			for _, v := range txlist2[0].Contact {
				if len(v.ContactPhoneNumber) > 0 {
					if v.ContactPhoneNumber[0] == phoneNum {
						this.Data["json"] = map[string]interface{}{"ret": 200, "relationship": v.ContactName, "phoneNum": phoneNum, "msg": "匹配成功"}
						return
					}
				}
			}
		} else {
			this.Abort("请开启获取手机通讯录权限")
		}
	}
	this.Data["json"] = map[string]interface{}{"ret": 400, "msg": "匹配不成功"}
	return
}

//匹配紧急联系人接口
func (this *UsersMetadataController) GetMatchCallRecord() {
	defer this.ServeJSON()
	uid, _ := this.GetInt("uid")
	if uid <= 0 {
		this.Data["json"] = map[string]interface{}{"ret": 403, "msg": "uid参数异常！"}
		return
	}
	//匹配手机号
	phoneNum := this.GetString("phoneNum")
	var match = models.MatchCall{}
	match.ConnectFlag = false
	match.MobileFlag = false
	txlist, err2 := models.QueryTelephonInfo(uid)
	if err2 != nil {
		cache.RecordLogs(this.User.Id, 0, this.User.Name, this.User.DisplayName, "获取通讯录信息异常", "个人信息/手机联系人PhoneLinkman", err2.Error(), this.Ctx.Input)
	}
	if len(txlist) > 0 {
		txlist2 := services.SortAddBook(txlist)
		if len(txlist2[0].Contact) > 0 {
			for _, v := range txlist2[0].Contact {
				if len(v.ContactPhoneNumber) > 0 {
					if v.ContactPhoneNumber[0] == phoneNum {
						match.ContactName = v.ContactName
						match.ContactPhoneNumber = phoneNum
						match.MobileFlag = true
					}
				}
			}
		} else {
			this.Abort("请开启获取手机通讯录权限")
		}
	}
	//联系次数
	mobileAuthType, err := models.GetMobileAuthTypeByUserId(uid)
	if err != nil && err.Error() != utils.ErrNoRow() {
		cache.RecordLogs(this.User.Id, 0, this.User.Name, this.User.DisplayName, "根据用户ID获取运营商授权类型失败", "个人信息/GetCallRecord", err.Error(), this.Ctx.Input)
	}
	//魔蝎
	if mobileAuthType == 1 {
		var mx_data_list []models.Call_Contact_Detail
		if utils.Rc.IsExist("xjfq_call_info:" + "_" + strconv.Itoa(uid)) {
			if data, err := utils.Rc.RedisBytes("xjfq_call_info:" + "_" + strconv.Itoa(uid)); err == nil {
				json.Unmarshal(data, &mx_data_list)
				//匹配处理
				for j := 0; j < len(mx_data_list); j++ {
					if mx_data_list[j].Peer_num == phoneNum {
						match.ConnectCount = mx_data_list[j].Total_Cnt
						match.ConnectFlag = true
					}
				}
			}
		} else {
			var userMxData models.MxrThreeReportData
			session := utils.GetSession()
			defer session.Close()
			err := session.DB(utils.MGO_DB).C("mxreportdata").Find(&models.MonGoQuery{Uid: uid}).Sort("-createtime").One(&userMxData)
			if err != nil {
				cache.RecordLogs(this.User.Id, 0, this.User.Name, this.User.DisplayName, "mgodb数据获取异常", "个人信息/运营商数据MxreportData", err.Error(), this.Ctx.Input)
				this.Abort("请开启获取手机通讯录权限" + err.Error())
				return
			}

			if len(userMxData.Rt.Ccd) == 0 {
				cache.RecordLogs(this.User.Id, 0, this.User.Name, this.User.DisplayName, "mgodb数据获取为空", "个人信息/运营商数据MxreportData", "", this.Ctx.Input)
				this.Abort("请开启获取手机通讯录权限")
				return
			}

			var mx_data_list []models.Call_Contact_Detail
			for i := 0; i < len(userMxData.Rt.Ccd); i++ {
				userMxData.Rt.Ccd[i].Total_Cnt = userMxData.Rt.Ccd[i].Dial_Cnt_3m + userMxData.Rt.Ccd[i].Dialed_Cnt_3m
				mx_data_list = append(mx_data_list, userMxData.Rt.Ccd[i])

			}
			sort_mx_data_list := services.SortMXCntTime(mx_data_list)
			if len(sort_mx_data_list) >= 15 {
				sort_mx_data_list = sort_mx_data_list[:15]
			}
			for j := 0; j < len(sort_mx_data_list); j++ {
				if sort_mx_data_list[j].Peer_num == phoneNum {
					match.ConnectCount = sort_mx_data_list[j].Total_Cnt
					match.ConnectFlag = true
				}
			}
		}
		//天玑
	} else if mobileAuthType == 2 {
		var tianjiList []models.CallLogInfo
		if utils.Rc.IsExist("xjfq_call_info:" + "_" + strconv.Itoa(uid)) {
			if data, err := utils.Rc.RedisBytes("xjfq_call_info:" + "_" + strconv.Itoa(uid)); err == nil {
				json.Unmarshal(data, &tianjiList)
				//匹配处理
				for j := 0; j < len(tianjiList); j++ {
					if tianjiList[j].Phone == phoneNum {
						match.ConnectCount = tianjiList[j].Talk_cnt
						match.ConnectFlag = true
					}
				}
			}
		} else {
			var tjr models.Tijireport
			session := utils.GetSession()
			defer session.Close()
			err := session.DB(utils.MGO_DB).C("tianjireport").Find(&models.MonGoQuery{Uid: uid}).Sort("-createtime").One(&tjr) //900000451
			if err != nil {
				cache.RecordLogs(this.User.Id, 0, this.User.Name, this.User.DisplayName, "mgodb数据获取异常", "个人信息/运营商数据OperatorData", err.Error(), this.Ctx.Input)
				this.Abort("mgodb数据获取tianjireport异常" + err.Error())
				return
			}
			if len(tjr.Tianji.CallLog) == 0 {
				cache.RecordLogs(this.User.Id, 0, this.User.Name, this.User.DisplayName, "mgodb数据获取为空", "个人信息/运营商数据OperatorData", err.Error(), this.Ctx.Input)
				this.Abort("mgodb数据获取tianjireport异常" + err.Error())
				return
			}
			var tianjiList []models.CallLogInfo
			for i := 0; i < len(tjr.Tianji.CallLog); i++ {
				tianjiList = append(tianjiList, tjr.Tianji.CallLog[i])
			}
			sort_tianji := services.SortTradeTime(tianjiList)
			if len(sort_tianji) >= 15 {
				sort_tianji = sort_tianji[:15]
			}
			for j := 0; j < len(sort_tianji); j++ {
				if sort_tianji[j].Phone == phoneNum {
					match.ConnectCount = tianjiList[j].Talk_cnt
					match.ConnectFlag = true
				}
			}
		}
	}

	//运营商
	phone_info := utils.QueryLocating(match.ContactPhoneNumber)
	match.MonType = strings.Split(phone_info, " ")[2]
	this.Data["json"] = map[string]interface{}{"ret": 200, "data": match}
	return
}

//短信记录
func (this *UsersMetadataController) GetMessRecord() {
	uid, _ := this.GetInt("uid")
	if uid <= 0 {
		cache.RecordLogs(this.User.Id, 0, this.User.Name, this.User.DisplayName, "uid参数传递错误", "信审工作平台/GetCallRecord", "", this.Ctx.Input)
		this.Abort("uid参数传递错误")
		return
	}
	page, _ := this.GetInt("page")
	if page <= 0 {
		page = 1
	}
	pageSize, _ := this.GetInt("pageSize")
	if pageSize < 1 {
		pageSize = 15
	}
	this.Data["pageSize"] = pageSize
	this.Data["currPage"] = page
	this.Data["count"] = 1
	this.TplName = "user/user_note_history.html"
}

//匹配联系人
func (this *UsersMetadataController) GetMatchLinkMan() {
	defer this.ServeJSON()
	session := utils.GetSession()
	defer session.Close()
	uid, _ := this.GetInt("uid")
	if uid <= 0 {
		cache.RecordLogs(this.User.Id, 0, this.User.Name, this.User.DisplayName, "uid参数传递错误", "信审工作平台/GetCallRecord", "", this.Ctx.Input)
		this.Abort("uid参数传递错误")
		return
	}
	phoneNum := this.GetString("phoneNum")
	mobileAuthType, err := models.GetMobileAuthTypeByUserId(uid)
	if err != nil {
		cache.RecordLogs(this.User.Id, 0, this.User.Name, this.User.DisplayName, "根据用户ID获取运营商授权类型失败", "个人信息/GetCallRecord", err.Error(), this.Ctx.Input)
	}

	if mobileAuthType == 1 { //魔蝎
		var userMxData models.MxrThreeReportData
		session := utils.GetSession()
		defer session.Close()
		err := session.DB(utils.MGO_DB).C("mxreportdata").Find(&models.MonGoQuery{Uid: uid}).Sort("-createtime").One(&userMxData)
		if err != nil {
			cache.RecordLogs(this.User.Id, 0, this.User.Name, this.User.DisplayName, "mgodb数据获取异常", "个人信息/运营商数据MxreportData", err.Error(), this.Ctx.Input)
			this.Abort("请开启获取手机通讯录权限" + err.Error())
			return
		}
		if len(userMxData.Rt.Ccd) == 0 {
			cache.RecordLogs(this.User.Id, 0, this.User.Name, this.User.DisplayName, "mgodb数据获取为空", "个人信息/运营商数据MxreportData", err.Error(), this.Ctx.Input)
			this.Abort("请开启获取手机通讯录权限" + err.Error())
			return
		}
		var mx_data_list []models.Call_Contact_Detail
		for i := 0; i < len(userMxData.Rt.Ccd); i++ {
			userMxData.Rt.Ccd[i].Total_Cnt = userMxData.Rt.Ccd[i].Dial_Cnt_3m + userMxData.Rt.Ccd[i].Dialed_Cnt_3m
			mx_data_list = append(mx_data_list, userMxData.Rt.Ccd[i])
		}
		sort_mx_data_list := services.SortMXCntTime(mx_data_list)
		if len(sort_mx_data_list) >= 15 {
			sort_mx_data_list = sort_mx_data_list[:15]
		}
		for _, v := range sort_mx_data_list {
			if v.Peer_num == phoneNum {
				this.Data["json"] = map[string]interface{}{"ret": 200, "phone": v.Peer_num, "city": v.City, "total_cnt": v.Total_Cnt}
				return
			}
		}
		this.Data["json"] = map[string]interface{}{"ret": 400, "msg": "无结果"}
		return
	} else if mobileAuthType == 2 {
		var tjr models.Tijireport
		session := utils.GetSession()
		defer session.Close()
		err := session.DB(utils.MGO_DB).C("tianjireport").Find(&models.MonGoQuery{Uid: uid}).Sort("-createtime").One(&tjr) //900000451
		if err != nil {
			cache.RecordLogs(this.User.Id, 0, this.User.Name, this.User.DisplayName, "mgodb数据获取异常", "个人信息/运营商数据OperatorData", err.Error(), this.Ctx.Input)
			this.Abort("mgodb数据获取tianjireport异常" + err.Error())
			return
		}
		if len(tjr.Tianji.CallLog) == 0 {
			cache.RecordLogs(this.User.Id, 0, this.User.Name, this.User.DisplayName, "mgodb数据获取为空", "个人信息/运营商数据OperatorData", err.Error(), this.Ctx.Input)
			this.Abort("mgodb数据获取tianjireport异常" + err.Error())
			return
		}
		var tianjiList []models.CallLogInfo
		for i := 0; i < len(tjr.Tianji.CallLog); i++ {
			tianjiList = append(tianjiList, tjr.Tianji.CallLog[i])
		}
		sort_tianji := services.SortTradeTime(tianjiList)
		if len(sort_tianji) >= 15 {
			sort_tianji = sort_tianji[:15]
		}
		for _, v := range sort_tianji {
			if v.Phone == phoneNum {
				this.Data["json"] = map[string]interface{}{"ret": 200, "phone": v.Phone, "city": v.Phone_location, "total_cnt": v.Talk_cnt}
				return
			}
		}
		this.Data["json"] = map[string]interface{}{"ret": 400, "msg": "无结果"}
		return
	}
}

//获取公积金信息
func (this *UsersMetadataController) GetUsersGJJInfo() {
	uid, _ := this.GetInt("uid")
	if uid <= 0 {
		cache.RecordLogs(this.User.Id, 0, this.User.Name, this.User.DisplayName, "uid参数传递错误", "信审工作平台/GetCallRecord", "", this.Ctx.Input)
		this.Abort("uid参数传递错误")
		return
	}
	this.Data["uid"] = uid
	//从Mongodb上获取
	var mdbMxGJJData models.MdbMxGJJData
	uidMap := make(map[string]interface{})
	uidMap["uid"] = uid
	session := utils.GetSession()
	defer session.Close()
	err := session.DB(utils.MGO_DB).C("mxfunddata").Find(uidMap).Sort("-createtime").Select(bson.M{"mxgjjdata.user_info": 1, "mxgjjdata.bill_record": 1}).Limit(1).One(&mdbMxGJJData)
	if err != nil && err.Error() != utils.MongdbErrNoRow() {
		cache.RecordLogs(this.User.Id, 0, this.User.Name, this.User.DisplayName, "mgodb数据获取mxfunddata异常", "信审工作平台/GetUsersGJJInfo", err.Error(), this.Ctx.Input)
		this.Abort("mgodb数据获取mxfunddata异常" + err.Error())
		return
	}
	var isGJJShow = true
	if err != nil && err.Error() == utils.MongdbErrNoRow() {
		isGJJShow = false
	}
	services.GJJQuickSort(mdbMxGJJData.MXGJJData.Bill_Record)
	if len(mdbMxGJJData.MXGJJData.Bill_Record) >= 12 {
		mdbMxGJJData.MXGJJData.Bill_Record = mdbMxGJJData.MXGJJData.Bill_Record[:12]
	}
	this.Data["isGJJShow"] = isGJJShow
	this.Data["data"] = mdbMxGJJData
	this.TplName = "user/user_gjj_info.html"
}

//获取公积金中的贷款信息
func (this *UsersMetadataController) GetGJJLoanInfo() {
	resultMap := make(map[string]interface{})
	resultMap["ret"] = 403
	defer func() {
		this.Data["json"] = resultMap
		this.ServeJSON()
	}()
	uid, _ := this.GetInt("uid")
	if uid <= 0 {
		cache.RecordLogs(this.User.Id, 0, this.User.Name, this.User.DisplayName, "uid参数传递错误", "信审工作平台/GetGJJLoanInfo", "", this.Ctx.Input)
		resultMap["err"] = "uid参数传递错误"
		return
	}
	//从Mongodb上获取
	var mdbMxGJJData models.MdbMxGJJData
	uidMap := make(map[string]interface{})
	uidMap["uid"] = uid
	session := utils.GetSession()
	defer session.Close()
	err := session.DB(utils.MGO_DB).C("mxfunddata").Find(uidMap).Sort("-createtime").Select(bson.M{"mxgjjdata.loan_info": 1}).Limit(1).One(&mdbMxGJJData)
	if err != nil && err.Error() != utils.MongdbErrNoRow() {
		cache.RecordLogs(this.User.Id, 0, this.User.Name, this.User.DisplayName, "mgodb数据获取mxfunddata异常", "信审工作平台/GetGJJLoanInfo", err.Error(), this.Ctx.Input)
		resultMap["err"] = "mgodb数据获取mxfunddata异常" + err.Error()
		return
	}
	page, _ := this.GetInt("page")
	if page < 1 {
		page = 1
	}
	count := len(mdbMxGJJData.MXGJJData.Loan_Info)
	reqLoanCount := utils.StartIndex(page, utils.PageSize10) + utils.PageSize10
	reqLoanStart := utils.StartIndex(page, utils.PageSize10)
	if len(mdbMxGJJData.MXGJJData.Loan_Info) >= reqLoanCount {
		mdbMxGJJData.MXGJJData.Loan_Info = mdbMxGJJData.MXGJJData.Loan_Info[reqLoanStart:reqLoanCount]
	} else {
		mdbMxGJJData.MXGJJData.Loan_Info = mdbMxGJJData.MXGJJData.Loan_Info[reqLoanStart:]
	}
	pageCount := utils.PageCount(count, utils.PageSize10)
	resultMap["ret"] = 200
	resultMap["uid"] = uid
	resultMap["data"] = mdbMxGJJData.MXGJJData.Loan_Info
	resultMap["pageSize"] = utils.PageSize10
	resultMap["count"] = count
	resultMap["currPage"] = page
	resultMap["pageCount"] = pageCount

}

//获取公积金中的还款信息
func (this *UsersMetadataController) GetGJJRepayInfo() {
	resultMap := make(map[string]interface{})
	resultMap["ret"] = 403
	defer func() {
		this.Data["json"] = resultMap
		this.ServeJSON()
	}()
	uid, _ := this.GetInt("uid")
	if uid <= 0 {
		cache.RecordLogs(this.User.Id, 0, this.User.Name, this.User.DisplayName, "uid参数传递错误", "信审工作平台/GetGJJRepayInfo", "", this.Ctx.Input)
		resultMap["err"] = "uid参数传递错误"
		return
	}
	resultMap["uid"] = uid
	//从Mongodb上获取
	var mdbMxGJJData models.MdbMxGJJData
	uidMap := make(map[string]interface{})
	uidMap["uid"] = uid
	session := utils.GetSession()
	defer session.Close()
	err := session.DB(utils.MGO_DB).C("mxfunddata").Find(uidMap).Sort("-createtime").Select(bson.M{"mxgjjdata.loan_repay_record": 1}).Limit(1).One(&mdbMxGJJData)
	if err != nil && err.Error() != utils.MongdbErrNoRow() {
		cache.RecordLogs(this.User.Id, 0, this.User.Name, this.User.DisplayName, "mgodb数据获取mxfunddata异常", "信审工作平台/GetGJJRepayInfo", err.Error(), this.Ctx.Input)
		resultMap["err"] = "mgodb数据获取mxfunddata异常" + err.Error()
		return
	}
	page, _ := this.GetInt("page")
	if page < 1 {
		page = 1
	}
	count := len(mdbMxGJJData.MXGJJData.Loan_Repay_Record)
	reqRepayCount := utils.StartIndex(page, utils.PageSize10) + utils.PageSize10
	reqRepayStart := utils.StartIndex(page, utils.PageSize10)
	if len(mdbMxGJJData.MXGJJData.Loan_Repay_Record) >= reqRepayCount {
		mdbMxGJJData.MXGJJData.Loan_Repay_Record = mdbMxGJJData.MXGJJData.Loan_Repay_Record[reqRepayStart:reqRepayCount]
	} else {
		mdbMxGJJData.MXGJJData.Loan_Repay_Record = mdbMxGJJData.MXGJJData.Loan_Repay_Record[reqRepayStart:]
	}
	pageCount := utils.PageCount(count, utils.PageSize10)
	resultMap["ret"] = 200
	resultMap["uid"] = uid
	resultMap["data"] = mdbMxGJJData.MXGJJData.Loan_Repay_Record
	resultMap["pageSize"] = utils.PageSize10
	resultMap["count"] = count
	resultMap["currPage"] = page
	resultMap["pageCount"] = pageCount
}

//获取支付宝基本信息
func (this *UsersMetadataController) GetUsersZFBInfo() {
	uid, _ := this.GetInt("uid")
	if uid <= 0 {
		cache.RecordLogs(this.User.Id, 0, this.User.Name, this.User.DisplayName, "uid参数传递错误", "信审工作平台/GetUsersZFBInfo", "", this.Ctx.Input)
		this.Abort("uid参数传递错误")
		return
	}
	//从Mongodb上获取
	var mdbMxZFBData models.MdbMxZFBData
	uidMap := make(map[string]interface{})
	uidMap["uid"] = uid
	session := utils.GetSession()
	defer session.Close()
	err := session.DB(utils.MGO_DB).C("mxalipaydata").Find(uidMap).Sort("-createtime").Select(bson.M{"mxzfbdata.userinfo": 1, "mxzfbdata.wealth": 1, "mxzfbdata.bankinfo": 1}).Limit(1).One(&mdbMxZFBData)
	if err != nil && err.Error() != utils.MongdbErrNoRow() {
		cache.RecordLogs(this.User.Id, 0, this.User.Name, this.User.DisplayName, "mgodb数据获取mxalipaydata异常", "信审工作平台/GetUsersZFBInfo", err.Error(), this.Ctx.Input)
		this.Abort("mgodb数据获取mxalipaydata异常" + err.Error())
		return
	}
	isZFBShow := true
	if err != nil && err.Error() == utils.MongdbErrNoRow() {
		isZFBShow = false
	}
	var mdbMxZFBZMData models.MdbMxZFBZMData
	err2 := session.DB(utils.MGO_DB).C("mxalipayzm").Find(uidMap).Sort("-createtime").Select(bson.M{"data": 1, "createtime": 1}).Limit(1).One(&mdbMxZFBZMData)
	if err2 != nil && err2.Error() != utils.MongdbErrNoRow() {
		cache.RecordLogs(this.User.Id, 0, this.User.Name, this.User.DisplayName, "mgodb数据获取mxalipayzm异常", "信审工作平台/GetUsersZFBInfo", err2.Error(), this.Ctx.Input)
		this.Abort("mgodb数据获取mxalipaydata异常" + err2.Error())
		return
	}
	isZMShow := true
	if err2 != nil && err2.Error() == utils.MongdbErrNoRow() {
		isZMShow = false
	}
	nearly3MonthZM := make([]models.MXZFBZMData, 0)
	if err2 == nil {
		if len(mdbMxZFBZMData.MXZFBZMData) == 0 {
			isZMShow = false
		} else {
			_, _, day := mdbMxZFBZMData.CreateTime.Date()
			startTime := mdbMxZFBZMData.CreateTime.Format("200601")
			before1MTime := mdbMxZFBZMData.CreateTime.AddDate(0, -1, 1-day).Format("200601")
			before2MTime := mdbMxZFBZMData.CreateTime.AddDate(0, -2, 1-day).Format("200601")
			for _, v := range mdbMxZFBZMData.MXZFBZMData {
				if len(v.Time) > 5 {
					compareTime := v.Time[:6]
					if compareTime == startTime || compareTime == before1MTime || compareTime == before2MTime {
						nearly3MonthZM = append(nearly3MonthZM, v)
					}
				}
			}
		}
	}
	this.Data["isZFBShow"] = isZFBShow
	this.Data["isZMShow"] = isZMShow
	this.Data["zfbZMData"] = nearly3MonthZM
	this.Data["data"] = mdbMxZFBData
	this.Data["uid"] = uid
	this.TplName = "user/user_zfb_info.html"
}

//获取支付宝交易记录
func (this *UsersMetadataController) GetUsersZFBTradeInfo() {
	resultMap := make(map[string]interface{})
	resultMap["ret"] = 403
	defer func() {
		this.Data["json"] = resultMap
		this.ServeJSON()
	}()
	uid, _ := this.GetInt("uid")
	if uid <= 0 {
		cache.RecordLogs(this.User.Id, 0, this.User.Name, this.User.DisplayName, "mgodb数据获取mxalipaydata异常", "信审工作平台/GetUsersZFBTradeinfo", "", this.Ctx.Input)
		resultMap["err"] = "mgodb数据获取mxalipaydata异常"
		return

	}
	//从Mongodb上获取
	var mdbMxZFBData models.MdbMxZFBData
	uidMap := make(map[string]interface{})
	uidMap["uid"] = uid
	session := utils.GetSession()
	defer session.Close()
	err := session.DB(utils.MGO_DB).C("mxalipaydata").Find(uidMap).Sort("-createtime").Select(bson.M{"mxzfbdata.tradeinfo": 1}).Limit(1).One(&mdbMxZFBData)
	if err != nil && err.Error() != utils.MongdbErrNoRow() {
		cache.RecordLogs(this.User.Id, 0, this.User.Name, this.User.DisplayName, "mgodb数据获取mxalipaydata异常", "信审工作平台/GetUsersZFBRecenttradersInfo", err.Error(), this.Ctx.Input)
		resultMap["err"] = "mgodb数据获取mxalipaydata异常" + err.Error()
		return

	}
	page, _ := this.GetInt("page")
	if page < 1 {
		page = 1
	}
	count := len(mdbMxZFBData.MXZFBData.Tradeinfo)
	reqTradeCount := utils.StartIndex(page, utils.PageSize20) + utils.PageSize20
	reqTradeStart := utils.StartIndex(page, utils.PageSize20)
	if len(mdbMxZFBData.MXZFBData.Tradeinfo) >= reqTradeCount {
		mdbMxZFBData.MXZFBData.Tradeinfo = mdbMxZFBData.MXZFBData.Tradeinfo[reqTradeStart:reqTradeCount]
	} else {
		mdbMxZFBData.MXZFBData.Tradeinfo = mdbMxZFBData.MXZFBData.Tradeinfo[reqTradeStart:]
	}
	pageCount := utils.PageCount(count, utils.PageSize20)
	resultMap["ret"] = 200
	resultMap["uid"] = uid
	resultMap["data"] = mdbMxZFBData.MXZFBData.Tradeinfo
	resultMap["pageSize"] = utils.PageSize20
	resultMap["count"] = count
	resultMap["currPage"] = page
	resultMap["pageCount"] = pageCount
}

//获取支付宝联系人信息
func (this *UsersMetadataController) GetUsersZFBAlipayContacts() {
	resultMap := make(map[string]interface{})
	resultMap["ret"] = 403
	defer func() {
		this.Data["json"] = resultMap
		this.ServeJSON()
	}()
	uid, _ := this.GetInt("uid")
	if uid <= 0 {
		cache.RecordLogs(this.User.Id, 0, this.User.Name, this.User.DisplayName, "uid参数传递错误", "信审工作平台/GetUsersZFBAlipayContacts", "", this.Ctx.Input)
		resultMap["err"] = "uid参数传递错误"
		return
	}
	//从Mongodb上获取
	var mdbMxZFBData models.MdbMxZFBData
	uidMap := make(map[string]interface{})
	uidMap["uid"] = uid
	session := utils.GetSession()
	defer session.Close()
	err := session.DB(utils.MGO_DB).C("mxalipaydata").Find(uidMap).Sort("-createtime").Select(bson.M{"mxzfbdata.alipaycontacts": 1}).Limit(1).One(&mdbMxZFBData)
	if err != nil && err.Error() != utils.MongdbErrNoRow() {
		cache.RecordLogs(this.User.Id, 0, this.User.Name, this.User.DisplayName, "mgodb数据获取mxalipaydata异常", "信审工作平台/GetUsersZFBAlipayContacts", err.Error(), this.Ctx.Input)
		resultMap["err"] = "mgodb数据获取mxalipaydata异常" + err.Error()
		return
	}
	page, _ := this.GetInt("page")
	if page < 1 {
		page = 1
	}
	count := len(mdbMxZFBData.MXZFBData.Alipaycontacts)
	reqContactsCount := utils.StartIndex(page, utils.PageSize20) + utils.PageSize20
	reqContactsStart := utils.StartIndex(page, utils.PageSize20)
	if len(mdbMxZFBData.MXZFBData.Alipaycontacts) >= reqContactsCount {
		mdbMxZFBData.MXZFBData.Alipaycontacts = mdbMxZFBData.MXZFBData.Alipaycontacts[reqContactsStart:reqContactsCount]
	} else {
		mdbMxZFBData.MXZFBData.Alipaycontacts = mdbMxZFBData.MXZFBData.Alipaycontacts[reqContactsStart:]
	}
	pageCount := utils.PageCount(count, utils.PageSize20)
	resultMap["ret"] = 200
	resultMap["uid"] = uid
	resultMap["data"] = mdbMxZFBData.MXZFBData.Alipaycontacts
	resultMap["pageSize"] = utils.PageSize20
	resultMap["count"] = count
	resultMap["currPage"] = page
	resultMap["pageCount"] = pageCount
}

//获取支付宝交易人信息
func (this *UsersMetadataController) GetUsersZFBRecenttradersInfo() {
	resultMap := make(map[string]interface{})
	resultMap["ret"] = 403
	defer func() {
		this.Data["json"] = resultMap
		this.ServeJSON()
	}()
	uid, _ := this.GetInt("uid")
	if uid <= 0 {
		cache.RecordLogs(this.User.Id, 0, this.User.Name, this.User.DisplayName, "uid参数传递错误", "信审工作平台/GetUsersZFBRecenttradersInfo", "", this.Ctx.Input)
		resultMap["err"] = "uid参数传递错误"
		return
	}
	//从Mongodb上获取
	var mdbMxZFBData models.MdbMxZFBData
	uidMap := make(map[string]interface{})
	uidMap["uid"] = uid
	session := utils.GetSession()
	defer session.Close()
	err := session.DB(utils.MGO_DB).C("mxalipaydata").Find(uidMap).Sort("-createtime").Select(bson.M{"mxzfbdata.recenttraders": 1}).Limit(1).One(&mdbMxZFBData)
	if err != nil && err.Error() != utils.MongdbErrNoRow() {
		cache.RecordLogs(this.User.Id, 0, this.User.Name, this.User.DisplayName, "mgodb数据获取mxalipaydata异常", "信审工作平台/GetUsersZFBRecenttradersInfo", err.Error(), this.Ctx.Input)
		resultMap["err"] = "mgodb数据获取mxalipaydata异常" + err.Error()
		return
	}
	page, _ := this.GetInt("page")
	if page < 1 {
		page = 1
	}
	count := len(mdbMxZFBData.MXZFBData.Recenttraders)
	reqRecenttradersCount := utils.StartIndex(page, utils.PageSize20) + utils.PageSize20
	reqRecenttradersStart := utils.StartIndex(page, utils.PageSize20)
	if len(mdbMxZFBData.MXZFBData.Recenttraders) >= reqRecenttradersCount {
		mdbMxZFBData.MXZFBData.Recenttraders = mdbMxZFBData.MXZFBData.Recenttraders[reqRecenttradersStart:reqRecenttradersCount]
	} else {
		mdbMxZFBData.MXZFBData.Recenttraders = mdbMxZFBData.MXZFBData.Recenttraders[reqRecenttradersStart:]
	}
	pageCount := utils.PageCount(count, utils.PageSize20)
	resultMap["ret"] = 200
	resultMap["uid"] = uid
	resultMap["data"] = mdbMxZFBData.MXZFBData.Recenttraders
	resultMap["pageSize"] = utils.PageSize20
	resultMap["count"] = count
	resultMap["currPage"] = page
	resultMap["pageCount"] = pageCount
}

//获取登录历史
func (this *UsersMetadataController) GetUsersLoginRecords() {
	page, _ := this.GetInt("page")
	if page < 1 {
		page = 1
	}
	uid, _ := this.GetInt("uid")
	if uid <= 0 {
		cache.RecordLogs(this.User.Id, 0, this.User.Name, this.User.DisplayName, "uid参数传递错误", "信审工作平台/GetUsersLoginRecords", "", this.Ctx.Input)
		this.Abort("uid参数传递错误")
		return
	}
	list, err := models.QueryUsersLoginRecords(uid, utils.StartIndex(page, utils.PageSize15), utils.PageSize15)
	if err != nil {
		cache.RecordLogs(this.User.Id, 0, this.User.Name, this.User.DisplayName, "查询用户登录历史出错", "信审工作平台/GetUsersLoginRecords", err.Error(), this.Ctx.Input)
		this.Abort("查询用户登录历史出错" + err.Error())
		return
	}
	for k, v := range list {
		if v.Province != "" || v.City != "" || v.District != "" || v.Street != "" {
			list[k].IsGPS = true
		}
	}
	count, err := models.QueryUsersLoginCount(uid)
	if err != nil {
		cache.RecordLogs(this.User.Id, 0, this.User.Name, this.User.DisplayName, "查询用户登录历史总数出错", "信审工作平台/GetUsersLoginRecords", err.Error(), this.Ctx.Input)
		this.Abort("查询用户登录历史总数出错" + err.Error())
		return
	}
	isManage, _ := this.GetBool("is_manage", false)
	pageCount := utils.PageCount(count, utils.PageSize15)
	this.Data["data"] = list
	this.Data["uid"] = uid
	this.Data["currPage"] = page
	this.Data["count"] = count
	this.Data["pageSize"] = utils.PageSize15
	this.Data["pageCount"] = pageCount
	this.Data["data"] = list
	if isManage {
		this.TplName = "manage/manage_user_login_records.html"
	} else {
		this.TplName = "user/user_login_records.html"
	}
}

//获取定位列表
func (this *UsersMetadataController) GetAddressList() {
	page, _ := this.GetInt("page")
	if page < 1 {
		page = 1
	}
	uid, _ := this.GetInt("uid")
	if uid <= 0 {
		cache.RecordLogs(this.User.Id, 0, this.User.Name, this.User.DisplayName, "uid参数传递错误", "信审工作平台/GetAddressList", "", this.Ctx.Input)
		this.Abort("uid参数传递错误")
		return
	}
	list, err := models.QueryAddressInfos(uid, utils.StartIndex(page, utils.PageSize15), utils.PageSize15)
	if err != nil {
		cache.RecordLogs(this.User.Id, 0, this.User.Name, this.User.DisplayName, "查询用户定位信息出错", "信审工作平台/GetAddressList", err.Error(), this.Ctx.Input)
		this.Abort("查询用户定位信息出错" + err.Error())
		return
	}
	count, err := models.QueryAddressCount(uid)
	if err != nil {
		cache.RecordLogs(this.User.Id, 0, this.User.Name, this.User.DisplayName, "查询用户定位信息总数出错", "信审工作平台/GetAddressList", err.Error(), this.Ctx.Input)
		this.Abort("查询用户借款记录总数出错" + err.Error())
		return
	}
	for k, v := range list {
		if v.Province != "" || v.City != "" || v.District != "" || v.Street != "" {
			list[k].IsGPS = true
		}
	}
	pageCount := utils.PageCount(count, utils.PageSize15)
	this.Data["data"] = list
	this.Data["uid"] = uid
	this.Data["currPage"] = page
	this.Data["count"] = count
	this.Data["pageSize"] = utils.PageSize15
	this.Data["pageCount"] = pageCount
	this.Data["data"] = list
	this.TplName = "user/user_address_list.html"
}

//认证用户详情的借款记录
func (this *UsersMetadataController) GetUsersManageLoanRecords() {
	page, _ := this.GetInt("page")
	if page < 1 {
		page = 1
	}
	pageSize, _ := this.GetInt("pageSize")
	if pageSize < 1 {
		pageSize = 15
	}
	uid, _ := this.GetInt("uid")
	if uid <= 0 {
		cache.RecordLogs(this.User.Id, 0, this.User.Name, this.User.DisplayName, "uid参数传递错误", "认证用户详情/借款记录GetUsersManageLoanRecords", "", this.Ctx.Input)
		this.Abort("uid参数传递错误")
		return
	}
	return_money, _ := this.GetInt("return_money")
	loanRecords, err := models.QueryUsersManageLoanRecords(uid, utils.StartIndex(page, pageSize), pageSize)
	if err != nil {
		cache.RecordLogs(this.User.Id, 0, this.User.Name, this.User.DisplayName, "查询用户借款记录出错", "认证用户详情/借款记录GetUsersManageLoanRecords", err.Error(), this.Ctx.Input)
		this.Abort("查询用户借款记录出错" + err.Error())
		return
	}
	loanApplyConfig, err := models.GetConfigCache("loan_apply")                //借款协议
	serviceaAgreementConfig, err := models.GetConfigCache("service_agreement") //平台服务协议
	if err != nil && err.Error() != utils.ErrNoRow() {
		cache.RecordLogs(this.User.Id, 0, this.User.Name, this.User.DisplayName, "根据agreement获取项目配置失败", "认证用户详情/借款记录GetUsersManageLoanRecords", err.Error(), this.Ctx.Input)
	}
	for k, v := range loanRecords {
		loanIdStr := strconv.Itoa(v.Id)
		uIdS := strconv.Itoa(v.Uid)
		paramsUrl := "?loanId=" + loanIdStr + "&uid=" + uIdS
		loanRecords[k].LoanAgreementUrl = loanApplyConfig.ConfigValue + paramsUrl
		loanRecords[k].ServiceAgreementUrl = serviceaAgreementConfig.ConfigValue + paramsUrl
	}
	count, err := models.QueryUsersManageLoanRecordCount(uid)
	if err != nil {
		cache.RecordLogs(this.User.Id, 0, this.User.Name, this.User.DisplayName, "查询用户借款记录总数出错", "认证用户详情/借款记录GetUsersManageLoanRecords", err.Error(), this.Ctx.Input)
		this.Abort("查询用户借款记录总数出错" + err.Error())
		return
	}
	pageCount := utils.PageCount(count, pageSize)
	this.Data["uid"] = uid
	this.Data["currPage"] = page
	this.Data["count"] = count
	this.Data["pageSize"] = pageSize
	this.Data["pageCount"] = pageCount
	this.Data["data"] = loanRecords
	this.Data["return_money"] = return_money
	this.TplName = "manage/manage_user_loan_records.html"
}

//还款
func (this *UsersMetadataController) RepaymentPost() {
	resultMap := make(map[string]interface{})
	resultMap["ret"] = 403
	defer func() {
		this.Data["json"] = resultMap
		this.ServeJSON()
	}()
	rsId, _ := this.GetInt("rs_id") //还款计划id
	if rsId <= 0 {
		cache.RecordLogs(this.User.Id, 0, this.User.Name, this.User.DisplayName, "rs_id参数传递错误", "认证用户详情/还款RepaymentPost", "", this.Ctx.Input)
		resultMap["msg"] = "rs_id参数传递错误"
		return
	}
	uid, _ := this.GetInt("uid") //用户id
	if uid <= 0 {
		cache.RecordLogs(this.User.Id, 0, this.User.Name, this.User.DisplayName, "uid参数有误", "认证用户详情/还款RepaymentPost", "", this.Ctx.Input)
		resultMap["msg"] = "uid参数有误"
		return
	}
	channel, _ := this.GetInt("channel")          //还款渠道
	repayMoney, _ := this.GetFloat("repay_money") //还款金额
	if repayMoney <= 0 {
		resultMap["msg"] = "金额有误"
		return
	}
	if channel != 5 && channel != 9 && channel != 8 && channel != 6 {
		resultMap["msg"] = "渠道有误"
		return
	}
	remark := this.GetString("remark")            //备注
	tradeNumber := this.GetString("trade_number") //订单号
	params := map[string]interface{}{"Uid": uid, "RepaymentScheduleId": rsId, "Channel": channel, "OperatorId": this.User.Id, "ReturnMoney": repayMoney, "Remark": remark, "OidPaybill": tradeNumber}
	b, err := services.PostApi(utils.Loan_Repayment, params)
	if err != nil {
		cache.RecordLogs(this.User.Id, 0, this.User.Name, this.User.DisplayName, "还款失败，请求接口失败", "认证用户详情/还款RepaymentPost", err.Error(), this.Ctx.Input)
		resultMap["msg"] = "还款失败，请求接口失败"
		return
	}
	beego.Info(string(b))
	var res models.BaseResponse
	json.Unmarshal(b, &res)
	if res.Ret == 200 {
		adjournHandlingDay, _ := this.GetInt("adjourn_handling_day") //延期处理天数
		err = models.UpdateAdjournHandlingDay(adjournHandlingDay, rsId)
		if err != nil {
			cache.RecordLogs(this.User.Id, 0, this.User.Name, this.User.DisplayName, "修改延期处理天数失败", "认证用户详情/还款RepaymentPost", err.Error(), this.Ctx.Input)
		}
	}
	resultMap["ret"] = res.Ret
	resultMap["msg"] = res.Msg
}

//获取还款链接
func (this *UsersMetadataController) RepaymentHref() {
	resultMap := make(map[string]interface{})
	resultMap["ret"] = 403
	defer func() {
		this.Data["json"] = resultMap
		this.ServeJSON()
	}()
	uid, _ := this.GetInt("uid")
	if uid <= 0 {
		cache.RecordLogs(this.User.Id, 0, this.User.Name, this.User.DisplayName, "uid参数传递错误", "认证用户详情/借款记录RepaymentHref", "", this.Ctx.Input)
		resultMap["err"] = "uid参数传递错误"
		return
	}
	repaymentScheduleId, _ := this.GetInt("repaymentScheduleId") //还款计划id
	loanId, _ := this.GetInt("loanId")                           //借款id
	urlShort, err := utils.SinaDispose("connectapi", uid, loanId, repaymentScheduleId)
	if err != nil {
		cache.RecordLogs(this.User.Id, 0, this.User.Name, this.User.DisplayName, "获取还款链接失败", "认证用户详情/借款记录RepaymentHref", err.Error(), this.Ctx.Input)
		resultMap["err"] = "获取还款链接失败" + err.Error()
		return
	}
	resultMap["ret"] = 200
	resultMap["hrefLink"] = urlShort
}

//获取借款记录
func (this *UsersMetadataController) GetUsersLoanRecords() {
	uid, _ := this.GetInt("uid")
	if uid <= 0 {
		cache.RecordLogs(this.User.Id, 0, this.User.Name, this.User.DisplayName, "uid参数传递错误", "信审工作平台/GetUsersLoanRecords", "", this.Ctx.Input)
		this.Abort("uid参数传递错误")
		return
	}
	loanRecords, err := models.QueryUsersLoanRecords(uid)
	if err != nil {
		cache.RecordLogs(this.User.Id, 0, this.User.Name, this.User.DisplayName, "查询用户借款记录出错", "信审工作平台/GetUsersLoanRecords", err.Error(), this.Ctx.Input)
		this.Abort("查询用户借款记录出错" + err.Error())
		return
	}
	this.Data["loanRecords"] = loanRecords
	this.TplName = "user/user_loan_repayment.html"
}

//获取用户还款记录
func (this *UsersMetadataController) GetUsersLoanRepaymentRecords() {
	resultMap := make(map[string]interface{})
	resultMap["ret"] = 403
	defer func() {
		this.Data["json"] = resultMap
		this.ServeJSON()
	}()
	loanId, _ := this.GetInt("loan_id")
	if loanId <= 0 {
		cache.RecordLogs(this.User.Id, 0, this.User.Name, this.User.DisplayName, "loan_id参数传递错误", "信审工作平台/GetUsersLoanRepaymentRecords", "", this.Ctx.Input)
		resultMap["err"] = "loan_id参数传递错误"
		return
	}
	paymentRecords, err := models.QueryPaymentRecords(loanId)
	if err != nil {
		cache.RecordLogs(this.User.Id, 0, this.User.Name, this.User.DisplayName, "查询用户还款记录出错", "信审工作平台/GetUsersLoanRepaymentRecords", err.Error(), this.Ctx.Input)
		resultMap["err"] = "查询用户还款记录出错" + err.Error()
		return
	}
	for k, v := range paymentRecords {
		switch v.Channel {
		case 1:
			paymentRecords[k].PaymentType = "连连支付"
		case 2:
			paymentRecords[k].PaymentType = "支付宝支付"
		case 3:
			paymentRecords[k].PaymentType = "合利宝代扣"
		case 5:
			paymentRecords[k].PaymentType = "支付宝转账"
		case 6:
			paymentRecords[k].PaymentType = "银行转账"
		case 9:
			paymentRecords[k].PaymentType = "费用减免"
		case 7:
			paymentRecords[k].PaymentType = "融宝代扣"
		case 8:
			paymentRecords[k].PaymentType = "畅捷代扣"
		case 10:
			paymentRecords[k].PaymentType = "先锋代扣"
		}
		if v.Days > 0 {
			paymentRecords[k].PaymentState = "逾期" + strconv.Itoa(v.Days) + "天还款"
		} else if v.Days < 0 {
			paymentRecords[k].PaymentState = "提前" + strconv.Itoa(-1*v.Days) + "天还款"
		} else {
			paymentRecords[k].PaymentState = "到期还款"
		}
	}
	resultMap["ret"] = 200
	resultMap["data"] = paymentRecords
}

//获取认证用户详情还款记录
func (this *UsersMetadataController) GetUsersRepaymentRecords() {
	page, _ := this.GetInt("page")
	if page < 1 {
		page = 1
	}
	pageSize, _ := this.GetInt("pageSize")
	if pageSize < 1 {
		pageSize = 15
	}
	uid, _ := this.GetInt("uid")
	if uid <= 0 {
		cache.RecordLogs(this.User.Id, 0, this.User.Name, this.User.DisplayName, "uid参数传递错误", "认证用户详情/还款记录GetUsersManageLoanRepaymentRecords", "", this.Ctx.Input)
		this.Abort("uid参数传递错误")
		return
	}
	loanRecords, err := models.QueryUsersManagePaymentRecords(uid, utils.StartIndex(page, pageSize), pageSize)
	if err != nil {
		cache.RecordLogs(this.User.Id, 0, this.User.Name, this.User.DisplayName, "查询用户还款记录出错", "认证用户详情/还款记录GetUsersManageLoanRepaymentRecords", err.Error(), this.Ctx.Input)
		this.Abort("查询用户借款记录出错" + err.Error())
		return
	}
	count, err := models.QueryUsersManagePaymentRecordCount(uid)
	if err != nil {
		cache.RecordLogs(this.User.Id, 0, this.User.Name, this.User.DisplayName, "查询用户还款记录总数出错", "认证用户详情/还款记录GetUsersManageLoanRepaymentRecords", err.Error(), this.Ctx.Input)
		this.Abort("查询用户还款记录总数出错" + err.Error())
		return
	}
	pageCount := utils.PageCount(count, pageSize)
	this.Data["uid"] = uid
	this.Data["currPage"] = page
	this.Data["count"] = count
	this.Data["pageSize"] = pageSize
	this.Data["pageCount"] = pageCount
	this.Data["data"] = loanRecords
	this.TplName = "user/user_repayment_record.html"
}

//获取用户还款计划
func (this *UsersMetadataController) GetUsersRepaymentSchedules() {
	resultMap := make(map[string]interface{})
	resultMap["ret"] = 403
	defer func() {
		this.Data["json"] = resultMap
		this.ServeJSON()
	}()
	loanId, _ := this.GetInt("loan_id")
	if loanId <= 0 {
		cache.RecordLogs(this.User.Id, 0, this.User.Name, this.User.DisplayName, "loan_id参数传递错误", "信审工作平台/GetUsersRepaymentSchedules", "", this.Ctx.Input)
		resultMap["err"] = "loan_id参数传递错误"
		return
	}
	paymentSchedules, err := models.QueryPaymentSchedules(loanId)
	if err != nil {
		cache.RecordLogs(this.User.Id, 0, this.User.Name, this.User.DisplayName, "查询用户还款计划出错", "信审工作平台/GetUsersRepaymentSchedules", err.Error(), this.Ctx.Input)
		resultMap["err"] = "查询用户还款计划出错" + err.Error()
		return
	}
	for k, v := range paymentSchedules {
		switch v.State {
		case "BACKING":
			paymentSchedules[k].State = "未还"
		case "OVERDUE":
			paymentSchedules[k].State = "逾期"
		case "FINISH":
			paymentSchedules[k].State = "结清"
		}
	}
	resultMap["ret"] = 200
	resultMap["data"] = paymentSchedules
}

//获取催收记录
func (this *UsersMetadataController) GetCollectionRecord() {
	this.IsNeedTemplate()
	this.TplName = "user/user_collection_record.html"
}

// 联系历史
func (this *UsersMetadataController) GetConnect() {
	this.IsNeedTemplate()
	page, _ := this.GetInt("page")
	if page < 1 {
		page = 1
	}
	pageSize, _ := this.GetInt("pageSize")
	if pageSize < 1 {
		pageSize = 15
	}
	uid, err := this.GetInt("uid")
	connectType := this.GetString("connectType")
	if err != nil {
		cache.RecordLogs(this.User.Id, 0, this.User.Name, this.User.DisplayName, "uid参数传递错误", "认证用户详情/GetConnect", "", this.Ctx.Input)
		this.Abort("uid参数传递错误")
	}
	connect, err := models.GetConnectRecordByUid(uid, connectType, utils.StartIndex(page, pageSize), pageSize)
	if err != nil {
		cache.RecordLogs(this.User.Id, 0, this.User.Name, this.User.DisplayName, "联系历史获取失败", "认证用户详情/GetConnect", "", this.Ctx.Input)
		this.Abort("联系历史获取失败")
	}
	count, err := models.GetConnectRecordCountByUid(uid, connectType)
	if err != nil {
		cache.RecordLogs(this.User.Id, 0, this.User.Name, this.User.DisplayName, "联系历史count获取失败", "认证用户详情/GetConnect", "", this.Ctx.Input)
		this.Abort("联系历史count获取失败")
	}
	pageCount := utils.PageCount(count, pageSize)
	this.Data["uid"] = uid
	this.Data["pageSize"] = pageSize
	this.Data["currPage"] = page
	this.Data["pageCount"] = pageCount
	this.Data["count"] = count
	this.Data["ConnectRecord"] = connect
	this.Data["ConnectType"] = connectType
	this.TplName = "user/user_connect_history.html"
}

//增加联系历史
func (this *UsersMetadataController) AddConnect() {
	defer this.ServeJSON()
	choose := this.GetString("choose")
	var connectObj string
	if choose != "" {
		connectObj = choose
	} else {
		man := this.GetString("man")
		connectObj = man
	}
	uid, err := this.GetInt("uid")
	if err != nil {
		cache.RecordLogs(this.User.Id, 0, this.User.Name, this.User.DisplayName, "uid参数传递错误", "认证用户详情/GetConnect", "", this.Ctx.Input)
		this.Data["json"] = map[string]interface{}{"ret": 304, "err": err.Error(), "msg": "uid参数传递错误"}
		return
	}
	context := this.GetString("content")
	connectType := this.GetString("type")
	if connectType == "COMPLANIT" {
		rcd := &models.Conn_record{Uid: uid, Conn_type: connectType, Content: context, Created_by: this.User.Id, DealStatus: 1}
		err = rcd.Insert() //增加投诉处理
		if err != nil {
			cache.RecordLogs(this.User.Id, 0, this.User.Name, this.User.DisplayName, "增加投诉处理历史失败", "催收管理/Handle", err.Error(), this.Ctx.Input)
			this.Data["json"] = map[string]interface{}{"ret": 304, "err": err.Error(), "msg": "增加投诉处理历史失败"}
			return
		} else {
			this.Data["json"] = map[string]interface{}{"ret": 200, "msg": "增加投诉处理历史成功"}
			return
		}
	} else {
		// var connectRecord models.CollectionRecord
		m := &models.ConnectRecord{Uid: uid, Context: context, ConnectType: connectType, Operator: this.User.DisplayName, CreateTime: time.Now().Local(), ConnectObj: connectObj}
		if id, err := models.AddConnectRecord(m); err == nil && id != 0 {
			this.Data["json"] = map[string]interface{}{"ret": 200, "msg": "增加联系成功"}
			return
		} else {
			cache.RecordLogs(this.User.Id, 0, this.User.Name, this.User.DisplayName, "增加联系历史失败", "认证用户详情/GetConnect", "", this.Ctx.Input)
			this.Data["json"] = map[string]interface{}{"ret": 304, "err": err.Error(), "msg": "增加联系历史失败"}
			return
		}
	}

}

//跳转到用户消息中心
func (this *UsersMetadataController) JumpToMessageCenter() {
	this.TplName = "user/user_message_center.html"
}

//获取用户消息
func (this *UsersMetadataController) GetUsersMessages() {
	resultMap := make(map[string]interface{})
	resultMap["ret"] = 403
	defer func() {
		this.Data["json"] = resultMap
		this.ServeJSON()
	}()
	uid, _ := this.GetInt("uid")
	if uid <= 0 {
		cache.RecordLogs(this.User.Id, 0, this.User.Name, this.User.DisplayName, "uid参数传递错误", "认证用户详情/消息中心GetUsersMessages", "", this.Ctx.Input)
		resultMap["err"] = "uid参数传递错误"
		return
	}
	//从Mongodb上获取
	var usersMessages []models.UsersMessage
	uidMap := make(map[string]interface{})
	uidMap["uid"] = uid
	session := utils.GetSession()
	defer session.Close()
	err := session.DB(utils.MGO_DB).C("users_message").Find(uidMap).Sort("-createhidetime").All(&usersMessages)
	if err != nil {
		cache.RecordLogs(this.User.Id, 0, this.User.Name, this.User.DisplayName, "mgodb数据获取users_message异常", "认证用户详情/消息中心GetUsersMessages", err.Error(), this.Ctx.Input)
		resultMap["uid"] = "mgodb数据获取users_message异常" + err.Error()
		return
	}
	page, _ := this.GetInt("page")
	if page < 1 {
		page = 1
	}
	count := len(usersMessages)
	reqRecenttradersCount := utils.StartIndex(page, utils.PageSize5) + utils.PageSize5
	reqRecenttradersStart := utils.StartIndex(page, utils.PageSize5)
	if len(usersMessages) >= reqRecenttradersCount {
		usersMessages = usersMessages[reqRecenttradersStart:reqRecenttradersCount]
	} else {
		usersMessages = usersMessages[reqRecenttradersStart:]
	}
	pageCount := utils.PageCount(count, utils.PageSize5)
	resultMap["ret"] = 200
	resultMap["data"] = usersMessages
	resultMap["pageSize"] = utils.PageSize5
	resultMap["currPage"] = page
	resultMap["pageCount"] = pageCount
}

//跳转到用户优惠券
func (this *UsersMetadataController) JumpToCoupons() {
	this.TplName = "user/user_coupon.html"
}

//获取用户优惠券
func (this *UsersMetadataController) GetUsersCoupons() {
	resultMap := make(map[string]interface{})
	resultMap["ret"] = 403
	defer func() {
		this.Data["json"] = resultMap
		this.ServeJSON()
	}()
	uid, _ := this.GetInt("uid")
	if uid <= 0 {
		cache.RecordLogs(this.User.Id, 0, this.User.Name, this.User.DisplayName, "uid参数传递错误", "认证用户详情/优惠券GetUsersCoupons", "", this.Ctx.Input)
		resultMap["err"] = "uid参数传递错误"
		return
	}
	page, _ := this.GetInt("page")
	if page < 1 {
		page = 1
	}
	usersCoupons, err := models.QueryUsersCoupons(uid, utils.StartIndex(page, utils.PageSize4), utils.PageSize4)
	if err != nil {
		cache.RecordLogs(this.User.Id, 0, this.User.Name, this.User.DisplayName, "查询用户优惠券出错", "认证用户详情/优惠券GetUsersCoupons", err.Error(), this.Ctx.Input)
		this.Abort("查询用户邀请记录总数出错" + err.Error())
		return
	}
	for k, v := range usersCoupons {
		if v.IsUsed == "DISABLE" {
			usersCoupons[k].IsUsed = "已使用"
		} else if usersCoupons[k].IsUsed == "USABLE" {
			usersCoupons[k].IsUsed = "可用的"
		} else if usersCoupons[k].IsUsed == "PAST" {
			usersCoupons[k].IsUsed = "已过期"
		} else if usersCoupons[k].IsUsed == "UNBEGIN" {
			usersCoupons[k].IsUsed = "未开始"
		}
	}
	count, err := models.QueryUsersCouponsCount(uid)
	if err != nil {
		cache.RecordLogs(this.User.Id, 0, this.User.Name, this.User.DisplayName, "查询用户优惠券总数出错", "认证用户详情/优惠券GetUsersCoupons", err.Error(), this.Ctx.Input)
		resultMap["err"] = "查询用户优惠券总数出错" + err.Error()
		return
	}
	pageCount := utils.PageCount(count, utils.PageSize4)
	resultMap["ret"] = 200
	resultMap["data"] = usersCoupons
	resultMap["pageSize"] = utils.PageSize4
	resultMap["currPage"] = page
	resultMap["pageCount"] = pageCount

}

//额度历史
func (this *UsersMetadataController) GetLimit() {
	this.IsNeedTemplate()
	page, _ := this.GetInt("page")
	if page < 1 {
		page = 1
	}
	uid, err := this.GetInt("uid")
	if err != nil {
		cache.RecordLogs(this.User.Id, 0, this.User.Name, this.User.DisplayName, "uid参数失败", "认证用户详情/GetLimit", err.Error(), this.Ctx.Input)
		this.Abort("uid参数失败" + err.Error())
		return
	}
	limit, err := models.GetLimitRecordByUid(uid, utils.StartIndex(page, utils.PageSize15), utils.PageSize15)
	if err != nil {
		cache.RecordLogs(this.User.Id, 0, this.User.Name, this.User.DisplayName, "limit获取失败", "认证用户详情/GetLimit", err.Error(), this.Ctx.Input)
		this.Abort("limit获取失败" + err.Error())
		return
	}
	count, err := models.GetLimitRecordCount(uid)
	if err != nil {
		cache.RecordLogs(this.User.Id, 0, this.User.Name, this.User.DisplayName, "limit获取失败", "认证用户详情/GetLimit", err.Error(), this.Ctx.Input)
		this.Abort("limit获取失败" + err.Error())
		return
	}
	pageCount := utils.StartIndex(page, utils.PageSize15)
	this.Data["uid"] = uid
	this.Data["page"] = page
	this.Data["pageSize"] = utils.PageSize15
	this.Data["currPage"] = page
	this.Data["pageCount"] = pageCount
	this.Data["count"] = count
	this.Data["Limit"] = limit
	this.TplName = "user/user_limit_history.html"
}

//获取用户邀请记录
func (this *UsersMetadataController) GetUsersInvitationRecords() {
	page, _ := this.GetInt("page")
	if page < 1 {
		page = 1
	}
	uid, _ := this.GetInt("uid")
	if uid <= 0 {
		cache.RecordLogs(this.User.Id, 0, this.User.Name, this.User.DisplayName, "uid参数传递错误", "认证用户详情/邀请记录GetUsersInvitationRecords", "", this.Ctx.Input)
		this.Abort("uid参数传递错误")
		return
	}
	list, err := models.QueryInvitationRecords(uid, utils.StartIndex(page, utils.PageSize15), utils.PageSize15)
	if err != nil {
		cache.RecordLogs(this.User.Id, 0, this.User.Name, this.User.DisplayName, "查询用户邀请记录出错", "认证用户详情/邀请记录GetUsersInvitationRecords", err.Error(), this.Ctx.Input)
		this.Abort("查询用户邀请记录出错" + err.Error())
		return
	}
	count, err := models.QueryInvitationRecordCount(uid)
	if err != nil {
		cache.RecordLogs(this.User.Id, 0, this.User.Name, this.User.DisplayName, "查询用户邀请记录总数出错", "认证用户详情/邀请记录GetUsersInvitationRecords", err.Error(), this.Ctx.Input)
		this.Abort("查询用户邀请记录总数出错" + err.Error())
		return
	}
	pageCount := utils.PageCount(count, utils.PageSize15)
	this.Data["uid"] = uid
	this.Data["currPage"] = page
	this.Data["count"] = count
	this.Data["pageSize"] = utils.PageSize15
	this.Data["pageCount"] = pageCount
	this.Data["data"] = list
	this.TplName = "user/user_invite_records.html"
}

//修改用户备注
func (this *UsersMetadataController) UpdateUsersSign() {
	resultMap := make(map[string]interface{})
	defer func() {
		this.Data["json"] = resultMap
		this.ServeJSON()
	}()
	resultMap["ret"] = 304
	uid, _ := this.GetInt("uid")
	if uid <= 0 {
		cache.RecordLogs(this.User.Id, 0, this.User.Name, this.User.DisplayName, "uid参数传递错误", "认证用户详情/UpdateUsersSign", "", this.Ctx.Input)
		resultMap["err"] = "uid参数传递错误"
		return
	}

	content := this.GetString("content")
	err := models.UpdateUsersSign(uid, content)
	if err != nil {
		cache.RecordLogs(this.User.Id, 0, this.User.Name, this.User.DisplayName, "修改用户备注出错", "认证用户详情/UpdateUsersSign", err.Error(), this.Ctx.Input)
		resultMap["err"] = "修改用户备注出错" + err.Error()
		return
	}
	resultMap["ret"] = 200
	resultMap["msg"] = "修改备注成功"
}

//跳转到用户银行卡
func (this *UsersMetadataController) JumpToBankcards() {
	this.TplName = "user/user_bankcard.html"
}

//获取用户银行卡
func (this *UsersMetadataController) GetUsersBankcards() {
	resultMap := make(map[string]interface{})
	resultMap["ret"] = 403
	defer func() {
		this.Data["json"] = resultMap
		this.ServeJSON()
	}()
	page, _ := this.GetInt("page")
	if page < 1 {
		page = 1
	}
	uid, _ := this.GetInt("uid")
	if uid <= 0 {
		cache.RecordLogs(this.User.Id, 0, this.User.Name, this.User.DisplayName, "uid参数传递错误", "认证用户详情/银行卡GetUsersBankcards", "", this.Ctx.Input)
		resultMap["err"] = "uid参数传递错误"
		return
	}
	isBind, err := models.QueryUsersBankcardBind(uid)
	if err != nil && err.Error() != utils.ErrNoRow() {
		cache.RecordLogs(this.User.Id, 0, this.User.Name, this.User.DisplayName, "查询用户银行卡绑定情况出错", "认证用户详情/银行卡GetUsersBankcards", err.Error(), this.Ctx.Input)
		resultMap["err"] = "查询用户银行卡绑定情况出错" + err.Error()
		return
	}
	if page != 1 { //非第一页都为未绑卡状态
		isBind = 1
	}
	list, err := models.QueryUsersBankcards(uid, utils.StartIndex(page, utils.PageSize4), utils.PageSize4)
	if err != nil {
		cache.RecordLogs(this.User.Id, 0, this.User.Name, this.User.DisplayName, "查询用户银行卡出错", "认证用户详情/银行卡GetUsersBankcards", err.Error(), this.Ctx.Input)
		resultMap["err"] = "查询用户银行卡出错" + err.Error()
		return
	}
	for k, v := range list {
		list[k].CardNumber = utils.Substr(v.CardNumber, -3, len(v.CardNumber))
		list[k].BankMobile = utils.MobileFilter(v.BankMobile)
		if len(v.Verifyrealname) >= 24 {
			list[k].Verifyrealname = "**" + v.Verifyrealname[6:]
		} else {
			list[k].Verifyrealname = "*" + v.Verifyrealname[3:]
		}
	}
	count, err := models.QueryUsersBankcardsCount(uid)
	if err != nil {
		cache.RecordLogs(this.User.Id, 0, this.User.Name, this.User.DisplayName, "查询用户银行卡总数出错", "认证用户详情/银行卡GetUsersBankcards", err.Error(), this.Ctx.Input)
		resultMap["err"] = "查询用户银行卡总数出错" + err.Error()
		return
	}
	pageCount := utils.PageCount(count, utils.PageSize4)
	resultMap["isBind"] = isBind
	resultMap["ret"] = 200
	resultMap["data"] = list
	resultMap["pageSize"] = utils.PageSize4
	resultMap["currPage"] = page
	resultMap["pageCount"] = pageCount
}

//解绑银行卡
func (this *UsersMetadataController) UnwrapBankcard() {
	resultMap := make(map[string]interface{})
	resultMap["ret"] = 403
	defer func() {
		this.Data["json"] = resultMap
		this.ServeJSON()
	}()
	uid, _ := this.GetInt("uid")
	if uid <= 0 {
		cache.RecordLogs(this.User.Id, 0, this.User.Name, this.User.DisplayName, "uid参数传递错误", "认证用户详情/银��卡UnwrapBankcard", "", this.Ctx.Input)
		resultMap["err"] = "uid参数传递错误"
		return
	}
	isBind, err := models.QueryUsersBankcardBind(uid)
	if err != nil && err.Error() != utils.ErrNoRow() {
		cache.RecordLogs(this.User.Id, 0, this.User.Name, this.User.DisplayName, "查询用户银行卡绑定情况出错", "认证用户详情/银行卡GetUsersBankcards", err.Error(), this.Ctx.Input)
		resultMap["err"] = "查询用户银行卡绑定情况出错" + err.Error()
		return
	}
	if isBind == 1 {
		resultMap["err"] = "该用户还未绑定银行卡"
		return
	}
	err = models.UpdateUsersBankcardBind(uid, 1)
	if err != nil {
		cache.RecordLogs(this.User.Id, 0, this.User.Name, this.User.DisplayName, "解绑银行卡出错", "认证用户详情/银行卡UnwrapBankcard", err.Error(), this.Ctx.Input)
		resultMap["err"] = "解绑银行卡出错" + err.Error()
		return
	}
	userCard, err := models.QueryUsersBankcardByUid(uid)
	if err != nil && err.Error() != utils.ErrNoRow() {
		cache.RecordLogs(this.User.Id, 0, this.User.Name, this.User.DisplayName, "查询用户绑定银行卡出错", "认证用户详情/银行卡UnwrapBankcard", err.Error(), this.Ctx.Input)
		resultMap["err"] = "查询用户绑定银行卡出错" + err.Error()
		return
	}
	err = models.UpdateUsersBankcard(uid)
	if err != nil {
		cache.RecordLogs(this.User.Id, 0, this.User.Name, this.User.DisplayName, "更新用户绑定银行卡出错", "认证用户详情/银行卡UnwrapBankcard", err.Error(), this.Ctx.Input)
		resultMap["err"] = "更新用户绑定银行卡出错" + err.Error()
		return
	}
	//删缓存
	if utils.Re == nil && utils.Rc.IsExist(utils.CACHE_KEY_USER_BANKCARD+strconv.Itoa(uid)) {
		err := utils.Rc.Delete(utils.CACHE_KEY_USER_BANKCARD + strconv.Itoa(uid))
		if err != nil {
			beego.Info("delete cache bandcard fail:", err)
		}
	}
	err = models.AddUsersBankcardLog(uid, userCard.Id, this.User.Id, userCard.CardNumber, this.User.DisplayName)
	if err != nil {
		cache.RecordLogs(this.User.Id, 0, this.User.Name, this.User.DisplayName, "添加解绑日志错误", "认证用户详情/银行卡UnwrapBankcard", err.Error(), this.Ctx.Input)
	}
	resultMap["ret"] = 200
	resultMap["msg"] = "解绑成功"
}

//获取用户发放记录
func (this *UsersMetadataController) GetTradeRecordWithLoanId() {
	resultMap := make(map[string]interface{})
	defer func() {
		this.Data["json"] = resultMap
		this.ServeJSON()
	}()
	resultMap["ret"] = 304
	loanId, _ := this.GetInt("loan_id")
	if loanId <= 0 {
		cache.RecordLogs(this.User.Id, 0, this.User.Name, this.User.DisplayName, "loan_id参数传递错误", "认证用户详情/GetTradeRecordWithLoanId", "", this.Ctx.Input)
		resultMap["err"] = "loan_id参数传递错误"
		return
	}
	tradeRecord, err := models.QueryTradeRecord(loanId)
	if err != nil && err.Error() != utils.ErrNoRow() {
		cache.RecordLogs(this.User.Id, 0, this.User.Name, this.User.DisplayName, "查询用户发放记录错误", "认证用户详情/GetTradeRecordWithLoanId", err.Error(), this.Ctx.Input)
		resultMap["err"] = "查询用户发放记录错误" + err.Error()
		return
	}
	if tradeRecord.State == "SUCCESS" {
		tradeRecord.State = "成功"
	} else if tradeRecord.State == "CONFIRM" {
		tradeRecord.State = "确认中"
	} else if tradeRecord.State == "FAIL" {
		tradeRecord.State = "失败"
	}
	resultMap["ret"] = 200
	resultMap["data"] = tradeRecord
}

//关闭订单
func (this *UsersMetadataController) ClosePaymentRecordOrder() {
	resultMap := make(map[string]interface{})
	resultMap["ret"] = 403
	defer func() {
		this.Data["json"] = resultMap
		this.ServeJSON()
	}()
	prId, _ := this.GetInt("pr_id")
	if prId <= 0 {
		cache.RecordLogs(this.User.Id, 0, this.User.Name, this.User.DisplayName, "pr_id参数传递错误", "认证用户详情/还款记录ClosePaymentRecordOrder", "", this.Ctx.Input)
		resultMap["err"] = "pr_id参数传递错误"
		return
	}
	err := models.UpdatePaymentRecordState(prId, "FAIL")
	if err != nil {
		cache.RecordLogs(this.User.Id, 0, this.User.Name, this.User.DisplayName, "关闭订单出错", "认证用户详情/还款记录ClosePaymentRecordOrder", err.Error(), this.Ctx.Input)
		this.Abort("关闭订单出错" + err.Error())
		return
	}
	resultMap["ret"] = 200
	resultMap["msg"] = "关闭订单成功"
}

//查看个人征信报告
func (this *UsersMetadataController) BuyRecordLook() {
	pageNum, _ := this.GetInt("page", 1) //分页信息（第几页）
	if pageNum < 1 {
		pageNum = 1
	}
	pageSize, _ := this.GetInt("pageSize")
	if pageSize < 1 {
		pageSize = 30
	}
	uid, _ := this.GetInt("uid")
	per_buy_list, err := models.GetPerCreditReport(utils.StartIndex(pageNum, pageSize), pageSize, uid)
	if err != nil {
		cache.RecordLogs(this.User.Id, 0, this.User.Name, this.User.DisplayName, "财务管理", "财务管理/还款管理list", err.Error(), this.Ctx.Input)
		beego.Info(err)
		this.Abort("查询还款管理失败")
		return
	}
	count, err := models.GetPerUsersCreditReportCount(uid)
	if err != nil {
		beego.Info(err)
		cache.RecordLogs(this.User.Id, 0, this.User.Name, this.User.DisplayName, "财务管理", "财务管理/还款管理list", err.Error(), this.Ctx.Input)
		this.Abort("获取还款管理总数异常")
		return
	}
	for k, v := range per_buy_list {
		if v.State == "SUCCESS" || v.State == "LOCKED" || (v.State == "USELESS" && v.LoanId != 0) {
			per_buy_list[k].StateStr = "成功"
		} else if v.State == "FAIL" || (v.State == "USELESS" && v.LoanId == 0) {
			per_buy_list[k].StateStr = "失败"
		} else if v.State == "CONFIRM" {
			per_buy_list[k].StateStr = "处理中"
		}
	}
	pageCount := utils.PageCount(count, pageSize)
	this.Data["uid"] = uid
	this.Data["list"] = per_buy_list
	this.Data["currPage"] = pageNum
	this.Data["count"] = count
	this.Data["pageSize"] = pageSize
	this.Data["pageCount"] = pageCount
	this.TplName = "user/user_report_buy_history.html"
}

//获取学信网信息
func (this *UsersMetadataController) GetUserXxwInfo() {
	status := "noAuth"
	defer func() {
		this.Data["status"] = status
		this.TplName = "user/user_xuexin.html"
	}()
	uid, _ := this.GetInt("uid")
	if uid <= 0 {
		cache.RecordLogs(this.User.Id, 0, this.User.Name, this.User.DisplayName, "uid参数传递错误", "认证用户详情/学信网信息GetGetUserXxwInfo", "", this.Ctx.Input)
		this.Abort("uid参数传递错误")
	}
	xxwIden, _ := models.QueryXxwIden(uid)
	//从Mongodb上获取
	var mgoMxXxwData models.MdbMxXXWData
	var mdbPyInfo models.MdbPyInfo
	session := utils.GetSession()
	defer session.Close()
	if xxwIden.Is_xxw_auth == 2 {
		err := session.DB(utils.MGO_DB).C("mxxxwdata").Find(&bson.M{"_id": bson.ObjectIdHex(xxwIden.Xxw_mgo_id)}).One(&mgoMxXxwData)
		if err != nil && err.Error() != utils.MongdbErrNoRow() {
			cache.RecordLogs(this.User.Id, 0, this.User.Name, this.User.DisplayName, "mgodb数据获取mxxxwdata异常", "认证用户详情/学信网信息GetGetUserXxwInfo", err.Error(), this.Ctx.Input)
			this.Abort("mgodb数据获取mxxxwdata异常" + err.Error())
		}
		status = "auth"
		if len(mgoMxXxwData.MXXXWData.StudentInfoList) > 0 {
			this.Data["student_data"] = mgoMxXxwData.MXXXWData.StudentInfoList[0]
		}
		if len(mgoMxXxwData.MXXXWData.EducationList) > 0 {
			this.Data["education_data"] = mgoMxXxwData.MXXXWData.EducationList[0]
		}
	} else if xxwIden.Is_yd_pscore_auth == 2 {
		err := session.DB(utils.MGO_DB).C("yd_personal_score").Find(&bson.M{"_id": bson.ObjectIdHex(xxwIden.Yd_pscore_mgo_id)}).One(&mdbPyInfo)
		if err != nil && err.Error() != utils.MongdbErrNoRow() {
			cache.RecordLogs(this.User.Id, 0, this.User.Name, this.User.DisplayName, "mgodb数据获取yd_personal_score异常", "认证用户详情/学信网信息GetGetUserXxwInfo", err.Error(), this.Ctx.Input)
			this.Abort("mgodb数据获取yd_personal_score异常" + err.Error())
		}
		status = "noXxwData"
		if len(mdbPyInfo.Data.ReturnValue.CisReport) > 0 {
			this.Data["data"] = mdbPyInfo.Data.ReturnValue.CisReport[0].PersonApplyScoreInfo.Score
		}
	}
}
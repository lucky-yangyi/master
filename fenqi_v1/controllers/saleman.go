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

	"fmt"

	"sync"

	"github.com/astaxie/beego"
)

type SalesmanController struct {
	BaseController
}

//业务员管理列表
func (c *SalesmanController) GetSalesmanList() {
	defer c.IsNeedTemplate()
	pars := []interface{}{}
	condition := ""
	pageNum, _ := c.GetInt("page", 1) //分页信息（第几页）
	if pageNum < 1 {
		pageNum = 1
	}
	pageSize, _ := c.GetInt("pageSize")
	if pageSize < 1 {
		pageSize = 30
	}

	region := c.GetString("district")
	if region != "" {
		condition += " AND s.region_id = ?"
		pars = append(pars, region)
	}

	place := c.GetString("province")
	if place != "" {
		condition += " AND s.place_id = ?"
		pars = append(pars, place)
	}

	operaDep := c.GetString("office")
	if operaDep != "" {
		condition += " AND s.dep_id = ?"
		pars = append(pars, operaDep)
	}

	saleman := c.GetString("name")
	if saleman != "" {
		condition += " AND s.saleman like ?"
		pars = append(pars, "%"+saleman+"%")
	}

	account := c.GetString("account")
	if account != "" {
		condition += " AND s.account = ?"
		pars = append(pars, account)
	}
	invite := c.GetString("invite_code")
	if invite != "" {
		condition += " AND s.invite_code like ?"
		pars = append(pars, invite+"%")
	}
	stationName := c.GetString("station_name")
	if stationName != "" {
		condition += " AND sta.name = ?"
		pars = append(pars, stationName)
	}
	createTime := c.GetString("auth_time")
	var startcreateTimeTime, endcreateTimeTime string
	if createTime != "" {
		createTimes := strings.Split(createTime, "~")
		startcreateTimeTime = createTimes[0] + " 00:00:00"
		endcreateTimeTime = createTimes[1] + " 23:59:59"
		pars = append(pars, startcreateTimeTime)
		pars = append(pars, endcreateTimeTime)
		condition += ` AND s.create_time >= ? AND s.create_time <= ?`
	}

	stationId := c.User.StationId
	orgStr, err := cache.GetCacheDataByStation(stationId)
	if err != nil {
		cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "业务员管理根据岗位获取数据权限失败", "业务员管理/GetSalesmanList", err.Error(), c.Ctx.Input)
	}

	salesmanlist, err := models.GetSalesmanLookList(utils.StartIndex(pageNum, pageSize), pageSize, condition, orgStr, pars...)
	if err != nil {
		cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "查询业务员列表失败", "业务员端/GetSalesmanList", err.Error(), c.Ctx.Input)
		c.Abort("查询业务员列表失败")
		return
	}
	s, err := cache.GetSysOrganization()
	salemanMap := make(map[string]models.SysOrganization, len(s))
	for _, v := range s {
		salemanMap[strconv.Itoa(v.Id)] = v
	}
	for _, v := range salesmanlist {
		//DepId 营业部id PlaceId 省级id RegionId 区运营中心id
		if v.RegionId != 0 {
			large_area_name, place_name, bus_dep_name := services.QueryOrgTypeDepartment(salemanMap, v.DepId, v.PlaceId, v.RegionId)
			v.Place = place_name
			v.OperaDep = bus_dep_name
			v.Region = large_area_name
		}
		//省运营中心id
		// } else if v.PlaceId != 0 {
		// 	a, b, _ := services.FindOrgTypePlace(salemanMap, strconv.Itoa(v.PlaceId))
		// 	v.Place = b
		// 	v.Region = a
		// //区运营中心id
		// } else if v.RegionId != 0 {
		// 	_, b, _ := services.FindOrgTypePlace(salemanMap, strconv.Itoa(v.RegionId))
		// 	v.Region = b
		// }
		//小组id
		if v.GroupId != 0 {
			group, _ := models.GetSysOrganizationById(v.GroupId)
			v.Group = group.Name
		}
	}
	count, err := models.GetSalesmanLookCount(condition, orgStr, pars...)
	if err != nil {
		cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "查询业务员总数失败", "业务员端/GetSalesmanList", err.Error(), c.Ctx.Input)
		c.Abort("查询业务员总数失败")
		return
	}
	pageCount := utils.PageCount(count, pageSize)
	//===============数据回填
	var provinces []models.SysOrganization
	var offices []models.SysOrganization
	var regions []models.SysOrganization
	r, _ := strconv.Atoi(region)
	p, _ := strconv.Atoi(place)
	o, _ := strconv.Atoi(operaDep)
	//region
	regions, err = models.GetParentSysOrganizationById(28)
	if err != nil {
		cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "查询业务员总数失败", "业务员端/GetSalesmanList", err.Error(), c.Ctx.Input)
		c.Abort("查询地区失败")
		return
	}

	//province
	if r != 0 {
		provinces, err = models.GetParentSysOrganizationById(r)
		if err != nil {
			cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "查询业务员总数失败", "业务员端/GetSalesmanList", err.Error(), c.Ctx.Input)
			c.Abort("查询省份失败")
			return
		}
	}
	//office
	if p != 0 {
		offices, err = models.GetParentSysOrganizationById(p)
		if err != nil {
			cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "查询业务员总数失败", "业务员端/GetSalesmanList", err.Error(), c.Ctx.Input)
			c.Abort("查询营业部失败")
			return
		}
	}
	//=================end
	c.Data["list"] = salesmanlist
	c.Data["currPage"] = pageNum
	c.Data["count"] = count
	c.Data["pageSize"] = pageSize
	c.Data["pageCount"] = pageCount
	//回填
	c.Data["region"] = r
	c.Data["regions"] = regions
	c.Data["province"] = p
	c.Data["provinces"] = provinces
	c.Data["office"] = o
	c.Data["offices"] = offices
	c.Data["saleman"] = saleman
	c.Data["account"] = account
	c.Data["invite"] = invite
	c.Data["createTime"] = createTime
	c.TplName = "clerk/salesman_list.html"

}

//根据id获取业务员
func (c *SalesmanController) GetSalesmanById() {
	defer c.IsNeedTemplate()
	id, err := c.GetInt("id")
	if err != nil {
		cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "获取id失败", "业务员端/GetSalesmanById", err.Error(), c.Ctx.Input)
		c.Abort("获取id失败")
		return
	}

	saleman, err := models.GetSalemanById(id)
	s, _ := models.GetOrgPid(saleman.DepId)
	if s != nil {
		saleman.OperaDep = s.Name
	}
	s, _ = models.GetOrgPid(saleman.PlaceId)
	if s != nil {
		saleman.Place = s.Name
	}
	s, _ = models.GetOrgPid(saleman.RegionId)
	if s != nil {
		saleman.Region = s.Name
	}
	if err != nil {
		cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "根据id查询业务员失败", "业务员端/GetSalesmanById", err.Error(), c.Ctx.Input)
		c.Abort("根据id查询业务员失败")
		return
	}
	c.Data["saleman"] = saleman
	c.TplName = "clerk/personal_information.html"

}

//已认证用户
func (c *SalesmanController) GetUsersIsAuth() {
	defer c.IsNeedTemplate()
	page, _ := c.GetInt("page")
	if page <= 0 {
		page = 1
	}
	pageSize, _ := c.GetInt("pageSize")
	if pageSize < 1 {
		pageSize = 15
	}
	id, err := c.GetInt("id")
	if err != nil {
		cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "获取id失败", "业务员端/GetUsersIsAuth", err.Error(), c.Ctx.Input)
		return
	}
	users, err := models.GetUsersIsAuth(page, pageSize, id)
	if err != nil {
		cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "获取已认证客户失败", "业务员端/GetUsersIsAuth", err.Error(), c.Ctx.Input)
		return
	}
	for _, v := range users {
		count := models.QueryUserCreditCount(v.Uid)
		if count > 1 {
			if v.CurrentAuthState == 0 {
				v.CurrentAuthState, _ = models.QueryUserLastCredit(v.Uid)
			}
		}
		if v.CurrentAuthState == 1 {
			v.State = "授信通过"
		} else if v.CurrentAuthState == 0 || v.CurrentAuthState == 3 {
			v.State = "未提交"
		} else if v.CurrentAuthState == 4 || v.CurrentAuthState == 5 || v.CurrentAuthState == 6 || v.CurrentAuthState == 10 ||
			v.CurrentAuthState == 11 || v.CurrentAuthState == 12 || v.CurrentAuthState == 13 || v.CurrentAuthState == 14 {
			v.State = "授信驳回"
		} else if v.CurrentAuthState == 8 {
			v.State = "授信关闭30天"
		} else if v.CurrentAuthState == 9 {
			v.State = "授信不通过"
		} else if v.CurrentAuthState == 2 {
			v.State = "授信中"
		}
		if v.RealNameTime.After(v.UsersBaseInfoTime) {
			v.SubmitTime = v.RealNameTime
		} else {
			v.SubmitTime = v.UsersBaseInfoTime
		}
	}
	count, _ := models.GetUsersIsAuthCount(id)
	pageCount := 0
	if count%pageSize == 0 {
		pageCount = (count / pageSize)
	} else {
		pageCount = (count / pageSize) + 1
	}
	c.Data["users"] = users
	c.Data["sid"] = id
	c.Data["page"] = page
	c.Data["pageCount"] = pageCount
	c.Data["pageSize"] = pageSize
	c.Data["count"] = count
	c.TplName = "user/user_certified.html"
}

//未认证用户
func (c *SalesmanController) GetUsersNotIsAuth() {
	defer c.IsNeedTemplate()
	page, _ := c.GetInt("page")
	if page <= 0 {
		page = 1
	}
	pageSize, _ := c.GetInt("pageSize")
	if pageSize < 1 {
		pageSize = 15
	}
	id, err := c.GetInt("id")
	if err != nil {
		cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "获取id失败", "业务员端/GetUsersNotIsAuth", err.Error(), c.Ctx.Input)
		return
	}
	users, err := models.GetUsersNotIsAuth(page, pageSize, id)
	if err != nil {
		beego.Info(err)
		cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "获取已认证客户失败", "业务员端/GetUsersNotIsAuth", err.Error(), c.Ctx.Input)
		return
	}
	for _, v := range users {
		count := models.QueryUserCreditCount(v.Uid)
		if count > 1 {
			if v.CurrentAuthState == 0 {
				v.CurrentAuthState, _ = models.QueryUserLastCredit(v.Uid)
			}
		}
		if v.CurrentAuthState == 1 {
			v.State = "授信通过"
		} else if v.CurrentAuthState == 0 || v.CurrentAuthState == 3 {
			v.State = "未提交"
		} else if v.CurrentAuthState == 4 || v.CurrentAuthState == 5 || v.CurrentAuthState == 6 || v.CurrentAuthState == 10 ||
			v.CurrentAuthState == 11 || v.CurrentAuthState == 12 || v.CurrentAuthState == 13 || v.CurrentAuthState == 14 {
			v.State = "授信驳回"
		} else if v.CurrentAuthState == 8 {
			v.State = "授信关闭30天"
		} else if v.CurrentAuthState == 9 {
			v.State = "授信不通过"
		} else if v.CurrentAuthState == 2 {
			v.State = "授信中"
		}
		if v.RealNameTime.After(v.UsersBaseInfoTime) {
			v.SubmitTime = v.RealNameTime
		} else {
			v.SubmitTime = v.UsersBaseInfoTime
		}
	}
	count, _ := models.GetUsersIsNotAuthCount(id)
	pageCount := 0
	if count%pageSize == 0 {
		pageCount = (count / pageSize)
	} else {
		pageCount = (count / pageSize) + 1
	}
	c.Data["users"] = users
	c.Data["nid"] = id
	c.Data["pageCount"] = pageCount
	c.Data["page"] = page
	c.Data["pageSize"] = pageSize
	c.Data["count"] = count
	c.TplName = "user/user_uncertified.html"

}

//根据id获取业务员
func (c *SalesmanController) GetSalesmanAllotmentById() {
	defer c.IsNeedTemplate()
	id, err := c.GetInt("id")
	if err != nil {
		cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "获取id失败", "业务员端/GetSalesmanById", err.Error(), c.Ctx.Input)
		c.Abort("获取id失败")
		return
	}

	saleman, err := models.GetSalemanById(id)
	s, _ := models.GetOrgPid(saleman.DepId)
	if s != nil {
		saleman.OperaDep = s.Name
	}
	s, _ = models.GetOrgPid(saleman.PlaceId)
	if s != nil {
		saleman.Place = s.Name
	}
	s, _ = models.GetOrgPid(saleman.RegionId)
	if s != nil {
		saleman.Region = s.Name
	}
	if err != nil {
		cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "根据id查询业务员失败", "业务员端/GetSalesmanById", err.Error(), c.Ctx.Input)
		c.Abort("根据id查询业务员失败")
		return
	}
	c.Data["saleman"] = saleman
	c.TplName = "clerk/clerk_data_allot_detail1.html"

}

//已联系用户
func (c *SalesmanController) GetUsersIsLink() {
	defer c.IsNeedTemplate()
	page, _ := c.GetInt("page")
	if page <= 0 {
		page = 1
	}
	pageSize, _ := c.GetInt("pageSize")
	if pageSize < 1 {
		pageSize = 15
	}
	id, err := c.GetInt("id")
	if err != nil {
		cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "获取id失败", "业务员端/GetUsersIsAuth", err.Error(), c.Ctx.Input)
		return
	}
	users, err := models.GetIslinkAlltomentList(page, pageSize, id)
	if err != nil {
		cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "获取已认证客户失败", "业务员端/GetUsersIsAuth", err.Error(), c.Ctx.Input)
		return
	}
	count, _ := models.GetIslinkAlltomentCount(id)
	pageCount := 0
	if count%pageSize == 0 {
		pageCount = (count / pageSize)
	} else {
		pageCount = (count / pageSize) + 1
	}
	c.Data["users"] = users
	c.Data["id"] = id
	c.Data["page"] = page
	c.Data["pageCount"] = pageCount
	c.Data["pageSize"] = pageSize
	c.Data["count"] = count
	c.TplName = "user/user_salesmanlinkman.html"
}

//未联系用户
func (c *SalesmanController) GetUsersNotIsLink() {
	defer c.IsNeedTemplate()
	page, _ := c.GetInt("page")
	if page <= 0 {
		page = 1
	}
	pageSize, _ := c.GetInt("pageSize")
	if pageSize < 1 {
		pageSize = 15
	}
	id, err := c.GetInt("id")
	beego.Info(id)
	if err != nil {
		cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "获取id失败", "业务员端/GetUsersNotIsAuth", err.Error(), c.Ctx.Input)
		return
	}
	users, err := models.GetNotIslinkAlltomentList(page, pageSize, id)
	if err != nil {
		cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "获取已认证客户失败", "业务员端/GetUsersNotIsAuth", err.Error(), c.Ctx.Input)
		return
	}
	count, _ := models.GetNotIslinkAlltomentCount(id)
	pageCount := 0
	if count%pageSize == 0 {
		pageCount = (count / pageSize)
	} else {
		pageCount = (count / pageSize) + 1
	}
	c.Data["users"] = users
	c.Data["id"] = id
	c.Data["pageCount"] = pageCount
	c.Data["page"] = page
	c.Data["pageSize"] = pageSize
	c.Data["count"] = count
	c.TplName = "user/user_salesmanlinkman.html"

}

//注册业务员账号页面
func (c *SalesmanController) ToAddSaleman() {
	defer c.IsNeedTemplate()
	c.Data["flag"] = 0
	c.TplName = "clerk/add_account.html"
}

//确定接口添加业务员
func (c *SalesmanController) AddSaleman() {
	defer c.ServeJSON()
	account := c.GetString("account")
	name := c.GetString("name")
	stationId, _ := c.GetInt("stationId")
	is_ok, _ := c.GetInt("account_status")
	invitationCode := c.GetString("invitationCode")

	salesman := &models.Salesman{
		Account:    account,
		Password:   "123456",
		CreateTime: time.Now(),
		Saleman:    name,
		IsOk:       is_ok,
		StaId:      stationId,
	}

	org_id, _ := models.GetOrgId(stationId)
	salesman.OrgId = org_id
	pid, _ := models.GetOrgPid(org_id)

	if strings.Contains(pid.Name, "组") {
		salesman.GroupId = pid.Id
		pid1, _ := models.GetOrgPid(pid.ParentId)
		pid2, _ := models.GetOrgPid(pid1.ParentId)
		pid3, _ := models.GetOrgPid(pid2.ParentId)
		salesman.DepId = pid1.Id
		salesman.PlaceId = pid2.Id
		salesman.RegionId = pid3.Id
		salesman.InviteCode = pid1.InvitationCodePrefix + invitationCode
	} else if strings.Contains(pid.Name, "部") {
		pid2, _ := models.GetOrgPid(pid.ParentId)
		pid3, _ := models.GetOrgPid(pid2.ParentId)
		salesman.DepId = pid.Id
		salesman.PlaceId = pid2.Id
		salesman.RegionId = pid3.Id
		salesman.InviteCode = pid.InvitationCodePrefix + invitationCode
	} else if strings.Contains(pid.Name, "省") {
		pid3, _ := models.GetOrgPid(pid.ParentId)
		salesman.PlaceId = pid.Id
		salesman.RegionId = pid3.Id
		salesman.InviteCode = pid.InvitationCodePrefix + invitationCode
	} else if strings.Contains(pid.Name, "区") {
		salesman.RegionId = pid.Id
		salesman.InviteCode = pid.InvitationCodePrefix + invitationCode
	}

	salemanList, err := models.GetSaleman()
	for _, v := range salemanList {
		if v.InviteCode == salesman.InviteCode || v.Account == salesman.Account {
			c.Data["json"] = map[string]interface{}{"ret": 403, "msg": "已经存在相同的账号或邀请码"}
			cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "已经存在相同的账号或邀请码", "业务员端/AddSaleman", "", c.Ctx.Input)
			return
		}
	}

	err = models.AddSaleman(salesman)
	if err != nil {
		c.Data["json"] = map[string]interface{}{"ret": 403, "msg": "增加业务员账号失败", "err": err.Error()}
		cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "增加业务员账号失败", "业务员端/AddSaleman", err.Error(), c.Ctx.Input)
		return
	}
	c.Data["json"] = map[string]interface{}{"ret": 200, "msg": "增加业务员账号成功"}
}

//业务员组织架构
//获取组织架构和岗位信息
func (c *SystemController) GetSalemanOrganizationStation() {
	s := models.QuerySalemanDisplayQn()
	o, err := models.GetSalemanOrganizationStations()
	if err != nil {
		cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "获取组织架构岗位信息失败", "系统管理/GetOrganizationStation", err.Error(), c.Ctx.Input)
	}

	// cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "获取组织架构岗位信息成功", "系统管理/GetOrganizationStation", "", c.Ctx.Input)
	c.Data["json"] = map[string]interface{}{"ret": 200, "organizationStationList": o, "station": s}
	c.ServeJSON()
}

//业务员分配
func (c *SalesmanController) GetSalesmanAllocation() {
	startTime := time.Now()
	//业务员管理列表
	defer c.IsNeedTemplate()
	pars := []interface{}{}
	condition := ""
	pageNum, _ := c.GetInt("page", 1) //分页信息（第几页）
	if pageNum < 1 {
		pageNum = 1
	}
	pageSize, _ := c.GetInt("pageSize")
	if pageSize < 1 {
		pageSize = 30
	}
	region := c.GetString("district")
	if region != "" {
		condition += " AND region_id = ?"
		pars = append(pars, region)
	}

	place := c.GetString("province")
	if place != "" {
		condition += " AND place_id = ?"
		pars = append(pars, place)
	}

	operaDep := c.GetString("office")
	if operaDep != "" {
		condition += " AND dep_id = ?"
		pars = append(pars, operaDep)
	}

	saleman := c.GetString("name")
	beego.Info(saleman)
	if saleman != "" {
		condition += " AND saleman = ?"
		pars = append(pars, saleman)
	}

	account := c.GetString("account")
	if account != "" {
		condition += " AND account = ?"
		pars = append(pars, account)
	}

	createTime := c.GetString("auth_time")
	var startcreateTimeTime, endcreateTimeTime string
	if createTime != "" {
		createTimes := strings.Split(createTime, "~")
		startcreateTimeTime = createTimes[0] + " 00:00:00"
		endcreateTimeTime = createTimes[1] + " 23:59:59"
		condition += ` AND allotment_time >= ? AND allotment_time <= ?`
		pars = append(pars, startcreateTimeTime)
		pars = append(pars, endcreateTimeTime)
	}

	stationId := c.User.StationId
	orgStr, err := cache.GetCacheDataByStation(stationId)
	if err != nil {
		cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "业务员管理根据岗位获取数据权限失败", "业务员管理/GetSalesmanAllocation", err.Error(), c.Ctx.Input)
	}

	salesmanlist, err := models.GetSalesmanList(utils.StartIndex(pageNum, pageSize), pageSize, condition, orgStr, pars...)
	s, err := cache.GetSysOrganization()
	salemanMap := make(map[string]models.SysOrganization, len(s))
	for _, v := range s {
		salemanMap[strconv.Itoa(v.Id)] = v
	}
	finish := sync.WaitGroup{}
	finish.Add(len(salesmanlist))
	for k := range salesmanlist {
		go services.ComputeValue(salesmanlist[k], salemanMap, &finish)
		//v.IsLinkAllotment, _ = models.GetIslinkAlltomentCount(v.Id)
		//v.NotIsLinkAllotment, _ = models.GetNotIslinkAlltomentCount(v.Id)
		//a, b, c, _ := services.FindOrgTypeDepartment(salemanMap, strconv.Itoa(v.DepId))
		//v.Place = a
		//v.OperaDep = b
		//v.Region = c
	}
	count, err := models.GetSalesmanCount(condition, orgStr, pars...)
	if err != nil {
		cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "查询业务员总数失败", "业务员端/GetSalesmanAllocation", err.Error(), c.Ctx.Input)
		c.Abort("查询业务员总数失败")
		return
	}
	pageCount := utils.PageCount(count, pageSize)
	//===============数据回填
	var provinces []models.SysOrganization
	var offices []models.SysOrganization
	var regions []models.SysOrganization
	r, _ := strconv.Atoi(region)
	p, _ := strconv.Atoi(place)
	o, _ := strconv.Atoi(operaDep)
	//region
	regions, err = models.GetParentSysOrganizationById(28) //大区
	if err != nil {
		cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "查询业务员总数失败", "业务员分配/GetSalesmanAllocation", err.Error(), c.Ctx.Input)
		c.Abort("查询地区失败")
		return
	}

	//province
	if r != 0 {
		provinces, err = models.GetParentSysOrganizationById(r) //省
		if err != nil {
			cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "查询业务员总数失败", "业务员分配/GetSalesmanAllocation", err.Error(), c.Ctx.Input)
			c.Abort("查询省份失败")
			return
		}
	}
	//office
	if p != 0 {
		offices, err = models.GetParentSysOrganizationById(p) //营业部
		if err != nil {
			cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "查询业务员总数失败", "业务员分配/GetSalesmanAllocation", err.Error(), c.Ctx.Input)
			c.Abort("查询营业部失败")
			return
		}
	}
	fmt.Println(c.User.Id)
	allotment, _ := models.GetSalesmanAllotmentCount(c.User.Id)
	//=================end
	c.Data["list"] = salesmanlist
	c.Data["currPage"] = pageNum
	c.Data["count"] = count
	c.Data["pageSize"] = pageSize
	c.Data["pageCount"] = pageCount
	c.Data["allotment"] = allotment
	//回填
	c.Data["region"] = r
	c.Data["regions"] = regions
	c.Data["province"] = p
	c.Data["provinces"] = provinces
	c.Data["office"] = o
	c.Data["offices"] = offices
	c.Data["saleman"] = saleman
	c.Data["account"] = account
	// c.Data["createTime"] = createTime
	finish.Wait()
	c.TplName = "clerk/clerk_data_allot_list.html"
	fmt.Println(time.Since(startTime))
}

//分配
func (c *SalesmanController) Allotment() {
	defer c.ServeJSON()
	var salesmans models.SalesmanResponse
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &salesmans)
	if err != nil {
		cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "参数解析失败", "业务员分配/Allotment", err.Error(), c.Ctx.Input)
		c.Data["json"] = map[string]interface{}{"ret": 304, "msg": "参数解析失败", "err": err.Error()}
		return
	}
	if len(salesmans.CheckIds) < 0 || salesmans.Num < 0 {
		cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "参数解析有误", "业务员分配/Allotment", err.Error(), c.Ctx.Input)
		c.Data["json"] = map[string]interface{}{"ret": 304, "msg": "参数解析有误", "err": err.Error()}
		return
	}
	list, err := models.GetAllotmentList(c.User.Id)
	if err != nil {
		cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "获取数据list失败", "业务员分配/Allotment", err.Error(), c.Ctx.Input)
		c.Data["json"] = map[string]interface{}{"ret": 304, "msg": "获取数据list失败", "err": err.Error()}
		return
	}
	beego.Info("1", list, err)
	shares := services.ShareEqual(salesmans.CheckIds, salesmans.Num, list)
	err = models.UpdataAllotment(shares)
	beego.Info(shares, err)
	if err != nil {
		cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "更新业务表失败", "业务员分配/Allotment", err.Error(), c.Ctx.Input)
		c.Data["json"] = map[string]interface{}{"ret": 304, "msg": "更新业务表失败", "err": err.Error()}
		return
	}

	c.Data["json"] = map[string]interface{}{"ret": 200, "msg": "分配成功"}
}

//业务员账号编辑
func (c *SalesmanController) ToUpdateSalesman() {
	defer c.IsNeedTemplate()
	id, err := c.GetInt("id")
	flag, err := c.GetInt("flag")
	beego.Info(id, err)
	if err != nil {
		cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "获取id失败", "业务员端/UpdateSalesman", err.Error(), c.Ctx.Input)
		return
	}
	saleman, err := models.GetSalemanStationById(id)
	if saleman.InviteCode != "" {
		saleman.InviteCode = saleman.InviteCode[4:]
	}
	if err != nil {
		beego.Info(err)
		cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "获取业务员信息失败", "业务员端/UpdateSalesman", err.Error(), c.Ctx.Input)
		return
	}
	c.Data["sid"] = id
	c.Data["flag"] = flag
	c.Data["saleman"] = saleman
	c.TplName = "clerk/add_account.html"
}

//确定接口添加业务员
func (c *SalesmanController) UpdateSaleman() {
	defer c.ServeJSON()
	sid, _ := c.GetInt("sid")
	account := c.GetString("account")
	name := c.GetString("name")
	stationId, _ := c.GetInt("stationId")
	is_ok, _ := c.GetInt("account_status")
	invitationCode := c.GetString("invitationCode")
	password := c.GetString("password")
	salesman := &models.Salesman{
		Id:         sid,
		Account:    account,
		Password:   password,
		CreateTime: time.Now(),
		Saleman:    name,
		IsOk:       is_ok,
		StaId:      stationId,
	}

	org_id, _ := models.GetOrgId(stationId)
	salesman.OrgId = org_id
	pid, _ := models.GetOrgPid(org_id)

	if strings.Contains(pid.Name, "组") {
		salesman.GroupId = pid.Id
		pid1, _ := models.GetOrgPid(pid.ParentId)
		pid2, _ := models.GetOrgPid(pid1.ParentId)
		pid3, _ := models.GetOrgPid(pid2.ParentId)
		salesman.DepId = pid1.Id
		salesman.PlaceId = pid2.Id
		salesman.RegionId = pid3.Id
		salesman.InviteCode = pid1.InvitationCodePrefix + invitationCode
	} else if strings.Contains(pid.Name, "部") {
		pid2, _ := models.GetOrgPid(pid.ParentId)
		pid3, _ := models.GetOrgPid(pid2.ParentId)
		salesman.DepId = pid.Id
		salesman.PlaceId = pid2.Id
		salesman.RegionId = pid3.Id
		salesman.InviteCode = pid.InvitationCodePrefix + invitationCode
	} else if strings.Contains(pid.Name, "省") {
		pid3, _ := models.GetOrgPid(pid.ParentId)
		salesman.PlaceId = pid.Id
		salesman.RegionId = pid3.Id
		salesman.InviteCode = pid.InvitationCodePrefix + invitationCode
	} else if strings.Contains(pid.Name, "区") {
		salesman.RegionId = pid.Id
		salesman.InviteCode = pid.InvitationCodePrefix + invitationCode
	}

	err := models.UpdateSalesman(salesman)
	if err != nil {
		c.Data["json"] = map[string]interface{}{"ret": 403, "msg": "更新业务员账号失败", "err": err.Error()}
		cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "增加业务员账号失败", "业务员端/AddSaleman", err.Error(), c.Ctx.Input)
		return
	}
	c.Data["json"] = map[string]interface{}{"ret": 200, "msg": "增加业务员账号成功"}
}

//数据统计
func (c *SalesmanController) DataStatistics() {
	defer c.IsNeedTemplate()

	c.TplName = "clerk/clerk_data_stat.html"
	pars := []interface{}{}
	condition := ""

	region := c.GetString("district")
	if region != "" {
		condition += " AND s.region_id = ?"
		pars = append(pars, region)
	}

	place := c.GetString("province")
	if place != "" {
		condition += " AND s.place_id = ?"
		pars = append(pars, place)
	}

	operaDep := c.GetString("office")
	if operaDep != "" {
		condition += " AND s.dep_id = ?"
		pars = append(pars, operaDep)
	}

	//===============数据回填
	var provinces []models.SysOrganization
	var offices []models.SysOrganization
	var regions []models.SysOrganization
	r, _ := strconv.Atoi(region)
	p, _ := strconv.Atoi(place)
	o, _ := strconv.Atoi(operaDep)
	//region
	regions, err := models.GetParentSysOrganizationById(28)
	if err != nil {
		cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "查询业务员总数失败", "业务员端/DataStatistics", err.Error(), c.Ctx.Input)
		c.Abort("查询地区失败")
		return
	}

	//province
	if r != 0 {
		provinces, err = models.GetParentSysOrganizationById(r)
		if err != nil {
			cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "查询业务员总数失败", "业务员端/DataStatistics", err.Error(), c.Ctx.Input)
			c.Abort("查询省份失败")
			return
		}
	}
	//office
	if p != 0 {
		offices, err = models.GetParentSysOrganizationById(p)
		if err != nil {
			cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "查询业务员总数失败", "业务员端/DataStatistics", err.Error(), c.Ctx.Input)
			c.Abort("查询营业部失败")
			return
		}
	}
	//=================end
	//回填
	c.Data["region"] = r
	c.Data["regions"] = regions
	c.Data["province"] = p
	c.Data["provinces"] = provinces
	c.Data["office"] = o
	c.Data["offices"] = offices

	//权限
	c.Data["permission"] = true

	//初始化返回值
	var rCount, ocrCount, postAuthCount, authPassCount, loanCount, loanSucceedCount, rCount2, ocrCount2, postAuthCount2, authPassCount2, loanCount2, loanSucceedCount2 int
	var moneyOv, authRate, loanRate, authRate2, loanRate2 float64

	c.Data["data"] = map[string]interface{}{
		"moneyOv":           moneyOv,
		"rCount":            rCount,
		"ocrCount":          ocrCount,
		"postAuthCount":     postAuthCount,
		"authPassCount":     authPassCount,
		"loanCount":         loanCount,
		"loanSucceedCount":  loanSucceedCount,
		"authRate":          authRate,
		"loanRate":          loanRate,
		"rCount2":           rCount2,
		"ocrCount2":         ocrCount2,
		"postAuthCount2":    postAuthCount2,
		"authPassCount2":    authPassCount2,
		"loanCount2":        loanCount2,
		"loanSucceedCount2": loanSucceedCount2,
		"authRate2":         authRate2,
		"loanRate2":         loanRate2,
	}
	// dataState := true
	// state := c.GetString("state")
	// if state == "2" {
	// 	dataState = false
	// }

	stationId := c.User.StationId
	jud := c.GetString("judge")
	if jud != "2" && stationId != 1 {
		orgid, err := models.QueryOrgIdByStationId(stationId)
		if err != nil {
			cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "业务员管理根据岗位获取组织id失败", "业务员管理/DataStatistics", err.Error(), c.Ctx.Input)
			c.Abort("业务员管理根据岗位获取组织id失败")
			return
		}
		if orgid != 28 {
			for i := 0; i < 3; i++ {
				parentId, err := models.QueryParentIdByOrgId(orgid)
				if err != nil {
					cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "业务员管理根据岗位获取组织父id失败", "业务员管理/DataStatistics", err.Error(), c.Ctx.Input)
					c.Abort("业务员管理根据岗位获取组织父id失败")
					return
				}
				if parentId == 28 {
					if i == 1 {
						o = 0
					}
					r = orgid
					break
				}
				if i == 0 {
					p = orgid
					o = orgid
				} else if i == 1 {
					p = orgid
				}
				orgid = parentId
			}
		}
		if r != 0 {
			condition += " AND s.region_id = ?"
			pars = append(pars, r)
		}
		if p != 0 {
			condition += " AND s.place_id = ?"
			pars = append(pars, p)
		}
		if o != 0 {
			condition += " AND s.dep_id = ?"
			pars = append(pars, o)
		}
	}
	//权限判断
	isB, err := models.QuerySysRrganizationIdByRPO(r, p, o, stationId)
	if err != nil {
		cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "业务员管理根据岗位获取省区营业部权限失败", "业务员管理/DataStatistics", err.Error(), c.Ctx.Input)
		c.Abort("业务员管理根据岗位获取省区营业部权限失败")
		return
	}
	if !isB {
		c.Data["permission"] = false
		return
	}
	//====
	orgStr, err := cache.GetCacheDataByStation(stationId)
	if err != nil {
		cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "业务员管理根据岗位获取数据权限失败", "业务员管理/DataStatistics", err.Error(), c.Ctx.Input)
		c.Abort("业务员管理根据岗位获取数据权限失败")
		return
	}

	uids, err := models.QueryUidsBySalesmanIds(condition, orgStr, true, pars...)
	if err != nil {
		cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "查询业务员列表失败", "业务员端/DataStatistics", err.Error(), c.Ctx.Input)
		c.Abort("查询业务员列表失败")
		return
	}

	str := strings.Join(uids, ",")
	if str != "" {
		//注册量
		rCount = len(uids)
		//身份认证通过量
		ocrCount, err = models.QuerySalesmanOcrAuthCount(str)
		if err != nil {
			cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "查询身份认证通过量失败", "业务员端/DataStatistics", err.Error(), c.Ctx.Input)
			c.Abort("查询身份认证通过量失败")
			return
		}
		//提交授信申请量
		postAuthCount, err = models.QuerySalesmanPostAuthCount(str)
		if err != nil {
			cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "查询提交授信申请量失败", "业务员端/DataStatistics", err.Error(), c.Ctx.Input)
			c.Abort("查询提交授信申请量失败")
			return
		}
		//授信通过量
		authPassCount, err = models.QuerySalesmanAuthPassCount(str)
		if err != nil {
			cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "查询授信通过量失败", "业务员端/DataStatistics", err.Error(), c.Ctx.Input)
			c.Abort("查询授信通过量失败")
			return
		}
		//申请借款笔数
		loanCount, err = models.QuerySalesmanLoanCount(str)
		if err != nil {
			cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "查询申请借款笔数失败", "业务员端/DataStatistics", err.Error(), c.Ctx.Input)
			c.Abort("查询申请借款笔数失败")
			return
		}
		//放款成功笔数
		loanSucceedCount, err = models.QuerySalesmanLoanSucceedCount(str)
		if err != nil {
			cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "查询放款成功笔数失败", "业务员端/DataStatistics", err.Error(), c.Ctx.Input)
			c.Abort("查询放款成功笔数失败")
			return
		}
		//授信通过率/日
		authRate = utils.Divide(authPassCount, rCount)
		//放款通过率/日
		loanRate = utils.Divide(loanSucceedCount, rCount)
		// if dataState {
		//当日逾期金额
		loanOverdue, err := models.QuerySalesmanLoanOverdueMoney(str)
		if err != nil {
			cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "查询当日逾期金额失败", "业务员端/DataStatistics", err.Error(), c.Ctx.Input)
			c.Abort("查询当日逾期金额失败")
			return
		}
		//当日应还金额
		loanAllMoney, err := models.QuerySalesmanLoanAllMoneyByToday(str)
		if err != nil {
			cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "查询当日应还金额失败", "业务员端/DataStatistics", err.Error(), c.Ctx.Input)
			c.Abort("查询当日应还金额失败")
			return
		}
		//金额首逾率/日
		moneyOv = utils.FloatToFloat(loanOverdue, loanAllMoney)
		// }
	}

	//全部数据
	uids, err = models.QueryUidsBySalesmanIds(condition, orgStr, false, pars...)
	if err != nil {
		cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "查询业务员列表失败", "业务员端/DataStatistics", err.Error(), c.Ctx.Input)
		c.Abort("查询业务员列表失败")
		return
	}
	str2 := strings.Join(uids, ",")
	if str2 != "" {
		//注册量
		rCount2 = len(uids)
		//身份认证通过量
		ocrCount2, err = models.QuerySalesmanOcrAuthCount(str2)
		if err != nil {
			cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "查询身份认证通过量失败", "业务员端/DataStatistics", err.Error(), c.Ctx.Input)
			c.Abort("查询身份认证通过量失败")
			return
		}
		//提交授信申请量
		postAuthCount2, err = models.QuerySalesmanPostAuthCount(str2)
		if err != nil {
			cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "查询提交授信申请量失败", "业务员端/DataStatistics", err.Error(), c.Ctx.Input)
			c.Abort("查询提交授信申请量失败")
			return
		}
		//授信通过量
		authPassCount2, err = models.QuerySalesmanAuthPassCount(str2)
		if err != nil {
			cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "查询授信通过量失败", "业务员端/DataStatistics", err.Error(), c.Ctx.Input)
			c.Abort("查询授信通过量失败")
			return
		}
		//申请借款笔数
		loanCount2, err = models.QuerySalesmanLoanCount(str2)
		if err != nil {
			cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "查询申请借款笔数失败", "业务员端/DataStatistics", err.Error(), c.Ctx.Input)
			c.Abort("查询申请借款笔数失败")
			return
		}
		//放款成功笔数
		loanSucceedCount2, err = models.QuerySalesmanLoanSucceedCount(str2)
		if err != nil {
			cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "查询放款成功笔数失败", "业务员端/DataStatistics", err.Error(), c.Ctx.Input)
			c.Abort("查询放款成功笔数失败")
			return
		}
		//授信通过率/日
		authRate2 = utils.Divide(authPassCount2, rCount2)
		//放款通过率/日
		loanRate2 = utils.Divide(loanSucceedCount2, rCount2)
	}

	c.Data["data"] = map[string]interface{}{
		"moneyOv":           utils.Float64ToString(moneyOv * 100),
		"rCount":            rCount,
		"ocrCount":          ocrCount,
		"postAuthCount":     postAuthCount,
		"authPassCount":     authPassCount,
		"loanCount":         loanCount,
		"loanSucceedCount":  loanSucceedCount,
		"authRate":          utils.Float64ToString(authRate * 100),
		"loanRate":          utils.Float64ToString(loanRate * 100),
		"rCount2":           rCount2,
		"ocrCount2":         ocrCount2,
		"postAuthCount2":    postAuthCount2,
		"authPassCount2":    authPassCount2,
		"loanCount2":        loanCount2,
		"loanSucceedCount2": loanSucceedCount2,
		"authRate2":         utils.Float64ToString(authRate2 * 100),
		"loanRate2":         utils.Float64ToString(loanRate2 * 100),
	}
}

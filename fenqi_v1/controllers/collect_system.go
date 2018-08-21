package controllers

import (
	"fenqi_v1/cache"
	"fenqi_v1/models"
	"fenqi_v1/services"
	"fenqi_v1/utils"
	"strconv"
	"strings"
)

type CollectSystemController struct {
	BaseController
}

func (c *CollectSystemController) UserList() {
	c.IsNeedTemplate()
	// fmt.Println()
	uidName := c.Ctx.GetCookie("uid")
	if uidName == "admin" {
		c.Data["isAdmin"] = true
	}
	page, _ := c.GetInt("page")
	if page < 1 {
		page = 1
	}

	condition := ""
	pars := []string{}
	if account := c.GetString("account"); account != "" {
		condition += ` and su.name like ?`
		account = "%" + account + "%"
		pars = append(pars, account)
	}
	if username := c.GetString("username"); username != "" {
		condition += ` and su.displayname like ?`
		username = "%" + username + "%"
		pars = append(pars, username)
	}
	if role := c.GetString("role"); role != "" {
		condition += ` and ss.name like ?`
		role = "%" + role + "%"
		pars = append(pars, role)
	}
	if status := c.GetString("accountstate"); status != "" {
		condition += " and su.accountstatus=?"
		pars = append(pars, status)
	}
	list, err := models.CollectSysUserList(condition, pars, utils.StartIndex(page, utils.PageSize20), utils.PageSize20)
	if err != nil {
		cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "查询系统用户列表失败", "系统管理/UserList", err.Error(), c.Ctx.Input)
	}
	count := models.CollectSysUserCount(condition, pars)
	pagecount := utils.PageCount(count, utils.PageSize20)

	c.Data["currpage"] = page
	c.Data["pagecount"] = pagecount

	// 二维码
	if len(list) > 0 {
		for i := 0; i < len(list); i++ {
			list[i].Secret = utils.CreateXjdSecret(list[i].Id)
			list[i].AuthURL = utils.CreateXjdAuthURLEscape(list[i].Secret, list[i].Displayname)
		}
	}

	c.Data["list"] = list
	c.Data["count"] = count
	c.TplName = "collect_system/system_account_list.html"
}

//@router /user [get]
func (c *CollectSystemController) GetUser() {
	c.IsNeedTemplate()
	uid_str := c.GetString("uid")
	if uid_str == "" {
		c.Data["user"] = &models.SysUserMini{}
	} else {
		uid, _ := c.GetInt("uid")
		user, err := models.CollectSysUserDetail(uid)
		if err != nil {
			cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "根据uid查询系统用户数据失败", "系统管理/GetUser", err.Error(), c.Ctx.Input)
		}
		c.Data["user"] = user
	}
	list, err := models.CollectSysRoleList()
	if err != nil {
		cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "查询角色管理列表失败", "系统管理/GetUser", err.Error(), c.Ctx.Input)
	}
	c.Data["list"] = list
	c.TplName = "collect_system/system_account_form.html"
}

//@router /user [post]
func (c *CollectSystemController) UserAdd() {
	defer c.ServeJSON()
	uid, _ := c.GetInt("uid")
	var err error
	var user *models.SysUserMini
	if uid == 0 {
		user = &models.SysUserMini{}
		user.Name = c.GetString("account")
	} else {
		user, _ = models.CollectSysUserDetail(uid)
	}
	pwd := c.GetString("password")
	if pwd == "" {
		pwd = "111111"
	}
	user.Password = utils.MD5(pwd)

	user.Displayname = c.GetString("name")
	user.Email = c.GetString("email")
	//user.Role_id, _ = c.GetInt("roleId")

	user.Accountstatus = c.GetString("account_status")
	user.Station_id, _ = c.GetInt("stationId")

	user.Role_id, err = models.QueryCollectRoleIdByStationId(user.Station_id)
	if err != nil {
		cache.RecordLogs(uid, uid, user.Name, user.Displayname, "查询角色id出错", "添加系统用户/UserAdd", err.Error(), c.Ctx.Input)
		c.Data["json"] = map[string]interface{}{"ret": 304, "msg": err.Error()}
		return
	}

	user.Qn_account = c.GetString("qnAccount")
	user.Qn_password = c.GetString("qnPassword")
	user.Place_role = c.GetString("placeRole")
	user.AccountType, _ = c.GetInt("account_type")
	if uid == 0 {
		count, err := models.CollectSysUserCounts(user.Name)
		if err != nil {
			cache.RecordLogs(uid, uid, user.Name, user.Displayname, "查询系统用户异常", "添加系统用户/UserAdd", err.Error(), c.Ctx.Input)
			c.Data["json"] = map[string]interface{}{"ret": 304, "msg": "异常"}
			return
		}
		if count > 0 {
			c.Data["json"] = map[string]interface{}{"ret": 304, "msg": "账号不能重复"}
			return
		} else {
			user.IsCollectAccount = 1
			err = user.CollectInsert()
		}
	} else {
		err = user.Update()
		userType, err := models.GetCollectUserType(user.Station_id)
		err = models.UpdateCollectionScheduleType(uid, userType)
		if err != nil {
			cache.RecordLogs(uid, uid, user.Name, user.Displayname, "更新催收排班人员类型失败", "添加系统用户/UserAdd", err.Error(), c.Ctx.Input)
		}
		//判断该uid下是否有催收案件(mtype>=25)
		count := models.GetCaseCount(uid)
		if count > 0 {
			//根据岗位id得到org_id
			orgId := models.GetOrgIdByStationId(user.Station_id)
			if orgId > 0 {
				err = models.UpdateLoanOrgId(orgId, uid)
				if err != nil {
					cache.RecordLogs(uid, uid, user.Name, user.Displayname, "更新组织架构id失败", "添加系统用户/UserAdd", err.Error(), c.Ctx.Input)
				}
			}
		}

	}
	if err != nil {
		cache.RecordLogs(uid, uid, user.Name, user.Displayname, "添加系统用户失败", "添加系统用户/UserAdd", err.Error(), c.Ctx.Input)
		c.Data["json"] = map[string]interface{}{"ret": 304, "msg": err.Error()}
	} else {
		c.Data["json"] = map[string]interface{}{"ret": 200}
	}
}

func (c *CollectSystemController) DelUser() {
	defer c.ServeJSON()
	uid, _ := c.GetInt("uid")
	if uid == 0 {
		c.Data["json"] = map[string]interface{}{"ret": 304, "msg": "id错误"}
		return
	}
	if uid == 1 { // 系统管理员不给删
		c.Data["json"] = map[string]interface{}{"ret": 200}
		return
	}
	err := models.CollectDeleteUser(uid)
	if err != nil {
		cache.RecordLogs(uid, uid, "", "", "删除系统用户失败", "删除系统用户/DelUser", err.Error(), c.Ctx.Input)
		c.Data["json"] = map[string]interface{}{"ret": 304, "msg": err.Error()}
	} else {
		c.Data["json"] = map[string]interface{}{"ret": 200}
	}
}

func (c *CollectSystemController) RoleList() {
	c.IsNeedTemplate()
	rolelist, err := models.CollectSysRoleList()
	if err != nil {
		cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "查询角色管理列表失败", "系统管理/RoleList", err.Error(), c.Ctx.Input)
	}
	c.Data["list"] = rolelist
	c.TplName = "collect_system/system_role_list.html"
}

func (c *CollectSystemController) SelectRoleList() {
	defer c.ServeJSON()
	condition := ""
	displayName := c.GetString("DisplayName")
	if displayName != "" {
		condition += " AND displayname LIKE '" + displayName + "%'"
	}
	list, err := models.CollectSelectRoleList(condition)
	if err != nil {
		cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "查询角色管理列表失败", "系统管理/SelectRoleList", err.Error(), c.Ctx.Input)
	}
	if list == nil {
		c.Data["json"] = map[string]interface{}{"ret": 403, "msg": "没有匹配的角色"}
	} else {
		c.Data["json"] = map[string]interface{}{"ret": 200, "list": list}
	}
}

func (c *CollectSystemController) RoleEdit() {
	c.IsNeedTemplate()
	rid_str := c.GetString("rid")
	var err error
	var role *models.SysRole
	if rid_str == "" {
		role = &models.SysRole{}
	} else {
		rid, _ := c.GetInt("rid")

		role, err = models.SysRoleByRid(rid)
		if err != nil {
			cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "根据rid查询角色管理列表失败", "系统管理/RoleEdit", err.Error(), c.Ctx.Input)
			c.Ctx.WriteString(err.Error())
			return
		}
	}
	c.Data["role"] = role
	c.TplName = "collect_system/system_role_edit.html"
}

func (c *CollectSystemController) RoleAdd() {
	defer c.ServeJSON()
	var role *models.SysRole
	var err error
	rid, _ := c.GetInt("rid")
	if rid == 0 {
		role = &models.SysRole{}
	} else {
		role, err = models.SysRoleByRid(rid)
		if err != nil {
			cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "根据rid查询角色管理列表失败", "系统角色/RoleAdd", err.Error(), c.Ctx.Input)
			c.Data["json"] = map[string]interface{}{"ret": 304, "msg": err.Error()}
			return
		}
	}
	ids_str := c.GetString("checkId") // 菜单权限
	ids := strings.Split(ids_str, ",")
	role.Displayname = c.GetString("account")
	if rid == 0 {
		utils.Rc.Delete(utils.CacheKeyRoleMenuMapTreePrefix + strconv.Itoa(rid))
		err = role.CollectInsert(ids)
		cache.RecordLogs(rid, rid, role.Displayname, "角色", "系统角色", "新增角色成功/RoleAdd", "", c.Ctx.Input)
	} else {
		utils.Rc.Delete(utils.CacheKeyRoleMenuMapTreePrefix + strconv.Itoa(rid))
		err = role.Update(ids)
		cache.RecordLogs(rid, rid, role.Displayname, "角色", "系统角色", "修改角色成功/RoleAdd", "", c.Ctx.Input)
	}
	if err != nil {
		cache.RecordLogs(rid, rid, role.Displayname, "角色", "系统角色", "修改角色失败/RoleAdd", err.Error(), c.Ctx.Input)
		c.Data["json"] = map[string]interface{}{"ret": 304, "msg": err.Error()}
		return
	}
	utils.Rc.Delete(utils.CacheKeySystemMenu)
	utils.Rc.Delete(utils.CacheKeyRoleMenuMapTreePrefix + strconv.Itoa(c.User.StationId))
	c.Data["json"] = map[string]interface{}{"ret": 200}
}

func (c *CollectSystemController) MenuData() {
	defer c.ServeJSON()
	rid_str := c.GetString("role_id")
	var list []models.SysMenu
	var err error
	if rid_str == "all" { // 所有菜单
		list, err = models.GetCollectSysMenuTreeAll()
	} else {
		rid, _ := c.GetInt("rid")
		// 该角色有的菜单
		list, err = models.GetCollectSysMenuTreeByRoleId(rid)
	}
	if err != nil && err.Error() != "<QuerySeter> no row found" {
		cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "MenuData查询菜单失败", "系统管理/MenuData", err.Error(), c.Ctx.Input)
		c.Data["json"] = map[string]interface{}{"ret": 304, "msg": err.Error()}
		return
	}
	m := services.GetSysMenuZTree(list)
	c.Data["json"] = map[string]interface{}{"ret": 200, "m": m}
}

//删除角色
func (c *CollectSystemController) DelRole() {
	defer c.ServeJSON()
	rid, _ := c.GetInt("rid")

	err := models.DelRole(rid)
	if err != nil {
		cache.RecordLogs(rid, rid, "", "角色", "系统角色", "删除角色成功失败/MenuData", "", c.Ctx.Input)
		c.Data["json"] = map[string]interface{}{"ret": 304, "msg": err.Error()}
	} else {
		cache.RecordLogs(rid, rid, "", "角色", "系统角色", "删除角色成功/MenuData", "", c.Ctx.Input)
		c.Data["json"] = map[string]interface{}{"ret": 200}
	}
}

func (c *CollectSystemController) Organization() {
	c.IsNeedTemplate()
	c.TplName = "collect_system/system_organization_list.html"
}

//获取组织架构 菜单 数据
func (c *CollectSystemController) BaseData() {
	o, err := services.GetCollectSysOrganizationZTree()
	if err != nil {
		cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "BaseData获取组织架构信息失败", "系统管理/BaseData", err.Error(), c.Ctx.Input)
	}
	m, err := services.GetCollectAllSysMenuZTree()
	if err != nil {
		cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "BaseData返回ztree格式菜单数据失败", "系统管理/BaseData", err.Error(), c.Ctx.Input)
	}
	c.Data["json"] = map[string]interface{}{"ret": 200, "organizationData": o, "menuData": m}
	c.ServeJSON()
}

//添加组织架构
func (c *CollectSystemController) AddOrganization() {
	defer c.ServeJSON()
	pid, _ := c.GetInt("pid")
	var err error
	var organization *models.SysOrganization
	if pid <= 0 {
		c.Data["json"] = map[string]interface{}{"ret": 304, "msg": "上级组织选择错误!"}
		return
	} else {
		organization, err = models.GetSysOrganizationById(pid)
		if err != nil {
			c.Data["json"] = map[string]interface{}{"ret": 304, "msg": "上级组织不存在!"}
			return
		}
	}
	name := c.GetString("organizationName")
	if name == "" {
		c.Data["json"] = map[string]interface{}{"ret": 304, "msg": "组织名称不能为空!"}
		return
	}
	organization, err = models.GetSysOrganizationByName(name)
	if err != nil && err.Error() != "<QuerySeter> no row found" {
		cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "根据名称获取组织架构失败", "组织架构/AddOrganization", err.Error(), c.Ctx.Input)
	}

	if organization != nil {
		c.Data["json"] = map[string]interface{}{"ret": 304, "msg": "组织机构已经存在!"}
		return
	}
	organization = &models.SysOrganization{}
	organization.ParentId = pid
	organization.Name = name
	organization.IsCollectOrg = 1
	err = organization.CollectInsert()
	if err != nil {
		cache.RecordLogs(pid, pid, name, "组织架构", "组织架构/AddOrganization", "添加组织架构失败", err.Error(), c.Ctx.Input)
		c.Data["json"] = map[string]interface{}{"ret": 304, "msg": err.Error()}
		return
	}
	c.Data["json"] = map[string]interface{}{"ret": 200, "msg": "添加成功!"}
	// cache.RecordLogs(pid, pid, name, "组织架构", "组织架构/AddOrganization", "添加组织架构成功", "", c.Ctx.Input)
}

//编辑组织架构
func (c *CollectSystemController) EditOrganization() {
	defer c.ServeJSON()
	var err error
	name := c.GetString("name")
	oid, _ := strconv.Atoi(c.GetString("oid"))

	var organization *models.SysOrganization
	if oid <= 0 {
		c.Data["json"] = map[string]interface{}{"ret": 304, "msg": "组织选择错误!"}
		return
	} else {
		organization, err = models.GetSysOrganizationById(oid)
		if err != nil && err.Error() != "<QuerySeter> no row found" {
			cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "组织机构不存在", err.Error(), c.Ctx.Input)
			c.Data["json"] = map[string]interface{}{"ret": 304, "msg": "组织机构不存在!"}
			return
		}
	}
	if name == "" {
		c.Data["json"] = map[string]interface{}{"ret": 304, "msg": "组织名称不能为空!"}
		return
	}
	organization, err = models.GetSysOrganizationByName(name)
	if organization != nil {
		c.Data["json"] = map[string]interface{}{"ret": 304, "msg": "组织机构已经存在!"}
		return
	}
	organization = &models.SysOrganization{}
	organization.Id = oid
	organization.Name = name
	err = models.UpdateOrganizationName(name, oid)
	if err != nil {
		cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "编辑失败", err.Error(), c.Ctx.Input)
		c.Data["json"] = map[string]interface{}{"ret": 304, "msg": err.Error()}
		return
	}
	utils.Rc.Delete(utils.CacheKeySystemOrganization)
	c.Data["json"] = map[string]interface{}{"ret": 200, "msg": "编辑成功!"}
	// cache.RecordLogs(oid, oid, name, "组织架构", "组织架构", "编辑组织架构成功", "", c.Ctx.Input)
}

//添加岗位
func (c *CollectSystemController) AddStation() {
	defer c.ServeJSON()
	var err error
	var station *models.SysStation
	station = &models.SysStation{}
	stationName := c.GetString("stationName")
	orgId, _ := c.GetInt("orgId")

	if orgId <= 0 {
		c.Data["json"] = map[string]interface{}{"ret": 304, "msg": "岗位ID错误!"}
		return
	}
	station, err = models.GetSysStationByName(stationName, orgId)
	if err != nil && err.Error() != "<QuerySeter> no row found" {
		cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "根据名称获取岗位失败", "岗位/AddStation", err.Error(), c.Ctx.Input)
	}
	if station != nil {
		c.Data["json"] = map[string]interface{}{"ret": 304, "msg": "岗位已经存在!"}
		return
	}
	roleId, _ := c.GetInt("roleId")

	if roleId <= 0 {
		c.Data["json"] = map[string]interface{}{"ret": 304, "msg": "角色ID错误!"}
		return
	} else {
		role, err := models.SysRoleByRid(roleId)
		if err != nil && err.Error() != "<QuerySeter> no row found" {
			cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "根据roleId查找角色管理列表失败", "岗位/AddStation", err.Error(), c.Ctx.Input)
		}
		if role == nil {
			c.Data["json"] = map[string]interface{}{"ret": 304, "msg": "角色不存在!"}
			return
		}
	}
	typeStr := c.GetString("typeStr")
	typeArr := strings.Split(typeStr, ",")
	dataStr := c.GetString("dataStr")
	dataArr := strings.Split(dataStr, ",")
	if station == nil {
		station = &models.SysStation{}
	}

	station.Name = stationName
	station.RoleId = roleId
	station.OrgId = orgId
	station.IsCollectStation = 1
	err = station.CollectInsert(typeArr, dataArr)
	if err != nil {
		cache.RecordLogs(orgId, orgId, stationName, "岗位", "岗位/AddStation", "添加岗位失败", err.Error(), c.Ctx.Input)
		c.Data["json"] = map[string]interface{}{"ret": 304, "msg": "新增失败!"}
	} else {
		cache.RecordLogs(orgId, orgId, stationName, "岗位", "岗位/AddStation", "添加岗位成功", "", c.Ctx.Input)
		c.Data["json"] = map[string]interface{}{"ret": 200, "msg": "新增成功!"}
	}
}

//编辑岗位
func (c *CollectSystemController) UpdateStation() {
	defer c.ServeJSON()
	var err error
	stationId, _ := c.GetInt("stationId")

	if stationId <= 0 {
		c.Data["json"] = map[string]interface{}{"ret": 304, "msg": "岗位ID错误!"}
		return
	}
	var station *models.SysStation
	station = &models.SysStation{}
	stationName := c.GetString("stationName")
	roleId, _ := c.GetInt("roleId")

	if roleId <= 0 {
		c.Data["json"] = map[string]interface{}{"ret": 304, "msg": "角色ID错误!"}
		return
	} else {
		role, err := models.SysRoleByRid(roleId)
		if err != nil && err.Error() != "<QuerySeter> no row found" {
			cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "根据roleId查找角色管理列表失败", "岗位/UpdateStation", err.Error(), c.Ctx.Input)
		}
		if role == nil {
			c.Data["json"] = map[string]interface{}{"ret": 304, "msg": "角色不存在!"}
			return
		}
	}
	typeStr := c.GetString("typeStr")
	typeArr := strings.Split(typeStr, ",")
	dataStr := c.GetString("dataStr")
	dataArr := strings.Split(dataStr, ",")
	if station == nil {
		station = &models.SysStation{}
	}
	orgId, _ := c.GetInt("orgId")

	if orgId <= 0 {
		c.Data["json"] = map[string]interface{}{"ret": 304, "msg": "岗位ID错误!"}
		return
	}
	station.Name = stationName
	station.RoleId = roleId
	station.Id = stationId
	station.OrgId = orgId

	err = models.UpdateSysUserRoleId(stationId, roleId)
	if err != nil {
		c.Data["json"] = map[string]interface{}{"ret": 304, "msg": "更新系统用户角色失败!"}
		return
	}

	err = station.Update(typeArr, dataArr)
	if err != nil {
		cache.RecordLogs(orgId, orgId, stationName, "岗位", "岗位/UpdateStation", "修改岗位失败", err.Error(), c.Ctx.Input)
		c.Data["json"] = map[string]interface{}{"ret": 304, "msg": "修改失败!"}
	} else {
		cache.RecordLogs(orgId, orgId, stationName, "岗位", "岗位/UpdateStation", "修改岗位成功", "", c.Ctx.Input)
		c.Data["json"] = map[string]interface{}{"ret": 200, "msg": "修改成功!"}
	}
}

//根据组织架构Id获取岗位信息
func (c *CollectSystemController) GetStationData() {
	defer c.ServeJSON()
	orgId, _ := c.GetInt("orgId")

	m, err := models.SysStationListByOrgId(orgId) //根据组织架构ID获取岗位信息
	if err != nil {
		cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "根据组织架构ID获取岗位信息失败", "岗位/GetStationData", err.Error(), c.Ctx.Input)
	}
	c.Data["json"] = map[string]interface{}{"ret": 200, "stationData": m}
}

//获取角色列表
func (c *CollectSystemController) GetRoleList() {
	defer c.ServeJSON()
	rolelist, err := models.CollectSysRoleList()
	if err != nil {
		cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "获取角色列表信息失败", "系统管理/GetRoleList", err.Error(), c.Ctx.Input)
	}
	c.Data["json"] = map[string]interface{}{"ret": 200, "roleListData": rolelist}
}

//删除岗位
func (c *CollectSystemController) DelStation() {
	defer c.ServeJSON()
	stationId, _ := c.GetInt("stationId")

	err := models.DelStation(stationId)
	if err != nil {
		cache.RecordLogs(stationId, stationId, "", "岗位", "岗位/DelStation", "删除岗位失败", err.Error(), c.Ctx.Input)
		c.Data["json"] = map[string]interface{}{"ret": 304, "msg": err.Error()}
	} else {
		cache.RecordLogs(stationId, stationId, "", "岗位", "岗位/DelStation", "删除岗位成功", "", c.Ctx.Input)
		c.Data["json"] = map[string]interface{}{"ret": 200}
	}
}

//根据岗位ID，获取岗位信息
func (c *CollectSystemController) GetStationById() {
	defer c.ServeJSON()
	stationId, _ := c.GetInt("stationId")

	result, err := models.SysStationById(stationId)
	if err != nil {
		cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "根据岗位ID获取岗位信息失败", "岗位/GetStationById", err.Error(), c.Ctx.Input)
		c.Data["json"] = map[string]interface{}{"ret": 304, "msg": err.Error()}
	} else {
		cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "根据岗位ID获取岗位信息成功", "岗位/GetStationById", "", c.Ctx.Input)
		c.Data["json"] = map[string]interface{}{"ret": 200, "data": result}
	}
}

//获取催收组织架构和岗位信息
func (c *CollectSystemController) GetOrganizationStation() {
	s := models.CollectQueryDisplayQn()
	o, err := models.GetCollectOrganizationStations()
	if err != nil {
		cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "获取组织架构岗位信息失败", "系统管理/GetOrganizationStation", err.Error(), c.Ctx.Input)
	}
	// cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "获取组织架构岗位信息成功", "系统管理/GetOrganizationStation", "", c.Ctx.Input)
	c.Data["json"] = map[string]interface{}{"ret": 200, "organizationStationList": o, "station": s}
	c.ServeJSON()
}

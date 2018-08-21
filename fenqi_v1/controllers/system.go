package controllers

import (
	"encoding/json"
	"fenqi_v1/cache"
	"fenqi_v1/models"
	"fenqi_v1/services"
	"fenqi_v1/utils"
	"strconv"
	"strings"
)

type SystemController struct {
	BaseController
}

func (c *SystemController) UserList() {
	c.IsNeedTemplate()
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
	list, err := models.SysUserList(condition, pars, utils.StartIndex(page, utils.PageSize20), utils.PageSize20)
	if err != nil {
		cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "查询系统用户列表失败", "系统管理/UserList", err.Error(), c.Ctx.Input)
	}
	count := models.SysUserCount(condition, pars)
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
	c.TplName = "system/system_account_list.html"
}

//@router /user [get]
func (c *SystemController) GetUser() {
	c.IsNeedTemplate()
	uid_str := c.GetString("uid")
	if uid_str == "" {
		c.Data["user"] = &models.SysUserMini{}
	} else {
		uid, _ := c.GetInt("uid")
		user, err := models.SysUserDetail(uid)
		if err != nil {
			cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "根据uid查询系统用户数据失败", "系统管理/GetUser", err.Error(), c.Ctx.Input)
		}
		c.Data["user"] = user
	}
	//list, err := models.SysRoleList()
	//if err != nil {
	//	cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "查询角色管理列表失败", "系统管理/GetUser", err.Error(), c.Ctx.Input)
	//}
	//c.Data["list"] = list
	c.TplName = "system/system_account_form.html"
}

//@router /user [post]
func (c *SystemController) UserAdd() {
	defer c.ServeJSON()
	uid, _ := c.GetInt("uid")
	var err error
	var user *models.SysUserMini
	if uid == 0 {
		user = &models.SysUserMini{}
		user.Name = c.GetString("account")
	} else {
		user, _ = models.SysUserDetail(uid)
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
	user.Role_id, err = models.QueryRoleIdByStationId(user.Station_id)
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
		err = user.Insert()
	} else {
		err = user.Update()
	}
	if err != nil {
		cache.RecordLogs(uid, uid, user.Name, user.Displayname, "添加系统用户失败", "添加系统用户/UserAdd", err.Error(), c.Ctx.Input)
		c.Data["json"] = map[string]interface{}{"ret": 304, "msg": err.Error()}
	} else {
		c.Data["json"] = map[string]interface{}{"ret": 200}
	}
}

func (c *SystemController) DelUser() {
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
	err := models.DeleteUser(uid)
	if err != nil {
		cache.RecordLogs(uid, uid, "", "", "删除系统用户失败", "删除系统用户/DelUser", err.Error(), c.Ctx.Input)
		c.Data["json"] = map[string]interface{}{"ret": 304, "msg": err.Error()}
	} else {
		c.Data["json"] = map[string]interface{}{"ret": 200}
	}
}

func (c *SystemController) RoleList() {
	c.IsNeedTemplate()
	rolelist, err := models.SysRoleList()
	if err != nil {
		cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "查询角色管理列表失败", "系统管理/RoleList", err.Error(), c.Ctx.Input)
	}
	c.Data["list"] = rolelist
	c.TplName = "system/system_role_list.html"
}

func (c *SystemController) SelectRoleList() {
	defer c.ServeJSON()
	condition := ""
	displayName := c.GetString("DisplayName")
	if displayName != "" {
		condition += " AND displayname LIKE '" + displayName + "%'"
	}
	list, err := models.SelectRoleList(condition)
	if err != nil {
		cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "查询角色管理列表失败", "系统管理/SelectRoleList", err.Error(), c.Ctx.Input)
	}
	if list == nil {
		c.Data["json"] = map[string]interface{}{"ret": 403, "msg": "没有匹配的角色"}
	} else {
		c.Data["json"] = map[string]interface{}{"ret": 200, "list": list}
	}
}

func (c *SystemController) RoleEdit() {
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
	c.TplName = "system/system_role_edit.html"
}

func (c *SystemController) RoleAdd() {
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
		err = role.Insert(ids)
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
	c.Data["json"] = map[string]interface{}{"ret": 200}
}

func (c *SystemController) MenuData() {
	defer c.ServeJSON()
	rid_str := c.GetString("role_id")
	var list []models.SysMenu
	var err error
	if rid_str == "all" { // 所有菜单
		list, err = models.GetSysMenuTreeAll()
	} else {
		rid, _ := c.GetInt("rid")
		// 该角色有的菜单
		list, err = models.GetSysMenuTreeByRoleId(rid)
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
func (c *SystemController) DelRole() {
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

func (c *SystemController) Organization() {
	c.IsNeedTemplate()
	c.TplName = "system/system_organization_list.html"
}

//获取组织架构信息
func (c *SystemController) GetOrganizationList() {
	o, err := services.GetSysOrganizationZTree()
	if err != nil {
		cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "获取组织架构信息失败", "系统管理/GetOrganizationList", err.Error(), c.Ctx.Input)
	}
	// cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "获取组织架构信息成功", "系统管理/GetOrganizationList", "", c.Ctx.Input)
	c.Data["json"] = map[string]interface{}{"ret": 200, "organizationList": o}
	c.ServeJSON()
}

//获取组织架构 菜单 数据
func (c *SystemController) BaseData() {
	o, err := services.GetSysOrganizationZTree()
	if err != nil {
		cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "BaseData获取组织架构信息失败", "系统管理/BaseData", err.Error(), c.Ctx.Input)
	}
	m, err := services.GetAllSysMenuZTree()
	if err != nil {
		cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "BaseData返回ztree格式菜单数据失败", "系统管理/BaseData", err.Error(), c.Ctx.Input)
	}
	c.Data["json"] = map[string]interface{}{"ret": 200, "organizationData": o, "menuData": m}
	c.ServeJSON()
}

//添加组织架构
func (c *SystemController) AddOrganization() {
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
	err = organization.Insert()
	if err != nil {
		cache.RecordLogs(pid, pid, name, "组织架构", "组织架构/AddOrganization", "添加组织架构失败", err.Error(), c.Ctx.Input)
		c.Data["json"] = map[string]interface{}{"ret": 304, "msg": err.Error()}
		return
	}
	c.Data["json"] = map[string]interface{}{"ret": 200, "msg": "添加成功!"}
	// cache.RecordLogs(pid, pid, name, "组织架构", "组织架构/AddOrganization", "添加组织架构成功", "", c.Ctx.Input)
}

//编辑组织架构
func (c *SystemController) EditOrganization() {
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
func (c *SystemController) AddStation() {
	defer c.ServeJSON()
	var err error
	var station *models.SysStation
	station = &models.SysStation{}
	stationName := c.GetString("stationNa" + "me")
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
	err = station.Insert(typeArr, dataArr)
	if err != nil {
		cache.RecordLogs(orgId, orgId, stationName, "岗位", "岗位/AddStation", "添加岗位失败", err.Error(), c.Ctx.Input)
		c.Data["json"] = map[string]interface{}{"ret": 304, "msg": "新增失败!"}
	} else {
		cache.RecordLogs(orgId, orgId, stationName, "岗位", "岗位/AddStation", "添加岗位成功", "", c.Ctx.Input)
		c.Data["json"] = map[string]interface{}{"ret": 200, "msg": "新增成功!"}
	}
}

//编辑岗位
func (c *SystemController) UpdateStation() {
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
func (c *SystemController) GetStationData() {
	defer c.ServeJSON()
	orgId, _ := c.GetInt("orgId")

	m, err := models.SysStationListByOrgId(orgId) //根据组织架构ID获取岗位信息
	if err != nil {
		cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "根据组织架构ID获取岗位信息失败", "岗位/GetStationData", err.Error(), c.Ctx.Input)
	}
	c.Data["json"] = map[string]interface{}{"ret": 200, "stationData": m}
}

//获取角色列表
func (c *SystemController) GetRoleList() {
	defer c.ServeJSON()
	rolelist, err := models.SysRoleList()
	if err != nil {
		cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "获取角色列表信息失败", "系统管理/GetRoleList", err.Error(), c.Ctx.Input)
	}
	c.Data["json"] = map[string]interface{}{"ret": 200, "roleListData": rolelist}
}

//删除岗位
func (c *SystemController) DelStation() {
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
func (c *SystemController) GetStationById() {
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

//获取组织架构和岗位信息
func (c *SystemController) GetOrganizationStation() {
	s := models.QueryDisplayQn()
	o, err := models.GetOrganizationStations()
	if err != nil {
		cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "获取组织架构岗位信息失败", "系统管理/GetOrganizationStation", err.Error(), c.Ctx.Input)
	}
	// cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "获取组织架构岗位信息成功", "系统管理/GetOrganizationStation", "", c.Ctx.Input)
	c.Data["json"] = map[string]interface{}{"ret": 200, "organizationStationList": o, "station": s}
	c.ServeJSON()
}

//获取催收用户列表
func (c *SystemController) CollectionUserList() {
	c.IsNeedTemplate()
	page, _ := c.GetInt("page")
	if page < 1 {
		page = 1
	}
	condition := ""
	pars := []string{}
	condition += " and su.place_role in('班长','普通')"
	if account := c.GetString("account"); account != "" {
		condition += " and su.name=?"
		pars = append(pars, account)
	}
	if username := c.GetString("username"); username != "" {
		condition += " and su.displayname=?"
		pars = append(pars, username)
	}
	if callType := c.GetString("callType"); callType != "" {
		if callType == "通话中" {
			condition += " and tr.call_type in('呼入','转接')"
		}
		if callType == "未通话" {
			condition += " and (tr.call_type is null or tr.call_type='挂断')"
		}

	}
	list, err := models.SysCollectionUserList(condition, pars, utils.StartIndex(page, utils.PageSize20), utils.PageSize20) // 系统用户列表
	if err != nil {
		cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "获取系统用户列表信息失败", "系统管理/CollectionUserList", err.Error(), c.Ctx.Input)
	}
	// fmt.Println(condition)
	// cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "获取系统用户列表信息成功", "系统管理/CollectionUserList", "", c.Ctx.Input)
	count := models.SysCollectionUserCount(condition, pars)
	pagecount := utils.PageCount(count, utils.PageSize20)
	c.Data["currpage"] = page
	c.Data["pagecount"] = pagecount
	c.Data["list"] = list
	c.Data["count"] = count
	c.TplName = "quality_list.html"
}

func (c *SystemController) CollectionUserDetails() {
	c.IsNeedTemplate()
	uid, _ := c.GetInt("uid")

	user, err := models.SysUserDetail(uid)
	if err != nil {
		cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "根据uid查询系统用户细节数据失败", "系统管理/CollectionUserDetails", err.Error(), c.Ctx.Input)
	}
	// cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "根据uid查询系统用户细节数据失败", "系统管理/CollectionUserDetails", "", c.Ctx.Input)
	c.Data["user"] = user
	c.TplName = "quality_detail.html"
}

//根据系统用户id更新登录状态
func (c *SystemController) UpdateSysUserLoginState() {
	uid, _ := c.GetInt("uid")
	loginState := c.GetString("login_state")
	resultMap := make(map[string]interface{})
	resultMap["ret"] = 403
	defer func() {
		c.Data["json"] = resultMap
		c.ServeJSON()
	}()
	if uid <= 0 {
		cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "uid参数传递错误", "信审工作平台/UpdateSysUserLoginState", "", c.Ctx.Input)
		resultMap["err"] = "uid参数传递错误"
		return
	}
	err := models.UpdateLoginStateById(uid, loginState)
	if err != nil {
		cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "更新用户登录状态出错", "信审工作平台/UpdateSysUserLoginState", err.Error(), c.Ctx.Input)
		resultMap["err"] = "更新用户登录状态出错"
		return
	}
	m, err := models.FindSysUserById(uid)
	if err == nil {
		//更新redis
		if data, err2 := json.Marshal(m); err2 == nil && utils.Re == nil {
			utils.Rc.Put(utils.CacheKeyUserPrefix+m.Name, data, utils.RedisCacheTime_User)
		}
	}
	resultMap["ret"] = 200
	resultMap["msg"] = "更新用户登录状态成功"
}

func (c *SystemController) BtnPermissions() {
	defer c.ServeJSON()
	//用户角色
	roleId := c.User.RoleId
	//按钮ID
	btnIds := c.GetString("btnIds")
	btnArr := strings.Split(btnIds, ",")
	where := ""

	for _, v := range btnArr {
		where += "'" + v + "',"
	}
	where = strings.TrimRight(where, ",")
	btns, err := models.GetBtnPermissions(where, roleId)
	if err != nil {
		c.Data["json"] = map[string]interface{}{"ret": 304, "msg": err.Error()}
		return
	}
	btnMap := make(map[string]int)
	for i := 0; i < len(btns); i++ { //当前角色有权限的按钮
		btnMap[btns[i].ControlUrl] = btns[i].Id
		btns[i].IsShow = true
	}
	for _, m := range btnArr {
		if _, ok := btnMap[m]; !ok { //当前角色没有权限的按钮
			btn := new(models.SysBtn)
			btn.IsShow = false
			btn.ControlUrl = m
			btns = append(btns, *btn)
		}
	}
	c.Data["json"] = map[string]interface{}{"ret": 200, "data": btns}
}

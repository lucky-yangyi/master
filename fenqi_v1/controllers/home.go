package controllers

import (
	"fenqi_v1/services"
)

// HomeController 主页
type HomeController struct {
	BaseController
}

// Get 主页Get
func (c *HomeController) Get() {
	c.IsNeedTemplate()
	c.TplName = "index.html"
}

// Post 主页获取数据
func (c *HomeController) Post() {

	m, err := services.GetSysMenuTreeByRoleId(c.User.StationId) //根据岗位ID获取菜单
	if err != nil {
		c.Data["json"] = map[string]interface{}{"ret": 304, "msg": err.Error()}

	} else {
		c.Data["json"] = map[string]interface{}{"ret": 200, "data": m}
	}
	c.ServeJSON()
}

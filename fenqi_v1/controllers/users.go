package controllers

import (
	"fenqi_v1/cache"
	"fenqi_v1/models"
	"fenqi_v1/utils"
)

type UsersController struct {
	BaseController
}

//修改密码
func (c *UsersController) ModifyPassword() {
	defer c.ServeJSON()
	originalPw := c.GetString("orgpwd")
	newPw := c.GetString("newpwd")
	confirmPw := c.GetString("newpwd2")

	originalPw = utils.MD5(originalPw)
	if originalPw != c.User.Password {
		c.Data["json"] = map[string]interface{}{"ret": 304, "err": "原始密码输入错误，请重新输入！"}
		return
	}
	if len(newPw) < 6 { //密码要求6位数以上
		c.Data["json"] = map[string]interface{}{"ret": 304, "err": "为了账户安全，请输入6位数及以上的新密码！"}
		return
	}

	if newPw != confirmPw {
		c.Data["json"] = map[string]interface{}{"ret": 304, "err": "新密码与确认密码不一致！"}
		return
	}

	flag, err := models.UpdatePassword(c.User.Id, utils.MD5(newPw), newPw)
	if flag && err == nil {
		cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "密码修改成功", "用户管理/ModifyPassword", "", c.Ctx.Input)
		c.Data["json"] = map[string]interface{}{"ret": 200}
		return
	}
	c.Data["json"] = map[string]interface{}{"ret": 304, "err": "密码修改失败！"}
}

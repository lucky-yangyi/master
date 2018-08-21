package controllers

import (
	"encoding/json"
	"fenqi_v1/cache"
	"fenqi_v1/models"
	"fenqi_v1/utils"
	"github.com/astaxie/beego"
)

// AccountController 登录控制器  qie
type AccountController struct {
	beego.Controller
}

// Login 登录
func (c *AccountController) Login() {
	c.TplName = "system/login.html"
}

// CheckPassword 验证密码
func (c *AccountController) CheckPassword() {
	defer c.ServeJSON()
	name := c.GetString("username")
	password := c.GetString("password")
	password = utils.MD5(password)
	verify_code := c.GetString("verify_code") // 将军令
	m, _ := models.Login(name, password)
	if m != nil {
		ip := c.Ctx.Input.IP()
		flag := models.GetIsUseByIp(ip)
		if utils.RunMode == "release" {
			if ip != "127.0.0.1" {
				// flag := true
				if !flag {
					if verify_code == "" {
						cache.RecordLogs(m.Id, 0, m.Name, m.DisplayName, "登录", "系统", "请输入验证码", c.Ctx.Input)
						c.Data["json"] = map[string]interface{}{"ret": 304, "msg": "请输入验证码"}
						return
					}
					result, _ := utils.Authenticate(m.Id, verify_code)
					if !result {
						cache.RecordLogs(m.Id, 0, m.Name, m.DisplayName, "登录", "系统", "登录失败3:验证码错误", c.Ctx.Input)
						c.Data["json"] = map[string]interface{}{"ret": 304, "msg": "登录失败3:验证码错误"}
						return
					}
				}
			}
		}
		if data, err2 := json.Marshal(m); err2 == nil && utils.Re == nil {
			//添加日志记录
			cache.RecordLogs(m.Id, 0, m.Name, m.DisplayName, "登录", "系统", "登录成功", c.Ctx.Input)
			utils.Rc.Put(utils.CacheKeyUserPrefix+m.Name, data, utils.RedisCacheTime_User)
			//保存用户缓存和cookie
			c.Ctx.SetCookie("uid", m.Name)
			password2 := utils.MD5(password + utils.PasswordEncryptKey)
			c.Ctx.SetCookie("pid", password2)
			// c.Ctx.SetCookie("userId", strconv.Itoa(m.Id))
			c.Data["json"] = map[string]interface{}{"ret": 200, "msg": "登录成功"}
		} else {
			c.Data["json"] = map[string]interface{}{"ret": 304, "msg": "登录失败2！" + err2.Error()}
		}

	} else {
		c.Data["json"] = map[string]interface{}{"ret": 404, "msg": "登录失败！用户名或密码不正确"}
	}
}

// LoginOut 退出登录
func (c *AccountController) LoginOut() {
	uid := c.Ctx.GetCookie("uid")
	if uid != "" {
		cache.RecordLogs(0, 0, uid, "", "退出登录", "系统", "退出登录成功", c.Ctx.Input)
		//清除cookie
		c.Ctx.SetCookie("uid", "0", -1)
		c.Ctx.SetCookie("pid", "0", -1)
		//清除redis
		utils.Rc.Delete(utils.CacheKeyUserPrefix + uid)
	}
	c.Ctx.Redirect(302, "/login")
}

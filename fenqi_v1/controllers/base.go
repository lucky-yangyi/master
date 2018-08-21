package controllers

import (
	"encoding/json"
	"fenqi_v1/cache"
	"fenqi_v1/models"
	"fenqi_v1/services"
	"fenqi_v1/utils"
	"github.com/astaxie/beego"
	"net/url"
	"strings"
	"time"
	"zcm_tools/log"
)

var v1Log *log.Log

func init() {
	v1Log = log.Init()
}

// BaseController 基础controller
type BaseController struct {
	beego.Controller
	User *models.SysUser
}

// // Prepare 验证用户登录信息
func (c *BaseController) Prepare() {
	verify := false
	uid := c.Ctx.GetCookie("uid")
	pid := c.Ctx.GetCookie("pid")
	if uid != "" && pid != "" {
		if utils.Re == nil && utils.Rc.IsExist(utils.CacheKeyUserPrefix+uid) {
			if data, err := utils.Rc.RedisBytes(utils.CacheKeyUserPrefix + uid); err == nil {
				err = json.Unmarshal(data, &c.User)
				if c.User != nil && utils.MD5(c.User.Password+utils.PasswordEncryptKey) == pid {
					isCheck, _ := c.GetBool("is_check") //是否需要校验按钮操作权限
					if isCheck {
						btnId := c.GetString("btn_id")
						isPermission := models.IsBtnPermission(btnId, c.User.RoleId)
						if !isPermission {
							c.Data["json"] = map[string]interface{}{"ret": 304, "msg": "你没有该按钮操作权限"}
							c.ServeJSON()
							return
						}
					}
					if c.Ctx.Input.URL() != "/" {
						menus, _ := cache.GetCacheSysMenu()
						if menu, ok := menus[c.Ctx.Input.URL()]; ok {
							m, _ := services.GetSysMenuByRoleId(c.User.StationId)
							_, o := m[menu.Id]
							if !o {
								c.Ctx.WriteString("你没有该权限!")
								return
							}
						}
					}
					ip := c.Ctx.Input.IP()
					requestBody, _ := url.QueryUnescape(string(c.Ctx.Input.RequestBody))
					v1Log.Println("请求地址：", c.Ctx.Input.URI(), "用户信息：", string(data), "RequestBody：", requestBody, "IP：", ip)
					//重新保存用户状态
					if ip != "127.0.0.1" || ip == "60.191.125.34" || ip == "60.191.37.251" {
						utils.Rc.Put(utils.CacheKeyUserPrefix+c.User.Name, data, utils.RedisCacheTime_TwoHour)
					} else {
						utils.Rc.Put(utils.CacheKeyUserPrefix+c.User.Name, data, utils.RedisCacheTime_User)
					}
					verify = true
					c.Data["Id"] = c.User.Id
					c.Data["OutPutPhone"] = "057186982392,057188309017,057188305261,057188305727,057188305583,057188306939,057188309753,057186983151"
					c.Data["DisplayName"] = c.User.DisplayName
					c.Data["QnAccount"] = c.User.QnAccount
					c.Data["QnPassword"] = c.User.QnPassword
					c.Data["AccountType"] = c.User.AccountType
					c.Data["LoginState"] = c.User.LoginState
					id := c.User.Id
					cache.RecordOperateTime(id, time.Now())
				}
			}
		}
	} else if c.Ctx.Input.IsUpload() { //上传文件跳过验证
		verify = true
	}
	if !verify {
		if c.Ctx.Input.IsAjax() {
			c.Ctx.Output.JSON(map[string]interface{}{"ret": 408, "msg": "timeout"}, false, false)
			c.StopRun()
		} else {
			c.Ctx.Redirect(302, "/login")
			c.StopRun()
		}
	}
	paramsUid := c.GetString("uid")
	paramsUid = strings.Replace(paramsUid, " ", "+", -1)
	if paramsUid != "" {
		uidStr, err := utils.MyDesBase64Decrypt(paramsUid)
		if err != nil {
			//email.SendEmail("xjd_v1后台解密uid失败", err.Error()+c.Ctx.Input.URI(), "qxw@zcmlc.com;yangzb@zcmlc.com;wy@zcmlc.com")
		} else {
			c.Ctx.Input.SetParam("uid", uidStr)
		}
	}
	paramsId := c.GetString("id")
	paramsId = strings.Replace(paramsId, " ", "+", -1)
	if paramsId != "" {
		idStr, err := utils.MyDesBase64Decrypt(paramsId)
		if err != nil {
			//email.SendEmail("xjd_v1后台解密id失败", err.Error()+c.Ctx.Input.URI(), "qxw@zcmlc.com;yangzb@zcmlc.com;wy@zcmlc.com")
		} else {
			c.Ctx.Input.SetParam("id", idStr)
		}
	}
}

//是否需要模板
func (c *BaseController) IsNeedTemplate() {
	pushstate := c.GetString("pushstate")
	if pushstate != "1" {
		c.Layout = "layout/layout.html"
	}
}

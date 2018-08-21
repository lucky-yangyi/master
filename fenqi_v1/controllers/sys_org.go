package controllers

import (
	"fenqi_v1/models"

	"github.com/astaxie/beego"
)

type SysOrgController struct {
	BaseController
}

func (c *SysOrgController) GetSysOrg() {
	oId, _ := c.GetInt("oId")
	sos, err := models.GetParentSysOrganizationById(oId)
	beego.Debug(err, oId)
	c.Data["json"] = map[string]interface{}{"ret": 200, "sos": sos}
	c.ServeJSON()
}

package controllers

import (
	"fenqi_v1/models"
	"fenqi_v1/utils"
)

type RiskAdviseController struct {
	BaseController
}

func (c *RiskAdviseController) GetRiskAdvise() {
	defer c.ServeJSON()
	uid, err := c.GetInt("uid")
	businessType, _ := c.GetInt("business_type", 1)
	if err != nil {
		c.Data["json"] = map[string]interface{}{"ret": 403, "err": err.Error(), "msg": "获取uid失败"}
		return
	}
	r := &models.RiskAdvise{}
	r, err = models.GetRiskAdvises(uid, businessType)
	if err != nil && err.Error() != utils.MongdbErrNoRow() {
		c.Data["json"] = map[string]interface{}{"ret": 403, "err": err.Error(), "msg": "获取风控意见失败"}
		return
	}
	for _, v := range r.RiseItems {
		if v.RiskAdvise == "" {
			v.State = 1
		}
	}
	if r.AuditAdvice == "" {
		r.AuditAdvice = "暂无风控总建议"
	}
	c.Data["json"] = map[string]interface{}{"ret": 200, "data": r, "msg": "获取风控意见成功"}
}

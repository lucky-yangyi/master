package controllers

import (
	"encoding/json"
	"fenqi_v1/cache"
	"fenqi_v1/models"
	"fenqi_v1/utils"
	"github.com/tealeg/xlsx"
	"os"
)

type TeleSaleController struct {
	BaseController
}

//电销管理
func (this *TeleSaleController) TeleSaleList() {
	this.IsNeedTemplate()
	pars := []interface{}{}
	condition := ""
	pageNum, _ := this.GetInt("page", 1) //分页信息（第几页）
	if pageNum < 1 {
		pageNum = 1
	}
	pageSize, _ := this.GetInt("pageSize")
	if pageSize < 1 {
		pageSize = 30
	}
	name := this.GetString("name") //姓名
	if name != "" {
		condition += " AND b.name = ?"
		pars = append(pars, name)
	}
	phone := this.GetString("phone") //手机号
	if phone != "" {
		condition += " AND b.phone = ?"
		pars = append(pars, phone)
	}
	platfrom := this.GetString("platfrom")//平台
	if platfrom != ""{
		condition += " AND b.platfrom = ?"
		pars = append(pars, platfrom)
	}
	callman := this.GetString("callman")//电销员
	if callman !=""{
		condition += " AND b.call_man = ?"
		pars = append(pars, callman)
	}
	phone_status := this.GetString("phone_status") //号码状态
	if phone_status == "registered"{
		condition += " AND a.id >0"
	}else{
		if phone_status != "" {
			condition += " AND b.phone_status = ?"
			pars = append(pars, phone_status)
		}
	}
	final_auth_state := this.GetString("final_auth_state") //授信状态
	if final_auth_state == "QUEUEING" ||  final_auth_state == "HANDING" || final_auth_state == "OUTQUEUE"  {
		condition += " AND f.state = ?"
		pars = append(pars, final_auth_state)
	} else {
		if final_auth_state == "0,3" {
			condition += ` AND (c.final_auth_state IN (` + final_auth_state + `)` + ` OR c.final_auth_state IS NUll )`
		} else {
			if final_auth_state != "" {
				condition += ` AND c.final_auth_state IN (` + final_auth_state + `)`
			}
		}
	}
	allomatcount, err := models.TeleSaleAllomatPower(this.User.Id)
	if err != nil {
		cache.RecordLogs(this.User.Id, 0, this.User.Name, this.User.DisplayName, "查询电销员权限失败", "电销管理/TeleSaleList", err.Error(), this.Ctx.Input)
		this.Abort("查询电销员权限失败")
		return
	}
	if allomatcount == 0 {
		condition += " AND b.tele_id = ?"
		pars = append(pars, this.User.Id)
	}
	telesalelist, err := models.TeleSaleList(utils.StartIndex(pageNum, pageSize), pageSize, condition, pars...)
	if err != nil {
		cache.RecordLogs(this.User.Id, 0, this.User.Name, this.User.DisplayName, "查询电销员列表失败", "电销管理/TeleSaleList", err.Error(), this.Ctx.Input)
		this.Abort("查询电销列表失败")
		return
	}
	count, err := models.TeleSaleListCount(condition, pars...)
	if err != nil {
		cache.RecordLogs(this.User.Id, 0, this.User.Name, this.User.DisplayName, "获取电销员总数异常", "电销管理/TeleSaleList", err.Error(), this.Ctx.Input)
		this.Abort("查询电销列表总数失败")
		return
	}
	teleoperators, err := models.TeleSaleListDisplayname()
	if err != nil {
		cache.RecordLogs(this.User.Id, 0, this.User.Name, this.User.DisplayName, "获取处理员异常", "电销管理/TeleSaleList", err.Error(), this.Ctx.Input)
		this.Abort("获取处理员异常")
		return
	}
	pageCount := utils.PageCount(count, pageSize)
	this.Data["list"] = telesalelist
	this.Data["operators"] = teleoperators
	this.Data["currPage"] = pageNum
	this.Data["count"] = count
	this.Data["pageSize"] = pageSize
	this.Data["pageCount"] = pageCount
	this.TplName = "telesale/telesale.html"
}

//修改状态
func (this *TeleSaleController) ModifyStatus() {
	defer this.ServeJSON()
	id, _ := this.GetInt("id")
	phone_status := this.GetString("phone_status")
	remark := this.GetString("remark")
	if phone_status == "" {
		this.Data["json"] = map[string]interface{}{"ret": 403, "msg": "请选择状态"}
		return
	}
	err := models.UpdateTeleSaleStatus(phone_status, remark, id)
	if err != nil {
		cache.RecordLogs(this.User.Id, 0, this.User.Name, this.User.DisplayName, "修改处理结果有误", "电销管理/ModifyStatus", err.Error(), this.Ctx.Input)
		this.Data["json"] = map[string]interface{}{"ret": 403, "msg": "修改处理结果失败" + err.Error()}
		return
	}
	this.Data["json"] = map[string]interface{}{"ret": 200, "msg": "修改成功"}
}

//分配
func (this *TeleSaleController) Allotment() {
	defer this.ServeJSON()
	call_man := this.GetString("call_man")
	number, _ := this.GetInt("number")
	tele_id, _ := this.GetInt("id")
	if call_man == "" {
		this.Data["json"] = map[string]interface{}{"ret": 403, "msg": "请选择电销员"}
		return
	}
	if number == 0{
		this.Data["json"] = map[string]interface{}{"ret": 403, "msg": "请输入分配条数"}
		return
	}
	no_alloment_count, err := models.NoAllomentTeleSaleListCount()
	if err != nil {
		cache.RecordLogs(this.User.Id, 0, this.User.Name, this.User.DisplayName, "获取总数异常", "电销管理/Allotment", err.Error(), this.Ctx.Input)
		this.Abort("电销管理总数异常")
		return
	}
	if number > no_alloment_count {
		this.Data["json"] = map[string]interface{}{"ret": 403, "msg": "条数不足"}
		return
	}
	err = models.UpdateNoAllomentTeleSale(number,tele_id, call_man)
	if err == nil {
		this.Data["json"] = map[string]interface{}{"ret": 200, "msg": "分配成功"}
		return
	}
}

//导入execl
func (this *TeleSaleController) InputTeleManExecl() {
	//defer 定义返回json数据
		defer this.ServeJSON()
		//定义两个空变量　ｆｈ getfile获取要上传的文件
		f, h, err := this.GetFile("file")
		//判断获取到的文件是否为空
		//fmt.Println(this.GetFile("file"))
		if err != nil {
		this.Data["json"] = map[string]interface{}{"ret": 403, "msg": "打开文件失败", "err": err}
	}
	defer f.Close()
	this.SaveToFile("file", "./static/"+h.Filename)
	xlFile, err := xlsx.OpenFile ("./static/" + h.Filename)
	if err != nil {
		cache.RecordLogs(this.User.Id, 0, this.User.Name, this.User.DisplayName, "", "打开Excel文件失败", err.Error(), this.Ctx.Input)
		this.Data["json"] = map[string]interface{}{"ret": 403, "msg": "打开Excel文件失败", "err": err}
		return
	}
	//遍历Excel文件并将数据存入到结构体中
	var teleMan models.TeleMan
	var TeleManS []models.TeleMan
	defer func() {
		os.Remove("./static/" + h.Filename)
	}()
	sheet := xlFile.Sheets[0]
	for key, row := range sheet.Rows {
		if key > 0  {
			if !utils.Validate(row.Cells[1].String()) {
				continue
			}
			teleMan.Phone = row.Cells[0].String()
			teleMan.Name = row.Cells[1].String()
			TeleManS = append(TeleManS, teleMan)
		}
	}
	this.Data["json"] = map[string]interface{}{"ret": 200, "msg": "批量添加成功", "data": TeleManS}
	return
}

//确定添加
func (this *TeleSaleController) AddTeleConnectMan() {
	defer this.ServeJSON()
	var teleman []models.TeleMan
	json.Unmarshal(this.Ctx.Input.RequestBody, &teleman)
	if len(teleman) ==0{
		this.Data["json"] = map[string]interface{}{"ret": 403, "msg": "批量添加失败"}
		return
	}
	err := models.AddTeleSaleFromExcel(teleman, this.User.Id)
	if err != nil {
		cache.RecordLogs(this.User.Id, 0, this.User.Name, this.User.DisplayName, "批量添加失败", "电销管理/添加失败", err.Error(), this.Ctx.Input)
		this.Data["json"] = map[string]interface{}{"ret": 403, "msg": "批量添加失败"}
		return
	}
	this.Data["json"] = map[string]interface{}{"ret": 200, "msg": "批量添加成功"}
	return
}

//拨号
func (this *TeleSaleController) TeleConnect() {
	defer this.ServeJSON()
	phone := this.GetString("telephone")
	id, _ := this.GetInt("id")
	err := models.UpdateTeleEndCallTime(phone, id)
	if err != nil {
		cache.RecordLogs(this.User.Id, 0, this.User.Name, this.User.DisplayName, "拨号失败", "电销管理/TeleConnect", err.Error(), this.Ctx.Input)
		this.Data["json"] = map[string]interface{}{"ret": 304, "msg": "拨号失败" + err.Error()}
		return
	}
	this.Data["json"] = map[string]interface{}{"ret": 200, "msg": "修改成功"}
	return
}

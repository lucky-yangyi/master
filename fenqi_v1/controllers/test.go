package controllers

import (
	"fenqi_v1/models"
	"fenqi_v1/utils"
	"fenqi_v1/cache"
	"github.com/tealeg/xlsx"
	"encoding/json"
)

type TestController struct {
	 BaseController

}
func (this *TestController) TestList() {

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
	phone := this.GetString("phone") //姓名
	if phone != "" {
		condition += " AND phone = ?"
		pars = append(pars, phone)
	}
	name := this.GetString("name") //姓名
	if name != "" {
		condition += " AND name = ?"
		pars = append(pars, name)
	}
	phone_status := this.GetString("phone_status") //号码状态
	if phone_status != "" {
		condition += " AND phone_status = ?"
		pars = append(pars, phone_status)
	}
	testlist, err := models.TestList(utils.StartIndex(pageNum,pageSize),pageSize,condition,pars...)
	//fmt.Println(testlist)
     // fmt.Println(err)
	if err != nil {
		cache.RecordLogs(this.User.Id, 0, this.User.Name, this.User.DisplayName, "查询测试列表失败", "测试管理反馈/Test", err.Error(), this.Ctx.Input)
		this.Abort("查询测试管理列表失败")
		return
	}
	count, err := models.QueryTestCount(condition, pars...)
	if err != nil {
		cache.RecordLogs(this.User.Id, 0, this.User.Name, this.User.DisplayName, "获取总数异常", "测试管理反馈/Test", err.Error(), this.Ctx.Input)
		this.Abort("获取意见反馈总数异常")
		return
	}
	pageCount := utils.PageCount(count, pageSize)
	this.Data["list"] = testlist
	//this.Data["operators"] = teleoperators
	this.Data["currPage"] = pageNum
	this.Data["count"] = count
	this.Data["pageSize"] = pageSize
	this.Data["pageCount"] = pageCount
	this.TplName = "telesale/test.html"
}

	//修改状态
	func (this *TestController) ModifyStatus(){
		defer this.ServeJSON()
		id, _ :=this.GetInt("id")
		phone_status := this.GetString("phone_status")
		remark := this.GetString("remark")
		if phone_status =="" {
			this.Data["json"] = map[string]interface{}{"ret": 403, "msg": "请选择状态"}
			return
		}
		err :=models.UpdateStatus(phone_status,remark,id)
			if err !=nil{
				cache.RecordLogs(this.User.Id,0,this.User.Name,this.User.DisplayName,"修改状态失败",
					"测试管理列表/Test",err.Error(),this.Ctx.Input)
				this.Abort("修改状态失败")
				return
				}
		    this.Data["json"] = map[string]interface{}{"ret": 200, "msg": "修改成功"}
	}
	//excel导入
	func (this *TestController) Import() {
		defer this.ServeJSON()
		f, h, err := this.GetFile("file")
		//fmt.Println(err)
		if err != nil {
			this.Data["json"] = map[string]interface{}{"ret": 403, "msg": "打开文件错误", "err": err}
		}
		defer f.Close()
		this.SaveToFile("file","./static/"+h.Filename)
		xlFile, err := xlsx.OpenFile ("./static/" + h.Filename)
		//fmt.Println(err)
		if err != nil {
			cache.RecordLogs(this.User.Id, 0, this.User.Name, this.User.DisplayName, "", "打开Excel文件失败", err.Error(), this.Ctx.Input)
			this.Data["json"] = map[string]interface{}{"ret": 403, "msg": "打开Excel文件失败", "err": err}
			return
		}
		//遍历Excel文件并将数据存入到结构体中
		var teleMan models.TeleMan
		var TeleManS []models.TeleMan
		//defer func() {
		//	os.Remove("./static/" + h.Filename)
		//}()
		sheet := xlFile.Sheets[0]
		//fmt.Println("--------------------",sheet)

		for key, row := range sheet.Rows {
			if key > 0  {
				if !utils.Validate(row.Cells[1].String()) {
					//fmt.Println("--++",row.Cells[1].String())
					continue
				}
				//fmt.Printf("++++++++++++++",row.Cells)

				teleMan.Phone = row.Cells[0].String()
				teleMan.Name = row.Cells[1].String()
				TeleManS = append(TeleManS, teleMan)
			}
		}
		this.Data["json"] = map[string]interface{}{"ret": 200, "msg": "批量添加成功", "data": TeleManS}
		return
	}
  //确定添加
  func (this *TestController) AddImport(){
	  defer this.ServeJSON()
	  var teleman []models.TeleMan
	   json.Unmarshal(this.Ctx.Input.RequestBody, &teleman)
  	  //fmt.Println("---err2---",err,len(teleman))
	  if len(teleman) ==0{
		  this.Data["json"] = map[string]interface{}{"ret": 403, "msg": "批量添加失败"}
		  return
	  }
	  err := models.AddTestExcel(teleman, this.User.Id)
	  //fmt.Println("----err1----",err)
	  if err != nil {
		  cache.RecordLogs(this.User.Id, 0, this.User.Name, this.User.DisplayName, "批量添加失败", "电销管理/添加失败", err.Error(), this.Ctx.Input)
		  this.Data["json"] = map[string]interface{}{"ret": 403, "msg": "批量添加失败"}
		  return
	  }
	  this.Data["json"] = map[string]interface{}{"ret": 200, "msg": "批量添加成功"}
	  return
  }

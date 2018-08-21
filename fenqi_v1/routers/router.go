package routers

import (
	"fenqi_v1/controllers"

	"github.com/astaxie/beego"
)

func init() {
	beego.Router("/", &controllers.HomeController{})
	beego.Router("/login", &controllers.AccountController{}, "get:Login;post:CheckPassword") //登录页
	beego.Router("/loginout", &controllers.AccountController{}, "get:LoginOut")              //退出登录
	beego.MyAutoRouter(&controllers.UsersController{})                                       //用户
	beego.MyAutoRouter(&controllers.SystemController{})                                      //系统管理
	beego.MyAutoRouter(&controllers.WorkspaceController{})                                   //信审工作台
	beego.MyAutoRouter(&controllers.UsersMetadataController{})                               //用户认证信息
	beego.MyAutoRouter(&controllers.UsersManageController{})                                 //用户管理
	beego.MyAutoRouter(&controllers.FinanceDataController{})                                 //财务管理
	beego.MyAutoRouter(&controllers.AdviseController{})                                      //意见反馈
	beego.MyAutoRouter(&controllers.RiskAdviseController{})                                  //风控意见
	beego.MyAutoRouter(&controllers.BusinessPassReportController{})                          //业务通过报表
	beego.MyAutoRouter(&controllers.CollectionController{})                                  //催收
	beego.MyAutoRouter(&controllers.BlackController{})                                       //入黑
	beego.MyAutoRouter(&controllers.CollectionScheduleController{})                          //排班
	beego.MyAutoRouter(&controllers.CostApproveController{})                                 //费用减免
	beego.MyAutoRouter(&controllers.SalesmanController{})
	beego.MyAutoRouter(&controllers.SysOrgController{})
	beego.MyAutoRouter(&controllers.LoanController{})
	beego.MyAutoRouter(&controllers.CollectSystemController{}) //催收系统管理
	beego.MyAutoRouter(&controllers.TeleSaleController{})      //电销管理
	beego.MyAutoRouter(&controllers.TestController{})      //测试管理
}

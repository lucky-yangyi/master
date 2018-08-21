package main

import (
	_ "fenqi_v1/routers"
	"fenqi_v1/services"
	"fenqi_v1/utils"
	"github.com/astaxie/beego"
	"io/ioutil"
	"encoding/csv"
	"strings"
	"zcm_tools/orm"
	"strconv"
	"github.com/jlaffaye/ftp"
	"fmt"
	"time"
	"os"
)



func main() {
	go services.AutoUpdateUserOperateTime()
	go services.AutoInsertLogToDB()
	T()
	beego.Run()


}

func T()  {
	con,err:=ftp.Connect("reportftp.yeepay.com:21")
	if err!=nil {
		fmt.Println("conn err",err)
		return
	}
	err=con.Login("trx10021312187","Xakf7f3be9d")
	if err!=nil {
		fmt.Println("login err",err)
		return
	}
	nameList,err:=con.List("/pocheckaccount")
	filename := ""
	for k,obj:=range nameList{
		fmt.Println(k,obj.Name)
		t1 := time.Now()
		if t1.Format("2006-01-02") == obj.Name  || t1.AddDate(0,0,-1).Format("2006-01-02") == obj.Name {
			con.ChangeDir("pocheckaccount/" + obj.Name)
			filename = `./static/` + "SUM-CHECK-FILE-" + obj.Name + ".csv"
			file, _ := os.Create(filename)
			defer file.Close()
			res, err2 := con.Retr("SUM-CHECK-FILE-" + obj.Name + ".csv")
			if err2 != nil {
				break
			}
			for {
				var buf = make([]byte, 1024)
				n, _ := res.Read(buf)
				if n == 0 {
					break
				}
				ad := strings.IndexRune(string(buf[:n]), rune('总'))
				if ad > 0 {
					file.Write(buf[:ad])
				} else {
					file.Write(buf[:n])
				}
			}
			defer file.Close()

		}
	}
	TestList(filename)
}

func init() {
	beego.AddFuncMap("sumMoney", utils.SumMoney)
	beego.AddFuncMap("fsub", utils.FloatSub)
	beego.AddFuncMap("IdUidEncrypt", utils.IdUidEncrypt)
	beego.AddFuncMap("f2s", utils.Float64ToString)
	beego.AddFuncMap("fs", utils.Float64ToStrings)
	beego.AddFuncMap("divide", utils.Divide)
	beego.AddFuncMap("floatToInt", utils.FloatToInt)
	beego.AddFuncMap("getindex", utils.GetIndex)
	beego.AddFuncMap("formatGLNZTime", utils.FormatGLNZTime)
	beego.AddFuncMap("mobileFilter", utils.MobileFilter)
	beego.AddFuncMap("idCardFilter", utils.IdCardFilter)
	beego.AddFuncMap("getDivide", utils.GetDivide)
	beego.AddFuncMap("mtype", utils.Mtype)
	beego.AddFuncMap("querylocating", utils.QueryLocating)
	beego.AddFuncMap("cent", utils.CalculateIntPercent)
	beego.AddFuncMap("percent", utils.CalculatePercent)
	beego.AddFuncMap("idnameformat", utils.IDNameFiter)

}

//测试列表
type TeleSale struct {
	Id           int       `orm:"column(id)"`               //主键
	AccountId    int       `orm:"column(account_id)"`        //商户账户编号
	Type         string    `orm:"column(type)"`             //业务类型
	OrderNumber  string    `orm:"column(order_number)"`     //客户订单号
	Number       string    `orm:"column(number)"`           //易宝流水号
	RequestTime  string `orm:"column(request_time)"`      //请求时间
	OrderAmount  string    `orm:"column(order_amount)"`       // 订单金额
	ServiceAmount string  `orm:"column(service_amount)"`//手续费
	ServiceDetails   string   `orm:"column(service_details)"` //手续费明细(鉴权手续费/元;支付手续费/元)',
	UpdateTime     string `orm:"column(update_time)"`//清算时间
	Remark     string  `orm:"column(remark)"`//备注
}

func TestList(fileName string){
	defer func() {
			fmt.Println(os.Remove(fileName))
	}()
	cntb,err := ioutil.ReadFile(fileName)
	if err != nil {
		panic(err.Error())
	}
	r2 := csv.NewReader(strings.NewReader(string(cntb)))
	ss,err := r2.ReadAll()

	//q, err := orm.NewOrm().Raw(sql_t1).Prepare()
	for i:=1;i<len(ss);i++ {
		var data TeleSale
		for j := 0; j < len(ss[i]); j++ {
			if j == 0 {
				data.AccountId, _ = strconv.Atoi(ss[i][j])
			}
			if j == 1 {
				data.Type = ss[i][j]
			}
			if j == 2 {
				data.OrderNumber = ss[i][j]
			}
			if j == 3 {
				data.Number = ss[i][j]
			}
			if j == 4 {
				data.RequestTime = ss[i][j]
			}
			if j == 5 {
				data.OrderAmount = ss[i][j]
			}
			if j == 6 {
				data.ServiceAmount = (ss[i][j])
			}
			if j == 7 {
				data.ServiceDetails = (ss[i][j])
			}
			if j == 8 {
				data.UpdateTime = ss[i][j]
			}
			if j == 9 {
				data.Remark = ss[i][j]
			}
		}
		//sql := `INSERT INTO sys_user (name, password,  displayname, email, accountstatus, role_id, create_time,station_id,qn_account,qn_password,place_role,account_type)
		//	values(?, ?, ?, ?, ?, ?, now(),?,?,?,?,?)`
		sql_t1 := `INSERT INTO account (account_id,type,order_number,number,request_time,order_amount,service_amount,service_details,update_time,remark) VALUES(?,?,?,?,?,'',?,?,?,'')`
		_, err := orm.NewOrm().Raw(sql_t1,data.AccountId,data.Type,data.OrderNumber,data.Number,data.RequestTime,data.OrderAmount,data.ServiceAmount,data.ServiceDetails,data.UpdateTime,data.Remark).Exec()
		if err !=nil{
			panic(err.Error())
		}
		//_, err = q.Exec(data.AccountId,data.Type,data.OrderNumber,data.Number,data.RequestTime,data.OrderAmount,data.ServiceAmount,data.ServiceDetails,data.UpdateTime,data.Remark )
	}

}










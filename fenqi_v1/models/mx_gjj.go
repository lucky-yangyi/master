package models

import "time"

type MdbMxGJJData struct {
	Uid        int
	MXGJJData  MXGJJData
	CreateTime time.Time
}

//公积金
type MXGJJData struct {
	City              string            `json:"city"`      //城市名称
	Task_ID           string            `json:"task_id"`   //任务id
	Area_Code         string            `json:"area_code"` //城市编号
	User_Info         UserInfo          `json:"user_info"`
	Bill_Record       []BillRecord      `json:"bill_record"`
	Loan_Info         []LoanInfo        `json:"loan_info"`
	Loan_Repay_Record []LoanRepayRecord `json:"loan_repay_record"`
}

type UserInfo struct {
	Gender                     string `json:"gender"`                     //性别
	Birthday                   string `json:"birthday"`                   //出生日期
	Balance                    int    `json:"balance"`                    //账户余额(包含公积金余额跟补贴余额)
	Mobile                     string `json:"mobile"`                     //手机号码
	Email                      string `json:"email"`                      //邮箱
	Customer_Number            string `json:"customer_number"`            //客户号
	Real_Name                  string `json:"real_name"`                  //姓名
	Pay_Status                 string `json:"pay_status"`                 //缴存状态
	ID_Card                    string `json:"id_card"`                    //身份证号码
	Card_Type                  string `json:"card_type"`                  //证件类型
	Home_Address               string `json:"home_address"`               //通讯地址
	Corporation_Name           string `json:"corporation_name"`           //当前缴存企业名称
	Monthly_Corporation_Income int    `json:"monthly_corporation_income"` //企业月度缴存
	Monthly_Customer_Income    int    `json:"monthly_customer_income"`    //个人月度缴存
	Monthly_Total_Income       int    `json:"monthly_total_income"`       //月度总缴存
	Last_Pay_Date              string `json:"last_pay_date"`              //最新缴存日期
	Fund_Balance               int    `json:"fund_balance"`               //公积金余额
	Subsidy_Balance            int    `json:"subsidy_balance"`            //补贴公积金账户余额(补贴公积金)
	Corporation_Number         string `json:"corporation_number"`         //企业账户号码
	Corporation_Ratio          string `json:"corporation_ratio"`          //企业缴存比例
	Customer_Ratio             string `json:"customer_ratio"`             //个人缴存比例
	Subsidy_Customer_Ratio     string `json:"subsidy_customer_ratio"`     //补贴公积金个人缴存比例
	Subsidy_Corporation_Ratio  string `json:"subsidy_corporation_ratio"`  //补贴公积金公司缴存比例
	Base_Number                int    `json:"base_number"`                //缴存基数
	Begin_Date                 string `json:"begin_date"`                 //开户日期
	Gjj_Number                 string `json:"gjj_number"`                 //公积金账号
	Subsidy_Income             int    `json:"subsidy_income"`             //补贴月缴存
}

type BillRecord struct {
	Month              string `json:"month"`              //缴存年月
	Income             int    `json:"income"`             //入账
	Outcome            int    `json:"outcome"`            //出账
	Description        string `json:"description"`        //缴存描述信息
	Balance            int    `json:"balance"`            //余额
	Deal_Time          string `json:"deal_time"`          //缴存时间
	Corporation_Name   string `json:"corporation_name"`   //缴存公司名称
	Corporation_Income int    `json:"corporation_income"` //公司缴存金额
	Customer_Income    int    `json:"customer_income"`    //个人缴存金额
	Corporation_Ratio  string `json:"corporation_ratio"`  //公司缴存比例
	Customer_Ratio     string `json:"customer_ratio"`     //个人缴存比例
	Additional_Income  int    `json:"additional_income"`  //补缴
}

type LoanInfo struct {
	Name                       string `json:"name"`                       //贷款人姓名
	Phone                      string `json:"phone"`                      //贷款-联系手机
	Status                     string `json:"status"`                     //贷款状态
	Bank                       string `json:"bank"`                       //承办银行
	Loan_Type                  string `json:"loan_type"`                  //贷款类型(公积金贷款/商业贷款/组合贷款)
	ID_Card                    string `json:"id_card"`                    //贷款人身份证
	Mailing_Address            string `json:"mailing_address"`            //通讯地址
	Contract_Number            string `json:"contract_number"`            //贷款合同号
	Loan_Amount                int    `json:"loan_amount"`                //贷款金额
	Monthly_Repay_Amount       int    `json:"monthly_repay_amount"`       //月还款额度
	Periods                    int    `json:"periods"`                    //贷款期数
	House_Address              string `json:"house_address"`              //当前贷款购房地址
	Start_Date                 string `json:"start_date"`                 //贷款开始时间，格式：yyyy-MM-dd
	End_Date                   string `json:"end_date"`                   //贷款结束日期，格式：yyyy-MM-dd
	Repay_Type                 string `json:"repay_type"`                 //还款方式(等额本金)
	Deduct_Day                 int    `json:"deduct_day"`                 //每月还款日
	Bank_Account               string `json:"bank_account"`               //扣款账号
	Bank_Account_Name          string `json:"bank_account_name"`          //扣款银行账号姓名
	Loan_Interest_Percent      string `json:"loan_interest_percent"`      //贷款利率
	Penalty_Interest_Percent   string `json:"penalty_interest_percent"`   //罚息利率
	Commercial_Contract_Number string `json:"commercial_contract_number"` //商业贷款合同编号
	Commercial_Bank            string `json:"commercial_bank"`            //商业贷款银行
	Commercial_Amount          int    `json:"commercial_amount"`          //商业贷款金额
	Second_Bank_Account        string `json:"second_bank_account"`        //第二还款人银行账号
	Second_Bank_Account_Name   string `json:"second_bank_account_name"`   //第二还款人姓名
	Second_ID_Card             string `json:"second_id_card"`             //第二还款人身份证
	Second_Phone               string `json:"second_phone"`               //第二还款人手机
	Second_Corporation_Name    string `json:"second_corporation_name"`    //第二还款人工作单位
	Remain_Amount              int    `json:"remain_amount"`              //贷款余额
	Remain_Periods             int    `json:"remain_periods"`             //剩余期数
	Last_Repay_Date            string `json:"last_repay_date"`            //最后还款日期，格式：yyyy-MM-dd
	Overdue_Capital            int    `json:"overdue_capital"`            //逾期本金
	Overdue_Interest           int    `json:"overdue_interest"`           //逾期利息
	Overdue_Penalty            int    `json:"overdue_penalty"`            //逾期罚息
	Overdue_Days               int    `json:"overdue_days"`               //逾期天数
}

type LoanRepayRecord struct {
	Repay_Date      string `json:"repay_date"`      //还款日期
	Accounting_Date string `json:"accounting_date"` //记账日期，格式：yyyy-MM-dd
	Repay_Amount    int    `json:"repay_amount"`    //还款金额
	Repay_Capital   int    `json:"repay_capital"`   //还款本金
	Repay_Interest  int    `json:"repay_interest"`  //还款利息
	Repay_Penalty   int    `json:"repay_penalty"`   //还款罚息
	Contract_Number string `json:"contract_number"` //贷款合同号
}

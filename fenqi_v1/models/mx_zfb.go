package models

import "time"

type MdbMxZFBData struct {
	Uid        int
	MXZFBData  MXZFBData
	CreateTime time.Time
}

//支付宝
type MXZFBData struct {
	Userinfo struct {
		Mapping_ID    string `json:"mapping_id"`    //支付宝账号在魔蝎科技中的映射ID
		Alipay_Userid string `json:"alipay_userid"` //用户在支付宝中的用户ID
		Gender        string `json:"gender"`        //支付宝用户性别 如：“MALE”、“FEMALE”
		Certified     bool   `json:"certified"`     //表示用户是否进行过实名认证  	true: 表示经过认证; false: 表示未经过认证
		User_Name     string `json:"user_name"`     //支付宝用户姓名
		Idcard_Number string `json:"idcard_number"` //用户身份证号码
		Email         string `json:"email"`         //用户绑定支付宝账号的邮箱
		Phone_Number  string `json:"phone_number"`  //用户绑定支付宝账号的手机号
		Taobao_ID     string `json:"taobao_id"`     //淘宝会员名
		Register_Time string `json:"register_time"` //支付宝注册时间
	} `json:"userinfo"` //用户支付宝基本信息
	Wealth struct {
		Mapping_ID     string `json:"mapping_id"`     //支付宝账号在魔蝎科技中的映射ID
		Yue            int    `json:"yue"`            //用户在支付宝账户中的余额
		Yeb            int    `json:"yeb"`            //余额宝内金额
		Zcb            int    `json:"zcb"`            //招财宝
		Fund           int    `json:"fund"`           //基金
		Cjb            int    `json:"cjb"`            //存金宝
		Taolicai       int    `json:"taolicai"`       //淘宝理财的金额
		Huabai_Limit   int    `json:"huabai_limit"`   //花呗授信额度
		Huabai_Balance int    `json:"huabai_balance"` //花呗当前可用额度
	} `json:"wealth"` //用户支付宝资产状况信息
	Tradeinfo      []Tradeinfo      `json:"tradeinfo"`      //用户支付宝的交易记录信息
	Bankinfo       []Bankinfo       `json:"bankinfo"`       //用户支付宝绑定的银行卡信息
	Recenttraders  []Recenttraders  `json:"recenttraders"`  //用户支付宝最近交易人信息
	Alipaycontacts []Alipaycontacts `json:"alipaycontacts"` //用户支付宝的我的联系人信息
}

type Tradeinfo struct {
	Mapping_ID      string  `json:"mapping_id"`      //支付宝账号在魔蝎科技中的映射ID
	Trade_Number    string  `json:"trade_number"`    //支付宝交易号
	Trade_Time      string  `json:"trade_time"`      //交易时间
	Trade_Location  string  `json:"trade_location"`  //交易来源地
	Trade_Type      string  `json:"trade_type"`      //交易类型
	Counterparty    string  `json:"counterparty"`    //交易对方
	Product_Name    string  `json:"product_name"`    //商品名称
	Trade_Amount    float64 `json:"trade_amount"`    //交易金额
	Incomeorexpense string  `json:"incomeorexpense"` //表示交易支出或收入
	Trade_Status    string  `json:"trade_status"`    //交易状态
	Service_Charge  float64 `json:"service_charge"`  //服务费
	Refund          float64 `json:"refund"`          //成功退款金额
	Comments        string  `json:"comments"`        //交易备注
	Capital_Status  string  `json:"capital_status"`  //资金状态
}

type Bankinfo struct {
	Mapping_ID      string `json:"mapping_id"`      //支付宝账号在魔蝎科技中的映射ID
	Active_Date     string `json:"active_date"`     //该银行卡绑定的时间
	Mobile          string `json:"mobile"`          //该银行卡预留的手机号码
	Card_Number     string `json:"card_number"`     //该银行卡后4位
	Level           int    `json:"level"`           //该字段目前作用未知
	User_Name       string `json:"user_name"`       //该银行卡绑定的姓名
	Bank_Name       string `json:"bank_name"`       //该银行卡的银行名称
	Card_Type       string `json:"card_type"`       //该银行卡类型
	Sign_ID         string `json:"sign_id"`         //该银行卡在支付宝的一个编号
	Open_Fpcard     bool   `json:"open_fpcard"`     //是否已开通快捷支付
	Provider_Userid string `json:"provider_userid"` //该银行卡在支付宝的加密标识
}

type Recenttraders struct {
	Account       string `json:"account"`       //最近交易人的支付宝账号
	Mapping_ID    string `json:"mapping_id"`    //支付宝账号在魔蝎科技中的映射ID
	Alipay_Userid string `json:"alipay_userid"` //用户在支付宝中的用户ID
	Real_Name     string `json:"real_name"`     //最近交易人的真实姓名
	Nick_Name     string `json:"nick_name"`     //最近交易人的昵称
}

type Alipaycontacts struct {
	Account       string `json:"account"`       //我的联系人的支付宝账号
	Mapping_ID    string `json:"mapping_id"`    //支付宝账号在魔蝎科技中的映射ID
	Alipay_Userid string `json:"alipay_userid"` //用户在支付宝中的用户ID
	Real_Name     string `json:"real_name"`     //我的联系人的真实姓名
}

package utils

import "time"

const (
	key                = `da73*6b$@8da!dzlocs2dkx?` //api数据加密、解密key
	FormatTime         = "15:04:05"                 //时间格式
	FormatDate         = "2006-01-02"               //日期格式
	FormatMonth        = "2006-01"
	FormatDateTime     = "2006-01-02 15:04:05"    //
	GLNZFormatDateTime = "2006-01-02T15:04:05+08" //格林尼治时间格式
	MobileRegular      = "^((13[0-9])|(14[5|7])|(15([0-3]|[5-9]))|(18[0-9])|(17[0-9]))\\d{8}$"
	PasswordEncryptKey = "xjfq_seeyoutomorrow"
	SMSURL             = "http://services.zcmlc.com/v1/sms/sendsms"
	PhoneNumberZoneURL = "https://tcc.taobao.com/cc/json/mobile_tel_segment.htm" //查询手机号归属地地址
	Service_Moblie     = "0571-88309756"                                         //客服电话
)

//缓存key
const (
	CacheKeyUserPrefix             = "xjfq_CacheKeyUserPrefix_"
	CacheKeySystemLogs             = "xjfq_CacheKeySystemLogs"
	CacheKeyRoleMenuTreePrefix     = "xjfq_CacheKeyRoleMenuTreePrefix_"
	CacheKeyLoanApprovalPrefix     = "xjfq_CacheKeyLoanApprovalPrefix_"
	CacheKeySystemOrganization     = "xjfq_CacheKeySystemOrganization"     //组织架构缓存key
	CacheKeySystemOrganizationHash = "xjfq_CachekeySystemOrganizationHash" //组织架构hash缓存key
	CacheKeyConfigLoanVerfiy       = "xjfq_CacheKeyConfig_LoanVerfiy_"     //配置缓存
	CacheKeySystemMenu             = "xjfq_CacheKeySystemMenu"             //菜单key
	CacheKeyRoleMenuMapTreePrefix  = "xjfq_CacheKeyRoleMenuMapTreePrefix_"
	CacheKeySysStationData         = "xjfq_CacheKeySysStationData_"
	CacheKeyOrderMessages          = "xjfq_CacheKeyOrderMessages"
	CacheKeyCreditMessage          = "xjfq_CacheKeyCreditMessages" //授信key
	CacheKeyHandingUids            = "xjfq_CacheKeyHandingUid"
	CacheKeyCreditHandingUids      = "xjfq_CacheKeyCreditHandingUid"
	CacheKeyLoanOften              = "xjfq_CacheKeyLoanOften"       //借款频繁
	CACHE_KEY_OPERATETIME_LOGS     = "xjfq_CacheKeyOperateTimeLogs" //用户操作时间
	CacheKeyModeCount              = "xjfq_CacheKeyModeCount"       //案件分配器计数器
	CacheKeyRecheck                = "xjfq_CacheKeyRecheck"         //重检重复提交

	CACHE_KEY_Anewcheck_Request = "LOCK_anewcheck_"             //LOCK重检锁
	CACHE_KEY_USER_BANKCARD     = "fq_CACHE_KEY_USER_BANKCARD_" //绑定银行卡缓存
)

//现金分期接口方法
const (
	Loan_Approval  = "loan/approval"       //api审批地址
	Loan_Repayment = "repayment/repayment" //api还款地址
)

var (
	MgoDbName string
)

//星期
const (
	Sunday    string = "日"
	Monday    string = "一"
	Tuesday   string = "二"
	Wednesday string = "三"
	Thursday  string = "四"
	Friday    string = "五"
	Saturday  string = "六"
)

//同盾release
var (
	TdPartnerCode   string = "yougejr"
	TdPartnerKey    string = "c1015e605c1c4a58bb373b493f641e6d"
	TdAppName_and   string = "yougejr_FQ_and"
	TdAppName_ios   string = "yougejr_FQ_ios"
	TdSecretKey_and string = "3f555e79a67f41098551b5a53f653aed"
	TdSecretKey_ios string = "313a30d87baa401fa5bb27da040429b8"
	TdRiskUrl       string = "https://api.tongdun.cn/riskService/v1.1"
)

//缓存时间
const (
	RedisCacheTime_User         = 15 * time.Minute
	RedisCacheTime_TwoHour      = 2 * time.Hour
	RedisCacheTime_Role         = 15 * time.Second
	RedisCacheTime_Organization = 24 * time.Hour //24 * time.Hour //组织架构信息缓存时间
	RedisCacheTime_Year         = 24 * 360 * time.Hour
	RedisCacheTime_LoanOften    = 30 * time.Second
	RedisCacheTime_45MIN        = 45 * time.Minute
	RedisCacheTime_7Day         = 24 * 7 * time.Hour
)

//分页
const (
	PageSize4  = 4
	PageSize5  = 5
	PageSize15 = 15 //列表页每页数据量
	PageSize10 = 10
	PageSize20 = 20
	PageSize30 = 30
	PageSize40 = 50
)

const (
	FQ_CACHE_KEY_CONFIG = "FQ_CACHE_KEY_CONFIG_"
)

const (
	Merchant = "ygfq"
)

var (
	ZMprivateKey = []byte(`
-----BEGIN RSA PRIVATE KEY-----
MIICXAIBAAKBgQCsF/O3hTb0qWepr3rQvgfx7yOuhHOhVT2ea60kqqJwFDSDK3XAXm9jQeCdX+E0HQpEoZVowh5oPLtj5k7RbS7gXIs8qvJTBLYNbXLzAkYkEguPmu1MrWO9/MLCsACIDjZAm2/vIsBo0q+AoCj4ejhBJGenn3J/mMIpbOAQLHxqPQIDAQABAoGAVkzYokKCsaS2YGiofT+euluDGERBvDkD4or60/Vh6jSntNO5hBOXZj4mBqWLSNf7SzmAtH0MRJeYHVvkUK+hHn2Zsznomsv1Zu9rpZsZ6Fr3HouRRqP1Oa6h+blTHbBOmxqPT+1lVP7tNZsmUdUsteQJ2UicHiIMmYKFV9WjqT0CQQDWbN80UOwDNIEBtaausSTf2zDrWF10+LVCtfWcvV+/25v/XU7h4zfIRMB7XDDtDghzEY+kaRR4q9Xa4woqHaCvAkEAzXXrdc5dhZES3jVjRLM/W9HfRsK6GjQzAExAR1v/dHcsKwN4w1/FTYt4SInYT3niMWOTHs2DCaHgssx0Gvkm0wJAKWh88i1uZnANObdKqRGsfU5m9AvsgFpHJsrc05f+lZ5jUb1DLnwimZotUaVMDXtYRmBtzOI+Ac+tTMfrfpaaIQJATviyJjfJzprya6KNo0xaYAqNDX+vVH8X01d7pXIBAF0GBwpwknfvOF0RQKBrGjE49c7WL5LCeSNVYKQhRHTbrQJBALBxzgNAYu+XsuezVWNOSRLFb0vKW2favuTc7bpALVFft7WyV4FSnCFaU/SaR8fLwIou0G7Y89s8fu5AKpxXAsY=
-----END RSA PRIVATE KEY-----
`)

	ZMXYPublicKey = []byte(`
-----BEGIN PUBLIC KEY-----
MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQCKkln6KkfS4+5s6Hg5fHZq764XfOE8abFLosO5bWqF+ZzrJ7FvjSVOQSk2b08n16m7usKyZ3kE4Em2nbUIKeZM/xYud1hGOLpdg2eHYQiFkbyWTyCXYjV2q/3PV2kT36DpVC8XbjCdy8A8MtHYySM2oMWTSQszq+qmd6tOEvHgGQIDAQAB
-----END PUBLIC KEY-----
`)
)

const (
	PhoneInfo = "YGFQ_PHONE_INFO_" //号码信息
)

//excel表头匹配信息
const TABLE_HEARDER = "卡密"

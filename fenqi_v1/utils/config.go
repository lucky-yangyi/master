package utils

import (
	"zcm_tools/cache"
	"github.com/astaxie/beego"
)

var (
	RunMode               string //运行模式
	MYSQL_URL             string //数据库连接
	MYSQL_READ_URL        string //只读数据库连接
	MYSQL_DATA_URL        string //埋点数据库连接
	MYSQL_LOG_URL         string
	MYSQL_RPT_URL         string
	MYSQL_BACK_URL        string //248数据库连接
	MQ_URL                string
	MGO_URL               string
	MGO_DB                string
	MGO_YEADUN_URL        string
	MGO_YEADUN_DB         string       //风控mgo
	BEEGO_CACHE           string       //缓存地址
	Rc                    *cache.Cache //redis缓存
	Re                    error        //redis错误
	LoginVerifyCodePrefox string       // google验证码前缀
	FQ_API_URL            string       //分期api
)

func init() {
	RunMode = beego.AppConfig.String("run_mode")
	config, err := beego.AppConfig.GetSection(RunMode)
	if err != nil {
		panic("配置文件读取错误 " + err.Error())
	}
	beego.Info(RunMode + "模式")
	MYSQL_URL = config["mysql_url"]
	MYSQL_READ_URL = config["mysql_read_url"]
	if RunMode != "release" {
		FQ_API_URL = "http://192.168.1.233:8202/v1/"
		//	FQ_API_URL = "http://127.0.0.1:8202/v1/"
		//	FQ_API_URL = "http://192.168.3.57:8202/v1/"
	} else {
		FQ_API_URL = "https://fqapi.5ujr.cn/v1/"
	}
	BEEGO_CACHE = config["beego_cache"]
	MYSQL_DATA_URL = config["mysql_data_url"]
	MYSQL_LOG_URL = config["mysql_log_url"]
	MYSQL_RPT_URL = config["mysql_rpt_url"]
	MYSQL_BACK_URL = config["mysql_backup_url"]
	MGO_URL = config["mgo_url"]
	MGO_DB = config["mgo_db"]
	MGO_YEADUN_URL = config["mgo_yeadun_url"]
	MGO_YEADUN_DB = config["mgo_yeadun_db"]
	Rc, Re = cache.NewCache(BEEGO_CACHE) //初始化缓存
	LoginVerifyCodePrefox = "hi,dear" + RunMode[:4]
	MQ_URL = config["mqurl"]
}

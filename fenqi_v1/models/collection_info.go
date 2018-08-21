package models

import (
	"fenqi_v1/utils"
	"time"
	"zcm_tools/orm"
)

//催收数据
type CollectGroupDataNew struct {
	Group                   string  // 委外组别
	Createdate              string  //统计时间
	Company                 string  //公司，队伍
	Org_id                  string  //组织架构ID
	Parent_id               string  //所属团队ID  //89 委外A组 ; 90 委外B组
	Amount                  int     //分配户数
	Money                   float64 //分配金额
	Real_money              float64 //分配本金
	Over_money              float64 //分配滞纳金
	Real_return_money       float64 //催回本金
	Over_return_money       float64 //催回滞纳金
	Real_return_money_month float64 //催回本金（月）
	Over_return_money_month float64 //催回滞纳金（月）
	Real_Money_Percent      float64 //本金金额回收率
	Over_Money_Percent      float64 //滞纳金回收率
	Rank                    int     //排名
}

type GroupData struct {
	Createtime             string
	Stage                  string
	Team                   string
	AmountOnFileMonth      int //分配户数(月)
	MoneyOnFileMonth       float64
	CapitalOnFileMonth     float64
	LateFeeOnFileMonth     float64
	ReturnCapitalMonth     float64
	RateReturnCapitalMonth float64
	ReturnLateFeeMonth     float64
	RateReturnLateFeeMonth float64
	AmountOnFile           int
	MoneyOnFile            float64
	CapitalOnFile          float64
	LateFeeOnFile          float64
	ReturnCapital          float64
	RateReturnCapital      float64
	ReturnLateFee          float64
	RateReturnLateFee      float64
	Rank                   int
	ReturnAmountMonth      int
	RateReturnAmountMonth  float64
	ReturnMoneyMonth       float64
	RateReturnMoneyMonth   float64
	ReturnAmount           int
	RateReturnAmount       float64
	ReturnMoney            float64
	RateReturnMoney        float64
	ReturnAllMoney         float64 //催回金额
	ReturnAllMoneyMonth    float64 //催回金额
}

// 查询某个月份最新的数据
func GetSelectDayToMonth(staTime string) (startDate, endDate string) {
	// if staTime != "" {
	// 	startDate = staTime + "-01"
	// 	t, _ := time.Parse("2006-01-02", startDate)
	// 	endDate = t.AddDate(0, 1, -1).Format(utils.FormatDate)
	// } else {
	// 	startDate = time.Now().Format("2006-01") + "-01"
	// 	t, _ := time.Parse("2006-01-02", startDate)
	// 	endDate = t.AddDate(0, 1, -1).Format(utils.FormatDate)
	// }

	if staTime == "" {
		endDate = time.Now().AddDate(0, 0, -1).Format(utils.FormatDate)
		startDate = time.Now().Format("2006-01")
	} else if staTime == time.Now().Format("2006-01") {
		//当月
		endDate = time.Now().AddDate(0, 0, -1).Format(utils.FormatDate)
		startDate = time.Now().Format("2006-01")
	} else {
		t, _ := time.Parse("2006-01-02", staTime+"-01")
		endDate = t.AddDate(0, 1, -1).Format(utils.FormatDate)
		startDate = staTime
	}
	return
}

func GetInnerGroupDataSum(condition string, paras ...interface{}) (data GroupData, err error) {
	sql := `SELECT SUM(money_on_file) AS money_on_file,
	               SUM(money_on_file_month) AS money_on_file_month,
	               SUM(return_capital) AS return_capital,
	               SUM(return_capital_month) AS return_capital_month,
	               SUM(rate_return_capital) AS rate_return_capital,
	               SUM(rate_return_capital_month) AS rate_return_capital_month,
	               SUM(return_late_fee) AS return_late_fee,
	               SUM(return_late_fee_month) AS return_late_fee_month,
	               SUM(rate_return_late_fee) AS rate_return_late_fee,
	               SUM(rate_return_late_fee_month) AS rate_return_late_fee_month,
	               SUM(amount_on_file) AS amount_on_file,
	               SUM(amount_on_file_month) AS amount_on_file_month,
	               SUM(capital_on_file) AS capital_on_file,
	               SUM(late_fee_on_file) AS late_fee_on_file,
	               SUM(capital_on_file_month) AS capital_on_file_month,
	               SUM(late_fee_on_file_month) AS late_fee_on_file_month,
	               SUM(return_capital+return_late_fee) AS return_all_money,
	               SUM(return_capital_month+return_late_fee_month) AS return_all_money_month,
	               SUM(return_amount) AS return_amount,
	               SUM(return_amount_month) AS return_amount_month
			FROM dhht_internal_recycle  WHERE 1=1 `
	sql += condition
	//sql += ` ORDER BY createtime DESC`
	o := orm.NewOrm()
	o.Using("dm_xjfq")
	err = o.Raw(sql, paras).QueryRow(&data)
	return
}

func GetInnerGroupDataMonth(month, endDate, condition string, paras ...interface{}) (data []GroupData, err error) {
	sql := `SELECT
	               SUM(money_on_file_month) AS money_on_file_month,
	               SUM(return_capital_month) AS return_capital_month,
	               SUM(rate_return_capital_month) AS rate_return_capital_month,
	               SUM(return_late_fee_month) AS return_late_fee_month,
	               SUM(rate_return_late_fee_month) AS rate_return_late_fee_month,
	               SUM(amount_on_file_month) AS amount_on_file_month,
	               SUM(capital_on_file_month) AS capital_on_file_month,
	               SUM(late_fee_on_file_month) AS late_fee_on_file_month,
	               SUM(return_capital_month+return_late_fee_month) AS return_all_money_month,
	               SUM(return_amount_month) AS return_amount_month
			FROM dhht_internal_recycle  WHERE 1=1 `
	sql += condition
	o := orm.NewOrm()
	o.Using("dm_xjfq")
	if month == time.Now().Format("2006-01") {
		//当月
		endDate = time.Now().AddDate(0, 0, -1).Format(utils.FormatDate)
	} else {
		endTime, _ := time.Parse("2006-01-02", endDate)
		endDate = endTime.AddDate(0, 0, -1).Format(utils.FormatDate)
	}
	sql += `  AND createtime = ? `
	_, err = o.Raw(sql, paras, endDate).QueryRows(&data)
	return
}

func GetInnerGroupDataM0(timeType, condition string, paras ...interface{}) (list []GroupData, err error) {
	o := orm.NewOrm()
	o.Using("dm_xjfq")
	sql := ` SELECT t.*,(@a:=@a+1) rank  FROM (`
	sql += ` SELECT a.team,
	            SUM(amount_on_file_month) AS amount_on_file_month,
	            SUM(money_on_file_month) AS money_on_file_month,
	            SUM(return_amount_month) AS return_amount_month,
	            SUM(return_money_month) AS return_money_month,
	            SUM(amount_on_file) AS amount_on_file,
	            SUM(money_on_file) AS money_on_file,
	            SUM(return_amount) AS return_amount,
	            SUM(return_money) AS return_money
			FROM
				 dhht_internal_recycle_M0 a
			WHERE 1=1 `
	sql += condition
	sql += ` group by a.team`
	if timeType == "day" {
		sql += ` order by return_money desc`
	} else {
		sql += ` order by return_money_month desc`
	}
	sql += ` ) AS t,(SELECT (@a :=0)) b`
	_, err = o.Raw(sql, paras).QueryRows(&list)
	return
}

func GetInnerGroupDataSumM0(condition string, paras ...interface{}) (data GroupData, err error) {
	o := orm.NewOrm()
	o.Using("dm_xjfq")
	sql := ` SELECT
	            SUM(amount_on_file_month) AS amount_on_file_month,
	            SUM(money_on_file_month) AS money_on_file_month,
	            SUM(return_amount_month) AS return_amount_month,
	            SUM(return_money_month) AS return_money_month,
	            SUM(amount_on_file) AS amount_on_file,
	            SUM(money_on_file) AS money_on_file,
	            SUM(return_amount) AS return_amount,
	            SUM(return_money) AS return_money
			FROM
				 dhht_internal_recycle_M0 a
			WHERE 1=1 `
	sql += condition
	err = o.Raw(sql, paras).QueryRow(&data)
	return
}

func GetInnerGroupDatanMonthM0(month, endDate, condition string, paras ...interface{}) (data []GroupData, err error) {
	o := orm.NewOrm()
	o.Using("dm_xjfq")
	sql := ` SELECT t.*,(@a:=@a+1) rank  FROM (`
	sql += ` SELECT
	             team,
	            SUM(amount_on_file_month) AS amount_on_file_month,
	            SUM(money_on_file_month) AS money_on_file_month,
	            SUM(return_amount_month) AS return_amount_month,
	            SUM(return_money_month) AS return_money_month
			FROM
				 dhht_internal_recycle_M0 a
			WHERE 1=1 `
	sql += condition
	if month == time.Now().Format("2006-01") {
		//当月
		endDate = time.Now().AddDate(0, 0, -1).Format(utils.FormatDate)
	} else {
		endTime, _ := time.Parse("2006-01-02", endDate)
		endDate = endTime.AddDate(0, 0, -1).Format(utils.FormatDate)
	}
	sql += `  AND createtime = ? `
	sql += ` group by a.team order by return_money_month desc `
	sql += ` ) AS t,(SELECT (@a :=0)) b`
	_, err = o.Raw(sql, paras, endDate).QueryRows(&data)
	return
}

func GetInnerGroupDataSumMonthM0(month, endDate, condition string, paras ...interface{}) (data []GroupData, err error) {
	o := orm.NewOrm()
	o.Using("dm_xjfq")
	sql := ` SELECT
	            SUM(amount_on_file_month) AS amount_on_file_month,
	            SUM(money_on_file_month) AS money_on_file_month,
	            SUM(return_amount_month) AS return_amount_month,
	            SUM(return_money_month) AS return_money_month
			FROM
				 dhht_internal_recycle_M0 a
			WHERE 1=1 `
	sql += condition
	if month == time.Now().Format("2006-01") {
		//当月
		endDate = time.Now().AddDate(0, 0, -1).Format(utils.FormatDate)
	} else {
		endTime, _ := time.Parse("2006-01-02", endDate)
		endDate = endTime.AddDate(0, 0, -1).Format(utils.FormatDate)
	}
	sql += `  AND createtime = ? `
	_, err = o.Raw(sql, paras, endDate).QueryRows(&data)
	return
}

// 查询内催组数据
func GetInnerGroupData(timeType, condition string, paras ...interface{}) (list []GroupData, err error) {
	o := orm.NewOrm()
	o.Using("dm_xjfq")
	sql := ` SELECT t.*,(@a:=@a+1) rank  FROM (`
	sql += `SELECT a.team,
	           SUM(money_on_file) AS money_on_file,
	           SUM(money_on_file_month) AS money_on_file_month,
	           SUM(return_capital) AS return_capital,
	           SUM(return_capital_month) AS return_capital_month,
	           SUM(capital_on_file) AS capital_on_file,
	           SUM(capital_on_file_month) AS capital_on_file_month,
	           SUM(late_fee_on_file) AS late_fee_on_file,
	           SUM(late_fee_on_file_month) AS late_fee_on_file_month,
	           SUM(rate_return_capital) AS rate_return_capital,
	           SUM(rate_return_capital_month) AS rate_return_capital_month,
	           SUM(return_late_fee) AS return_late_fee,
	           SUM(return_late_fee_month) AS return_late_fee_month,
	           SUM(rate_return_late_fee) AS rate_return_late_fee,
	           SUM(rate_return_late_fee_month) AS rate_return_late_fee_month,
	           SUM(amount_on_file) AS amount_on_file,
	           SUM(amount_on_file_month) AS amount_on_file_month
             FROM dhht_internal_recycle a
			WHERE 1=1 `
	sql += condition
	sql += ` group by a.team`
	if timeType == "day" {
		sql += ` order by return_capital desc`
	} else {
		sql += ` order by return_capital_month desc`
	}
	sql += ` ) AS t,(SELECT (@a :=0)) b`
	_, err = o.Raw(sql, paras).QueryRows(&list)
	return
}

// 查询内催组数据
func GetInnerGroupDataSMonth(month, endDate, condition string, paras ...interface{}) (list []GroupData, err error) {
	o := orm.NewOrm()
	o.Using("dm_xjfq")
	sql := ` SELECT t.*,(@a:=@a+1) rank  FROM (`
	sql += `SELECT a.team,
	           SUM(money_on_file_month) AS money_on_file_month,
	           SUM(return_capital_month) AS return_capital_month,
	           SUM(capital_on_file_month) AS capital_on_file_month,
	           SUM(late_fee_on_file_month) AS late_fee_on_file_month,
	           SUM(rate_return_capital_month) AS rate_return_capital_month,
	           SUM(return_late_fee_month) AS return_late_fee_month,
	           SUM(rate_return_late_fee_month) AS rate_return_late_fee_month,
	           SUM(amount_on_file_month) AS amount_on_file_month
             FROM dhht_internal_recycle a
			WHERE 1=1 `
	sql += condition
	if month == time.Now().Format("2006-01") {
		//当月
		endDate = time.Now().AddDate(0, 0, -1).Format(utils.FormatDate)
	} else {
		endTime, _ := time.Parse("2006-01-02", endDate)
		endDate = endTime.AddDate(0, 0, -1).Format(utils.FormatDate)
	}
	sql += `  AND createtime = ? `
	sql += ` group by a.team`
	sql += ` order by return_capital_month desc`
	sql += ` ) AS t,(SELECT (@a :=0)) b`
	_, err = o.Raw(sql, paras, endDate).QueryRows(&list)
	return
}

func GetInnerGroupDataSumMonth(month, endDate, condition string, paras ...interface{}) (data []GroupData, err error) {
	o := orm.NewOrm()
	o.Using("dm_xjfq")
	sql := ` SELECT
	           SUM(money_on_file_month) AS money_on_file_month,
	           SUM(return_capital_month) AS return_capital_month,
	           SUM(capital_on_file_month) AS capital_on_file_month,
	           SUM(late_fee_on_file_month) AS late_fee_on_file_month,
	           SUM(rate_return_capital_month) AS rate_return_capital_month,
	           SUM(return_late_fee_month) AS return_late_fee_month,
	           SUM(rate_return_late_fee_month) AS rate_return_late_fee_month,
	           SUM(amount_on_file_month) AS amount_on_file_month
			FROM
				 dhht_internal_recycle a
			WHERE 1=1 `
	sql += condition
	if month == time.Now().Format("2006-01") {
		//当月
		endDate = time.Now().AddDate(0, 0, -1).Format(utils.FormatDate)
	} else {
		endTime, _ := time.Parse("2006-01-02", endDate)
		endDate = endTime.AddDate(0, 0, -1).Format(utils.FormatDate)
	}
	sql += `  AND createtime = ? `
	_, err = o.Raw(sql, paras, endDate).QueryRows(&data)
	return
}

// 查询委外组数据
func GetOutGroupData(timeType, condition string, paras ...interface{}) (list []GroupData, err error) {
	o := orm.NewOrm()
	o.Using("dm_xjfq")
	sql := ` SELECT t.*,(@a:=@a+1) rank  FROM (`
	sql += `SELECT a.team,
	           SUM(money_on_file) AS money_on_file,
	           SUM(money_on_file_month) AS money_on_file_month,
	           SUM(return_capital) AS return_capital,
	           SUM(return_capital_month) AS return_capital_month,
	           SUM(capital_on_file) AS capital_on_file,
	           SUM(capital_on_file_month) AS capital_on_file_month,
	           SUM(late_fee_on_file) AS late_fee_on_file,
	           SUM(late_fee_on_file_month) AS late_fee_on_file_month,
	           SUM(rate_return_capital) AS rate_return_capital,
	           SUM(rate_return_capital_month) AS rate_return_capital_month,
	           SUM(return_late_fee) AS return_late_fee,
	           SUM(return_late_fee_month) AS return_late_fee_month,
	           SUM(rate_return_late_fee) AS rate_return_late_fee,
	           SUM(rate_return_late_fee_month) AS rate_return_late_fee_month,
	           SUM(amount_on_file) AS amount_on_file,
	           SUM(amount_on_file_month) AS amount_on_file_month
             FROM dhht_out_recycle a
			WHERE 1=1 `
	sql += condition
	sql += ` group by a.team `
	if timeType == "day" {
		sql += ` order by return_capital desc`
	} else {
		sql += ` order by return_capital_month desc`
	}
	sql += ` ) AS t,(SELECT (@a :=0)) b`
	_, err = o.Raw(sql, paras).QueryRows(&list)
	return
}

// 查询委外组数据
func GetOutGroupDataMonth(month, endDate, condition string, paras ...interface{}) (list []GroupData, err error) {
	o := orm.NewOrm()
	o.Using("dm_xjfq")
	sql := ` SELECT t.*,(@a:=@a+1) rank  FROM (`
	sql += `SELECT a.team,
	           SUM(money_on_file_month) AS money_on_file_month,
	           SUM(return_capital_month) AS return_capital_month,
	           SUM(capital_on_file_month) AS capital_on_file_month,
	           SUM(late_fee_on_file_month) AS late_fee_on_file_month,
	           SUM(rate_return_capital_month) AS rate_return_capital_month,
	           SUM(return_late_fee_month) AS return_late_fee_month,
	           SUM(rate_return_late_fee_month) AS rate_return_late_fee_month,
	           SUM(amount_on_file_month) AS amount_on_file_month
             FROM dhht_out_recycle a
			WHERE 1=1 `
	sql += condition
	if month == time.Now().Format("2006-01") {
		//当月
		endDate = time.Now().AddDate(0, 0, -1).Format(utils.FormatDate)
	} else {
		endTime, _ := time.Parse("2006-01-02", endDate)
		endDate = endTime.AddDate(0, 0, -1).Format(utils.FormatDate)
	}
	sql += `  AND createtime = ? `
	sql += ` group by a.team `
	sql += ` order by return_capital_month desc`
	sql += ` ) AS t,(SELECT (@a :=0)) b`
	_, err = o.Raw(sql, paras, endDate).QueryRows(&list)
	return
}

func GetOutGroupDataSum(condition string, paras ...interface{}) (data GroupData, err error) {
	sql := `SELECT SUM(money_on_file) AS money_on_file,
	               SUM(money_on_file_month) AS money_on_file_month,
	               SUM(return_capital) AS return_capital,
	               SUM(return_capital_month) AS return_capital_month,
	               SUM(rate_return_capital) AS rate_return_capital,
	               SUM(rate_return_capital_month) AS rate_return_capital_month,
	               SUM(return_late_fee) AS return_late_fee,
	               SUM(return_late_fee_month) AS return_late_fee_month,
	               SUM(rate_return_late_fee) AS rate_return_late_fee,
	               SUM(rate_return_late_fee_month) AS rate_return_late_fee_month,
	               SUM(amount_on_file) AS amount_on_file,
	               SUM(amount_on_file_month) AS amount_on_file_month,
	               SUM(capital_on_file) AS capital_on_file,
	               SUM(late_fee_on_file) AS late_fee_on_file,
	               SUM(capital_on_file_month) AS capital_on_file_month,
	               SUM(late_fee_on_file_month) AS late_fee_on_file_month
			FROM dhht_out_recycle  WHERE 1=1 `
	sql += condition
	//sql += ` ORDER BY createtime DESC`
	o := orm.NewOrm()
	o.Using("dm_xjfq")
	err = o.Raw(sql, paras).QueryRow(&data)
	return
}

func GetOutGroupDataSumMonth(month, endDate, condition string, paras ...interface{}) (data []GroupData, err error) {
	sql := `SELECT
	               SUM(money_on_file_month) AS money_on_file_month,
	               SUM(return_capital_month) AS return_capital_month,
	               SUM(rate_return_capital_month) AS rate_return_capital_month,
	               SUM(return_late_fee_month) AS return_late_fee_month,
	               SUM(rate_return_late_fee_month) AS rate_return_late_fee_month,
	               SUM(amount_on_file_month) AS amount_on_file_month,
	               SUM(capital_on_file_month) AS capital_on_file_month,
	               SUM(late_fee_on_file_month) AS late_fee_on_file_month
			FROM dhht_out_recycle  WHERE 1=1 `
	sql += condition
	if month == time.Now().Format("2006-01") {
		//当月
		endDate = time.Now().AddDate(0, 0, -1).Format(utils.FormatDate)
	} else {
		endTime, _ := time.Parse("2006-01-02", endDate)
		endDate = endTime.AddDate(0, 0, -1).Format(utils.FormatDate)
	}
	sql += `  AND createtime = ? `
	o := orm.NewOrm()
	o.Using("dm_xjfq")
	_, err = o.Raw(sql, paras, endDate).QueryRows(&data)
	return
}

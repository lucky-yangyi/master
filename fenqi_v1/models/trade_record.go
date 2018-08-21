package models

import "zcm_tools/orm"

//发放记录
type TradeRecord struct {
	Id             int     `description:"主键"`
	OrderNumber    string  `description:"发放订单号"`
	CapitalAmount  float64 `description:"发放金额"`
	UserCardNumber string  `description:"放款卡号"`
	State          string  `description:"放款状态"`
	CreateTime     string  `description:"放款时间"`
	DtOrder        string  `description:"第三方订单时间"`
}

//查询发放记录
func QueryTradeRecord(loanId int) (tr TradeRecord, err error) {
	o := orm.NewOrm()
	sql := `SELECT
				id,
				order_number,
				capital_amount,
				user_card_number,
				create_time,
				state,
				dt_order
			FROM trade_record
			WHERE loan_id = ?`
	err = o.Raw(sql, loanId).QueryRow(&tr)
	return
}

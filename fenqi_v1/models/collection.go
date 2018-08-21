package models

import (
	"fenqi_v1/utils"
	"github.com/astaxie/beego"
	"strconv"
	"strings"
	"time"
	"zcm_tools/orm"
)

//催收处理
type CollectionHandle struct {
	Id            int       `orm:"column(int);null"`            //催收处理主键
	LoanId        int       `orm:"column(loan_id);null"`        //借款ID
	HandleUserId  int       `orm:"column(handle_user_id);null"` //催收处理主键
	ActionType    string    `orm:"column(action_type);null"`    //行动分类
	HandleTime    time.Time `orm:"column(handle_time);null"`    //处理时间
	Remark        string    `orm:"column(remark);null"`         //备注
	PromiseMoney  float64   `orm:"column(promise_money);null"`  //承诺金额，当行动分类为PTP有值
	Type          int       `orm:"column(type);null"`           //状态:1催收M1,2催收M2，3催收M3
	CompositeDate string    `orm:"column(composite_date);null"` //复合日期
	ConnRecordId  int       `orm:"column(conn_record_id);null"` //联系历史ID
}

//催收管理列表
type NewCollectMList struct {
	Id                             int       //还款计划Id
	LoanId                         int       `orm:"column(loan_id);null"`                             //借款ID
	Money                          float64   `orm:"column(money);null"`                               //借款金额
	Uid                            int       `orm:"column(uid);null"`                                 //借款人ID
	Account                        string    `orm:"column(account);null"`                             //手机号码
	VerifyRealName                 string    `orm:"column(verifyrealname);null"`                      //姓名
	CollectionName                 string    `orm:"column(collection_name);null"`                     //催收员用户名
	Displayname                    string    `orm:"column(displayname);null"`                         //催收员姓名
	LoanReturnDate                 time.Time `orm:"column(loan_return_date);null"`                    //应还日期
	CapitalAmount                  float64   `orm:"column(capital_amount);null"`                      //应收本金
	TaxAmount                      float64   `orm:"column(tax_amount);null"`                          //应收利息
	OverdueMoneyAmount             float64   `orm:"column(overdue_money_amount);null"`                //应收延期违约金
	OverdueBreachOfAmount          float64   `orm:"column(overdue_breach_of_amount);null"`            //应收滞纳金
	AheadOfTimeClearedAmount       float64   `orm:"column(ahead_of_time_cleared_amount);null"`        //应收提前结清违约金
	State                          string    `orm:"column(state);null"`                               //状态，未还,结清
	ReturnDate                     time.Time `orm:"column(return_date);null"`                         //实还日期
	ReturnTaxAmount                float64   `orm:"column(return_tax_amount);null"`                   //实还利息
	ReturnCapitalAmount            float64   `orm:"column(return_capital_amount);null"`               //实还本金
	OverdueDays                    int       `orm:"column(overdue_days);null"`                        //逾期天数
	ReturnOverdueMoneyAmount       float64   `orm:"column(return_overdue_money_amount);null"`         //实还滞纳金
	ReturnOverdueBreachOfAmount    float64   `orm:"column(return_overdue_breach_of_amount);null"`     //实还延期违约金
	ReturnAheadOfTimeClearedAmount float64   `orm:"column(return_ahead_of_time_cleared_amount);null"` //实还提前结清违约金
	RemainMoneyChargeUpAmount      float64   `orm:"column(remain_money_charge_up_amount);null"`       //本次挂账余额
	ReturnDataServiceFee           float64
	LoanTermCount                  int       `orm:"column(loan_term_count);null"` //借款期限
	Actiontype                     string    `orm:"column(action_type);null"`     //行动分类
	CompositeDate                  string    `orm:"column(composite_date);null"`  //复合时间
	LoanDate                       time.Time `orm:"column(loan_date);null"`       //借款时间
	Ctype                          int
	HandleId                       int
	HandleTime                     string
	OrgId                          int
	DerateAmount                   float64 // 费用减免
	ProductId                      int
	ReturnMoney                    float64
	SumReturnMoney                 float64
	Mtype                          int //所处的阶段
	Name                           string
	FinishDate                     string
	Sumloanmoney                   float64
	Finishmoney                    float64
	CostFee                        float64
	DataServiceFee                 float64
	MaxOverdueDays                 int
	ActionType                     string
	TermNo                         int
	NeedPayment                    float64 //到期待还或逾期待还
	ExpirationTime                 string  //到期时间
	IdCard                         string
	ContactPhone                   string
	ContractCode                   string
	Flag                           bool
}

//催收管理列表
func GetList(where, orgStr string, begin, count, mtype int, pars ...string) (mlist []*NewCollectMList, err error) {
	sql := `SELECT
			  a.id AS loan_id,
			  a.money,
			  a.account,
			  a.verifyrealname,
			  a.handle_time,
			  a.collection_name,
			  a.action_type,
			  a.composite_date,
			  a.expiration_time,
			  a.max_overdue_days
			FROM loan AS a,repayment_schedule b
			WHERE a.id = b.loan_id AND a.state='BACKING' AND b.state='BACKING' AND a.mtype=? AND b.loan_return_date<=curdate()  `
	if where != "" {
		sql += where
	}
	if orgStr != "" {
		sql += ` AND (a.collection_user_id=? OR a.org_id IN(` + orgStr + `))`
	} else {
		sql += ` AND a.collection_user_id=?`
	}
	sql += "   ORDER BY  a.collection_handle_id  ASC   LIMIT ?,?"
	o := orm.NewOrm()
	o.Using("read")
	_, err = o.Raw(sql, mtype, pars, begin, count).QueryRows(&mlist)
	return
}

//获取M管理总记录数
func GetCount(where, orgStr string, mtype int, pars ...string) []int {
	sql := `SELECT COUNT(1)
			FROM loan AS a
			INNER JOIN repayment_schedule b ON a.id =b.loan_id
			WHERE 1=1 AND a.state='BACKING' AND a.mtype=? AND b.state='BACKING' AND b.loan_return_date<=?  `
	if where != "" {
		sql += where
	}
	if orgStr != "" {
		sql += ` AND (a.collection_user_id=? OR a.org_id IN(` + orgStr + `))`
	} else {
		sql += ` AND a.collection_user_id=? `
	}
	sql += ` GROUP BY a.id`
	o := orm.NewOrm()
	o.Using("read")
	var count []int
	today := time.Now().Format(utils.FormatDate)
	o.Raw(sql, mtype, today, pars).QueryRows(&count)
	return count
}

//催收人员
type CollectionUser struct {
	Id             int    `orm:"column(id);null"`             //用户ID
	Account        string `orm:"column(account);null"`        //用户手机号码
	Verifyrealname string `orm:"column(verifyrealname);null"` //用户真实姓名
	LimitCount     int    `orm:"column(limit_count);null"`    //限制条数
	AccountType    int
}

//获取加权限的催收员
func GetCollectUser(orgstr string, mtype int) []CollectionUser {
	sql := `SELECT DISTINCT u.id,u.displayname AS verifyrealname,u.account_type
            FROM sys_user AS u
            INNER JOIN sys_station_type AS t ON u.station_id=t.station_id
            INNER JOIN sys_station AS ss ON t.station_id=ss.id
            WHERE t.type =?  AND ss.org_id in(` + orgstr + `) `
	sql += " and u.accountstatus='启用' "
	users := []CollectionUser{}
	o := orm.NewOrm()
	o.Using("read")
	o.Raw(sql, mtype).QueryRows(&users)
	return users
}

//回收案件根据催收员获取其组织架构id和type
func GetInfoByCollector(cid int) (mtype int) {
	o := orm.NewOrm()
	o.Using("read")
	sql := `	SELECT sst.type
            FROM sys_user su
            INNER JOIN sys_station ss ON su.station_id = ss.id
            INNER JOIN sys_organization so ON ss.org_id=so.id
            INNER JOIN sys_station_type sst ON  sst.station_id=ss.id
            WHERE su.id =? `
	o.Raw(sql, cid).QueryRow(&mtype)
	return mtype
}

//费用减免人员配置
type CostConfig struct {
	Id                          int
	CollectionDistributeRoleId  string
	StationName                 string
	CollectionExeclLetterRoleId string
	OutsourceBlackSignRoleId    string
	CollectionExeclPhoneRoleId  string
	CollectionSignRoleId        string
	RequestRoleId               string
	ApproveId                   string
	DisposeId                   string
	StationId                   string
	CostExeclRoleId             string
}

//查询所有费用减免人员配置信息
func QueryByIdCostConfig(id int) (config CostConfig, err error) {
	o := orm.NewOrm()
	sql := `select *  from cost_config  where id=?`
	err = o.Raw(sql, id).QueryRow(&config)
	return
}

//  回收案件
func GetRecycleCase(begin, size int, condition, orgStr string, pars ...string) (list []NewCollectMList, err error) {
	sql := `SELECT
	             a.id as loan_id,
	             a.uid,
	             a.account,
	             a.verifyrealname as name,
	             b.capital_amount,
	             b.term_no,
	             b.return_date,
	             b.overdue_days,
	             b.mtype,
                 b.tax_amount,
                 b.overdue_money_amount,
                 b.overdue_breach_of_amount,
                 b.data_service_fee,
                 b.derate_amount,
                 b.collection_name
            FROM
	             loan AS a
            INNER JOIN repayment_schedule AS b ON a.id = b.loan_id
	        WHERE  b.state = 'FINISH'    AND b.mtype >=11    `
	if condition != "" {
		sql += condition
	}
	if orgStr != "" {
		sql += ` AND (b.collection_user_id=? OR b.org_id IN(` + orgStr + `))`
	} else {
		sql += ` AND b.collection_user_id=?`
	}
	sql += `    LIMIT ?, ?  `
	o := orm.NewOrm()
	o.Using("read")
	_, err = o.Raw(sql, pars, begin, size).QueryRows(&list)
	if err != nil && err != orm.ErrNoRows {
		return nil, err
	}
	return list, nil
}

//得到所有催收员催收员
func GetAllCollectUser(orgstr string) []CollectionUser {
	sql := `SELECT DISTINCT u.id,u.displayname AS verifyrealname,u.account_type
            FROM sys_user AS u
            INNER JOIN sys_station_type AS t ON u.station_id=t.station_id
            INNER JOIN sys_station AS ss ON t.station_id=ss.id
            WHERE t.type >=11   AND ss.org_id in(` + orgstr + `) `
	sql += " and u.accountstatus='启用' "
	users := []CollectionUser{}
	o := orm.NewOrm()
	o.Using("read")
	o.Raw(sql).QueryRows(&users)
	return users
}

//回收案件的数量
func GetRecycleCount(condition, orgStr string, pars ...string) int {
	sql := `SELECT count(1) FROM loan AS a
	        INNER JOIN repayment_schedule AS b ON a.id = b.loan_id
	        WHERE  b.state = 'FINISH'  AND b.mtype >=11   `
	if condition != "" {
		sql += condition
	}
	if orgStr != "" {
		sql += ` AND (b.collection_user_id=? OR b.org_id IN(` + orgStr + `))`
	} else {
		sql += ` AND b.collection_user_id=?`
	}
	var count int
	o := orm.NewOrm()
	o.Using("read")
	o.Raw(sql, pars).QueryRow(&count)
	return count
}

//借款
type Loan struct {
	Id                     int
	Uid                    int
	Account                string
	ContractCode           string
	Money                  float64
	ProductId              int
	LoanTaxFee             float64
	LoanServiceFee         float64
	LoanTax                float64
	LoanService            float64
	LoanTermCount          int
	UsersCouponId          int
	CouponValue            int
	LoanType               int
	UserBankId             int
	LoanOverdueFee         float64
	LoanOverdueBreachOfFee float64
	State                  string
	LoanDate               time.Time
	CreateTime             time.Time
	RepaymentScheduleId    int
	Remark                 string
	CurrState              int
	ReplyTime              time.Time
	Channel                int
	RealMoney              float64
	CardNumber             string
	Overdue_days           int
	EndDate                time.Time
	FinishDate             time.Time
	Lotion                 string
	OidPaybill             string
	DisplayName            string `orm:"column(displayname);null"` //系统用户真实姓名
	Loan_agreement_url     string
	Service_agreement_url  string
	UrgentFee              float64
	PayMoney               float64
	IsRenew                int
	Name                   string
	CollectionUserId       int
	PkgType                int //平台来源  0花无忧，1有个钱包
}

//获取借款详情
func GetLoanById(id int) (l *Loan, err error) {
	o := orm.NewOrm()
	sql := `SELECT *
			from loan as a
			where  a.id= ? `
	err = o.Raw(sql, id).QueryRow(&l)
	return
}

//判断目前为止案件未还期数
func GetCaseStateCount(id int) (count int, err error) {
	o := orm.NewOrm()
	sql := `SELECT COUNT(1)
			FROM repayment_schedule rs
			WHERE rs.loan_id=? AND rs.loan_return_date<=? AND rs.state ='BACKING' `
	today := time.Now().Format(utils.FormatDate)
	err = o.Raw(sql, id, today).QueryRow(&count)
	return
}

//根据借款ID，获取催收信息
func GetHandleListById(loanId int) (collection []*NewCollectMList, err error) {
	sql := `SELECT
			a.id as loan_id,a.money,a.uid,b.overdue_money_amount,a.loan_date,b.overdue_days,b.overdue_breach_of_amount,b.capital_amount,
            b.tax_amount,b.remain_money_charge_up_amount,b.term_no,b.return_overdue_money_amount,b.return_capital_amount,b.return_tax_amount,
            b.return_data_service_fee,b.data_service_fee,b.state,b.loan_return_date,b.return_overdue_breach_of_amount
			FROM loan a
			INNER JOIN repayment_schedule AS b
			ON a.id=b.loan_id
			WHERE 1=1 AND a.state='BACKING'   AND a.id=?  `
	o := orm.NewOrm()
	o.Using("read")
	_, err = o.Raw(sql, loanId).QueryRows(&collection)
	return
}

// 查看-个人信息
func GetUsersMetadata(uid int) (um *UsersMetadata, err error) {
	sql := ` SELECT um.uid,um.verifyrealname,um.id_card,um.sex,um.account,um.is_verifyemail,um.verify_time,
	        um.mail_address,um.birth_date,um.balance,um.use_balance,um.isnt_sys_quota,um.assess_time,um.finger_key,um.init_audit,
            um.zm_score,
 	        u.create_time, u.pkg_type,u.source,yu.address,u.tag_type,u.tag_sys_id,u.is_tag,u.remark,
 	        ub.live_address , ub.live_detail_address,u.mobile_type_recent,ub.live_time,ub.marriage
			FROM users_metadata um
			INNER JOIN users u ON um.uid=u.id
			LEFT join users_base_info ub on
			um.uid = ub.uid
			LEFT join ocr_info yu on
			u.id=yu.uid
			WHERE um.uid=? `
	err = orm.NewOrm().Raw(sql, uid).QueryRow(&um)
	sql = `SELECT
				create_time AS operationtime
			FROM
				login_record
			WHERE
				uid = ?
			ORDER BY
				id DESC
			LIMIT 1 `
	orm.NewOrm().Raw(sql, uid).QueryRow(&um.OperationTime)
	return
}

//银行卡
type CardMini struct {
	Id          int
	Uid         int
	Card_number string
	Bank_name   string
	Bank_mobile string
	Bind_time   time.Time
	State       string //绑卡状态
}

// 个人信息 银行卡
func BankcardInfo(uid int) (card *CardMini, err error) {
	sql := `SELECT id, uid, card_number, bank_name, bank_mobile, create_time bind_time,state
			FROM users_bankcards
			WHERE uid=? AND state = 'USING' limit 1`
	err = orm.NewOrm().Raw(sql, uid).QueryRow(&card)
	if err != nil {
		return nil, err
	}
	return
}

//行动分类为‘盗办’，‘死亡’，’坐牢‘，的用户自动入黑
func UpdateBlackByUserId(uid, id int, actionType string) (err error) {
	o := orm.NewOrm()
	sql := `UPDATE users SET tag_type=?,is_tag=? ,tag_sys_id=?,black_time =now()  where id=? `
	var isTag int
	if actionType == "盗办" {
		isTag = 4
	} else if actionType == "死亡" {
		isTag = 5
	} else if actionType == "坐牢" {
		isTag = 6
	}
	_, err = o.Raw(sql, actionType, isTag, id, uid).Exec()
	return
}

// 联系历史
type Conn_record struct {
	Id             int
	Uid            int
	Conn_type      string
	Content        string
	Created_by     int
	Create_time    time.Time
	Modify_by      int
	Modify_by_name string `orm:"column(modify_by_name);null"` //操作人
	Modify_time    time.Time
	Action_type    string
	DealStatus     int
	DealResult     string
}

//新增联系历史
func (record *Conn_record) Insert() error {
	sql := `INSERT INTO conn_record(uid, conn_type, content, created_by, create_time,deal_status)
			VALUES(?, ?, ?, ?, now(), ?)`
	rec, err := orm.NewOrm().Raw(sql, record.Uid, record.Conn_type, record.Content, record.Created_by, record.DealStatus).Exec()
	beego.Info(err)
	lastId, err := rec.LastInsertId()
	record.Id = int(lastId)
	return err
}

//催收处理
func (handle *CollectionHandle) HandleInsert() error {
	o := orm.NewOrm()
	sql := `INSERT INTO collection_handle(loan_id,handle_user_id,action_type,handle_time,remark,promise_money,type,composite_date,conn_record_id)
	VALUES(?,?,?,?,?,?,?,?,?)`
	res, err := o.Raw(sql, handle.LoanId, handle.HandleUserId, handle.ActionType,
		time.Now(), handle.Remark, handle.PromiseMoney, handle.Type, handle.CompositeDate, handle.ConnRecordId).Exec()
	if err == nil {
		id, _ := res.LastInsertId()
		o.Raw("update loan set collection_handle_id=? where id=?", id, handle.LoanId).Exec()

	}
	return err
}

//更新loan表的行动分类
func (handle *CollectionHandle) UpdateLoan() error {
	sql := `UPDATE loan SET handle_time=?,action_type=?,composite_date=?,handle_user_id=? WHERE id=?`
	_, err := orm.NewOrm().Raw(sql, time.Now(), handle.ActionType, handle.CompositeDate, handle.HandleUserId, handle.LoanId).Exec()
	return err
}

//根据用户ID获取组织架构信息
func GetOrganizationByUserId(userId int) (organization *SysOrganization, err error) {
	sql := `SELECT o.id,o.parent_id,o.name,o.remark FROM sys_user AS u
			INNER JOIN sys_station AS s
			ON u.station_id=s.id
			INNER JOIN sys_organization AS o
			ON s.org_id=o.id
			WHERE u.id=?`
	err = orm.NewOrm().Raw(sql, userId).QueryRow(&organization)
	return
}

//手动分配订单给相应的催收人员
func InManualDistributionModels(loanIds []string, ctype, operation_id, collectionUserId, orgId int) error {
	o := orm.NewOrm()
	//修改借款表中，催收人员信息
	su, err := GetColletionNameById(collectionUserId)
	if err != nil {
		return err
	}
	sql := `UPDATE loan SET collection_user_id=?,org_id=?,collection_name=?  WHERE id=?`
	//添加手动分发催收记录
	sqlAdd := `INSERT INTO installment_log(loan_id,uid,type,createTime,org_id,money,operation_id,rid,term_no,receive_money)VALUES(?,?,?,?,?,?,?,?,?,?)`
	up, err := o.Raw(sql).Prepare()
	if err != nil {
		return err
	}
	add, err := o.Raw(sqlAdd).Prepare()
	if err != nil {
		return err
	}
	for _, v := range loanIds {
		// 根据loan_id 获取对接的该阶段的分期案件
		rsGroup, err := GetFqInstanceByLoanId(v, ctype)
		if err != nil {
			return err
		}
		var receive_money float64
		sql := `SELECT
				SUM(rs.capital_amount) - SUM(rs.return_capital_amount) AS receive_money
			FROM
				repayment_schedule rs
			WHERE
				rs.state = "BACKING" AND rs.loan_id=?`
		orm.NewOrm().Raw(sql, v).QueryRow(&receive_money)
		for _, rs := range rsGroup {
			_, err = up.Exec(collectionUserId, orgId, su.Displayname, v)
			_, err = add.Exec(v, collectionUserId, ctype, time.Now(), orgId, rs.Shouldmoney, operation_id, rs.Id, rs.Term_no, receive_money)
		}
	}
	up.Close()
	add.Close()
	return err
}

//系统用户
type SysUsers struct {
	Displayname string
	AccountType int
}

//根据uid得到系统用户的名字和账号类型
func GetColletionNameById(id int) (su *SysUsers, err error) {
	sql := `SELECT displayname,account_type FROM sys_user WHERE id=?`
	o := orm.NewOrm()
	err = o.Raw(sql, id).QueryRow(&su)
	return
}

//根据loan_id 获取对接的该阶段的分期案件
func GetFqInstanceByLoanId(loanid string, mtype int) (rs []*RepaymentSchedule, err error) {
	o := orm.NewOrm()
	o.Using("read")
	sql := `SELECT  b.id,b.term_no, a.id AS loan_id, IFNULL(b.capital_amount, 0) + IFNULL(b.tax_amount, 0) + IFNULL(b.overdue_breach_of_amount, 0) + IFNULL(b.overdue_money_amount, 0) - IFNULL(b.derate_amount, 0) AS shouldmoney
			FROM loan a
			INNER JOIN repayment_schedule b ON a.id = b.loan_id
			WHERE a.state = 'backing' AND b.state = 'BACKING' AND b.loan_return_date <= ? AND a.mtype = ? AND a.id = ?`
	date := time.Now().Format(utils.FormatDate)
	_, err = o.Raw(sql, date, mtype, loanid).QueryRows(&rs)
	return
}

//得到到期待还或逾期待还金额
func GetNeedPayment(loanId int) (money float64) {
	sql := ` SELECT SUM(IFNULL(rs.capital_amount,0)
			        + IFNULL(rs.tax_amount,0)
					+ IFNULL(rs.overdue_money_amount,0)
					+ IFNULL(rs.overdue_breach_of_amount,0)
					+ IFNULL(rs.data_service_fee,0)
					- IFNULL(rs.remain_money_charge_up_amount,0)) AS money
			FROM repayment_schedule rs
			WHERE rs.loan_id=? AND rs.loan_return_date<=? AND rs.state='BACKING' `
	today := time.Now().Format(utils.FormatDate)
	orm.NewOrm().Raw(sql, loanId, today).QueryRow(&money)
	return
}

//联系历史数量
func ConnRcdsCount(uid int, condition string, pars []string) int {
	sql := `SELECT count(1)
			FROM conn_record cr
			LEFT JOIN sys_user su ON cr.modify_by=su.id
			LEFT JOIN collection_handle ch ON cr.id=ch.conn_record_id
			WHERE cr.uid=? AND cr.conn_type="COLLECTION"`
	sql += condition
	var count int
	orm.NewOrm().Raw(sql, uid, pars).QueryRow(&count)
	return count
}

//联系历史展示
func ConnRcdsList(uid int, condition string, pars []string, begin, count int) (list []Conn_record, err error) {
	sql := `SELECT cr.deal_status,cr.deal_result,cr.id, cr.conn_type, cr.content, su.displayname as modify_by_name,cr.modify_by,cr.uid, cr.modify_time,cr.create_time,ch.action_type
			FROM conn_record cr
			LEFT JOIN sys_user su ON cr.created_by=su.id
			LEFT JOIN collection_handle ch ON cr.id=ch.conn_record_id
			AND ch.conn_record_id>0
			WHERE cr.uid=? AND cr.conn_type="COLLECTION"`
	sql += condition
	sql += " ORDER BY cr.create_time desc limit ?, ?"
	_, err = orm.NewOrm().Raw(sql, uid, pars, begin, count).QueryRows(&list)
	if err == nil {
		return list, nil
	} else if err == orm.ErrNoRows {
		return nil, nil
	}
	return nil, err
}

//投诉处理
func ComplaintHandlingList(uid, start, pageSize int) (list []Conn_record, err error) {
	sql := `SELECT
				cr.id,
				cr.uid,
				cr.content,
				su.displayname AS modify_by_name,
				cr.create_time
			FROM
				conn_record cr
			LEFT JOIN sys_user su ON cr.created_by = su.id
			WHERE
				cr.uid =? AND cr.conn_type="COMPLANIT"
			ORDER BY
				cr.create_time DESC
			LIMIT ?,?`
	_, err = orm.NewOrm().Raw(sql, uid, start, pageSize).QueryRows(&list)
	return
}

//投诉处理
func ComplaintHandlingListCount(uid int) (count int, err error) {
	sql := `SELECT
				COUNT(1)
			FROM
				conn_record cr
			LEFT JOIN sys_user su ON cr.created_by = su.id
			WHERE
				cr.uid =? AND cr.conn_type="COMPLANIT"
			ORDER BY
				cr.create_time DESC`
	err = orm.NewOrm().Raw(sql, uid).QueryRow(&count)
	return
}

//根据借款ID，获取催收信息
func GetINstallHandleListById(loanId int) (collection []*NewCollectMList, err error) {
	sql := `SELECT
			a.id AS loan_id,
			a.money,
			a.uid,
			a.loan_term_count,
			a.loan_date,
			a.contract_code,
			a.account,
			a.verifyrealname AS name,
			a.collection_name,
			a.action_type,
	        a.composite_date,
			b.loan_return_date,
			b.capital_amount,
			b.tax_amount,
			b.overdue_money_amount,
	        b.overdue_breach_of_amount,
			b.ahead_of_time_cleared_amount,
			b.state,
	        b.return_date,
			b.return_tax_amount,
			b.return_capital_amount,
			b.overdue_days,
	        b.return_overdue_money_amount,
			b.return_overdue_breach_of_amount,
			b.return_ahead_of_time_cleared_amount,
	        b.remain_money_charge_up_amount,
			b.data_service_fee,
			c.id_card
		FROM loan AS a
		INNER JOIN repayment_schedule AS b ON a.id=b.loan_id
		INNER JOIN users_metadata AS c ON a.uid=c.uid
		WHERE 1=1 AND a.state='BACKING' AND b.state='BACKING' AND a.id=? AND b.loan_return_date <=? `
	o := orm.NewOrm()
	o.Using("read")
	date := time.Now().Format(utils.FormatDate)
	_, err = o.Raw(sql, loanId, date).QueryRows(&collection)
	return
}

//用户联系人信息
type UsersTelephoneDirectory struct {
	Id                 int64     `description:"主键"`
	Uid                int64     `description:"用户uid"`
	ContactName        string    `description:"联系人姓名"`
	ContactPhoneNumber string    `description:"联系人手机号"`
	IsUrgent           int64     `description:"是否紧急联系人0否,1是"`
	Relation           string    `description:"关系"`
	Code               string    `description:"code编码"`
	CreateTime         time.Time `description:"创建时间"`
	SignState          int64     //标记状态
	SignRemark         string    //标记备注
	SignTime           time.Time //标记时间
}

// 查看联系人count
func QueryTelephonCount(uid int, phoneNum string) (count int, err error) {
	session := utils.GetSession()
	defer session.Close()
	smap := map[string]interface{}{}
	c := session.DB(utils.MGO_DB).C("mailList")
	if uid > 0 {
		smap["uid"] = uid
	}
	if phoneNum != "" {
		smap["contactphonenumber"] = phoneNum
	}
	count, err = c.Find(&smap).Count()
	return
}

// 查看联系人
func QueryTelephon(uid int, phoneNum string) (list []MailList, err error) {
	session := utils.GetSession()
	defer session.Close()
	smap := map[string]interface{}{}
	c := session.DB(utils.MGO_DB).C("mailList")
	if uid > 0 {
		smap["uid"] = uid
	}
	if phoneNum != "" {
		smap["contactphonenumber"] = phoneNum
	}
	err = c.Find(&smap).All(&list)
	return
}

// 查看联系人
func QueryTelephonInfo(uid int) (list []MailList, err error) {
	session := utils.GetSession()
	defer session.Close()
	smap := map[string]interface{}{}
	c := session.DB(utils.MGO_DB).C("mailList")
	if uid > 0 {
		smap["uid"] = uid
	}

	err = c.Find(&smap).All(&list)
	return
}

// 查看联系人count
func QueryTelephonInfoCount(uid int, phoneNum string) (count int, err error) {
	session := utils.GetSession()
	defer session.Close()
	smap := map[string]interface{}{}
	c := session.DB(utils.MGO_DB).C("mailList")
	if uid > 0 {
		smap["uid"] = uid
	}

	count, err = c.Find(&smap).Count()
	return
}

// 查看联系人不分页
func QueryTelephonS(uid int, phoneNum string) (list []UsersTelephoneDirectory, err error) {
	session := utils.GetSession()
	defer session.Close()
	smap := map[string]interface{}{}
	c := session.DB(utils.MGO_DB).C("users_telephone_directory")
	if uid > 0 {
		smap["uid"] = uid
	}
	if phoneNum != "" {
		smap["contactphonenumber"] = phoneNum
	}
	err = c.Find(&smap).Sort("-signtime").All(&list)
	return
}

//给用户通讯录去重
func DistinctList(list []UsersTelephoneDirectory) (distinctList []UsersTelephoneDirectory) {
	var phones []string
	for _, v := range list {
		if !strings.Contains(strings.Join(phones, ","), v.ContactPhoneNumber) {
			phones = append(phones, v.ContactPhoneNumber)
			distinctList = append(distinctList, v)
		}
	}
	return
}

//手机通讯录
func QueryTelephonList(uid int) {

}

// 查看紧急联系人
func QueryInstancyLinkman(uid int) (linklist []UsersLinkman, err error) {
	o := orm.NewOrm()
	sql := "select id,uid,relation,contact_phone_number,is_normal,abnormal_result,source,result_code,sign_state,sign_remark,sign_time from users_linkman where uid=?"
	_, err = o.Raw(sql, uid).QueryRows(&linklist)
	return
}

/*//
func UpdateCountUsersModifyLinkmanById(id, signState int, makeContent string) (num int64, err error) {
	o := orm.NewOrm()
	sql := ` UPDATE users_modify_linkman SET sign_state=?,sign_remark = ? ,sign_time = NOW() WHERE id =? `
	res, err := o.Raw(sql, signState, makeContent, id).Exec()
	if err != nil {
		return 0, err
	}
	num, err = res.RowsAffected()
	return
}*/

//更新users_linkman
func SignUsersTelephoneDirectory(id, state int, Abnormal_result string) (num int64, err error) {
	o := orm.NewOrm()
	sql := `UPDATE users_linkman SET sign_state = ?, sign_remark = ?,sign_time = NOW() WHERE id = ?`
	res, err := o.Raw(sql, state, Abnormal_result, id).Exec()
	if err != nil {
		return 0, err
	}
	num, err = res.RowsAffected()
	return
}

//根据用户ID获取运营商授权类型1:魔蝎,2:天机
func GetMobileAuthTypeByUserId(id int) (mobileAuthType int, err error) {
	sql := `SELECT mobile_auth_type FROM users_auth WHERE uid = ? AND is_valid=1 `
	err = orm.NewOrm().Raw(sql, id).QueryRow(&mobileAuthType)
	return
}

func UpdateMobileState(userReportResult UserReportResult, mobils string, signState int, signRemark string) bool {
	isFind := false
	for key, v := range userReportResult.ReportResult.Data.Call_contact_detail {
		if v.Peer_num == mobils {
			userReportResult.ReportResult.Data.Call_contact_detail[key].SignState = signState
			userReportResult.ReportResult.Data.Call_contact_detail[key].SignRemark = signRemark
			userReportResult.ReportResult.Data.Call_contact_detail[key].SignTime = time.Now()
			isFind = true
			break
		}
	}
	return isFind
}

//修改常用人标记信息
func ModifyCommonContacts(id, state int, signRemark string) error {
	o := orm.NewOrm()
	sql := `UPDATE bqs_mno_onemonth_commonly_connect_mobiles SET sign_state = ?, sign_remark = ?,sign_time = NOW() WHERE id = ?`
	_, err := o.Raw(sql, state, signRemark, id).Exec()
	return err
}

// 查看联系人2
func QueryTelephon2(uid int, phoneNum string) (list *UsersTelephoneDirectory, err error) {
	session := utils.GetSession()
	defer session.Close()
	smap := map[string]interface{}{}
	c := session.DB(utils.MGO_DB).C("users_telephone_directory")
	if uid > 0 {
		smap["uid"] = uid
	}
	if phoneNum != "" {
		smap["contactphonenumber"] = phoneNum
	}
	err = c.Find(&smap).One(&list)
	return
}

//TOP15最近一个月通话
type BqsMnoOnemonthCommonlyConnectMobiles struct {
	Id                   int       `description:"编号"`
	Uid                  int64     `description:"用户编号"`
	Mobile               string    `description:"关联电话"`
	MonType              string    `description:"号码类型(电信,移动等)"`
	BeginTime            time.Time `description:"开始时间"`
	EndTime              time.Time `description:"结束时间"`
	BelongTo             string    `description:"归属地"`
	ConnectCount         string    `description:"拨打次数"`
	ConnectTime          string    `description:"拨打时间"`
	OriginatingCallCount string    `description:"打电话次数"`
	TerminatingCallCount string    `description:"接电话次数"`
	CreateTime           time.Time `description:"创建时间"`
	SignState            int64
	SignRemark           string
	SignTime             time.Time
}

type CommonlyConnectMobiles struct {
	Id                     int
	Mobile                 string `bson:"mobile"`
	Begin_time             string `bson:"begintime"`
	End_time               string `bson:"endtime"`
	Belong_to              string `bson:"belongto"`
	Connect_count          string `bson:"connectcount"`
	Connect_time           string `bson:"connecttime"`
	Originating_call_count string `bson:"originatingcallcount"`
	Terminating_call_count string `bson:"terminatingcallcount"`
	Mon_type               string `bson:"montype"`
	SignState              int64
	SignRemark             string
	SignTime               time.Time
	Uid                    int
	CreateTime             time.Time `bson:"_id"`
	Period_type            string
}

type CommonlyConnectMobiles_2 struct {
	Id                     int
	Mobile                 string
	Begin_time             string
	End_time               string
	Belong_to              string
	Connect_count          int
	Connect_time           string
	Originating_call_count string
	Terminating_call_count string
	Mon_type               string
	SignState              int64
	SignRemark             string
	SignTime               time.Time
	Uid                    int
	CreateTime             time.Time
	Period_type            string
}

type MonGoQuery struct {
	Uid int
}

//催收查询
func GetMSearch(where string, pars []string) (mlist []*NewCollectMList, err error) {
	sql := `	SELECT
				a.id AS loan_id,
				a.money,
				a.max_overdue_days,
				a.uid,
				a.account,
				a.verifyrealname,
				a.collection_name,
				a.action_type,
				a.handle_time,
				a.mtype,
				a.composite_date
			FROM
				loan AS a
			INNER JOIN repayment_schedule AS b ON a.id = b.loan_id
			WHERE
				1 = 1
			AND a.state = 'BACKING'
			AND b.state = 'BACKING' AND b.mtype > 0`
	if where != "" {
		sql += where
	}
	sql += `  group by  a.id `
	o := orm.NewOrm()
	o.Using("read")
	_, err = o.Raw(sql, pars).QueryRows(&mlist)
	return
}

//投诉处理
type Comlanit struct {
	Id           int
	Uid          int
	Content      string
	ModifyByName string
	ModifyTime   string
	CreateTime   string
	DealStatus   int
	DealResult   string
	Account      string
	ModifyBy     int
}

func GetComplanitInfo(startDate, endDate, mobile string, status, begin, count int) ([]Comlanit, error) {
	var com []Comlanit
	var params []string
	sql := `SELECT cr.id,cr.uid, cr.content,cr.modify_by,su.displayname AS modify_by_name,cr.modify_time, cr.create_time,cr.deal_status, cr.deal_result,u.account FROM  conn_record cr
        LEFT JOIN sys_user su ON cr.created_by = su.id
        LEFT JOIN collection_handle ch ON cr.id = ch.conn_record_id
        LEFT JOIN users u ON cr.uid = u.id WHERE  cr.deal_status >0 `
	if startDate != "" {
		sql += ` AND cr.create_time >= ? `
		params = append(params, startDate)
	}
	if endDate != "" {
		sql += ` AND cr.create_time <= ? `
		params = append(params, endDate)
	}
	if mobile != "" {
		sql += ` AND u.account = ?`
		params = append(params, mobile)
	}
	if status > 0 {
		sql += ` AND cr.deal_status = ? `
		params = append(params, strconv.Itoa(status))
	}
	sql += " order by cr.create_time desc limit ?, ?"
	o := orm.NewOrm()
	o.Using("read")
	_, err := o.Raw(sql, params, begin, count).QueryRows(&com)
	return com, err
}

func GetComplanitCount(startDate, endDate, mobile string, status int) (count int, err error) {
	var params []string
	sql := `SELECT COUNT(1) FROM  conn_record cr
        LEFT JOIN sys_user su ON cr.created_by = su.id
        LEFT JOIN collection_handle ch ON cr.id = ch.conn_record_id
        LEFT JOIN users u ON cr.uid = u.id WHERE  cr.deal_status > 0 `
	if startDate != "" {
		sql += ` AND cr.create_time >= ? `
		params = append(params, startDate)
	}
	if endDate != "" {
		sql += ` AND cr.create_time <= ? `
		params = append(params, endDate)
	}
	if mobile != "" {
		sql += ` AND u.account = ?`
		params = append(params, mobile)
	}
	if status > 0 {
		sql += ` AND cr.deal_status = ? `
		params = append(params, strconv.Itoa(status))
	}
	o := orm.NewOrm()
	o.Using("read")
	err = o.Raw(sql, params).QueryRow(&count)
	return
}

func GetUsersCount(id int) (count int, err error) {
	sql := `SELECT COUNT(1) FROM users_metadata WHERE uid=?`
	err = orm.NewOrm().Raw(sql, id).QueryRow(&count)
	return
}

func ConnRcDeal(id, flagId, sysId int, res string) error {
	_, err := orm.NewOrm().Raw(`UPDATE conn_record SET deal_result = ?,deal_status = ?,modify_by=?,modify_time=NOW() WHERE id = ?`, res, flagId, sysId, id).Exec()
	return err
}

//投诉操作处理记录
type ComplainOperationLog struct {
	Id              int    `description:"主键"`
	Uid             int    `description:"用户ID"`
	ConnId          int    `description:"投诉 id"`
	OperationId     int    `description:"操作人 id"`
	OperationPerson string `description:"操作人"`
	DealResult      string `description:"内容"`
	DealStatus      int    `description:"2  更加 3 处理困难 4 处理完毕"`
	OperateTime     string `description:"操作时间"`
}

//添加
func AddComplainOperationLog(uid, connId, operationId, dealStatus int, operationPerson, dealResult string) (err error) {
	o := orm.NewOrm()
	insertSql := `insert into complain_operation_log(uid,conn_id,operation_id,operation_person,deal_result,deal_status,operate_time)values(?,?,?,?,?,?,?)`
	_, err = o.Raw(insertSql, uid, connId, operationId, operationPerson, dealResult, dealStatus, time.Now()).Exec()
	return
}

//条件查询列表
func QueryComplainOperationLog(id int) (list []ComplainOperationLog, err error) {
	o := orm.NewOrm()
	sql := "select  *  from complain_operation_log where conn_id= ? order by id desc  "
	_, err = o.Raw(sql, id).QueryRows(&list)
	return
}

//从mobile_directory_sign查询记录是否存在
func GetSignCountByUidAndPhone(uid int, name, moblie, stype string) (count int) {
	o := orm.NewOrm()
	sql := "select  count(1)  from mobile_directory_sign where uid= ? AND contact_phone= ? AND sign_type=?  "
	if name != "" {
		sql += "  AND contact_name=?"
		o.Raw(sql, uid, moblie, stype, name).QueryRow(&count)
	} else {
		o.Raw(sql, uid, moblie, stype).QueryRow(&count)
	}
	return
}

//mobile_directory_sign插入记录
func InsertMoblieDirectory(uid, singId, signState int, name, makeContent, mobile, stype string) (err error) {
	o := orm.NewOrm()
	sql := " INSERT INTO mobile_directory_sign (uid, contact_name, contact_phone, sign_state, sign_remark, sign_time, sign_user_id,sign_type)VALUES (?,?, ?, ?, ?, ?, ?,?) "
	_, err = o.Raw(sql, uid, name, mobile, signState, makeContent, time.Now(), singId, stype).Exec()
	return
}

func SignMoblieDirectory(uid, singId, signState int, name, makeContent, mobile string) (num int64, err error) {
	o := orm.NewOrm()
	sql := `UPDATE mobile_directory_sign SET sign_state = ?, sign_remark = ?,sign_time = NOW(),sign_user_id=? WHERE  uid= ? AND contact_phone= ?  AND contact_name=? `
	res, err := o.Raw(sql, signState, makeContent, singId, uid, mobile, name).Exec()
	if err != nil {
		return 0, err
	}
	num, err = res.RowsAffected()
	return
}

//从mobile_directory_sign查询记录
func GetInfoByUidAndPhone(uid int, moblie, name string) (contact Contact, err error) {
	o := orm.NewOrm()
	sql := "select id as sign_id, sign_state,sign_remark  from mobile_directory_sign where uid= ? AND contact_phone= ?  AND contact_name=? AND sign_type=2 limit 1"
	err = o.Raw(sql, uid, moblie, name).QueryRow(&contact)
	return
}

//从mobile_directory_sign查询c常用联系人记录
func GetMxInfo(uid int, moblie string) (contact CallContactDetail, err error) {
	o := orm.NewOrm()
	sql := "select id as sign_id, sign_state,sign_remark  from mobile_directory_sign where uid= ? AND contact_phone= ? AND sign_type=3  limit 1"
	err = o.Raw(sql, uid, moblie).QueryRow(&contact)
	return
}

// 查询最大逾期天数
func GetMaxOverdueDays(loanId int) error {
	o := orm.NewOrm()
	var max_overdue_days int
	sql := `SELECT
			MAX(b.overdue_days) AS max_overdue_days
			FROM  repayment_schedule b
			WHERE b.loan_id=? AND b.state='BACKING' `
	o.Raw(sql, loanId).QueryRow(&max_overdue_days)
	upsql := ` UPDATE loan SET max_overdue_days= ?  WHERE id=? `
	_, err := o.Raw(upsql, max_overdue_days, loanId).Exec()
	return err
}

type NewMSum struct {
	Sumloanmoney float64 //借款总本金
	SumWaitMoney float64 //逾期待还总金额
	CostFee      float64
	Finishmoney  float64
}

//M管理-获取逾期应还总金额
func GetNewSumMoney(where, orgStr string, mtype int, pars []string) (m NewMSum) {
	sql := `SELECT
			SUM(b.capital_amount
				+ b.tax_amount
				+ b.overdue_money_amount
				+ b.overdue_breach_of_amount
				+ b.data_service_fee
				- b.remain_money_charge_up_amount
				- b.derate_amount
				) AS sum_wait_money,SUM(b.capital_amount) as sumloanmoney
			FROM loan AS a
			INNER JOIN repayment_schedule AS b ON a.id = b.loan_id
			WHERE 1=1 AND a.state='BACKING'  AND b.state='BACKING' AND a.mtype=? AND b.loan_return_date<=curdate() `
	if where != "" {
		sql += where
	}
	if orgStr != "" {
		sql += ` AND (a.collection_user_id=? OR a.org_id IN(` + orgStr + `))`
	} else {
		sql += ` AND a.collection_user_id=?`
	}
	o := orm.NewOrm()
	o.Using("read")
	o.Raw(sql, mtype, pars).QueryRow(&m)
	return
}

func GetNewRecycleMoney(condition, orgStr string, pars []string) (m NewMSum) {
	sql := `SELECT  SUM(b.capital_amount) AS sumloanmoney,
	                SUM(b.derate_amount) AS cost_fee,
	                SUM(IFNULL(b.return_overdue_money_amount,0)
				 + IFNULL(b.return_overdue_breach_of_amount,0)
				 + IFNULL(b.return_data_service_fee,0)
				 + IFNULL(b.return_tax_amount,0)
				 + IFNULL(b.return_capital_amount,0)
				 - IFNULL(b.derate_amount,0)) AS finishmoney
			FROM
				 loan AS a
			INNER JOIN repayment_schedule AS b ON a.id = b.loan_id
			WHERE  b.state = 'FINISH'    AND b.mtype >=11    `
	if condition != "" {
		sql += condition
	}
	if orgStr != "" {
		sql += ` AND (b.collection_user_id=? OR b.org_id IN(` + orgStr + `))`
	} else {
		sql += ` AND b.collection_user_id=?`
	}
	o := orm.NewOrm()
	o.Using("read")
	o.Raw(sql, pars).QueryRow(&m)
	return
}

//获取当前阶段每个催收员的案件数
func GetCaseCount(uid int) int {
	o := orm.NewOrm()
	sql := ` SELECT COUNT(1)
			FROM loan l
			WHERE l.collection_user_id=? AND l.state='backing' AND l.mtype >= 25 `
	var count int
	o.Raw(sql, uid).QueryRow(&count)
	return count
}
func UpdateLoanOrgId(orgId, uid int) (err error) {
	sql := `UPDATE loan SET org_id = ? WHERE  collection_user_id= ?  AND state='backing' `
	_, err = orm.NewOrm().Raw(sql, orgId, uid).Exec()
	return
}

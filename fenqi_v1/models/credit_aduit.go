package models

import (
	"fenqi_v1/utils"
	"strconv"
	"strings"
	"time"
	"zcm_tools/orm"
)

//授信列表
type CreditAduit struct {
	Id             int       `orm:"column(id)"`              //主键
	PhoneNo        string    `orm:"column(phone_no)"`        //电话号码
	Uid            int       `orm:"column(uid)"`             //用户id
	BalanceMoney   int       `orm:"column(balance_money)"`   //额度
	State          string    `orm:"column(state)"`           //信审状态
	CreateTime     time.Time `orm:"column(create_time)"`     //创建时间
	HandlingTime   time.Time `orm:"column(handling_time)"`   //处理时间
	OperatorId     int       `orm:"column(operator_id)"`     //处理人ID
	UserName       string    `orm:"column(user_name)"`       //用户姓名
	Displayname    string    `orm:"column(displayname)"`     //处理人
	OrderTime      time.Time `orm:"column(order_time)"`      //预约处理时间
	Remark         string    `orm:"column(remark)"`          //备注
	InqueueTime    time.Time `orm:"column(inqueue_time)"`    //入列时间
	InqueueType    int       `orm:"column(inqueue_type)"`    //入列状态
	AllocationTime time.Time `orm:"column(allocation_time)"` //分配时间
	IsOther        int       `orm:"column(is_other)"`        //其他资料
	IsMetadata     int       `orm:"column(is_myinfo)"`       //本人资料
	IsLinkMan      int       `orm:"column(is_linkman)"`      //联系人资料
	IsAuth         int       `orm:"column(is_auth)"`         //认证资料
	CreditAduitId  int       `orm:"column(credit_aduit_id)"` //授信主键
}

//获取授信列表
func GetCreditAduitList(start, pageSize int, condition string, paras ...interface{}) (list []CreditAduit, err error) {
	o := orm.NewOrm()
	sql := ` SELECT
				ca.id,
				ca.phone_no,
				ca.uid,
				ca.balance_money,
				ca.state,
				ca.create_time,
				ca.allocation_time,
				ca.handling_time,
				ca.inqueue_time,
				ca.user_name,
				ca.displayname
			FROM credit_aduit AS ca
			WHERE ca.is_now = 1`
	if condition != "" {
		sql += condition
	}
	sql += ` ORDER BY ca.create_time DESC LIMIT ?,?`
	_, err = o.Raw(sql, paras, start, pageSize).QueryRows(&list)
	return
}

//获取授信列表总数
func GetCreditAduitCount(condition string, paras ...interface{}) (count int, err error) {
	o := orm.NewOrm()
	sql := ` SELECT
				COUNT(1)
			FROM credit_aduit AS ca
			WHERE ca.is_now = 1`
	if condition != "" {
		sql += condition
	}
	err = o.Raw(sql, paras).QueryRow(&count)
	return
}

//获取信审人员
func GetCreditOperators() (displayname []string, err error) {
	sql := `SELECT distinct displayname FROM credit_aduit WHERE displayname IS NOT NULL`
	_, err = orm.NewOrm().Raw(sql).QueryRows(&displayname)
	return
}

//获取退回列表
func GetCreditOutQueueList(start, pageSize int, condition string, paras ...interface{}) (list []CreditAduit, err error) {
	o := orm.NewOrm()
	sql := ` SELECT
				ca.id,
				ca.uid,
				ca.phone_no,
				ca.balance_money,
				ca.state,
				ca.user_name,
				ca.handling_time,
				ca.inqueue_time,
				ca.remark,
				ca.create_time,
				su.displayname
			FROM credit_aduit AS ca
			INNER JOIN sys_user AS su
			ON ca.operator_id = su.id
			WHERE ca.is_now = 1 AND ca.state ="OUTQUEUE"
			AND ((ca.inqueue_time IS NOT NULL
			AND ca.inqueue_time > NOW())
			OR (ca.inqueue_type = 0))`
	if condition != "" {
		sql += condition
	}
	sql += ` ORDER BY ca.create_time DESC LIMIT ?,?`
	_, err = o.Raw(sql, paras, start, pageSize).QueryRows(&list)
	return
}

//获取退回列表总数
func GetCreditOutQueueCount(condition string, paras ...interface{}) (count int, err error) {
	o := orm.NewOrm()
	sql := ` SELECT
				COUNT(1)
			FROM credit_aduit AS ca
			INNER JOIN sys_user AS su
			ON ca.operator_id = su.id
			WHERE ca.is_now = 1 AND ca.state ="OUTQUEUE"
			AND ((ca.inqueue_time IS NOT NULL
			AND ca.inqueue_time > NOW())
			OR (ca.inqueue_type = 0))`
	if condition != "" {
		sql += condition
	}
	err = o.Raw(sql, paras).QueryRow(&count)
	return
}

//获取授信中预约的数量
func GetNumCreditQueueAp() (count int) {
	sql := `SELECT
			COUNT(1)
		FROM
			sys_user AS a,
			credit_aduit AS b
		WHERE
			a.id = b.operator_id
		AND b.state = "OUTQUEUE"`
	orm.NewOrm().Raw(sql).QueryRow(&count)
	return
}

//获取预约时间队列id
func GetAppointmentQueueId() (v CreditAduit, err error) {
	sql := `SELECT
				b.id,
				b.uid,
				b.inqueue_type,
				b.inqueue_time
			FROM
				credit_aduit AS b
			WHERE
				b.state = "OUTQUEUE"
			AND b.is_zone = 0
			AND b.inqueue_time IS NOT NULL
			AND b.inqueue_time <= NOW()
			AND b.inqueue_type > 0
			ORDER BY
				b.inqueue_time
			LIMIT 1`
	err = orm.NewOrm().Raw(sql).QueryRow(&v)
	return
}

//获取授信排队中信息id
func GetCreditQueueUpId() (v CreditAduit, err error) {
	sql := `SELECT
				b.id,
				b.uid
			FROM
				credit_aduit AS b
			WHERE
				b.state = "QUEUEING"
			AND b.is_zone = 0
			ORDER BY
				b.queue_time
			LIMIT 1`
	err = orm.NewOrm().Raw(sql).QueryRow(&v)
	return
}

//获取排队中数量
func GetCreditQueueUpIdCount() (count int) {
	sql := `SELECT COUNT(1) FROM credit_aduit as b WHERE b.state = 'QUEUEING' ORDER BY b.queue_time`
	orm.NewOrm().Raw(sql).QueryRow(&count)
	return
}

//更新授信中状态
func UpdateCueditQueueStatusInqueueTime(state, displayname string, operator_id, id int) error {
	sql := `UPDATE credit_aduit SET state=?,displayname=?,inqueue_time=NULL,operator_id=? WHERE id =? `
	_, err := orm.NewOrm().Raw(sql, state, displayname, operator_id, id).Exec()
	return err
}

//更新授信中记录
func UpdateCueditQueueRemark(remark string, id int) error {
	sql := `UPDATE credit_aduit SET remark = ? WHERE id =? `
	_, err := orm.NewOrm().Raw(sql, remark, id).Exec()
	return err
}

//更新授信中处理人状态
func UpdateCueditQueueStatusOp(state string, id, operator_id int) error {
	sql := `UPDATE credit_aduit SET state=?,operator_id=?,handling_time=NOW() WHERE id =? `
	_, err := orm.NewOrm().Raw(sql, state, operator_id, id).Exec()
	return err
}

//更新授信中处理人状态
func UpdateCutCueditQueueStatusOp(state, displayname string, id, operator_id int, handling_time time.Time) error {
	sql := `UPDATE credit_aduit SET state=?,operator_id=?,handling_time=?,displayname=? WHERE id =? `
	_, err := orm.NewOrm().Raw(sql, state, operator_id, handling_time, displayname, id).Exec()
	return err
}

//Pass状态下更新授信
func UpdateCueditQueuePassStatus(state string, balance_money, operator_id, id int) error {
	sql := `UPDATE credit_aduit SET state=?,balance_money=?,operator_id=?,handling_time=NOW() WHERE id=? `
	_, err := orm.NewOrm().Raw(sql, state, balance_money, operator_id, id).Exec()
	return err
}

//Pass状态下更新授信
func UpdateCutCueditQueuePassStatus(state, displayname string, balance_money, operator_id, id int, handling_time time.Time) error {
	sql := `UPDATE credit_aduit SET state=?,balance_money=?,operator_id=?,handling_time=?,displayname=? WHERE id=? `
	_, err := orm.NewOrm().Raw(sql, state, balance_money, operator_id, handling_time, displayname, id).Exec()
	return err
}

//更新退回队列中入队类型(排队,插队)
func UpdateInqueueType(id int, queue_time time.Time) (err error) {
	sql := `UPDATE credit_aduit SET state = "QUEUEING",queue_time=? WHERE id =?`
	_, err = orm.NewOrm().Raw(sql, queue_time, id).Exec()
	return
}

//更新退回队列时间与状态
func UpdateCreditOutqueueTime(id, operator_id int, inqueue_time string) (err error) {
	sql := `UPDATE credit_aduit SET state = "OUTQUEUE"`
	if inqueue_time == "" {
		sql += `,inqueue_time = NULL,operator_id=?,inqueue_type=0,handling_time=NOW() WHERE id =?`
		_, err = orm.NewOrm().Raw(sql, operator_id, id).Exec()

	} else {
		sql += `,inqueue_time=?,operator_id=?,inqueue_type=1,handling_time=NOW() WHERE id =?`
		_, err = orm.NewOrm().Raw(sql, inqueue_time, operator_id, id).Exec()
	}
	return err
}

//更新分配時間
func UpdateCreditAlloctionTime(id int) error {
	sql := `UPDATE credit_aduit SET allocation_time = NOW(),inqueue_time = NULL WHERE id = ? `
	_, err := orm.NewOrm().Raw(sql, id).Exec()
	return err
}

//统计45分钟之内
func QueryCreditHandingCountIn(id int) (count int) {
	sql := `SELECT COUNT(1) FROM credit_aduit WHERE id =? AND allocation_time >= DATE_SUB(NOW(),INTERVAL 45 MINUTE)`
	orm.NewOrm().Raw(sql, id).QueryRow(&count)
	return
}

//超时清缓存并入队列
func UpdateCreditQueueing(id int) error {
	sql := `UPDATE credit_aduit SET state = "QUEUEING",queue_time = NOW() WHERE id = ? `
	_, err := orm.NewOrm().Raw(sql, id).Exec()
	return err
}

//超时分配查询attime时间
func QueryCreditAttime(id int) (v CreditAduit) {
	sql := "SELECT allocation_time FROM credit_aduit WHERE id=?"
	orm.NewOrm().Raw(sql, id).QueryRow(&v)
	return
}

//超时分配查询attime时间
func UpdateCreditAttime(id int, allocation_time time.Time) (err error) {
	sql := "UPDATE credit_aduit SET allocation_time=? WHERE id=?"
	_, err = orm.NewOrm().Raw(sql, allocation_time, id).Exec()
	return
}

//选择退回队列中需要授信名单 更新入列状态
func UpadatecCredintInQueue(inqueue_time string, inqueue_type, id int) (err error) {
	o := orm.NewOrm()
	err = o.Begin()
	if err != nil {
		return err
	}
	sql := `UPDATE credit_aduit SET inqueue_type = ? ,`
	if inqueue_time == "" {
		sql += `inqueue_time= NULL `
	} else {
		sql += `inqueue_time= ? `
	}
	sql += ` WHERE id = ?`
	p, err := o.Raw(sql).Prepare()
	defer p.Close()
	if inqueue_time == "" {
		_, err = p.Exec(inqueue_type, id)
	} else {
		_, err = p.Exec(inqueue_type, inqueue_time, id)
	}
	if err != nil {
		o.Rollback()
		return err
	}
	o.Commit()
	return nil
}

//选择退回队列中需要授信名单 更新入列状态
func UpadatecCutCredintInQueue(inqueue_time string, inqueue_type, id int) (err error) {

	var table = [5]string{"credit_aduit", "credit_auth_queue", "credit_linkman_queue", "credit_myinfo_queue", "credit_other_queue"}
	for i := 0; i < 5; i++ {
		o := orm.NewOrm()
		err = o.Begin()
		if err != nil {
			return err
		}
		sql := `UPDATE ` + table[i] + ` SET inqueue_type = ?,`
		if table[i] == "credit_aduit" {
			sql += `state="OUTQUEUE",is_myinfo=0,is_auth=0,is_linkman=0,is_other=0,`
		}
		if table[i] == "credit_auth_queue" {
			sql += `auth_state="OUTQUEUE",auth_name="",remark="",handling_time=NULL,`
		}
		if table[i] == "credit_linkman_queue" {
			sql += `linkman_state="OUTQUEUE",linkman_name="",remark="",handling_time=NULL,`
		}
		if table[i] == "credit_myinfo_queue" {
			sql += `myinfo_state="OUTQUEUE",myinfo_name="",remark="",handling_time=NULL,`
		}
		if table[i] == "credit_other_queue" {
			sql += `other_state="OUTQUEUE",other_name="",remark="",handling_time=NULL,`
		}
		if inqueue_time == "" {
			sql += `inqueue_time= NULL`
		} else {
			sql += `inqueue_time= ?`
		}
		if table[i] == "credit_aduit" {
			sql += ` WHERE id = ?`
		} else {
			sql += ` WHERE credit_aduit_id = ?`
		}
		p, err := o.Raw(sql).Prepare()
		defer p.Close()
		if inqueue_time == "" {
			_, err = p.Exec(inqueue_type, id)
		} else {
			_, err = p.Exec(inqueue_type, inqueue_time, id)
		}
		if err != nil {
			o.Rollback()
			return err
		}
		o.Commit()
	}
	return nil
}

//45之内 处于睡眠 关机
func CreditHandingLogOut() {
	var id []int
	sql := `SELECT id FROM credit_aduit WHERE state = "HANDING"`
	orm.NewOrm().Raw(sql).QueryRows(&id)
	for _, v := range id {
		if !utils.Rc.IsExist("xjfq:" + utils.CacheKeyCreditMessage + "_" + strconv.Itoa(v)) {
			sql := `UPDATE credit_aduit SET state = "QUEUEING",queue_time = NOW() WHERE id = ? `
			orm.NewOrm().Raw(sql, v).Exec()
		}
	}
}

func SelectSystemIDHanding(operator_id int) (id, uid int) {
	sql := `SELECT id FROM credit_aduit WHERE state = "HANDING" AND operator_id =?`
	orm.NewOrm().Raw(sql, operator_id).QueryRow(&id)
	if id > 0 {
		sql := `SELECT uid FROM credit_aduit WHERE id =? `
		orm.NewOrm().Raw(sql, id).QueryRow(&uid)
	}
	return
}

//授信信息
type CreditInfo struct {
	Id           int
	Uid          int     //用户ID
	BalanceMoney float64 //审批额度
	Remark       string  //备注
}

//根据id获取授信信息
func GetCreditInfoByUid(id int) (info CreditInfo, err error) {
	sql := `SELECT id,uid,balance_money,remark FROM credit_aduit WHERE id = ?`
	err = orm.NewOrm().Raw(sql, id).QueryRow(&info)
	return
}

//裁分授信id查询状态
func GetCreditIsOk(id int) (flag int) {
	var v CreditAduit
	sql := `SELECT is_myinfo,is_auth,is_linkman,is_other FROM credit_aduit WHERE id = ? AND is_zone != 0`
	orm.NewOrm().Raw(sql, id).QueryRow(&v)
	return v.IsOther + v.IsMetadata + v.IsLinkMan + v.IsAuth
}

//授信信息
type CreditCutState struct {
	AuthState     string `orm:"column(auth_state)"`    //认证状态
	LinkManState  string `orm:"column(linkman_state)"` //联系人状态
	OtherState    string `orm:"column(other_state)"`   //其他状态
	MetaDataState string `orm:"column(myinfo_state)"`  //基本信息状态
}

//授信money
type CreditAduitMoney struct {
	AuthBalanceMoney    int `orm:"column(auth_balance_money)"`    //额度认证money
	LinkManBalanceMoney int `orm:"column(linkman_balance_money)"` //额度联系人money
	MetaBalanceMoney    int `orm:"column(myinfo_balance_money)"`  //额度基本信息money
	OtherBalanceMoney   int `orm:"column(other_balance_money)"`   //额度其他money
}

//处理拆分id状态
func GetCreditIsCutStatus(credit_aduit_id int) (cutstate, displayname string) {
	var v CreditCutState
	sql := `SELECT
				a.auth_state,
				b.linkman_state,
				c.myinfo_state,
				d.other_state
			FROM
				credit_auth_queue a,
				credit_linkman_queue b,
				credit_myinfo_queue c,
				credit_other_queue d
			WHERE
				a.credit_aduit_id = b.credit_aduit_id
			AND c.credit_aduit_id = b.credit_aduit_id AND d.credit_aduit_id = c.credit_aduit_id
			AND d.credit_aduit_id=?`
	orm.NewOrm().Raw(sql, credit_aduit_id).QueryRow(&v)
	cutstate = v.AuthState + v.LinkManState + v.OtherState + v.MetaDataState
	return cutstate, GetDisplayNameStatus(credit_aduit_id)
}

//获取拆分下状态
func GetCreditCutState(cutstate string) (state string) {
	var status = []string{"CLOSE", "PAUSE", "REJECT", "OUTQUEUE", "PASS"}
	for _, v := range status {
		if strings.Contains(cutstate, v) {
			return v
		}
	}
	return
}

//裁分方案全部通过下额度
func GetCreditCutPassMoney(credit_aduit_id int) (balance_money int) {
	var b CreditAduitMoney
	sql := `SELECT
				a.auth_balance_money,
				b.linkman_balance_money,
				c.myinfo_balance_money,
				d.other_balance_money
			FROM
				credit_auth_queue a,
				credit_linkman_queue b,
				credit_myinfo_queue c,
				credit_other_queue d
			WHERE
				a.credit_aduit_id = b.credit_aduit_id
				AND c.credit_aduit_id = b.credit_aduit_id AND d.credit_aduit_id = c.credit_aduit_id
				AND d.credit_aduit_id =?`
	orm.NewOrm().Raw(sql, credit_aduit_id).QueryRow(&b)
	pass_money := getcreditcutpassmoney(b)
	return pass_money
}

//处理裁分后Pass状态下money
func getcreditcutpassmoney(b CreditAduitMoney) (balance_money int) {
	var pass_money int = 0
	if b.AuthBalanceMoney > pass_money {
		pass_money = b.AuthBalanceMoney
	}
	if b.LinkManBalanceMoney < pass_money {
		pass_money = b.LinkManBalanceMoney
	}
	if b.OtherBalanceMoney < pass_money {
		pass_money = b.OtherBalanceMoney
	}
	if b.MetaBalanceMoney < pass_money {
		pass_money = b.MetaBalanceMoney
	}
	return pass_money
}

//查看授信订单属于那种方案
func GetCreditIsCut(id int) (is_zone int) {
	sql := `SELECT is_zone FROM credit_aduit WHERE id = ?`
	orm.NewOrm().Raw(sql, id).QueryRow(&is_zone)
	return
}

type OpHandlingTime struct {
	HandlingTime   time.Time `orm:"column(handling_time)"`   //处理时间
	OperatorId     int       `orm:"column(operator_id)"`     //处理人ID
	Displayname    string    `orm:"column(displayname)"`     //处理人
	Remark         string    `orm:"column(remark)"`          //备注
	AllocationTime time.Time `orm:"column(allocation_time)"` //分配时间
}

//授信处理
func getauthhandtime(state string, uid, credit_aduit_id int) (op OpHandlingTime, err error) {
	sql := `SELECT handling_time,operator_id,auth_name as displayname,remark,allocation_time FROM credit_auth_queue WHERE auth_state=? AND uid=? And credit_aduit_id=?`
	orm.NewOrm().Raw(sql, state, uid, credit_aduit_id).QueryRow(&op)
	return
}

func getlinkmanhandtime(state string, uid, credit_aduit_id int) (op OpHandlingTime, err error) {
	sql := `SELECT handling_time,operator_id,linkman_name as displayname,remark,allocation_time FROM credit_linkman_queue WHERE linkman_state=? AND uid=? AND credit_aduit_id=?`
	orm.NewOrm().Raw(sql, state, uid, credit_aduit_id).QueryRow(&op)
	return
}

func getmyinfohandtime(state string, uid, credit_aduit_id int) (op OpHandlingTime, err error) {
	sql := `SELECT handling_time,operator_id,myinfo_name as displayname,remark,allocation_time FROM credit_myinfo_queue WHERE myinfo_state=? AND uid=? AND credit_aduit_id=?`
	orm.NewOrm().Raw(sql, state, uid, credit_aduit_id).QueryRow(&op)
	return
}

func getotherhandtime(state string, uid, credit_aduit_id int) (op OpHandlingTime, err error) {
	sql := `SELECT handling_time,operator_id,other_name as displayname,remark,allocation_time FROM credit_other_queue WHERE other_state=? AND uid=? And credit_aduit_id=?`
	orm.NewOrm().Raw(sql, state, uid, credit_aduit_id).QueryRow(&op)
	return
}

func HandlingTimeOp(state string, uid, credit_aduit_id int) (handling_time, at_time time.Time, oper_id int, displayname, remark string) {
	ap, _ := getauthhandtime(state, uid, credit_aduit_id)
	if ap.HandlingTime.Format(utils.FormatDate) != "0001-01-01" && ap.OperatorId != 0 {
		handling_time = ap.HandlingTime
		oper_id = ap.OperatorId
		displayname = ap.Displayname
		remark = ap.Remark
		at_time = ap.AllocationTime
	}

	lp, _ := getlinkmanhandtime(state, uid, credit_aduit_id)
	if lp.HandlingTime.Format(utils.FormatDate) != "0001-01-01" && lp.HandlingTime.After(handling_time) {

		handling_time = lp.HandlingTime
		oper_id = lp.OperatorId
		displayname = lp.Displayname
		remark = lp.Remark
		at_time = lp.AllocationTime
	}
	mp, _ := getmyinfohandtime(state, uid, credit_aduit_id)
	if mp.HandlingTime.Format(utils.FormatDate) != "0001-01-01" && mp.HandlingTime.After(handling_time) {
		handling_time = mp.HandlingTime
		oper_id = mp.OperatorId
		displayname = mp.Displayname
		remark = mp.Remark
		at_time = mp.AllocationTime
	}

	op, _ := getotherhandtime(state, uid, credit_aduit_id)
	if op.HandlingTime.Format(utils.FormatDate) != "0001-01-01" && op.HandlingTime.After(handling_time) {
		handling_time = op.HandlingTime
		oper_id = op.OperatorId
		displayname = op.Displayname
		remark = op.Remark
		at_time = op.AllocationTime
	}
	return handling_time, at_time, oper_id, displayname, remark
}

type CreditAdvise struct {
	HandlingTime time.Time `orm:"column(handling_time)"` //处理时间
	Remark       string    `orm:"column(remark)"`        //备注
	Displayname  string    `orm:"column(displayname)"`   //处理人
	State        string    `orm:"column(state)"`         //信审状态
	Msg          string    //审核信息
}

//授信处理
func getauthadvise(credit_aduit_id int) (op CreditAdvise, err error) {
	sql := `SELECT
				a.handling_time,
				a.remark,
				a.auth_name AS displayname,
				a.auth_state AS state
			FROM
				credit_auth_queue a
			WHERE
				a.credit_aduit_id = ?`
	orm.NewOrm().Raw(sql, credit_aduit_id).QueryRow(&op)
	return
}

func getlinkmanadvise(credit_aduit_id int) (op CreditAdvise, err error) {
	sql := `SELECT
				a.handling_time,
				a.remark,
				a.linkman_name AS displayname,
				a.linkman_state AS state

			FROM
				credit_linkman_queue a
			WHERE
				a.credit_aduit_id = ?`
	orm.NewOrm().Raw(sql, credit_aduit_id).QueryRow(&op)
	return
}

func getmyinfoadvise(credit_aduit_id int) (op CreditAdvise, err error) {
	sql := `SELECT
				a.handling_time,
				a.remark,
				a.myinfo_name AS displayname,
				a.myinfo_state AS state
			FROM
				credit_myinfo_queue a
			WHERE
				a.credit_aduit_id =?`
	orm.NewOrm().Raw(sql, credit_aduit_id).QueryRow(&op)
	return
}

func getotheradvise(credit_aduit_id int) (op CreditAdvise, err error) {
	sql := `SELECT
				a.handling_time,
				a.remark,
				a.other_name AS displayname,
				a.other_state AS state
			FROM
				credit_other_queue a
			WHERE
				a.credit_aduit_id = ?`
	orm.NewOrm().Raw(sql, credit_aduit_id).QueryRow(&op)
	return
}

func GetCreditAdvise(credit_aduit_id int) (ca []CreditAdvise, err error) {
	ap, _ := getauthadvise(credit_aduit_id)
	if ap.State != "" && (ap.HandlingTime.Format(utils.FormatDate) != "0001-01-01") {
		ap.Msg = "授权信息"
		ca = append(ca, ap)
	}
	lp, _ := getlinkmanadvise(credit_aduit_id)
	if lp.State != "" && (lp.HandlingTime.Format(utils.FormatDate) != "0001-01-01") {
		lp.Msg = "联系人信息"
		ca = append(ca, lp)
	}
	mp, _ := getmyinfoadvise(credit_aduit_id)
	if mp.State != "" && (mp.HandlingTime.Format(utils.FormatDate) != "0001-01-01") {
		mp.Msg = "基本信息"
		ca = append(ca, mp)
	}
	op, _ := getotheradvise(credit_aduit_id)
	if op.State != "" && (op.HandlingTime.Format(utils.FormatDate) != "0001-01-01") {
		op.Msg = "其他信息"
		ca = append(ca, op)
	}
	return
}

func GetDisplayNameStatus(credit_aduit_id int) (displayname string) {
	ca, _ := GetCreditAdvise(credit_aduit_id)
	for k, v := range ca {
		if v.Displayname != "" && !strings.Contains(displayname, v.Displayname) {
			if k < len(ca)-1 {
				displayname += v.Displayname + "、"
			} else {
				displayname += v.Displayname
			}

		}
	}
	return displayname
}

func GetCreditId(uid int) (id int, err error) {
	sql := `select id from credit_aduit where uid =? and is_now = 1`
	err = orm.NewOrm().Raw(sql, uid).QueryRow(&id)
	return
}

func GetCreditState(uid int) (state string, err error) {
	sql := `select state from credit_aduit where uid =? and is_now= 1`
	err = orm.NewOrm().Raw(sql, uid).QueryRow(&state)
	return
}

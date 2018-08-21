package models

import (
	"fenqi_v1/utils"
	"github.com/astaxie/beego"
	"strconv"
	"time"
	"zcm_tools/orm"
)

type Salesman struct {
	Id                 int       `orm:"column(id);pk"`
	Account            string    `orm:"column(account)"`     // 业务员手机号
	Password           string    `orm:"column(password)"`    // 密码
	InviteCode         string    `orm:"column(invite_code)"` // 邀请码
	CreateTime         time.Time `orm:"column(create_time)"`
	AllotmentTime      time.Time
	Saleman            string `orm:"column(saleman)"` // 名字
	IsOk               int
	RegionId           int
	PlaceId            int
	DepId              int
	Region             string // 大区
	Place              string // 地区
	OperaDep           string // 营业部
	NowAllotment       int
	AllAllotment       int
	AlltomentTime      time.Time
	IsLinkAllotment    int
	NotIsLinkAllotment int
	OrgId              int
	Name               string
	Group              string
	StaId              int
	GroupId            int
}

//获取大区
func GetSalemanRegion() (s []Salesman, err error) {
	sql := `SELECT name FROM sys_organization WHERE parent_id = (SELECT id FROM sys_organization WHERE name = "业务员") `
	_, err = orm.NewOrm().Raw(sql).QueryRows(&s)
	return
}

//获取省份和运营部
func GetSalemanDep(id int) (s []Salesman, err error) {
	sql := `SELECT region,place,opera_dep FROM sys_organization WHERE parent_id = ?`
	_, err = orm.NewOrm().Raw(sql, id).QueryRows(&s)
	return
}

//获取业务员列表
func GetSalesmanList(page, pageSize int, condition string, orgStr string, para ...interface{}) (v []*Salesman, err error) {
	o := orm.NewOrm()
	sql := `SELECT * FROM salesman s
			WHERE 1 = 1 `
	if condition != "" {
		sql += condition
	}
	if orgStr != "" {
		sql += ` AND s.org_id IN(` + orgStr + `)`
	}
	sql += ` ORDER BY create_time DESC LIMIT ?,?`
	_, err = o.Raw(sql, para, page, pageSize).QueryRows(&v)
	beego.Info(err, sql, para)
	return
}

//获取业务员列表
func GetSalesmanLookList(page, pageSize int, condition string, orgStr string, para ...interface{}) (v []*Salesman, err error) {
	o := orm.NewOrm()
	sql := `SELECT s.id as id,s.account,s.invite_code,s.create_time,s.saleman,s.dep_id,s.place_id,s.is_ok,s.region_id,sta.name,s.group_id,s.sta_id FROM salesman s
			LEFT JOIN sys_station sta ON s.sta_id = sta.id
			WHERE 1 = 1 `
	if condition != "" {
		sql += condition
	}
	if orgStr != "" {
		sql += ` AND s.org_id IN(` + orgStr + `)`
	}
	sql += ` ORDER BY s.create_time DESC LIMIT ?,?`
	_, err = o.Raw(sql, para, page, pageSize).QueryRows(&v)
	return
}

//获取业务员count
func GetSalesmanLookCount(condition string, orgStr string, para ...interface{}) (c int, err error) {
	o := orm.NewOrm()
	sql := `SELECT COUNT(1) FROM salesman s
			LEFT JOIN sys_station sta ON s.sta_id = sta.id 
			WHERE 1 = 1 `
	if condition != "" {
		sql += condition
	}
	if orgStr != "" {
		sql += ` AND s.org_id IN(` + orgStr + `)`
	}
	err = o.Raw(sql, para).QueryRow(&c)
	return
}

//获取业务员count
func GetSalesmanCount(condition string, orgStr string, para ...interface{}) (c int, err error) {
	o := orm.NewOrm()
	sql := `SELECT COUNT(1) FROM salesman
			WHERE 1 = 1 `
	if condition != "" {
		sql += condition
	}
	if orgStr != "" {
		sql += ` AND org_id IN(` + orgStr + `)`
	}
	err = o.Raw(sql, para).QueryRow(&c)
	return
}

//根据id查询业务员
func GetSalemanById(id int) (v Salesman, err error) {
	o := orm.NewOrm()
	sql := `SELECT * FROM salesman WHERE id = ?`
	err = o.Raw(sql, id).QueryRow(&v)
	beego.Info(id, err)
	return
}

func GetSalemanStationById(id int) (v Salesman, err error) {
	o := orm.NewOrm()
	sql := `SELECT s.* FROM salesman s 
			WHERE s.id = ?`
	err = o.Raw(sql, id).QueryRow(&v)
	return
}

type SalesmanUsers struct {
	Uid               int
	CreateTime        time.Time
	Verifyrealname    string
	InviteAccount     string
	Money             string
	LoanTermCount     string
	CurrentAuthState  int
	IsRealName        int
	State             string
	UsersBaseInfoTime time.Time
	RealNameTime      time.Time
	SubmitTime        time.Time
}

func GetUsersIsAuth(page, pageSize, id int) (u []*SalesmanUsers, err error) {
	o := orm.NewOrm()
	sql := `SELECT i.uid,i.create_time,i.invite_account,l.loan_term_count,l.money,u.current_auth_state,u.submit_time,um.verifyrealname,u.real_name_time,u.users_base_info_time,u.is_real_name FROM salesman AS s 
			LEFT JOIN salesman_invite AS i ON s.id = i.salesman_id 
			LEFT JOIN loan AS l ON i.uid = l.uid  AND l.is_now = 1 
			INNER JOIN users_auth AS u ON i.uid = u.uid  AND u.is_valid= 1
			INNER JOIN users_metadata AS um ON i.uid = um.uid AND um.verifyrealname IS NOT NULL AND um.id_card IS NOT NULL
			WHERE s.id= ? and i.state = 1 ORDER BY i.create_time DESC LIMIT ?,?`
	beego.Info(id, page, pageSize)
	_, err = o.Raw(sql, id, (page-1)*pageSize, pageSize).QueryRows(&u)
	return
}
func GetUsersIsAuthCount(id int) (count int, err error) {
	o := orm.NewOrm()
	sql := `SELECT count(1) FROM salesman AS s 
			LEFT JOIN salesman_invite AS i ON s.id = i.salesman_id 
			LEFT JOIN loan AS l ON i.uid = l.uid  AND l.is_now = 1 
			INNER JOIN users_auth AS u ON i.uid = u.uid    AND u.is_valid= 1
			INNER JOIN users_metadata AS um ON i.uid = um.uid AND um.verifyrealname IS NOT NULL AND um.id_card IS NOT NULL
			WHERE s.id= ? and i.state = 1`
	err = o.Raw(sql, id).QueryRow(&count)
	return
}
func GetUsersNotIsAuth(page, pageSize, id int) (u []SalesmanUsers, err error) {
	o := orm.NewOrm()
	sql := `SELECT i.uid,i.create_time,i.invite_account,l.loan_term_count,l.money,um.verifyrealname  FROM salesman AS s 
			LEFT JOIN salesman_invite AS i ON s.id = i.salesman_id 
			LEFT JOIN loan AS l ON i.uid = l.uid  AND l.is_now = 1 
			INNER JOIN users_metadata AS um ON i.uid = um.uid  AND um.verifyrealname IS NULL AND um.id_card IS NULL
			WHERE s.id= ? and i.state = 1  ORDER BY i.create_time DESC LIMIT ?,?`
	_, err = o.Raw(sql, id, (page-1)*pageSize, pageSize).QueryRows(&u)
	return
}
func GetUsersIsNotAuthCount(id int) (count int, err error) {
	o := orm.NewOrm()
	sql := `SELECT count(1) FROM salesman AS s 
			LEFT JOIN salesman_invite AS i ON s.id = i.salesman_id 
			LEFT JOIN loan AS l ON i.uid = l.uid AND l.is_now = 1
			INNER JOIN users_metadata AS um ON i.uid = um.uid  AND um.verifyrealname IS  NULL AND um.id_card IS NULL
			WHERE s.id= ? and i.state = 1`
	err = o.Raw(sql, id).QueryRow(&count)
	return
}

//增加用户
func AddSaleman(s *Salesman) error {
	o := orm.NewOrm()
	sql := `insert into salesman(account,password,invite_code,create_time,saleman,dep_id,is_ok,region_id,place_id,org_id,sta_id,group_id) values(?,?,?,?,?,?,?,?,?,?,?,?)`
	_, err := o.Raw(sql, s.Account, utils.MD5(s.Password), s.InviteCode, s.CreateTime, s.Saleman, s.DepId, s.IsOk, s.RegionId, s.PlaceId, s.OrgId, s.StaId, s.GroupId).Exec()
	beego.Info(sql, err, s)
	return err
}

//编辑业务员
func UpdateSalesman(s *Salesman) error {
	o := orm.NewOrm()
	sql := `update salesman set account=?,password=?,invite_code=?,saleman=?,dep_id=?,is_ok=?,region_id=?,place_id=?,org_id=?,sta_id=?,group_id=? where id = ?`
	_, err := o.Raw(sql, s.Account, utils.MD5(s.Password), s.InviteCode, s.Saleman, s.DepId, s.IsOk, s.RegionId, s.PlaceId, s.OrgId, s.StaId, s.GroupId, s.Id).Exec()
	beego.Info(sql, err, s)
	return err
}

//业务员组织架构
func GetSalemanOrganizationStations() ([]map[string]interface{}, error) {
	sql := `SELECT id,parent_id,invitation_code_prefix,name,1 nocheck  FROM sys_organization WHERE id >27
			UNION ALL 
			SELECT id+100000 id,org_id AS parent_id ,"" AS invitation_code_prefix,name,0 nocheck FROM sys_station  WHERE id >28 `
	list := []SysOrganization{}

	_, err := orm.NewOrm().Raw(sql).QueryRows(&list)
	l := len(list)
	var org []map[string]interface{}
	for i := 0; i < l; i++ {
		org = append(org, map[string]interface{}{"id": list[i].Id, "pId": list[i].ParentId, "open": false, "name": list[i].Name, "invitation_code_prefix": list[i].InvitationCodePrefix, "nocheck": list[i].NoCheck})
	}
	return org, err
}

func QuerySalemanDisplayQn() (sList []SysStation) {
	var sys_1 []SysOrganization
	o := orm.NewOrm()
	inStr := ""
	sql := `select id,name,parent_id from sys_organization 
	        where parent_id in (select id from sys_organization where 1=1)`
	o.Raw(sql).QueryRows(&sys_1)
	if len(sys_1) > 0 {
		for i := 0; i < len(sys_1); i++ {
			inStr += strconv.Itoa(sys_1[i].Id) + ","
		}
		inStr2 := inStr[:len(inStr)-1]
		sql2 := `select id,name,parent_id from sys_organization 
	        where parent_id in(` + inStr2 + `)`
		var sys_2 []SysOrganization
		o.Raw(sql2).QueryRows(&sys_2)
		if len(sys_2) > 0 {
			inStr3 := ""
			for i := 0; i < len(sys_2); i++ {
				inStr += strconv.Itoa(sys_2[i].Id) + ","
				inStr3 += strconv.Itoa(sys_2[i].Id) + ","
			}
			inStr3 = inStr3[:len(inStr3)-1]
			sql3 := `select id,name,parent_id from sys_organization 
	        where parent_id in(` + inStr3 + `)`
			var sys_3 []SysOrganization
			o.Raw(sql3).QueryRows(&sys_3)
			if len(sys_3) > 0 {
				inStr4 := ""
				for i := 0; i < len(sys_3); i++ {
					inStr += strconv.Itoa(sys_3[i].Id) + ","
					inStr4 += strconv.Itoa(sys_3[i].Id) + ","
				}
				inStr4 = inStr4[:len(inStr4)-1]
				sql4 := `select id,name,parent_id from sys_organization
				         where parent_id in(` + inStr4 + `)`
				var sys_4 []SysOrganization
				o.Raw(sql4).QueryRows(&sys_4)
				if len(sys_4) > 0 {
					// inStr5:=""
					for i := 0; i < len(sys_4); i++ {
						inStr += strconv.Itoa(sys_4[i].Id) + ","
						// inStr5 += strconv.Itoa(sys_4[i].Id) + ","
					}
					// inStr5 = inStr5[:len(inStr4)-1]
					// sql5 := `select id,name,parent_id from sys_organization
					//         where parent_id in(` + inStr5 + `)`
					// var sys_5 []SysOrganization
					// o.Raw(sql5).QueryRows(&sys_5)
					// for i := 0; i < len(sys_5); i++ {
					// 	inStr += strconv.Itoa(sys_5[i].Id) + ","
					// }
				}
			}
		}
	}
	var sys SysOrganization
	o.Raw(`select id from sys_organization where 1=1 `).QueryRow(&sys)
	inStr += strconv.Itoa(sys.Id)
	stationSql := `select * from sys_station where org_id in (` + inStr + `)`
	o.Raw(stationSql).QueryRows(&sList)
	return
}

func GetOrgPid(id int) (s *SysOrganization, err error) {
	sql := `select * from sys_organization where id = ?`
	err = orm.NewOrm().Raw(sql, id).QueryRow(&s)
	return
}

func GetSaleman() (v []Salesman, err error) {
	sql := `SELECT * FROM salesman`
	_, err = orm.NewOrm().Raw(sql).QueryRows(&v)
	return
}

//获取业务员分配列表
func GetSalesmanAllotmentList(page, pageSize int, condition string, para ...interface{}) (v []*Salesman, err error) {
	o := orm.NewOrm()
	sql := `SELECT s.id,s.account,s.dep_id,s.region_id,s.place_id s.create_time,s.saleman ,s.all_allotment,s.now_allotment FROM salesman s
			WHERE 1 = 1 `
	if condition != "" {
		sql += condition
	}
	sql += ` group by sa.salesman_id  ORDER BY s.create_time DESC LIMIT ?,?`
	_, err = o.Raw(sql, para, page, pageSize).QueryRows(&v)
	beego.Info(err, sql, para)
	return
}

//获取业务员分配总数
func GetAllotmentCount(condition string, para ...interface{}) (count int, err error) {
	o := orm.NewOrm()
	sql := `SELECT COUNT(1) FROM salesman s
			LEFT JOIN salesman_allotment sa ON s.id = sa.salesman_id
			WHERE 1 = 1 AND sa.pm_id= ?`
	if condition != "" {
		sql += condition
	}
	sql += ` group by sa.salesman_id  ORDER BY sa.create_time`
	err = o.Raw(sql, para).QueryRow(&count)
	beego.Info(err, sql, para)
	return
}

func GetSalesmanAll(dep_id int) (s *[]Salesman, err error) {
	o := orm.NewOrm()
	sql := `SELECT * FROM salesman s
			WHERE dep_id = ? `
	_, err = o.Raw(sql, dep_id).QueryRows(&s)
	return
}

func GetSalesmanAllotmentCount(pmId int) (count int, err error) {
	o := orm.NewOrm()
	sql := `SELECT COUNT(1) FROM salesman_allotment WHERE pm_id = ? AND salesman_id=0`

	err = o.Raw(sql, pmId).QueryRow(&count)

	return
}

func GetAllotmentList(pmId int) (sa []*SalesmanAllotment, err error) {
	o := orm.NewOrm()
	sql := `SELECT * FROM salesman_allotment WHERE pm_id = ? AND salesman_id=0`
	_, err = o.Raw(sql, pmId).QueryRows(&sa)
	return
}

func UpdataAllotment(shares []Share) (err error) {
	o := orm.NewOrm()
	err = o.Begin()
	if err != nil {
		beego.Info(err)
		return err
	}
	sql := `UPDATE salesman SET all_allotment = ? ,now_allotment= ?,allotment_time= ? WHERE id = ?`
	p, err := o.Raw(sql).Prepare()
	if err != nil {
		beego.Info(err)
		return err
	}
	defer p.Close()

	beego.Info(shares)
	for _, v := range shares {
		a := v.Id
		b := v.New
		c := v.All
		d := v.AllotmentId
		f := "(" + d + ")"
		sql1 := `UPDATE salesman_allotment SET create_time = ? ,salesman_id = ?  WHERE id IN  ` + f
		p1, err := o.Raw(sql1).Prepare()
		if err != nil {
			beego.Info(err)
			return err
		}
		defer p1.Close()
		_, err = p.Exec(c, b, time.Now(), a)
		_, err = p1.Exec(time.Now(), a)
		if err != nil {
			o.Rollback()
			return err
		}
	}
	if err == nil {
		o.Commit()
	}
	return
}

type SalesmanAllotment struct {
	Id         int
	SalesmanId int
	CreateTime time.Time
	PmId       int
	Mark       string
	PhoneNum   string
	UserName   string
	Place      string
	State      int
	HandleTime time.Time
}

type Share struct {
	Id          int
	New         int
	All         int
	AllotmentId string
}

func GetIslinkAlltomentList(page, pageSize, id int) (s []SalesmanAllotment, err error) {
	o := orm.NewOrm()
	sql := `select * from salesman_allotment where salesman_id = ? and mark <> ""`
	_, err = o.Raw(sql, id).QueryRows(&s)
	return
}
func GetNotIslinkAlltomentList(page, pageSize, id int) (s []SalesmanAllotment, err error) {
	o := orm.NewOrm()
	sql := `select * from salesman_allotment where salesman_id = ? and mark = ""`
	_, err = o.Raw(sql, id).QueryRows(&s)
	return
}
func GetIslinkAlltomentCount(id int) (count int, err error) {
	o := orm.NewOrm()
	sql := `select count(1) from salesman_allotment where salesman_id = ? and mark <> ""`
	err = o.Raw(sql, id).QueryRow(&count)
	return
}

func GetNotIslinkAlltomentCount(id int) (count int, err error) {
	o := orm.NewOrm()
	sql := `select count(1) from salesman_allotment where salesman_id = ? and mark = ""`
	err = o.Raw(sql, id).QueryRow(&count)
	return
}

//根据用户uid获取邀请码
func QueryUserInviteCodeByUid(uid int) (inviteCode string, err error) {
	o := orm.NewOrm()
	sql := `SELECT invite_code FROM salesman_invite WHERE uid = ?`
	err = o.Raw(sql, uid).QueryRow(&inviteCode)
	return
}

//查询所有业务员id
func QueryUidsBySalesmanIds(condition string, orgStr string, isData bool, para ...interface{}) (v []string, err error) {
	o := orm.NewOrm()
	sql := `SELECT 
				uid 
			FROM
				salesman_invite 
			WHERE salesman_id IN (
				SELECT s.id FROM salesman s
						LEFT JOIN sys_station sta ON s.sta_id = sta.id
						WHERE 1 = 1 `
	if condition != "" {
		sql += condition
	}
	if orgStr != "" {
		sql += ` AND s.org_id IN(` + orgStr + `)`
	}
	sql += `) AND state = 1 `
	if isData {
		sql += ` AND DATE(create_time) = CURDATE()`
	}
	_, err = o.Raw(sql, para).QueryRows(&v)
	return
}

//身份认证通过量
func QuerySalesmanOcrAuthCount(uids string) (count int, err error) {
	o := orm.NewOrm()
	sql := `SELECT COUNT(1) FROM users_auth WHERE is_valid = 1 AND is_ocr = 2 AND uid IN (` + uids + `)`
	err = o.Raw(sql).QueryRow(&count)
	return
}

//提交授信申请量
func QuerySalesmanPostAuthCount(uids string) (count int, err error) {
	o := orm.NewOrm()
	sql := `SELECT COUNT(DISTINCT uid) FROM credit_aduit WHERE uid IN (` + uids + `)`
	err = o.Raw(sql).QueryRow(&count)
	return
}

//授信通过量
func QuerySalesmanAuthPassCount(uids string) (count int, err error) {
	o := orm.NewOrm()
	sql := `SELECT COUNT(DISTINCT uid) FROM credit_aduit WHERE uid IN (` + uids + `) AND state = "PASS"`
	err = o.Raw(sql).QueryRow(&count)
	return
}

//申请借款笔数
func QuerySalesmanLoanCount(uids string) (count int, err error) {
	o := orm.NewOrm()
	sql := `SELECT COUNT(1) FROM loan WHERE uid IN (` + uids + `)`
	err = o.Raw(sql).QueryRow(&count)
	return
}

//放款成功笔数
func QuerySalesmanLoanSucceedCount(uids string) (count int, err error) {
	o := orm.NewOrm()
	sql := `SELECT COUNT(1) FROM loan WHERE uid IN (` + uids + `) AND state IN ("FINISH","BACKING")`
	err = o.Raw(sql).QueryRow(&count)
	return
}

//当日逾期金额
func QuerySalesmanLoanOverdueMoney(uids string) (money float64, err error) {
	o := orm.NewOrm()
	sql := `SELECT SUM(capital_amount+tax_amount+overdue_breach_of_amount+overdue_money_amount+data_service_fee-remain_money_charge_up_amount) FROM repayment_schedule WHERE uid IN (` + uids + `) AND state="OVERDUE"`
	err = o.Raw(sql).QueryRow(&money)
	return
}

//当日应还金额
func QuerySalesmanLoanAllMoneyByToday(uids string) (money float64, err error) {
	o := orm.NewOrm()
	sql := `SELECT SUM(capital_amount+tax_amount+overdue_breach_of_amount+overdue_money_amount+data_service_fee) FROM repayment_schedule WHERE uid IN (` + uids + `) AND DATE(loan_return_date) <= CURDATE()`
	err = o.Raw(sql).QueryRow(&money)
	return
}

package models

import (
	"fenqi_v1/utils"
	"time"
	"zcm_tools/orm"
)

//用户基本信息
type UsersBaseInfo struct {
	Uid               int       //用户ID
	Source            string    //渠道
	Verifyrealname    string    //真实姓名
	Account           string    //手机号
	IdCard            string    //身份证
	Education         string    //学历
	Marriage          string    //婚姻状态
	Province          string    //省
	City              string    //市
	District          string    //区或县
	LiveDetailAddress string    //详细地址
	Longitude         string    //经度
	Latitude          string    //纬度
	AssessTime        time.Time //信审时间
	UbiId             int       //user_base_info的id
	Sex               string    //性别
	State             int       //用户状态:0-正常 1-冻结
	Remark            string    //备注
	Balance           float64   //总额度
	UseBalance        float64   //已使用额度
	MobileTypeRecent  string    //最近操作的手机型号
}

type UsersMetadata struct {
	Id                  int       `orm:"column(id);auto"`
	Uid                 int       `orm:"column(uid);null" description:"用户id"`
	Verifyrealname      string    `orm:"column(verifyrealname);size(20);null" description:"实名"`
	IdCard              string    `orm:"column(id_card);size(20);null" description:"身份证号"`
	Verifytradepassword string    `orm:"column(verifytradepassword);size(32);null"`
	Sex                 string    `orm:"column(sex);null" description:"性别"`
	Account             string    `orm:"column(account);size(11);null" description:"手机号码"`
	IsVerifyemail       int8      `orm:"column(is_verifyemail);null" description:"邮箱是否认证"`
	IsVerifymobile      int8      `orm:"column(is_verifymobile);null" description:"手机是否认证"`
	VerifyTime          time.Time `orm:"column(verify_time);type(timestamp);null;auto_now_add" description:"实名认证时间"`
	MailAddress         string    `orm:"column(mail_address);size(200);null" description:"邮箱"`
	BirthDate           time.Time `orm:"column(birth_date);type(date);null" description:"出生日期"`
	Balance             float64   `orm:"column(balance);null;digits(10);decimals(2)" description:"可借总额"`
	UseBalance          float64   `orm:"column(use_balance);null;digits(10);decimals(2)" description:"已借金额"`
	IsntSysQuota        int8      `orm:"column(isnt_sys_quota);null" description:"0:未出系统额度 1:已出系统额度 2：人工"`
	AssessTime          time.Time `orm:"column(assess_time);type(datetime);null" description:"评估提交时间"`
	FingerKey           string    `orm:"column(finger_key);size(255);null"`
	Nation              string    `orm:"column(nation);size(20);null" description:"民族"`
	SesameCredit2       int       `orm:"column(sesame_credit2);null" description:"芝麻信用600分是否达标,0默认,1:达标,2: 不达标,3为:无法评估该用户的信用4:接口请求失败"`
	ProvinceCode        string    `orm:"column(province_code);size(5);null"`
	IsCommunication     int       `orm:"column(is_communication);null" description:"是否获取了通讯录,1.没有,2.获取"`
	InitAudit           int       `orm:"column(init_audit);null" description:"0:老用户,1新用户等待审核,2未通过强指标,3未通过弱指标,4通过,5授信未通过"`
	ZmScore             string    `orm:"column(zm_score);size(10);null" description:"芝麻信用分"`
	EqualToPetitioner   string    `orm:"column(equal_to_petitioner);size(50);null" description:"是否本人"`
	Location            string    `orm:"column(location);size(100);null"`
	Address             string    `orm:"column(address);size(255);null" description:"详细地址"`
	IpLocation          string    `orm:"column(ip_location);size(100);null"`
	IpAddress           string    `orm:"column(ip_address);size(255);null" description:"详细地址"`
	PassRealName        string    `orm:"column(pass_real_name);size(20);null" description:"是否实名认证"`
	LinkNormal          int       `orm:"column(link_normal);null" description:"紧急联系人是否异常:0:正常,1:异常"`
	Province            string    `orm:"column(province);size(30);null"`
	TdResult            string    `orm:"column(td_result);size(10);null" description:"同盾决策结果"`
	PayPwd              string    `orm:"column(pay_pwd);size(200);null" description:"支付密码"`
	AuthState           int8      `orm:"column(auth_state);null" description:"0:未授信;1:授信中;2:授信成功;3:授信资料过期;4:授信驳回;8:授信关闭30天;9:授信永久关闭;"`
	AuthStateValidTime  time.Time `orm:"column(auth_state_valid_time);type(datetime);null" description:"允许授信时间"`
	LoanState           int8      `orm:"column(loan_state);null" description:"1:允许借款8:借款关闭30天;9:借款永久关闭;"`
	LoanStateValidTime  time.Time `orm:"column(loan_state_valid_time);type(datetime);null" description:"允许借款时间"`
	OptState            int8      `orm:"column(opt_state);null" description:"1:允许8:关闭30天;9:永久关闭;"`
	OptStateValidTime   time.Time `orm:"column(opt_state_valid_time);type(datetime);null"`
	YdReport            int8      `orm:"column(yd_report);null"`
	IsTag               int
	YdAddress           string
	LiveAddress         string
	Marriage            string
	LiveDetailAddress   string
	TagType             string
	OperationTime       string
}

//有盾数据
func GetYoudunReportUsersMetadata(uid int) (v UsersMetadata, err error) {
	o := orm.NewOrm()
	err = o.Raw(`SELECT * from users_metadata where uid = ? `, uid).QueryRow(&v)
	return
}

func UpdataYoudunReportState(uid, state int) bool {
	o := orm.NewOrm()
	sql := `UPDATE users_metadata set yd_report =?  where id = ?`
	_, err := o.Raw(sql, state, uid).Exec()
	if err == nil {
		return true
	}
	return false
}

//同盾
//更新用户的状态
func UpdateTdResult(Result string, uid int) bool {
	o := orm.NewOrm()
	if res, err := o.Raw(`UPDATE users_metadata set td_result = ? where uid = ? `, Result, uid).Exec(); err == nil {
		if num, err := res.RowsAffected(); err == nil {
			return num > 0
		}
	}
	return false
}

//获取基本信息
func QueryUsersBaseInfo(uid int) (userInfo *UsersBaseInfo, err error) {
	sql := `SELECT
				u.source,
				u.state,
				u.remark,
				u.mobile_type_recent,
				um.uid,
				um.verifyrealname,
				um.account,
				um.id_card,
				um.balance,
				um.use_balance,
				um.assess_time,
				um.sex,
				ubi.education,
				ubi.marriage,
				ubi.province,
				ubi.city,
				ubi.district,
				ubi.live_detail_address,
				ubi.lng AS longitude,
				ubi.lat AS latitude,
				ubi.id AS ubi_id
			FROM users AS u
			INNER JOIN users_metadata AS um
			ON u.id = um.uid
			LEFT JOIN users_base_info AS ubi
			ON um.uid = ubi.uid
			WHERE u.id = ?
			ORDER BY ubi.create_time DESC
			LIMIT 1`
	o := orm.NewOrm()
	err = o.Raw(sql, uid).QueryRow(&userInfo)
	return
}

//更新用户基本信息
func UpdateUsersBaseInfo(ubiId int, education, marriage, province, city, district, live_detail_address string) (err error) {
	sql := `UPDATE users_base_info SET education = ?,marriage = ?,province = ?,city = ?,district = ?,live_detail_address = ? WHERE id = ?`
	o := orm.NewOrm()
	_, err = o.Raw(sql, education, marriage, province, city, district, live_detail_address, ubiId).Exec()
	return
}

//实名认证信息
type UserAuthInfo struct {
	Uid      int    //用户ID
	IdName   string //真实姓名
	IdNumber string //身份证
	Address  string //住址
	Nation   string //民族
	Gender   string //性别
}

//获取实名认证信息
func QueryRealnameAuthInfo(uid int) (userInfo UserAuthInfo, err error) {
	sql := `SELECT
				uid,
				id_name,
				id_number,
				address,
				nation,
				gender
			FROM ocr_info
			WHERE uid = ?
			ORDER BY create_time DESC
			LIMIT 1`
	o := orm.NewOrm()
	err = o.Raw(sql, uid).QueryRow(&userInfo)
	return
}

//芝麻信用
type ZMXYInfo struct {
	Id         int
	Uid        int       //用户ID
	ZmScore    int       //芝麻分
	CreateTime time.Time //时间
}

//查询芝麻信用
func QueryZmxyInfos(uid, start, pageSize int) (list []ZMXYInfo, err error) {
	sql := `SELECT id,uid,zm_score,create_time FROM zm_score_record WHERE uid = ? ORDER BY create_time DESC LIMIT ?,?`
	o := orm.NewOrm()
	o.Using("fq_log")
	_, err = o.Raw(sql, uid, start, pageSize).QueryRows(&list)
	return
}

//查询芝麻信用总数
func QueryZmxyCount(uid int) (count int, err error) {
	sql := `SELECT COUNT(1) FROM zm_score_record WHERE uid = ?`
	o := orm.NewOrm()
	o.Using("fq_log")
	err = o.Raw(sql, uid).QueryRow(&count)
	return
}

//联系人
type Linkman struct {
	Id                 int    `json:"id"`
	Uid                int    `json:"uid"`                  //用户ID
	Relation           string `json:"relation"`             //关系
	LinkmanName        string `json:"linkman_name"`         //联系人姓名
	ContactPhoneNumber string `json:"contact_phone_number"` //联系人手机号
	RiskRecordType     int
	Zone               string  //归属地
	DialCount          int     //主叫次数
	DialedCount        int     //被叫次数
	DurationTotal      float64 //联系时间（秒）
}

//查询联系人
func QueryLinkmans(uid int) (list []Linkman, err error) {
	sql := `(SELECT
				id,
				uid,
				relation,
				linkman_name,
				contact_phone_number,
				location AS zone
			FROM users_linkman
			WHERE uid = ? AND sort_no = 1
			ORDER BY id DESC LIMIT 1)
			UNION 
			(SELECT
				id,
				uid,
				relation,
				linkman_name,
				contact_phone_number,
				location AS zone
			FROM users_linkman
			WHERE uid = ? AND sort_no = 2
			ORDER BY id DESC LIMIT 1)`
	o := orm.NewOrm()
	_, err = o.Raw(sql, uid, uid).QueryRows(&list)
	return
}

//通讯录
type UserPhoneRercord struct {
	ContactName        string `description:"联系人姓名"`
	ContactPhoneNumber string `description:"联系人手机号"`
}

//查询手机通讯录
func QueryPhoneRecord(uid, start, pageSize int) (list []UserPhoneRercord, err error) {
	session := utils.GetSession()
	defer session.Close()
	smap := map[string]interface{}{}
	c := session.DB(utils.MGO_DB).C("mailList")
	if uid > 0 {
		smap["uid"] = uid
	}
	err = c.Find(&smap).Skip(start).Limit(pageSize).One(&list)
	return
}

//查询手机通讯录数量
func QueryPhoneRecordCount(uid int) (count int, err error) {
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

//获取授权提交时间
func GetAssessTime(uid int) (assess_time time.Time, err error) {
	sql := `SELECT assess_time FROM users_metadata WHERE id = ?`
	err = orm.NewOrm().Raw(sql, uid).QueryRow(&assess_time)
	return
}

//更新用户联系人
func UpdateLinkmans(linkmans []Linkman) (err error) {
	sql := `UPDATE users_linkman SET relation = ?,linkman_name = ?,contact_phone_number = ? WHERE id = ?`
	o := orm.NewOrm()
	update, err := o.Raw(sql).Prepare()
	defer update.Close()
	if err == nil {
		for _, linkman := range linkmans {
			_, err = update.Exec(linkman.Relation, linkman.LinkmanName, linkman.ContactPhoneNumber, linkman.Id)
			if err != nil {
				o.Rollback()
				return err
			}
		}
		o.Commit()
	}
	return
}

//登录记录
type LoginRecord struct {
	Uid           int       //用户ID
	CreateTime    time.Time //登录时间
	MobileType    string    //手机机型
	MobileVersion string    //系统版本
	AppVersion    string    //APP版本
	Address       string    //定位
	SilentLogin   int       //是否静默登陆 1:是
	Devicecode    string    //设备号
	IpProvince    string    //IP定位-省
	IpCity        string    //IP定位-市
	Province      string    //GPS定位-省
	City          string    //GPS定位-市
	District      string    //GPS定位-县或区
	Street        string    //GPS定位-街
	IsGPS         bool      //是否开启GPS
	Ip            string    //Ip
}

//查询登录历史
func QueryUsersLoginRecords(uid, start, pageSize int) (list []LoginRecord, err error) {
	sql := `SELECT
				uid,
				create_time,
				mobile_type,
				mobile_version,
				app_version,
				address,
				devicecode,
				silent_login,
				ip_province,
				ip_city,
				province,
				city,
				district,
				street,
				ip
			FROM login_record
			WHERE uid = ?
			ORDER BY create_time DESC
			LIMIT ?,?`
	o := orm.NewOrm()
	_, err = o.Raw(sql, uid, start, pageSize).QueryRows(&list)
	return
}

//查询登录历史总数
func QueryUsersLoginCount(uid int) (count int, err error) {
	sql := `SELECT COUNT(1) FROM login_record WHERE uid = ?`
	o := orm.NewOrm()
	err = o.Raw(sql, uid).QueryRow(&count)
	return
}

//定位信息
type Address struct {
	Uid        int       //用户ID
	CreateTime time.Time //获取时间
	Longitude  string    //经度
	Latitude   string    //纬度
	IpProvince string    //IP定位-省
	IpCity     string    //IP定位-市
	Province   string    //GPS定位-省
	City       string    //GPS定位-市
	District   string    //GPS定位-县或区
	Street     string    //GPS定位-街
	IsGPS      bool      //是否开启GPS
	Ip         string    //Ip
}

//查询定位信息
func QueryAddressInfos(uid, start, pageSize int) (list []Address, err error) {
	sql := `SELECT
				uid,
				create_time,
				longitude,
				latitude,
				ip_province,
				ip_city,
				province,
				city,
				district,
				street,
				ip
			FROM users_location
			WHERE uid = ?
			ORDER BY create_time DESC
			LIMIT ?,?`
	o := orm.NewOrm()
	o.Using("fq_log")
	_, err = o.Raw(sql, uid, start, pageSize).QueryRows(&list)
	return
}

//查询定位信息总数
func QueryAddressCount(uid int) (count int, err error) {
	sql := `SELECT COUNT(1) FROM users_location WHERE uid = ?`
	o := orm.NewOrm()
	o.Using("fq_log")
	err = o.Raw(sql, uid).QueryRow(&count)
	return
}

//人脸比对信息
type IdentifyInfo struct {
	Uid         int       //用户ID
	CreateTime  time.Time //时间
	Similarity  string    //相似度
	AuthResult  string    //认证结果
	LivingPhoto string    //活体清晰照
}

//获取用户人脸比对信息
func GetUsersIdentifyInfo(uid int) (list []IdentifyInfo, err error) {
	sql := `SELECT
				uid,
				create_time,
				living_photo,
				similarity,
				auth_result
			FROM living_info
			WHERE uid = ?
			ORDER BY create_time DESC`
	o := orm.NewOrm()
	_, err = o.Raw(sql, uid).QueryRows(&list)
	return
}

//OCR信息
type OCRInfo struct {
	Uid                 int       //用户ID
	CreateTime          time.Time //时间
	IdcardFrontPhoto    string    //身份证正面照
	IdcardBackPhoto     string    //身份证反面照
	IdcardPortraitPhoto string    //头像照
}

//获取用户OCR信息
func GetUsersOCRInfo(uid int) (list []OCRInfo, err error) {
	sql := `SELECT
				uid,
				create_time,
				idcard_front_photo,
				idcard_back_photo,
				idcard_portrait_photo
			FROM ocr_info
			WHERE uid = ?
			ORDER BY create_time DESC`
	o := orm.NewOrm()
	_, err = o.Raw(sql, uid).QueryRows(&list)
	return
}

//更新Users_metadata  信审状态
func UpdateUserMetadataState(auth_state int, uid int) error {
	sql := `update users_metadata set auth_state = ? where uid = ?`
	_, err := orm.NewOrm().Raw(sql, auth_state, uid).Exec()
	return err
}

//更新Users_metadata money初始值
func UpdateUserMetadataInitMoney(auth_state, balance, uid int) error {
	sql := `update users_metadata set auth_state = ?,balance = ?,assess_time=NOW(),loan_state= 1 where uid = ?`
	_, err := orm.NewOrm().Raw(sql, auth_state, balance, uid).Exec()
	return err
}

//更新User_metadata 授信状态 授信有效时间
func UpdatUserMetadataStatusVaildTime(auth_state, uid int) error {
	sql := `update users_metadata set auth_state = ?,auth_state_valid_time=DATE_ADD(NOW(), INTERVAL 30 DAY) where uid = ?`
	_, err := orm.NewOrm().Raw(sql, auth_state, uid).Exec()
	return err
}

//检查该人员是否有冻结账户的权限
func IsHasForzenPermission(sysId int) (isFrozen bool) {
	o := orm.NewOrm()
	sql := `SELECT COUNT(1) FROM sys_user s WHERE s.id = ? AND s.accountstatus = '启用' AND s.station_id IN (1)`
	var count int
	err := o.Raw(sql, sysId).QueryRow(&count)
	if err != nil {
		return false
	}
	if count > 0 {
		return true
	}
	return false
}

//更新用户状态（冻结和取消冻结）
func UpdateUsersFrozenState(id, state int) (err error) {
	o := orm.NewOrm()
	sql := `UPDATE users SET state = ? WHERE id = ?`
	_, err = o.Raw(sql, state, id).Exec()
	return
}

//插入冻结和取消冻结记录
func InsertUsersFrozenLog(uid int, operator, operatorType string) (err error) {
	o := orm.NewOrm()
	sql := `INSERT INTO users_frozen_log (uid,operator,operator_time,operator_type) VALUES(?,?,NOW(),?)`
	o.Using("fq_log")
	_, err = o.Raw(sql, uid, operator, operatorType).Exec()
	return
}

//更新用户备注
func UpdateUsersSign(uid int, content string) (err error) {
	o := orm.NewOrm()
	sql := `UPDATE users SET remark = ? WHERE id = ?`
	_, err = o.Raw(sql, content, uid).Exec()
	return err
}

//查询用户认证状态
func GetUserAuthState(uid int) (result int, err error) {
	sql := `SELECT auth_state FROM users_metadata WHERE uid=?`
	err = orm.NewOrm().Raw(sql, uid).QueryRow(&result)
	return
}

//通过手机号码确定是否存在用户
func GetUsersMetadataByaccount(account string) (uid int, err error) {
	sql := `SELECT uid  from users_metadata where  account = ?`
	err = orm.NewOrm().Raw(sql, account).QueryRow(&uid)
	return
}

//用户学信接口标识
type UserXxwIden struct {
	Is_xxw_auth       int    `orm:"column(is_xxw_auth);tinyint(2);null"`
	Is_yd_pscore_auth int    `orm:"column(is_yd_pscore_auth);tinyint(2);null"`
	Xxw_mgo_id        string `orm:"column(xxw_mgo_id);varchar(45);null"`
	Yd_pscore_mgo_id  string `orm:"column(yd_pscore_mgo_id);varchar(45);null"`
}

//学信接口标识
func QueryXxwIden(uid int) (v *UserXxwIden, err error) {
	o := orm.NewOrm()
	sql := `SELECT is_xxw_auth,is_yd_pscore_auth,xxw_mgo_id,yd_pscore_mgo_id FROM users_auth  WHERE is_valid = 1 AND uid= ?`
	err = o.Raw(sql, uid).QueryRow(&v)
	return
}

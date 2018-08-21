package utils

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math/rand"
	"regexp"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"
	"zcm_tools/email"

	"github.com/tealeg/xlsx"

	"github.com/astaxie/beego/httplib"
	"github.com/axgle/mahonia"
)

//随机数种子
var rnd = rand.New(rand.NewSource(time.Now().UnixNano()))

func IdUidEncrypt(id interface{}) (idStr string) {
	idS := ""
	if idString, ok := id.(string); ok {
		idS = idString
	} else if idInt64, ok := id.(int64); ok {
		idS = strconv.FormatInt(idInt64, 10)
	} else if idInt, ok := id.(int); ok {
		idS = strconv.Itoa(idInt)
	}
	idStr, err := MyDesBase64Encrypt(idS)
	// DescEnLog.Info("加密前参数===" + idS + ",加密后参数===" + idStr)
	if err != nil {
		email.Send("xjd_v1后台加密失败", err.Error(), "qxw@zcmlc.com", "huawuyou")
		return
	}
	return
}

type ParamsSina struct {
	Url_short string
	Url_long  string
	Type      int
}

func SinaDispose(urlParams string, uid, loanId, repaymentScheduleId int) (urlShort string, err error) {
	var list []ParamsSina
	ActivityUrl := ""
	params := `{"uid": ` + strconv.Itoa(uid) + `,"loanId":` + strconv.Itoa(loanId) +
		`,"repaymentScheduleId":` + strconv.Itoa(repaymentScheduleId) + `}`
	// params := "uid=" + strconv.Itoa(uid) + "&loanId=" + strconv.Itoa(loanId) +
	// "&repaymentScheduleId=" + strconv.Itoa(repaymentScheduleId)
	fmt.Println(params)
	params = string(DesBase64Encrypt([]byte(params)))
	urlParams += "?params=" + params
	fmt.Println(ActivityUrl + urlParams)
	url := "http://api.t.sina.com.cn/short_url/shorten.json?source=3271760578&url_long=" + ActivityUrl + urlParams
	req := httplib.Get(url)
	err = req.ToJSON(&list)
	if err == nil {
		urlShort = list[0].Url_short
	}
	return
}

func FloatSub(x, y float64) float64 {
	return x - y
}

func SumMoney(n int, nums ...float64) string {
	result := 0.0
	for k, v := range nums {
		if k > (n - 2) {
			result -= v
		} else {
			result += v
		}
	}
	return fmt.Sprintf("%.2f", result)
}

//float64转字符串 保留2位小数
func Float64ToString(m float64) string {
	return fmt.Sprintf("%.2f", m)
}

func Float64ToStrings(m float64) string {
	return fmt.Sprintf("%.0f", m)
}

func StringsToJSON(str string) string {
	rs := []rune(str)
	jsons := ""
	for _, r := range rs {
		rint := int(r)
		if rint < 128 {
			jsons += string(r)
		} else {
			jsons += "\\u" + strconv.FormatInt(int64(rint), 16) // json
		}
	}
	return jsons
}

//编码转换
func ConvertToString(src string, srcCode string, tagCode string) string {
	srcCoder := mahonia.NewDecoder(srcCode)
	srcResult := srcCoder.ConvertString(src)
	tagCoder := mahonia.NewDecoder(tagCode)
	_, cdata, _ := tagCoder.Translate([]byte(srcResult), true)
	result := string(cdata)
	return result
}

//数据没有记录
func ErrNoRow() string {
	return "<QuerySeter> no row found"
}

//Mongdb没有记录
func MongdbErrNoRow() string {
	return "not found"
}

//序列化
func ToString(v interface{}) string {
	data, _ := json.Marshal(v)
	return string(data)
}

//md5加密
func MD5(data string) string {
	m := md5.Sum([]byte(data))
	return hex.EncodeToString(m[:])
}

func GetToday(format string) string {
	today := time.Now().Format(format)
	return today
}

//获取今天剩余秒数
func GetTodayLastSecond() time.Duration {
	today := GetToday(FormatDate) + " 23:59:59"
	end, _ := time.ParseInLocation(FormatDateTime, today, time.Local)
	return time.Duration(end.Unix()-time.Now().Local().Unix()) * time.Second
}

//截取小数点后几位
func SubFloatToString(f float64, m int) string {
	n := strconv.FormatFloat(f, 'f', -1, 64)
	if n == "" {
		return ""
	}
	if m >= len(n) {
		return n
	}
	newn := strings.Split(n, ".")
	if m == 0 {
		return newn[0]
	}
	if len(newn) < 2 || m >= len(newn[1]) {
		return n
	}
	return newn[0] + "." + newn[1][:m]
}

//截取小数点后几位
func SubFloatToFloat(f float64, m int) float64 {
	newn := SubFloatToString(f, m)
	newf, _ := strconv.ParseFloat(newn, 64)
	return newf
}

//获取相差时间-年
func GetYearDiffer(start_time, end_time string) int64 {
	var Age int64
	t1, err := time.ParseInLocation("2006-01-02", start_time, time.Local)
	t2, err := time.ParseInLocation("2006-01-02", end_time, time.Local)
	if err == nil && t1.Before(t2) {
		diff := t2.Unix() - t1.Unix()
		Age = diff / (3600 * 365 * 24)
		return Age
	} else {
		return Age
	}
}

func FixFloat(f float64, m int) float64 {
	newn := SubFloatToString(f+0.00000001, m)
	newf, _ := strconv.ParseFloat(newn, 64)
	return newf
}

var whoareyou = make(map[string]string)

func init() {
	// Rc, Re = cache.NewCache("redis", BEEGO_CACHE)
	var yidong []string = []string{"134", "135", "136", "137", "138", "139", "147", "150", "151", "152", "157", "158", "159", "178", "182", "183", "184", "187", "188"}
	var liantong []string = []string{"130", "131", "132", "145", "155", "156", "176", "185", "186"}
	var dianxin []string = []string{"133", "153", "177", "180", "181", "189", "173"}
	for i := 0; i < len(yidong); i++ {
		whoareyou[yidong[i]] = "P100081"
	}
	for i := 0; i < len(liantong); i++ {
		whoareyou[liantong[i]] = "P100080"
	}
	for i := 0; i < len(dianxin); i++ {
		whoareyou[dianxin[i]] = "P100082"
	}
}

//验证是否是手机号
func Validate(mobileNum string) bool {
	reg := regexp.MustCompile(MobileRegular)
	return reg.MatchString(mobileNum)
}

func PageCount(count, pagesize int) int {
	if count%pagesize > 0 {
		return count/pagesize + 1
	} else {
		return count / pagesize
	}
}

func StartIndex(page, pagesize int) int {
	if page > 1 {
		return (page - 1) * pagesize
	}
	return 0
}

//获取省市县
func GetAddressName(address string) (provinces, city, counties string) {
	//直辖市,特别行政区
	cities := []string{"上海", "天津", "北京", "重庆", "香港", "澳门"}
	for _, v := range cities {
		if strings.Contains(address, v) {
			provinces = v
			city = v
			if strings.Contains(address, "市") {
				if len(address) > strings.Index(address, "市")+3 {
					address = address[strings.Index(address, "市")+3:]
					counties = address[:9]
				}
			} else if len(address) > 6 {
				address = address[6:]
				counties = address[:9]
			}
			return
		}
	}
	//获取省份
	if provincesIndex := strings.Index(address, "省"); provincesIndex != -1 { //省
		provinces = address[:provincesIndex+3]
		if len(address) > provincesIndex+3 {
			address = address[provincesIndex+3:]
		} else {
			return
		}
	}
	if provincesIndex := strings.Index(address, "自治区"); provincesIndex != -1 { //自治区
		provinces = address[:provincesIndex+9]
		if len(address) > provincesIndex+9 {
			address = address[provincesIndex+9:]
		} else {
			return
		}
	}
	//获取市
	if cityIndex := strings.Index(address, "市"); cityIndex != -1 { //市
		city = address[:cityIndex+3]
		if len(address) > cityIndex+3 {
			address = address[cityIndex+3:]
		} else {
			return
		}
	}
	if cityIndex := strings.Index(address, "盟"); cityIndex != -1 { //盟
		city = address[:cityIndex+3]
		if len(address) > cityIndex+3 {
			address = address[cityIndex+3:]
		} else {
			return
		}
	}
	if cityIndex := strings.Index(address, "自治州"); cityIndex != -1 { //自治州
		city = address[:cityIndex+9]
		if len(address) > cityIndex+9 {
			address = address[cityIndex+9:]
		} else {
			return
		}
	}
	if cityIndex := strings.Index(address, "地区"); cityIndex != -1 { //地区
		city = address[:cityIndex+6]
		if len(address) > cityIndex+6 {
			address = address[cityIndex+6:]
		} else {
			return
		}
	}
	//获取县
	if countiesIndex := strings.Index(address, "县"); countiesIndex != -1 {
		counties = address[:countiesIndex+3]
		if len(address) > countiesIndex+3 {
			address = address[countiesIndex+3:]
		}
	}
	if countiesIndex := strings.Index(address, "区"); countiesIndex != -1 {
		counties = address[:countiesIndex+3]
		if len(address) > countiesIndex+3 {
			address = address[countiesIndex+3:]
		}
	}
	if countiesIndex := strings.Index(address, "市"); countiesIndex != -1 {
		counties = address[:countiesIndex+3]
		if len(address) > countiesIndex+3 {
			address = address[countiesIndex+3:]
		}
	}
	return
}

func AddRemark(orderState string, operator string) (remark string) {
	if orderState == "OUTQUEUE" {
		remark = "订单退回———" + operator
	} else if orderState == "PASS" {
		remark = "借款订单通过———" + operator
	} else if orderState == "REJECT" {
		remark = "驳回订单———" + operator
	} else if orderState == "PAUSE" {
		remark = "关闭30天————" + operator
	} else if orderState == "CLOSE" {
		remark = "永久关闭————" + operator
	} else if orderState == "CANCEL" {
		remark = "正常关闭————" + operator
	}
	return
}

func GenerateContent(orderState string, operator string) (remark string) {
	if orderState == "OUTQUEUE" {
		remark = "订单退回———" + operator
	} else if orderState == "PASS" {
		remark = "借款订单通过———" + operator
	} else if orderState == "REJECT" {
		remark = "驳回订单———" + operator
	} else if orderState == "PAUSE" {
		remark = "关闭30天————" + operator
	} else if orderState == "CLOSE" {
		remark = "永久关闭————" + operator
	} else if orderState == "CANCEL" {
		remark = "正常关闭————" + operator
	}
	return
}

//Users_metadata 状态
func StateToUsersMetadate(orderState string) int {
	var authState int
	if orderState == "PAUSE" {
		authState = 8
	}
	if orderState == "CLOSE" {
		authState = 9
	}

	if orderState == "PASS" || orderState == "CANCEL" {
		authState = 1
	}
	return authState
}

//两int相除
func Divide(num1, num2 int) float64 {
	if num2 == 0 {
		return 0.0
	}
	a := float64(num1)
	b := float64(num2)
	return a / b
}

//两int相除
func FloatToInt(num1 float64, num2 int) float64 {
	if num2 == 0 {
		return 0.0
	}
	a := float64(num1)
	b := float64(num2)
	return a / b
}

//格式化格林尼治时间
func FormatGLNZTime(glnzTime string) string {
	formatTime, _ := time.ParseInLocation(GLNZFormatDateTime, glnzTime, time.Local)
	return formatTime.Format("2006-01-02 15:04:05")
}

func StrToTime(ti string) time.Time {
	local, _ := time.LoadLocation("Local")
	t, _ := time.ParseInLocation("2006-01-02 15:04:05", ti, local)
	return t
}

func UsersCreditMetaDataState(state string) int {
	if state == "CLOSE" {
		return 9
	}
	if state == "PAUSE" {
		return 8
	}
	if state == "PASS" {
		return 2
	}
	if state == "OUTQUEUE" {
		return 0
	}
	if state == "REJECT" {
		return 4
	}
	return 0
}

func UsersAuthState(state string) int {
	if state == "PASS" {
		return 2
	}
	if state == "CLOSE" {
		return 5
	}
	if state == "PAUSE" {
		return 4
	}
	if state == "REJECT" {
		return 6
	}
	if state == "OUTQUEUE" {
		return 1
	}
	return 0
}

//授信队列通道变量
type ShareDdata struct {
	Uid     int
	Flag    bool
	Id      int
	TimeDff float64
}

func GetIndex(index int) int {
	return index + 1
}

func ShareFlag(tmp_a ShareDdata, tmp_p ShareDdata) ShareDdata {
	var tmpdata ShareDdata
	if tmp_a.Flag {
		tmpdata = tmp_a
	}
	if tmp_p.Flag {
		tmpdata = tmp_p
	}
	return tmpdata
}

func Substr(str string, start, length int) string {
	rs := []rune(str)
	rl := len(rs)
	end := 0

	if start < 0 {
		start = rl - 1 + start
	}
	end = start + length

	if start > end {
		start, end = end, start
	}

	if start < 0 {
		start = 0
	}
	if start > rl {
		start = rl
	}
	if end < 0 {
		end = 0
	}
	if end > rl {
		end = rl
	}
	return string(rs[start:end])
}

//手机号脱敏
func MobileFilter(str string) string {
	if len(str) == 11 {
		return str[0:3] + "****" + str[7:]
	} else {
		return str
	}
}

//身份证脱敏
func IdCardFilter(str string) string {
	if len(str) == 15 || len(str) == 18 {
		return str[0:3] + "***********" + str[14:]
	} else {
		return str
	}
}

//姓名脱敏
func IDNameFiter(str string) (c string) {
	if str != "" {
		b := []rune(str)
		c = string(b[0])
		for i := 1; i < len(b); i++ {
			c = c + `*`
		}
	}
	return
}

//  方案比例随机
type RadioWeight struct {
	Rad    string
	Weight int
}

func WeightRoundRobin(ws []RadioWeight) RadioWeight {
	var randPool []RadioWeight
	var n int
	for _, v := range ws {
		if v.Weight != 0 {
			for i := 0; i < v.Weight; i++ {
				randPool = append(randPool, v)
			}
			n = n + v.Weight
		}
	}
	r := RandNumber(n)
	return randPool[r]
}

// 获取随机数
func RandNumber(n int) int {
	a := rnd.Intn(n)
	return a
}

//相除
func GetDivide(num, num1 int) string {
	if num1 == 0 {
		return "0.00%"
	}
	a := (float64)(num)
	b := (float64)(num1)
	data := fmt.Sprintf("%.2f", (a/b)*100)
	return data + "%"
}

//导出到excel下载
func ExportToExcel(data [][]string, colWidth []float64, fileName string, fontSize ...int) (filename string, err error) {
	//遍历exportToExcel
	var file *xlsx.File
	var sheet *xlsx.Sheet
	var row *xlsx.Row
	var cell *xlsx.Cell

	style := xlsx.NewStyle()
	style.Alignment.Vertical = "center"
	style.Alignment.Horizontal = "center"
	style.Font.Size = 16
	if len(fontSize) != 0 {
		style.Font.Size = fontSize[0]
	}
	file = xlsx.NewFile()
	sheet, err = file.AddSheet("Sheet1")
	if err != nil {
		fmt.Printf(err.Error())
	}
	//遍历添加数据
	for _, v := range data {
		row = sheet.AddRow()
		for k, t := range v {
			cell = row.AddCell()
			cell.SetStyle(style)
			sheet.Cols[k].Width = colWidth[k]
			cell.Value = t
		}
	}
	filename = fileName + ".xlsx"
	err = file.Save(filename)
	return filename, err
}

//流程转化率导出EXCEL
func ProcessConversionToExcel(data [][]string, colWidth []float64, fileName string) (filename string, err error) {
	var file *xlsx.File
	var sheet *xlsx.Sheet
	var row *xlsx.Row
	var cell *xlsx.Cell
	style := xlsx.NewStyle()
	style.Alignment.Vertical = "center"
	style.Alignment.Horizontal = "center"
	style.Font.Size = 16
	style.Font.Name = "Verdana"
	file = xlsx.NewFile()
	sheet, err = file.AddSheet("Sheet1")
	if err != nil {
		fmt.Printf(err.Error())
	}
	r0 := sheet.AddRow()
	for i := 0; i < 16; i++ {
		c0 := r0.AddCell()
		c0.SetStyle(style)
		sheet.Cols[i].Width = colWidth[i]
		switch i {
		case 0:
			c0.Value = "日期"
			c0.Merge(0, 1)
		case 1:
			c0.Value = "注册"
			c0.Merge(0, 1)
		case 2:
			c0.Value = "认证申请"
			c0.Merge(0, 1)
		case 3:
			c0.Value = "认证通过"
			c0.Merge(0, 1)
		case 4:
			c0.Value = "活体通过"
			c0.Merge(0, 1)
		case 5:
			c0.Value = "个人信息补充"
			c0.Merge(0, 1)
		case 6:
			c0.Value = "基础授信"
			c0.Merge(3, 0)
		case 10:
			c0.Value = "完整授信"
			c0.Merge(1, 0)
		case 12:
			c0.Value = "授信申请"
			c0.Merge(0, 1)
		case 13:
			c0.Value = "授信通过"
			c0.Merge(0, 1)
		case 14:
			c0.Value = "首次借款申请"
			c0.Merge(0, 1)
		case 15:
			c0.Value = "首次借款通过"
			c0.Merge(0, 1)
		default:
			c0.Value = ""
		}
	}

	//遍历添加数据
	for _, v := range data {
		row = sheet.AddRow()
		for k, t := range v {
			cell = row.AddCell()
			cell.SetStyle(style)
			sheet.Cols[k].Width = colWidth[k]
			cell.Value = t
		}
	}
	filename = fileName + ".xlsx"
	err = file.Save(filename)
	return filename, err
}

/*
*辗转相除法：最大公约数
*非递归写法
 */
func MaxCommonDivisor(x, y int) int {
	var tmp int
	for {
		tmp = x % y
		if tmp > 0 {
			x = y
			y = tmp
		} else {
			return y
		}
	}
}

//遍历是否存在数据库中
func CheckIsSysId(loginSysId int, needids string) (flag bool) {
	str := strings.Split(needids, ",")
	sid := strconv.Itoa(loginSysId)
	for i := 0; i < len(str); i++ {
		if sid == str[i] {
			flag = true
			break
		}
	}
	return
}

//过滤emoji
func FilterEmoji(content string) string {
	new_content := ""
	for _, value := range content {
		_, size := utf8.DecodeRuneInString(string(value))
		if size <= 3 {
			new_content += string(value)
		}
	}
	return new_content
}

// 数组去重
func RemoveDuplicatesAndEmpty(a []string) (ret []string) {
	a_len := len(a)
	for i := 0; i < a_len; i++ {
		if (i > 0 && a[i-1] == a[i]) || len(a[i]) == 0 {
			continue
		}
		ret = append(ret, a[i])
	}
	return
}

func Mtype(mtype int) (ctype string) {
	if mtype == 11 {
		ctype = "M0"
	} else if mtype == 12 {
		ctype = "S1"
	} else if mtype == 13 {
		ctype = "S2"
	} else if mtype == 14 {
		ctype = "S3"
	} else if mtype == 15 {
		ctype = "M2"
	} else if mtype == 16 {
		ctype = "M3"
	} else {
		ctype = ""
	}
	return
}

func OverudDays(day int) (ctype string) {
	if day == 0 {
		ctype = "M0"
	} else if day >= 1 && day <= 3 {
		ctype = "S1"
	} else if day >= 4 && day <= 15 {
		ctype = "S2"
	} else if day >= 16 && day <= 30 {
		ctype = "S3"
	} else if day >= 31 && day <= 60 {
		ctype = "M2"
	} else if day >= 61 && day <= 90 {
		ctype = "M3"
	} else {
		ctype = ""
	}
	return
}

//眼见的长度，中文:2 /英文数字:1 长度
func LenOfSee(str string) int {
	rs := []rune(str)
	return (len(rs) + len(str)) / 2
}

//获取i长度的空格串
func GetBlank(i int) string {
	blank := ""
	if i == 0 {
		return blank
	}
	shi := i / 10
	ge := i % 10
	if shi > 0 {
		for num := 0; num < shi; num++ {
			blank += "          "
		}
	}
	if ge > 0 {
		for num := 0; num < ge; num++ {
			blank += " "
		}
	}

	return blank
}
func Strlen(s string) int {
	rs := []rune(s)
	rl := len(rs)
	return rl
}

func GetSeriesMonth(start, end time.Time) (monthSlice []string) {
	monthSlice = make([]string, 0)
	monthSlice = append(monthSlice, Substr(start.Format("2006-01-02"), 8, 2)+"("+GetWeekDay(start.Weekday())+")")
	for start.Before(end) {
		start = start.AddDate(0, 0, 1)
		addDay := start.Format("2006-01-02")
		monthSlice = append(monthSlice, Substr(addDay, 8, 2)+"("+GetWeekDay(start.Weekday())+")")
	}
	return
}
func GetWeekDay(weekDay time.Weekday) string {
	switch weekDay {
	case time.Sunday:
		return Sunday
	case time.Monday:
		return Monday
	case time.Tuesday:
		return Tuesday
	case time.Wednesday:
		return Wednesday
	case time.Thursday:
		return Thursday
	case time.Friday:
		return Friday
	case time.Saturday:
		return Saturday
	default:
		return ""
	}
}

func IdCradDispose(idCard string) string {
	var err error
	var reg *regexp.Regexp
	var placeStr1, placeStr2 string
	if len(idCard) == 15 {
		reg, err = regexp.Compile("^(\\d{4})(\\d{2})(\\d{2})(\\d{5})(.*)")
		placeStr1 = "**"
		placeStr2 = "*****"
	} else if len(idCard) == 18 {
		reg, err = regexp.Compile("^(\\d{4})(\\d{2})(\\d{4})(\\d{6})(.*)")
		// reg, err = regexp.Compile("^(\\d{4})(\\d{12})(.*)")
		placeStr1 = "**"
		placeStr2 = "******"
	} else {
		return ""
	}
	if err != nil {
		return ""
	}
	if reg.MatchString(idCard) == true {
		submatch := reg.FindStringSubmatch(idCard)
		return submatch[1] + placeStr1 + submatch[3] + placeStr2 + submatch[5]
	}
	return ""
}

//两int相除
func FloatToFloat(num1 float64, num2 float64) float64 {
	if num2 == 0.0 {
		return 0.0
	}
	return num1 / num2
}

func CalculatePercent(a, b float64) (per string) {
	if b == 0 { //分母不能为0
		return per
	} else {
		p := Float64ToString((a * 100) / b)
		per = p + "%"
		return per
	}
}
func CalculateIntPercent(a, b int) (per string) {
	if b == 0 { //分母不能为0
		return per
	} else {
		p := strconv.Itoa((a * 100) / b)
		per = p + "%"
		return per
	}
}

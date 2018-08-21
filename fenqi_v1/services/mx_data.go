package services

import (
	"fenqi_v1/models"
	"fenqi_v1/utils"
	"sort"
	"strconv"
	"time"

	"sync"

	"github.com/astaxie/beego"
)

//公积金账单信息排序（从大到小）
func GJJQuickSort(list []models.BillRecord) {
	if len(list) <= 1 {
		return //递归终止条件，slice变为0为止。
	}
	mid := list[0]
	i := 1 //list[0]为中间值mid，所以要从1开始比较
	head, tail := 0, len(list)-1
	for head < tail {
		if list[i].Deal_Time < mid.Deal_Time {
			list[i], list[tail] = list[tail], list[i] //交换
			tail--
		} else {
			list[i], list[head] = list[head], list[i]
			head++
			i++
		}
	}
	list[head] = mid
	GJJQuickSort(list[:head])
	GJJQuickSort(list[head+1:])
}

//公积金账单信息排序（从大到小）
func AdviseSort(list []models.CreditAdvise) {
	if len(list) <= 1 {
		return //递归终止条件，slice变为0为止。
	}
	mid := list[0]
	i := 1 //list[0]为中间值mid，所以要从1开始比较
	head, tail := 0, len(list)-1
	for head < tail {
		if list[i].HandlingTime.Before(mid.HandlingTime) {
			list[i], list[tail] = list[tail], list[i] //交换
			tail--
		} else {
			list[i], list[head] = list[head], list[i]
			head++
			i++
		}
	}
	list[head] = mid
	AdviseSort(list[:head])
	AdviseSort(list[head+1:])
}

//获取和联系人近一个月的通话情况
func YYSMonthCallData(calls []models.Calls, beginTime time.Time, phone string, link *models.Linkman, finish *sync.WaitGroup) {
	before30 := beginTime.AddDate(0, 0, -29) //前30天
	interval := int(beginTime.Month() - before30.Month())
	var beginMonthItem []models.Items   //本月通话记录
	var before1MonthItem []models.Items //上个月通话记录
	var before2MonthItem []models.Items //上上个月通话记录
	if interval == 0 {                  //取1个月（本月）
		for _, v := range calls {
			if v.BillMonth == beginTime.Format(utils.FormatMonth) {
				beginMonthItem = v.Items
			}
		}
	} else if interval == 1 || interval < 0 { //取2个月（本月、上月）
		for _, v := range calls {
			if v.BillMonth == beginTime.Format(utils.FormatMonth) {
				beginMonthItem = v.Items
			} else if v.BillMonth == beginTime.AddDate(0, -1, 0).Format(utils.FormatMonth) {
				before1MonthItem = v.Items
			}
		}
	} else if interval == 2 { //取3个月（本月、上月、上上月）
		for _, v := range calls {
			if v.BillMonth == beginTime.Format(utils.FormatMonth) {
				beginMonthItem = v.Items
			} else if v.BillMonth == beginTime.AddDate(0, -1, 0).Format(utils.FormatMonth) {
				before1MonthItem = v.Items
			} else if v.BillMonth == beginTime.AddDate(0, -2, 0).Format(utils.FormatMonth) {
				before2MonthItem = v.Items
			}
		}
	}
	var dialCount int         //主叫次数
	var dialedCount int       //被叫次数
	var durationTotal float64 //联系时间
	for _, v := range beginMonthItem {
		if v.PeerNumber == phone && v.Time >= before30.Format(utils.FormatDate) && v.Time <= beginTime.Format(utils.FormatDate)+" 23:59:59" {
			if v.DialType == "DIAL" {
				dialCount++
			} else if v.DialType == "DIALED" {
				dialedCount++
			}
			durationTotal += v.Duration
		}
	}
	for _, v := range before1MonthItem {
		if v.PeerNumber == phone && v.Time >= before30.Format(utils.FormatDate) && v.Time <= beginTime.Format(utils.FormatDate)+" 23:59:59" {
			if v.DialType == "DIAL" {
				dialCount++
			} else if v.DialType == "DIALED" {
				dialedCount++
			}
			durationTotal += v.Duration
		}
	}
	for _, v := range before2MonthItem {
		if v.PeerNumber == phone && v.Time >= before30.Format(utils.FormatDate) && v.Time <= beginTime.Format(utils.FormatDate)+" 23:59:59" {
			if v.DialType == "DIAL" {
				dialCount++
			} else if v.DialType == "DIALED" {
				dialedCount++
			}
			durationTotal += v.Duration
		}
	}
	link.DialCount, link.DialedCount, link.DurationTotal = dialCount, dialedCount, durationTotal
	finish.Done()
}

// 通讯录排序
func SortAddressBook(one_month []models.CommonlyConnectMobiles) (one_month_mn []models.CommonlyConnectMobiles_2) {
	for _, v := range one_month {
		var a int
		var err error
		if v.Connect_count == "" {
			a = 0
		} else {
			a, err = strconv.Atoi(v.Connect_count)
			if err != nil {
				beego.Debug(err.Error())
				return
			}
		}
		temp := models.CommonlyConnectMobiles_2{}
		temp.Connect_count = a
		temp.Mobile = v.Mobile
		temp.Belong_to = v.Belong_to
		temp.Connect_time = v.Connect_time
		temp.Originating_call_count = v.Originating_call_count
		temp.Terminating_call_count = v.Terminating_call_count
		temp.Mon_type = v.Mon_type
		temp.Uid = v.Uid
		temp.Id = v.Id
		one_month_mn = append(one_month_mn, temp)
	}
	sort.Sort(SortMonth(one_month_mn))
	return
}

// 通话记录
type SortMonth []models.CommonlyConnectMobiles_2

func (a SortMonth) Len() int      { return len(a) }
func (a SortMonth) Swap(i, j int) { a[i], a[j] = a[j], a[i] }

func (a SortMonth) Less(i, j int) bool {
	return a[i].Connect_count > a[j].Connect_count
}

//新手机通讯录排序
func SortAddBook(mail []models.MailList) (mail2 []models.MailList2) {
	for _, v := range mail {
		//var a int
		//var err error
		/*if len(v.Contact) == 0 {
			a = 0
		} else {
			a = len(v.Contact)
		}*/
		temp := models.MailList2{}
		temp.Id = v.Id
		temp.CreateTime = v.Id.Time()
		temp.Contact = v.Contact
		mail2 = append(mail2, temp)
	}
	sort.Sort(SortMail(mail2))
	return
}

type SortMail []models.MailList2

func (a SortMail) Len() int      { return len(a) }
func (a SortMail) Swap(i, j int) { a[i], a[j] = a[j], a[i] }

func (a SortMail) Less(i, j int) bool {
	return a[i].CreateTime.After(a[j].CreateTime)
}

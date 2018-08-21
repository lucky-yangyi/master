package services

import (
	"fenqi_v1/models"
	"sort"
)

//魔蝎根据通话次数排序
type SortMXCnt []models.Call_Contact_Detail

func (a SortMXCnt) Len() int      { return len(a) }
func (a SortMXCnt) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a SortMXCnt) Less(i, j int) bool {
	return a[i].Total_Cnt > a[j].Total_Cnt
}

//魔蝎根据通话次数排序
func SortMXCntTime(mx []models.Call_Contact_Detail) (sort_mx []models.Call_Contact_Detail) {
	for _, v := range mx {
		temp := models.Call_Contact_Detail{}
		temp.Peer_num = v.Peer_num
		temp.City = v.City
		temp.Total_Cnt = v.Total_Cnt
		temp.Dial_Cnt_3m = v.Dial_Cnt_3m
		temp.Dialed_Cnt_3m = v.Dialed_Cnt_3m
		temp.Call_Time_3m = v.Call_Time_3m
		sort_mx = append(sort_mx, temp)
	}
	sort.Sort(SortMXCnt(sort_mx))
	return
}

//天机根据通话次数排序
type SortTianjiTalkCnt []models.CallLogInfo

func (a SortTianjiTalkCnt) Len() int      { return len(a) }
func (a SortTianjiTalkCnt) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a SortTianjiTalkCnt) Less(i, j int) bool {
	return a[i].Talk_cnt > a[j].Talk_cnt
}

// 天机根据通话次数排序
func SortTradeTime(tianji []models.CallLogInfo) (sort_tianji []models.CallLogInfo) {
	for _, v := range tianji {
		temp := models.CallLogInfo{}
		temp.Phone = v.Phone
		temp.Phone_location = v.Phone_location
		temp.Talk_cnt = v.Talk_cnt
		temp.Call_cnt = v.Call_cnt
		temp.Called_cnt = v.Called_cnt
		temp.Talk_seconds = v.Talk_seconds
		sort_tianji = append(sort_tianji, temp)
	}
	sort.Sort(SortTianjiTalkCnt(sort_tianji))
	return
}

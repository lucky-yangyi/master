package services

import (
	"fenqi_v1/cache"
	"fenqi_v1/models"
	"fenqi_v1/utils"
	//	"fenqi_v1/controllers"
	"strconv"
	"time"

	"github.com/astaxie/beego"
)

//校验有没有正在处理的状态
func CheckHanding(id int, displayName string) (timeDiff float64, loanId, uid int, err error) {
	uid = cache.GetCacheHandingUids(id)
	loanId = cache.GetCacheOrderMessage(uid)
	if uid > 0 && loanId > 0 {
		//count, err := models.QueryHandingCount(loanId)
		//if err != nil {
		//	return timeDiff, loanId, uid, err
		//}
		oc, err := models.QueryOrderByUid(loanId) //订单
		timeDiff = oc.Atime.Add(45 * time.Minute).Sub(time.Now()).Seconds()
		if timeDiff <= 0 {
			content := "【" + displayName + "】" + "订单处理超时"
			err = models.AddLoanAduitRecord(uid, loanId, content)
			if err != nil {
				beego.Info(err)
			}
			err = models.UpdateQueueing(loanId) //delete redis(处理后清楚缓存)
			if err != nil {
				beego.Info(err)
			}
			err = cache.DeleteOrderRedis(uid, id)
			if err != nil {
				beego.Info(err)
			}
			return 0, 0, 0, err
		} else {
			oc, err := models.QueryOrderByUid(loanId) //订单
			timeDiff = oc.Atime.Add(45 * time.Minute).Sub(time.Now()).Seconds()
			if err != nil {
				return timeDiff, loanId, uid, err
			}
		}
	}
	return
}

//排队逻辑
func Line(num int) (oc models.OrderCredit, err error) {
	// var oc models.OrderCredit
	num += 1
	if num > 10 {
		return oc, nil
	}
	oc, err = models.GetQueueing()
	if err != nil && err.Error() == utils.ErrNoRow() {
		return
	}
	//是否在处理
	loanId := cache.GetCacheOrderMessage(oc.Uid)
	if loanId != oc.Id {
		return oc, err
	} else {
		oc, err = Line(num)
	}
	return oc, err
}

//预约逻辑
func Order(num int) (oc models.OrderCredit, err error) {
	num += 1
	if num > 10 {
		return oc, nil
	}
	oc, err = models.GetOrdering()
	if err != nil && err.Error() == utils.ErrNoRow() {
		oc, err = Line(num)
		return
	}
	//预约订单不为空
	if oc.Id >= 0 {
		//入列状态
		if oc.InqueueType == 1 {
			//有无在处理
			loanId := cache.GetCacheOrderMessage(oc.Uid)
			if loanId != oc.Id {
				//err := models.UpdateOrderState(loanId, "HANDING")
				return oc, err
			} else {
				oc, err = Order(num)
			}
		}
		if oc.InqueueType == 2 {
			err = models.UpdateQueueTime(oc.Id, oc.InqueueTime)
			if err != nil {
				return oc, err
			} else {
				oc, err = Order(num)
			}
		}
	}
	return
}

func InsertOrderCache(id, loanId, uid int, operator, orderState string) (err error) {
	err = models.UpdateOrderState(loanId, id, operator, orderState)
	if err != nil {
		return
	}
	//insert redis
	if utils.Re == nil {
		err = utils.Rc.Put("xjfq:"+utils.CacheKeyHandingUids+"_"+strconv.Itoa(id), uid, utils.RedisCacheTime_Year)
		if err != nil {
			return err
		}
		err = utils.Rc.Put("xjfq:"+utils.CacheKeyOrderMessages+"_"+strconv.Itoa(uid), loanId, utils.RedisCacheTime_Year)
		if err != nil {
			return err
		}
	}
	return err
}

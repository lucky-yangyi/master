package cache

import (
	"fenqi_v1/utils"
	"strconv"
)

//加缓存锁(存在在处理状态)
func GetCacheOrderMessage(uid int) int {
	if utils.Re == nil && utils.Rc.IsExist("xjfq:"+utils.CacheKeyOrderMessages+"_"+strconv.Itoa(uid)) {
		loanId, _ := utils.Rc.RedisInt("xjfq:" + utils.CacheKeyOrderMessages + "_" + strconv.Itoa(uid))
		if loanId > 0 {
			return loanId
		}
	}
	return 0
}

//系统用户在处理的订单
func GetCacheHandingUids(id int) (uid int) {
	if utils.Re == nil && utils.Rc.IsExist("xjfq:"+utils.CacheKeyHandingUids+"_"+strconv.Itoa(id)) {
		uid, _ := utils.Rc.RedisInt("xjfq:" + utils.CacheKeyHandingUids + "_" + strconv.Itoa(id))
		if uid > 0 {
			return uid
		}
	}
	return
}

//授信存锁
func GetCacheCreditMessage(id int) (int, bool) {
	if utils.Re == nil && utils.Rc.IsExist("xjfq:"+utils.CacheKeyCreditMessage+"_"+strconv.Itoa(id)) {
		uid, _ := utils.Rc.RedisInt("xjfq:" + utils.CacheKeyCreditMessage + "_" + strconv.Itoa(id))
		if uid > 0 {
			return uid, false
		}
	}
	return 0, true
}

//系统用户在处理的授信数据
func GetCacheCreditHandingUids(tid int) (id int) {
	if utils.Re == nil && utils.Rc.IsExist("xjfq:"+utils.CacheKeyCreditHandingUids+"_"+strconv.Itoa(tid)) {
		id, _ := utils.Rc.RedisInt("xjfq:" + utils.CacheKeyCreditHandingUids + "_" + strconv.Itoa(tid))
		if id > 0 {
			return id
		}

	}
	return
}

//授信存锁
func GetAuthCacheCreditMessage(id int) (int, bool) {
	if utils.Re == nil && utils.Rc.IsExist("auth:"+utils.CacheKeyCreditMessage+"_"+strconv.Itoa(id)) {
		uid, _ := utils.Rc.RedisInt("auth:" + utils.CacheKeyCreditMessage + "_" + strconv.Itoa(id))
		if uid > 0 {
			return uid, false
		}
	}
	return 0, true
}

//系统用户在处理的授信数据
func GetAuthCacheCreditHandingUids(tid int) (id int) {
	if utils.Re == nil && utils.Rc.IsExist("auth:"+utils.CacheKeyCreditHandingUids+"_"+strconv.Itoa(tid)) {
		id, _ := utils.Rc.RedisInt("auth:" + utils.CacheKeyCreditHandingUids + "_" + strconv.Itoa(tid))
		if id > 0 {
			return id
		}

	}
	return
}

//授信存锁
func GetLinkManCacheCreditMessage(id int) (int, bool) {
	if utils.Re == nil && utils.Rc.IsExist("linkman:"+utils.CacheKeyCreditMessage+"_"+strconv.Itoa(id)) {
		uid, _ := utils.Rc.RedisInt("linkman:" + utils.CacheKeyCreditMessage + "_" + strconv.Itoa(id))
		if uid > 0 {
			return uid, false
		}
	}
	return 0, true
}

//系统用户在处理的授信数据
func GetLinkManCacheCreditHandingUids(tid int) (id int) {
	if utils.Re == nil && utils.Rc.IsExist("linkman:"+utils.CacheKeyCreditHandingUids+"_"+strconv.Itoa(tid)) {
		id, _ := utils.Rc.RedisInt("linkman:" + utils.CacheKeyCreditHandingUids + "_" + strconv.Itoa(tid))
		if id > 0 {
			return id
		}

	}
	return
}

//授信存锁
func GetMyInfoCacheCreditMessage(id int) (int, bool) {
	if utils.Re == nil && utils.Rc.IsExist("myinfo:"+utils.CacheKeyCreditMessage+"_"+strconv.Itoa(id)) {
		uid, _ := utils.Rc.RedisInt("myinfo:" + utils.CacheKeyCreditMessage + "_" + strconv.Itoa(id))
		if uid > 0 {
			return uid, false
		}
	}
	return 0, true
}

//系统用户在处理的授信数据
func GetMyInfoCacheCreditHandingUids(tid int) (id int) {
	if utils.Re == nil && utils.Rc.IsExist("myinfo:"+utils.CacheKeyCreditHandingUids+"_"+strconv.Itoa(tid)) {
		id, _ := utils.Rc.RedisInt("myinfo:" + utils.CacheKeyCreditHandingUids + "_" + strconv.Itoa(tid))
		if id > 0 {
			return id
		}

	}
	return
}

//授信存锁
func GetOtherCacheCreditMessage(id int) (int, bool) {
	if utils.Re == nil && utils.Rc.IsExist("other:"+utils.CacheKeyCreditMessage+"_"+strconv.Itoa(id)) {
		uid, _ := utils.Rc.RedisInt("other:" + utils.CacheKeyCreditMessage + "_" + strconv.Itoa(id))
		if uid > 0 {
			return uid, false
		}
	}
	return 0, true
}

//系统用户在处理的授信数据
func GetOtherCacheCreditHandingUids(tid int) (id int) {
	if utils.Re == nil && utils.Rc.IsExist("other:"+utils.CacheKeyCreditHandingUids+"_"+strconv.Itoa(tid)) {
		id, _ := utils.Rc.RedisInt("other:" + utils.CacheKeyCreditHandingUids + "_" + strconv.Itoa(tid))
		if id > 0 {
			return id
		}

	}
	return
}

//delete redis(处理后清楚缓存)
func DeleteOrderRedis(uid, id int) error {
	err := utils.Rc.Delete("xjfq:" + utils.CacheKeyOrderMessages + "_" + strconv.Itoa(uid))
	if err != nil {
		return err
	}
	err = utils.Rc.Delete("xjfq:" + utils.CacheKeyHandingUids + "_" + strconv.Itoa(id))
	if err != nil {
		return err
	}
	return err
}

//频繁提交
func GetCacheLoanOften(loanId int) bool {
	if utils.Re == nil && utils.Rc.IsExist(utils.CacheKeyLoanOften+"_"+strconv.Itoa(loanId)) {
		return true
	}
	return false
}

//重检频繁提交
func GetCacheRecheckOften(uid int) bool {
	if utils.Re == nil && utils.Rc.IsExist(utils.CacheKeyRecheck+"_"+strconv.Itoa(uid)) {
		return true
	}
	return false
}

//重检频繁提交
func GetCacheLOCKRecheck(uid int, loanId int) bool {
	if utils.Re == nil && utils.Rc.IsExist(utils.CACHE_KEY_Anewcheck_Request+"_"+strconv.Itoa(uid)+strconv.Itoa(loanId)) {
		return true
	}
	return false
}

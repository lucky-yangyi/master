package cache

import (
	"fenqi_v1/models"
	"fenqi_v1/utils"
	"github.com/astaxie/beego/context"
	"net/url"
	"time"
)

func RecordLogs(user_id, business_id int, username, displayname, action, logger, message string, input *context.BeegoInput) bool {
	ip := input.IP()
	urlpath := input.URL()
	querystrings := input.URI()
	fromparams, _ := url.QueryUnescape(string(input.RequestBody))
	log := &models.SysLog{UserId: user_id, UserName: username, UserDisplayName: displayname, UserIp: ip, Action: action, Logger: logger, UrlPath: urlpath, Message: message, FromParams: fromparams, QueryStrings: querystrings, CreateTime: time.Now(), BusinessId: business_id}
	if utils.Re == nil {
		utils.Rc.LPush(utils.CacheKeySystemLogs, log)
		return true
	}
	return false
}

func RecordOperateTime(user_id int, operateTime time.Time) bool {
	log := &models.OperateTimeLogs{UserId: user_id, OperationTime: operateTime}
	if utils.Re == nil {
		utils.Rc.LPush(utils.CACHE_KEY_OPERATETIME_LOGS, log)
		return true
	}
	return false
}

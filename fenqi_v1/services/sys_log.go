package services

import (
	"encoding/json"
	"fenqi_v1/models"
	"fenqi_v1/utils"
	"fmt"
)

//插入日志到数据库
func AutoInsertLogToDB() {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("[AutoInsertLogToDB]", err)
		}
	}()
	for {
		utils.Rc.Brpop(utils.CacheKeySystemLogs, func(b []byte) {
			var log models.SysLog
			if err := json.Unmarshal(b, &log); err != nil {
				fmt.Println("json unmarshal wrong!")
			}
			if _, err := models.AddLogs(&log); err != nil {
				fmt.Println(err.Error(), log)
			}
		})
	}
}

//更新用户操作时间到数据库
func AutoUpdateUserOperateTime() {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("[AutoUpdateUserOperateTime]", err)
		}
	}()
	for {
		utils.Rc.Brpop(utils.CACHE_KEY_OPERATETIME_LOGS, func(b []byte) {
			var log models.OperateTimeLogs
			if err := json.Unmarshal(b, &log); err != nil {
				fmt.Println("json unmarshal wrong!")
			}
			if err := models.UpdateLastOperationTime(log.UserId, log.OperationTime); err != nil {
				fmt.Println(err.Error(), log)
			}
		})
	}
}

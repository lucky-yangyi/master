package models

import (
	"fenqi_v1/utils"
	"time"
)

type UserMessage struct {
	Title          string
	Content        string
	Uid            int
	CreateTime     string
	MsgType        int8
	IsRead         int8
	CreateHideTime string
}

func AddUsersMessage(um *UserMessage) bool {
	session := utils.GetSession()
	defer session.Close()
	umsg := session.DB(utils.MGO_DB).C("users_message")
	um.CreateTime = time.Now().Format(utils.FormatDate)
	um.CreateHideTime = time.Now().Format(utils.FormatDateTime)
	err := umsg.Insert(um)
	if err != nil {
		return false
	}
	return true
}

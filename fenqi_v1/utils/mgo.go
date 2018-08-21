package utils

import (
	"fmt"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

/**
 * 公共方法，获取session，如果存在则拷贝一份
 */

var MgoSession, MgoSessionYeadun *mgo.Session

func GetSession() *mgo.Session {
	if MgoSession == nil {
		var err error
		MgoSession, err = mgo.Dial(MGO_URL)
		if err != nil {
			panic(err) //直接终止程序运行
		}
	} else {
		err := MgoSession.Ping()
		if err != nil {
			MgoSession, err = mgo.Dial(MGO_URL)
			if err != nil {
				fmt.Println(err)
			}
		}
	}
	//最大连接池默认为4096
	return MgoSession.Clone()
}

func GetYeadunSession() *mgo.Session {
	if MgoSessionYeadun == nil {
		var err error
		MgoSessionYeadun, err = mgo.Dial(MGO_YEADUN_URL)
		if err != nil {
			panic(err) //直接终止程序运行
		}
	} else {
		err := MgoSessionYeadun.Ping()
		if err != nil {
			MgoSessionYeadun, err = mgo.Dial(MGO_YEADUN_URL)
			if err != nil {
				fmt.Println(err)
			}
		}
	}
	//最大连接池默认为4096
	return MgoSessionYeadun.Clone()
}

//验证id是否符合mongodbId类型
func CheckIsBsonId(_id string) bool {
	return bson.IsObjectIdHex(_id)
}

//公共方法，获取collection对象
func WitchCollection(collection string, s func(*mgo.Collection) error) error {
	session := GetSession()
	defer session.Close()
	c := session.DB(MGO_DB).C(collection)
	return s(c)
}

func init() {
	//MgoSession = GetSession()
	//MgoSessionYeadun = GetYeadunSession()
}

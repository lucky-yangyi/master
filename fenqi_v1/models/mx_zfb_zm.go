package models

import "time"

type MdbMxZFBZMData struct {
	Uid         int
	MXZFBZMData []MXZFBZMData `bson:"data"`
	CreateTime  time.Time
}

type MXZFBZMData struct {
	Time    string  `bson:"time"`
	ZmScore float64 `bson:"zm_score"`
}

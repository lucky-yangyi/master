package models

//用户消息
type UsersMessage struct {
	Uid            int    //用户id
	Title          string //标题
	Content        string //内容
	Isread         int    //是否已读
	Createhidetime string //创建时间
}

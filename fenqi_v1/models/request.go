package models

type Loading struct {
	LoanIds     []int  `json:"loan_ids"`     //uid数组
	InqueueTime string `json:"inqueue_time"` //入列时间
	InqueueType int    `json:"inqueue_type"` //入列类型
}

//api请求体
type LoanApprovalRequest struct {
	Ret int    `description:"状态码"`
	Msg string `description:"返回消息"`
	Err string `description:"错误信息"`
}
type CreditQueueLine struct {
	Ids         []int  `json:"Ids"`         //id数组
	InqueueTime string `json:"inqueueTime"` //入列时间
	InqueueType int    `json:"inqueueType"` //入列类型
}

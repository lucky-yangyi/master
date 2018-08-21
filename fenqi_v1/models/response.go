package models

type BaseResponse struct {
	Ret int    `json:"ret"`
	Err string `json:"err"`
	Msg string `json:"msg"`
}

type SalesmanResponse struct {
	CheckIds []int `json:"checkedId"`
	Num      int   `json:"num"`
}

package dto

type Base struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}
type Response struct {
	Base Base        `json:"base"`
	Data interface{} `json:"data"`
}

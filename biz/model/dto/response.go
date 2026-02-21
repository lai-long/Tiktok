package dto

type Base struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}
type Response struct {
	Base Base        `json:"base"`
	Data interface{} `json:"data"`
}
type LoginResponse struct {
	Response
	RefreshToken string `json:"refresh_token"`
	AccessToken  string `json:"access_token"`
}
type Items struct {
	Video   []Video   `json:"video"`
	Comment []Comment `json:"comment"`
	Total   int       `json:"total"`
}

package dto

type Base struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}
type Response struct {
	Base         Base        `json:"base"`
	Data         interface{} `json:"data"`
	AccessToken  string      `json:"access_token"`
	RefreshToken string      `json:"refresh_token"`
}
type LoginResponse struct {
	Response
	RefreshToken string `json:"refresh_token"`
	AccessToken  string `json:"access_token"`
}

type Data struct {
	Items interface{} `json:"items"`
	Total int         `json:"total"`
}

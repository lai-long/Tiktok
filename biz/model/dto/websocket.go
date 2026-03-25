package dto

type SendMsg struct {
	Type    string `json:"type"`
	Content string `json:"content"`
	GroupId string `json:"group_id"`
}
type ReplyMsg struct {
	From    string `json:"from"`
	Code    int    `json:"code"`
	Content string `json:"content"`
}

type Message struct {
	Sender    string `json:"sender"`
	Recipient string `json:"recipient"`
	Content   string `json:"content"`
}

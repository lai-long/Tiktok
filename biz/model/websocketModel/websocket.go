package websocketModel

type SendMsg struct {
	Type    string `json:"type"`
	Content string `json:"content"`
}
type ReplyMsg struct {
	From    string `json:"from"`
	Code    string `json:"code"`
	Content string `json:"content"`
}

type Message struct {
	Sender    string `json:"sender"`
	Recipient string `json:"recipient"`
	Content   string `json:"content"`
}

package service

import (
	"Tiktok/biz/model/ai"
	"Tiktok/pkg/config"
	"bytes"
	"encoding/json"
	"io"
	"net/http"

	"github.com/pkg/errors"
)

func ChatWithAi(content string) (*ai.GLMResp, error) {
	url := "https://open.bigmodel.cn/api/paas/v4/chat/completions"
	reqBody := &ai.GLMReq{
		Model: "glm-4.7-flash",
		Messages: []*ai.ChatMessage{
			{
				Role:    "system",
				Content: "你是一个有用的AI助手",
			},
			{
				Role:    "user",
				Content: content,
			},
			{},
		},
		Stream:      false,
		Temperature: 1,
	}
	payload, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(payload))
	req.Header.Add("Authorization", "Bearer "+config.Cfg.Api.ApiKey)
	req.Header.Add("Content-Type", "application/json")
	res, _ := http.DefaultClient.Do(req)
	defer res.Body.Close()
	body, _ := io.ReadAll(res.Body)
	var resp ai.GLMResp
	if err := json.Unmarshal(body, &resp); err == nil {
		return &resp, nil
	}
	var errResp ai.GLMResp
	err := json.Unmarshal(body, &errResp)
	if err != nil {
		return nil, err
	}
	return nil, errors.New("未知错误:" + string(body))
}

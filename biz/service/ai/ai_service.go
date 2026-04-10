package ai

import (
	"Tiktok/biz/model/ai"
	"Tiktok/pkg/config"
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"time"

	"github.com/pkg/errors"
)

type AiService struct{}

func NewAiService() *AiService {
	return &AiService{}
}

func (s *AiService) ChatWithAi(content string) (*ai.GLMResp, error) {
	url := "https://open.bigmodel.cn/api/paas/v4/chat/completions"
	reqBody := &ai.GLMReq{
		Model: "glm-4-air",
		Messages: []*ai.ChatMessage{
			{
				Role:    "system",
				Content: "你是一个有用的AI助手",
			},
			{
				Role:    "user",
				Content: content,
			},
		},
		Stream:      false,
		Temperature: 1,
	}
	payload, err := json.Marshal(reqBody)
	if err != nil {
		return nil, errors.Wrap(err, "json marshal error")
	}
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(payload))
	if err != nil {
		return nil, errors.Wrap(err, "create request error")
	}
	req.Header.Add("Authorization", "Bearer "+config.Cfg.Api.ApiKey)
	req.Header.Add("Content-Type", "application/json")
	client := &http.Client{
		Timeout: 180 * time.Second,
	}
	res, err := client.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "http request error")
	}
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, errors.Wrap(err, "read response body error")
	}
	if res.StatusCode != http.StatusOK {
		return nil, errors.Errorf("API error: status code %d, body: %s", res.StatusCode, string(body))
	}
	var resp ai.GLMResp
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, errors.Wrap(err, "json unmarshal error")
	}
	return &resp, nil
}

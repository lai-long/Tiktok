package ai

import (
	"Tiktok/pkg/config"
	"context"
	"log"
	"testing"
)

func TestAI(t *testing.T) {
	_, err := config.Load([]string{"/home/lai-long/Tiktok/pkg/config"})
	if err != nil {
		log.Println("config load err", err)
	}
	ctx := context.Background()
	agent := NewAgent(ctx)
	resp := agent.StartAction("能不能告诉我永泰最近的天气怎么样，把你的全部思考过程和工具调用情况都告诉我")
	log.Println("action resp:", resp)
}

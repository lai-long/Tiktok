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
	resp := agent.StartAction("1+1等于多少")
	log.Println("action resp:", resp)
}

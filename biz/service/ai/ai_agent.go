package ai

import (
	"context"
	"log"
)

type Agent struct {
	chatCli *ChatClient
	ctx     context.Context
}

func NewAgent(ctx context.Context) *Agent {
	llm, err := NewChatClient(ctx)
	if err != nil {
		log.Println("mcp NewChatClient fail", err)
		return nil
	}
	if llm == nil {
		log.Println("mcp NewChatClient fail llm is nil")
		return nil
	}
	return &Agent{
		chatCli: llm,
		ctx:     ctx,
	}
}

func (a *Agent) StartAction(prompt string) string {
	if a.chatCli == nil {
		log.Println("start action err short of chat client")
		return ""
	}
	resp, err := a.chatCli.Chat(prompt)
	if err != nil {
		log.Println("start action fail", err)
		return ""
	}
	return resp
}

func (a *Agent) StopAction() {
	if a.chatCli == nil {
		log.Println("stop action err short of chat client")
		return
	}
	a.chatCli.Close()
}

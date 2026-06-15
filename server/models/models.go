package models

import (
	"fmt"
	"sync"
	"time"

	"github.com/coder/websocket"
)

type Bot struct {
	Mu sync.RWMutex

	ID       string `json:"id"`
	OS       string `json:"os"`
	HostName string `json:"hostname"`
	LastSeen time.Time
	Active   bool
	Con      *websocket.Conn
	// Command Command
}

func (b *Bot) UpdateLastseen() {
	b.Mu.Lock()
	defer b.Mu.Unlock()

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

type Ostype struct {
	Ext  string
	Arch string
}

type BotCreds struct {
	ID        string
	SecretKey string
}

func (bc *BotCreds) CompileBot(ostype Ostype) {
	ldflagsvalue := fmt.Sprintf("-X main.AgentID=%s -X 'main.Registration=%s'", bc.ID, bc.SecretKey)

}

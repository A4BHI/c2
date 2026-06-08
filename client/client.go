package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/coder/websocket"
	"github.com/coder/websocket/wsjson"
)

type c2 struct {
	con *websocket.Conn
}

func (c c2) connectToserver() {

}

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	c, _, err := websocket.Dial(ctx, "ws://localhost:8080/api", nil)
	if err != nil {
		log.Fatal(err)
	}
	defer c.CloseNow()
	var msg string
	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Println("enter your message: ")
		scanner.Scan()
		msg = scanner.Text()
		if msg == "exit" {
			break
		}
		err = wsjson.Write(ctx, c, msg)
		if err != nil {
			fmt.Println(err)
		}
		var v string
		wsjson.Read(ctx, c, v)
		fmt.Println("recieved : ", v)
	}

	c.Close(websocket.StatusNormalClosure, "")
}

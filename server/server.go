package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/coder/websocket"
	"github.com/coder/websocket/wsjson"
)

type Bot struct {
	ID       string `json:"id"`
	OS       string `jsin:"os"`
	IP       string `json:"ip"`
	LastSeen string
	// Command Command
}

// type Command struct {
// 	Cmdname string
// 	Result  string
// }

// type c2 struct {
// 	bots map[string]Bot
// }

func (b *Bot) connectBot(w http.ResponseWriter, r *http.Request) {
	var con *websocket.Conn
	var err error

	if con, err = websocket.Accept(w, r, nil); err != nil {
		log.Println(err)
		return
	}

	ctx := context.Background()
	// read
	go func() {
		for {

			if err := wsjson.Read(ctx, con, b); err != nil {
				log.Println(err)
				return

			}
			b.LastSeen = time.Now().Format("03:04:05PM")
			fmt.Print(b)
		}
	}()

	con.Close(websocket.StatusNormalClosure, "")

}

func (b *Bot) heartBeat() {}
func main() {
	b := Bot{}
	http.HandleFunc("/connect", b.connectBot)
	http.ListenAndServe(":4444", nil)

}

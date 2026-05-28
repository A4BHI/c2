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
	Active   bool
	// Command Command
}

// type Command struct {
// 	Cmdname string
// 	Result  string
// }

// type c2 struct {
// 	bots map[string]Bot
// }

func connectBot(w http.ResponseWriter, r *http.Request) {
	b := Bot{
		Active: false,
	}
	var con *websocket.Conn
	var err error
	if con, err = websocket.Accept(w, r, nil); err != nil {
		log.Println(err)
		return
	}

	defer con.Close(websocket.StatusNormalClosure, "")
	ctx := context.Background()
	// read
	go heartBeat(con, &b)
	for {
		if err := wsjson.Read(ctx, con, b); err != nil {
			log.Println(err)
			return
		}
		b.LastSeen = time.Now().Format("03:04:05PM")
		b.Active = true
		fmt.Print(b)
	}

}

func heartBeat(con *websocket.Conn, bot *Bot) {
	if bot.Active {
		time.Sleep(1 * time.Minute)
	}
}
func main() {

	http.HandleFunc("/connect", connectBot)
	http.ListenAndServe(":4444", nil)

}

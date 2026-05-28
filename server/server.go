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
	LastSeen time.Time
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

	go heartBeat(con, &b)
	for {
		if err := wsjson.Read(ctx, con, b); err != nil {
			log.Println(err)
			return
		}
		b.LastSeen = time.Now()
		b.Active = true
		fmt.Print(b)
	}

}

func heartBeat(con *websocket.Conn, bot *Bot) {

	timer := time.Tick(10 * time.Second)
	inactiveSince := time.Now()
	for range timer {
		if !bot.Active {
			if time.Since(inactiveSince) > 15*time.Second {
				fmt.Println("Closing now...")
				con.CloseNow()
				return
			}
			continue
		}

		if time.Since(bot.LastSeen) > 1*time.Minute {
			fmt.Println("Closing now..")
			con.CloseNow()
			return
		}

	}

}
func main() {

	http.HandleFunc("/connect", connectBot)
	http.ListenAndServe(":4444", nil)

}

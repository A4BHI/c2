package main

type Bot struct {
	ID int    `json:"id"`
	OS string `json:"os"`
	HostName string `json:"hostname"`
}

func getSysInfo() *Bot {

	return &Bot{
		ID: ,
	}
}

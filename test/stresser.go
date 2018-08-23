package main

import (
	"fmt"
	"time"

	"github.com/gorilla/websocket"
)

func main() {

	mockWebsocket, _, err := websocket.DefaultDialer.Dial("wss://engineerbeard.com/bchatws", nil)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	mockWebsocket.SetReadDeadline(time.Now().Add(time.Second * time.Duration(3)))
	mockWebsocket.SetWriteDeadline(time.Now().Add(time.Second * time.Duration(3)))
	mockWebsocket.Close()
}

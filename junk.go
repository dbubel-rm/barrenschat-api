package main

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/gorilla/websocket"
)

func main() {
	u := "ws" + strings.TrimPrefix("http://localhost:9000/bchatws", "http")
	ws, _, err := websocket.DefaultDialer.Dial(u, nil)
	if err != nil {
		fmt.Println(err.Error())
	}
	// ws.SetPongHandler(func(string) error {
	// 	time.Sleep(time.Second * 1)
	// 	ws.SetReadDeadline(time.Now().Add(time.Second * 10))
	// 	log.Println("PONG rec")
	// 	err := ws.WriteMessage(websocket.PingMessage, nil)
	// 	if err != nil {
	// 		log.Println(err.Error())
	// 	}

	// 	return nil
	// })

	ws.SetPingHandler(func(string) error {

		ws.SetWriteDeadline(time.Now().Add(time.Second * 10))
		ws.SetReadDeadline(time.Now().Add(time.Second * 10))
		log.Println("PING rec")
		err := ws.WriteMessage(websocket.PongMessage, nil)
		if err != nil {
			log.Println(err.Error())
		}
		return nil
	})

	ws.WriteMessage(websocket.PingMessage, nil)
	for {
		msgType, msg, err := ws.ReadMessage()
		log.Println(msgType, msg, err.Error())
	}

}

package client

import (
	"time"

	"log"

	"github.com/engineerbeard/barrenschat-api/hub"
	"github.com/gorilla/websocket"
)

const (
	writeWait      = 10 * time.Second
	pongWait       = 60 * time.Second
	pingPeriod     = (pongWait * 9) / 10
	maxMessageSize = 512
)

func NewClient(c *websocket.Conn, h *hub.Hub) {
	go readPump(c, h)
}

func readPump(c *websocket.Conn, h *hub.Hub) {
	defer func() {
		log.Printf("Closing connection for [%s]\n", c.RemoteAddr().String())
		//c.hub.MsgRecvr <- hub.Message{MsgType: "client disconnect"}
		c.Close()
	}()
	c.SetReadLimit(maxMessageSize)
	c.SetReadDeadline(time.Now().Add(pongWait))
	c.SetPongHandler(func(string) error { c.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	for {

		var msg hub.Message
		err := c.ReadJSON(&msg)
		if err != nil {
			h.ClientDisconnect <- c
			log.Println(err.Error())
			break
		}

		switch msgType := msg.MsgType; msgType {
		case "newconnection":
			h.NewConnection <- c
			err = h.RedisClient.Publish("newconnection", "user data").Err()
			if err != nil {
				panic(err)
			}
		case "newmessage":
			err = h.RedisClient.Publish("newmessage", "new message").Err()
			if err != nil {
				panic(err)
			}
		default:
			log.Println("bad message")
		}
	}
}

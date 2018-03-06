package client

import (
	"time"

	"github.com/engineerbeard/barrenschat-api/hub"
	"github.com/gorilla/websocket"
	logging "github.com/op/go-logging"
)

var log = logging.MustGetLogger("barrenschat-api")

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
		log.Debugf("Closing connection for [%s]\n", c.RemoteAddr().String())
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
			log.Error(err.Error())
			break
		}
		//log.Debug("Message Type", mType, string(message))

		switch msgType := msg.MsgType; msgType {
		case "new connection":
			h.NewConnection <- c
		case "new message":
			h.NewMessage <- msg.Data.(string)
		case "client info":
			h.ClientInfo <- msg.Data.(string)
		default:
			log.Debug("bad message")
		}
	}
}

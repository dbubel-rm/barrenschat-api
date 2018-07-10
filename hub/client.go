package hub

import "github.com/gorilla/websocket"

type Client struct {
	conn        *websocket.Conn
	closeChan   chan int
	newMsgChan  chan string
	channelName string
	Claims      map[string]string
}

// TODO: send message

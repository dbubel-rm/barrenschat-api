package hub

import (
	"encoding/json"
	"log"

	"github.com/go-redis/redis"
)

// Hub maintains the set of active clients and broadcasts messages to the
// clients.
type Hub struct {
	channels map[string]map[*Client]bool

	// Inbound messages from the clients.
	broadcast chan []byte

	// Register requests from the clients.
	clientConnect chan *Client

	// Unregister requests from clients.
	clientDisconnect chan *Client

	// Channel that messages from other hubs come from
	pubSubRecv chan []byte

	msgRouter map[string]func(rawMessage)
}

const (
	redisPubSubChannel       string = "datapipe"
	MESSAGE_TYPE_NEW         string = "message_new"
	MESSAGE_TEXT             string = "message_text"
	MESSAGE_TYPE_NEW_CHANNEL string = "message_new_channel"
)

var redisClient *redis.Client

func init() {
	redisClient = redis.NewClient(&redis.Options{
		Addr:     "redis:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	redisClient.Set("hi", "dean", 0)
}

// NewHub used to create a new hub instance
func NewHub() *Hub {
	return &Hub{
		//clients:    make(map[*Client]bool),
		channels:         make(map[string]map[*Client]bool),
		clientConnect:    make(chan *Client),
		clientDisconnect: make(chan *Client),
		pubSubRecv:       make(chan []byte),
		broadcast:        make(chan []byte),
		msgRouter:        make(map[string]func(rawMessage)),
	}
}

func (h *Hub) pubSubListen(pubSub string) {
	// defer h.RedisClient.Close()
	pSub := redisClient.Subscribe(pubSub)
	for {
		msg, err := pSub.ReceiveMessage()
		if err != nil {
			log.Println(err.Error())
		}
		h.pubSubRecv <- []byte(msg.Payload)
	}
}

func (h *Hub) getChannels() {
	channelList := redisClient.PubSubChannels("*")
	for _, j := range channelList.Val() {
		log.Println(j)
	}
}

func (h *Hub) createChannel(s string) {
	if _, ok := h.channels[s]; !ok {
		h.channels[s] = make(map[*Client]bool)
	}
}

// Run starts the hub listening on its channels
func (h *Hub) Run() {
	h.addHandler(MESSAGE_TYPE_NEW, h.handleClientMessage)
	h.addHandler(MESSAGE_TYPE_NEW_CHANNEL, h.handleCreateNewChannel)

	go h.pubSubListen(redisPubSubChannel)

	// Main program loop, listens for messages from clients and from redis
	for {
		select {
		case client := <-h.clientConnect:
			log.Println("in connect", client.channelsSubscribedTo)
			for _, clientChannel := range client.channelsSubscribedTo {
				if _, ok := h.channels[clientChannel]; !ok {
					h.channels[clientChannel] = make(map[*Client]bool)
				}
				h.channels[clientChannel][client] = true
			}

		case client := <-h.clientDisconnect:
			for channel := range h.channels {
				if _, ok := h.channels[channel][client]; ok {
					delete(h.channels[channel], client)
					close(client.send)
					return
				}
			}
		case message := <-h.broadcast:
			// We received a message from a client connected to this hub
			result := redisClient.Publish(redisPubSubChannel, message)
			if result.Err() != nil {
				log.Println(result.Err().Error())
			}
		case message := <-h.pubSubRecv:
			log.Println(message)
			// We received a message from redis
			var m rawMessage
			err := json.Unmarshal(message, &m)
			if err != nil {
				log.Println(err.Error())
			}
			// log.Println("New msg", m)

			if handler, found := h.findHandler(m.MsgType); found {
				handler(m)
			}
		}
	}
}

// package hub

// import (
// 	"bytes"
// 	"log"
// 	"math/rand"
// 	"net/http"
// 	"sync"
// 	"time"

// 	"github.com/go-redis/redis"
// 	"github.com/gorilla/websocket"
// )

// var redisClient *redis.Client

// // Message Types
// const TYPE_NEW_MESSAGE = "message_new"

// // Payloads
// const MESSAGE_TEXT = "message_text"
// const MESSAGE_USER = "user_name"

// type message struct {
// 	MsgType string                 `json:"msgType"`
// 	Payload map[string]interface{} `json:"payload"`
// }
// type newClient struct {
// 	Ws     *websocket.Conn
// 	Claims map[string]string
// }
// type hub struct {
// 	NewConnection    chan newClient
// 	ClientDisconnect chan *websocket.Conn
// 	RoomList         map[string][]*Client
// 	RedisClient      *redis.Client
// 	Router           map[string]func(map[string]interface{})
// 	m                sync.Mutex
// }

// func NewHub() *hub {
// 	// TODO: fail if redis isnt started
// 	rand.Seed(time.Now().Unix())
// 	x := &hub{
// 		Router:           make(map[string]func(map[string]interface{})),
// 		NewConnection:    make(chan newClient),
// 		ClientDisconnect: make(chan *websocket.Conn, 1),
// 		RoomList:         make(map[string][]*Client),
// 		RedisClient:      redisClient,
// 	}
// 	log.Println("Redis:", x.RedisClient.Ping().String())
// 	x.listenForNewMessages()
// 	return x
// }
// func init() {
// 	redisClient = redis.NewClient(&redis.Options{
// 		Addr:     "redis:6379",
// 		Password: "", // no password set
// 		DB:       0,  // use default DB
// 	})
// }

// func sendDBLog(token string) {
// 	payload := []byte(`{"name":"Dealer Z","destinationDnis":"4151112222"}`)
// 	req, _ := http.NewRequest("POST", "https://barrenschat-27212.firebaseio.com/message_list.json", bytes.NewBuffer(payload))

// 	q := req.URL.Query()
// 	q.Add("auth", token)
// 	req.URL.RawQuery = q.Encode()
// 	log.Println(req.URL.String())
// 	cc := &http.Client{}
// 	res, e := cc.Do(req)
// 	if e != nil {
// 		log.Println(e.Error())
// 	} else {
// 		log.Println(res.Body)
// 	}
// }

// func (h *hub) newClientConnection(c newClient) {
// 	log.Printf("New client connection [%s]\n", c.Ws.RemoteAddr().String())

// 	newClient := &Client{
// 		conn:        c.Ws,
// 		channelName: "main",
// 		newMsgChan:  make(chan string),
// 		closeChan:   make(chan int),
// 		Claims:      c.Claims,
// 	}
// 	//log.Println(newClient.Claims["user_id"])
// 	h.RoomList["main"] = append(h.RoomList["main"], newClient)

// 	//Reader
// 	go func(c *websocket.Conn, h *hub, client *Client) {

// 		defer func() {
// 			log.Printf("Closing reader for [%s]\n", c.RemoteAddr().String())
// 			c.Close()
// 			//client.closeChan <- 1
// 		}()

// 		for {
// 			mmsg := message{}
// 			err := c.ReadJSON(&mmsg)
// 			log.Println("Msg RECV:", mmsg)

// 			if err != nil {
// 				log.Println(err.Error())
// 				break
// 			}
// 			if handler, found := h.findHandler(mmsg.MsgType); found {
// 				handler(mmsg.Payload)
// 			} else {
// 				log.Println(mmsg.MsgType, "Not found")
// 			}
// 		}
// 	}(c.Ws, h, newClient)

// 	// Writer
// 	go func(c *websocket.Conn, h *hub, client *Client) {
// 		defer func() {
// 			log.Printf("Closing writer for [%s]\n", c.RemoteAddr().String())
// 			h.ClientDisconnect <- c
// 		}()
// 		ticker := time.NewTicker(time.Second * 5)
// 		for {
// 			select {
// 			case <-client.closeChan:
// 				log.Println("Stopping ticker")
// 				ticker.Stop()
// 				return
// 			case sendMsg := <-client.newMsgChan:
// 				d := make(map[string]interface{})
// 				d[MESSAGE_TEXT] = sendMsg
// 				packet := message{MsgType: TYPE_NEW_MESSAGE, Payload: d}
// 				c.WriteJSON(packet)
// 			case <-ticker.C:
// 				err := c.WriteMessage(websocket.PingMessage, nil)
// 				if err != nil {
// 					log.Println(err.Error())
// 					return
// 				}
// 			}
// 		}
// 	}(c.Ws, h, newClient)
// }

// func (h *hub) removeCLient(c *websocket.Conn) {
// 	// h.m.Lock()
// 	// defer h.m.Unlock()
// 	// for _, j := range h.RoomList {
// 	// 	for i := 0; i < len(j); i++ {
// 	// 		if c == j[i].conn {
// 	// 			close(j[i].closeChan)
// 	// 			close(j[i].newMsgChan)
// 	// 			log.Printf("Removed [%s]\n", c.RemoteAddr().String())
// 	// 			h.RoomList[j[i].channelName] = append(h.RoomList[j[i].channelName][:i], h.RoomList[j[i].channelName][i+1:]...)
// 	// 			return
// 	// 		}
// 	// 	}
// 	// }
// }

// func (h *hub) Run() {
// 	h.handleMsg(TYPE_NEW_MESSAGE, h.handleClientMessage)
// 	h.handleMsg("client_info", h.handleUpdateClientInfo)
// 	h.handleMsg("command_who", h.handleWhoCommand)
// 	for {
// 		select {
// 		case c := <-h.NewConnection:
// 			h.newClientConnection(c)
// 		case c := <-h.ClientDisconnect:
// 			h.removeCLient(c)
// 		}
// 	}
// }

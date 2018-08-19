package hub

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/go-redis/redis"
)

// Hub maintains the set of active clients and broadcasts messages to the
// clients.
type Hub struct {
	blocker  chan bool
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

	msgRouterClient map[string]func(rawMessage, *Client)
}

const (
	redisPubSubChannel        string = "datapipe"
	MESSAGE_TYPE_NEW          string = "message_new"
	MESSAGE_TEXT              string = "message_text"
	MESSAGE_TYPE_NEW_CHANNEL  string = "message_new_channel"
	MESSAGE_TYPE_JOIN_CHANNEL string = "message_join_channel"
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
		msgRouterClient:  make(map[string]func(rawMessage, *Client)),
		blocker:          make(chan bool, 1),
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
		log.Println("made new channel")
	}
	<-h.blocker
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
					break
				}
			}
		case message := <-h.broadcast:
			// We received a message from a client connected to this hub
			result := redisClient.Publish(redisPubSubChannel, message)
			if result.Err() != nil {
				log.Println(result.Err().Error())
			}
		case message := <-h.pubSubRecv:
			log.Println(string(message))
			// We received a message from redis
			var m rawMessage
			err := json.Unmarshal(message, &m)
			if err != nil {
				log.Println(err.Error())
			}
			// log.Println("New msg", m)

			if handler, found := h.findHandler(m.MsgType); found {
				fmt.Println("handler found")
				handler(m)
			}
		}
	}
}

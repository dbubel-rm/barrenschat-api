package hub

import (
	"encoding/json"
	"log"

	"github.com/engineerbeard/barrenschat-api/config"
	"github.com/go-redis/redis"
)

type Hub struct {
	blocker          chan bool
	clients          map[string]*Client     // client IDs to Client
	channelListeners map[string]chan []byte // Map of channel names to redis pubsub stream
	channelMembers   map[string][]*Client   // Map of channel names to clients
	broadcast        chan []byte
	clientConnect    chan *Client
	clientDisconnect chan *Client
	msgRouter        map[string]func(rawMessage)
}

const (
	redisPubSubChannel    string = "datapipe"
	MessageTypeNew        string = "message_new"
	MessageText           string = "message_text"
	MessageTypeNewChannel string = "message_new_channel"
	// MESSAGE_TYPE_JOIN_CHANNEL string = "message_join_channel"
	// MESSAGE_TYPE_GET_CHANNELS string = "message_get_channels"
)

var redisClient *redis.Client

func init() {
	redisClient = redis.NewClient(&redis.Options{
		Addr:     config.RedisURL,
		Password: "", // no password set
		DB:       0,  // use default DB
	})
}

// NewHub used to create a new hub instance
func NewHub() *Hub {
	return &Hub{
		clients:          make(map[string]*Client),
		channelListeners: make(map[string]chan []byte),
		clientConnect:    make(chan *Client),
		clientDisconnect: make(chan *Client),
		channelMembers:   make(map[string][]*Client),
		broadcast:        make(chan []byte),
		msgRouter:        make(map[string]func(rawMessage)),
		blocker:          make(chan bool, 1),
	}
}

// func (h *Hub) pubSubListen(pubSub string) {
// 	// defer h.RedisClient.Close()
// 	pSub := redisClient.Subscribe(pubSub)
// 	for {
// 		msg, err := pSub.ReceiveMessage()
// 		if err != nil {
// 			log.Println(err.Error())
// 		}
// 		h.pubSubRecv <- []byte(msg.Payload)
// 	}
// }

// func (h *Hub) getChannels() {
// 	channelList := redisClient.PubSubChannels("*")
// 	for _, j := range channelList.Val() {
// 		log.Println(j)
// 	}
// }

// func (h *Hub) createChannel(s string) {
// 	if _, ok := h.channels[s]; !ok {
// 		h.channels[s] = make(map[*Client]bool)
// 		log.Println("made new channel")
// 	}
// 	<-h.blocker
// }
func (h *Hub) newChannelListener(clientChannel string) {
	pSub := redisClient.Subscribe(clientChannel) // h.channelProducers[clientChannel] = pSub
	cc := make(chan []byte)
	h.channelListeners[clientChannel] = cc

	go func(c chan []byte, ps *redis.PubSub) {
		for {
			msg, err := ps.ReceiveMessage()
			if err != nil {
				log.Println(err.Error())
			}

			var m rawMessage
			err = json.Unmarshal([]byte(msg.Payload), &m)
			if err != nil {
				log.Println(err.Error())
			}

			if handler, found := h.findHandler(m.MsgType); found {
				handler(m)
			}
		}
	}(cc, pSub)
}

// Run starts the hub listening on its channels
func (h *Hub) Run() {
	h.addHandler(MessageTypeNew, h.handleClientMessage)
	// h.addHandler(MESSAGE_TYPE_NEW_CHANNEL, h.handleCreateNewChannel)

	// go h.pubSubListen(redisPubSubChannel)

	// Main program loop, listens for messages from clients and from redis
	for {
		select {
		case client := <-h.clientConnect:
			for _, clientChannel := range client.channelsSubscribedTo {
				if _, ok := h.channelListeners[clientChannel]; !ok {
					h.newChannelListener(clientChannel)
				}
				h.clients[client.ID] = client
				h.channelMembers[clientChannel] = []*Client{}
				h.channelMembers[clientChannel] = append(h.channelMembers[clientChannel], client)
				// client
			}

		// case client := <-h.clientDisconnect:
		// 	for channel := range h.channels {
		// 		if _, ok := h.channels[channel][client]; ok {
		// 			delete(h.channels[channel], client)
		// 			close(client.send)
		// 			break
		// 		}
		// 	}
		case message := <-h.broadcast:
			log.Println(string(message))
			// We received a message from a client connected to this hub
			// result := redisClient.Publish(redisPubSubChannel, message)
			// if result.Err() != nil {
			// 	log.Println(result.Err().Error())
			// }
			// case message := <-h.pubSubRecv:
			// 	log.Println(string(message))
			// 	// We received a message from redis
			var m rawMessage
			err := json.Unmarshal(message, &m)
			if err != nil {
				log.Println(err.Error())
			}
			// log.Println("New msg", m.Payload["channel"].(string))
			redisClient.Publish(m.Payload["channel"].(string), message)
			// h.channelProducers[m.Payload["channel"].(string)].p <- )

			// if handler, found := h.findHandler(m.MsgType); found {
			// 	fmt.Println("handler found")
			// 	handler(m)
			// }
		}
	}
}

package hub

import (
	"encoding/json"
	"log"

	"github.com/engineerbeard/barrenschat-api/config"
	"github.com/go-redis/redis"
)

type Hub struct {
	locker                chan bool
	clients               map[string]*Client     // Map of client IDs to *Client
	channelMembers        map[string][]*Client   // Map of channel names to []*Client
	topicChannels         map[string]chan []byte // Map of channel names to redis pubsub stream
	incomingClientMessags chan []byte
	clientConnect         chan *Client
	clientDisconnect      chan *Client
	msgRouter             map[string]func(rawMessage) // Map of message type to handler function

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
		clients:               make(map[string]*Client),
		clientConnect:         make(chan *Client),
		clientDisconnect:      make(chan *Client), // todo: remove in favor of a message
		topicChannels:         make(map[string]chan []byte),
		channelMembers:        make(map[string][]*Client),
		incomingClientMessags: make(chan []byte),
		msgRouter:             make(map[string]func(rawMessage)),
		locker:                make(chan bool, 1),
	}
}

func (h *Hub) getClients() map[string]*Client {
	h.locker <- true
	c := h.clients
	<-h.locker
	return c
}

func (h *Hub) createNewChannel(name string) error {
	return nil
}

// func (h *Hub) getChannels() {
// 	channelList := redisClient.PubSubChannels("*")
// 	for _, j := range channelList.Val() {
// 		log.Println(j)
// 	}
// }

func (h *Hub) newChannelListener(clientChannel string) {
	pSub := redisClient.Subscribe(clientChannel)
	cc := make(chan []byte)
	h.topicChannels[clientChannel] = cc

	go func(c chan []byte, ps *redis.PubSub) {

		for {
			msg, err := ps.ReceiveMessage()
			if err != nil {
				log.Println("ERROR:", err.Error())
			}

			var m rawMessage
			err = json.Unmarshal([]byte(msg.Payload), &m)
			if err != nil {
				log.Println("ERROR:", err.Error())
			}

			if handler, found := h.findHandler(m.MsgType); found {
				handler(m)
			} else {
				log.Println("WARN:", "No message type found")
			}
		}
	}(cc, pSub)
}

// Run starts the hub listening on its channels
func (h *Hub) Run() {
	h.addHandler(MessageTypeNew, h.handleClientMessage)

	for {
		select {
		case client := <-h.clientConnect:
			h.locker <- true
			for _, clientChannel := range client.channelsSubscribedTo {
				if _, ok := h.topicChannels[clientChannel]; !ok {
					h.newChannelListener(clientChannel)
				}
				h.clients[client.getClientID()] = client
				h.channelMembers[clientChannel] = []*Client{}
				h.channelMembers[clientChannel] = append(h.channelMembers[clientChannel], client)
			}
			<-h.locker
		case client := <-h.clientDisconnect:
			// log.Println("Before remove client", h.clients)
			delete(h.getClients(), client.getClientID())
			h.locker <- true
			// log.Println("After remove client", h.clients)
			for _, channel := range client.channelsSubscribedTo {
				for i := range h.channelMembers[channel] {
					if client == h.channelMembers[channel][i] {
						copy(h.channelMembers[channel][i:], h.channelMembers[channel][i+1:])
						h.channelMembers[channel][len(h.channelMembers[channel])-1] = nil
						h.channelMembers[channel] = h.channelMembers[channel][:len(h.channelMembers[channel])-1]
					}
				}
			}
			<-h.locker
		case message := <-h.incomingClientMessags:

			log.Println(string(message))
			var m rawMessage
			err := json.Unmarshal(message, &m)
			if err != nil {
				log.Println(err.Error())
			}
			redisClient.Publish(m.Payload["channel"].(string), message)
		}
	}
}

package hub

import (
	"fmt"
	"time"

	"github.com/go-redis/redis"
	"github.com/gorilla/websocket"
	logging "github.com/op/go-logging"
)

var log = logging.MustGetLogger("barrenschat-api")

type Message struct {
	MsgType string      `json:"msgType"`
	Data    interface{} `json:"data"`
}

type Hub struct {
	NewConnection    chan *websocket.Conn
	ClientInfo       chan string
	ClientDisconnect chan *websocket.Conn
	NewMessage       chan string
	RoomList         map[string][]*Client
	RedisClient      *redis.Client
}

type Client struct {
	conn *websocket.Conn
	name string
	room string
}

func NewHub() *Hub {

	return &Hub{
		NewConnection:    make(chan *websocket.Conn),
		ClientDisconnect: make(chan *websocket.Conn),
		NewMessage:       make(chan string),
		RoomList:         make(map[string][]*Client),
		RedisClient: redis.NewClient(&redis.Options{
			Addr:     "redis:6379",
			Password: "", // no password set
			DB:       0,  // use default DB
		}),
	}
}

func (h *Hub) GetChannels() {
	channelList := h.RedisClient.PubSubChannels("*")
	for _, j := range channelList.Val() {
		fmt.Println(j)
	}
}

func (h *Hub) publishMessage(channel, message string) {

	p := h.RedisClient.Subscribe("mychannel1")
	//fmt.Println(h.RedisClient.PubSubChannels("*"))
	defer h.RedisClient.Close()

	// Wait for subscription to be created before publishing message.
	subscr, err := p.ReceiveTimeout(time.Second)
	if err != nil {
		panic(err)
	}
	fmt.Println(subscr)

	err = h.RedisClient.Publish("mychannel1", "hello").Err()
	if err != nil {
		panic(err)
	}

	msg, err := p.ReceiveMessage()
	if err != nil {
		panic(err)
	}

	fmt.Println(msg.Channel, msg.Payload)
}

func (h *Hub) newClient(c *websocket.Conn) {
	h.GetChannels()
	h.RoomList["main"] = append(h.RoomList["main"], &Client{conn: c, name: "dean", room: "main"})
	// for x, i := range h.RoomList["main"] {
	// 	fmt.Println(x, *i)
	// }
}

func (h *Hub) removeCLient(c *websocket.Conn) {
	for _, j := range h.RoomList {
		for i := 0; i < len(j); i++ {
			if c == j[i].conn {
				h.RoomList[j[i].room] = append(h.RoomList[j[i].room][:i], h.RoomList[j[i].room][i+1:]...)
				return
			}
		}
	}
}

func (h *Hub) Run() {
	for {
		select {
		case c := <-h.NewConnection:
			h.newClient(c)
			log.Debugf("New client: %s", c.RemoteAddr())
		case msg := <-h.NewMessage:
			log.Debug("New message recv:", msg)
		case c := <-h.ClientDisconnect:
			h.removeCLient(c)
		}
	}
}

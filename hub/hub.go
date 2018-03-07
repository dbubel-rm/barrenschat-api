package hub

import (
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
	x := &Hub{
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
	x.listenRedis()
	return x
}

func (h *Hub) GetChannels() {
	channelList := h.RedisClient.PubSubChannels("*")
	for _, j := range channelList.Val() {
		_ = j
		//log.Debug(j)
	}
}

func (h *Hub) listenRedis() {
	log.Debug("Listening redis...")

	go func() {
		p := h.RedisClient.Subscribe("newconnection")
		for {
			log.Debug("Waiting on new connections...")
			msg, err := p.ReceiveMessage()
			if err != nil {
				panic(err)
			}
			log.Debug("New Client from redis:", msg)
		}
	}()

	go func() {
		p := h.RedisClient.Subscribe("newmessage")
		for {
			log.Debug("Waiting on new messages...")
			msg, err := p.ReceiveMessage()
			if err != nil {
				panic(err)
			}
			h.NewMessage <- msg.String()
			log.Debug("New message from redis:", msg)
		}
	}()

	//fmt.Println(h.RedisClient.PubSubChannels("*"))
	// defer h.RedisClient.Close()

	// // Wait for subscription to be created before publishing message.
	// subscr, err := p.ReceiveTimeout(time.Second)
	// if err != nil {
	// 	panic(err)
	// }
	// fmt.Println(subscr)

	// err = h.RedisClient.Publish("mychannel1", "hello").Err()
	// if err != nil {
	// 	panic(err)
	// }

	// fmt.Println(msg.Channel, msg.Payload)
}
func (h *Hub) Broadcast(msg string) {
	for _, j := range h.RoomList {
		for i := 0; i < len(j); i++ {
			data := struct {
				Paste            string
				KeepAlive        bool
				BurnAfterReading bool
			}{
				"FUCK YEA",
				true,
				true,
			}
			j[i].conn.WriteJSON(data)
		}
	}
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
			h.Broadcast(msg)
		case c := <-h.ClientDisconnect:
			h.removeCLient(c)
		}
	}
}

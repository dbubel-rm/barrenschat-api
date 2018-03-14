package hub

import (
	"encoding/json"
	"log"
	"time"

	"github.com/go-redis/redis"
	"github.com/gorilla/websocket"
)

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
	conn      *websocket.Conn
	closeChan chan int
	sendChan  chan string
	name      string
	room      string
}

func NewHub() *Hub {

	// TODO: fail if redis isnt started

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
	log.Println("Redis:", x.RedisClient.Ping().String())
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

type msgRouter struct {
	routes map[string]func()
}

func (r *msgRouter) addRoute(route string, f func()) {
	r.routes[route] = f
}

func (h *Hub) listenRedis() {
	go func() {
		defer h.RedisClient.Close()
		pSub := h.RedisClient.Subscribe("datapipe")

		// if subscr, err := pSub.ReceiveTimeout(time.Second); err == nil {
		// 	log.Println(subscr)
		// } else {
		// 	log.Println(err.Error())
		// 	panic(err)
		// }

		for {
			if msg, err := pSub.ReceiveMessage(); err == nil {
				log.Println("New msg from datapipe:", msg.Payload)
				h.Broadcast(msg.Payload)
			} else {
				break
			}
		}
	}()
}

func (h *Hub) Broadcast(msg string) {
	for _, j := range h.RoomList {
		for i := 0; i < len(j); i++ {
			data := struct {
				Paste string
			}{
				msg,
			}
			b, _ := json.Marshal(data)
			j[i].sendChan <- string(b)
		}
	}
}

func (h *Hub) newClient(c *websocket.Conn) {
	log.Printf("New client connection [%s]\n", c.RemoteAddr().String())
	//h.GetChannels()
	cc := make(chan string)
	closeChan := make(chan int)
	h.RoomList["main"] = append(h.RoomList["main"], &Client{
		conn:      c,
		name:      "dean",
		room:      "main",
		sendChan:  cc,
		closeChan: closeChan,
	})

	//Reader
	go func(c *websocket.Conn, h *Hub, close chan int) {

		defer func() {
			log.Printf("Closing reader for [%s]\n", c.RemoteAddr().String())
			c.Close()
			close <- 1
			h.ClientDisconnect <- c
		}()
		for {
			_, msg, err := c.ReadMessage()

			if err != nil {
				log.Println(err.Error())
				break
			}
			err = h.RedisClient.Publish("datapipe", msg).Err()
			//log.Println("RECV:", msgType, string(msg))
		}
	}(c, h, closeChan)

	// Writer
	go func(c *websocket.Conn, h *Hub, send chan string, close chan int) {
		defer log.Printf("Closing writer for [%s]\n", c.RemoteAddr().String())

		ticker := time.NewTicker(time.Second * 5)
		for {
			select {
			case <-close:
				log.Println("Stopping ticker")
				ticker.Stop()
				return
			case sendMsg := <-send:
				// data := struct {
				// 	Paste            string
				// 	KeepAlive        bool
				// 	BurnAfterReading bool
				// }{
				// 	sendMsg,
				// 	true,
				// 	true,
				// }
				c.WriteJSON(sendMsg)

			case <-ticker.C:
				err := c.WriteMessage(websocket.PingMessage, nil)
				if err != nil {
					log.Println(err.Error())
					return
				}
			}
		}
	}(c, h, cc, closeChan)

}

func (h *Hub) removeCLient(c *websocket.Conn) {
	for _, j := range h.RoomList {
		for i := 0; i < len(j); i++ {
			if c == j[i].conn {
				log.Printf("Removed [%s]\n", c.RemoteAddr().String())
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
		case c := <-h.ClientDisconnect:
			h.removeCLient(c)
			//log.Println("New client:", c.RemoteAddr())
			// case msg := <-h.NewMessage:
			// 	log.Println("New message recv:", msg)
			// 	h.Broadcast(msg)
			// case c := <-h.ClientDisconnect:
			// 	h.removeCLient(c)
		}
	}
}

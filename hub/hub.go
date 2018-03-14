package hub

import (
	"bytes"
	"log"
	"math/rand"
	"net/http"
	"time"

	"github.com/go-redis/redis"
	"github.com/gorilla/websocket"
)

var redisClient *redis.Client

type Message struct {
	MsgType string
	Data    map[string]interface{}
}
type NewConn struct {
	Ws     *websocket.Conn
	Claims map[string]string
}
type Hub struct {
	NewConnection    chan NewConn
	ClientDisconnect chan *websocket.Conn
	RoomList         map[string][]*Client
	RedisClient      *redis.Client
	Router           map[string]func(map[string]interface{})
}

type Client struct {
	conn        *websocket.Conn
	closeChan   chan int
	newMsgChan  chan string
	channelName string
	Claims      map[string]string
}

func NewHub() *Hub {

	// TODO: fail if redis isnt started
	rand.Seed(time.Now().Unix())
	x := &Hub{
		Router:           make(map[string]func(map[string]interface{})),
		NewConnection:    make(chan NewConn),
		ClientDisconnect: make(chan *websocket.Conn),
		RoomList:         make(map[string][]*Client),
		RedisClient:      redisClient,
	}
	log.Println("Redis:", x.RedisClient.Ping().String())
	x.listenForNewMessages()
	return x
}
func init() {
	redisClient = redis.NewClient(&redis.Options{
		Addr:     "redis:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
}

func (h *Hub) GetChannels() {
	channelList := h.RedisClient.PubSubChannels("*")
	for _, j := range channelList.Val() {
		_ = j
	}
}

func (h *Hub) listenForNewMessages() {
	go func() {
		defer h.RedisClient.Close()
		pSub := h.RedisClient.Subscribe("datapipe")
		for {
			msg, err := pSub.ReceiveMessage()
			if err != nil {
				log.Println(err.Error())
			}

			for _, j := range h.RoomList {
				for i := 0; i < len(j); i++ {
					j[i].newMsgChan <- msg.Payload
				}
			}
		}
	}()
}

func sendDBLog(token string) {
	payload := []byte(`{"name":"Dealer Z","destinationDnis":"4151112222"}`)
	req, _ := http.NewRequest("POST", "https://barrenschat-27212.firebaseio.com/message_list.json", bytes.NewBuffer(payload))

	q := req.URL.Query()
	q.Add("auth", token)
	req.URL.RawQuery = q.Encode()
	log.Println(req.URL.String())
	cc := &http.Client{}
	res, e := cc.Do(req)
	if e != nil {
		log.Println(e.Error())
	} else {
		log.Println(res.Body)
	}
}

func (h *Hub) newClientConnection(c NewConn) {
	log.Printf("New client connection [%s]\n", c.Ws.RemoteAddr().String())

	newClient := &Client{
		conn:        c.Ws,
		channelName: "main",
		newMsgChan:  make(chan string),
		closeChan:   make(chan int),
		Claims:      c.Claims,
	}
	log.Println(newClient.Claims["user_id"])
	h.RoomList["main"] = append(h.RoomList["main"], newClient)

	//Reader
	go func(c *websocket.Conn, h *Hub, client *Client) {

		defer func() {
			log.Printf("Closing reader for [%s]\n", c.RemoteAddr().String())
			c.Close()
			client.closeChan <- 1
			h.ClientDisconnect <- c
		}()

		for {
			mmsg := Message{}
			err := c.ReadJSON(&mmsg)

			if err != nil {
				log.Println(err.Error())
				break
			}
			if handler, found := h.findHandler(mmsg.MsgType); found {
				handler(mmsg.Data)
			}
		}
	}(c.Ws, h, newClient)

	// Writer
	go func(c *websocket.Conn, h *Hub, client *Client) {
		defer log.Printf("Closing writer for [%s]\n", c.RemoteAddr().String())

		ticker := time.NewTicker(time.Second * 5)
		for {
			select {
			case <-client.closeChan:
				log.Println("Stopping ticker")
				ticker.Stop()
				return
			case sendMsg := <-client.newMsgChan:
				packet := struct {
					Txt string
				}{
					sendMsg,
				}
				c.WriteJSON(packet)
			case <-ticker.C:
				err := c.WriteMessage(websocket.PingMessage, nil)
				if err != nil {
					log.Println(err.Error())
					return
				}
			}
		}
	}(c.Ws, h, newClient)
}

func (h *Hub) removeCLient(c *websocket.Conn) {
	for _, j := range h.RoomList {
		for i := 0; i < len(j); i++ {
			if c == j[i].conn {
				close(j[i].closeChan)
				close(j[i].newMsgChan)
				log.Printf("Removed [%s]\n", c.RemoteAddr().String())
				h.RoomList[j[i].channelName] = append(h.RoomList[j[i].channelName][:i], h.RoomList[j[i].channelName][i+1:]...)
				return
			}
		}
	}
}

func (h *Hub) Run() {

	// h.HandleMsg("client_info", h.handleClientInfo)
	for {
		select {
		case c := <-h.NewConnection:
			h.newClientConnection(c)
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

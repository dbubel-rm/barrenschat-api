package handler

import (
	"log"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/engineerbeard/barrenschat-api/hub"
	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/assert"
)

type Message struct {
	MsgType string      `json:"msgType"`
	Data    interface{} `json:"data"`
}

func init() {
	//log.SetOutput(ioutil.Discard)
}

// func TestNewConnection(t *testing.T) {
// 	h := hub.NewHub()
// 	go h.Run()

// 	testEngine := GetEngine(h)
// 	s := httptest.NewServer(testEngine)
// 	defer s.Close()

// 	u := "ws" + strings.TrimPrefix(s.URL, "http")
// 	ws, _, err := websocket.DefaultDialer.Dial(u, nil)
// 	defer ws.Close()

// 	assert.NoError(t, err)

// 	data := struct {
// 		Paste            bool
// 		KeepAlive        bool
// 		BurnAfterReading bool
// 	}{
// 		true,
// 		true,
// 		true,
// 	}

// 	var m = Message{MsgType: "newconnection", Data: data}

// 	err = ws.WriteJSON(m)
// 	assert.NoError(t, err)
// }

func TestPingPong(t *testing.T) {
	h := hub.NewHub()
	go h.Run()

	testEngine := GetEngine(h)
	s := httptest.NewServer(testEngine)
	defer s.Close()

	u := "ws" + strings.TrimPrefix("http://localhost:9000/bchatws", "http")
	ws, _, err := websocket.DefaultDialer.Dial(u, nil)
	defer ws.Close()

	assert.NoError(t, err)

	ws.SetPingHandler(func(string) error {
		time.Sleep(time.Second * 1)
		ws.SetWriteDeadline(time.Now().Add(time.Second * 10))
		ws.SetReadDeadline(time.Now().Add(time.Second * 10))
		log.Println("PING rec")
		err := ws.WriteMessage(websocket.PongMessage, nil)
		if err != nil {
			log.Println(err.Error())
		}
		return nil
	})

	log.Println("first ping")
	err = ws.WriteMessage(websocket.PingMessage, nil)
	if err != nil {
		log.Println(err.Error())
	}
	time.Sleep(time.Second * 2)

	msgType, msg, err := ws.ReadMessage()
	_ = msgType
	_ = msg
	// i := 0
	// for i < 10 {
	// 	log.Println(i)
	// 	i++
	// 	time.Sleep(time.Second)
	// 	msgType, msg, err := ws.ReadMessage()
	// 	_ = msgType
	// 	_ = msg

	// 	if err != nil {
	// 		log.Println(err.Error())
	// 		break
	// 	}
	// 	//err = h.RedisClient.Publish("datapipe", msg).Err()
	// 	log.Println("MSG RECV:", msgType, string(msg))
	// }
}

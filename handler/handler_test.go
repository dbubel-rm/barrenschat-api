package handler

import (
	"net/http/httptest"
	"strings"
	"testing"

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

func TestConnect(t *testing.T) {
	h := hub.NewHub()
	go h.Run()

	testEngine := GetEngine(h)
	s := httptest.NewServer(testEngine)
	defer s.Close()

	u := "ws" + strings.TrimPrefix(s.URL, "http")
	ws, _, err := websocket.DefaultDialer.Dial(u, nil)
	// defer ws.Close()

	assert.NoError(t, err)
	d := struct {
		data string
	}{
		data: "HI",
	}
	msg := Message{MsgType: "new_connection", Data: d}
	err = ws.WriteJSON(msg)

	// msgType, msg, err := ws.ReadMessage()
	// _ = msgType
	// _ = msg
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

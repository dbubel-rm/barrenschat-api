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

func fakeAuth(s string) (map[string]string, error) {
	var c map[string]string
	c = make(map[string]string)
	c["user_id"] = "test"
	return c, nil
}

func TestConnectBadAuth(t *testing.T) {
	h := hub.NewHub()
	go h.Run()

	testEngine := GetEngine(h, fakeAuth)
	s := httptest.NewServer(testEngine)
	defer s.Close()

	u := "ws" + strings.TrimPrefix(s.URL, "http")
	ws, res, err := websocket.DefaultDialer.Dial(u, nil)
	_ = res
	defer ws.Close()
}
func TestConnect(t *testing.T) {
	h := hub.NewHub()
	go h.Run()

	testEngine := GetEngine(h, fakeAuth)
	s := httptest.NewServer(testEngine)
	defer s.Close()

	u := "ws" + strings.TrimPrefix(s.URL, "http")
	ws, res, err := websocket.DefaultDialer.Dial(u, nil)
	_ = res
	defer ws.Close()
	// defer ws.Close()

	assert.NoError(t, err)
	// d := struct {
	// 	data string
	// }{
	// 	data: "HI",
	// }
	// msg := Message{MsgType: "new_connection", Data: d}
	// err = ws.WriteJSON(msg)

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

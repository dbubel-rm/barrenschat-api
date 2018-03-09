package handler

import (
	"io/ioutil"
	"log"
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
	log.SetOutput(ioutil.Discard)
}
func TestNewConnection(t *testing.T) {

	h := hub.NewHub()
	go h.Run()

	testEngine := GetEngine(h)

	s := httptest.NewServer(testEngine)

	defer s.Close()

	// Convert http://127.0.0.1 to ws://127.0.0.
	u := "ws" + strings.TrimPrefix(s.URL, "http")

	// Connect to the server
	ws, _, err := websocket.DefaultDialer.Dial(u, nil)
	if err != nil {
		t.Fatalf("%v", err)
	}
	defer ws.Close()

	data := struct {
		Paste            bool
		KeepAlive        bool
		BurnAfterReading bool
	}{
		true,
		true,
		true,
	}

	var m = Message{MsgType: "newconnection", Data: data}

	if err := ws.WriteJSON(m); err != nil {
		assert.Error(t, err)
	}
	for i := 0; i < 20; i++ {
		var z = Message{MsgType: "newmessage", Data: data}
		if err := ws.WriteJSON(z); err != nil {
			assert.NoError(t, err)
		}

	}

	// err = c.Close()
	// assert.NoError(t, err)

	// objmap := make(map[string]string)
	// err = json.Unmarshal(resp.Body.Bytes(), &objmap)
	// assert.NoError(t, err)

}

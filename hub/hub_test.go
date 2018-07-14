package hub

import (
	"log"
	"net/http/httptest"
	"strings"
	"testing"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/assert"
)

func init() {
	//log.SetOutput(ioutil.Discard)
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
}

func fakeAuth(s string) (jwt.MapClaims, error) {
	var c jwt.MapClaims
	c = make(jwt.MapClaims)
	c["user_id"] = "test"
	return c, nil
}

func TestSendMessage(t *testing.T) {
	mockHub := NewHub()
	go mockHub.Run()
	mockMux := GetMux(mockHub, fakeAuth)
	mockServer := httptest.NewServer(mockMux)
	mockURL := "ws" + strings.TrimPrefix(mockServer.URL, "http")
	mockWebsocket, res, err := websocket.DefaultDialer.Dial(mockURL, nil)
	_ = res
	_ = err
	_ = mockWebsocket
	//defer mockWebsocket.Close()

	assert.NoError(t, err)
	assert.Equal(t, 101, res.StatusCode)

	// payload := make(map[string]interface{})

	// payload["message_text"] = "Sample message text"
	// payload["channel"] = "main"
	// m := rawMessage{MsgType: "message_new", Payload: payload}

	// mockWebsocket.WriteJSON(m)
	// _, recvJSON, _ := mockWebsocket.ReadMessage()

	// log.Println(string(recvJSON))
	// expectedJSON := `{"msgType":"message_new","payload":{"channel":"main","message_text":"Sample message text"}}`
	// assert.JSONEq(t, expectedJSON, string(recvJSON))
}

// func BenchmarkSendMessage(b *testing.B) {
// 	mockHub := NewHub()
// 	go mockHub.Run()
// 	mockMux := GetMux(mockHub, fakeAuth)
// 	mockServer := httptest.NewServer(mockMux)
// 	fmt.Println(mockServer.URL)
// 	mockURL := "ws" + strings.TrimPrefix(mockServer.URL, "http")
// 	mockWs, _, _ := websocket.DefaultDialer.Dial(mockURL, nil)
// 	defer mockWs.Close()

// 	payload := make(map[string]interface{})
// 	payload["message_text"] = "Sample message text"
// 	payload["channel"] = "main"
// 	m := rawMessage{MsgType: "message_new", Payload: payload}

// 	log.SetOutput(ioutil.Discard)
// 	for n := 0; n < b.N; n++ {
// 		mockWs.WriteJSON(m)
// 		_, recvJSON, _ := mockWs.ReadMessage()
// 		expectedJSON := `{"msgType":"message_new","payload":{"channel":"main","message_text":"Sample message text"}}`
// 		assert.JSONEq(b, expectedJSON, string(recvJSON))
// 	}
// 	log.SetOutput(os.Stdout)
// }

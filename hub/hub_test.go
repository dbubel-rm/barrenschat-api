package hub

import (
	"fmt"
	"log"
	"net/http/httptest"
	"strings"
	"sync"
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

// func TestSendMessage(t *testing.T) {
// 	mockHub := NewHub()
// 	go mockHub.Run()
// 	mockMux := GetMux(mockHub, fakeAuth)
// 	mockServer := httptest.NewServer(mockMux)
// 	mockURL := "ws" + strings.TrimPrefix(mockServer.URL, "http")
// 	mockWebsocket, res, err := websocket.DefaultDialer.Dial(mockURL, nil)

// 	assert.NoError(t, err)
// 	assert.Equal(t, 101, res.StatusCode)

// 	payload := make(map[string]interface{})

// 	payload[MESSAGE_TEXT] = "Sample message text"
// 	payload["channel"] = "main"
// 	m := rawMessage{MsgType: MESSAGE_TYPE_NEW, Payload: payload}

// 	mockWebsocket.WriteJSON(m)
// 	_, recvJSON, _ := mockWebsocket.ReadMessage()

// 	expectedJSON := `{"msgType":"message_new","payload":{"channel":"main","message_text":"Sample message text"}}`
// 	assert.JSONEq(t, expectedJSON, string(recvJSON))
// }

func TestSendMessageDiffferentChannels(t *testing.T) {
	mockHub := NewHub()
	go mockHub.Run()
	mockMux := GetMux(mockHub, fakeAuth)
	mockServer := httptest.NewServer(mockMux)
	mockURL := "ws" + strings.TrimPrefix(mockServer.URL, "http")
	mockWebsocket, res, err := websocket.DefaultDialer.Dial(mockURL, nil)

	assert.NoError(t, err)
	assert.Equal(t, 101, res.StatusCode)

	payload := make(map[string]interface{})

	payload["message_text"] = "Sample message text A"
	payload["channel"] = "main"
	m := rawMessage{MsgType: MESSAGE_TYPE_NEW, Payload: payload}

	// recvJSON := rawMessage{}
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		_, b, _ := mockWebsocket.ReadMessage()
		expectedJSON := fmt.Sprintf(`{"msgType":"%s","payload":{"channel":"main","message_text":"Sample message text A"}}`, MESSAGE_TYPE_NEW)
		log.Println(string(b))
		assert.JSONEq(t, expectedJSON, string(b))
		wg.Done()
	}()
	mockWebsocket.WriteJSON(m)
	wg.Wait()

	payload["message_text"] = "Sample message text B"
	payload["channel"] = "cool"
	m = rawMessage{MsgType: MESSAGE_TYPE_NEW, Payload: payload}

	wg.Add(1)
	go func() {
		_, b, _ := mockWebsocket.ReadMessage()
		expectedJSON := fmt.Sprintf(`{"msgType":"%s","payload":{"channel":"cool","message_text":"Sample message text B"}}`, MESSAGE_TYPE_NEW)
		log.Println(string(b))
		assert.JSONEq(t, expectedJSON, string(b))
		wg.Done()
	}()
	mockWebsocket.WriteJSON(m)
	wg.Wait()

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

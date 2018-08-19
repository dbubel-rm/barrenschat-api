package hub

import (
	"log"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

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
func setupTestConn() (*websocket.Conn, error) {
	mockHub := NewHub()
	go mockHub.Run()
	mockMux := GetMux(mockHub, fakeAuth)
	mockServer := httptest.NewServer(mockMux)
	mockURL := "ws" + strings.TrimPrefix(mockServer.URL, "http")
	mockWebsocket, _, err := websocket.DefaultDialer.Dial(mockURL, nil)
	mockWebsocket.SetReadDeadline(time.Now().Add(time.Second * 3))
	mockWebsocket.SetWriteDeadline(time.Now().Add(time.Second))

	return mockWebsocket, err
}

func TestSendMessage(t *testing.T) {
	mockWebsocket, err := setupTestConn()
	payload := make(map[string]interface{})
	payload[MESSAGE_TEXT] = "Sample message text"
	payload["channel"] = "main"
	m := rawMessage{MsgType: MESSAGE_TYPE_NEW, Payload: payload}

	var mm rawMessage
	mockWebsocket.WriteJSON(m)
	time.Sleep(time.Millisecond * 10)
	err = mockWebsocket.ReadJSON(&mm)

	assert.NoError(t, err)

	// v := make(map[string]interface{})
	// v["channel"] = "main"
	// v["message_text"] = "Sample message text"
	// a := rawMessage{MsgType: "message_new", Payload: v}
	// assert.Equal(t, a, mm)

}

// func TestCreateChannel(t *testing.T) {
// 	mockWebsocket, err := setupTestConn()
// 	assert.NoError(t, err)
// 	payload := make(map[string]interface{})
// 	// payload[MESSAGE_TEXT] = "Sample message text"
// 	payload["channel_name"] = "NewChannel"
// 	m := rawMessage{MsgType: MESSAGE_TYPE_NEW_CHANNEL, Payload: payload}

// 	var mm rawMessage
// 	mockWebsocket.WriteJSON(m)
// 	time.Sleep(time.Millisecond * 10)
// 	err = mockWebsocket.ReadJSON(&mm)

// 	time.Sleep(time.Second)
// 	fmt.Println(mock)

// 	// assert.NoError(t, err)

// 	// v := make(map[string]interface{})
// 	// v["channel"] = "main"
// 	// v["message_text"] = "Sample message text"
// 	// a := rawMessage{MsgType: "message_new", Payload: v}
// 	// assert.Equal(t, a, mm)

// }

// func BenchmarkSendMessage(b *testing.B) {
// 	log.SetOutput(ioutil.Discard)
// 	mockHub := NewHub()
// 	go mockHub.Run()
// 	mockMux := GetMux(mockHub, fakeAuth)
// 	mockServer := httptest.NewServer(mockMux)
// 	mockURL := "ws" + strings.TrimPrefix(mockServer.URL, "http")
// 	mockWebsocket, res, err := websocket.DefaultDialer.Dial(mockURL, nil)

// 	assert.NoError(b, err)
// 	assert.Equal(b, 101, res.StatusCode)

// 	payload := make(map[string]interface{})

// 	payload[MESSAGE_TEXT] = "Sample message text"
// 	payload["channel"] = "main"
// 	m := rawMessage{MsgType: MESSAGE_TYPE_NEW, Payload: payload}
// 	mockWebsocket.SetReadDeadline(time.Now().Add(time.Second * 3))
// 	mockWebsocket.SetWriteDeadline(time.Now().Add(time.Second))
// 	for n := 0; n < b.N; n++ {
// 		var mm rawMessage
// 		mockWebsocket.WriteJSON(m)
// 		err := mockWebsocket.ReadJSON(&mm)
// 		assert.NoError(b, err)

// 		v := make(map[string]interface{})
// 		v["channel"] = "main"
// 		v["message_text"] = "Sample message text"
// 		t := rawMessage{MsgType: "message_new", Payload: v}

// 		// expectedJSON := `{"msgType":"message_new","payload":{"channel":"main","message_text":"Sample message text"}}`
// 		assert.Equal(b, t, mm)
// 	}
// 	log.SetOutput(os.Stdout)
// }

// func TestCreateNewChannel(t *testing.T) {
// 	mockHub := NewHub()
// 	go mockHub.Run()
// 	mockMux := GetMux(mockHub, fakeAuth)
// 	mockServer := httptest.NewServer(mockMux)
// 	mockURL := "ws" + strings.TrimPrefix(mockServer.URL, "http")
// 	mockWebsocket, res, err := websocket.DefaultDialer.Dial(mockURL, nil)

// 	assert.NoError(t, err)
// 	assert.Equal(t, 101, res.StatusCode)

// 	payload := make(map[string]interface{})

// 	// payload["message_text"] = "Sample message text A"
// 	payload["channel_name"] = "NewChannelName"
// 	m := rawMessage{MsgType: MESSAGE_TYPE_NEW_CHANNEL, Payload: payload}
// 	mockWebsocket.WriteJSON(m)
// 	// fmt.Println(mockHub.channels)
// 	// fmt.Println(mockHub.channels)
// 	// _, b, _ := mockWebsocket.ReadMessage()
// 	// fmt.Println(string(b))

// }

// func TestSendMessageDiffferentChannels(t *testing.T) {
// 	mockHub := NewHub()
// 	go mockHub.Run()
// 	mockMux := GetMux(mockHub, fakeAuth)
// 	mockServer := httptest.NewServer(mockMux)
// 	mockURL := "ws" + strings.TrimPrefix(mockServer.URL, "http")
// 	mockWebsocket, res, err := websocket.DefaultDialer.Dial(mockURL, nil)

// 	assert.NoError(t, err)
// 	assert.Equal(t, 101, res.StatusCode)

// 	payload := make(map[string]interface{})

// 	payload["message_text"] = "Sample message text A"
// 	payload["channel"] = "main"
// 	m := rawMessage{MsgType: MESSAGE_TYPE_NEW, Payload: payload}

// 	// recvJSON := rawMessage{}
// 	var wg sync.WaitGroup
// 	wg.Add(1)
// 	go func() {
// 		_, b, _ := mockWebsocket.ReadMessage()
// 		expectedJSON := fmt.Sprintf(`{"msgType":"%s","payload":{"channel":"main","message_text":"Sample message text A"}}`, MESSAGE_TYPE_NEW)
// 		log.Println(string(b))

// 		assert.JSONEq(t, expectedJSON, string(b))
// 		wg.Done()
// 	}()
// 	mockWebsocket.WriteJSON(m)
// 	wg.Wait()

// 	// payload["message_text"] = "Sample message text B"
// 	// payload["channel"] = "cool"
// 	// m = rawMessage{MsgType: MESSAGE_TYPE_NEW, Payload: payload}

// 	// wg.Add(1)
// 	// go func() {
// 	// 	_, b, _ := mockWebsocket.ReadMessage()
// 	// 	expectedJSON := fmt.Sprintf(`{"msgType":"%s","payload":{"channel":"cool","message_text":"Sample message text B"}}`, MESSAGE_TYPE_NEW)
// 	// 	log.Println(string(b))
// 	// 	assert.JSONEq(t, expectedJSON, string(b))
// 	// 	wg.Done()
// 	// }()
// 	// mockWebsocket.WriteJSON(m)
// 	// wg.Wait()

// }

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

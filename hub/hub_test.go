package hub

// import (
// 	"fmt"
// 	"io/ioutil"
// 	"log"
// 	"net/http/httptest"
// 	"strings"
// 	"testing"

// 	"github.com/gorilla/websocket"
// 	"github.com/stretchr/testify/assert"
// )

// func init() {
// 	log.SetOutput(ioutil.Discard)
// 	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
// }

// func fakeAuth(s string) (map[string]string, error) {
// 	var c map[string]string
// 	c = make(map[string]string)
// 	c["user_id"] = "test"
// 	return c, nil
// }
// func TestSendMessage(t *testing.T) {
// 	mockHub := NewHub()
// 	go mockHub.Run()
// 	mockMux := GetMux(mockHub, fakeAuth)
// 	mockServer := httptest.NewServer(mockMux)
// 	mockURL := "ws" + strings.TrimPrefix(mockServer.URL, "http")
// 	mockWebsocket, res, err := websocket.DefaultDialer.Dial(mockURL, nil)
// 	defer mockWebsocket.Close()

// 	assert.NoError(t, err)
// 	assert.Equal(t, 101, res.StatusCode)

// 	payload := make(map[string]interface{})
// 	payload[MESSAGE_USER] = "test"
// 	payload[MESSAGE_TEXT] = "Sample message text"
// 	m := message{MsgType: TYPE_NEW_MESSAGE, Payload: payload}

// 	mockWebsocket.WriteJSON(m)
// 	_, recvJSON, _ := mockWebsocket.ReadMessage()

// 	expectedJSON := `{"msgType":"message_new","payload":{"message_text":"Sample message text"}}`
// 	assert.JSONEq(t, expectedJSON, string(recvJSON))

// 	payload = make(map[string]interface{})
// 	payload[MESSAGE_USER] = "test"
// 	payload[MESSAGE_TEXT] = "Sample Message two"
// 	m = message{MsgType: TYPE_NEW_MESSAGE, Payload: payload}

// 	mockWebsocket.WriteJSON(m)
// 	_, recvJSON, _ = mockWebsocket.ReadMessage()

// 	expectedJSON = `{"msgType":"message_new","payload":{"message_text":"Sample Message two"}}`
// 	assert.JSONEq(t, expectedJSON, string(recvJSON))
// }

// // func getClient() *websocket.Conn {
// // 	mockHub := NewHub()
// // 	go mockHub.Run()
// // 	mockMux := GetMux(mockHub, fakeAuth)
// // 	mockServer := httptest.NewServer(mockMux)
// // 	mockURL := "ws" + strings.TrimPrefix(mockServer.URL, "http")
// // 	mockWs, _, err := websocket.DefaultDialer.Dial(mockURL, nil)
// // 	if err != nil {
// // 		log.Println(err.Error())
// // 		return nil
// // 	}
// // 	return mockWs
// // }
// func BenchmarkMessages(b *testing.B) {
// 	mockHub := NewHub()
// 	go mockHub.Run()
// 	mockMux := GetMux(mockHub, fakeAuth)
// 	mockServer := httptest.NewServer(mockMux)
// 	fmt.Println(mockServer.URL)
// 	mockURL := "ws" + strings.TrimPrefix(mockServer.URL, "http")
// 	mockWs, _, _ := websocket.DefaultDialer.Dial(mockURL, nil)
// 	defer mockWs.Close()

// 	payload := make(map[string]interface{})
// 	payload[MESSAGE_USER] = "test"
// 	payload[MESSAGE_TEXT] = "Sample message text"
// 	m := message{MsgType: TYPE_NEW_MESSAGE, Payload: payload}

// 	//log.SetOutput(ioutil.Discard)
// 	for n := 0; n < b.N; n++ {
// 		mockWs.WriteJSON(m)
// 		_, recvJSON, _ := mockWs.ReadMessage()
// 		expectedJSON := `{"msgType":"message_new","payload":{"message_text":"Sample message text"}}`
// 		assert.JSONEq(b, expectedJSON, string(recvJSON))
// 	}
// 	//log.SetOutput(os.Stdout)
// }

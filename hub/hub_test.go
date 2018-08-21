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
	c["localId"] = "12345qwerty"
	return c, nil
}
func setupTestConn() (*websocket.Conn, error) {
	mockHub := NewHub()
	go mockHub.Run()
	mockMux := GetMux(mockHub, fakeAuth)
	mockServer := httptest.NewServer(mockMux)
	mockURL := "ws" + strings.TrimPrefix(mockServer.URL, "http")
	mockWebsocket, _, err := websocket.DefaultDialer.Dial(mockURL, nil)
	mockWebsocket.SetReadDeadline(time.Now().Add(time.Second * time.Duration(3)))
	mockWebsocket.SetWriteDeadline(time.Now().Add(time.Second * time.Duration(10)))

	return mockWebsocket, err
}

func TestSendMessage(t *testing.T) {
	mockWebsocket, err := setupTestConn()
	payload := make(map[string]interface{})
	payload[MessageText] = "Sample message text"
	payload["channel"] = "main"
	m := rawMessage{MsgType: MessageTypeNew, Payload: payload}

	var mm rawMessage
	mockWebsocket.WriteJSON(m)
	err = mockWebsocket.ReadJSON(&mm)

	assert.NoError(t, err)

	v := make(map[string]interface{})
	v["channel"] = "main"
	v["message_text"] = "Sample message text"
	a := rawMessage{MsgType: "message_new", Payload: v}
	assert.Equal(t, a, mm)

}

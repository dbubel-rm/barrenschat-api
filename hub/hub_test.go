package hub

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http/httptest"
	"os"
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
	c["localId"] = RandStringBytes(32)

	return c, nil
}

func RandStringBytes(n int) string {
	const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}

func setupTestConn() (*Hub, *websocket.Conn, error) {
	mockHub := NewHub()
	go mockHub.Run()
	mockMux := GetMux(mockHub, fakeAuth)
	mockServer := httptest.NewServer(mockMux)
	mockURL := "ws" + strings.TrimPrefix(mockServer.URL, "http")
	mockWebsocket, _, err := websocket.DefaultDialer.Dial(mockURL, nil)
	mockWebsocket.SetReadDeadline(time.Now().Add(time.Second * time.Duration(3)))
	mockWebsocket.SetWriteDeadline(time.Now().Add(time.Second * time.Duration(3)))

	return mockHub, mockWebsocket, err
}

func TestClientConnect(t *testing.T) {

	mockHub := NewHub()
	go mockHub.Run()
	mockMux := GetMux(mockHub, fakeAuth)
	mockServer := httptest.NewServer(mockMux)
	mockURL := "ws" + strings.TrimPrefix(mockServer.URL, "http")

	mockWebsocket, _, err := websocket.DefaultDialer.Dial(mockURL, nil)
	assert.NoError(t, err)
	mockWebsocket.SetReadDeadline(time.Now().Add(time.Second * time.Duration(3)))
	mockWebsocket.SetWriteDeadline(time.Now().Add(time.Second * time.Duration(3)))
	time.Sleep(time.Millisecond * 200)
	assert.Equal(t, 1, len(mockHub.getClients()))
}

func BenchmarkClientConnect(b *testing.B) {
	mockHub := NewHub()
	go mockHub.Run()
	mockMux := GetMux(mockHub, fakeAuth)
	mockServer := httptest.NewServer(mockMux)
	mockURL := "ws" + strings.TrimPrefix(mockServer.URL, "http")

	for i := 0; i < b.N; i++ {
		websocket.DefaultDialer.Dial(mockURL, nil)
	}
}

func BenchmarkUnmarshalJSON(b *testing.B) {
	type Pool struct {
		PoolID      int        `gorm:"primary_key" json:"poolId"`
		Name        string     `json:"name" binding:"required"`
		ClientKey   string     `json:"clientKey"` // This key is used by the JavaScript client to lease numbers from the pool
		LogoURL     string     `json:"logoUrl"`
		RVCompanyID int        `json:"rvCompanyId"`
		TenantID    string     `json:"tenantId"`
		CreatedAt   time.Time  `json:"createdAt"`
		UpdatedAt   time.Time  `json:"updatedAt"`
		DeletedAt   *time.Time `json:"-"`
	}
	var p Pool
	jsons := `
		{
			"poolId": 1,
			"name": "Verizon FiOS",
			"clientKey": "hxgy46G4P6FN3MMNPkesV3",
			"logoUrl": "fivz_logo.png",
			"rvCompanyId": 54,
			"tenantId": "tenantId1",
			"createdAt": "2000-01-01T00:00:00Z",
			"updatedAt": "2000-01-01T00:00:00Z",
			"phoneNumbersTotal": 3,
			"phoneNumbersAvailable": 3,
			"defaultMarketingCode": null,
			"permaleases": 0
		}
	`
	for i := 0; i < b.N; i++ {
		json.Unmarshal([]byte(jsons), &p)
	}
}
func BenchmarkDecodeJSON(b *testing.B) {
	type Pool struct {
		PoolID      int        `gorm:"primary_key" json:"poolId"`
		Name        string     `json:"name" binding:"required"`
		ClientKey   string     `json:"clientKey"` // This key is used by the JavaScript client to lease numbers from the pool
		LogoURL     string     `json:"logoUrl"`
		RVCompanyID int        `json:"rvCompanyId"`
		TenantID    string     `json:"tenantId"`
		CreatedAt   time.Time  `json:"createdAt"`
		UpdatedAt   time.Time  `json:"updatedAt"`
		DeletedAt   *time.Time `json:"-"`
	}
	var p Pool
	jsons := `
		{
			"poolId": 1,
			"name": "Verizon FiOS",
			"clientKey": "hxgy46G4P6FN3MMNPkesV3",
			"logoUrl": "fivz_logo.png",
			"rvCompanyId": 54,
			"tenantId": "tenantId1",
			"createdAt": "2000-01-01T00:00:00Z",
			"updatedAt": "2000-01-01T00:00:00Z",
			"phoneNumbersTotal": 3,
			"phoneNumbersAvailable": 3,
			"defaultMarketingCode": null,
			"permaleases": 0
		}
	`
	for i := 0; i < b.N; i++ {
		dec := json.NewDecoder(bufio.NewReader(bytes.NewBuffer([]byte(jsons))))
		dec.Decode(&p)
	}

}

func TestSendMessage(t *testing.T) {
	_, mockWebsocket, err := setupTestConn()
	assert.NoError(t, err)
	payload := make(map[string]interface{})
	payload[MessageText] = "Sample message text"
	payload["channel"] = "main"
	m := rawMessage{MsgType: MessageTypeChat, Payload: payload}

	var mm rawMessage
	mockWebsocket.WriteJSON(m)
	err = mockWebsocket.ReadJSON(&mm)

	assert.NoError(t, err)

	v := make(map[string]interface{})
	v["channel"] = "main"
	v["message_text"] = "Sample message text"
	a := rawMessage{MsgType: "message_new", Payload: v, ClientID: "test"}
	assert.Equal(t, a, mm)

}
func BenchmarkSendMessage(b *testing.B) {

	log.SetOutput(ioutil.Discard)
	_, mockWebsocket, err := setupTestConn()
	assert.NoError(b, err)
	payload := make(map[string]interface{})
	payload[MessageText] = "Sample message text"
	payload["channel"] = "main"
	m := rawMessage{MsgType: MessageTypeChat, Payload: payload}

	var mm rawMessage

	for i := 0; i < b.N; i++ {
		mockWebsocket.WriteJSON(m)
		mockWebsocket.ReadJSON(&mm)
	}
	log.SetOutput(os.Stdout)
}

func TestCreateNewChannel(t *testing.T) {
	mockHub := NewHub()
	go mockHub.Run()
	mockMux := GetMux(mockHub, fakeAuth)
	mockServer := httptest.NewServer(mockMux)
	mockURL := "ws" + strings.TrimPrefix(mockServer.URL, "http")

	mockWebsocket, _, err := websocket.DefaultDialer.Dial(mockURL, nil)
	assert.NoError(t, err)
	mockWebsocket.SetReadDeadline(time.Now().Add(time.Second * time.Duration(3)))
	mockWebsocket.SetWriteDeadline(time.Now().Add(time.Second * time.Duration(3)))

	payload := make(map[string]interface{})

	payload["channel"] = "New Channel"
	m := rawMessage{MsgType: CommandNewChannel, Payload: payload}

	var mm rawMessage

	mockWebsocket.WriteJSON(m)
	mockWebsocket.ReadJSON(&mm)

	assert.Equal(t, 2, len(mockHub.getTopicChannels()))

	payload = make(map[string]interface{})
	payload[MessageText] = "stuff from new channel"
	payload["channel"] = "New Channel"
	m = rawMessage{MsgType: MessageTypeChat, Payload: payload}

	mockWebsocket.WriteJSON(m)
	err = mockWebsocket.ReadJSON(&mm)

}

func BenchmarkCreateNewChannel(b *testing.B) {

	// log.SetOutput(ioutil.Discard)
	mockHub := NewHub()
	go mockHub.Run()
	mockMux := GetMux(mockHub, fakeAuth)
	mockServer := httptest.NewServer(mockMux)
	mockURL := "ws" + strings.TrimPrefix(mockServer.URL, "http")

	mockWebsocket, _, err := websocket.DefaultDialer.Dial(mockURL, nil)
	assert.NoError(b, err)
	mockWebsocket.SetReadDeadline(time.Now().Add(time.Second * time.Duration(3)))
	mockWebsocket.SetWriteDeadline(time.Now().Add(time.Second * time.Duration(3)))

	payload := make(map[string]interface{})

	for i := 0; i < b.N; i++ {
		s := RandStringBytes(25)
		payload["channel"] = s
		m := rawMessage{MsgType: CommandNewChannel, Payload: payload}

		var mm rawMessage

		mockWebsocket.WriteJSON(m)
		mockWebsocket.ReadJSON(&mm)

		// assert.Equal(t, 2, len(mockHub.getTopicChannels()))

		payload = make(map[string]interface{})
		payload[MessageText] = "stuff from new channel"
		payload["channel"] = s
		m = rawMessage{MsgType: MessageTypeChat, Payload: payload}

		mockWebsocket.WriteJSON(m)
		err = mockWebsocket.ReadJSON(&mm)
		fmt.Println(mm)
		// assert.Equal(b, s, mm.Payload["channel"].(string))
	}
	// log.SetOutput(os.Stdout)
}

package hub

import (
	"bufio"
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
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
	c["localId"] = "12345qwerty"
	return c, nil
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
	_, _, err := setupTestConn()
	assert.NoError(t, err)
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
		b.StartTimer()
		json.Unmarshal([]byte(jsons), &p)
		b.StopTimer()
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
		b.StartTimer()
		dec := json.NewDecoder(bufio.NewReader(bytes.NewBuffer([]byte(jsons))))
		dec.Decode(&p)
		b.StopTimer()
	}

}
func TestSendMessage(t *testing.T) {
	_, mockWebsocket, err := setupTestConn()
	assert.NoError(t, err)
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

func BenchmarkSendMessage(b *testing.B) {

	log.SetOutput(ioutil.Discard)
	_, mockWebsocket, err := setupTestConn()
	assert.NoError(b, err)
	payload := make(map[string]interface{})
	payload[MessageText] = "Sample message text"
	payload["channel"] = "main"
	m := rawMessage{MsgType: MessageTypeNew, Payload: payload}

	var mm rawMessage

	// fmt.Println(unsafe.Sizeof(mm))
	for i := 0; i < b.N; i++ {

		b.StartTimer()
		mockWebsocket.WriteJSON(m)
		mockWebsocket.ReadJSON(&mm)
		b.StopTimer()
	}
	log.SetOutput(os.Stdout)
}

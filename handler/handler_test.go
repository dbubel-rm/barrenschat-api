package handler

import (
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"testing"
	"time"

	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/assert"
)

type Message struct {
	MsgType string      `json:"msgType"`
	Data    interface{} `json:"data"`
}

func TestGetVersion(t *testing.T) {

	url := "http://load-balancer/version"

	spaceClient := http.Client{
		Timeout: time.Second * 2, // Maximum of 2 secs
	}

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		log.Fatal(err)
	}

	res, getErr := spaceClient.Do(req)
	if getErr != nil {
		log.Fatal(getErr)
	}

	body1, readErr := ioutil.ReadAll(res.Body)
	if readErr != nil {
		log.Fatal(readErr)
	}

	res, getErr = spaceClient.Do(req)
	if getErr != nil {
		log.Fatal(getErr)
	}

	body, readErr := ioutil.ReadAll(res.Body)
	if readErr != nil {
		log.Fatal(readErr)
	}
	s := string(body) + string(body1)
	assert.Contains(t, s, "1")
	assert.Contains(t, s, "2")

	// objmap := make(map[string]string)
	// err = json.Unmarshal(resp.Body.Bytes(), &objmap)
	// assert.NoError(t, err)

}
func TestNewConnection(t *testing.T) {
	u := url.URL{Scheme: "ws", Host: "load-balancer", Path: "/bchatws"}
	d := websocket.Dialer{}
	c, _, err := d.Dial(u.String(), nil)

	if err != nil {
		assert.NoError(t, err)
	}

	// cm := websocket.FormatCloseMessage(websocket.CloseNormalClosure, "add your message here")
	// if err := c.WriteMessage(websocket.CloseMessage, cm); err != nil {
	// 	fmt.Println(err)
	// }
	// c.Close()

	//cm := websocket.FormatCloseMessage(websocket.CloseNormalClosure, "add your message here")
	data := struct {
		Paste            bool
		KeepAlive        bool
		BurnAfterReading bool
	}{
		true,
		true,
		true,
	}
	for i := 0; i < 100; i++ {
		var m = Message{MsgType: "new connection", Data: data}

		if err := c.WriteJSON(m); err != nil {
			assert.NoError(t, err)
		}
	}

	err = c.Close()
	assert.NoError(t, err)

	// objmap := make(map[string]string)
	// err = json.Unmarshal(resp.Body.Bytes(), &objmap)
	// assert.NoError(t, err)

}

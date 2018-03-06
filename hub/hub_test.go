package hub

// import (
// 	"fmt"
// 	"net/url"
// 	"testing"
// 	"time"

// 	"github.com/gorilla/websocket"
// )

// type msg struct {
// 	MsgType string      `json:"msgType"`
// 	Data    interface{} `json:"data"`
// }

// func TestRedis(t *testing.T) {
// 	h := NewHub()
// 	psub := h.RedisClient.Subscribe("test chan")
// 	//p := h.RedisClient.Subscribe("mychannel1")
// 	fmt.Println(h.RedisClient.PubSubChannels("*"))
// 	defer h.RedisClient.Close()

// 	subscr, err := psub.ReceiveTimeout(time.Second)
// 	if err != nil {
// 		panic(err)
// 	}
// 	fmt.Println(subscr)

// 	err = h.RedisClient.Publish("test chan", "hello").Err()
// }
// func TestNewConnection(t *testing.T) {
// 	u := url.URL{Scheme: "ws", Host: "load-balancer", Path: "/bchatws"}
// 	d := websocket.Dialer{}
// 	c, _, err := d.Dial(u.String(), nil)

// 	if err != nil {
// 		fmt.Errorf(err.Error())
// 	}

// 	// cm := websocket.FormatCloseMessage(websocket.CloseNormalClosure, "add your message here")
// 	// if err := c.WriteMessage(websocket.CloseMessage, cm); err != nil {
// 	// 	fmt.Println(err)
// 	// }
// 	// c.Close()

// 	//cm := websocket.FormatCloseMessage(websocket.CloseNormalClosure, "add your message here")
// 	var m = msg{MsgType: "new connection", Data: "Data goes here"}

// 	if err := c.WriteJSON(m); err != nil {
// 		fmt.Println(err)
// 	}
// 	c.Close()

// 	// objmap := make(map[string]string)
// 	// err = json.Unmarshal(resp.Body.Bytes(), &objmap)
// 	// assert.NoError(t, err)

// }

package handler

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/engineerbeard/barrenschat-api/hub"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

const (
	writeWait      = 10 * time.Second
	pongWait       = 60 * time.Second
	pingPeriod     = (pongWait * 9) / 10
	maxMessageSize = 512
)

func init() {
	f, err := os.OpenFile("hub_log.txt", os.O_APPEND|os.O_CREATE|os.O_RDWR, 0666)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	mw := io.MultiWriter(os.Stdout)
	log.SetOutput(mw)
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
}

func GetEngine(h *hub.Hub) *gin.Engine {
	router := gin.New()

	router.GET("/version", func(c *gin.Context) {
		c.String(http.StatusOK, fmt.Sprint("Barrenschat API OK v", os.Getenv("NAME")))
	})

	router.GET("/", func(c *gin.Context) {
		wshandler(c.Writer, c.Request, h)
	})
	return router
}

func wshandler(w http.ResponseWriter, r *http.Request, h *hub.Hub) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Failed to upgrade ws: ", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	go func(c *websocket.Conn, h *hub.Hub) {
		defer func() {
			log.Printf("Closing connection for [%s]\n", c.RemoteAddr().String())
			//c.hub.MsgRecvr <- hub.Message{MsgType: "client disconnect"}
			c.Close()
		}()
		// c.SetReadLimit(maxMessageSize)
		// c.SetReadDeadline(time.Now().Add(pongWait))
		// c.SetPongHandler(func(string) error { c.SetReadDeadline(time.Now().Add(pongWait)); return nil })
		for {
			msgType, msg, err := c.ReadMessage()
			log.Println(msgType, string(msg))
			if err != nil {
				log.Println(err.Error())
				break
			}
		}
	}(conn, h)
}

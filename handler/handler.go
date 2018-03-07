package handler

import (
	"fmt"
	"net/http"
	"os"

	"github.com/engineerbeard/barrenschat-api/client"
	"github.com/engineerbeard/barrenschat-api/hub"
	logging "github.com/op/go-logging"

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
var log = logging.MustGetLogger("barrenschat-api")

func init() {

	//
	var format = logging.MustStringFormatter(`%{color}%{time:15:04:05} %{shortfile} %{level:.4s} %{id:03x}%{color:reset} %{message}`)
	logBackend := logging.NewLogBackend(os.Stdout, "", 0)
	backendFormatter := logging.NewBackendFormatter(logBackend, format)
	logging.SetBackend(backendFormatter)
	logging.SetLevel(logging.DEBUG, "barrenschat-api")
}

func GetEngine(h *hub.Hub) *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()

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
	_ = conn
	if err != nil {
		fmt.Println("Failed to upgrade ws: ", err)
		return
	}
	go client.NewClient(conn, h)
}

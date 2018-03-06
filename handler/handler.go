package handler

import (
	"fmt"
	"net/http"
	"os"

	"github.com/engineerbeard/barrenschat-api/client"
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

func GetEngine(h *hub.Hub) *gin.Engine {
	gin.SetMode(gin.DebugMode)
	router := gin.Default()

	router.GET("/version", func(c *gin.Context) {
		c.String(http.StatusOK, fmt.Sprint("Barrenschat API OK v", os.Getenv("NAME")))
	})

	router.GET("/", func(c *gin.Context) {
		wshandler(c.Writer, c.Request, h)
	})
	return router
}
func middleMan() gin.HandlerFunc {
	return func(c *gin.Context) {
		fmt.Println("HI")
		return
	}
}

func wshandler(w http.ResponseWriter, r *http.Request, h *hub.Hub) {
	conn, err := upgrader.Upgrade(w, r, nil)

	if err != nil {
		fmt.Println("Failed to upgrade ws: ", err)
		return
	}
	go client.NewClient(conn, h)
}

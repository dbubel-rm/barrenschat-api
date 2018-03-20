package handler

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/engineerbeard/barrenschat-api/hub"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

var connTimeout = 60 * time.Second

func init() {
	if os.Getenv("ENV_NAME") == "test" {
		connTimeout = 2 * time.Second
	}

}

// GetEngine returns router for the API
func GetEngine(h *hub.Hub) *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("/version", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(fmt.Sprint("Barrenschat API OK v", os.Getenv("NAME"))))
	})

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		ws, err := upgrader.Upgrade(w, r, nil)

		if err != nil {
			log.Println(err.Error())
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		ws.SetReadLimit(1024)
		ws.SetPongHandler(func(string) error {
			ws.SetWriteDeadline(time.Now().Add(connTimeout))
			ws.SetReadDeadline(time.Now().Add(connTimeout))
			log.Println("Pong rec")
			return nil
		})
		h.NewConnection <- ws
	})

	return mux
}

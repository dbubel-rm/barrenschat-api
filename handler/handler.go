package handler

import (
	"fmt"
	"io"
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
	f, err := os.OpenFile("hub_log.txt", os.O_APPEND|os.O_CREATE|os.O_RDWR, 0666)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	mw := io.MultiWriter(os.Stdout)
	log.SetOutput(mw)
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
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

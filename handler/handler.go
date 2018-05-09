package handler

import (
	"crypto/rsa"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/engineerbeard/barrenschat-api/hub"
	"github.com/gorilla/websocket"
)

var connTimeout = 60 * time.Second

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024 * 1024,
	WriteBufferSize: 1024 * 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}
var publicKey *rsa.PublicKey

func wsStart(h *hub.Hub, authFunc func(string) (map[string]string, error)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var claims map[string]string
		var err error

		//claims, err = authFunc
		claims, err = authFunc(r.URL.Query().Get("params"))

		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
		}
		ws, err := upgrader.Upgrade(w, r, nil)

		if err != nil {
			log.Println(err.Error())
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		// TODO: Check for duplicate connection
		ws.SetReadLimit(1024 * 1024)
		ws.SetPongHandler(func(string) error {
			ws.SetWriteDeadline(time.Now().Add(connTimeout))
			ws.SetReadDeadline(time.Now().Add(connTimeout))
			log.Println("Pong rec")
			return nil
		})

		h.NewConnection <- struct {
			Ws     *websocket.Conn
			Claims map[string]string
		}{
			ws,
			claims,
		}
	}
}

// GetEngine returns router for the API
func GetEngine(h *hub.Hub, authFunc func(string) (map[string]string, error)) *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("/version", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(fmt.Sprint("Barrenschat API OK:", os.Getenv("NAME"))))
	})

	mux.Handle("/", wsStart(h, authFunc))
	return mux
}

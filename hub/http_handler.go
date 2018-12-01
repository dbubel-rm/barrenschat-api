package hub

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"runtime"

	"github.com/dgrijalva/jwt-go"

	"github.com/gorilla/websocket"
)

func bToMb(b uint64) uint64 {
	return b / 1024 / 1024
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func wsStart(h *Hub, authUser func(string) (jwt.MapClaims, error)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			http.Error(w, "Not allowed", http.StatusMethodNotAllowed)
			return
		}
		var err error
		var ws *websocket.Conn

		// Grab jwt from query param
		var claimsMap jwt.MapClaims
		claimsMap, err = authUser(r.URL.Query().Get("params"))

		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		ws, err = upgrader.Upgrade(w, r, nil)

		if err != nil {
			log.Println(err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		channels := []string{"main"}
		client := &Client{
			Hub:                  h,
			conn:                 ws,
			send:                 make(chan []byte),
			channelsSubscribedTo: channels,
			claims:               claimsMap,
			ID:                   claimsMap["user_id"].(string),
			locker:               make(chan bool, 1),
		}
		client.Hub.clientConnect <- client
		go client.writeWorker()
		go client.readWorker()
	}
}

func health(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Not allowed", http.StatusMethodNotAllowed)
		return
	}

	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprint(
		"Barrenschat API: OK", os.Getenv("NAME"), "\n",
		"Total Allocated:", bToMb(m.Alloc), "M\n",
		"Total Sys:", bToMb(m.Sys), "M\n",
		"Total Allocations:", bToMb(m.TotalAlloc), "\n",
		"Live Objects:", m.Mallocs-m.Frees,
	)))
}

// GetMux returns router for the API
func GetMux(h *Hub, authFunc func(string) (jwt.MapClaims, error)) *http.ServeMux {
	mux := http.NewServeMux()

	mux.Handle("/", wsStart(h, authFunc))
	mux.HandleFunc("/version", health)

	return mux
}

package handler

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"runtime"
	"time"

	"github.com/engineerbeard/barrenschat-api/hub"
	"github.com/engineerbeard/barrenschat-api/middleware"
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

func wsStart(h *hub.Hub, authUser func(string) (map[string]string, error)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var claims map[string]string
		var err error
		var ws *websocket.Conn

		// Grab jwt from query param
		claims, err = authUser(r.URL.Query().Get("params"))

		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		ok := middleware.ValidateClaims(claims)

		if ok {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		// jwt signature and claims ok so upgrade user to websocket
		ws, err = upgrader.Upgrade(w, r, nil)

		if err != nil {
			log.Println(err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		ws.SetReadLimit(1024 * 1024)
		ws.SetPongHandler(func(string) error {
			ws.SetWriteDeadline(time.Now().Add(connTimeout))
			ws.SetReadDeadline(time.Now().Add(connTimeout))
			log.Println("Pong rec")
			return nil
		})

		// TODO: Check for duplicate connection
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

		var m runtime.MemStats
		runtime.ReadMemStats(&m)
		// For info on each, see: https://golang.org/pkg/runtime/#MemStats
		stats := fmt.Sprintf("\tAlloc = %v MiB", bToMb(m.Alloc))

		fmt.Sprintf(stats, fmt.Sprintf("\tTotalAlloc = %v MiB", bToMb(m.TotalAlloc)))
		// fmt.Sprintf(stats, "\tSys = %v MiB", bToMb(m.Sys))
		// fmt.Sprintf(stats, "\tNumGC = %v\n", m.NumGC)
		w.Write([]byte(fmt.Sprint(
			"Barrenschat API: OK", os.Getenv("NAME"), "\n",
			"Total Allocated:", bToMb(m.Alloc), "M\n",
			"Total Sys:", bToMb(m.Sys), "M\n",
			"Total Allocations:", bToMb(m.TotalAlloc), "\n",
			"Live Objects:", m.Mallocs-m.Frees,
		)))
	})

	mux.Handle("/", wsStart(h, authFunc))
	return mux
}

func bToMb(b uint64) uint64 {
	return b / 1024 / 1024
}

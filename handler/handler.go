package handler

import (
	"crypto/rsa"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"reflect"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
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

func wsStart(h *hub.Hub) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// TODO: Check for duplicate connection
		var tok *jwt.Token
		var err error

		resp, err := http.Get("https://www.googleapis.com/robot/v1/metadata/x509/securetoken@system.gserviceaccount.com")

		if err != nil {
			log.Println(err)
		}
		defer resp.Body.Close()

		respBody, err := ioutil.ReadAll(resp.Body)

		var publicPEM map[string]string
		err = json.Unmarshal(respBody, &publicPEM)
		if err != nil {
			log.Println(err.Error())
		}
		for _, v := range publicPEM {
			tok, err = jwt.Parse(r.URL.Query().Get("params"), func(token *jwt.Token) (interface{}, error) {
				publicKey, err = jwt.ParseRSAPublicKeyFromPEM([]byte(v))
				return publicKey, err
			})
			if err == nil {
				break
			}
		}

		if err != nil {
			log.Println(err.Error())
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		claims := tok.Claims.(jwt.MapClaims)

		// TODO: validate claims
		log.Println(reflect.TypeOf(claims))

		ws, err := upgrader.Upgrade(w, r, nil)

		if err != nil {
			log.Println(err.Error())
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		ws.SetReadLimit(1024 * 1024)
		ws.SetPongHandler(func(string) error {
			ws.SetWriteDeadline(time.Now().Add(connTimeout))
			ws.SetReadDeadline(time.Now().Add(connTimeout))
			log.Println("Pong rec")
			return nil
		})

		h.NewConnection <- struct {
			Ws     *websocket.Conn
			Claims jwt.MapClaims
		}{
			ws,
			claims,
		}
	}
}

// GetEngine returns router for the API
func GetEngine(h *hub.Hub) *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("/version", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(fmt.Sprint("Barrenschat API OK:", os.Getenv("NAME"))))
	})
	mux.Handle("/", wsStart(h))
	return mux
}

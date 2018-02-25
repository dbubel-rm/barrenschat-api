package main

import (
	"math/rand"
	"net/http"
	"os"

	"github.com/op/go-logging"
)

var log = logging.MustGetLogger("barrenschat-api")
var format = logging.MustStringFormatter(`%{color}%{time:15:04:05} %{shortfile} %{level:.4s} %{id:03x}%{color:reset} %{message}`)

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func RandStringBytes(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}
func sayHello(w http.ResponseWriter, r *http.Request) {
	log.Debug(os.Getenv("NAME"))
	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte(os.Getenv("NAME") + " " + RandStringBytes(1024)))
}

func main() {
	backend1 := logging.NewLogBackend(os.Stdout, "", 0)
	backend2Formatter := logging.NewBackendFormatter(backend1, format)
	logging.SetBackend(backend2Formatter)
	//logging.SetLevel(logging.INFO, "barrenschat-api")

	http.HandleFunc("/", sayHello)
	log.Info("Listening ..." + os.Getenv("NAME") + os.Getenv("PORT"))
	if err := http.ListenAndServe(":9000", nil); err != nil {
		panic(err)
	}
}

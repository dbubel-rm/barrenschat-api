package main

import (
	"fmt"
	"os"

	"github.com/engineerbeard/barrenschat-api/handler"
	"github.com/engineerbeard/barrenschat-api/hub"
	"github.com/op/go-logging"
)

func main() {
	var log = logging.MustGetLogger("barrenschat-api")
	var format = logging.MustStringFormatter(`%{color}%{time:15:04:05} %{shortfile} %{level:.4s} %{id:03x}%{color:reset} %{message}`)
	logBackend := logging.NewLogBackend(os.Stdout, "", 0)
	backendFormatter := logging.NewBackendFormatter(logBackend, format)
	logging.SetBackend(backendFormatter)
	logging.SetLevel(logging.DEBUG, "barrenschat-api")
	h := hub.NewHub()
	go h.Run()

	ginEngine := handler.GetEngine(h)
	log.Info("Server running")
	ginErr := ginEngine.Run(fmt.Sprintf(":9000"))
	log.Fatal(ginErr)

	// http.HandleFunc("/healthcheck", handlers.HealthCheck)
	// Log.Info("Listening ..." + os.Getenv("NAME") + os.Getenv("PORT"))
	// if err := http.ListenAndServe(":9000", nil); err != nil {
	// 	panic(err)
	// }
}

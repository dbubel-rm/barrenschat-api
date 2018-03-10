package main

import (
	"log"
	"net/http"

	"github.com/engineerbeard/barrenschat-api/handler"
	"github.com/engineerbeard/barrenschat-api/hub"
)

func main() {

	hubHandle := hub.NewHub()
	go hubHandle.Run()

	serverMux := handler.GetEngine(hubHandle)
	log.Println("Server running")
	http.ListenAndServe(":9000", serverMux)
}

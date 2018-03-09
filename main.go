package main

import (
	"fmt"
	"log"

	"github.com/engineerbeard/barrenschat-api/handler"
	"github.com/engineerbeard/barrenschat-api/hub"
)

func main() {

	h := hub.NewHub()
	go h.Run()

	ginEngine := handler.GetEngine(h)
	log.Println("Server running")
	ginErr := ginEngine.Run(fmt.Sprintf(":9000"))
	log.Fatal(ginErr)

	// http.HandleFunc("/healthcheck", handlers.HealthCheck)
	// Log.Info("Listening ..." + os.Getenv("NAME") + os.Getenv("PORT"))
	// if err := http.ListenAndServe(":9000", nil); err != nil {
	// 	panic(err)
	// }
}

package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"runtime"
	"time"

	"github.com/engineerbeard/barrenschat-api/handler"
	"github.com/engineerbeard/barrenschat-api/hub"
)

func PrintMemUsage() {
	for {
		var m runtime.MemStats
		runtime.ReadMemStats(&m)
		// For info on each, see: https://golang.org/pkg/runtime/#MemStats
		fmt.Printf("Alloc = %v MiB", bToMb(m.Alloc))
		fmt.Printf("\tTotalAlloc = %v MiB", bToMb(m.TotalAlloc))
		fmt.Printf("\tSys = %v MiB", bToMb(m.Sys))
		fmt.Printf("\tNumGC = %v\n", m.NumGC)
		time.Sleep(time.Second * 5)
	}

}

func bToMb(b uint64) uint64 {
	return b / 1024 / 1024
}

// TODO: benchcmp
func main() {
	f, err := os.OpenFile("hub_log.txt", os.O_APPEND|os.O_CREATE|os.O_RDWR, 0666)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	mw := io.MultiWriter(os.Stdout, f)
	log.SetOutput(mw)
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)

	hubHandle := hub.NewHub()
	hubHandle.HandleMsg("message_new", hubHandle.HandleClientMessage)
	go hubHandle.Run()

	serverMux := handler.GetEngine(hubHandle)
	log.Println("Server running port 9000")
	//go PrintMemUsage()
	log.Fatalln(http.ListenAndServe(":9000", serverMux))
}

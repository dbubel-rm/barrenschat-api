package main

import (
	"fmt"
	"time"

	"github.com/gorilla/websocket"
)

func main() {

	mockWebsocket, _, err := websocket.DefaultDialer.Dial("ws://localhost:9000/bchatws?params=eyJhbGciOiJSUzI1NiIsImtpZCI6IjMxNjAwMjk1MjI3ODA5M2RmODA3YzkxMGNjYTBmODE3YmI4ODcxY2YifQ.eyJpc3MiOiJodHRwczovL3NlY3VyZXRva2VuLmdvb2dsZS5jb20vYmFycmVuc2NoYXQtMjcyMTIiLCJhdWQiOiJiYXJyZW5zY2hhdC0yNzIxMiIsImF1dGhfdGltZSI6MTUzNTQwMzA2NiwidXNlcl9pZCI6IjBLbGtRR1lBcjdRR2tRQ0JYZzYwZUw4U2lzcTEiLCJzdWIiOiIwS2xrUUdZQXI3UUdrUUNCWGc2MGVMOFNpc3ExIiwiaWF0IjoxNTM1NDAzMDY2LCJleHAiOjE1MzU0MDY2NjYsImVtYWlsIjoiZGVhbkB0ZXN0LmNvbSIsImVtYWlsX3ZlcmlmaWVkIjpmYWxzZSwiZmlyZWJhc2UiOnsiaWRlbnRpdGllcyI6eyJlbWFpbCI6WyJkZWFuQHRlc3QuY29tIl19LCJzaWduX2luX3Byb3ZpZGVyIjoicGFzc3dvcmQifX0.zfhY6O54mOJdU42t3_6t6DN3hmLdYjT80zYM3zSjo9p9nuMDNrkCpHpzly_p5ytEJTXJEAeVEE0au_8SZ1E9GHNyz8zUsmODx7ijUJJBvWsP9n_BrnBjP1J0NLcvcLhWPwtFO48FD1Hk15l1iaflzKYtSssTmPnStO5JlI7qkvhK5jh4XhZlg2LW9KT0mYW8BiyQIK1PH4B-4wKCpU0Q4OdxnqForioFXtFtnv_k7bivt9VOvgQjldPP9MIAEcaig50b2M103NPoq_j6ipx1dwSGGms2Dy3Im0jetHHAbsUQJmr6XRKZ3Kiyso1DSbxlU1zDHvvEedOPri0SGTwljg", nil)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	mockWebsocket.SetReadDeadline(time.Now().Add(time.Second * time.Duration(3)))
	mockWebsocket.SetWriteDeadline(time.Now().Add(time.Second * time.Duration(3)))
	time.Sleep(time.Second * 30)
	mockWebsocket.Close()
}

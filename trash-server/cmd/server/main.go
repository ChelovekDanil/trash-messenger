package main

import (
	"fmt"
	"net/http"

	"github.com/ChelovekDanil/trash-messenger/internal/server"
)

func main() {
	http.HandleFunc("/ws", server.WsHandler)
	fmt.Println("server started")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		panic(err)
	}
}

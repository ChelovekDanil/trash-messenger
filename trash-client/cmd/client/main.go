package main

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/gorilla/websocket"
)

type message struct {
	Username string `json:"username"`
	Message  string `json:"message"`
}

const (
	url             = "ws://localhost:8080/ws"
	textMessageType = 1
)

func main() {
	conn, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	ctx, cancel := context.WithCancel(context.Background())

	go recipient(cancel, conn)
	go sender(cancel, conn)
	fmt.Println("send message like: username@text message$")
	<-ctx.Done()
	log.Println("connection close(")
}

func recipient(cancel context.CancelFunc, conn *websocket.Conn) {
	for {
		_, messageData, err := conn.ReadMessage()
		if err != nil {
			log.Println(err)
			cancel()
			return
		}
		var req message
		if err := json.Unmarshal(messageData, &req); err != nil {
			log.Println(err)
			return
		}
		fmt.Printf("%s: %s\n", req.Username, req.Message)
	}
}

func sender(cancel context.CancelFunc, conn *websocket.Conn) {
	for {
		text, err := bufio.NewReader(os.Stdin).ReadString('$')
		if err != nil {
			log.Println(err)
			return
		}

		messageData := strings.Split(text, "@")
		messageJson, err := json.Marshal(message{Username: messageData[0], Message: messageData[1]})
		if err != nil {
			log.Println(err)
			return
		}

		if err := conn.WriteMessage(textMessageType, messageJson); err != nil {
			log.Println(err)
			cancel()
			return
		}
	}
}

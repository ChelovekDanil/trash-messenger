package server

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/websocket"
)

type databusReq struct {
	messageType int
	Message
}

type Message struct {
	Username string `json:"username"`
	Message  string `json:"message"`
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

var userinmem = make(map[string]*websocket.Conn)
var count = 1

func WsHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	userinmem[strconv.Itoa(count)] = conn
	count++
	databus := make(chan databusReq)
	ctx, cancel := context.WithCancel(context.Background())

	go recipient(cancel, conn, databus)
	go sender(cancel, databus)
	fmt.Println(userinmem)
	<-ctx.Done()
	log.Println("connection close (0_0)")
}

// recipient - принимает сообщение
func recipient(cancel context.CancelFunc, conn *websocket.Conn, databus chan<- databusReq) {
	for {
		messageType, messageData, err := conn.ReadMessage()
		if err != nil {
			log.Println(err)
			cancel()
			return
		}
		var mess Message
		if err = json.Unmarshal(messageData, &mess); err != nil {
			log.Println(err)
			return
		}
		fmt.Println(mess)
		databus <- databusReq{messageType, mess}
	}
}

// sender - отправляет сообщение получателю
func sender(cancel context.CancelFunc, databus <-chan databusReq) {
	for {
		req := <-databus
		reqJson, err := json.Marshal(Message{req.Username, string(req.Message.Message)})
		if err != nil {
			log.Println(err)
			cancel()
			return
		}
		conn := userinmem[req.Username]
		if err := conn.WriteMessage(req.messageType, reqJson); err != nil {
			log.Println(err)
			cancel()
			return
		}
	}
}

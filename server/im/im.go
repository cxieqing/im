package im

import (
	"fmt"
	"im/pkg"
	"im/pkg/config"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

const Version = "1.0.0"

var mu sync.Mutex

var upgrader = websocket.Upgrader{
	ReadBufferSize:    4096,
	WriteBufferSize:   4096,
	EnableCompression: true,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func CreateImServer() {
	config := config.NewConfig()
	addr := fmt.Sprintf("%s:%d", config.ImHost, config.ImPort)
	http.HandleFunc("/", acceptRequest)
	http.ListenAndServe(addr, nil)
}

func acceptRequest(w http.ResponseWriter, r *http.Request) {
	//conn.SetReadDeadline(time.Now().Add(s.ReadDeadline))
	server := pkg.NewImServer()
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		//@todo log
		return
	}
	server.WaitConnetClinet <- conn
}

func handleClinet() {
	s := pkg.NewImServer()
	for {
		select {
		case conn := <-s.WaitConnetClinet:
			go func() {

			}()
		}
	}

	for {
		select {
		case token := <-s.NewClient:
			client := s.Clients[token]
			go client.HandReadMessage(s)
			go client.HandWriteMessage(s)
		case token := <-s.WaitCloseClient:
			client := s.Clients[token]
			client.Free()
			delete(s.Clients, token)
		}
	}
}

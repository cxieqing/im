package im

import (
	"fmt"
	"im/pkg"
	"im/pkg/config"
	"net/http"
	"sync"
	"time"

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
	go handleClinet()

	config := config.NewConfig()
	addr := fmt.Sprintf("%s:%d", config.ImHost, config.ImPort)
	s := &http.Server{
		Addr:           addr,
		Handler:        http.HandlerFunc(acceptRequest),
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	s.ListenAndServe()
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
			go pkg.WaitClientReady(conn, s)
		case client := <-s.NewClient:
			s.Clients[client.Token] = client
			client.HandMessage(s)
		case token := <-s.WaitCloseClient:
			client := s.Clients[token]
			client.Free()
			delete(s.Clients, token)
		}
	}
}

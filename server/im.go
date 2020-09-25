package server

import (
	"im/pkg"
	"im/pkg/config"
	"im/pkg/message"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

var imInstance ImService
var once sync.Once
var mu sync.Mutex

type ImService struct {
	Clinets       map[string]*pkg.Clinet
	MaxClinetNum  int
	ReadDeadline  int64
	WriteDeadline int64
}

type ImMessage struct {
	Token       string              `json:"token"`
	Content     string              `json:"content"`
	ContentType message.ContenType  `json:"contentType"`
	To          int                 `json:"to"`
	MessageType message.MessageType `json:"messageType"`
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:    4096,
	WriteBufferSize:   4096,
	EnableCompression: true,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func getInstance(c *Config) *ImMessage {
	once.Do(func() {
		imInstance = new(ImMessage)
		imInstance.MaxClinetNum = c.MaxClinetNum
		imInstance.ReadDeadline = c.ReadDeadline
		imInstance.WriteDeadline = c.WriteDeadline
	})
	return imInstance
}

func wsService(w http.ResponseWriter, r *http.Request) {
	config := config.NewConfig()
	im := getInstance(config)
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		//@todo log
		return
	}
	go addClinet(conn, im)
}

func addClinet(conn websocket.Conn, s ImService) {
	conn.SetReadDeadline(time.Now().Add(s.ReadDeadline))
	var linkMessage = new(ImMessage)
	err := conn.ReadJSON(linkMessage)
	if err != nil {
		conn.Close()
		return
	}
	mu.lock()
	s.Clinets[linkMessage.Token] = 
}

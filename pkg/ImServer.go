package pkg

import (
	"im/pkg/config"
	"sync"

	"github.com/gorilla/websocket"
)

var once sync.Once

var imInstance *ImServer

type ImServer struct {
	Clients          map[string]*Client
	MaxClientNum     int
	ReadDeadline     int
	WriteDeadline    int
	NewClient        chan string
	WaitCloseClient  chan string
	WaitConnetClinet chan *websocket.Conn
	//Log           log
}

func NewImServer() *ImServer {
	once.Do(func() {
		imInstance := new(ImServer)
		config := config.NewConfig()
		imInstance.Clients = make(map[string]*Client, 0)
		imInstance.MaxClientNum = config.MaxClientNum
		imInstance.ReadDeadline = config.ReadDeadline
		imInstance.WriteDeadline = config.WriteDeadline
		imInstance.NewClient = make(chan string)
		imInstance.WaitCloseClient = make(chan string)
	})
	return imInstance
}

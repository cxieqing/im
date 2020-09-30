package server

import (
	"fmt"
	"im/pkg"
	"im/pkg/config"
	"im/pkg/message"
	"im/pkg/redis"
	"im/pkg/user"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

var imInstance ImService
var once sync.Once
var mu sync.Mutex

type ImService struct {
	Clinets       map[string]*pkg.Clinet
	ClinetNum     int
	ReadDeadline  int64
	WriteDeadline int64
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:    4096,
	WriteBufferSize:   4096,
	EnableCompression: true,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func GetImInstance(c *config.Config) *ImService {
	once.Do(func() {
		imInstance := new(ImService)
		imInstance.ReadDeadline = c.ReadDeadline
		imInstance.WriteDeadline = c.WriteDeadline
	})
	return &imInstance
}

func wsService(w http.ResponseWriter, r *http.Request) {
	config := config.NewConfig()
	im := GetImInstance(config)
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		//@todo log
		return
	}
	go addClinet(conn, im)
}

func addClinet(conn *websocket.Conn, s *ImService) {
	//conn.SetReadDeadline(time.Now().Add(s.ReadDeadline))
	var linkMessage = new(pkg.ImMessage)
	err := conn.ReadJSON(linkMessage)
	if err != nil {
		conn.Close()
		return
	}
	clinetPtr := pkg.NewClinet(conn, linkMessage.Token)
	clinetPtr.Init()
	mu.Lock()
	s.Clinets[linkMessage.Token] = clinetPtr
	mu.Unlock()
}

func getOnlineUserById(id int) *pkg.Clinet {
	cacheKey := "userid:token:" + fmt.Sprint(id)
	redis := redis.NewRedis()
	token, err := redis.Client.Get(cacheKey).Result()
	if err != nil {
		return nil
	}
	clinetptr, ok := imInstance.Clinets[token]
	if !ok {
		return nil
	}
	return clinetptr
}

func ReadMessage(c *pkg.Clinet) {
	for {
		var msg = pkg.ImMessage{}
		err := c.Connet.ReadJSON(msg)
		if err != nil {
			//@todo log
		} else {
			if msg.MessageType == message.GroupMessage {
				group, ok := c.Group[msg.To]
				if !ok {
					continue
				}
				GroupMessageSend(group, msg)
			} else {
				SendUserMessage(msg.To, msg)
			}

		}
	}
}

func GroupMessageSend(g *user.Group, msg pkg.ImMessage) {
	for _, memberId := range g.Members {
		SendUserMessage(memberId, msg)
	}
}

func SendUserMessage(toUserId int, msg pkg.ImMessage) {
	clinet := getOnlineUserById(toUserId)
	if clinet == nil {
		if msg.MessageType == message.UserMessage {
			message.SaveUnsendUserMsg(pkg.ImTransTomsg(msg))
		} else {
			message.SaveUnsendGroupMsg(pkg.ImTransTomsg(msg))
		}
	}
	clinet.MessageChn <- msg
}

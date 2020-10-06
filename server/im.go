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
	"time"

	"github.com/gorilla/websocket"
)

const Version = "1.0.0"

var imInstance ImService
var once sync.Once
var mu sync.Mutex
var pingPeriod = 5 * time.Second

type ImService struct {
	Clinets       map[string]*pkg.Clinet
	ClinetNum     int
	ReadDeadline  int
	WriteDeadline int
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

func ImServer(w http.ResponseWriter, r *http.Request) {
	config := config.NewConfig()
	im := GetImInstance(config)
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		//@todo log
		return
	}
	clinetPtr := addClinet(conn, im)
	if clinetPtr != nil {
		go ReadMessage(clinetPtr)
		go WriteMessage(clinetPtr)
	}
}

func addClinet(conn *websocket.Conn, s *ImService) *pkg.Clinet {
	//conn.SetReadDeadline(time.Now().Add(s.ReadDeadline))
	var linkMessage = new(pkg.ImMessage)
	err := conn.ReadJSON(linkMessage)
	if err != nil {
		conn.Close()
		return nil
	}
	clinetPtr := pkg.NewClinet(conn, linkMessage.Token)
	clinetPtr.Init()
	mu.Lock()
	s.Clinets[linkMessage.Token] = clinetPtr
	mu.Unlock()
	return clinetPtr
}

func getOnlineUserById(id uint) *pkg.Clinet {
	token, err := getUserTokenById(id)
	if err != nil {
		return nil
	}
	clinetptr, ok := imInstance.Clinets[token]
	if !ok {
		return nil
	}
	return clinetptr
}

func getUserTokenById(id uint) (string, error) {
	cacheKey := "userid:token:" + fmt.Sprint(id)
	redis := redis.NewRedis()
	return redis.Client.Get(cacheKey).Result()
}

func ReadMessage(c *pkg.Clinet) {
	for {
		msg, err := c.ReadMessage()
		if err != nil {
			//@todo log
		} else {
			if msg.MessageType == message.Group {
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

func WriteMessage(c *pkg.Clinet) {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.Connet.Close()
		token, err := getUserTokenById(c.UserInfo.Id)
		if err != nil {
			delete(imInstance.Clinets, token)
		}
	}()
	writeDeadline := time.Duration(imInstance.WriteDeadline)
	for {
		select {
		case msg, ok := <-c.MessageChn:
			if !ok {
				c.Connet.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			c.Connet.SetWriteDeadline(time.Now().Add())
			c.Connet.WriteJSON(msg)
		case <-ticker.C:
			err := c.Connet.WriteControl(websocket.PingMessage, []byte("ping"), time.Now().Add(writeDeadline))
			if err != nil {
				c.Connet.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
		}
	}
}

func GroupMessageSend(g *user.Group, msg pkg.ImMessage) {
	for _, memberId := range g.Members {
		SendUserMessage(memberId, msg)
	}
}

func SendUserMessage(toUserId uint, msg pkg.ImMessage) {
	clinet := getOnlineUserById(toUserId)
	if clinet == nil {
		if msg.MessageType == message.User {
			message.SaveUserMsg(pkg.ImTransTomsg(msg), message.UnSend)
		} else {
			message.SaveUnsendGroupMsg(pkg.ImTransTomsg(msg))
		}
		return
	}
	clinet.WriteMessage(msg)
}
